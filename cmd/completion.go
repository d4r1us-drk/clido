package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(clido completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ clido completion bash > /etc/bash_completion.d/clido
  # macOS:
  $ clido completion bash > /usr/local/etc/bash_completion.d/clido

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ clido completion zsh > "${fpath[1]}/_clido"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ clido completion fish | source

  # To load completions for each session, execute once:
  $ clido completion fish > ~/.config/fish/completions/clido.fish

PowerShell:

  PS> clido completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> clido completion powershell > clido.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
