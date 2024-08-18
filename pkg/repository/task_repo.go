package repository

import (
	"github.com/d4r1us-drk/clido/pkg/models"
)

func (r *Repository) CreateTask(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *Repository) GetTaskByID(id int) (*models.Task, error) {
	var task models.Task
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *Repository) GetAllTasks() ([]*models.Task, error) {
	var tasks []*models.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *Repository) GetTasksByProjectID(projectID int) ([]*models.Task, error) {
	var tasks []*models.Task
	err := r.db.Where("project_id = ?", projectID).Find(&tasks).Error
	return tasks, err
}

func (r *Repository) GetSubtasks(parentTaskID int) ([]*models.Task, error) {
	var tasks []*models.Task
	err := r.db.Where("parent_task_id = ?", parentTaskID).Find(&tasks).Error
	return tasks, err
}

func (r *Repository) UpdateTask(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *Repository) DeleteTask(id int) error {
	return r.db.Delete(&models.Task{}, id).Error
}

func (r *Repository) GetNextTaskID() (int, error) {
	var maxID int
	err := r.db.Model(&models.Task{}).Select("COALESCE(MAX(id), 0) + 1").Scan(&maxID).Error
	return maxID, err
}
