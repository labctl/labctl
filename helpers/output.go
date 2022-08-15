package helpers

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type ResultOutput interface {
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

// Implements ResultOutput
type WebSocketOutput struct {
	Conn *websocket.Conn
}

func (ws *WebSocketOutput) Info(node string, msg string) {
	ws.Conn.WriteJSON(WebSocketData{
		Code: WsConfigOk,
		Node: node,
		Msg:  msg,
	})
}

func (ws *WebSocketOutput) Error(node string, msg string) {
	ws.Conn.WriteJSON(WebSocketData{
		Code: WsConfigErr,
		Node: node,
		Msg:  msg,
	})
}
