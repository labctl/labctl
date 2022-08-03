// The discovery package allows test processing of router CLI output
// It leverages TextFSMto parse the output and apply a golang template to the output
package discovery

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/labctl/labctl/utils"
	"github.com/sirikothe/gotextfsm"
	"gopkg.in/yaml.v2"
)

type DiscoverTemplate struct {
	Command        string `yaml:"command"`
	ParseTemplate  string `yaml:"parse"`
	OutputTemplate string `yaml:"output"`
	Kind, Name     string
}

func NewDiscoveryConfig(name string, kind string, dirs ...string) (*DiscoverTemplate, error) {
	path := utils.KindPath{
		Kind: kind,
		Name: name,
		Ext:  ".yml",
	}
	c := &DiscoverTemplate{
		Name: name,
		Kind: kind,
	}
	fn, err := path.Resolve(dirs...)
	if err != nil {
		return nil, err
	}
	return c, c.Load(fn)
}

// Load a yaml filename with the conversion
func (c *DiscoverTemplate) Load(filename string) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return fmt.Errorf("invalid yaml: %v", err)
	}

	if c.Command == "" {
		return fmt.Errorf("discovery template %s requires a command: to execute", filename)
	}
	if c.ParseTemplate == "" {
		return fmt.Errorf("discovery template %s requires a parse: textfsm template", filename)
	}
	if c.OutputTemplate == "" {
		return fmt.Errorf("discovery template %s requires a output: golang template", filename)
	}

	return nil
}

// Parses the show command using the TextFSM template
// Returns TextFSM's parser.Dict
func (c *DiscoverTemplate) ParseShow(show string) ([]map[string]interface{}, error) {
	fsm := gotextfsm.TextFSM{}
	err := fsm.ParseString(c.ParseTemplate)
	if err != nil {
		return nil, fmt.Errorf("error while parsing template: %s", err.Error())
	}

	parser := gotextfsm.ParserOutput{}

	err = parser.ParseTextString(show, fsm, true)
	if err != nil {
		return nil, fmt.Errorf("error while parsing input '%s'", err.Error())
	}

	return parser.Dict, nil
}

// ProcessShow parses the show (ParseShow) and then transforms it using the OutputTemplate
func (c *DiscoverTemplate) ProcessShow(show, custom_output string) ([]map[string]interface{}, error) {
	res, err := c.ParseShow(show)
	if err != nil {
		return nil, err
	}

	// The template (array expected) gets a result: key, which is stripped at the end

	outputt := "result:\n"
	if custom_output != "" {
		outputt += custom_output
	} else {
		outputt += c.OutputTemplate
	}

	t, err := template.New(c.Kind).Parse(outputt)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	err = t.Execute(buf, res)
	if err != nil {
		return nil, err
	}

	yaml_result := make(map[string]interface{})

	err = yaml.Unmarshal(buf.Bytes(), yaml_result)
	if err != nil {
		return nil, fmt.Errorf("invalid yaml, check yor template\n%v,\n\ntemplate:\n%v", err, c.OutputTemplate)
	}

	// take off the result to process the inner array
	no_result_key := yaml_result["result"]

	return utils.ArrayMapify(no_result_key)
}
