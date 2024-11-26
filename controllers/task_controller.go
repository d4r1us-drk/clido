package controllers

import (
	"errors"
	"time"

	"github.com/d4r1us-drk/clido/models"
	"github.com/d4r1us-drk/clido/repository"
	"github.com/d4r1us-drk/clido/utils"
)

// Error constants for task operations.
var (
	ErrNoTaskName         = errors.New("task name is required")
	ErrNoProject          = errors.New("project name or numeric ID is required")
	ErrNoProjectFound     = errors.New("project not found")
	ErrInvalidParentTask  = errors.New("parent task must be identified by a numeric ID")
	ErrInvalidDueDate     = errors.New("invalid due date format")
	ErrTaskNotFound       = errors.New("task not found")
	ErrParentTaskNotFound = errors.New("parent task not found")
)

// TaskController manages the task-related business logic.
type TaskController struct {
	repo *repository.Repository
}

// NewTaskController creates and returns a new instance of TaskController.
func NewTaskController(repo *repository.Repository) *TaskController {
	return &TaskController{repo: repo}
}

// CreateTask handles the creation of a new task.
func (tc *TaskController) CreateTask(
	name, description, projectIdentifier, parentTaskIdentifier, dueDateStr string,
	priority int,
) error {
	// Validate mandatory arguments
	if name == "" {
		return ErrNoTaskName
	}
	if projectIdentifier == "" {
		return ErrNoProject
	}

	// Try to parse projectIdentifier as an integer (ID), otherwise get the project by name
	projectID, projectErr := utils.ParseIntOrError(projectIdentifier)
	if projectErr != nil {
		project, lookupErr := tc.repo.GetProjectByName(projectIdentifier)
		if lookupErr != nil || project == nil {
			return ErrNoProjectFound
		}
		projectID = project.ID
	} else {
		// Check if the project exists using the ID
		project, getProjectErr := tc.repo.GetProjectByID(projectID)
		if getProjectErr != nil || project == nil {
			return ErrNoProjectFound
		}
	}

	// Get parent task ID (optional)
	var parentTaskID *int
	if parentTaskIdentifier != "" {
		id, parentErr := utils.ParseIntOrError(parentTaskIdentifier)
		if parentErr != nil {
			return ErrInvalidParentTask
		}
		parentTaskID = &id
	}

	// Parse due date (optional)
	var dueDate *time.Time
	if dueDateStr != "" {
		parsedDate, dueDateErr := utils.ParseDueDate(dueDateStr)
		if dueDateErr != nil {
			return ErrInvalidDueDate
		}
		dueDate = parsedDate
	}

	// Create a new task
	task := &models.Task{
		Name:         name,
		Description:  description,
		ProjectID:    projectID,
		DueDate:      dueDate,
		Priority:     priority,
		ParentTaskID: parentTaskID,
	}

	// Store the task in the repository
	if createErr := tc.repo.CreateTask(task); createErr != nil {
		return createErr
	}

	return nil
}

// EditTask handles updating an existing task by its ID.
func (tc *TaskController) EditTask(
	id int,
	name, description, dueDateStr string,
	priority int,
	parentTaskIdentifier string,
) error {
	task, getTaskErr := tc.repo.GetTaskByID(id)
	if getTaskErr != nil {
		return ErrTaskNotFound
	}

	// Apply updates
	if name != "" {
		task.Name = name
	}
	if description != "" {
		task.Description = description
	}
	if dueDateStr != "" {
		dueDate, dueDateErr := utils.ParseDueDate(dueDateStr)
		if dueDateErr != nil {
			return ErrInvalidDueDate
		}
		task.DueDate = dueDate
	}
	if priority != 0 {
		task.Priority = priority
	}
	if parentTaskIdentifier != "" {
		parentTaskID, parentErr := utils.ParseIntOrError(parentTaskIdentifier)
		if parentErr != nil {
			return ErrInvalidParentTask
		}
		task.ParentTaskID = &parentTaskID
	}

	// Update the task in the repository
	if updateErr := tc.repo.UpdateTask(task); updateErr != nil {
		return updateErr
	}

	return nil
}

