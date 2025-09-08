package helpers

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"strings"

	"dario.cat/mergo"
	"github.com/labctl/labctl/utils"
	"github.com/labctl/labctl/utils/colorize"
	"github.com/labctl/labctl/utils/tx"
	"gopkg.in/yaml.v3"

	"github.com/charmbracelet/log"
)

type Settings struct {
	Colors []*colorize.Colorize `yaml:"colorize"`

	Kindmap map[string]tx.KindDef `yaml:"kindmap"`
}

//go:embed settings.yml
var defaultSettings embed.FS

func NewSettings(defaults ...bool) *Settings {
	s := &Settings{}
	if len(defaults) > 0 {
		in, err := defaultSettings.ReadFile("settings.yml")
		if err != nil {
			log.Fatal(err)
		}
		err = s.unmarshal(in)
		if err != nil {
			log.Errorf("built in settings:\n%s", string(in))
			log.Fatal(err)
		}
	}
	return s
}

func (s *Settings) AddSettings(path string, silent bool) error {
	new := NewSettings()
	err := new.load(path, silent)
	if err == nil {
		err = mergo.MergeWithOverwrite(s, new)
	}
	return err
}

func (s *Settings) load(path string, silent bool) error {
	p := utils.NewPathExpandUser(path)
	exist, err := p.Exists()
	if !exist {
		return nil
	}
	if err != nil {
		log.Error("could not load file", "name", path, "err", err)
		return err
	}
	setByte, err := p.ReadFile()
	if err != nil {
		if silent && errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return s.unmarshal(setByte)
}

func (s *Settings) unmarshal(bb []byte) error {
	err := yaml.Unmarshal(bb, &s)
	if err != nil {
		if strings.Contains(err.Error(), "unknown escape") {
			err = fmt.Errorf("%s (regex requires single quote escape characters '')", err)
		}
		return err
	}
	return nil
}

func (s *Settings) InitColors() error {
	for _, c := range s.Colors {
		err := c.Init()
		if err != nil {
			return err
		}

	}
	log.Debugf("colorize: %v rules\n%v\n", len(s.Colors), s.Colors)
	return nil
}
