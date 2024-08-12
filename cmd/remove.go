package cmd

import (
	"fmt"
	"strconv"

	"github.com/d4r1us-drk/clido/pkg/repository"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [project|task] <id>",
	Short: "Remove a project or task",
	Long:  `Remove an existing project or task identified by its ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Insufficient arguments. Use 'remove project <id>' or 'remove task <id>'.")
			return
		}

		repo, err := repository.NewRepository()
		if err != nil {
			fmt.Printf("Error initializing repository: %v\n", err)
			return
		}
		defer repo.Close()

		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid ID. Please provide a numeric ID.")
			return
		}

		switch args[0] {
		case "project":
			removeProject(repo, id)
		case "task":
			removeTask(repo, id)
		default:
			fmt.Println("Invalid option. Use 'remove project <id>' or 'remove task <id>'.")
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func removeProject(repo *repository.Repository, id int) {
	project, err := repo.GetProjectByID(id)
	if err != nil {
		fmt.Printf("Error retrieving project: %v\n", err)
		return
	}

	err = repo.DeleteProject(id)
	if err != nil {
		fmt.Printf("Error removing project: %v\n", err)
		return
	}

	fmt.Printf("Project '%s' (ID: %d) removed successfully.\n", project.Name, id)
}

func removeTask(repo *repository.Repository, id int) {
	task, err := repo.GetTaskByID(id)
	if err != nil {
		fmt.Printf("Error retrieving task: %v\n", err)
		return
	}

	err = repo.DeleteTask(id)
	if err != nil {
		fmt.Printf("Error removing task: %v\n", err)
		return
	}

	fmt.Printf("Task '%s' (ID: %d) removed successfully.\n", task.Name, id)
}