// ListTasks returns all tasks stored in the repository.
func (tc *TaskController) ListTasks() ([]*models.Task, error) {
	tasks, getAllErr := tc.repo.GetAllTasks()
	if getAllErr != nil {
		return nil, getAllErr
	}
	return tasks, nil
}

// ListTasksByProjectFilter returns tasks filtered by project.
func (tc *TaskController) ListTasksByProjectFilter(
	projectFilter string,
) ([]*models.Task, *models.Project, error) {
	if projectFilter == "" {
		// Return all tasks if no filter is provided
		tasks, getAllErr := tc.repo.GetAllTasks()
		return tasks, nil, getAllErr
	}

	// Try to parse the projectFilter as a numeric ID first
	projectID, projectErr := utils.ParseIntOrError(projectFilter)
	if projectErr != nil {
		// If parsing fails, assume it's a project name and get the project by name
		project, lookupErr := tc.repo.GetProjectByName(projectFilter)
		if lookupErr != nil || project == nil {
			return nil, nil, ErrNoProjectFound
		}
		projectID = project.ID
	}

	// Retrieve the project by ID
	project, getProjectErr := tc.repo.GetProjectByID(projectID)
	if getProjectErr != nil || project == nil {
		return nil, nil, ErrNoProjectFound
	}

	// Get tasks by project ID
	tasks, getTasksErr := tc.repo.GetTasksByProjectID(project.ID)
	if getTasksErr != nil {
		return nil, nil, getTasksErr
	}

	return tasks, project, nil
}

// ToggleTaskCompletion toggles the completion status of a task.
// If recursive is true, it also toggles the completion status of all subtasks.
func (tc *TaskController) ToggleTaskCompletion(id int, recursive bool) (string, error) {
	// Retrieve the task by its ID
	task, getTaskErr := tc.repo.GetTaskByID(id)
	if getTaskErr != nil {
		return "", ErrTaskNotFound
	}

	// Toggle the completion status
	completion := "not completed"
	if task.TaskCompleted {
		task.TaskCompleted = false
		task.CompletionDate = nil
	} else {
		task.TaskCompleted = true
		now := time.Now()
		task.CompletionDate = &now
		completion = "completed"
	}

	// Update the task in the repository
	updateErr := tc.repo.UpdateTask(task)
	if updateErr != nil {
		return "", updateErr
	}

	// If recursive flag is set, toggle the completion status of all subtasks
	if recursive {
		subtasks, getSubtasksErr := tc.ListSubtasks(id)
		if getSubtasksErr != nil {
			return "", getSubtasksErr
		}

		for _, subtask := range subtasks {
			if _, toggleErr := tc.ToggleTaskCompletion(subtask.ID, true); toggleErr != nil {
				return "", toggleErr
			}
		}
	}

	return completion, nil
}

// RemoveTask handles the recursive removal of a task and all its subtasks.
func (tc *TaskController) RemoveTask(id int) error {
	// Get all subtasks for the given task
	subtasks, getSubtasksErr := tc.repo.GetSubtasks(id)
	if getSubtasksErr != nil {
		return getSubtasksErr
	}

	// Recursively remove subtasks
	for _, subtask := range subtasks {
		if removeErr := tc.RemoveTask(subtask.ID); removeErr != nil {
			return removeErr
		}
	}

	// Remove the parent task
	if deleteErr := tc.repo.DeleteTask(id); deleteErr != nil {
		return deleteErr
	}

	return nil
}

// GetTaskByID returns the task details for a given task ID.
func (tc *TaskController) GetTaskByID(id int) (*models.Task, error) {
	task, getTaskErr := tc.repo.GetTaskByID(id)
	if getTaskErr != nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (tc *TaskController) GetTaskProjectName(id int) (*string, error) {
	task, getTaskErr := tc.GetTaskByID(id)
	if getTaskErr != nil {
		return nil, getTaskErr
	}
	return &task.Project.Name, nil
}

// ListSubtasks returns the subtasks for a given task ID.
func (tc *TaskController) ListSubtasks(taskID int) ([]*models.Task, error) {
	subtasks, getSubtasksErr := tc.repo.GetSubtasks(taskID)
	if getSubtasksErr != nil {
		return nil, getSubtasksErr
	}
	return subtasks, nil
}
