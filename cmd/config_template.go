package cmd

import (
	"github.com/labctl/labctl/helpers"
	"github.com/spf13/cobra"
	"github.com/srl-labs/containerlab/clab/config"

	log "github.com/sirupsen/logrus"
)

// Show the template variable and exit
var templateVarOnly, templateVar bool

// configCmd represents the config command
var configTemplateCmd = &cobra.Command{
	Use:          "template",
	Short:        "render a template",
	Long:         "render a template based on variables from the topology definition file\nreference: https://labctl.github.io/cmd/config/template",
	Aliases:      []string{"conf"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		allConfig, err := helpers.LoadAndPrep()
		if err != nil {
			return err
		}

		if templateVarOnly {
			for _, n := range helpers.NodeFilter {
				allConfig[n].Print(true, false)
			}
			log.Infof("aaz")
			return nil
		}

		config.TemplateNames = helpers.TemplateNames
		config.TemplatePaths = helpers.TemplatePaths
		err = config.RenderAll(allConfig)
		if err != nil {
			return err
		}

		for _, n := range helpers.NodeFilter {
			allConfig[n].Print(templateVar, true)
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configTemplateCmd)
	configTemplateCmd.Flags().AddFlagSet(configCmd.Flags())
	configTemplateCmd.Flags().BoolVarP(&templateVar, "vars", "v", false, "show variables used for template rendering")
	configTemplateCmd.Flags().BoolVarP(&templateVarOnly, "onlyvars", "o", false, "only show variables used for template rendering and exit")
	configTemplateCmd.Flags().SortFlags = false
}
