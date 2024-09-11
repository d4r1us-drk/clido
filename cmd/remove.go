package cmd

import (
	"strconv"

	"github.com/d4r1us-drk/clido/pkg/repository"
	"github.com/spf13/cobra"
)

// NewRemoveCmd creates and returns the 'remove' command for deleting projects or tasks.
//
// The command allows users to remove an existing project or task by its ID. 
// It also ensures that all associated subprojects or subtasks are recursively removed.
//
// Usage:
//   clido remove project <id>    # Remove a project and its subprojects
//   clido remove task <id>       # Remove a task and its subtasks
func NewRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [project|task] <id>",  // Specifies valid options: 'project' or 'task' followed by an ID
		Short: "Remove a project or task along with all its subprojects or subtasks",  // Short description of the command
		Long: `Remove an existing project or task identified by its ID. This will also remove all associated subprojects 
  or subtasks.`,  // Extended description with clarification on subprojects and subtasks being removed recursively
		Run: func(cmd *cobra.Command, args []string) {
			// Ensure sufficient arguments (either 'project' or 'task' followed by an ID)
			if len(args) < MinArgsLength {
				cmd.Println(
					"Insufficient arguments. Use 'remove project <id>' or 'remove task <id>'.",
				)
				return
			}

			// Initialize the repository for database operations
			repo, err := repository.NewRepository()
			if err != nil {
				cmd.Printf("Error initializing repository: %v\n", err)
				return
			}
			defer repo.Close()

			// Parse the ID argument into an integer
			id, err := strconv.Atoi(args[1])
			if err != nil {
				cmd.Println("Invalid ID. Please provide a numeric ID.")
				return
			}

			// Determine whether the user wants to remove a project or a task
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

// removeProject handles the recursive removal of a project and all its subprojects.
// It first retrieves and removes all subprojects, and then deletes the parent project.
func removeProject(cmd *cobra.Command, repo *repository.Repository, id int) {
	// Retrieve all subprojects associated with the project
	subprojects, err := repo.GetSubprojects(id)
	if err != nil {
		cmd.Printf("Error retrieving subprojects: %v\n", err)
		return
	}

	// Recursively remove all subprojects
	for _, subproject := range subprojects {
		removeProject(cmd, repo, subproject.ID)
	}

	// Remove the parent project after all subprojects have been removed
	err = repo.DeleteProject(id)
	if err != nil {
		cmd.Printf("Error removing project: %v\n", err)
		return
	}

	cmd.Printf("Project (ID: %d) and all its subprojects removed successfully.\n", id)
}

// removeTask handles the recursive removal of a task and all its subtasks.
// It first retrieves and removes all subtasks, and then deletes the parent task.
func removeTask(cmd *cobra.Command, repo *repository.Repository, id int) {
	// Retrieve all subtasks associated with the task
	subtasks, err := repo.GetSubtasks(id)
	if err != nil {
		cmd.Printf("Error retrieving subtasks: %v\n", err)
		return
	}

	// Recursively remove all subtasks
	for _, subtask := range subtasks {
		removeTask(cmd, repo, subtask.ID)
	}

	// Remove the parent task after all subtasks have been removed
	err = repo.DeleteTask(id)
	if err != nil {
		cmd.Printf("Error removing task: %v\n", err)
		return
	}

	cmd.Printf("Task (ID: %d) and all its subtasks removed successfully.\n", id)
}
