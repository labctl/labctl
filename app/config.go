package app

import (
	"fmt"
	"strings"

	"github.com/labctl/labctl/helpers"
	"github.com/labctl/labctl/utils"
	"github.com/labctl/labctl/utils/tx"
	"github.com/srl-labs/containerlab/clab/config"

	log "github.com/sirupsen/logrus"
)

// Send commands or commit/compare configuration to a single node
func ConfigTx(configs map[string]*config.NodeConfig, ctx *helpers.Context) error {
	tx.DebugCount = ctx.DebugCount

	act, err := tx.StringToAction(ctx.Command)
	if err != nil {
		return err
	}
	for nn, cfg := range configs {
		if !utils.Contains(ctx.NodeFilter, nn) {
			continue
		}

		resp, err := ConfigTx1(cfg, act)
		if err != nil {
			ctx.Output.Error(cfg.TargetNode.ShortName, err.Error())
		}

		ctx.Output.LogResponses(resp, cfg)
	}

	return nil
}

const vkTransport = "transport"

// Send commands or commit/compare configuration to a single node
func ConfigTx1(c *config.NodeConfig, action tx.Action) ([]*tx.Response, error) {
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
			log.Infof("no config for %s", c.TargetNode.ShortName)
			return nil, nil
		}
		return s.Send(action)
	}

	return nil, fmt.Errorf("transport '%s' not implemented", transport)
}

// Load the topo files and prepare the variables
func LoadAndPrep(nodeFilter *[]string, topoFile string, render bool) (map[string]*config.NodeConfig, error) {
	var topo helpers.Topo
	err := topo.Load(topoFile)
	if err != nil {
		return nil, err
	}

	err = helpers.ValidateNodeFilter(nodeFilter, topo.Nodes)
	if err != nil {
		return nil, err
	}

	configs := config.PrepareVars(topo.Nodes, topo.Links)

	if render {
		config.DebugCount = Ctx.DebugCount

		// Call containerlab's render
		err = config.RenderAll(configs)
		if err != nil {
			return nil, err
		}

		// Validate render
		err = validateRender(configs)
		if err != nil {
			return nil, err
		}
	}

	return configs, nil
}

// Validate the rendered template
// Check for any occurrences of "<no value>" / "<nil>"
func validateRender(allConfig map[string]*config.NodeConfig) error {
	fail_count := 0
	for node, c := range allConfig {
		for idx, ic := range c.Data {
			if strings.Contains(ic, "<no value>") {
				log.Warnf("%s: template %s rendered with <no value>", node, c.Info[idx])
				fail_count += 1
			}
			if strings.Contains(ic, "<nil>") {
				log.Warnf("%s: template %s rendered with <no value>", node, c.Info[idx])
				fail_count += 1
			}
		}
	}

	if fail_count > 0 {
		return fmt.Errorf("%d badly rendered templates", fail_count)
	}
	return nil
}
