package helpers

import (
	"github.com/labctl/labctl/utils/tx"
)

type WsConfig struct {
	Cmd     string         `json:"cmd"`
	Node    string         `json:"node"`
	Sresult string         `json:"sresult"`
	Result  []*tx.Response `json:"result"`
}
