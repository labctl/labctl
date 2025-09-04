package helpers

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"time"

	"github.com/labctl/labctl/core/config"
	clabcore "github.com/srl-labs/containerlab/core"

	"github.com/charmbracelet/log"
)

func LoadTopo(topoFile string) (*clabcore.CLab, error) {
	dl := log.GetLevel()
	defer log.SetLevel(dl)
	log.SetLevel(log.InfoLevel)
	c, err := clabcore.NewContainerLab(
		clabcore.WithTimeout(time.Second*30),
		clabcore.WithTopoPath(topoFile, ""),
	)
	log.SetLevel(dl)
	if err != nil {
		log.Errorf("failed to create containerlab instance: %v", err)
		return nil, err
	}
	log.Debugf("Loaded %s", topoFile)
	return c, nil
}

type Vars map[string]map[string]any

func PrepareVarsAsJson(c *clabcore.CLab) (Vars, error) {
	nc := config.PrepareVars(c)
	v := make(Vars)
	for _, node := range c.Nodes {
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
				v[name]["clab_management_ipv6"] = ipv6[0]
			} else {
				v[name]["clab_management_ipv6"] = ipv6
			}
		}

	}
	return v, nil
}

func PrintPreparedVars(c *clabcore.CLab) error {
	j, err := PrepareVarsAsJson(c)
	if err != nil {
		return err
	}

	by, _ := json.MarshalIndent(j, "", "  ")
	_, err = os.Stdout.Write(by)
	return err
}
