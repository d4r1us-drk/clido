package cmd

import (
	"errors"
	"strconv"

	"github.com/d4r1us-drk/clido/controllers"
	"github.com/d4r1us-drk/clido/utils"
	"github.com/spf13/cobra"
)

// NewEditCmd creates and returns the 'edit' command for editing projects or tasks.
func NewEditCmd(
	projectController *controllers.ProjectController,
	taskController *controllers.TaskController,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit [project|task] <id>",
		Short: "Edit an existing project or task",
		Long:  "Edit the details of an existing project or task identified by its ID.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Ensure the command receives sufficient arguments (either "project" or "task" followed by an ID)
			if len(args) < MinArgsLength {
				return errors.New(
					"insufficient arguments. Use 'edit project <id>' or 'edit task <id>'",
				)
			}

			// Parse the ID argument into an integer
			id, err := strconv.Atoi(args[1])
			if err != nil {
				return errors.New("invalid ID. Please provide a numeric ID")
			}

			// Determine whether the user wants to edit a project or a task
			switch args[0] {
			case "project":
				return editProject(cmd, projectController, id)
			case "task":
				return editTask(cmd, taskController, id)
			default:
				return errors.New("invalid option. Use 'edit project <id>' or 'edit task <id>'")
			}
		},
	}

	// Define flags for the edit command, allowing users to specify what fields they want to update
	cmd.Flags().StringP("name", "n", "", "New name")
	cmd.Flags().StringP("description", "d", "", "New description")
	cmd.Flags().StringP("project", "p", "", "New parent project name or ID")
	cmd.Flags().StringP("task", "t", "", "New parent task ID for subtasks")
	cmd.Flags().StringP("due", "D", "", "New due date for task (format: YYYY-MM-DD HH:MM)")
	cmd.Flags().
		IntP("priority", "P", 0, "New priority for task (1: High, 2: Medium, 3: Low, 4: None)")

	return cmd
}

// editProject handles updating an existing project by its ID.
func editProject(
	cmd *cobra.Command,
	projectController *controllers.ProjectController,
	id int,
) error {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	parentProjectIdentifier, _ := cmd.Flags().GetString("project")

	// Check if any fields are provided for update
	if name == "" && description == "" && parentProjectIdentifier == "" {
		return errors.New("no fields provided for update. " +
			"Use flags to update the name, description, or parent project")
	}

	// Call the controller to edit the project
	err := projectController.EditProject(id, name, description, parentProjectIdentifier)
	if err != nil {
		return errors.New("error updating project: " + err.Error())
	}

	cmd.Println("Project with ID '" + strconv.Itoa(id) + "' updated successfully.")
	return nil
}

// editTask handles updating an existing task by its ID.
func editTask(cmd *cobra.Command, taskController *controllers.TaskController, id int) error {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	dueDateStr, _ := cmd.Flags().GetString("due")
	priority, _ := cmd.Flags().GetInt("priority")
	parentTaskIdentifier, _ := cmd.Flags().GetString("task")

	// Validate priority if provided
	if priority != 0 && (priority < PriorityHigh || priority > PriorityNone) {
		return errors.New(
			"invalid priority. Use 1 for High, 2 for Medium, 3 for Low, or 4 for None",
		)
	}

	// Check if any fields are provided for update
	if name == "" && description == "" && dueDateStr == "" && priority == 0 &&
		parentTaskIdentifier == "" {
		return errors.New("no fields provided for update. " +
			"Use flags to update the name, description, due date, priority, or parent task")
	}

	// Call the controller to edit the task
	err := taskController.EditTask(
		id,
		name,
		description,
		dueDateStr,
		priority,
		parentTaskIdentifier,
	)
	if err != nil {
		return errors.New("error updating task: " + err.Error())
	}

	// Format and display the new details
	priorityStr := utils.GetPriorityString(priority)
	formattedDueDate := "None"
	if dueDateStr != "" {
		parsedDueDate, _ := utils.ParseDueDate(dueDateStr)
		formattedDueDate = utils.FormatDate(parsedDueDate)
	}

	cmd.Println("Task with ID '" + strconv.Itoa(id) + "' updated successfully.")
	cmd.Println("New details: Priority: " + priorityStr + ", Due Date: " + formattedDueDate)
	return nil
}
