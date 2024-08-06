package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(dbPath string) (*Repository, error) {
	homePath := os.Getenv("HOME")
	if homePath == "" {
		return nil, fmt.Errorf("The HOME environment variable is not set.")
	}

	dbDir := fmt.Sprintf("%s/.local/share/clido", homePath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("Error creating database directory: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	repo := &Repository{db: db}
	err = repo.init()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *Repository) init() error {
	createProjectTable := `
	CREATE TABLE IF NOT EXISTS Projects (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Name TEXT NOT NULL,
		Description TEXT,
		CreationDate DATETIME NOT NULL,
		LastModifiedDate DATETIME NOT NULL
	);`

	createTaskTable := `
	CREATE TABLE IF NOT EXISTS Tasks (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Name TEXT NOT NULL,
		Description TEXT,
		ProjectID INTEGER NOT NULL,
		TaskCompleted BOOLEAN NOT NULL,
		DueDate DATETIME,
		CompletionDate DATETIME,
		CreationDate DATETIME NOT NULL,
		LastUpdatedDate DATETIME NOT NULL,
		FOREIGN KEY (ProjectID) REFERENCES Projects(ID)
	);`

	_, err := r.db.Exec(createProjectTable)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(createTaskTable)
	return err
}

func (r *Repository) Close() {
	r.db.Close()
}

// Project CRUD operations
func (r *Repository) CreateProject(project *Project) error {
	project.CreationDate = time.Now()
	project.LastModifiedDate = time.Now()
	result, err := r.db.Exec(`INSERT INTO Projects (Name, Description, CreationDate, LastModifiedDate) VALUES (?, ?, ?, ?)`,
		project.Name, project.Description, project.CreationDate, project.LastModifiedDate)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	project.ID = int(id)
	return nil
}

func (r *Repository) GetProjectByID(id int) (*Project, error) {
	project := &Project{}
	err := r.db.QueryRow(`SELECT ID, Name, Description, CreationDate, LastModifiedDate FROM Projects WHERE ID = ?`, id).
		Scan(&project.ID, &project.Name, &project.Description, &project.CreationDate, &project.LastModifiedDate)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *Repository) GetProjectByName(name string) (*Project, error) {
	project := &Project{}
	err := r.db.QueryRow(`SELECT ID, Name, Description, CreationDate, LastModifiedDate FROM Projects WHERE Name = ?`, name).
		Scan(&project.ID, &project.Name, &project.Description, &project.CreationDate, &project.LastModifiedDate)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *Repository) GetAllProjects() ([]*Project, error) {
	rows, err := r.db.Query(`SELECT ID, Name, Description, CreationDate, LastModifiedDate FROM Projects`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*Project
	for rows.Next() {
		project := &Project{}
		err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.CreationDate, &project.LastModifiedDate)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *Repository) UpdateProject(project *Project) error {
	project.LastModifiedDate = time.Now()
	_, err := r.db.Exec(`UPDATE Projects SET Name = ?, Description = ?, LastModifiedDate = ? WHERE ID = ?`,
		project.Name, project.Description, project.LastModifiedDate, project.ID)
	return err
}

func (r *Repository) DeleteProject(id int) error {
	_, err := r.db.Exec(`DELETE FROM Projects WHERE ID = ?`, id)
	return err
}

// Task CRUD operations
func (r *Repository) CreateTask(task *Task) error {
	task.CreationDate = time.Now()
	task.LastUpdatedDate = time.Now()
	result, err := r.db.Exec(`INSERT INTO Tasks (Name, Description, ProjectID, TaskCompleted, DueDate, CompletionDate, CreationDate, LastUpdatedDate) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		task.Name, task.Description, task.ProjectID, task.TaskCompleted, task.DueDate, task.CompletionDate, task.CreationDate, task.LastUpdatedDate)
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

func (r *Repository) GetTaskByID(id int) (*Task, error) {
	task := &Task{}
	err := r.db.QueryRow(`SELECT ID, Name, Description, ProjectID, TaskCompleted, DueDate, CompletionDate, CreationDate, LastUpdatedDate FROM Tasks WHERE ID = ?`, id).
		Scan(&task.ID, &task.Name, &task.Description, &task.ProjectID, &task.TaskCompleted, &task.DueDate, &task.CompletionDate, &task.CreationDate, &task.LastUpdatedDate)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *Repository) GetAllTasks() ([]*Task, error) {
	rows, err := r.db.Query(`SELECT ID, Name, Description, ProjectID, TaskCompleted, DueDate, CompletionDate, CreationDate, LastUpdatedDate FROM Tasks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.ProjectID, &task.TaskCompleted, &task.DueDate, &task.CompletionDate, &task.CreationDate, &task.LastUpdatedDate)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *Repository) GetTasksByProjectID(projectID int) ([]*Task, error) {
	rows, err := r.db.Query(`SELECT ID, Name, Description, ProjectID, TaskCompleted, DueDate, CompletionDate, CreationDate, LastUpdatedDate FROM Tasks WHERE ProjectID = ?`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.ProjectID, &task.TaskCompleted, &task.DueDate, &task.CompletionDate, &task.CreationDate, &task.LastUpdatedDate)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *Repository) UpdateTask(task *Task) error {
	task.LastUpdatedDate = time.Now()
	_, err := r.db.Exec(`UPDATE Tasks SET Name = ?, Description = ?, ProjectID = ?, TaskCompleted = ?, DueDate = ?, CompletionDate = ?, LastUpdatedDate = ? WHERE ID = ?`,
		task.Name, task.Description, task.ProjectID, task.TaskCompleted, task.DueDate, task.CompletionDate, task.LastUpdatedDate, task.ID)
	return err
}

func (r *Repository) DeleteTask(id int) error {
	_, err := r.db.Exec(`DELETE FROM Tasks WHERE ID = ?`, id)
	return err
}
