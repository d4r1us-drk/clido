package cmd

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/d4r1us-drk/clido/pkg/models"
	"github.com/d4r1us-drk/clido/pkg/repository"
	"github.com/d4r1us-drk/clido/pkg/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// NewListCmd creates and returns the 'list' command for displaying projects or tasks.
//
// The command allows users to list all projects or tasks. Tasks can be optionally filtered by project.
// Users can display the output in a table, tree view, or JSON format.
//
// Usage:
//   clido list projects      # List all projects
//   clido list tasks         # List all tasks
//   clido list tasks -p 1    # List tasks filtered by project ID
//   clido list tasks -p MyProject # List tasks filtered by project name
//   clido list tasks -t      # Display tasks in a tree-like structure
//   clido list projects -j   # Output projects in JSON format
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [projects|tasks]",  // Specifies the options: 'projects' or 'tasks'
		Short: "List projects or tasks",  // Short description of the command
		Long:  `List all projects or tasks, optionally filtered by project for tasks.`,  // Extended description
		Run: func(cmd *cobra.Command, args []string) {
			// Ensure at least one argument (either 'projects' or 'tasks') is provided
			if len(args) < 1 {
				cmd.Println("Insufficient arguments. Use 'list projects' or 'list tasks'.")
				return
			}

			// Initialize the repository for database operations
			repo, err := repository.NewRepository()
			if err != nil {
				cmd.Printf("Error initializing repository: %v\n", err)
				return
			}
			defer repo.Close()

			// Retrieve flags for output format
			outputJSON, _ := cmd.Flags().GetBool("json")
			treeView, _ := cmd.Flags().GetBool("tree")

			// Determine whether to list projects or tasks
			switch args[0] {
			case "projects":
				listProjects(cmd, repo, outputJSON, treeView)
			case "tasks":
				projectFilter, _ := cmd.Flags().GetString("project")
				listTasks(cmd, repo, projectFilter, outputJSON, treeView)
			default:
				cmd.Println("Invalid option. Use 'list projects' or 'list tasks'.")
			}
		},
	}

	// Define flags for the list command
	cmd.Flags().StringP("project", "p", "", "Filter tasks by project name or ID")  // Filter tasks by project
	cmd.Flags().BoolP("json", "j", false, "Output list in JSON format")  // Output as JSON
	cmd.Flags().BoolP("tree", "t", false, "Display projects or tasks in a tree-like structure")  // Display in tree view

	return cmd
}

// listProjects lists all projects in either table, tree view, or JSON format.
func listProjects(cmd *cobra.Command, repo *repository.Repository, outputJSON bool, treeView bool) {
	projects, err := repo.GetAllProjects()
	if err != nil {
		cmd.Printf("Error listing projects: %v\n", err)
		return
	}

	switch {
	case outputJSON:
		// Output projects in JSON format
		var jsonData []byte
		jsonData, err = json.MarshalIndent(projects, "", "  ")
		if err != nil {
			cmd.Printf("Error marshalling projects to JSON: %v\n", err)
			return
		}
		cmd.Println(string(jsonData))

	case treeView:
		// Display projects in tree view
		printProjectTree(cmd, projects, nil, 0)

	default:
		// Display projects in a table
		printProjectTable(cmd, repo, projects)
	}
}

// listTasks lists tasks, optionally filtered by a project, in table, tree view, or JSON format.
func listTasks(
	cmd *cobra.Command,
	repo *repository.Repository,
	projectFilter string,
	outputJSON bool,
	treeView bool,
) {
	tasks, project, err := getTasks(repo, projectFilter)
	if err != nil {
		cmd.Println(err)
		return
	}

	if !outputJSON {
		// Print header for tasks
		printTaskHeader(cmd, project)
	}

	switch {
	case outputJSON:
		// Output tasks in JSON format
		printTasksJSON(cmd, tasks)
	case treeView:
		// Display tasks in tree view
		printTaskTree(cmd, tasks, nil, 0)
	default:
		// Display tasks in a table
		printTaskTable(repo, tasks)
	}
}

// getTasks retrieves tasks filtered by a project or returns all tasks if no filter is provided.
func getTasks(
	repo *repository.Repository,
	projectFilter string,
) ([]*models.Task, *models.Project, error) {
	if projectFilter == "" {
		// Return all tasks if no project filter is specified
		tasks, err := repo.GetAllTasks()
		return tasks, nil, err
	}

	// Get tasks filtered by project
	project, err := getProject(repo, projectFilter)
	if err != nil {
		return nil, nil, err
	}

	tasks, err := repo.GetTasksByProjectID(project.ID)
	return tasks, project, err
}

