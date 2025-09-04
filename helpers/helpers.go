package helpers

import (
	"fmt"
	"strings"

	clabnodes "github.com/srl-labs/containerlab/nodes"
)

// Debug Level
var DebugCount int

// Ensure the NodeFilter contains valid nodes
// If empty, ensure it contains ALL the nodes
func ValidateNodeFilter(NodeFilter *[]string, nodes map[string]clabnodes.Node) error {
	if len(*NodeFilter) == 0 {
		for n := range nodes {
			*NodeFilter = append(*NodeFilter, n)
		}
		return nil
	}
	mis := []string{}
	for _, nn := range *NodeFilter {
		if _, ok := nodes[nn]; !ok {
			mis = append(mis, nn)
		}
	}
	if len(mis) > 0 {
		return fmt.Errorf("invalid nodes in filter: %s", strings.Join(mis, ", "))
	}
	return nil
}
