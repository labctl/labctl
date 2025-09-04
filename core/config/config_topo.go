package config

import (
	"github.com/charmbracelet/log"
	clabcore "github.com/srl-labs/containerlab/core"
	clablinks "github.com/srl-labs/containerlab/links"
	clabtypes "github.com/srl-labs/containerlab/types"
)

// ConfigTopo is a struct that contains the config input (name, node&links with vars)
// from the topology file. Can be marshalled as JSON.
type ConfigTopo struct {
	Name   string                `json:"name"`
	Prefix string                `json:"prefix"`
	Path   *clabtypes.TopoPaths  `json:"-"`
	Nodes  map[string]ConfigNode `json:"nodes,omitempty"`
	Links  map[int]ConfigLink    `json:"links,omitempty"`
}

type ConfigNode struct {
	clabtypes.ContainerDetails
	Vars map[string]any `json:"vars,omitempty"`
}

type ConfigLink struct {
	clabcore.Link
	Vars map[string]any `json:"vars,omitempty"`
}

func NewConfigTopo(c *clabcore.CLab) ConfigTopo {
	res := ConfigTopo{
		Name:   c.Config.Name,
		Path:   c.TopoPaths,
		Prefix: *c.Config.Prefix,
		Nodes:  make(map[string]ConfigNode, 0),
		Links:  make(map[int]ConfigLink, 0),
	}

	for _, n := range c.Nodes {
		nn := n.Config().ShortName
		res.Nodes[nn] = ConfigNode{
			ContainerDetails: clabtypes.ContainerDetails{
				Name:        nn,
				Kind:        n.Config().Kind,
				Image:       n.Config().Image,
				Group:       n.Config().Group,
				IPv4Address: n.Config().MgmtIPv4Address,
				IPv6Address: n.Config().MgmtIPv6Address,
			},
			Vars: n.Config().Config.Vars,
		}
	}
	if len(c.Links) == 0 {
		err := c.ResolveLinks()
		if err != nil {
			log.Error("Could not resolve topofile links", "err", err)
		}
		if len(c.Links) == 0 {
			log.Warn("No links defined in the topology")
		}
	}
	for i, l := range c.Links {
		res.Links[i] = *NewConfigLink(l)
	}
	return res
}

func NewConfigLink(lnk clablinks.Link) *ConfigLink {
	res := &ConfigLink{
		Link: clabcore.Link{
			Source:         lnk.GetEndpoints()[0].GetNode().GetShortName(),
			SourceEndpoint: lnk.GetEndpoints()[0].GetIfaceName(),
			Target:         lnk.GetEndpoints()[1].GetNode().GetShortName(),
			TargetEndpoint: lnk.GetEndpoints()[1].GetIfaceName(),
		},
	}

	// casting required until the containerlab interface is fixed to expose the link vars
	if v, ok := lnk.(*clablinks.LinkVEth); ok {
		res.Vars = v.Vars
	} else {
		log.Warnf("Link %s-%s is not a virtual Ethernet link", res.Source, res.Target)
	}
	return res
}

func (l ConfigLink) String() string {
	return l.Source + "-" + l.Target
}
