package helpers

import (
	"fmt"
	"strings"

	"github.com/labctl/labctl/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Used by the frontend to request a template to be rendered
type WsTemplate struct {
	Id         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Template   string                 `json:"template,omitempty"`
	Vars       map[string]interface{} `json:"vars,omitempty"`
	Result     string                 `json:"result,omitempty"`
	ResultYaml map[string]interface{} `json:"resulty,omitempty"`
}

// Render the template
func (t *WsTemplate) Render(ctx *Context) error {
	var buf strings.Builder

	if t.Name == "" {
		return fmt.Errorf("template should have a name")
	}
	if ctx.Template == nil {
		ctx.Template = utils.NewTemplate()
	}
	tmp := ctx.Template.Lookup(t.Name)
	if tmp == nil {
		tmp = ctx.Template.New(t.Name)
	}
	if t.Template != "" {
		_, err := tmp.Parse(t.Template)
		if err != nil {
			t.Result = fmt.Sprintf("%s", err)
			return err
		}
	}
	err := ctx.Template.ExecuteTemplate(&buf, t.Name, t.Vars)
	if err != nil {
		t.Result = fmt.Sprintf("%s", err)
		return err
	}
	t.Result = buf.String()

	// try parse to object
	t.ResultYaml = make(map[string]interface{})
	err = yaml.Unmarshal([]byte(t.Result), t.ResultYaml)
	if err != nil {
		log.Errorf("%s", t.Result)
	}

	log.Infof("Rendered: %s", t.Result)

	return nil
}

// Clear the template command's input (typically before replying)
func (t *WsTemplate) ClearInput() {
	t.Template = ""
	t.Vars = nil
}
