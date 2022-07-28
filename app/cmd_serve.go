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
	ctx.Load(r.Topo)
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
	log.Infof("a")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func json_response(w http.ResponseWriter, j interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(j)
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
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server,
"Send" to send a message to the server and "Close" to close the connection.
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
