package helpers

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type layoutsNodes map[string]struct {
	X, Y float32
}

type UiData struct {
	Layouts struct {
		Nodes layoutsNodes
	} // layouts is directly for v-network-graph
	Options map[string]interface{} // A generic options dictionary for UI related options
}

func NewUiData() UiData {
	return UiData{
		Layouts: struct{ Nodes layoutsNodes }{
			Nodes: make(layoutsNodes, 0),
		},
		Options: make(map[string]interface{}, 0),
	}
}

func (u *UiData) Print() {
	log.Printf("%v", u)
	// for k, v := range u.Layouts.Nodes {
	// 	log.Infof("%s x:%v y:%v", k, v.X, v.Y)
	// }
}

type WebSocketData struct {
	Code int
	Msg  string
	Data UiData `json:"data,omitempty"`
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

// type Dict struct {
// 	Data map[string]interface{}
// }

// func (d *Dict) Decode(data []byte) error {
// 	err := json.Unmarshal([]byte(data), &d.Data)
// 	return err
// }
// func (d *Dict) Print() {
// 	for k := range d.Data {
// 		log.Infof(k)
// 	}
// 	//log.Printf("%v", u.Layouts.Nodes)
// }
