package app

import (
	"fmt"

	"github.com/labctl/labctl/helpers"
	"github.com/labctl/labctl/utils/colorize"
)

type CmdColor struct {
	Cmd []string `arg:""`
}

func (r *CmdColor) Run(ctx *helpers.Context) error {
	tst := len(r.Cmd) == 0 || r.Cmd[0] == "--test"
	if tst || ctx.DebugCount > 0 {
		fmt.Println("Colors from settings:")
		for _, r := range ctx.Settings.Colors {
			fmt.Printf("- %s %s, %s\n", r.Replace, r.ColorStr, r.Regexp.String())
		}
	}
	if tst {
		colorize.TestPattern()
		return nil
	}
	return colorize.Spawn(ctx.Settings.Colors, r.Cmd[0], r.Cmd[1:]...)
}
