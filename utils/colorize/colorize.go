package colorize

import (
	"fmt"
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
	Regexp   *Regexp `yaml:"regex"`
	ColorStr string  `yaml:"color"`
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

// Converts a comma separated list of colors into color.Attributes
func StringToAttributes(colorStr string) ([]color.Attribute, error) {
	slst := strings.Split(colorStr, ",")
	var attr []color.Attribute
	for _, s := range slst {
		col, err := stringToAttribute(s)
		if err != nil {
			return nil, err
		}
		attr = append(attr, col)
	}
	return attr, nil
}

// Convert a single color string into color.Attribute
// Used by StringToAttributes
func stringToAttribute(colorString string) (color.Attribute, error) {
	colorStr := strings.Trim(strings.ToLower(colorString), " \t")

	// Search through styles
	var styles = []string{"reset", "bold", "faint", "italic", "blinkslow", "blinkrapid", "reversevideo", "concealed", "crossedout"}
	for i, s := range styles {
		if s == colorStr {
			return color.Attribute(i), nil
		}
	}

	// Strip fg, bg, hi
	var BgHi func(string) (bool, bool, string)
	BgHi = func(s string) (bool, bool, string) {
		if strings.HasPrefix(s, "fg") {
			return BgHi(s[2:])
		}
		if strings.HasPrefix(s, "bg") {
			_, hi, res := BgHi(s[2:])
			return true, hi, res
		}
		if strings.HasPrefix(s, "hi") {
			bg, _, res := BgHi(s[2:])
			return bg, true, res
		}
		return false, false, s
	}

	// starts at: "fg": 30, "bg": 40, "fghi": 90, "bghi": 100
	bg, hi, cres := BgHi(colorStr)
	offset := color.FgBlack
	if bg { // background +10
		offset += color.BgBlack - color.FgBlack
	}
	if hi { // highlight +60
		offset += color.FgHiBlack - color.FgBlack
	}

	// search through colors
	var colors = []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"}
	for i, c := range colors {
		if c == cres {
			return offset + color.Attribute(i), nil
		}
	}
	return 0, fmt.Errorf("invalid color: %s. Valid colors: %s, valid styles: %s prefix with: fg,bg,hi", colorString, strings.Join(colors, ","), strings.Join(styles, ","))
}
