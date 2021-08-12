package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "generate completion script",
	Long: `To load completions:

Bash:

  $ source <(labctl completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ labctl completion bash > /etc/bash_completion.d/labctl
  # macOS:
  $ labctl completion bash > /usr/local/etc/bash_completion.d/labctl

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ labctl completion zsh > "${fpath[1]}/_labctl"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ labctl completion fish | source

  # To load completions for each session, execute once:
  $ labctl completion fish > ~/.config/fish/completions/labctl.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			_ = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			_ = cmd.Root().GenFishCompletion(os.Stdout, true)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
