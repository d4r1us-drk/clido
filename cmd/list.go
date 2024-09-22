package cmd

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/d4r1us-drk/clido/controllers"
	"github.com/d4r1us-drk/clido/models"
	"github.com/d4r1us-drk/clido/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
)

// NewListCmd creates and returns the 'list' command for displaying projects or tasks.
func NewListCmd(
	projectController *controllers.ProjectController,
	taskController *controllers.TaskController,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [projects|tasks]",
		Short: "List projects or tasks",
		Long:  "List all projects or tasks, optionally filtered by project for tasks.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Println("Insufficient arguments. Use 'list projects' or 'list tasks'.")
				return
			}

			// Retrieve flags for output format
			outputJSON, _ := cmd.Flags().GetBool("json")
			treeView, _ := cmd.Flags().GetBool("tree")

			switch args[0] {
			case "projects":
				listProjects(cmd, projectController, outputJSON, treeView)
			case "tasks":
				projectFilter, _ := cmd.Flags().GetString("project")
				listTasks(
					cmd,
					taskController,
					projectController,
					projectFilter,
					outputJSON,
					treeView,
				)
			default:
				cmd.Println("Invalid option. Use 'list projects' or 'list tasks'.")
			}
		},
	}

	cmd.Flags().StringP("project", "p", "", "Filter tasks by project name or ID")
	cmd.Flags().BoolP("json", "j", false, "Output list in JSON format")
	cmd.Flags().BoolP("tree", "t", false, "Display projects or tasks in a tree-like structure")

	return cmd
}

// listProjects lists all projects in either table, tree view, or JSON format.
func listProjects(
	cmd *cobra.Command,
	projectController *controllers.ProjectController,
	outputJSON bool,
	treeView bool,
) {
	projects, err := projectController.ListProjects()
	if err != nil {
		cmd.Printf("Error listing projects: %v\n", err)
		return
	}

	switch {
	case outputJSON:
		printProjectsJSON(cmd, projects)
	case treeView:
		printProjectTree(cmd, projects)
	default:
		printProjectTable(cmd, projects)
	}
}

// listTasks lists tasks, optionally filtered by a project, in table, tree view, or JSON format.
func listTasks(
	cmd *cobra.Command,
	taskController *controllers.TaskController,
	projectController *controllers.ProjectController,
	projectFilter string,
	outputJSON bool,
	treeView bool,
) {
	tasks, project, err := taskController.ListTasksByProjectFilter(projectFilter)
	if err != nil {
		cmd.Printf("Error listing tasks: %v\n", err)
		return
	}

	if !outputJSON {
		printTaskHeader(cmd, project)
	}

	switch {
	case outputJSON:
		printTasksJSON(cmd, tasks)
	case treeView:
		printTaskTree(cmd, tasks)
	default:
		printTaskTable(taskController, projectController, tasks)
	}
}

// printTaskHeader prints the header for the task list, either all tasks or tasks within a specific project.
func printTaskHeader(cmd *cobra.Command, project *models.Project) {
	if project != nil {
		cmd.Printf("Tasks in project '%s':\n", project.Name)
	} else {
		cmd.Println("All Tasks:")
	}
}

func printProjectsJSON(cmd *cobra.Command, projects []*models.Project) {
	jsonData, jsonErr := json.MarshalIndent(projects, "", "  ")
	if jsonErr != nil {
		cmd.Printf("Error marshalling projects to JSON: %v\n", jsonErr)
		return
	}
	cmd.Println(string(jsonData))
}

// printTasksJSON outputs the tasks in JSON format.
func printTasksJSON(cmd *cobra.Command, tasks []*models.Task) {
	jsonData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		cmd.Printf("Error marshalling tasks to JSON: %v\n", err)
		return
	}
	cmd.Println(string(jsonData))
}

// printProjectTable displays the list of projects in a table format.
func printProjectTable(cmd *cobra.Command, projects []*models.Project) {
	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"ID", "Name", "Description", "Type", "Child Of"})
	table.SetRowLine(true)

	for _, project := range projects {
		typeField := "Parent"
		parentChildField := "None"
		if project.ParentProjectID != nil {
			typeField = "Child"
		}
		table.Append([]string{
			strconv.Itoa(project.ID),
			utils.WrapText(project.Name, MaxProjectNameLength),
			utils.WrapText(project.Description, MaxProjectDescLength),
			typeField,
			parentChildField,
		})
	}

	cmd.Println("Projects:")
	table.Render()
}

// printTaskTable displays the list of tasks in a table format.
func printTaskTable(
	taskController *controllers.TaskController,
	projectController *controllers.ProjectController,
	tasks []*models.Task,
) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"ID", "Name", "Description", "Due Date", "Completed", "Past Due", "Priority", "Project", "Type", "Parent/Child Of",
	})
	table.SetRowLine(true)

	for _, task := range tasks {
		typeField := "Parent"
		parentChildField := "None"

		if task.ParentTaskID != nil {
			typeField = "Child"
			parentTask, _ := taskController.GetTaskByID(*task.ParentTaskID)
			if parentTask != nil {
				parentChildField = parentTask.Name
			}
		} else {
			subtasks, _ := taskController.ListSubtasks(task.ID)
			if len(subtasks) > 0 {
				typeField = "Parent"
			}
		}

		project, _ := projectController.GetProjectByID(task.ProjectID)
		projectName := ""
		if project != nil {
			projectName = project.Name
		}

		table.Append([]string{
			strconv.Itoa(task.ID),
			utils.WrapText(task.Name, MaxTaskNameLength),
			utils.WrapText(task.Description, MaxTaskDescLength),
			utils.FormatDate(task.DueDate),
			strconv.FormatBool(task.TaskCompleted),
			utils.ColoredPastDue(task.DueDate, task.TaskCompleted),
			utils.GetPriorityString(task.Priority),
			utils.WrapText(projectName, MaxProjectNameWrapLength),
			typeField,
			parentChildField,
		})
	}

	table.Render()
}

// printProjectTree displays projects in a tree view using treeprint.
func printProjectTree(cmd *cobra.Command, projects []*models.Project) {
	tree := treeprint.New()
	projectMap := make(map[int]treeprint.Tree)

	for _, project := range projects {
		projectLabel := formatProjectLabel(project)
		if project.ParentProjectID != nil {
			parentNode, exists := projectMap[*project.ParentProjectID]
			if exists {
				projectMap[project.ID] = parentNode.AddBranch(projectLabel)
			}
		} else {
			projectMap[project.ID] = tree.AddBranch(projectLabel)
		}
	}

	cmd.Println(tree.String())
}

// printTaskTree displays tasks in a tree view using treeprint.
func printTaskTree(cmd *cobra.Command, tasks []*models.Task) {
	tree := treeprint.New()
	taskMap := make(map[int]treeprint.Tree)

	for _, task := range tasks {
		taskLabel := formatTaskLabel(task)
		if task.ParentTaskID != nil {
			parentNode, exists := taskMap[*task.ParentTaskID]
			if exists {
				taskMap[task.ID] = parentNode.AddBranch(taskLabel)
			}
		} else {
			taskMap[task.ID] = tree.AddBranch(taskLabel)
		}
	}

	cmd.Println(tree.String())
}

// formatProjectLabel creates a label for each project node.
func formatProjectLabel(project *models.Project) string {
	return project.Name + " (ID: " + strconv.Itoa(project.ID) + ")"
}

// formatTaskLabel creates a label for each task node.
func formatTaskLabel(task *models.Task) string {
	return task.Name + " (ID: " + strconv.Itoa(task.ID) + ")"
}
