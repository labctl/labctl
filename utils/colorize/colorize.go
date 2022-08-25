package colorize

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// a Marshallable Regexp
type Regexp struct {
	*regexp.Regexp
}

// implements the yaml.Unmarshaler interface for Regexp
func (re *Regexp) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	regex, err := regexp.Compile(s)
	if err != nil {
		return err
	}
	re.Regexp = regex
	return nil
}

// implements the yaml.Marshaler interface for Regexp
func (re *Regexp) MarshalYAML() (interface{}, error) {
	return re.Regexp.String(), nil
}

// The colorize struct contains a regex to match and the teminal color to replace to
type Colorize struct {
	ColorStr string  `yaml:"color"`
	Regexp   *Regexp `yaml:"regex"`
	Replace  []byte
	attr     []color.Attribute
}

func NewColorize(re string, colorStr string) (*Colorize, error) {
	res, err := regexp.Compile(re)
	if err != nil {
		return nil, err
	}

	c := &Colorize{
		Regexp:   &Regexp{Regexp: res},
		ColorStr: colorStr,
	}
	return c, nil
}

func (c *Colorize) Color() *color.Color {
	return color.New(c.attr...)
}

func (c *Colorize) Init() error {
	attr, err := StringToAttributes(c.ColorStr)
	if err != nil {
		return err
	}
	format := make([]string, len(attr))
	for i, v := range attr {
		format[i] = strconv.Itoa(int(v))
	}
	c.Replace = []byte("\x1b[" + strings.Join(format, ";") + "m$0\x1b[m")
	return nil
}
