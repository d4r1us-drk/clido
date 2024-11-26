package cmd

import (
	"errors"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			// Ensure at least one argument (the task ID) is provided
			if len(args) < 1 {
				return errors.New("insufficient arguments. Use 'toggle <task_id>'")
			}

			// Parse the task ID argument into an integer
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.New("invalid task ID. Please provide a numeric ID")
			}

			// Check if the recursive flag was provided
			recursive, _ := cmd.Flags().GetBool("recursive")

			// Toggle task completion status using the controller
			completionStatus, toggleErr := taskController.ToggleTaskCompletion(id, recursive)
			if toggleErr != nil {
				return errors.New("error toggling task: " + toggleErr.Error())
			}

			// Print the result based on the recursive flag
			if recursive {
				cmd.Println(
					"Task (ID: " + strconv.Itoa(
						id,
					) + ") and its subtasks (if any) have been set as " + completionStatus + ".",
				)
			} else {
				cmd.Println("Task (ID: " + strconv.Itoa(id) + ") has been set as " + completionStatus + ".")
			}

			return nil
		},
	}

	// Add flag for recursive toggle, allowing users to recursively toggle all subtasks
	cmd.Flags().BoolP("recursive", "r", false, "Recursively toggle subtasks")

	return cmd
}
