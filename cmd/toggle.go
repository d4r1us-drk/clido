package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/d4r1us-drk/clido/pkg/repository"
	"github.com/spf13/cobra"
)

var toggleCmd = &cobra.Command{
	Use:   "toggle <task_id>",
	Short: "Toggle task completion status",
	Long:  `Toggle the completion status of a task identified by its ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Insufficient arguments. Use 'toggle <task_id>'.")
			return
		}

		repo, err := repository.NewRepository()
		if err != nil {
			fmt.Printf("Error initializing repository: %v\n", err)
			return
		}
		defer repo.Close()

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid task ID. Please provide a numeric ID.")
			return
		}

		toggleTask(repo, id)
	},
}

func init() {
	rootCmd.AddCommand(toggleCmd)
}

func toggleTask(repo *repository.Repository, id int) {
	task, err := repo.GetTaskByID(id)
	if err != nil {
		fmt.Printf("Error retrieving task: %v\n", err)
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
		fmt.Printf("Error updating task: %v\n", err)
		return
	}

	status := "completed"
	if !task.TaskCompleted {
		status = "uncompleted"
	}

	fmt.Printf("Task '%s' (ID: %d) marked as %s.\n", task.Name, id, status)
}
