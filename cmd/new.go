package cmd

import (
	"errors"
	"strconv"
	"time"

	"github.com/d4r1us-drk/clido/pkg/models"
	"github.com/d4r1us-drk/clido/pkg/repository"
	"github.com/d4r1us-drk/clido/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	// ErrNoParentTask is returned when no parent task is provided during task creation.
	ErrNoParentTask = errors.New("no parent task specified")

	// ErrNoDueDate is returned when no due date is provided during task creation.
	ErrNoDueDate = errors.New("no due date specified")
)

// NewNewCmd creates and returns the 'new' command for creating projects or tasks.
//
// The command allows users to create new projects or tasks with the specified details,
// such as name, description, parent project, parent task, due date, and priority.
//
// Usage:
//   clido new project    # Create a new project
//   clido new task       # Create a new task
func NewNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [project|task]",  // Specifies valid options: 'project' or 'task'
		Short: "Create a new project or task",  // Short description of the command
		Long:  `Create a new project or task with the specified details.`,  // Extended description
		Run: func(cmd *cobra.Command, args []string) {
			// Ensure sufficient arguments (either 'project' or 'task')
			if len(args) < 1 {
				cmd.Println("Insufficient arguments. Use 'new project' or 'new task'.")
				return
			}

			// Initialize the repository for database operations
			repo, err := repository.NewRepository()
			if err != nil {
				cmd.Printf("Error initializing repository: %v\n", err)
				return
			}
			defer repo.Close()

			// Determine whether the user wants to create a project or a task
			switch args[0] {
			case "project":
				createProject(cmd, repo)
			case "task":
				createTask(cmd, repo)
			default:
				cmd.Println("Invalid option. Use 'new project' or 'new task'.")
			}
		},
	}

	// Define flags for the new command, allowing users to specify details
	cmd.Flags().StringP("name", "n", "", "Name of the project or task")  // Name of the project/task (required)
	cmd.Flags().StringP("description", "d", "", "Description of the project or task")  // Description
	cmd.Flags().StringP("project", "p", "", "Parent project name or ID for subprojects or tasks")  // Parent project
	cmd.Flags().StringP("task", "t", "", "Parent task ID for subtasks")  // Parent task for subtasks
	cmd.Flags().StringP("due", "D", "", "Due date for the task (format: YYYY-MM-DD HH:MM)")  // Due date
	cmd.Flags().IntP("priority", "P", 0, "Priority of the task (1: High, 2: Medium, 3: Low, 4: None)")  // Task priority

	return cmd
}

// createProject handles the creation of a new project.
// It retrieves input from flags (name, description, and parent project), validates them,
// and saves the new project to the database.
func createProject(cmd *cobra.Command, repo *repository.Repository) {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	parentProjectIdentifier, _ := cmd.Flags().GetString("project")

	// Ensure the project name is provided
	if name == "" {
		cmd.Println("Project name is required.")
		return
	}

	var parentProjectID *int
	if parentProjectIdentifier != "" {
		// Determine whether the parent project is specified by ID or name
		if utils.IsNumeric(parentProjectIdentifier) {
			id, _ := strconv.Atoi(parentProjectIdentifier)
			parentProjectID = &id
		} else {
			parentProject, err := repo.GetProjectByName(parentProjectIdentifier)
			if err != nil || parentProject == nil {
				cmd.Printf("Parent project '%s' not found.\n", parentProjectIdentifier)
				return
			}
			parentProjectID = &parentProject.ID
		}
	}

	// Create the new project
	project := &models.Project{
		Name:            name,
		Description:     description,
		ParentProjectID: parentProjectID,
	}

	// Save the new project to the database
	err := repo.CreateProject(project)
	if err != nil {
		cmd.Printf("Error creating project: %v\n", err)
		return
	}

	cmd.Printf("Project '%s' created successfully.\n", name)
}

// createTask handles the creation of a new task.
// It retrieves input from flags (name, description, project, parent task, due date, and priority),
// validates them, and saves the new task to the database.
func createTask(cmd *cobra.Command, repo *repository.Repository) {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	projectIdentifier, _ := cmd.Flags().GetString("project")
	parentTaskIdentifier, _ := cmd.Flags().GetString("task")
	dueDateStr, _ := cmd.Flags().GetString("due")
	priority, _ := cmd.Flags().GetInt("priority")

	// Ensure the task name is provided
	if name == "" {
		cmd.Println("Task name is required.")
		return
	}

	// Validate and retrieve the project ID
	projectID, err := getProjectID(projectIdentifier, repo)
	if err != nil {
		cmd.Println(err)
		return
	}

	// Validate and retrieve the parent task ID, if provided
	parentTaskID, err := getParentTaskID(parentTaskIdentifier)
	if err != nil && !errors.Is(err, ErrNoParentTask) {
		cmd.Println(err)
		return
	}

	// Validate and parse the due date, if provided
	dueDate, err := parseDueDate(dueDateStr)
	if err != nil && !errors.Is(err, ErrNoDueDate) {
		cmd.Println("Invalid date format. Using no due date.")
	}

	// Create the new task
	task := &models.Task{
		Name:         name,
		Description:  description,
		ProjectID:    projectID,
		DueDate:      dueDate,
		Priority:     utils.Priority(priority),
		ParentTaskID: parentTaskID,
	}

	// Save the new task to the database
	if err = repo.CreateTask(task); err != nil {
		cmd.Printf("Error creating task: %v\n", err)
		return
	}

	cmd.Printf(
		"Task '%s' created successfully with priority %s.\n",
		name,
		utils.GetPriorityString(utils.Priority(priority)),
	)
}

// getProjectID retrieves the project ID based on the identifier (either name or ID).
// It returns an error if the project does not exist.
func getProjectID(identifier string, repo *repository.Repository) (int, error) {
	if identifier == "" {
		return 0, errors.New("task must be associated with a project")
	}

	if utils.IsNumeric(identifier) {
		return strconv.Atoi(identifier)
	}

	project, err := repo.GetProjectByName(identifier)
	if err != nil || project == nil {
		return 0, errors.New("project '" + identifier + "' not found")
	}

	return project.ID, nil
}

// getParentTaskID retrieves the parent task ID based on the identifier.
// It returns an error if the identifier is not a numeric ID.
func getParentTaskID(identifier string) (*int, error) {
	if identifier == "" {
		return nil, ErrNoParentTask
	}

	if !utils.IsNumeric(identifier) {
		return nil, errors.New("parent task must be identified by a numeric ID")
	}

	id, _ := strconv.Atoi(identifier)
	return &id, nil
}

// parseDueDate parses a string into a time.Time object.
// It returns an error if the string is not in the expected format.
func parseDueDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, ErrNoDueDate
	}

	parsedDate, err := time.Parse("2006-01-02 15:04", dateStr)
	if err != nil {
		return nil, err
	}

	return &parsedDate, nil
}
