package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/log"
)

type WsMsgCode string

const (
	WscError    WsMsgCode = "error"    // .Msg
	WscWarn     WsMsgCode = "warn"     // .Msg
	WscTemplate WsMsgCode = "template" // .Template
	WscUiData   WsMsgCode = "uidata"   // .UiData
	WscConfig   WsMsgCode = "config"   // .Config
	WscFsChange WsMsgCode = "fschange" // .Msg contains the filename
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

// WsMessage stringer
func (w *WsMessage) String() string {
	res, err := json.Marshal(&w)
	if err != nil {
		log.Error(err)
		return ""
	}
	return string(res)
}

func WsErrorf(ws chan<- any, msg string, args ...any) {
	ws <- &WsMessage{Code: WscError, Msg: fmt.Sprintf(msg, args...)}
}

func WsWarnf(ws chan<- any, msg string, args ...any) {
	ws <- &WsMessage{Code: WscWarn, Msg: fmt.Sprintf(msg, args...)}
}

// a FSChange message
func WsFsChange(msg string) *WsMessage {
	return &WsMessage{Code: WscFsChange, Msg: msg}
}
