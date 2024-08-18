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

	switch {
	case outputJSON:
		var jsonData []byte
		jsonData, err = json.MarshalIndent(projects, "", "  ")
		if err != nil {
			fmt.Printf("Error marshalling projects to JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))

	case treeView:
		printProjectTree(projects, nil, 0)

	default:
		printProjectTable(repo, projects)
	}
}

func listTasks(repo *repository.Repository, projectFilter string, outputJSON bool, treeView bool) {
	tasks, project, err := getTasks(repo, projectFilter)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !outputJSON {
		printTaskHeader(project)
	}

	switch {
	case outputJSON:
		printTasksJSON(tasks)
	case treeView:
		printTaskTree(tasks, nil, 0)
	default:
		printTaskTable(repo, tasks)
	}
}

func getTasks(repo *repository.Repository, projectFilter string) ([]*models.Task, *models.Project, error) {
	if projectFilter == "" {
		tasks, err := repo.GetAllTasks()
		return tasks, nil, err
	}

	project, err := getProject(repo, projectFilter)
	if err != nil {
		return nil, nil, err
	}

	tasks, err := repo.GetTasksByProjectID(project.ID)
	return tasks, project, err
}

func getProject(repo *repository.Repository, projectFilter string) (*models.Project, error) {
	var project *models.Project
	var err error

	if utils.IsNumeric(projectFilter) {
		projectID, _ := strconv.Atoi(projectFilter)
		project, err = repo.GetProjectByID(projectID)
	} else {
		project, err = repo.GetProjectByName(projectFilter)
	}

	if err != nil || project == nil {
		return nil, fmt.Errorf("project '%s' not found", projectFilter)
	}

	return project, nil
}

func printTaskHeader(project *models.Project) {
	if project != nil {
		fmt.Printf("Tasks in project '%s':\n", project.Name)
	} else {
		fmt.Println("All Tasks:")
	}
}

func printTasksJSON(tasks []*models.Task) {
	jsonData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling tasks to JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

func printProjectTable(repo *repository.Repository, projects []*models.Project) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Description", "Type", "Child Of"})
	table.SetRowLine(true)

	for _, project := range projects {
		typeField := "Parent"
		parentChildField := "None"
		if project.ParentProjectID != nil {
			typeField = "Child"
			parentProject, _ := repo.GetProjectByID(*project.ParentProjectID)
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
			utils.WrapText(project.Name, MaxProjectNameLength),
			utils.WrapText(project.Description, MaxProjectDescLength),
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
		if task.ParentTaskID != nil {
			typeField = "Child"
			parentTask, _ := repo.GetTaskByID(*task.ParentTaskID)
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

type TreeNode interface {
	GetID() int
	GetParentID() *int
	GetName() string
}

type ProjectNode struct {
	*models.Project
}

func (p ProjectNode) GetID() int        { return p.ID }
func (p ProjectNode) GetParentID() *int { return p.ParentProjectID }
func (p ProjectNode) GetName() string   { return p.Name }

type TaskNode struct {
	*models.Task
}

func (t TaskNode) GetID() int        { return t.ID }
func (t TaskNode) GetParentID() *int { return t.ParentTaskID }
func (t TaskNode) GetName() string   { return t.Name }

func printTree(nodes []TreeNode, parentID *int, level int, printDetails func(TreeNode, string)) {
	indent := strings.Repeat("│  ", level)
	for i, node := range nodes {
		if (parentID == nil && node.GetParentID() == nil) ||
			(parentID != nil && node.GetParentID() != nil && *node.GetParentID() == *parentID) {
			prefix := "├──"
			if i == len(nodes)-1 {
				prefix = "└──"
			}
			fmt.Printf("%s%s %s (ID: %d)\n", indent, prefix, node.GetName(), node.GetID())
			if printDetails != nil {
				printDetails(node, indent+"    ")
			}
			nodeID := node.GetID()
			printTree(nodes, &nodeID, level+1, printDetails)
		}
	}
}

func printProjectTree(projects []*models.Project, parentID *int, level int) {
	nodes := make([]TreeNode, len(projects))
	for i, p := range projects {
		nodes[i] = ProjectNode{p}
	}
	printTree(nodes, parentID, level, nil)
}

func printTaskTree(tasks []*models.Task, parentID *int, level int) {
	nodes := make([]TreeNode, len(tasks))
	for i, t := range tasks {
		nodes[i] = TaskNode{t}
	}
	printTree(nodes, parentID, level, func(node TreeNode, indent string) {
		task := node.(TaskNode).Task
		fmt.Printf("%sDescription: %s\n", indent, task.Description)
		fmt.Printf(
			"%sDue Date: %s, Completed: %v, Priority: %s\n",
			indent,
			utils.FormatDate(task.DueDate),
			task.TaskCompleted,
			utils.GetPriorityString(task.Priority),
		)
	})
}
