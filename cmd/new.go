package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/d4r1us-drk/clido/pkg/models"
	"github.com/d4r1us-drk/clido/pkg/repository"
	"github.com/d4r1us-drk/clido/pkg/utils"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [project|task]",
	Short: "Create a new project or task",
	Long:  `Create a new project or task with the specified details.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Insufficient arguments. Use 'new project' or 'new task'.")
			return
		}

		repo, err := repository.NewRepository()
		if err != nil {
			fmt.Printf("Error initializing repository: %v\n", err)
			return
		}
		defer repo.Close()

		switch args[0] {
		case "project":
			createProject(cmd, repo)
		case "task":
			createTask(cmd, repo)
		default:
			fmt.Println("Invalid option. Use 'new project' or 'new task'.")
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringP("name", "n", "", "Name of the project or task")
	newCmd.Flags().StringP("description", "d", "", "Description of the project or task")
	newCmd.Flags().StringP("project", "p", "", "Project name or ID for the task")
	newCmd.Flags().StringP("due", "D", "", "Due date for the task (format: YYYY-MM-DD HH:MM)")
	newCmd.Flags().
		IntP("priority", "r", 4, "Priority of the task (1: High, 2: Medium, 3: Low, 4: None)")
}

func createProject(cmd *cobra.Command, repo *repository.Repository) {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")

	if name == "" {
		fmt.Println("Project name is required.")
		return
	}

	project := &models.Project{
		Name:        name,
		Description: description,
	}

	err := repo.CreateProject(project)
	if err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		return
	}

	fmt.Printf("Project '%s' created successfully.\n", name)
}

func createTask(cmd *cobra.Command, repo *repository.Repository) {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	projectIdentifier, _ := cmd.Flags().GetString("project")
	dueDateStr, _ := cmd.Flags().GetString("due")
	priority, _ := cmd.Flags().GetInt("priority")

	if name == "" || projectIdentifier == "" {
		fmt.Println("Task name and project identifier are required.")
		return
	}

	var project *models.Project
	var err error

	if utils.IsNumeric(projectIdentifier) {
		projectID, _ := strconv.Atoi(projectIdentifier)
		project, err = repo.GetProjectByID(projectID)
	} else {
		project, err = repo.GetProjectByName(projectIdentifier)
	}

	if err != nil || project == nil {
		fmt.Printf("Project '%s' not found.\n", projectIdentifier)
		return
	}

	var dueDate *time.Time
	if dueDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02 15:04", dueDateStr)
		if err == nil {
			dueDate = &parsedDate
		} else {
			fmt.Println("Invalid date format. Using no due date.")
		}
	}

	task := &models.Task{
		Name:        name,
		Description: description,
		ProjectID:   project.ID,
		DueDate:     dueDate,
		Priority:    priority,
	}

	err = repo.CreateTask(task)
	if err != nil {
		fmt.Printf("Error creating task: %v\n", err)
		return
	}

	fmt.Printf(
		"Task '%s' created successfully with priority %s.\n",
		name,
		utils.GetPriorityString(priority),
	)
}
