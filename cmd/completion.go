package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// NewCompletionCmd creates and returns the 'completion' command,
// which generates shell completion scripts for Bash, Zsh, Fish, and PowerShell.
//
// The command allows users to generate and load completion scripts for their preferred shell.
// Completion scripts help users auto-complete command-line inputs for 'clido'.
//
// This command supports the following shells:
// - Bash
// - Zsh
// - Fish
// - PowerShell
func NewCompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",  // Defines the valid subcommands for shell types
		Short: "Generate completion script",  // Brief description of the command
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
`, // Detailed usage instructions for each shell

		DisableFlagsInUseLine: true,  // Disables flag usage display in the command usage line
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},  // Specifies valid arguments for shell types
		Args: func(cmd *cobra.Command, args []string) error {
			// Ensures exactly one argument (shell type) is provided
			if len(args) != 1 {
				cmd.PrintErrln(
					"Error: requires exactly one argument: bash, zsh, fish, or powershell",
				)
				return cobra.NoArgs(cmd, args)  // Returns an error if no arguments are provided
			}
			return cobra.OnlyValidArgs(cmd, args)  // Validates the argument
		},

		Run: func(cmd *cobra.Command, args []string) {
			// Switch case to handle shell type provided as argument
			switch args[0] {
			case "bash":
				// Generate Bash completion script and output it to stdout
				if err := cmd.Root().GenBashCompletion(os.Stdout); err != nil {
					cmd.PrintErrf("Error generating bash completion: %v\n", err)
					os.Exit(1)  // Exit with error code 1 if there is a failure
				}

			case "zsh":
				// Generate Zsh completion script and output it to stdout
				if err := cmd.Root().GenZshCompletion(os.Stdout); err != nil {
					cmd.PrintErrf("Error generating zsh completion: %v\n", err)
					os.Exit(1)  // Exit with error code 1 if there is a failure
				}

			case "fish":
				// Generate Fish completion script and output it to stdout
				if err := cmd.Root().GenFishCompletion(os.Stdout, true); err != nil {
					cmd.PrintErrf("Error generating fish completion: %v\n", err)
					os.Exit(1)  // Exit with error code 1 if there is a failure
				}

			case "powershell":
				// Generate PowerShell completion script and output it to stdout
				if err := cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout); err != nil {
					cmd.PrintErrf("Error generating PowerShell completion: %v\n", err)
					os.Exit(1)  // Exit with error code 1 if there is a failure
				}
			}
		},
	}
}
