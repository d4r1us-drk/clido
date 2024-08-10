package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func handleNew(repo *Repository, args []string) {
	if len(args) < 1 {
		color.Red("Insufficient arguments for 'new' command.")
		return
	}

	switch args[0] {
	case "project":
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

	case "task":
		name := getValueFromArgs(args, "-n")
		description := getValueFromArgs(args, "-d")
		dueDateStr := getValueFromArgs(args, "-D")
		priorityStr := getValueFromArgs(args, "-p")
		var dueDate *time.Time
		if dueDateStr != "" {
			parsedDate, err := time.Parse("2006-01-02 15:04", dueDateStr)
			if err == nil {
				dueDate = &parsedDate
			}
		}
		priority, err := strconv.Atoi(priorityStr)
		if err != nil || priority < 1 || priority > 4 {
			priority = 4 // default priority
		}

		projectIdentifier := getValueFromArgs(args, "-p")
		if projectIdentifier == "" {
			color.Red("Project name or number is required for creating a task.")
			return
		}

		var project *Project
		if isNumeric(projectIdentifier) {
			projectID, _ := strconv.Atoi(projectIdentifier)
			project, err = repo.GetProjectByID(projectID)
			if err != nil || project == nil {
				color.Red("Project with number %d not found.\n", projectID)
				return
			}
		} else {
			project, err = repo.GetProjectByName(projectIdentifier)
			if err != nil || project == nil {
				color.Red("Project '%s' not found.\n", projectIdentifier)
				return
			}
		}

		task := &Task{
			Name:        name,
			Description: description,
			ProjectID:   project.ID,
			DueDate:     dueDate,
			Priority:    priority,
		}
		err = repo.CreateTask(task)
		if err != nil {
			color.Red("Error creating task: %v\n", err)
			return
		}
		color.Green("Task '%s' created successfully with priority %d.\n", name, priority)

	default:
		color.Red("Invalid option for 'new' command.")
		handleHelp()
	}
}

