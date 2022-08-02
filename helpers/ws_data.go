package helpers

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type WebSocketData struct {
	Code     int               `json:"code,omitempty"`
	Msg      string            `json:"msg,omitempty"`
	Data     UiData            `json:"data,omitempty"`
	Template WebSocketTemplate `json:"template,omitempty"`
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
