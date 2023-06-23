package helpers

import (
	"os"
	"path/filepath"

	"github.com/labctl/labctl/utils"
	log "github.com/sirupsen/logrus"
)

type Template struct {
	Name    string   `json:"name"`
	P       string   `json:"p"`
	ShadowP []string `json:"shadow"`
	Value   string   `json:"value"`
}

type Templates map[string]*Template

func (t *Template) Load(ctx *Context) error {
	p, ok := ctx.TemplatePaths.Get(t.P)
	if !ok {
		log.Fatal("x")
	}
	fn := filepath.Join(p, t.Name)

	b, err := os.ReadFile(fn)
	if err != nil {
		return err
	}
	t.Value = string(b[:])
	// log.Errorf("%s --> %d", t.Name, len(t.Value))
	return nil
}

// Load the template from the ctx.TemplatePaths
func LoadTemplates(ctx *Context) (Templates, error) {
	res := make(Templates)

	for pp := ctx.TemplatePaths.Oldest(); pp != nil; pp = pp.Next() {
		// log.Warnf("%s", pp.Value)
		p := utils.Path{Path: pp.Value}
		lres, err := p.ReadFiles("*.tmpl")
		if err != nil {
			log.Errorf("Error reading templates in %s: %s", pp.Value, err)
			continue
		}
		for fname, fcontent := range lres {
			t, ok := res[fname]
			if ok {
				// it exists...
				t.ShadowP = append(t.ShadowP, t.P)
				t.P = pp.Key
				t.Value = fcontent
			} else {
				res[fname] = &Template{
					Name:    fname,
					P:       pp.Key,
					ShadowP: make([]string, 0),
					Value:   fcontent,
				}
			}
		}
	}
	return res, nil
}
