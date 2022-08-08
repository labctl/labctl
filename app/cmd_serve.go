package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labctl/labctl/helpers"
	"github.com/labctl/labctl/helpers/frontend"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

type CmdServe struct {
	Topo          string   `short:"t" help:"Topology file" type:"existingfile"`
	TemplatePaths []string `short:"p" help:"Paths to search for templates" type:"existingdir"`

	Addr string `help:"Serve on addr." default:":8080"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
} // use default options

var Ctx *helpers.Context

func (r *CmdServe) Run(ctx *helpers.Context) error {
	err := ctx.Load(r.Topo)
	if err != nil {
		return err
	}
	ctx.TemplatePaths, err = helpers.InitTemplatePaths(r.TemplatePaths)
	if err != nil {
		return err
	}
	Ctx = ctx

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

func websock(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	wsdata := helpers.NewWebSocketData()
	err = wsdata.Data.ReadFile(Ctx)
	if err != nil {
		log.Debugf("Could not read lab file %s", err.Error())
	} else {
		wsdata.Code = 100
		err = c.WriteJSON(wsdata)
		if err != nil {
			log.Errorf("could not write to the websocket: %s", err)
		}
	}

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		err = wsdata.Unmarshal(message)
		if err != nil {
			continue
		}
		switch wsdata.Code {
		case 1: // hearbeat
			log.Debugf(wsdata.Msg)
		case 2: // echo
			err = c.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
			}
		case 100: // save settings
			wsdata.Data.WriteFile(Ctx)
			Ctx.Template, err = helpers.ParseTemplates(wsdata.Data.Templates)
			if err != nil {
				log.Errorf("cannot parse tmeplate: %s", err.Error())
			}
		case 300: // render template
			err := wsdata.Template.Render(Ctx)
			if err != nil {
				log.Errorf("%s", err)
			} else {
				// if successful, we can clear the template & vars
				wsdata.Template.ClearInput()
			}
			err = c.WriteJSON(wsdata)
			if err != nil {
				log.Errorf("could not write to the websocket: %s", err)
			}

		default:
			log.Infof("recv %v", wsdata)
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
