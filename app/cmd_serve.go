package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labctl/labctl/helpers"
	"github.com/labctl/labctl/helpers/frontend"
	"github.com/labctl/labctl/utils"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

type CmdServe struct {
	Topo          string   `short:"t" help:"Topology file" type:"existingfile"`
	TemplatePaths []string `short:"p" help:"Paths to search for templates" type:"path"`

	Addr string `help:"Serve on addr." default:":8080"`
}

func (r *CmdServe) Run(ctx *helpers.Context) error {
	err := ctx.Load(r.Topo)
	if err != nil {
		return err
	}
	ctx.TemplatePaths, err = helpers.InitTemplatePaths(r.TemplatePaths)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.Handle("/favicon.ico", http.RedirectHandler("/labctl/favicon.ico", http.StatusSeeOther))
	mux.Handle("/", http.RedirectHandler("/labctl", http.StatusSeeOther))
	mux.HandleFunc("/labctl/ws", websock)
	mux.HandleFunc("/labctl/topo", http_topo)
	mux.HandleFunc("/labctl/vars", http_vars)
	mux.HandleFunc("/labctl/templates", http_templates)
	mux.HandleFunc("/error", http_error)

	frontendServer := frontend.LabctlFileServer()
	mux.Handle("/labctl", frontendServer)
	mux.Handle("/labctl/", frontendServer)

	handler := cors.Default().Handler(&frontend.SlashFix{Mux: mux})

	log.Infof("Serve on %s", r.Addr)
	return http.ListenAndServe(r.Addr, handler)
}

var wsupgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func websock(w http.ResponseWriter, r *http.Request) {
	c, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	// immediately send a UI update
	helpers.WsSendUiUpdate(c, Ctx)

	var wsmsg helpers.WsMessage

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		err = wsmsg.UnmarshalJson(message)
		if err != nil {
			log.Warnf("Ws Rx: %s: %v", err.Error(), message)
			continue
		}
		switch wsmsg.Code {
		case helpers.WscHeartbeat:
		case helpers.WscUiData: // save UI settings
			wsmsg.UiData.WriteFile(Ctx)
			Ctx.Template, err = utils.ParseTemplates(wsmsg.UiData.Templates)
			if err != nil {
				helpers.WsErrorf(c, "cannot parse template: %s", err)
			}

		case helpers.WscTemplate: // render template
			err := wsmsg.Template.Render(Ctx)
			if err != nil {
				helpers.WsWarnf(c, err.Error())
				continue
			}
			wsmsg.Template.ClearInput()
			err = c.WriteJSON(wsmsg)
			if err != nil {
				log.Errorf("could not write to the websocket: %s", err)
			}

		case helpers.WscConfig: // run the config command
			err = RunWebConfig(c, wsmsg.Config.Cmd)
			if err != nil {
				helpers.WsErrorf(c, err.Error())
			}

		default:
			log.Warnf("unhandled websocket message: %v", wsmsg)
		}
	}
}

type jsonResponse struct {
	Ok      bool        `json:"ok"`
	Message string      `json:"msg,omitempty"`
	Data    interface{} `json:"data"`
}

func json_response(w http.ResponseWriter, j interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(jsonResponse{
		Ok:   true,
		Data: j,
	})
	if err != nil {
		log.Errorf("could not encode json: %s", err)
	}
}

func json_message(w http.ResponseWriter, message string, ok ...bool) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(jsonResponse{
		Ok:      len(ok) > 0 && ok[0],
		Message: message,
		Data:    struct{}{},
	})
	if err != nil {
		log.Errorf("could not encode json: %s", err)
	}
}

func http_topo(w http.ResponseWriter, req *http.Request) {
	j, err := Ctx.Topo.AsJson()
	if err != nil {
		log.Error(err)
		json_message(w, err.Error())
		return
	}
	json_response(w, j)
}

func http_vars(w http.ResponseWriter, req *http.Request) {
	j, err := Ctx.Topo.VarsAsJson()
	if err != nil {
		log.Error(err)
		json_message(w, err.Error())
		return
	}
	json_response(w, j)
}

func http_templates(w http.ResponseWriter, req *http.Request) {
	t, err := helpers.LoadTemplates(Ctx)
	if err != nil {
		log.Error(err)
		json_message(w, err.Error())
		return
	}
	json_response(w, t)
}

func http_error(w http.ResponseWriter, req *http.Request) {
	json_message(w, "error!")
}
