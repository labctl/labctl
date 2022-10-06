package helpers

import (
	"os"
	"path/filepath"
)

// Load the template from the ctx.TemplatePaths
func LoadFiles(ctx *Context) (map[string]string, error) {
	res := make(map[string]string)
	path := filepath.Dir(ctx.TopoFile)
	glob, err := filepath.Glob(filepath.Join(path, "*.md"))
	if err != nil {
		res["readme.md"] = err.Error()
		return res, nil
	}
	for _, fn := range glob {
		b, err := os.ReadFile(fn)
		if err != nil {
			res[fn+" *"] = err.Error()
			continue
		}
		res[fn] = string(b[:])
	}
	return res, nil
}
