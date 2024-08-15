package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

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

		outputJSON, _ := cmd.Flags().GetBool("json")
		treeView, _ := cmd.Flags().GetBool("tree")

		switch args[0] {
		case "projects":
			listProjects(repo, outputJSON, treeView)
		case "tasks":
			projectFilter, _ := cmd.Flags().GetString("project")
			listTasks(repo, projectFilter, outputJSON, treeView)
		default:
			fmt.Println("Invalid option. Use 'list projects' or 'list tasks'.")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("project", "p", "", "Filter tasks by project name or ID")
	listCmd.Flags().BoolP("json", "j", false, "Output list in JSON format")
	listCmd.Flags().BoolP("tree", "t", false, "Display projects or tasks in a tree-like structure")
}

func listProjects(repo *repository.Repository, outputJSON bool, treeView bool) {
	projects, err := repo.GetAllProjects()
	if err != nil {
		fmt.Printf("Error listing projects: %v\n", err)
		return
	}

	if outputJSON {
		jsonData, err := json.MarshalIndent(projects, "", "  ")
		if err != nil {
			fmt.Printf("Error marshalling projects to JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	} else if treeView {
		printProjectTree(repo, projects, nil, 0)
	} else {
		printProjectTable(repo, projects)
	}
}

func listTasks(repo *repository.Repository, projectFilter string, outputJSON bool, treeView bool) {
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

		if !outputJSON {
			fmt.Printf("Tasks in project '%s':\n", project.Name)
		}
	} else {
		tasks, err = repo.GetAllTasks()
		if err != nil {
			fmt.Printf("Error listing tasks: %v\n", err)
			return
		}

		if !outputJSON {
			fmt.Println("All Tasks:")
		}
	}

	if outputJSON {
		jsonData, err := json.MarshalIndent(tasks, "", "  ")
		if err != nil {
			fmt.Printf("Error marshalling tasks to JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	} else if treeView {
		printTaskTree(repo, tasks, nil, 0)
	} else {
		printTaskTable(repo, tasks)
	}
}

func printProjectTable(repo *repository.Repository, projects []*models.Project) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Description", "Type", "Child Of"})
	table.SetRowLine(true)

	for _, project := range projects {
		typeField := "Parent"
		parentChildField := "None"
		if project.ParentProjectId != nil {
			typeField = "Child"
			parentProject, _ := repo.GetProjectByID(*project.ParentProjectId)
			if parentProject != nil {
				parentChildField = parentProject.Name
			}
		} else {
			subprojects, _ := repo.GetSubprojects(project.ID)
			if len(subprojects) > 0 {
				typeField = "Parent"
			}
		}

		table.Append([]string{
			strconv.Itoa(project.ID),
			utils.WrapText(project.Name, 30),
			utils.WrapText(project.Description, 50),
			typeField,
			parentChildField,
		})
	}

	fmt.Println("Projects:")
	table.Render()
}

func printTaskTable(repo *repository.Repository, tasks []*models.Task) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"ID", "Name", "Description", "Due Date", "Completed", "Past Due", "Priority", "Project", "Type", "Parent/Child Of",
	})
	table.SetRowLine(true)

	for _, task := range tasks {
		typeField := "Parent"
		parentChildField := "None"
		if task.ParentTaskId != nil {
			typeField = "Child"
			parentTask, _ := repo.GetTaskByID(*task.ParentTaskId)
			if parentTask != nil {
				parentChildField = parentTask.Name
			}
		} else {
			subtasks, _ := repo.GetSubtasks(task.ID)
			if len(subtasks) > 0 {
				typeField = "Parent"
			}
		}

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
			typeField,
			parentChildField,
		})
	}

	table.Render()
}

func printProjectTree(
	repo *repository.Repository,
	projects []*models.Project,
	parentID *int,
	level int,
) {
	indent := strings.Repeat("│  ", level)
	for i, project := range projects {
		if (parentID == nil && project.ParentProjectId == nil) ||
			(parentID != nil && project.ParentProjectId != nil && *project.ParentProjectId == *parentID) {
			prefix := "├──"
			if i == len(projects)-1 {
				prefix = "└──"
			}
			fmt.Printf("%s%s %s (ID: %d)\n", indent, prefix, project.Name, project.ID)
			printProjectTree(repo, projects, &project.ID, level+1)
		}
	}
}

func printTaskTree(repo *repository.Repository, tasks []*models.Task, parentID *int, level int) {
	indent := strings.Repeat("│  ", level)
	for i, task := range tasks {
		if (parentID == nil && task.ParentTaskId == nil) ||
			(parentID != nil && task.ParentTaskId != nil && *task.ParentTaskId == *parentID) {
			prefix := "├──"
			if i == len(tasks)-1 {
				prefix = "└──"
			}
			fmt.Printf("%s%s %s (ID: %d)\n", indent, prefix, task.Name, task.ID)
			fmt.Printf("%s    Description: %s\n", indent, task.Description)
			fmt.Printf(
				"%s    Due Date: %s, Completed: %v, Priority: %s\n",
				indent,
				utils.FormatDate(task.DueDate),
				task.TaskCompleted,
				utils.GetPriorityString(task.Priority),
			)
			printTaskTree(repo, tasks, &task.ID, level+1)
		}
	}
}
