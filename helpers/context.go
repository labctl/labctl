package helpers

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/labctl/labctl/utils"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Context struct {
	// the Kong Command()
	Command    string
	DebugCount int
	Settings   *Settings
	Async      bool

	// Topology filename, used by config, serve
	TopoFile       string
	LabctlFilename string
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

type ContextJson struct {
	Command       string   `json:"command"`
	TopoFile      string   `json:"topofile"`
	TopoError     string   `json:"topoerror"`
	TemplatePaths []string `json:"template_paths"`
}

func (ctx *Context) AsJson() *ContextJson {
	return &ContextJson{
		Command:       ctx.Command,
		TopoFile:      ctx.TopoFile,
		TopoError:     ctx.Topo.TopoError,
		TemplatePaths: ctx.TemplatePathsSlice(),
	}
}

func (c *Context) Load() error {
	err := c.Topo.Load(c.TopoFile)
	if err != nil {
		c.Topo.TopoError = err.Error()
	} else {
		c.Topo.TopoError = ""
	}
	return err
}

// Convert the TemplatePath stored in an orderedmap to a []string
func (c *Context) TemplatePathsSlice() []string {
	res := make([]string, 0, c.TemplatePaths.Len())
	for pair := c.TemplatePaths.Oldest(); pair != nil; pair = pair.Next() {
		res = append(res, pair.Value)
	}
	return res
}

// Initialize the TopFile and TemplatePaths
func (c *Context) InitPaths(topofile string, paths []string) (string, error) {
	var err error
	c.TopoFile, err = ensureTopo(topofile)
	if err != nil {
		return "", err
	}
	p := utils.Path{Path: c.TopoFile}
	err = p.Resolve()
	if err != nil {
		log.Fatalf("cannot access topo file %s: %s", topofile, err)
	}
	c.TopoFile = p.Path

	// Init the matching labctl file
	c.LabctlFilename = labctlFilename(c.TopoFile)

	c.TemplatePaths, err = initTemplatePaths(paths)
	if err != nil {
		return "", err
	}
	log.Debugf("Topo file: %s, Template paths: %s", c.TopoFile, c.TemplatePathsSlice())
	return c.TopoFile, nil
}

// Try find a local topo file
func ensureTopo(topo string) (string, error) {
	if topo != "" {
		return topo, nil
	}
	files, err := filepath.Glob("*.clab.y*ml")
	if err != nil {
		return "", fmt.Errorf("could not find local topo files: %v", err)
	}
	if len(files) == 0 {
		return "", fmt.Errorf("no local topo files")
	}
	if len(files) > 1 {
		return "", fmt.Errorf("multiple topo files found: %v", strings.Join(files, ", "))
	}
	return files[0], nil
}

// Resolve TemplatePaths and return an ordered map with unique ID
// The unique ID is the key to the map & will be used by the UI
func initTemplatePaths(paths []string) (*orderedmap.OrderedMap[string, string], error) {
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

// Get the labctl filename & mode matching the topology file (.clab.yaml)
func labctlFilename(p string) string {
	ext := path.Ext(p)
	p = p[0 : len(p)-len(ext)]
	p = p[0 : len(p)-len(path.Ext(p))]
	return p + ".labctl" + ext
}
