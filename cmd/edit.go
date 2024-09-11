package cmd

import (
	"strconv"
	"time"

	"github.com/d4r1us-drk/clido/pkg/models"
	"github.com/d4r1us-drk/clido/pkg/repository"
	"github.com/d4r1us-drk/clido/pkg/utils"
	"github.com/spf13/cobra"
)

// NewEditCmd creates and returns the 'edit' command for editing projects or tasks.
//
// The command allows users to modify existing projects or tasks by their unique ID.
// It supports editing the name, description, parent project, parent task, due date, and priority of a task.
//
// Usage:
//   clido edit project <id>    # Edit a project by ID
//   clido edit task <id>       # Edit a task by ID
func NewEditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit [project|task] <id>",  // Specifies the valid options: project or task followed by an ID
		Short: "Edit an existing project or task",  // Short description of the command
		Long:  `Edit the details of an existing project or task identified by its ID.`,  // Extended description
		Run: func(cmd *cobra.Command, args []string) {
			// Ensure the command receives sufficient arguments (either "project" or "task" followed by an ID)
			if len(args) < MinArgsLength {
				cmd.Println("Insufficient arguments. Use 'edit project <id>' or 'edit task <id>'.")
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

			// Determine whether the user wants to edit a project or a task
			switch args[0] {
			case "project":
				editProject(cmd, repo, id)
			case "task":
				editTask(cmd, repo, id)
			default:
				cmd.Println("Invalid option. Use 'edit project <id>' or 'edit task <id>'.")
			}
		},
	}

	// Define flags for the edit command, allowing users to specify what fields they want to update
	cmd.Flags().StringP("name", "n", "", "New name")  // Option to change the name of the project/task
	cmd.Flags().StringP("description", "d", "", "New description")  // Option to change the description
	cmd.Flags().StringP("project", "p", "", "New parent project name or ID")  // Option to change the parent project (for projects)
	cmd.Flags().StringP("task", "t", "", "New parent task ID for subtasks")  // Option to change the parent task (for tasks)
	cmd.Flags().StringP("due", "D", "", "New due date for task (format: YYYY-MM-DD HH:MM)")  // Option to set a new due date
	cmd.Flags().IntP("priority", "P", 0, "New priority for task (1: High, 2: Medium, 3: Low, 4: None)")  // Option to set a new priority level

	return cmd
}

// editProject handles updating an existing project by its ID.
// The function retrieves the project from the repository, applies updates (name, description, parent project),
// and saves the changes back to the database.
//
// If the user provides a new parent project, it validates whether the project exists by name or ID.
func editProject(cmd *cobra.Command, repo *repository.Repository, id int) {
	// Retrieve the project by its ID
	project, err := repo.GetProjectByID(id)
	if err != nil {
		cmd.Printf("Error retrieving project: %v\n", err)
		return
	}

	// Get the new values from the command flags
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	parentProjectIdentifier, _ := cmd.Flags().GetString("project")

	// Apply the updates to the project, if specified by the user
	if name != "" {
		project.Name = name
	}
	if description != "" {
		project.Description = description
	}
	if parentProjectIdentifier != "" {
		if utils.IsNumeric(parentProjectIdentifier) {
			parentID, _ := strconv.Atoi(parentProjectIdentifier)
			project.ParentProjectID = &parentID
		} else {
			// If the parent project is provided by name, fetch it by name
			var parentProject *models.Project
			parentProject, err = repo.GetProjectByName(parentProjectIdentifier)
			if err != nil || parentProject == nil {
				cmd.Printf("Parent project '%s' not found.\n", parentProjectIdentifier)
				return
			}
			project.ParentProjectID = &parentProject.ID
		}
	}

	// Save the updated project to the database
	err = repo.UpdateProject(project)
	if err != nil {
		cmd.Printf("Error updating project: %v\n", err)
		return
	}

	cmd.Printf("Project '%s' updated successfully.\n", project.Name)
}

// editTask handles updating an existing task by its ID.
// The function retrieves the task from the repository, applies updates (name, description, due date, priority, parent task),
// and saves the changes back to the database.
//
// If the user provides a new parent task, it validates whether the task exists by ID.
func editTask(cmd *cobra.Command, repo *repository.Repository, id int) {
	// Retrieve the task by its ID
	task, err := repo.GetTaskByID(id)
	if err != nil {
		cmd.Printf("Error retrieving task: %v\n", err)
		return
	}

	// Get the new values from the command flags
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	dueDateStr, _ := cmd.Flags().GetString("due")
	priority, _ := cmd.Flags().GetInt("priority")
	parentTaskIdentifier, _ := cmd.Flags().GetString("task")

	// Apply the updates to the task, if specified by the user
	if name != "" {
		task.Name = name
	}
	if description != "" {
		task.Description = description
	}
	if dueDateStr != "" {
		// Parse the new due date
		var parsedDate time.Time
		parsedDate, err = time.Parse("2006-01-02 15:04", dueDateStr)
		if err == nil {
			task.DueDate = &parsedDate
		} else {
			cmd.Println("Invalid date format. Keeping the existing due date.")
		}
	}
	if priority != 0 {
		// Validate the priority (must be between 1 and 4)
		if priority >= 1 && priority <= 4 {
			task.Priority = utils.Priority(priority)
		} else {
			cmd.Println("Invalid priority. Keeping the existing priority.")
		}
	}
	if parentTaskIdentifier != "" {
		// Validate the parent task ID (must be numeric)
		if utils.IsNumeric(parentTaskIdentifier) {
			parentID, _ := strconv.Atoi(parentTaskIdentifier)
			task.ParentTaskID = &parentID
		} else {
			cmd.Println("Parent task must be identified by a numeric ID.")
			return
		}
	}

	// Save the updated task to the database
	err = repo.UpdateTask(task)
	if err != nil {
		cmd.Printf("Error updating task: %v\n", err)
		return
	}

	cmd.Printf("Task '%s' updated successfully.\n", task.Name)
	cmd.Printf("New details: Priority: %s, Due Date: %s\n",
		utils.GetPriorityString(task.Priority),
		utils.FormatDate(task.DueDate))
}
