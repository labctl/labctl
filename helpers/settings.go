package helpers

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/labctl/labctl/utils"
	"github.com/labctl/labctl/utils/colorize"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

type Settings struct {
	Colors []*colorize.Colorize `yaml:"colorize"`
}

func (s *Settings) Load() error {
	path := utils.Path{Path: "~/set.yml"}
	_ = path.ExpandUser()
	setByte, err := ioutil.ReadFile(path.Path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(setByte, &s)
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
	log.Debugf("%v rules%v\n\n", len(s.Colors), s.Colors)
	return nil
}
