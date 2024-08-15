package cmd

import (
	"fmt"
	"strconv"

	"github.com/d4r1us-drk/clido/pkg/repository"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [project|task] <id>",
	Short: "Remove a project or task along with all its subprojects or subtasks",
	Long:  `Remove an existing project or task identified by its ID. This will also remove all associated subprojects or subtasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < MinArgsLength {
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
	// First remove all subprojects
	subprojects, err := repo.GetSubprojects(id)
	if err != nil {
		fmt.Printf("Error retrieving subprojects: %v\n", err)
		return
	}
	for _, subproject := range subprojects {
		removeProject(repo, subproject.ID)
	}

	// Now remove the parent project
	err = repo.DeleteProject(id)
	if err != nil {
		fmt.Printf("Error removing project: %v\n", err)
		return
	}

	fmt.Printf("Project (ID: %d) and all its subprojects removed successfully.\n", id)
}

func removeTask(repo *repository.Repository, id int) {
	// First remove all subtasks
	subtasks, err := repo.GetSubtasks(id)
	if err != nil {
		fmt.Printf("Error retrieving subtasks: %v\n", err)
		return
	}
	for _, subtask := range subtasks {
		removeTask(repo, subtask.ID)
	}

	// Now remove the parent task
	err = repo.DeleteTask(id)
	if err != nil {
		fmt.Printf("Error removing task: %v\n", err)
		return
	}

	fmt.Printf("Task (ID: %d) and all its subtasks removed successfully.\n", id)
}
