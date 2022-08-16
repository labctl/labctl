package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type WsMsgCode string

const (
	WscHeartbeat WsMsgCode = "."
	WscError     WsMsgCode = "error"
	WscTemplate  WsMsgCode = "template"
	WscUiData    WsMsgCode = "uidata"
	WscConfig    WsMsgCode = "config"
)

type WsMessage struct {
	Code     WsMsgCode   `json:"code"`
	Error    string      `json:"error,omitempty"`
	UiData   *WsUiData   `json:"uidata,omitempty"`
	Template *WsTemplate `json:"template,omitempty"`
	Config   *WsConfig   `json:"config,omitempty"`
}

func (w *WsMessage) Unmarshal(data []byte) error {
	// Clear the websocket data, ready to populate the newly received string
	w.Code = ""
	// for k := range w.UiData.Layouts.Nodes {
	// 	delete(w.UiData.Layouts.Nodes, k)
	// }

	err := json.Unmarshal(data, &w)
	if err != nil {
		log.Errorf("%s: %v", err.Error(), data)
		return err
	}
	return nil
}

// Implement stringer
func (w *WsMessage) String() string {
	res, err := json.Marshal(&w)
	if err != nil {
		log.Error(err)
		return ""
	}
	return string(res)
}

func WsErrorf(conn *websocket.Conn, msg string, args ...interface{}) {
	m := fmt.Sprintf(msg, args...)
	wsmsg := &WsMessage{
		Code:  WscError,
		Error: m,
	}
	WsWriteJSON(conn, wsmsg)
}
