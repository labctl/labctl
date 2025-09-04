package app

import (
	"fmt"
	"strings"
	"sync"

	"github.com/labctl/labctl/core/config"
	"github.com/labctl/labctl/helpers"
	"github.com/labctl/labctl/utils/tx"
	"golang.org/x/exp/slices"

	"github.com/charmbracelet/log"
)

// Send commands or commit/compare configuration to a single node
func ConfigTx(configs map[string]*config.NodeConfig, ctx *helpers.Context) error {
	tx.DebugCount = ctx.DebugCount

	act, err := tx.StringToAction(ctx.Command)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for nn, cfg := range configs {
		if !slices.Contains(ctx.NodeFilter, nn) {
			continue
		}
		if !ctx.Async {
			ConfigTx1PP(ctx, act, nn, cfg)
			continue
		}
		wg.Add(1)
		go func(nn string, cfg *config.NodeConfig) {
			defer wg.Done()
			ConfigTx1PP(ctx, act, nn, cfg)
		}(nn, cfg)
	}
	wg.Wait()

	return nil
}

func ConfigTx1PP(ctx *helpers.Context, act tx.Action, nn string, cfg *config.NodeConfig) {
	resp, err := ConfigTx1(cfg, act)
	if err != nil {
		ctx.Output.Error(cfg.TargetNode.ShortName, err.Error())
	}

	if Ctx.Settings != nil && Ctx.Settings.Colors != nil {
		helpers.ColorResults(resp, ctx.Settings.Colors)
	}

	if len(resp) == 0 {
		// Add a response if this node has no templates/responses
		msg := ""
		if len(cfg.Data) == 0 {
			msg += fmt.Sprintf("No templates for node %s", nn)
		} else {
			msg += fmt.Sprintf("All ok from node %s", nn)
		}
		resp = []*tx.Response{{Node: nn, Response: msg}}
	}
	ctx.Output.LogResponses(resp, cfg)
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
		tx.KindMapOverride(Ctx.Settings.Kindmap)
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
func LoadAndPrep(nodeFilter *[]string, topoFile string, render bool, templateList []string, templatePaths []string) (map[string]*config.NodeConfig, error) {
	c, err := helpers.LoadTopo(topoFile)
	if err != nil {
		return nil, err
	}

	err = helpers.ValidateNodeFilter(nodeFilter, c.Nodes)
	if err != nil {
		return nil, err
	}

	configs := config.PrepareVars(c)

	if render {
		config.DebugCount = Ctx.DebugCount

		// Call containerlab's render
		err = config.RenderAll(configs, templateList, templatePaths)
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
func validateRender(configs map[string]*config.NodeConfig) error {
	fail_count := 0
	for node, c := range configs {
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
