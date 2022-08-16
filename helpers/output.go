package helpers

import (
	"github.com/gorilla/websocket"
	"github.com/labctl/labctl/utils/tx"

	log "github.com/sirupsen/logrus"
)

type OutputAsLog interface{}

type OutputAsJson interface{}

type ResultOutput interface {
	Obj([]*tx.Response)
	Info(string, string)
	Error(string, string)
}

// The standard output to the log. Implements ResultOutput
type LogOutput struct{}

func (l *LogOutput) Info(node string, msg string) {
	log.Infof("%s: %s", node, msg)
}

func (l *LogOutput) Error(node string, msg string) {
	log.Errorf("%s: %s", node, msg)
}

func (l *LogOutput) Obj(obj []*tx.Response) {
	for _, r := range obj {
		switch r.Level {
		case -1:
			r.Log(log.DebugLevel)
		case 0:
			r.Log(log.InfoLevel)
		case 1:
			r.Log(log.WarnLevel)
		default:
			r.Log(log.ErrorLevel)
		}
	}
}

// Implements ResultOutput
type WebSocketOutput struct {
	Conn *websocket.Conn
}

func (ws *WebSocketOutput) Info(node string, msg string) {
	wsmsg := &WsMessage{
		Code: WscConfig,
		Config: &WsConfig{
			Node:    node,
			Sresult: msg,
		},
	}
	WsWriteJSON(ws.Conn, wsmsg)
}

func (ws *WebSocketOutput) Error(node string, msg string) {
	wsmsg := &WsMessage{
		Code: WscConfig,
		Config: &WsConfig{
			Node:    node,
			Sresult: msg,
		},
	}
	WsWriteJSON(ws.Conn, wsmsg)
}

func (ws *WebSocketOutput) Obj(obj []*tx.Response) {
	wsmsg := &WsMessage{
		Code: WscConfig,
		Config: &WsConfig{
			Result: obj,
		},
	}
	WsWriteJSON(ws.Conn, wsmsg)
}
