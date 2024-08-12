package repository

import (
	"time"

	"github.com/d4r1us-drk/clido/pkg/models"
)

func (r *Repository) CreateTask(task *models.Task) error {
	task.CreationDate = time.Now()
	task.LastUpdatedDate = time.Now()
	result, err := r.db.Exec(
		`INSERT INTO Tasks (Name, Description, ProjectID, TaskCompleted, DueDate, CompletionDate, CreationDate, LastUpdatedDate, Priority) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		task.Name,
		task.Description,
		task.ProjectID,
		task.TaskCompleted,
		task.DueDate,
		task.CompletionDate,
		task.CreationDate,
		task.LastUpdatedDate,
		task.Priority,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = int(id)
	return nil
}

func (r *Repository) GetTaskByID(id int) (*models.Task, error) {
	task := &models.Task{}
	err := r.db.QueryRow(`SELECT ID, Name, Description, ProjectID, TaskCompleted, DueDate, CompletionDate, CreationDate, LastUpdatedDate, Priority FROM Tasks WHERE ID = ?`, id).
		Scan(&task.ID, &task.Name, &task.Description, &task.ProjectID, &task.TaskCompleted, &task.DueDate, &task.CompletionDate, &task.CreationDate, &task.LastUpdatedDate, &task.Priority)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *Repository) GetAllTasks() ([]*models.Task, error) {
	rows, err := r.db.Query(
		`SELECT ID, Name, Description, ProjectID, TaskCompleted, DueDate, CompletionDate, CreationDate, LastUpdatedDate, Priority FROM Tasks`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Description,
			&task.ProjectID,
			&task.TaskCompleted,
			&task.DueDate,
			&task.CompletionDate,
			&task.CreationDate,
			&task.LastUpdatedDate,
			&task.Priority,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *Repository) GetTasksByProjectID(projectID int) ([]*models.Task, error) {
	rows, err := r.db.Query(
		`SELECT ID, Name, Description, ProjectID, TaskCompleted, DueDate, CompletionDate, CreationDate, LastUpdatedDate, Priority FROM Tasks WHERE ProjectID = ?`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Description,
			&task.ProjectID,
			&task.TaskCompleted,
			&task.DueDate,
			&task.CompletionDate,
			&task.CreationDate,
			&task.LastUpdatedDate,
			&task.Priority,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *Repository) UpdateTask(task *models.Task) error {
	task.LastUpdatedDate = time.Now()
	_, err := r.db.Exec(
		`UPDATE Tasks SET Name = ?, Description = ?, ProjectID = ?, TaskCompleted = ?, DueDate = ?, CompletionDate = ?, LastUpdatedDate = ?, Priority = ? WHERE ID = ?`,
		task.Name,
		task.Description,
		task.ProjectID,
		task.TaskCompleted,
		task.DueDate,
		task.CompletionDate,
		task.LastUpdatedDate,
		task.Priority,
		task.ID,
	)
	return err
}

func (r *Repository) DeleteTask(id int) error {
	_, err := r.db.Exec(`DELETE FROM Tasks WHERE ID = ?`, id)
	return err
}
