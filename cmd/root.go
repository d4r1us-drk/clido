package cmd

import (
	"os"

	"github.com/d4r1us-drk/clido/internal/version"
	"github.com/spf13/cobra"
)

// Constants for table printing. These are not user constraints
const (
	MaxProjectNameLength     = 30  // Maximum length for project names
	MaxProjectDescLength     = 50  // Maximum length for project descriptions
	MaxTaskNameLength        = 20  // Maximum length for task names
	MaxTaskDescLength        = 30  // Maximum length for task descriptions
	MaxProjectNameWrapLength = 20  // Maximum length for wrapping project names in the UI
	MinArgsLength            = 2   // Minimum required arguments for certain commands
)

// NewRootCmd creates and returns the root command for the CLI application.
// This is the entry point of the application, which is responsible for managing
// various subcommands like version, new, edit, list, remove, and toggle commands.
// 
// Returns a *cobra.Command which acts as the root command for all subcommands.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "clido",
		Short: "Clido is an awesome CLI to-do list management application",
		Long: `Clido is a simple yet powerful CLI tool designed to help you manage 
  your projects and tasks effectively from the terminal.`,
	}

	// Adding subcommands to rootCmd
	rootCmd.AddCommand(NewVersionCmd())  // Version command to display the app version
	rootCmd.AddCommand(NewCompletionCmd())  // Completion command to generate shell autocompletion scripts
	rootCmd.AddCommand(NewNewCmd())  // New command to add a new project or task
	rootCmd.AddCommand(NewEditCmd())  // Edit command to modify an existing project or task
	rootCmd.AddCommand(NewListCmd())  // List command to display projects or tasks
	rootCmd.AddCommand(NewRemoveCmd())  // Remove command to delete a project or task
	rootCmd.AddCommand(NewToggleCmd())  // Toggle command to change the status of a task

	return rootCmd
}

// NewVersionCmd creates and returns the version command.
// This command prints the current version of clido, using the version package.
// It helps users check which version of the tool they are running.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Clido",
		Run: func(cmd *cobra.Command, _ []string) {
			// Print the full version of the application, stored in the version package.
			cmd.Println(version.FullVersion())
		},
	}
}

// Execute runs the root command of the application, which triggers
// the appropriate subcommand based on user input.
//
// If an error occurs during execution (such as invalid command usage),
// the application prints the error message and exits with a non-zero status code.
func Execute() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		rootCmd.Println(err)
		os.Exit(1)  // Exit the application with an error status code if the command fails
	}
}
