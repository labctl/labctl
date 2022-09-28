package helpers

import (
	"text/template"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Context struct {
	// the Kong Command()
	Command    string
	DebugCount int
	Settings   *Settings
	Async      bool

	// used by config, serve
	TopoFile string
	// used by config, serve
	TemplatePaths *orderedmap.OrderedMap[string, string]
	// used by config
	TemplateList []string
	// Used by config
	NodeFilter []string

	Topo Topo

	Template *template.Template

	// Output of any config command. Either to the terminal or websocket
	Output ResultOutput
}

func (c *Context) Load(filename string) error {
	c.TopoFile = filename
	return c.Topo.Load(c.TopoFile)
}

// Convert the TemplatePath stored in an orderedmap to a []string
func (c *Context) TemplatePathsSlice() []string {
	res := make([]string, 0, c.TemplatePaths.Len())
	for pair := c.TemplatePaths.Oldest(); pair != nil; pair = pair.Next() {
		res = append(res, pair.Value)
	}
	return res
}
