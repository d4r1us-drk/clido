package cmd

import (
	"strconv"
	"time"

	"github.com/d4r1us-drk/clido/pkg/repository"
	"github.com/spf13/cobra"
)

// NewToggleCmd creates and returns the 'toggle' command for marking tasks as completed or uncompleted.
//
// The command allows users to toggle the completion status of a task by its ID. 
// If the task is marked as completed, it updates the completion date.
// The command also supports recursively toggling the status of all subtasks.
//
// Usage:
//   clido toggle <task_id>    # Toggle the completion status of a task
//   clido toggle <task_id> -r # Recursively toggle the completion status of the task and its subtasks
func NewToggleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "toggle <task_id>",  // Specifies the required task ID as the argument
		Short: "Toggle task completion status",  // Short description of the command
		Long:  `Toggle the completion status of a task identified by its ID.`,  // Extended description
		Run: func(cmd *cobra.Command, args []string) {
			// Ensure at least one argument (the task ID) is provided
			if len(args) < 1 {
				cmd.Println("Insufficient arguments. Use 'toggle <task_id>'.")
				return
			}

			// Initialize the repository for database operations
			repo, err := repository.NewRepository()
			if err != nil {
				cmd.Printf("Error initializing repository: %v\n", err)
				return
			}
			defer repo.Close()

			// Parse the task ID argument into an integer
			id, err := strconv.Atoi(args[0])
			if err != nil {
				cmd.Println("Invalid task ID. Please provide a numeric ID.")
				return
			}

			// Check if the recursive flag was provided
			recursive, _ := cmd.Flags().GetBool("recursive")
			toggleTask(cmd, repo, id, recursive)
		},
	}

	// Add flag for recursive toggle, allowing users to recursively toggle all subtasks
	cmd.Flags().BoolP("recursive", "r", false, "Recursively toggle subtasks")

	return cmd
}

// toggleTask toggles the completion status of a task and optionally its subtasks.
//
// If the task is currently marked as incomplete, it will be marked as completed, 
// and the completion date will be updated. If the task is marked as completed, 
// it will be toggled to incomplete and the completion date will be cleared.
//
// If the recursive flag is set to true, this function will also toggle the completion status of all subtasks.
func toggleTask(cmd *cobra.Command, repo *repository.Repository, id int, recursive bool) {
	// Retrieve the task by its ID
	task, err := repo.GetTaskByID(id)
	if err != nil {
		cmd.Printf("Error retrieving task: %v\n", err)
		return
	}

	// Toggle the task's completion status
	task.TaskCompleted = !task.TaskCompleted
	task.LastUpdatedDate = time.Now()

	// Set the completion date if the task is marked as completed, or clear it if uncompleted
	if task.TaskCompleted {
		task.CompletionDate = &task.LastUpdatedDate
	} else {
		task.CompletionDate = nil
	}

	// Update the task in the repository
	err = repo.UpdateTask(task)
	if err != nil {
		cmd.Printf("Error updating task: %v\n", err)
		return
	}

	// Display the updated status of the task
	status := "completed"
	if !task.TaskCompleted {
		status = "uncompleted"
	}
	cmd.Printf("Task '%s' (ID: %d) marked as %s.\n", task.Name, id, status)

	// If the recursive flag is set, retrieve and toggle the completion status of all subtasks
	if recursive {
		subtasks, _ := repo.GetSubtasks(id)
		for _, subtask := range subtasks {
			toggleTask(cmd, repo, subtask.ID, recursive)  // Recursively toggle subtasks
		}
	}
}