// getProject retrieves a project by its name or ID, based on the projectFilter input.
func getProject(repo *repository.Repository, projectFilter string) (*models.Project, error) {
	var project *models.Project
	var err error

	if utils.IsNumeric(projectFilter) {
		// Retrieve project by ID
		projectID, _ := strconv.Atoi(projectFilter)
		project, err = repo.GetProjectByID(projectID)
	} else {
		// Retrieve project by name
		project, err = repo.GetProjectByName(projectFilter)
	}

	if err != nil || project == nil {
		return nil, err
	}

	return project, nil
}

// printTaskHeader prints the header for the task list, either all tasks or tasks within a specific project.
func printTaskHeader(cmd *cobra.Command, project *models.Project) {
	if project != nil {
		cmd.Printf("Tasks in project '%s':\n", project.Name)
	} else {
		cmd.Println("All Tasks:")
	}
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
func printProjectTable(
	cmd *cobra.Command,
	repo *repository.Repository,
	projects []*models.Project,
) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Description", "Type", "Child Of"})
	table.SetRowLine(true)

	for _, project := range projects {
		// Determine whether the project is a parent or child project
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

		// Add project details to the table
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
func printTaskTable(repo *repository.Repository, tasks []*models.Task) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"ID", "Name", "Description", "Due Date", "Completed", "Past Due", "Priority", "Project", "Type", "Parent/Child Of",
	})
	table.SetRowLine(true)

	for _, task := range tasks {
		// Determine whether the task is a parent or child task
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

		// Get project name for the task
		project, _ := repo.GetProjectByID(task.ProjectID)
		projectName := ""
		if project != nil {
			projectName = project.Name
		}

		// Add task details to the table
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

// printProjectTree displays projects in a tree view.
func printProjectTree(cmd *cobra.Command, projects []*models.Project, parentID *int, level int) {
	nodes := make([]TreeNode, len(projects))
	for i, p := range projects {
		nodes[i] = ProjectNode{p}
	}
	printTree(cmd, nodes, parentID, level, nil)
}

// printTaskTree displays tasks in a tree view.
func printTaskTree(cmd *cobra.Command, tasks []*models.Task, parentID *int, level int) {
	nodes := make([]TreeNode, len(tasks))
	for i, t := range tasks {
		nodes[i] = TaskNode{t}
	}
	printTree(cmd, nodes, parentID, level, func(node TreeNode, indent string) {
		// Print task details in the tree view
		task := node.(TaskNode).Task
		cmd.Printf("%sDescription: %s\n", indent, task.Description)
		cmd.Printf(
			"%sDue Date: %s, Completed: %v, Priority: %s\n",
			indent,
			utils.FormatDate(task.DueDate),
			task.TaskCompleted,
			utils.GetPriorityString(task.Priority),
		)
	})
}

// TreeNode represents a node in a tree structure (for both projects and tasks).
type TreeNode interface {
	GetID() int
	GetParentID() *int
	GetName() string
}

// ProjectNode represents a project in the tree view.
type ProjectNode struct {
	*models.Project
}

func (p ProjectNode) GetID() int        { return p.ID }
func (p ProjectNode) GetParentID() *int { return p.ParentProjectID }
func (p ProjectNode) GetName() string   { return p.Name }

// TaskNode represents a task in the tree view.
type TaskNode struct {
	*models.Task
}

func (t TaskNode) GetID() int        { return t.ID }
func (t TaskNode) GetParentID() *int { return t.ParentTaskID }
func (t TaskNode) GetName() string   { return t.Name }

// printTree prints the tree structure for projects or tasks.
func printTree(
	cmd *cobra.Command,
	nodes []TreeNode,
	parentID *int,
	level int,
	printDetails func(TreeNode, string),
) {
	indent := strings.Repeat("│  ", level)
	for i, node := range nodes {
		if (parentID == nil && node.GetParentID() == nil) ||
			(parentID != nil && node.GetParentID() != nil && *node.GetParentID() == *parentID) {
			// Use appropriate tree symbols for formatting
			prefix := "├──"
			if i == len(nodes)-1 {
				prefix = "└──"
			}
			cmd.Printf("%s%s %s (ID: %d)\n", indent, prefix, node.GetName(), node.GetID())
			if printDetails != nil {
				printDetails(node, indent+"    ")
			}
			nodeID := node.GetID()
			// Recursively print child nodes
			printTree(cmd, nodes, &nodeID, level+1, printDetails)
		}
	}
}
