package helpers

import (
	"github.com/gorilla/websocket"
	"github.com/labctl/labctl/utils/colorize"
	"github.com/labctl/labctl/utils/tx"
	"github.com/srl-labs/containerlab/clab/config"

	log "github.com/sirupsen/logrus"
)

type OutputAsLog interface{}

type OutputAsJson interface{}

type ResultOutput interface {
	LogResponses([]*tx.Response, *config.NodeConfig)
	Info(string, string)
	Error(string, string)
	PreferStdout() bool
}

// The standard output to the log. Implements ResultOutput
type LogOutput struct{}

func (l *LogOutput) Info(node string, msg string) {
	log.Infof("%s: %s", node, msg)
}

func (l *LogOutput) Error(node string, msg string) {
	log.Errorf("%s: %s", node, msg)
}

func (l *LogOutput) LogResponses(obj []*tx.Response, _ *config.NodeConfig) {
	for _, r := range obj {
		r.Log(r.Level)
	}
}

func (l *LogOutput) PreferStdout() bool {
	return true
}

// Implements ResultOutput
type WebSocketOutput struct {
	Conn *websocket.Conn
}

func (ws *WebSocketOutput) Info(node string, msg string) {
	wsmsg := &WsMessage{
		Code: WscConfig,
		Config: &WsConfig{
			Results: []*tx.Response{{Node: node, Response: msg, Level: log.DebugLevel}},
		},
	}
	WsWriteJSON(ws.Conn, wsmsg)
}

func (ws *WebSocketOutput) Error(node string, msg string) {
	wsmsg := &WsMessage{
		Code: WscConfig,
		Config: &WsConfig{
			Results: []*tx.Response{{Node: node, Response: msg, Level: log.ErrorLevel}},
		},
	}
	WsWriteJSON(ws.Conn, wsmsg)
}

func (ws *WebSocketOutput) LogResponses(obj []*tx.Response, nc *config.NodeConfig) {
	wsmsg := &WsMessage{
		Code: WscConfig,
		Config: &WsConfig{
			Results: obj,
			Input:   nc,
		},
	}
	WsWriteJSON(ws.Conn, wsmsg)
}

func (l *WebSocketOutput) PreferStdout() bool {
	return false
}

// Color the templates
func ColorResults(responses []*tx.Response, regexes []*colorize.Colorize) {
	for i, r := range responses {
		// Replace for every regex
		for _, rex := range regexes {
			responses[i].Response = rex.Regexp.ReplaceAllString(r.Response, string(rex.Replace))
		}
	}
}
