package cmd

import (
	"strconv"
	"time"

	"github.com/d4r1us-drk/clido/pkg/models"
	"github.com/d4r1us-drk/clido/pkg/repository"
	"github.com/d4r1us-drk/clido/pkg/utils"
	"github.com/spf13/cobra"
)

// NewEditCmd creates and returns the edit command.
func NewEditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit [project|task] <id>",
		Short: "Edit an existing project or task",
		Long:  `Edit the details of an existing project or task identified by its ID.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < MinArgsLength {
				cmd.Println("Insufficient arguments. Use 'edit project <id>' or 'edit task <id>'.")
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
				editProject(cmd, repo, id)
			case "task":
				editTask(cmd, repo, id)
			default:
				cmd.Println("Invalid option. Use 'edit project <id>' or 'edit task <id>'.")
			}
		},
	}

	// Define flags for the edit command
	cmd.Flags().StringP("name", "n", "", "New name")
	cmd.Flags().StringP("description", "d", "", "New description")
	cmd.Flags().StringP("project", "p", "", "New parent project name or ID")
	cmd.Flags().StringP("task", "t", "", "New parent task ID for subtasks")
	cmd.Flags().StringP("due", "D", "", "New due date for task (format: YYYY-MM-DD HH:MM)")
	cmd.Flags().
		IntP("priority", "P", 0, "New priority for task (1: High, 2: Medium, 3: Low, 4: None)")

	return cmd
}

func editProject(cmd *cobra.Command, repo *repository.Repository, id int) {
	project, err := repo.GetProjectByID(id)
	if err != nil {
		cmd.Printf("Error retrieving project: %v\n", err)
		return
	}

	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	parentProjectIdentifier, _ := cmd.Flags().GetString("project")

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
			var parentProject *models.Project
			parentProject, err = repo.GetProjectByName(parentProjectIdentifier)
			if err != nil || parentProject == nil {
				cmd.Printf("Parent project '%s' not found.\n", parentProjectIdentifier)
				return
			}
			project.ParentProjectID = &parentProject.ID
		}
	}

	err = repo.UpdateProject(project)
	if err != nil {
		cmd.Printf("Error updating project: %v\n", err)
		return
	}

	cmd.Printf("Project '%s' updated successfully.\n", project.Name)
}

func editTask(cmd *cobra.Command, repo *repository.Repository, id int) {
	task, err := repo.GetTaskByID(id)
	if err != nil {
		cmd.Printf("Error retrieving task: %v\n", err)
		return
	}

	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	dueDateStr, _ := cmd.Flags().GetString("due")
	priority, _ := cmd.Flags().GetInt("priority")
	parentTaskIdentifier, _ := cmd.Flags().GetString("task")

	if name != "" {
		task.Name = name
	}
	if description != "" {
		task.Description = description
	}
	if dueDateStr != "" {
		var parsedDate time.Time
		parsedDate, err = time.Parse("2006-01-02 15:04", dueDateStr)
		if err == nil {
			task.DueDate = &parsedDate
		} else {
			cmd.Println("Invalid date format. Keeping the existing due date.")
		}
	}
	if priority != 0 {
		if priority >= 1 && priority <= 4 {
			task.Priority = utils.Priority(priority)
		} else {
			cmd.Println("Invalid priority. Keeping the existing priority.")
		}
	}
	if parentTaskIdentifier != "" {
		if utils.IsNumeric(parentTaskIdentifier) {
			parentID, _ := strconv.Atoi(parentTaskIdentifier)
			task.ParentTaskID = &parentID
		} else {
			cmd.Println("Parent task must be identified by a numeric ID.")
			return
		}
	}

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
