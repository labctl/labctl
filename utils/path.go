package utils

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/chigopher/pathlib"
)

// KindPath represents a name__kind.ext path
type KindPath struct {
	Name string
	Kind string
	Ext  string
	dir  string
}

// parse a path into a KindPath
func (f *KindPath) Parse(path string) error {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, f.Ext)
	parts := strings.Split(name, "__")
	if len(parts) < 2 {
		return fmt.Errorf("filename needs a '__' separator: %s", path)
	}
	kind := parts[len(parts)-1]
	name = strings.Join(parts[:len(parts)-2], "__")

	f.Ext = ext
	f.Name = name
	if d := filepath.Dir(path); d != "." {
		f.dir = d
	}
	f.Kind = kind
	return nil
}

func (f *KindPath) String() string {
	fn := fmt.Sprintf("%s__%s.%s", f.Name, f.Kind, strings.Trim(f.Ext, "."))
	if f.dir == "" {
		return fn
	}
	return filepath.Join(f.dir, fn)
}

// resolve the KindPath to a path that exists, searching backward through a list of dirs
func (f *KindPath) Resolve(dirs ...string) (string, error) {
	if len(dirs) == 0 {
		p := NewPathExpandUser("~/go/src/github.com/labctl/labctl/templates")
		p, err := p.ResolveAll()
		if err != nil {
			log.Warn("invalid path", "err", err)
		}
		dirs = []string{p.String(), "/etc/labctl/templates", "."}
	}
	p := pathlib.NewPath(f.String())
	exist, err := p.Exists()
	if exist && err == nil {
		return p.String(), nil
	}
	if f.dir == "" {
		for i := len(dirs) - 1; i >= 0; i-- {
			_p := pathlib.NewPath(dirs[i]).JoinPath(p)
			exist, err := _p.Exists()
			if exist && err == nil {
				return _p.String(), nil
			}
		}
	}
	return "", fmt.Errorf("file not found %s", p.String())
}

func ReadFiles(path *pathlib.Path, pattern string) (map[string]string, error) {
	res := make(map[string]string)
	paths, err := path.Glob(pattern)
	if err != nil {
		res["readme.md"] = err.Error()
		return res, nil
	}

	for _, pth := range paths {
		content, err := pth.ReadFile()
		if err != nil {
			res[pth.Name()+" *"] = err.Error()
			continue
		}
		res[pth.Name()] = string(content[:])
	}
	return res, nil
}

func NewPathExpandUser(path string) *pathlib.Path {
	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err == nil {
			return pathlib.NewPath(usr.HomeDir).Join(path[2:])
		}
		log.Warn("cannot expand home directory", "err", err)
	}
	return pathlib.NewPath(path)
}
