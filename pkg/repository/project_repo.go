package repository

import (
	"time"

	"github.com/d4r1us-drk/clido/pkg/models"
)

func (r *Repository) CreateProject(project *models.Project) error {
	project.CreationDate = time.Now()
	project.LastModifiedDate = time.Now()
	result, err := r.db.Exec(
		`INSERT INTO Projects (Name, Description, CreationDate, LastModifiedDate) VALUES (?, ?, ?, ?)`,
		project.Name,
		project.Description,
		project.CreationDate,
		project.LastModifiedDate,
	)
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

func (r *Repository) GetProjectByID(id int) (*models.Project, error) {
	project := &models.Project{}
	err := r.db.QueryRow(`SELECT ID, Name, Description, CreationDate, LastModifiedDate FROM Projects WHERE ID = ?`, id).
		Scan(&project.ID, &project.Name, &project.Description, &project.CreationDate, &project.LastModifiedDate)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *Repository) GetProjectByName(name string) (*models.Project, error) {
	project := &models.Project{}
	err := r.db.QueryRow(`SELECT ID, Name, Description, CreationDate, LastModifiedDate FROM Projects WHERE Name = ?`, name).
		Scan(&project.ID, &project.Name, &project.Description, &project.CreationDate, &project.LastModifiedDate)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *Repository) GetAllProjects() ([]*models.Project, error) {
	rows, err := r.db.Query(
		`SELECT ID, Name, Description, CreationDate, LastModifiedDate FROM Projects`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.CreationDate,
			&project.LastModifiedDate,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *Repository) UpdateProject(project *models.Project) error {
	project.LastModifiedDate = time.Now()
	_, err := r.db.Exec(
		`UPDATE Projects SET Name = ?, Description = ?, LastModifiedDate = ? WHERE ID = ?`,
		project.Name,
		project.Description,
		project.LastModifiedDate,
		project.ID,
	)
	return err
}

func (r *Repository) DeleteProject(id int) error {
	_, err := r.db.Exec(`DELETE FROM Projects WHERE ID = ?`, id)
	return err
}
