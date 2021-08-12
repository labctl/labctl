package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/srl-labs/containerlab/clab"
	"github.com/srl-labs/containerlab/nodes"
	"github.com/srl-labs/containerlab/types"
)

// Debug Level
var DebugCount int

// path to the topology file
var TopoFiles []string

func LoadTopoFiles(topoFiles ...string) (map[string]nodes.Node, map[int]*types.Link, error) {

	var nodes map[string]nodes.Node
	var links map[int]*types.Link

	if len(topoFiles) == 0 {
		return nil, nil, fmt.Errorf("specify at least one topology file (--topo / -t)")
	}

	for i, topo := range topoFiles {
		c, err := clab.NewContainerLab(
			clab.WithTimeout(time.Second*30),
			clab.WithTopoFile(topo),
		)
		if err != nil {
			return nil, nil, err
		}
		if i == 0 {
			nodes = c.Nodes
			links = c.Links
			continue
		}
		// Do some merging
	}

	err := ValidateNodeFilter(nodes)
	if err != nil {
		return nil, nil, err
	}
	return nodes, links, nil
}

// Ensure the NodeFilter contains valid nodes
// If empty, ensure it contains ALL the nodes
func ValidateNodeFilter(nodes map[string]nodes.Node) error {
	if len(NodeFilter) == 0 {
		for n := range nodes {
			NodeFilter = append(NodeFilter, n)
		}
		return nil
	}
	mis := []string{}
	for _, nn := range NodeFilter {
		if _, ok := nodes[nn]; !ok {
			mis = append(mis, nn)
		}
	}
	if len(mis) > 0 {
		return fmt.Errorf("invalid nodes in filter: %s", strings.Join(mis, ", "))
	}
	return nil
}
