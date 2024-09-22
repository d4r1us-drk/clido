package cmd

import (
	"os"

	"github.com/d4r1us-drk/clido/controllers"
	"github.com/d4r1us-drk/clido/internal/version"
	"github.com/spf13/cobra"
)

// Constants for table printing. These are not user constraints.
const (
	MaxProjectNameLength     = 30 // Maximum length for project names
	MaxProjectDescLength     = 50 // Maximum length for project descriptions
	MaxTaskNameLength        = 20 // Maximum length for task names
	MaxTaskDescLength        = 30 // Maximum length for task descriptions
	MaxProjectNameWrapLength = 20 // Maximum length for wrapping project names in the UI
	MinArgsLength            = 2  // Minimum required arguments for certain commands
	PriorityHigh             = 1
	PriorityMedium           = 2
	PriorityLow              = 3
	PriorityNone             = 4
	PriorityEmpty            = 0
)

// NewRootCmd creates and returns the root command for the CLI application.
func NewRootCmd(
	projectController *controllers.ProjectController,
	taskController *controllers.TaskController,
) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "clido",
		Short: "clido is an awesome CLI to-do list management application",
		Long: "clido is a simple yet powerful CLI tool designed to help you manage " +
			"your projects and tasks effectively from the terminal.",
	}

	// Add subcommands and pass the controllers
	rootCmd.AddCommand(NewVersionCmd()) // Version command to display the app version
	rootCmd.AddCommand(NewCompletionCmd())
	rootCmd.AddCommand(NewNewCmd(projectController, taskController))
	rootCmd.AddCommand(NewEditCmd(projectController, taskController))
	rootCmd.AddCommand(NewListCmd(projectController, taskController))
	rootCmd.AddCommand(NewRemoveCmd(projectController, taskController))
	rootCmd.AddCommand(NewToggleCmd(taskController))

	return rootCmd
}

// NewVersionCmd creates the version command.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Clido",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Println(version.FullVersion())
		},
	}
}

// Execute runs the root command.
func Execute() {
	rootCmd := NewRootCmd(nil, nil)
	if err := rootCmd.Execute(); err != nil {
		rootCmd.Println(err)
		os.Exit(1)
	}
}
