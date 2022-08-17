package helpers

import (
	"github.com/labctl/labctl/utils/tx"
	"github.com/srl-labs/containerlab/clab/config"
)

type WsConfig struct {
	Cmd     string             `json:"cmd,omitempty"`
	Results []*tx.Response     `json:"results,omitempty"`
	Input   *config.NodeConfig `json:"env"`
}
