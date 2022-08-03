package app

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labctl/labctl/helpers"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

type CmdServe struct {
	Topo          string   `short:"t" help:"Topology file" type:"existingfile"`
	TemplatePaths []string `short:"p" help:"Paths to search for templates" type:"existingdir"`

	Addr string `help:"Serve on addr." default:":8080"`
}

var upgrader = websocket.Upgrader{} // use default options

var Ctx *helpers.Context

func (r *CmdServe) Run(ctx *helpers.Context) error {
	err := ctx.Load(r.Topo)
	if err != nil {
		return err
	}
	Ctx = ctx

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", websock)
	mux.HandleFunc("/topo", http_topo)
	mux.HandleFunc("/vars", http_vars)
	mux.HandleFunc("/", home)

	handler := cors.Default().Handler(mux)

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

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
		case 1: // echo
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
		case 200: // hearbeat
			log.Debugf(wsdata.Msg)
		case 300: // render template
			err := wsdata.Template.Render(Ctx)
			if err != nil {
				log.Errorf("%s", err)
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

func json_response(w http.ResponseWriter, j interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(j)
	if err != nil {
		log.Errorf("could not encode json: %s", err)
	}
}

func http_topo(w http.ResponseWriter, req *http.Request) {
	j, err := Ctx.Topo.AsJson()
	if err != nil {
		log.Error(err)
	}
	json_response(w, j)
}

func http_vars(w http.ResponseWriter, req *http.Request) {
	if (*req).Method == "OPTIONS" {
		return
	}

	j, err := Ctx.Topo.VarsAsJson()
	if err != nil {
		log.Error(err)
	}
	json_response(w, j)
}

func home(w http.ResponseWriter, r *http.Request) {
	err := homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
	if err != nil {
		log.Error(err)
	}
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
</head>
<body>
No implemented
</body>
</html>
`))
