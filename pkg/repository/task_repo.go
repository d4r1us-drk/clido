package repository

import (
	"github.com/d4r1us-drk/clido/pkg/models"
)

// CreateTask inserts a new task into the database.
//
// Parameters:
//   - task: A pointer to the task model to be created.
//
// Returns:
//   - An error if the operation fails; nil if successful.
func (r *Repository) CreateTask(task *models.Task) error {
	return r.db.Create(task).Error
}

// GetTaskByID retrieves a task from the database by its ID.
//
// Parameters:
//   - id: The unique ID of the task to retrieve.
//
// Returns:
//   - A pointer to the retrieved task if found.
//   - An error if the task could not be found or another issue occurred.
func (r *Repository) GetTaskByID(id int) (*models.Task, error) {
	var task models.Task
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetAllTasks retrieves all tasks from the database.
//
// Returns:
//   - A slice of pointers to all retrieved tasks.
//   - An error if the operation fails.
func (r *Repository) GetAllTasks() ([]*models.Task, error) {
	var tasks []*models.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

// GetTasksByProjectID retrieves all tasks associated with a specific project by the project's ID.
//
// Parameters:
//   - projectID: The ID of the project to retrieve tasks for.
//
// Returns:
//   - A slice of pointers to all tasks associated with the specified project.
//   - An error if the operation fails.
func (r *Repository) GetTasksByProjectID(projectID int) ([]*models.Task, error) {
	var tasks []*models.Task
	err := r.db.Where("project_id = ?", projectID).Find(&tasks).Error
	return tasks, err
}

// GetSubtasks retrieves all subtasks that have the given parent task ID.
//
// Parameters:
//   - parentTaskID: The ID of the parent task to retrieve subtasks for.
//
// Returns:
//   - A slice of pointers to all subtasks under the specified parent task.
//   - An error if the operation fails.
func (r *Repository) GetSubtasks(parentTaskID int) ([]*models.Task, error) {
	var tasks []*models.Task
	err := r.db.Where("parent_task_id = ?", parentTaskID).Find(&tasks).Error
	return tasks, err
}

// UpdateTask updates an existing task in the database.
//
// Parameters:
//   - task: A pointer to the task model to be updated.
//
// Returns:
//   - An error if the update operation fails; nil if successful.
func (r *Repository) UpdateTask(task *models.Task) error {
	return r.db.Save(task).Error
}

// DeleteTask removes a task from the database by its ID.
//
// Parameters:
//   - id: The unique ID of the task to be deleted.
//
// Returns:
//   - An error if the deletion fails; nil if successful.
func (r *Repository) DeleteTask(id int) error {
	return r.db.Delete(&models.Task{}, id).Error
}

// GetNextTaskID retrieves the next available task ID in the database.
// It selects the maximum task ID and adds 1 to determine the next available ID.
//
// Returns:
//   - The next available task ID as an integer.
//   - An error if the operation fails.
func (r *Repository) GetNextTaskID() (int, error) {
	var maxID int
	err := r.db.Model(&models.Task{}).Select("COALESCE(MAX(id), 0) + 1").Scan(&maxID).Error
	return maxID, err
}
