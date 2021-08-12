package cmd

import (
	"fmt"

	"github.com/labctl/labctl/utils"
	"github.com/sirikothe/gotextfsm"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var infile string

// configCmd represents the config command
var discoverCmd = &cobra.Command{
	Use:          "discover",
	Short:        "discover a lab topology",
	Long:         "discover a lab topology from running nodes\nreference: https://labctl.github.io/labctl/cmd/discover/",
	Aliases:      []string{"d"},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("unexpected arguments: %s", args)
		}

		template := `Value Required,Key interface (\w\S+)
Value admin (Up|Down)
Value oper_v4 (\S+)
Value oper_v6 (\S+)
Value mode (\S+)
Value port_sap (\S+)
Value List address ((\d+\.|[a-f\d]+:)\S+)
Value List pfx_state (\S+)


Start
  ^[|\s]*.*?\s+(Up|Down) -> Continue.Record
  ^[|\s]*${interface}\s+${admin}\s+${oper_v4}/${oper_v6}\s+${mode}\s+${port_sap}
  ^[|\s]*${address}\s+${pfx_state}
`
		log.Infof("Template\n%s", template)
		input := utils.LoadFile(infile)
		//log.Infof("Input\n%s", input)

		fsm := gotextfsm.TextFSM{}
		err := fsm.ParseString(template)
		if err != nil {
			return fmt.Errorf("error while parsing template: %s", err.Error())
		}

		parser := gotextfsm.ParserOutput{}

		err = parser.ParseTextString(input, fsm, true)
		if err != nil {
			return fmt.Errorf("error while parsing input '%s'", err.Error())
		}

		utils.LogPretty("result:", parser.Dict)

		// fmt.Printf("Parsed output: %+v\n", parser.Dict)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(discoverCmd)
	discoverCmd.Flags().StringVar(&infile, "file", "", "filename")
	_ = discoverCmd.MarkFlagFilename("file")
}
