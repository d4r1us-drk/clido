package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func handleNew(repo *Repository, args []string) {
	if len(args) < 1 {
		color.Red("Insufficient arguments for 'new' command.")
		return
	}

	if args[0] == "-p" {
		name := getValueFromArgs(args, "-n")
		description := getValueFromArgs(args, "-d")
		project := &Project{
			Name:        name,
			Description: description,
		}
		err := repo.CreateProject(project)
		if err != nil {
			color.Red("Error creating project: %v\n", err)
			return
		}
		color.Green("Project '%s' created successfully.\n", name)
	} else if args[0] == "-t" {
		name := getValueFromArgs(args, "-n")
		description := getValueFromArgs(args, "-d")
		dueDateStr := getValueFromArgs(args, "-D")
		var dueDate *time.Time
		if dueDateStr != "" {
			parsedDate, err := time.Parse("2006-01-02 15:04", dueDateStr)
			if err == nil {
				dueDate = &parsedDate
			}
		}

		var project *Project
		var err error
		if strings.Contains(strings.Join(args, " "), "-p") {
			projectName := getValueFromArgs(args, "-p")
			project, err = repo.GetProjectByName(projectName)
			if err != nil || project == nil {
				color.Red("Project '%s' not found.\n", projectName)
				return
			}
		} else {
			projectID, _ := strconv.Atoi(getValueFromArgs(args, "-P"))
			project, err = repo.GetProjectByID(projectID)
			if err != nil || project == nil {
				color.Red("Project with ID %d not found.\n", projectID)
				return
			}
		}

		task := &Task{
			Name:        name,
			Description: description,
			ProjectID:   project.ID,
			DueDate:     dueDate,
		}
		err = repo.CreateTask(task)
		if err != nil {
			color.Red("Error creating task: %v\n", err)
			return
		}
		color.Green("Task '%s' created successfully.\n", name)
	}
}

func handleEdit(repo *Repository, args []string) {
	if len(args) < 1 {
		color.Red("Insufficient arguments for 'edit' command.")
		return
	}

	if args[0] == "-p" {
		id, _ := strconv.Atoi(getValueFromArgs(args, "-i"))
		newName := getValueFromArgs(args, "-n")
		newDescription := getValueFromArgs(args, "-d")

		project, err := repo.GetProjectByID(id)
		if err != nil || project == nil {
			color.Red("Project with ID %d not found.\n", id)
			return
		}

		project.Name = newName
		project.Description = newDescription
		err = repo.UpdateProject(project)
		if err != nil {
			color.Red("Error updating project: %v\n", err)
			return
		}
		color.Green("Project '%s' updated successfully.\n", project.Name)
	} else if args[0] == "-t" {
		id, _ := strconv.Atoi(getValueFromArgs(args, "-i"))
		newName := getValueFromArgs(args, "-n")
		newDescription := getValueFromArgs(args, "-d")
		newDueDateStr := getValueFromArgs(args, "-D")
		var newDueDate *time.Time
		if newDueDateStr != "" {
			parsedDate, err := time.Parse("2006-01-02 15:04", newDueDateStr)
			if err == nil {
				newDueDate = &parsedDate
			}
		}

		task, err := repo.GetTaskByID(id)
		if err != nil || task == nil {
			color.Red("Task with ID %d not found.\n", id)
			return
		}

		task.Name = newName
		task.Description = newDescription
		task.DueDate = newDueDate
		err = repo.UpdateTask(task)
		if err != nil {
			color.Red("Error updating task: %v\n", err)
			return
		}
		color.Green("Task '%s' updated successfully.\n", task.Name)
	}
}

