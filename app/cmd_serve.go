package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labctl/labctl/helpers"
	"github.com/labctl/labctl/helpers/frontend"
	"github.com/labctl/labctl/utils"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

type CmdServe struct {
	Topo          string   `short:"t" help:"Topology file" type:"existingfile" predictor:"topo"`
	TemplatePaths []string `short:"p" help:"Paths to search for templates" type:"path" predictor:"dir"`

	Addr string `help:"Serve on addr." default:":8080"`
	Url  string `help:"Serve on a custom URL path" default:"/labctl"`
}

func (r *CmdServe) Run(ctx *helpers.Context) error {
	var err error
	r.Topo, err = utils.EnsureTopo(r.Topo)
	if err != nil {
		return err
	}

	err = ctx.Load(r.Topo)
	if err != nil {
		return err
	}
	ctx.TemplatePaths, err = helpers.InitTemplatePaths(r.TemplatePaths)
	if err != nil {
		return err
	}

	// Create new FS watcher
	watcher := helpers.WatchFS(ctx)
	defer watcher.Close()

	// Start the web server
	mux := http.NewServeMux()

	mux.Handle("/favicon.ico", http.RedirectHandler(r.Url+"/favicon.ico", http.StatusSeeOther))
	mux.Handle("/", http.RedirectHandler(r.Url, http.StatusSeeOther))
	mux.HandleFunc(r.Url+"/ws", websock)
	mux.HandleFunc(r.Url+"/topo", http_topo)
	mux.HandleFunc(r.Url+"/vars", http_vars)
	mux.HandleFunc(r.Url+"/templates", http_templates)
	mux.HandleFunc("/error", http_error)

	frontendServer := frontend.LabctlFileServer(r.Url)
	mux.Handle(r.Url, frontendServer)
	mux.Handle(r.Url+"/", frontendServer)

	handler := cors.Default().Handler(&frontend.SlashFix{Mux: mux})

	// Check and print the web URL
	host, _, port := utils.Partition(r.Addr, ":")
	if port == "" {
		return fmt.Errorf("you need to specify a port")
	}
	url := "http://" + r.Addr
	if host == "" {
		url = fmt.Sprintf("http://localhost%s or http://%s%s", r.Addr, utils.GetOutboundIP(), r.Addr)
	}
	logLatestVersion(5)
	log.Infof("Access the web server on %s", url)
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
	helpers.WsLogf(c, helpers.WscWarn, "websocket connected. version %s", version)

	helpers.WsSendUiUpdate(c, Ctx)

	var wsmsg helpers.WsMessage

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		if bytes.Equal(message, []byte("ping")) {
			_ = c.WriteMessage(0, message)
			continue
		}

		err = wsmsg.UnmarshalJson(message)
		if err != nil {
			log.Warnf("Ws Rx: %s: %v", err.Error(), message)
			continue
		}
		switch wsmsg.Code {
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
			WebConfigDone(c, 3)

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
