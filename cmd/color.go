package cmd

import (
	"fmt"

	"github.com/labctl/labctl/helpers"
	"github.com/labctl/labctl/utils/colorize"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

// completionCmd represents the completion command
var colorCmd = &cobra.Command{
	Use:                   "color",
	Short:                 "add color to",
	Aliases:               []string{"c"},
	Long:                  ``,
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) < 1 {
			return fmt.Errorf("no argument supplied")
		}
		mycmd := args[0]
		args = args[1:]

		s := helpers.Settings{}
		var err error

		if err=s.Load(); err!= nil {
			return err
		}
		if err=s.InitColors(); err!= nil {
			return err
		}
		log.Debugf("%v rules%v\n\n", len(s.Colors), s.Colors)
		for _, r := range s.Colors {
			fmt.Printf("- %s %s, %s\n", r.Replace, r.ColorStr, r.Regexp.String())
		}

		return colorize.Spawn(s.Colors, mycmd, args...)
	},
}

func init() {
	rootCmd.AddCommand(colorCmd)
}
