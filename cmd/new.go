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
		case "task":
			createTask(cmd, repo)
		default:
			fmt.Println("Invalid option. Use 'new project' or 'new task'.")
		}
	},
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
