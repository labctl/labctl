package helpers

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type layoutsNodes map[string]struct {
	X float32 `json:"x" yaml:"x,flow"`
	Y float32 `json:"y" yaml:"y,flow"`
}
type layouts struct {
	Nodes layoutsNodes `json:"nodes"`
}

type UiData struct {
	Options   map[string]interface{} `json:"options"`   // A generic options dictionary for persistent UI options
	Layouts   layouts                `json:"layouts"`   // layouts is a property of v-network-graph
	Templates map[string]string      `json:"templates"` // Templates used by the UI etc
}

func NewUiData() UiData {
	return UiData{
		Layouts: layouts{
			Nodes: make(layoutsNodes, 0),
		},
		Options:   make(map[string]interface{}, 0),
		Templates: make(map[string]string),
	}
}

func (u *UiData) Print() {
	log.Printf("%v", u)
	// for k, v := range u.Layouts.Nodes {
	// 	log.Infof("%s x:%v y:%v", k, v.X, v.Y)
	// }
}

// remove extension .clab.yaml
func labFileName(p string) (string, fs.FileMode) {
	s, err := os.Lstat(p)
	fm := s.Mode()
	if err != nil {
		log.Errorf("Unable to read %s: %s", p, err.Error())
		fm = 0o664
	}
	ext := path.Ext(p)
	p = p[0 : len(p)-len(ext)]
	p = p[0 : len(p)-len(path.Ext(p))]
	return p + ".labctl" + ext, fm
}

func (u *UiData) WriteFile(ctx *Context) {
	labfn, fm := labFileName(ctx.TopoFile)
	log.Infof("Saving as %s", labfn)

	data, err := yaml.Marshal(&u)
	if err != nil {
		log.Error(err)
		return
	}
	err = ioutil.WriteFile(labfn, data, fm)
	if err != nil {
		log.Error(err)
	}
}

func (u *UiData) ReadFile(ctx *Context) error {
	labfn, _ := labFileName(ctx.TopoFile)
	data, err := ioutil.ReadFile(labfn)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, u)
	if err != nil {
		return err
	}
	return nil
}
