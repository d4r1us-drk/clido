package repository

import (
	"time"

	"github.com/d4r1us-drk/clido/pkg/models"
)

func (r *Repository) CreateTask(task *models.Task) error {
	task.CreationDate = time.Now()
	task.LastUpdatedDate = time.Now()

	// Find the lowest unused ID
	var id int
	err := r.db.QueryRow(`
		SELECT COALESCE(MIN(t1.ID + 1), 1)
		FROM Tasks t1
		LEFT JOIN Tasks t2 ON t1.ID + 1 = t2.ID
		WHERE t2.ID IS NULL`).Scan(&id)
	if err != nil {
		return err
	}

	// Insert the task with the found ID
	_, err = r.db.Exec(
		`INSERT INTO Tasks (ID, Name, Description, ProjectID, TaskCompleted, DueDate, CompletionDate, CreationDate,
    LastUpdatedDate, Priority, ParentTaskID) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id,
		task.Name,
		task.Description,
		task.ProjectID,
		task.TaskCompleted,
		task.DueDate,
		task.CompletionDate,
		task.CreationDate,
		task.LastUpdatedDate,
		task.Priority,
		task.ParentTaskID,
	)
	if err != nil {
		return err
	}

	task.ID = id
	return nil
}

func (r *Repository) GetTaskByID(id int) (*models.Task, error) {
	task := &models.Task{}
	err := r.db.QueryRow(
		`SELECT ID, Name, Description, ProjectID, TaskCompleted, DueDate, CompletionDate, CreationDate, LastUpdatedDate,
    Priority, ParentTaskID FROM Tasks WHERE ID = ?`,
		id,
	).
		Scan(
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
			&task.ParentTaskID,
		)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *Repository) GetAllTasks() ([]*models.Task, error) {
	rows, err := r.db.Query(
		`SELECT ID, Name, Description, ProjectID, TaskCompleted, DueDate, CompletionDate, CreationDate, LastUpdatedDate,
    Priority, ParentTaskID FROM Tasks`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		err = rows.Scan(
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
			&task.ParentTaskID,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) GetTasksByProjectID(projectID int) ([]*models.Task, error) {
	rows, err := r.db.Query(
		`SELECT ID, Name, Description, ProjectID, TaskCompleted, DueDate, CompletionDate, CreationDate, LastUpdatedDate, 
    Priority, ParentTaskID FROM Tasks WHERE ProjectID = ?`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		err = rows.Scan(
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
			&task.ParentTaskID,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) GetSubtasks(parentTaskID int) ([]*models.Task, error) {
	rows, err := r.db.Query(
		`SELECT ID, Name, Description, ProjectID, TaskCompleted, DueDate, CompletionDate, CreationDate, LastUpdatedDate, 
    Priority, ParentTaskID FROM Tasks WHERE ParentTaskID = ?`,
		parentTaskID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		err = rows.Scan(
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
			&task.ParentTaskID,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) UpdateTask(task *models.Task) error {
	task.LastUpdatedDate = time.Now()
	_, err := r.db.Exec(
		`UPDATE Tasks SET Name = ?, Description = ?, ProjectID = ?, TaskCompleted = ?, DueDate = ?, CompletionDate = ?, 
    LastUpdatedDate = ?, Priority = ?, ParentTaskID = ? WHERE ID = ?`,
		task.Name,
		task.Description,
		task.ProjectID,
		task.TaskCompleted,
		task.DueDate,
		task.CompletionDate,
		task.LastUpdatedDate,
		task.Priority,
		task.ParentTaskID,
		task.ID,
	)
	return err
}

func (r *Repository) DeleteTask(id int) error {
	_, err := r.db.Exec(`DELETE FROM Tasks WHERE ID = ?`, id)
	return err
}
