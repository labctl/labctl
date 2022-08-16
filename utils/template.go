package utils

import (
	"context"
	"text/template"

	"github.com/hairyhenderson/gomplate/v3"
	gomData "github.com/hairyhenderson/gomplate/v3/data"
	jT "github.com/kellerza/template"
)

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
