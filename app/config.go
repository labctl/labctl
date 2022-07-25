package app

import (
	"fmt"
	"strings"

	"github.com/labctl/labctl/helpers"
	"github.com/labctl/labctl/utils/tx"
	"github.com/srl-labs/containerlab/clab/config"

	log "github.com/sirupsen/logrus"
)

// Paths to search for templates
//var TemplatePaths []string

// Template names to use
//var TemplateNames []string

// Render configuration and perform the action
func ConfigRun(actionStr string, ctx *Context) error {

	action, err := tx.ParseActionString(actionStr)
	if err != nil {
		return err
	}

	tx.DebugCount = ctx.DebugCount
	config.DebugCount = ctx.DebugCount

	allConfig, err := LoadAndPrep(ctx.NodeFilter, ctx.TopoFile)
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

	for _, node := range ctx.NodeFilter {

		cs, ok := allConfig[node]
		if !ok {
			log.Fatalf("Invalid node in filter: %s", node)
		}

		_, err = ConfigSend(cs, action)
		if err != nil {
			log.Errorf("%s: %s", cs.TargetNode.ShortName, err)
		}

	}

	return nil
}

const vkTransport = "transport"

// Send configuration to the action
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

func ConfigView(actionStr string, ctx *Context) error {
	allConfig, err := LoadAndPrep(ctx.NodeFilter, ctx.TopoFile)
	if err != nil {
		return err
	}
	log.Printf("a")
	if ctx.Command == "vars" {
		for _, n := range ctx.NodeFilter {
			allConfig[n].Print(true, false)
		}
		return nil
	}
	log.Printf("b")
	ll := log.GetLevel()
	log.SetLevel(log.ErrorLevel)
	err = config.RenderAll(allConfig)
	log.SetLevel(ll)
	if err != nil {
		return err
	}
	log.Printf("c")
	for _, n := range ctx.NodeFilter {
		allConfig[n].Print(false, true)
	}
	return nil

}

// Load the topo files and prepare the variables
// topoFiles is OPTIONAL. If not supplied, use the config flag TopoFiles
func LoadAndPrep(nodeFilter []string, topoFiles ...string) (map[string]*config.NodeConfig, error) {

	nodes, links, err := helpers.LoadTopoFiles(topoFiles...)
	if err != nil {
		return nil, err
	}

	err = helpers.ValidateNodeFilter(nodeFilter, nodes)
	if err != nil {
		return nil, err
	}

	return config.PrepareVars(nodes, links), nil
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
