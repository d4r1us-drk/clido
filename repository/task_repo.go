package repository

import (
	"github.com/d4r1us-drk/clido/models"
)

// CreateTask inserts a new task into the database.
func (r *Repository) CreateTask(task *models.Task) error {
	return r.db.Create(task).Error
}

// GetTaskByID retrieves a task from the database by its ID.
func (r *Repository) GetTaskByID(id int) (*models.Task, error) {
	var task models.Task
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetAllTasks retrieves all tasks from the database.
func (r *Repository) GetAllTasks() ([]*models.Task, error) {
	var tasks []*models.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

// GetTasksByProjectID retrieves all tasks associated with a specific project by the project's ID.
func (r *Repository) GetTasksByProjectID(projectID int) ([]*models.Task, error) {
	var tasks []*models.Task
	err := r.db.Where("project_id = ?", projectID).Find(&tasks).Error
	return tasks, err
}

// GetSubtasks retrieves all subtasks that have the given parent task ID.
func (r *Repository) GetSubtasks(parentTaskID int) ([]*models.Task, error) {
	var tasks []*models.Task
	err := r.db.Where("parent_task_id = ?", parentTaskID).Find(&tasks).Error
	return tasks, err
}

// UpdateTask updates an existing task in the database.
func (r *Repository) UpdateTask(task *models.Task) error {
	return r.db.Save(task).Error
}

// DeleteTask removes a task from the database by its ID.
func (r *Repository) DeleteTask(id int) error {
	return r.db.Delete(&models.Task{}, id).Error
}

// GetNextTaskID retrieves the next available task ID in the database.
// It selects the maximum task ID and adds 1 to determine the next available ID.
func (r *Repository) GetNextTaskID() (int, error) {
	var maxID int
	err := r.db.Model(&models.Task{}).Select("COALESCE(MAX(id), 0) + 1").Scan(&maxID).Error
	return maxID, err
}
