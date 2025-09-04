package helpers

import (
	"github.com/labctl/labctl/core/config"
	"github.com/labctl/labctl/utils/tx"
)

type WsConfig struct {
	Cmd     string             `json:"cmd,omitempty"`
	Results []*tx.Response     `json:"results,omitempty"`
	Input   *config.NodeConfig `json:"env"`
}
