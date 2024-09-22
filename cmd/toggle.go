package cmd

import (
	"strconv"

	"github.com/d4r1us-drk/clido/controllers"
	"github.com/spf13/cobra"
)

// NewToggleCmd creates and returns the 'toggle' command for marking tasks as completed or uncompleted.
func NewToggleCmd(taskController *controllers.TaskController) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "toggle <task_id>",
		Short: "Toggle task completion status",
		Long:  "Toggle the completion status of a task identified by its ID.",
		Run: func(cmd *cobra.Command, args []string) {
			// Ensure at least one argument (the task ID) is provided
			if len(args) < 1 {
				cmd.Println("Insufficient arguments. Use 'toggle <task_id>'.")
				return
			}

			// Parse the task ID argument into an integer
			id, err := strconv.Atoi(args[0])
			if err != nil {
				cmd.Println("Invalid task ID. Please provide a numeric ID.")
				return
			}

			// Check if the recursive flag was provided
			recursive, _ := cmd.Flags().GetBool("recursive")

			// Toggle task completion status using the controller
			completionStatus, toggleErr := taskController.ToggleTaskCompletion(id, recursive)
			if toggleErr != nil {
				cmd.Printf("Error toggling task: %v\n", toggleErr)
				return
			}

			if recursive {
				cmd.Printf(
					"Task (ID: %d) and its subtasks (if any) have been set as %s.\n",
					id,
					completionStatus,
				)
			} else {
				cmd.Printf("Task (ID: %d) has been set as %s.\n", id, completionStatus)
			}
		},
	}

	// Add flag for recursive toggle, allowing users to recursively toggle all subtasks
	cmd.Flags().BoolP("recursive", "r", false, "Recursively toggle subtasks")

	return cmd
}
