package tx

import (
	"fmt"

	"github.com/srl-labs/containerlab/clab/config"
)

var DebugCount = 0

type Action int

const (
	AUnknown Action = iota
	ASend
	ACompare
	ACommit
)

var actionStrs = []string{"unknown", "send", "compare", "commit"}

func (c Action) String() string {
	return actionStrs[c]
}

// Parse a string and return an action
func StringToAction(str string) (Action, error) {
	for i, a := range actionStrs {
		if i != int(AUnknown) && a == str {
			return Action(i), nil
		}
	}
	return AUnknown, fmt.Errorf("unknown action %s", str)
}

// Tx Interface to prepare config and to send it
type Tx interface {
	Prepare(c *config.NodeConfig) (int, error)
	Send(action Action) ([]*Response, error)
}
