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
	WscWarn      WsMsgCode = "warn"
	WscTemplate  WsMsgCode = "template"
	WscUiData    WsMsgCode = "uidata"
	WscConfig    WsMsgCode = "config"
)

type WsMessage struct {
	Code     WsMsgCode   `json:"code"`
	Msg      string      `json:"msg,omitempty"`
	UiData   *WsUiData   `json:"uidata,omitempty"`
	Template *WsTemplate `json:"template,omitempty"`
	Config   *WsConfig   `json:"config,omitempty"`
}

func (w *WsMessage) UnmarshalJson(data []byte) error {
	// Clear the websocket data, ready to populate the newly received string
	w.Code = ""
	w.Msg = ""
	w.UiData = nil
	w.Template = nil
	w.Config = nil

	err := json.Unmarshal(data, &w)
	if err != nil {
		return err
	}

	// Check that we have the correct key and no extra keys
	checks := []struct {
		code WsMsgCode
		ok   bool
	}{
		{WscConfig, w.Config != nil},
		{WscTemplate, w.Template != nil},
		{WscUiData, w.UiData != nil},
	}
	for _, c := range checks {
		if c.code == w.Code && !c.ok {
			return fmt.Errorf("expected '%s' in %s", c.code, string(data))
		}
		if c.code != w.Code && c.ok {
			return fmt.Errorf("did not expect '%s' in %s", c.code, string(data))
		}

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
	WsLogf(conn, WscError, msg, args...)
}

func WsWarnf(conn *websocket.Conn, msg string, args ...interface{}) {
	WsLogf(conn, WscWarn, msg, args...)
}

func WsLogf(conn *websocket.Conn, code WsMsgCode, msg string, args ...interface{}) {
	m := fmt.Sprintf(msg, args...)
	wsmsg := &WsMessage{
		Code: code,
		Msg:  m,
	}
	WsWriteJSON(conn, wsmsg)
}
