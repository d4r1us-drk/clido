package cmd

import (
	"strconv"

	"github.com/d4r1us-drk/clido/pkg/repository"
	"github.com/spf13/cobra"
)

// NewRemoveCmd creates and returns the remove command.
func NewRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [project|task] <id>",
		Short: "Remove a project or task along with all its subprojects or subtasks",
		Long: `Remove an existing project or task identified by its ID. This will also remove all associated subprojects 
  or subtasks.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < MinArgsLength {
				cmd.Println(
					"Insufficient arguments. Use 'remove project <id>' or 'remove task <id>'.",
				)
				return
			}

			repo, err := repository.NewRepository()
			if err != nil {
				cmd.Printf("Error initializing repository: %v\n", err)
				return
			}
			defer repo.Close()

			id, err := strconv.Atoi(args[1])
			if err != nil {
				cmd.Println("Invalid ID. Please provide a numeric ID.")
				return
			}

			switch args[0] {
			case "project":
				removeProject(cmd, repo, id)
			case "task":
				removeTask(cmd, repo, id)
			default:
				cmd.Println("Invalid option. Use 'remove project <id>' or 'remove task <id>'.")
			}
		},
	}

	return cmd
}

func removeProject(cmd *cobra.Command, repo *repository.Repository, id int) {
	// First remove all subprojects
	subprojects, err := repo.GetSubprojects(id)
	if err != nil {
		cmd.Printf("Error retrieving subprojects: %v\n", err)
		return
	}
	for _, subproject := range subprojects {
		removeProject(cmd, repo, subproject.ID)
	}

	// Now remove the parent project
	err = repo.DeleteProject(id)
	if err != nil {
		cmd.Printf("Error removing project: %v\n", err)
		return
	}

	cmd.Printf("Project (ID: %d) and all its subprojects removed successfully.\n", id)
}

func removeTask(cmd *cobra.Command, repo *repository.Repository, id int) {
	// First remove all subtasks
	subtasks, err := repo.GetSubtasks(id)
	if err != nil {
		cmd.Printf("Error retrieving subtasks: %v\n", err)
		return
	}
	for _, subtask := range subtasks {
		removeTask(cmd, repo, subtask.ID)
	}

	// Now remove the parent task
	err = repo.DeleteTask(id)
	if err != nil {
		cmd.Printf("Error removing task: %v\n", err)
		return
	}

	cmd.Printf("Task (ID: %d) and all its subtasks removed successfully.\n", id)
}
