package helpers

import (
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/hairyhenderson/gomplate/v3"
	gomData "github.com/hairyhenderson/gomplate/v3/data"
	jT "github.com/kellerza/template"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type WebSocketTemplate struct {
	Name       string                 `json:"name,omitempty"`
	Template   string                 `json:"template,omitempty"`
	Vars       map[string]interface{} `json:"vars,omitempty"`
	Result     string                 `json:"result,omitempty"`
	ResultYaml map[string]interface{} `json:"resulty,omitempty"`
}

func NewTemplate() *template.Template {
	// gomplate overrides the built-in *slice* function. You can still use *coll.Slice*
	gfuncs := gomplate.CreateFuncs(context.Background(), new(gomData.Data))
	delete(gfuncs, "slice")
	return template.New("").Funcs(gfuncs).Funcs(jT.Funcs)
}

func ParseTemplates(temps map[string]string) (*template.Template, error) {
	var err error
	res := NewTemplate()

	// Parse all templates
	for k, t := range temps {
		_, err = res.New(k).Parse(t)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// Render the template
func (t *WebSocketTemplate) Render(ctx *Context) error {
	var buf strings.Builder

	if t.Name == "" {
		return fmt.Errorf("template should have a name")
	}
	if ctx.Template == nil {
		ctx.Template = NewTemplate()
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
