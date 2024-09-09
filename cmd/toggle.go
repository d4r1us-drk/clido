package cmd

import (
	"strconv"
	"time"

	"github.com/d4r1us-drk/clido/pkg/repository"
	"github.com/spf13/cobra"
)

// NewToggleCmd creates and returns the toggle command.
func NewToggleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "toggle <task_id>",
		Short: "Toggle task completion status",
		Long:  `Toggle the completion status of a task identified by its ID.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Println("Insufficient arguments. Use 'toggle <task_id>'.")
				return
			}

			repo, err := repository.NewRepository()
			if err != nil {
				cmd.Printf("Error initializing repository: %v\n", err)
				return
			}
			defer repo.Close()

			id, err := strconv.Atoi(args[0])
			if err != nil {
				cmd.Println("Invalid task ID. Please provide a numeric ID.")
				return
			}

			recursive, _ := cmd.Flags().GetBool("recursive")
			toggleTask(cmd, repo, id, recursive)
		},
	}

	// Add flag for recursive toggle
	cmd.Flags().BoolP("recursive", "r", false, "Recursively toggle subtasks")

	return cmd
}

func toggleTask(cmd *cobra.Command, repo *repository.Repository, id int, recursive bool) {
	task, err := repo.GetTaskByID(id)
	if err != nil {
		cmd.Printf("Error retrieving task: %v\n", err)
		return
	}

	task.TaskCompleted = !task.TaskCompleted
	task.LastUpdatedDate = time.Now()

	if task.TaskCompleted {
		task.CompletionDate = &task.LastUpdatedDate
	} else {
		task.CompletionDate = nil
	}

	err = repo.UpdateTask(task)
	if err != nil {
		cmd.Printf("Error updating task: %v\n", err)
		return
	}

	status := "completed"
	if !task.TaskCompleted {
		status = "uncompleted"
	}

	cmd.Printf("Task '%s' (ID: %d) marked as %s.\n", task.Name, id, status)

	if recursive {
		subtasks, _ := repo.GetSubtasks(id)
		for _, subtask := range subtasks {
			toggleTask(cmd, repo, subtask.ID, recursive)
		}
	}
}