func handleList(repo *Repository, args []string) {
	if len(args) < 1 {
		color.Red("Insufficient arguments for 'list' command.")
		return
	}

	if args[0] == "-p" {
		projects, err := repo.GetAllProjects()
		if err != nil {
			color.Red("Error listing projects: %v\n", err)
			return
		}

		color.Cyan("Projects:")
		for _, project := range projects {
			fmt.Printf("  ID: %d, Name: %s, Description: %s\n", project.ID, project.Name, project.Description)
		}
	} else if args[0] == "-t" {
		if strings.Contains(strings.Join(args, " "), "-p") {
			projectName := getValueFromArgs(args, "-p")
			project, err := repo.GetProjectByName(projectName)
			if err != nil || project == nil {
				color.Red("Project '%s' not found.\n", projectName)
				return
			}

			tasks, err := repo.GetTasksByProjectID(project.ID)
			if err != nil {
				color.Red("Error listing tasks: %v\n", err)
				return
			}

			color.Cyan("Tasks:")
			for _, task := range tasks {
				fmt.Printf("  ID: %d, Name: %s, Description: %s, Project ID: %d\n", task.ID, task.Name, task.Description, task.ProjectID)
			}
		} else if strings.Contains(strings.Join(args, " "), "-P") {
			projectID, _ := strconv.Atoi(getValueFromArgs(args, "-P"))
			tasks, err := repo.GetTasksByProjectID(projectID)
			if err != nil {
				color.Red("Error listing tasks: %v\n", err)
				return
			}

			color.Cyan("Tasks:")
			for _, task := range tasks {
				fmt.Printf("  ID: %d, Name: %s, Description: %s, Project ID: %d\n", task.ID, task.Name, task.Description, task.ProjectID)
			}
		} else {
			tasks, err := repo.GetAllTasks()
			if err != nil {
				color.Red("Error listing tasks: %v\n", err)
				return
			}

			color.Cyan("Tasks:")
			for _, task := range tasks {
				fmt.Printf("  ID: %d, Name: %s, Description: %s, Project ID: %d\n", task.ID, task.Name, task.Description, task.ProjectID)
			}
		}
	}
}

func handleRemove(repo *Repository, args []string) {
	if len(args) < 1 {
		color.Red("Insufficient arguments for 'remove' command.")
		return
	}

	if args[0] == "-p" {
		id, _ := strconv.Atoi(getValueFromArgs(args, "-i"))
		err := repo.DeleteProject(id)
		if err != nil {
			color.Red("Error removing project: %v\n", err)
			return
		}
		color.Green("Project with ID %d removed successfully.\n", id)
	} else if args[0] == "-t" {
		id, _ := strconv.Atoi(getValueFromArgs(args, "-i"))
		err := repo.DeleteTask(id)
		if err != nil {
			color.Red("Error removing task: %v\n", err)
			return
		}
		color.Green("Task with ID %d removed successfully.\n", id)
	}
}

func handleSet(repo *Repository, args []string) {
	if len(args) < 1 {
		color.Red("Insufficient arguments for 'set' command.")
		return
	}

	id, _ := strconv.Atoi(getValueFromArgs(args, "-i"))
	task, err := repo.GetTaskByID(id)
	if err != nil || task == nil {
		color.Red("Task with ID %d not found.\n", id)
		return
	}

	task.TaskCompleted = !task.TaskCompleted
	task.LastUpdatedDate = time.Now()
	err = repo.UpdateTask(task)
	if err != nil {
		color.Red("Error updating task: %v\n", err)
		return
	}
	color.Green("Task '%s' marked as %v.\n", task.Name, task.TaskCompleted)
}

func getValueFromArgs(args []string, key string) string {
	for i, arg := range args {
		if arg == key && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func handleHelp() {
	fmt.Println(`CLI Todo Application

Usage:
  cli-todo <command> [arguments]

Commands:
  new      Create a new project or task
  edit     Edit an existing project or task
  list     List projects or tasks
  remove   Remove a project or task
  set      Toggle task completion
  help     Show this help message

Options:
  new:
    -p                   Create a new project
      -n <name>           Project name (required)
      -d <description>    Project description (optional)
    -t                   Create a new task
      -n <name>           Task name (required)
      -d <description>    Task description (optional)
      -D <dueDate>        Task due date (optional, format: YYYY-MM-DD HH:MM)
      -p <projectName>    Project name (required if -P is not used)
      -P <projectID>      Project ID (required if -p is not used)

  edit:
    -p                   Edit an existing project
      -i <projectID>      Project ID (required)
      -n <newName>        New project name (optional)
      -d <newDescription> New project description (optional)
    -t                   Edit an existing task
      -i <taskID>         Task ID (required)
      -n <newName>        New task name (optional)
      -d <newDescription> New task description (optional)
      -D <newDueDate>     New task due date (optional, format: YYYY-MM-DD HH:MM)

  list:
    -p                   List all projects
    -t                   List all tasks
      -p <projectName>    Filter tasks by project name (optional)
      -P <projectID>      Filter tasks by project ID (optional)

  remove:
    -p                   Remove a project
      -i <projectID>      Project ID (required)
    -t                   Remove a task
      -i <taskID>         Task ID (required)

  set:
    -i <taskID>           Task ID (required) to toggle completion

Examples:
  cli-todo new -p -n "New Project" -d "Project Description"
  cli-todo new -t -n "New Task" -d "Task Description" -D "2024-08-15 23:00" -p "Existing Project"
  cli-todo edit -p -i 1 -n "Updated Project Name" -d "Updated Description"
  cli-todo list -p
  cli-todo list -t -P 1
  cli-todo remove -p -i 1
  cli-todo set -i 1
`)
}
