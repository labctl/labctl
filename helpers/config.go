package helpers

import (
	"fmt"
	"strings"

	"github.com/srl-labs/containerlab/clab/config"
	"github.com/srl-labs/containerlab/clab/config/transport"

	log "github.com/sirupsen/logrus"
)

// Filter specific nodes
var NodeFilter []string

// Paths to search for templates
var TemplatePaths []string

// Template names to use
var TemplateNames []string

// Render configuration and perform the action
func ConfigRun(action string) error {

	err := validateSendAction(action)
	if err != nil {
		return err
	}

	transport.DebugCount = DebugCount
	config.DebugCount = DebugCount

	allConfig, err := LoadAndPrep()
	if err != nil {
		return err
	}

	config.TemplateNames = TemplateNames
	config.TemplatePaths = TemplatePaths
	err = config.RenderAll(allConfig)
	if err != nil {
		return err
	}

	err = validateRender(allConfig)
	if err != nil {
		return err
	}

	for _, node := range NodeFilter {

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

// Load the topo files and prepare the variables
// topoFiles is OPTIONAL. If not supplied, use the config flag TopoFiles
func LoadAndPrep(topoFiles ...string) (map[string]*config.NodeConfig, error) {

	nodes, links, err := LoadTopoFiles(topoFiles...)
	if err != nil {
		return nil, err
	}

	return config.PrepareVars(nodes, links), nil
}

func validateRender(allConfig map[string]*config.NodeConfig) error {
	// Check for any occurrences of "<no value>" / "<nil>"
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
