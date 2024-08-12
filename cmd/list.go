package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/d4r1us-drk/clido/pkg/models"
	"github.com/d4r1us-drk/clido/pkg/repository"
	"github.com/d4r1us-drk/clido/pkg/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [projects|tasks]",
	Short: "List projects or tasks",
	Long:  `List all projects or tasks, optionally filtered by project for tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Insufficient arguments. Use 'list projects' or 'list tasks'.")
			return
		}

		repo, err := repository.NewRepository()
		if err != nil {
			fmt.Printf("Error initializing repository: %v\n", err)
			return
		}
		defer repo.Close()

		switch args[0] {
		case "projects":
			listProjects(repo)
		case "tasks":
			projectFilter, _ := cmd.Flags().GetString("project")
			listTasks(repo, projectFilter)
		default:
			fmt.Println("Invalid option. Use 'list projects' or 'list tasks'.")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("project", "p", "", "Filter tasks by project name or ID")
}

func listProjects(repo *repository.Repository) {
	projects, err := repo.GetAllProjects()
	if err != nil {
		fmt.Printf("Error listing projects: %v\n", err)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Description"})
	table.SetRowLine(true)

	for _, project := range projects {
		table.Append([]string{
			strconv.Itoa(project.ID),
			utils.WrapText(project.Name, 30),
			utils.WrapText(project.Description, 50),
		})
	}

	fmt.Println("Projects:")
	table.Render()
}

func listTasks(repo *repository.Repository, projectFilter string) {
	var tasks []*models.Task
	var err error

	if projectFilter != "" {
		var project *models.Project
		if utils.IsNumeric(projectFilter) {
			projectID, _ := strconv.Atoi(projectFilter)
			project, err = repo.GetProjectByID(projectID)
		} else {
			project, err = repo.GetProjectByName(projectFilter)
		}

		if err != nil || project == nil {
			fmt.Printf("Project '%s' not found.\n", projectFilter)
			return
		}

		tasks, err = repo.GetTasksByProjectID(project.ID)
		if err != nil {
			fmt.Printf("Error listing tasks: %v\n", err)
			return
		}

		fmt.Printf("Tasks in project '%s':\n", project.Name)
	} else {
		tasks, err = repo.GetAllTasks()
		if err != nil {
			fmt.Printf("Error listing tasks: %v\n", err)
			return
		}

		fmt.Println("All Tasks:")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(
		[]string{
			"ID",
			"Name",
			"Description",
			"Due Date",
			"Completed",
			"Past Due",
			"Priority",
			"Project",
		},
	)
	table.SetRowLine(true)

	for _, task := range tasks {
		project, _ := repo.GetProjectByID(task.ProjectID)
		projectName := ""
		if project != nil {
			projectName = project.Name
		}

		table.Append([]string{
			strconv.Itoa(task.ID),
			utils.WrapText(task.Name, 20),
			utils.WrapText(task.Description, 30),
			utils.FormatDate(task.DueDate),
			fmt.Sprintf("%v", task.TaskCompleted),
			utils.ColoredPastDue(task.DueDate, task.TaskCompleted),
			utils.GetPriorityString(task.Priority),
			utils.WrapText(projectName, 20),
		})
	}

	table.Render()
}
