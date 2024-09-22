package cmd

import (
	"errors"

	"github.com/d4r1us-drk/clido/controllers"
	"github.com/spf13/cobra"
)

// NewNewCmd creates and returns the 'new' command for creating projects or tasks.
func NewNewCmd(
	projectController *controllers.ProjectController,
	taskController *controllers.TaskController,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [project|task]",
		Short: "Create a new project or task",
		Long:  "Create a new project or task with the specified details.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("insufficient arguments. Use 'new project' or 'new task'")
			}

			// Create project or task based on the argument
			switch args[0] {
			case "project":
				return createProject(cmd, projectController)
			case "task":
				return createTask(cmd, taskController)
			default:
				return errors.New("invalid option. Use 'new project' or 'new task'")
			}
		},
	}

	// Define flags for project and task creation
	cmd.Flags().StringP("name", "n", "", "Name of the project or task")
	cmd.Flags().StringP("description", "d", "", "Description of the project or task")
	cmd.Flags().StringP("project", "p", "", "Parent project name or ID for subprojects or tasks")
	cmd.Flags().StringP("task", "t", "", "Parent task ID for subtasks")
	cmd.Flags().StringP("due", "D", "", "Due date for the task (format: YYYY-MM-DD HH:MM)")
	cmd.Flags().
		IntP("priority", "P", PriorityEmpty, "Priority of the task (1: High, 2: Medium, 3: Low, 4: None)")

	return cmd
}

func createProject(cmd *cobra.Command, projectController *controllers.ProjectController) error {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	parentProjectIdentifier, _ := cmd.Flags().GetString("project")

	// Ensure project name is provided
	if name == "" {
		return errors.New("project name is required")
	}

	// Call the controller to create the project
	err := projectController.CreateProject(name, description, parentProjectIdentifier)
	if err != nil {
		return errors.New("error creating project: " + err.Error())
	}

	cmd.Println("Project '" + name + "' created successfully.")
	return nil
}

func createTask(cmd *cobra.Command, taskController *controllers.TaskController) error {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	projectIdentifier, _ := cmd.Flags().GetString("project")
	parentTaskIdentifier, _ := cmd.Flags().GetString("task")
	dueDateStr, _ := cmd.Flags().GetString("due")
	priority, _ := cmd.Flags().GetInt("priority")

	// Ensure task name is provided
	if name == "" {
		return errors.New("task name is required")
	}

	// Validate priority
	if priority != 0 && (priority < 1 || priority > 4) {
		return errors.New(
			"invalid priority. Use 1 for High, 2 for Medium, 3 for Low, or 4 for None",
		)
	}

	// Call the controller to create the task
	err := taskController.CreateTask(
		name,
		description,
		projectIdentifier,
		parentTaskIdentifier,
		dueDateStr,
		priority,
	)
	if err != nil {
		return errors.New("error creating task: " + err.Error())
	}

	cmd.Println("Task '" + name + "' created successfully.")
	return nil
}
