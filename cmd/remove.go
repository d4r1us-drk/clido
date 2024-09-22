package cmd

import (
	"errors"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			// Ensure sufficient arguments (either 'project' or 'task' followed by an ID)
			if len(args) < MinArgsLength {
				return errors.New(
					"insufficient arguments. Use 'remove project <id>' or 'remove task <id>'",
				)
			}

			// Parse the ID argument into an integer
			id, err := strconv.Atoi(args[1])
			if err != nil {
				return errors.New("invalid ID. Please provide a numeric ID")
			}

			// Determine whether the user wants to remove a project or a task
			switch args[0] {
			case "project":
				return removeProject(cmd, projectController, id)
			case "task":
				return removeTask(cmd, taskController, id)
			default:
				return errors.New("invalid option. Use 'remove project <id>' or 'remove task <id>'")
			}
		},
	}

	return cmd
}

// removeProject handles the recursive removal of a project and all its subprojects.
// It uses the ProjectController to handle the deletion.
func removeProject(
	cmd *cobra.Command,
	projectController *controllers.ProjectController,
	id int,
) error {
	err := projectController.RemoveProject(id)
	if err != nil {
		return errors.New("error removing project: " + err.Error())
	}

	cmd.Println(
		"Project (ID: " + strconv.Itoa(id) + ") and all its subprojects removed successfully.",
	)
	return nil
}

// removeTask handles the recursive removal of a task and all its subtasks.
// It uses the TaskController to handle the deletion.
func removeTask(cmd *cobra.Command, taskController *controllers.TaskController, id int) error {
	err := taskController.RemoveTask(id)
	if err != nil {
		return errors.New("error removing task: " + err.Error())
	}

	cmd.Println("Task (ID: " + strconv.Itoa(id) + ") and all its subtasks removed successfully.")
	return nil
}
