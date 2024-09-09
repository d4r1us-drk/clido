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
	ErrNoParentTask = errors.New("no parent task specified")
	ErrNoDueDate    = errors.New("no due date specified")
)

// NewNewCmd creates and returns the new command.
func NewNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [project|task]",
		Short: "Create a new project or task",
		Long:  `Create a new project or task with the specified details.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Println("Insufficient arguments. Use 'new project' or 'new task'.")
				return
			}

			repo, err := repository.NewRepository()
			if err != nil {
				cmd.Printf("Error initializing repository: %v\n", err)
				return
			}
			defer repo.Close()

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

	// Define flags for the new command
	cmd.Flags().StringP("name", "n", "", "Name of the project or task")
	cmd.Flags().StringP("description", "d", "", "Description of the project or task")
	cmd.Flags().StringP("project", "p", "", "Parent project name or ID for subprojects or tasks")
	cmd.Flags().StringP("task", "t", "", "Parent task ID for subtasks")
	cmd.Flags().StringP("due", "D", "", "Due date for the task (format: YYYY-MM-DD HH:MM)")
	cmd.Flags().
		IntP("priority", "P", 0, "Priority of the task (1: High, 2: Medium, 3: Low, 4: None)")

	return cmd
}

func createProject(cmd *cobra.Command, repo *repository.Repository) {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	parentProjectIdentifier, _ := cmd.Flags().GetString("project")

	if name == "" {
		cmd.Println("Project name is required.")
		return
	}

	var parentProjectID *int
	if parentProjectIdentifier != "" {
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

	project := &models.Project{
		Name:            name,
		Description:     description,
		ParentProjectID: parentProjectID,
	}

	err := repo.CreateProject(project)
	if err != nil {
		cmd.Printf("Error creating project: %v\n", err)
		return
	}

	cmd.Printf("Project '%s' created successfully.\n", name)
}

func createTask(cmd *cobra.Command, repo *repository.Repository) {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	projectIdentifier, _ := cmd.Flags().GetString("project")
	parentTaskIdentifier, _ := cmd.Flags().GetString("task")
	dueDateStr, _ := cmd.Flags().GetString("due")
	priority, _ := cmd.Flags().GetInt("priority")

	if name == "" {
		cmd.Println("Task name is required.")
		return
	}

	projectID, err := getProjectID(projectIdentifier, repo)
	if err != nil {
		cmd.Println(err)
		return
	}

	parentTaskID, err := getParentTaskID(parentTaskIdentifier)
	if err != nil && !errors.Is(err, ErrNoParentTask) {
		cmd.Println(err)
		return
	}

	dueDate, err := parseDueDate(dueDateStr)
	if err != nil && !errors.Is(err, ErrNoDueDate) {
		cmd.Println("Invalid date format. Using no due date.")
	}

	task := &models.Task{
		Name:         name,
		Description:  description,
		ProjectID:    projectID,
		DueDate:      dueDate,
		Priority:     utils.Priority(priority),
		ParentTaskID: parentTaskID,
	}

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
