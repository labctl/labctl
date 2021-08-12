package cmd

import (
	"fmt"
	"strings"

	"github.com/labctl/labctl/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const unexpected = "unexpected command: %v"

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:          "config",
	Short:        "configure a lab",
	Long:         "configure a lab based on templates and variables from the topology definition file\nreference: https://labctl.github.io/labctl/cmd/config/",
	Aliases:      []string{"conf"},
	ValidArgs:    []string{"commit", "send", "compare", "template"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf(unexpected, args)
		}
		if len(args) == 0 {
			if len(helpers.TemplateNames) == 1 && strings.HasPrefix(helpers.TemplateNames[0], "show-") {
				log.Debugf("Template name start with show- so default command will be 'send'")
				return helpers.ConfigRun("send")
			}
			return fmt.Errorf("command required: commit, send, compare, template")
		}
		return helpers.ConfigRun(args[0])
	},
}

var configSendCmd = &cobra.Command{
	Use:          "send",
	Short:        "send raw commands to a lab",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf(unexpected, args)
		}
		return helpers.ConfigRun(helpers.ActionSend)
	},
}

var configCompareCmd = &cobra.Command{
	Use:          "compare",
	Short:        "compare configuration to a running lab",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf(unexpected, args)
		}
		return helpers.ConfigRun(helpers.ActionCompare)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringSliceVarP(&helpers.TemplatePaths, "template-path", "p", []string{}, "comma separated list of paths to search for templates")
	_ = configCmd.MarkFlagDirname("template-path")
	configCmd.Flags().StringSliceVarP(&helpers.TemplateNames, "template-list", "l", []string{}, "comma separated list of template names to render")
	configCmd.Flags().StringSliceVarP(&helpers.NodeFilter, "filter", "f", []string{}, "comma separated list of nodes to include")
	configCmd.Flags().SortFlags = false

	configCmd.AddCommand(configSendCmd)
	configSendCmd.Flags().AddFlagSet(configCmd.Flags())

	configCmd.AddCommand(configCompareCmd)
	configCompareCmd.Flags().AddFlagSet(configCmd.Flags())
}
