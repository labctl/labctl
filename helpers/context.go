package helpers

import (
	"text/template"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Context struct {
	Command    string // the Kong Command()
	DebugCount int
	Settings   Settings

	TopoFile      string   // used by config, serve
	NodeFilter    []string // Used by config
	TemplatePaths *orderedmap.OrderedMap[string, string]

	Topo Topo

	Template *template.Template
}

func (c *Context) Load(filename string) error {
	c.TopoFile = filename
	return c.Topo.Load(c.TopoFile)
}
