package helpers

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/imdario/mergo"
	"github.com/labctl/labctl/utils"
	"github.com/labctl/labctl/utils/colorize"
	"gopkg.in/yaml.v3"

	log "github.com/sirupsen/logrus"
)

type Settings struct {
	Colors []*colorize.Colorize `yaml:"colorize"`
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
	p := utils.Path{Path: "~/set.yml"}
	_ = p.ExpandUser()
	setByte, err := os.ReadFile(p.Path)
	if err != nil {
		if silent && errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return s.unmarshal(setByte)
}

func (s *Settings) unmarshal(b []byte) error {
	err := yaml.Unmarshal(b, &s)
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
