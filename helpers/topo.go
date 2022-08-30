package helpers

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"time"

	"github.com/srl-labs/containerlab/clab"
	"github.com/srl-labs/containerlab/clab/config"
	"github.com/srl-labs/containerlab/nodes"
	"github.com/srl-labs/containerlab/types"

	log "github.com/sirupsen/logrus"
)

type Topo struct {
	Name  string
	Nodes map[string]nodes.Node
	Links map[int]*types.Link
}

func (topo *Topo) Load(topoFile string) error {
	dl := log.GetLevel()
	defer log.SetLevel(dl)
	log.SetLevel(log.InfoLevel)
	c, err := clab.NewContainerLab(
		clab.WithTimeout(time.Second*30),
		clab.WithTopoFile(topoFile, ""),
	)
	log.SetLevel(dl)
	log.Debugf("Loaded %s", topoFile)
	if err != nil {
		return err
	}
	topo.Links = c.Links
	topo.Nodes = c.Nodes
	topo.Name = c.Config.Name
	return nil
}

type NodeJson struct {
	types.ContainerDetails
	Vars map[string]interface{} `json:"vars,omitempty"`
}

type LinkJson struct {
	clab.Link
	Vars map[string]interface{} `json:"vars,omitempty"`
}

type TopoJson struct {
	Name  string              `json:"name"`
	Nodes map[string]NodeJson `json:"nodes,omitempty"`
	Links map[int]LinkJson    `json:"links,omitempty"`
}

type Vars map[string]map[string]interface{}

func (topo *Topo) VarsAsJson() (Vars, error) {
	nc := config.PrepareVars(topo.Nodes, topo.Links)
	v := make(Vars)
	for _, node := range topo.Nodes {
		name := node.Config().ShortName
		v[name] = nc[name].Vars

		host := node.Config().LongName
		if ipv4, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", host); err == nil {
			if len(ipv4) == 1 {
				v[name]["clab_management_ipv4"] = ipv4[0]
			} else {
				v[name]["clab_management_ipv4"] = ipv4
			}
		}
		if ipv6, err := net.DefaultResolver.LookupIP(context.Background(), "ip6", host); err == nil {
			if len(ipv6) == 1 {
				v[name]["clab_management_ipv4"] = ipv6[0]
			} else {
				v[name]["clab_management_ipv6"] = ipv6
			}
		}

	}
	return v, nil
}

func (topo *Topo) AsJson() (TopoJson, error) {
	res := TopoJson{
		Name:  topo.Name,
		Nodes: make(map[string]NodeJson, 0),
		Links: make(map[int]LinkJson, 0),
	}

	for _, n := range topo.Nodes {
		nn := n.Config().ShortName
		res.Nodes[nn] = NodeJson{
			ContainerDetails: types.ContainerDetails{
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
	for i, l := range topo.Links {
		res.Links[i] = LinkJson{
			Link: clab.Link{
				Source:         l.A.Node.ShortName,
				SourceEndpoint: l.A.EndpointName,
				Target:         l.B.Node.ShortName,
				TargetEndpoint: l.B.EndpointName,
			},
			Vars: l.Vars,
		}
	}

	return res, nil
}

func (topo *Topo) Print() error {
	// j, err := topo.AsJson()
	j, err := topo.VarsAsJson()
	if err != nil {
		return err
	}

	by, _ := json.MarshalIndent(j, "", "  ")
	_, err = os.Stdout.Write(by)
	return err
}
