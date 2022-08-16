package app

import (
	"fmt"
	"strings"

	"github.com/labctl/labctl/helpers"
	"github.com/labctl/labctl/utils/tx"
	"github.com/srl-labs/containerlab/clab/config"

	log "github.com/sirupsen/logrus"
)

// Render configuration and perform the action
func ConfigRun(actionStr string, ctx *helpers.Context) error {
	action, err := tx.StringToAction(actionStr)
	if err != nil {
		return err
	}

	tx.DebugCount = ctx.DebugCount
	config.DebugCount = ctx.DebugCount

	allConfig, err := LoadAndPrep(&ctx.NodeFilter, ctx.TopoFile)
	if err != nil {
		return err
	}

	err = config.RenderAll(allConfig)
	if err != nil {
		return err
	}

	err = validateRender(allConfig)
	if err != nil {
		return err
	}

	if len(ctx.NodeFilter) == 0 {
		// Ad al the nodes to the filter
		for n := range allConfig {
			ctx.NodeFilter = append(ctx.NodeFilter, n)
		}
	}

	for _, node := range ctx.NodeFilter {

		cs, ok := allConfig[node]
		if !ok {
			return fmt.Errorf("invalid node in filter: %s", node)
		}

		resp, err := ConfigSend(cs, action)
		if err != nil {
			ctx.Output.Error(cs.TargetNode.ShortName, err.Error())
		}
		ctx.Output.Info("a", "zz")

		LogResults(resp, ctx)

	}

	return nil
}

// Display all results
func LogResults(results []*tx.Response, ctx *helpers.Context) {
}

const vkTransport = "transport"

// Send commands or commit/compare configuration
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
			log.Infof("no config for %s", c.TargetNode.ShortName)
			return nil, nil
		}
		return s.Send(action)
	}

	return nil, fmt.Errorf("transport '%s' not implemented", transport)
}

func ConfigView(actionStr string, ctx *helpers.Context) error {
	allConfig, err := LoadAndPrep(&ctx.NodeFilter, ctx.TopoFile)
	if err != nil {
		return err
	}
	if ctx.Command == "vars" {
		for _, n := range ctx.NodeFilter {
			allConfig[n].Print(true, false)
		}
		return nil
	}
	err = config.RenderAll(allConfig)
	if err != nil {
		return err
	}
	for _, n := range ctx.NodeFilter {
		allConfig[n].Print(false, true)
	}
	return nil
}

// Load the topo files and prepare the variables
func LoadAndPrep(nodeFilter *[]string, topoFile string) (map[string]*config.NodeConfig, error) {
	var topo helpers.Topo
	err := topo.Load(topoFile)
	if err != nil {
		return nil, err
	}

	err = helpers.ValidateNodeFilter(nodeFilter, topo.Nodes)
	if err != nil {
		return nil, err
	}

	return config.PrepareVars(topo.Nodes, topo.Links), nil
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