func handleEdit(repo *Repository, args []string) {
	if len(args) < 2 {
		color.Red("Insufficient arguments for 'edit' command.")
		return
	}

	switch args[0] {
	case "project":
		id, _ := strconv.Atoi(args[1])
		newName := getValueFromArgs(args, "-n")
		newDescription := getValueFromArgs(args, "-d")

		project, err := repo.GetProjectByID(id)
		if err != nil || project == nil {
			color.Red("Project with number %d not found.\n", id)
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

	case "task":
		id, _ := strconv.Atoi(args[1])
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
			color.Red("Task with number %d not found.\n", id)
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

	default:
		color.Red("Invalid option for 'edit' command.")
		handleHelp()
	}
}

func handleList(repo *Repository, args []string) {
	if len(args) < 1 {
		color.Red("Insufficient arguments for 'list' command.")
		return
	}

	switch args[0] {
	case "projects":
		projects, err := repo.GetAllProjects()
		if err != nil {
			color.Red("Error listing projects: %v\n", err)
			return
		}

		color.Cyan("Projects:")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Number", "Name", "Description"})
		table.SetRowLine(true)

		for _, project := range projects {
			table.Append([]string{strconv.Itoa(project.ID), wrapText(project.Name, 20), wrapText(project.Description, 30)})
		}
		table.Render()

	case "tasks":
		if len(args) > 2 && (args[1] == "-p" || args[1] == "-P") {
			projectIdentifier := args[2]
			var project *Project
			var err error
			if isNumeric(projectIdentifier) {
				projectID, _ := strconv.Atoi(projectIdentifier)
				project, err = repo.GetProjectByID(projectID)
				if err != nil || project == nil {
					color.Red("Project with number %d not found.\n", projectID)
					return
				}
			} else {
				project, err = repo.GetProjectByName(projectIdentifier)
				if err != nil || project == nil {
					color.Red("Project '%s' not found.\n", projectIdentifier)
					return
				}
			}

			tasks, err := repo.GetTasksByProjectID(project.ID)
			if err != nil {
				color.Red("Error listing tasks: %v\n", err)
				return
			}

			color.Cyan("Tasks in project '%s':", project.Name)
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Number", "Name", "Description", "Due Date", "Is Completed", "Is Past Due"})
			table.SetRowLine(true)

			for _, task := range tasks {
				isCompleted := "no"
				if task.TaskCompleted {
					isCompleted = "yes"
				}

				var dueDate string
				isPastDue := color.GreenString("no")
				if task.DueDate != nil {
					dueDate = task.DueDate.Format("2006-01-02 15:04")
					if task.DueDate.Before(time.Now()) {
						if !task.TaskCompleted {
							isPastDue = color.RedString("yes")
						} else {
							isPastDue = color.GreenString("yes")
						}
					}
				} else {
					dueDate = "None"
				}

				table.Append([]string{
					strconv.Itoa(task.ID),
					wrapText(task.Name, 20),
					wrapText(task.Description, 30),
					wrapText(dueDate, 20),
					isCompleted,
					isPastDue,
				})
			}
			table.Render()

		} else {
			tasks, err := repo.GetAllTasks()
			if err != nil {
				color.Red("Error listing tasks: %v\n", err)
				return
			}

			color.Cyan("Tasks:")
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Number", "Name", "Description", "Due Date", "Is Completed", "Is Past Due", "Project"})
			table.SetRowLine(true)

			for _, task := range tasks {
				isCompleted := "no"
				if task.TaskCompleted {
					isCompleted = "yes"
				}

				var dueDate string
				isPastDue := color.GreenString("no")
				if task.DueDate != nil {
					dueDate = task.DueDate.Format("2006-01-02 15:04")
					if task.DueDate.Before(time.Now()) {
						if !task.TaskCompleted {
							isPastDue = color.RedString("yes")
						} else {
							isPastDue = color.GreenString("yes")
						}
					}
				} else {
					dueDate = "None"
				}

				project, err := repo.GetProjectByID(task.ProjectID)
				if err == nil {
					table.Append([]string{
						strconv.Itoa(task.ID),
						wrapText(task.Name, 20),
						wrapText(task.Description, 30),
						wrapText(dueDate, 20),
						isCompleted,
						isPastDue,
						wrapText(project.Name, 20),
					})
				}
			}
			table.Render()
		}

	default:
		color.Red("Invalid option for 'list' command.")
		handleHelp()
	}
}

func handleRemove(repo *Repository, args []string) {
	if len(args) < 2 {
		color.Red("Insufficient arguments for 'remove' command.")
		return
	}

	switch args[0] {
	case "project":
		id, _ := strconv.Atoi(args[1])
		err := repo.DeleteProject(id)
		if err != nil {
			color.Red("Error removing project: %v\n", err)
			return
		}
		color.Green("Project with number %d removed successfully.\n", id)

	case "task":
		id, _ := strconv.Atoi(args[1])
		err := repo.DeleteTask(id)
		if err != nil {
			color.Red("Error removing task: %v\n", err)
			return
		}
		color.Green("Task with number %d removed successfully.\n", id)

	default:
		color.Red("Invalid option for 'remove' command.")
		handleHelp()
	}
}

func handleToggle(repo *Repository, args []string) {
	if len(args) < 1 {
		color.Red("Insufficient arguments for 'set-completed' command.")
		return
	}

	id, _ := strconv.Atoi(args[1])
	task, err := repo.GetTaskByID(id)
	if err != nil || task == nil {
		color.Red("Task with number %d not found.\n", id)
		return
	}

	task.TaskCompleted = !task.TaskCompleted
	task.LastUpdatedDate = time.Now()
	err = repo.UpdateTask(task)
	if err != nil {
		color.Red("Error updating task: %v\n", err)
		return
	}
	message := "uncompleted"
	if task.TaskCompleted {
		message = "completed"
	}
	color.Green("Task '%s' marked as %v.\n", task.Name, message)
}

func getValueFromArgs(args []string, key string) string {
	for i, arg := range args {
		if arg == key && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func wrapText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}

	var result string
	words := strings.Fields(text)
	line := ""

	for _, word := range words {
		if len(line)+len(word)+1 > maxLength {
			if len(result) > 0 {
				result += "\n"
			}
			result += line
			line = word
		} else {
			if len(line) > 0 {
				line += " "
			}
			line += word
		}
	}

	if len(line) > 0 {
		if len(result) > 0 {
			result += "\n"
		}
		result += line
	}

	return result
}

func handleHelp() {
	// Define colors
	headerColor := color.New(color.FgMagenta).SprintFunc()
	commandColor := color.New(color.FgGreen).SprintFunc()
	optionColor := color.New(color.FgYellow).SprintFunc()
	argColor := color.New(color.FgCyan).SprintFunc()
	valueColor := color.New(color.FgWhite).SprintFunc()

	helpMessage := `
    clido: An awesome cli to-do list management application

    ` + headerColor("Usage :") + `
        cli-todo ` + commandColor("<command>") + ` ` + optionColor("[options]") + ` ` + argColor("[params]") + ` ` + valueColor("value") + `

    ` + headerColor("Commands :") + `
        ` + commandColor("new") + `      Create a new project or task
        ` + commandColor("edit") + `     Edit an existing project or task
        ` + commandColor("list") + `     List projects or tasks
        ` + commandColor("remove") + `   Remove a project or task
        ` + commandColor("toggle") + `   Toggle task completion
        ` + commandColor("help") + `     Show this help message

    ` + headerColor("Options :") + `
        ` + commandColor("new") + ` :
            ` + optionColor("project") + ` :
                ` + argColor("-n") + ` ` + valueColor("<name>") + `                        Project name (required)
                ` + argColor("-d") + ` ` + valueColor("<description>") + `                 Project description (optional)
            ` + optionColor("task") + ` :
                ` + argColor("-n") + ` ` + valueColor("<name>") + `                        Task name (required)
                ` + argColor("-d") + ` ` + valueColor("<description>") + `                 Task description (optional)
                ` + argColor("-D") + ` ` + valueColor("<dueDate>") + `                     Task due date (optional, format: YYYY-MM-DD HH:MM)
                ` + argColor("-p") + ` ` + valueColor("<projectName>/<projectNumber>") + ` Project name or number (required)
        ` + commandColor("edit") + ` :
            ` + optionColor("project") + ` ` + valueColor("<projectNumber>") + ` :
                ` + argColor("-n") + ` ` + valueColor("<newName>") + `                     New project name (optional)
                ` + argColor("-d") + ` ` + valueColor("<newDescription>") + `              New project description (optional)
            ` + optionColor("task") + ` ` + valueColor("<taskNumber>") + ` :
                ` + argColor("-n") + ` ` + valueColor("<newName>") + `                     New task name (optional)
                ` + argColor("-d") + ` ` + valueColor("<newDescription>") + `              New task description (optional)
                ` + argColor("-D") + ` ` + valueColor("<newDueDate>") + `                  New task due date (optional, format: YYYY-MM-DD HH:MM)
        ` + commandColor("list") + ` :
            ` + optionColor("projects") + `                             List all projects
            ` + optionColor("tasks") + `
                ` + argColor("-p") + ` ` + valueColor("<projectName>/<projectNumber>") + ` Filter tasks by project name or number (optional)
        ` + commandColor("remove") + ` :
            ` + optionColor("project") + ` ` + valueColor("<projectNumber>") + `              Remove a project
            ` + optionColor("task") + ` ` + valueColor("<taskNumber>") + `                    Remove a task
        ` + commandColor("toggle") + ` ` + valueColor("<taskNumber>") + `                      Set a task as completed

    ` + headerColor("Examples :") + `
        cli-todo new project -n "New Project" -d "Project Description"
        cli-todo new task -n "New Task" -d "Task Description" -D "2024-08-15 23:00" -p "Existing Project"
        cli-todo edit project 1 -n "Updated Project Name" -d "Updated Description"
        cli-todo list projects
        cli-todo list tasks -P 1
        cli-todo remove project 1
        cli-todo toggle 1`

	fmt.Println(helpMessage)
}
