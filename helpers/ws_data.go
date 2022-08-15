package helpers

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type WsCode int

const (
	WsRender    WsCode = 300
	WsSave      WsCode = 100
	WsEcho      WsCode = 1
	WsConfigCmd WsCode = 400
	WsConfigOk  WsCode = 401
	WsConfigErr WsCode = 402
)

type WebSocketData struct {
	Code WsCode `json:"code"`
	Msg  string `json:"msg,omitempty"`
	// Nodes are used by config template returns
	Node     string             `json:"node,omitempty"`
	Data     *UiData            `json:"data,omitempty"`
	Template *WebSocketTemplate `json:"template,omitempty"`
}

func NewWebSocketData() WebSocketData {
	return WebSocketData{
		Data: NewUiData(),
	}
}

func (w *WebSocketData) Unmarshal(data []byte) error {
	// Clear the websocket data, ready to populate the newly received string
	w.Code = 0
	w.Msg = ""
	for k := range w.Data.Layouts.Nodes {
		delete(w.Data.Layouts.Nodes, k)
	}

	err := json.Unmarshal(data, &w)
	if err != nil {
		log.Errorf("%s: %v", err.Error(), data)
		return err
	}
	return nil
}

func (w *WebSocketData) AsJson() []byte {
	res, err := json.Marshal(&w)
	if err != nil {
		log.Error(err)
		return []byte("")
	}
	return res
}
