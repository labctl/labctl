package app

import (
	"io/ioutil"

	"github.com/labctl/labctl/helpers/discovery"
	"github.com/labctl/labctl/utils/tx"
	log "github.com/sirupsen/logrus"
)

func DiscoverTopo(templateName string, ctx *Context) (map[string][]map[string]interface{}, error) {

	allConfig, err := LoadAndPrep(ctx.NodeFilter, templateName)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]map[string]interface{})

	// discovery templates per kind
	allDis := make(map[string]*discovery.DiscoverTemplate)

	for _, node := range ctx.NodeFilter {

		nc := allConfig[node]

		// Get the discovery template per kind, or load if not there yet
		dis, ok := allDis[nc.TargetNode.Kind]
		if !ok {
			dis, err = discovery.NewDiscoveryConfig(templateName, nc.TargetNode.Kind)
			if err != nil {
				return nil, err
			}
			allDis[nc.TargetNode.Kind] = dis
		}

		nc.Data = []string{dis.Command}
		nc.Info = []string{dis.Command}

		cli, err := ConfigSend(nc, tx.ASend)
		if err != nil {
			return nil, err
		}

		if len(cli) == 0 {
			// command likely failed?
			log.Errorf("No reply from %s", node)
			continue
		}

		res, err := dis.ProcessShow(cli[0].Response, "")
		if err != nil {
			return nil, err
		}
		result[node] = res
	}

	return result, nil
}

func DiscoverFile(template, infile string) ([]map[string]interface{}, error) {
	dis := discovery.DiscoverTemplate{}
	err := dis.Load(template)
	if err != nil {
		return nil, err
	}
	input, err := ioutil.ReadFile(infile)
	if err != nil {
		return nil, err
	}
	res, err := dis.ProcessShow(string(input), "")
	if err != nil {
		return nil, err
	}
	return res, nil
}
