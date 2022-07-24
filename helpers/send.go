package helpers

import (
	"fmt"

	"github.com/labctl/labctl/utils/tx"

	"github.com/srl-labs/containerlab/clab/config"
)

const vkTransport = "transport"

// Send configuration prepared by containerlab according to the action
func ConfigSend(c *config.NodeConfig, action tx.Action) ([]*tx.Response, error) {

	// Get the transport
	transport, ok := c.TargetNode.Config.Vars[vkTransport]
	if !ok {
		transport = "ssh"
	}

	if transport == "ssh" {
		s := tx.SSHTx{}
		i, err := s.Prepare(c)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			return nil, nil
		}
		return s.Send(action)
	}

	return nil, fmt.Errorf("transport '%s' not implemented", transport)
}
