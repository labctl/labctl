package app

// import (
// 	"fmt"

// 	"github.com/labctl/labctl/helpers"
// 	"github.com/labctl/labctl/utils"
// 	"github.com/spf13/cobra"
// )

// // Filename for the output, else to stdout
// var outputFile string

// const (
// 	input    = "input"
// 	output   = "output"
// 	template = "template"
// )

// // configCmd represents the config command
// var discoverCmd = &cobra.Command{
// 	Use:          "discover",
// 	Short:        "discover lab variables",
// 	Long:         "discover lab variables from running nodes\nreference: https://labctl.github.io/labctl/cmd/discover/",
// 	Aliases:      []string{"d"},
// 	SilenceUsage: true,
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		if len(args) > 0 {
// 			return fmt.Errorf("unexpected arguments: %s", args)
// 		}

// 		res, err := helpers.DiscoverTopo("interfaces")
// 		if err != nil {
// 			return err
// 		}

// 		return utils.WriteYAML(outputFile, res)
// 	}}

// var discoverFileCmd = &cobra.Command{
// 	Use:          "file",
// 	Short:        "discover lab variables",
// 	Long:         "discover lab variables from a file containing CLI output\nreference: https://labctl.github.io/labctl/cmd/discover/",
// 	Aliases:      []string{"f"},
// 	SilenceUsage: true,
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		if len(args) > 0 {
// 			return fmt.Errorf("unexpected arguments: %s", args)
// 		}
// 		file_in := cmd.Flag(input).Value.String()
// 		file_tmpl := cmd.Flag(template).Value.String()

// 		res, err := helpers.DiscoverFile(file_tmpl, file_in)
// 		if err != nil {
// 			return err
// 		}

// 		return utils.WriteYAML(outputFile, res)
// 	},
// }

// func init() {
// 	rootCmd.AddCommand(discoverCmd)
// 	discoverCmd.Flags().SortFlags = false
// 	discoverCmd.Flags().StringSliceVarP(&helpers.NodeFilter, "filter", "f", []string{}, "comma separated list of nodes to include")
// 	discoverCmd.Flags().StringVarP(&outputFile, output, "o", "", "write result to this file")
// 	_ = discoverCmd.MarkFlagFilename(output)

// 	discoverCmd.AddCommand(discoverFileCmd)
// 	discoverFileCmd.Flags().SortFlags = false
// 	discoverFileCmd.Flags().StringP(input, "i", "", "read cli from this file")
// 	_ = discoverFileCmd.MarkFlagFilename(input, "*.txt", "*.log")
// 	_ = discoverFileCmd.MarkFlagRequired(input)
// 	discoverFileCmd.Flags().String(template, "", "file with the discovery template")
// 	_ = discoverFileCmd.MarkFlagFilename(template, "*.yaml")
// 	_ = discoverFileCmd.MarkFlagRequired(template)
// 	discoverFileCmd.Flags().StringVarP(&outputFile, output, "o", "", "write result to this file")
// 	_ = discoverFileCmd.MarkFlagFilename(output)
// }
