package utils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Path struct {
	Path string
}

// KndPath represents a name__kind.ext path
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
		p := Path{Path: "~/go/src/github.com/labctl/labctl/templates"}
		err := p.ExpandUser()
		if err != nil {
			return "", err
		}
		dirs = []string{p.Path, "/etc/labctl/templates", "."}
	}
	p := Path{Path: f.String()}
	if p.Exists() == nil {
		return p.Path, nil
	}
	if f.dir == "" {
		for i := len(dirs) - 1; i >= 0; i-- {
			_p := Path{Path: filepath.Join(dirs[i], p.Path)}
			if _p.Exists() == nil {
				return _p.Path, nil
			}
		}
	}
	return "", fmt.Errorf("file not found %s", p.Path)
}

func (p *Path) String() string {
	return p.Path
}

// make this an absolute path (remove ..), expand user and test if exists (strict)
// Returns nil if all ok
func (p *Path) Resolve(strict ...bool) error {
	if len(strict) == 0 {
		strict = []bool{false}
	}
	err := p.ExpandUser()
	if err != nil {
		return err
	}

	if !filepath.IsAbs(p.Path) {
		p.Path, err = filepath.Abs(p.Path)
		if err != nil {
			return err
		}
	}
	if strict[0] {
		return p.Exists()
	}
	return nil
}

// returns nil if file exists, else error
func (p *Path) Exists() error {
	if p.Path == "" {
		return fmt.Errorf("empty path")
	}
	_, err := os.Stat(p.Path)
	if err == nil {
		return nil
	} else if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", p.Path)
	}
	return fmt.Errorf("cannot access file %s: %s", p.Path, err)
}

// expand the user's home directory in a file path
func (p *Path) ExpandUser() error {
	if p.Path == "" || p.Path[0] != '~' {
		return nil
	}
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("cannot expand home directory: %s", err)
	}
	p.Path = filepath.Join(usr.HomeDir, p.Path[1:])
	return nil
}
