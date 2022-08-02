package helpers

import (
	"context"
	"strings"
	"text/template"

	"github.com/hairyhenderson/gomplate/v3"
	gomData "github.com/hairyhenderson/gomplate/v3/data"
	jT "github.com/kellerza/template"
)

type WebSocketTemplate struct {
	Name     string                 `json:"name,omitempty"`
	Template string                 `json:"template,omitempty"`
	Vars     map[string]interface{} `json:"vars,omitempty"`
	Result   string                 `json:"result,omitempty"`
}

func ParseTemplates(temps map[string]string) (*template.Template, error) {
	var err error
	// gomplate overrides the built-in *slice* function. You can still use *coll.Slice*
	gfuncs := gomplate.CreateFuncs(context.Background(), new(gomData.Data))
	delete(gfuncs, "slice")

	res := template.New("").Funcs(gfuncs).Funcs(jT.Funcs)

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
	err := ctx.Template.ExecuteTemplate(&buf, t.Name, t.Vars)
	if err != nil {
		return err
	}
	t.Result = buf.String()
	return nil
}
