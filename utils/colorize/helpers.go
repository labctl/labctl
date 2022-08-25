package colorize

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

func center(s string, w int) string {
	return fmt.Sprintf("%*s", -w, fmt.Sprintf("%*s", (w+len(s))/2, s))
}

func ColorStr(col string, text string) string {
	cc := Colorize{ColorStr: col}
	err := cc.Init()
	if err != nil {
		log.Error(err)
		return fmt.Sprintf("%12s", "")
	} else {
		return strings.Replace(string(cc.Replace), "$0", text, 1)
	}
}

func TestPattern() {
	sep := ColorStr("faint", strings.Repeat("-", 12*8))
	fmt.Printf("\n%s\nColors & prefixes\n%s\n", sep, sep)
	for _, pre := range []string{"", "Fg", "Hi", "FgHi", "Bg", "BgHi"} {
		for _, c := range colors {
			fmt.Print(ColorStr(pre+c, center(pre+c, 12)))
		}
		fmt.Print("\n")
	}
	fmt.Printf("\nStyles\n%s\n", sep)
	for _, s := range styles {
		if s == "reset" {
			continue
		}
		fmt.Print(ColorStr(s, center(s, 12)))
	}
	fmt.Printf("\n\nCombinations (you can combine multiple attributes in a list)\n%s\n", sep)
	for _, pre := range []string{"red,bg", "cyan,bg", "bgred,", "bgcyan,", "bold,", "italic,", "blinkslow,"} {
		for _, c := range colors {
			fmt.Print(ColorStr(pre+c, center(pre+c, 13)))
		}
		fmt.Print("\n")
	}
	fmt.Printf("%s\n\n", sep)
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

var (
	styles = []string{"reset", "bold", "faint", "italic", "blinkslow", "blinkrapid", "reversevideo", "concealed", "crossedout"}
	colors = []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white"}
)

// Convert a single color string into color.Attribute
// Used by StringToAttributes
func stringToAttribute(colorString string) (color.Attribute, error) {
	colorStr := strings.Trim(strings.ToLower(colorString), " \t")

	// Search through styles
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
	for i, c := range colors {
		if c == cres {
			return offset + color.Attribute(i), nil
		}
	}
	return 0, fmt.Errorf("invalid color: %s. Valid colors: %s, valid styles: %s prefix with: fg,bg,hi",
		colorString, strings.Join(colors, ","), strings.Join(styles, ","))
}
