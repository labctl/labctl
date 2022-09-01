package helpers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/labctl/labctl/utils"
	log "github.com/sirupsen/logrus"
	orderedmap "github.com/wk8/go-ordered-map/v2"
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
		log.Warnf("%s", pp.Value)
		p := filepath.Join(pp.Value, "*.tmpl")
		files, err := filepath.Glob(p)
		if err != nil {
			return nil, err
		}
		for _, fname := range files {
			_, fn := filepath.Split(fname)
			t, ok := res[fn]
			if ok {
				// it exists...
				t.ShadowP = append(t.ShadowP, t.P)
				t.P = pp.Key
			} else {
				res[fn] = &Template{
					Name:    fn,
					P:       pp.Key,
					ShadowP: make([]string, 0),
				}
			}
		}
	}
	for _, t := range res {
		err := t.Load(ctx)
		if err != nil {
			return nil, err
		}
		// log.Infof("%s --> %d", t.Name, len(t.Value))

	}
	return res, nil
}

// Resolve TemplatePaths and return an ordered map with unique ID
// The unique ID is the key to the map & will be used by the UI
func InitTemplatePaths(paths []string) (*orderedmap.OrderedMap[string, string], error) {
	res := orderedmap.New[string, string]()
	for _, ps := range paths {
		p := utils.Path{Path: ps}
		err := p.Resolve()
		if err != nil {
			return nil, fmt.Errorf("path %s: %s", ps, err)
		}
		n := filepath.Base(p.Path)
		i := 1
		pval, ok := res.Get(n)
		for ok && pval != p.Path { // ensure a unique name & does not already exist
			tmp := fmt.Sprintf("%s_%d", n, i)
			pval, ok = res.Get(tmp)
			if !ok {
				n = tmp
			}
			i++
		}
		res.Set(n, p.Path)
		log.Debugf("--template-path %s [%s]", p.Path, n)
	}
	return res, nil
}
