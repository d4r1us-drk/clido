package cmd

import (
	"os"

	"github.com/d4r1us-drk/clido/internal/version"
	"github.com/spf13/cobra"
)

const (
	MaxProjectNameLength     = 30
	MaxProjectDescLength     = 50
	MaxTaskNameLength        = 20
	MaxTaskDescLength        = 30
	MaxProjectNameWrapLength = 20
	MinArgsLength            = 2
)

// NewRootCmd creates and returns the root command.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "clido",
		Short: "Clido is an awesome CLI to-do list management application",
		Long: `Clido is a simple yet powerful CLI tool designed to help you manage 
  your projects and tasks effectively from the terminal.`,
	}

	// Add subcommands to rootCmd here
	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewCompletionCmd())
	rootCmd.AddCommand(NewNewCmd())
	rootCmd.AddCommand(NewEditCmd())
	rootCmd.AddCommand(NewListCmd())
	rootCmd.AddCommand(NewRemoveCmd())
	rootCmd.AddCommand(NewToggleCmd())

	return rootCmd
}

// NewVersionCmd creates and returns the version command.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Clido",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Println(version.FullVersion())
		},
	}
}

// Execute runs the root command, which includes all subcommands.
func Execute() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		rootCmd.Println(err)
		os.Exit(1)
	}
}
