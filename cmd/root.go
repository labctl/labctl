package cmd

import (
	"os"

	"github.com/labctl/labctl/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ce", // "labctl"
	Short: "control a lab environment",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if helpers.DebugCount > 0 {
			log.SetLevel(log.DebugLevel)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SilenceUsage = true
	rootCmd.PersistentFlags().CountVarP(&helpers.DebugCount, "debug", "d", "enable debug mode")
	rootCmd.PersistentFlags().StringSliceVarP(&helpers.TopoFiles, "topo", "t", []string{}, "path to the file(s) with topology information")
	_ = rootCmd.MarkPersistentFlagFilename("topo", "*.yaml", "*.yml")
}
