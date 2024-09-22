package cmd

import (
	"strconv"

	"github.com/d4r1us-drk/clido/controllers"
	"github.com/spf13/cobra"
)

// NewRemoveCmd creates and returns the 'remove' command for deleting projects or tasks.
func NewRemoveCmd(
	projectController *controllers.ProjectController,
	taskController *controllers.TaskController,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [project|task] <id>",
		Short: "Remove a project or task along with all its subprojects or subtasks",
		Long:  "Remove a project or task by ID, along with all its sub-items.",
		Run: func(cmd *cobra.Command, args []string) {
			// Ensure sufficient arguments (either 'project' or 'task' followed by an ID)
			if len(args) < MinArgsLength {
				cmd.Println(
					"Insufficient arguments. Use 'remove project <id>' or 'remove task <id>'.",
				)
				return
			}

			// Parse the ID argument into an integer
			id, err := strconv.Atoi(args[1])
			if err != nil {
				cmd.Println("Invalid ID. Please provide a numeric ID.")
				return
			}

			// Determine whether the user wants to remove a project or a task
			switch args[0] {
			case "project":
				removeProject(cmd, projectController, id)
			case "task":
				removeTask(cmd, taskController, id)
			default:
				cmd.Println("Invalid option. Use 'remove project <id>' or 'remove task <id>'.")
			}
		},
	}

	return cmd
}

// removeProject handles the recursive removal of a project and all its subprojects.
// It uses the ProjectController to handle the deletion.
func removeProject(cmd *cobra.Command, projectController *controllers.ProjectController, id int) {
	err := projectController.RemoveProject(id)
	if err != nil {
		cmd.Printf("Error removing project: %v\n", err)
		return
	}

	cmd.Printf("Project (ID: %d) and all its subprojects removed successfully.\n", id)
}

// removeTask handles the recursive removal of a task and all its subtasks.
// It uses the TaskController to handle the deletion.
func removeTask(cmd *cobra.Command, taskController *controllers.TaskController, id int) {
	err := taskController.RemoveTask(id)
	if err != nil {
		cmd.Printf("Error removing task: %v\n", err)
		return
	}

	cmd.Printf("Task (ID: %d) and all its subtasks removed successfully.\n", id)
}
