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
	if err := ctx.Settings.InitColors(); err != nil {
		return err
	}
	if ctx.DebugCount > 0 {
		for _, r := range ctx.Settings.Colors {
			fmt.Printf("- %s %s, %s\n", r.Replace, r.ColorStr, r.Regexp.String())
		}
	}
	return colorize.Spawn(ctx.Settings.Colors, r.Cmd[0], r.Cmd[1:]...)
}
