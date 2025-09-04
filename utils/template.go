package utils

import (
	"maps"
	"text/template"

	jT "github.com/kellerza/template"
	clabutils "github.com/srl-labs/containerlab/utils"
)

func AllFuncs() template.FuncMap {
	res := clabutils.CreateFuncs()
	maps.Copy(res, jT.Funcs)
	// log.Infof("Loaded %d template functions: %v", len(res), res)
	return res
}

func NewTemplate() *template.Template {
	return template.New("").Funcs(AllFuncs())
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
