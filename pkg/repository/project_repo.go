package repository

import (
	"time"

	"github.com/d4r1us-drk/clido/pkg/models"
)

func (r *Repository) CreateProject(project *models.Project) error {
	project.CreationDate = time.Now()
	project.LastModifiedDate = time.Now()

	// Find the lowest unused ID
	var id int
	err := r.db.QueryRow(`
		SELECT COALESCE(MIN(p1.ID + 1), 1)
		FROM Projects p1
		LEFT JOIN Projects p2 ON p1.ID + 1 = p2.ID
		WHERE p2.ID IS NULL`).Scan(&id)
	if err != nil {
		return err
	}

	// Insert the project with the found ID
	_, err = r.db.Exec(
		`INSERT INTO Projects (ID, Name, Description, CreationDate, LastModifiedDate, ParentProjectID)
    VALUES (?, ?, ?, ?, ?, ?)`,
		id,
		project.Name,
		project.Description,
		project.CreationDate,
		project.LastModifiedDate,
		project.ParentProjectID,
	)
	if err != nil {
		return err
	}

	project.ID = id
	return nil
}

func (r *Repository) GetProjectByID(id int) (*models.Project, error) {
	project := &models.Project{}
	err := r.db.QueryRow(
		`SELECT ID, Name, Description, CreationDate, LastModifiedDate, ParentProjectID FROM Projects WHERE ID = ?`,
		id,
	).
		Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.CreationDate,
			&project.LastModifiedDate,
			&project.ParentProjectID,
		)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *Repository) GetProjectByName(name string) (*models.Project, error) {
	project := &models.Project{}
	err := r.db.QueryRow(
		`SELECT ID, Name, Description, CreationDate, LastModifiedDate, ParentProjectID FROM Projects WHERE Name = ?`,
		name,
	).
		Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.CreationDate,
			&project.LastModifiedDate,
			&project.ParentProjectID,
		)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *Repository) GetAllProjects() ([]*models.Project, error) {
	rows, err := r.db.Query(
		`SELECT ID, Name, Description, CreationDate, LastModifiedDate, ParentProjectID FROM Projects`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		err = rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.CreationDate,
			&project.LastModifiedDate,
			&project.ParentProjectID,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *Repository) GetSubprojects(parentProjectID int) ([]*models.Project, error) {
	rows, err := r.db.Query(
		`SELECT ID, Name, Description, CreationDate, LastModifiedDate, ParentProjectID FROM Projects WHERE
    ParentProjectID = ?`,
		parentProjectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		project := &models.Project{}
		err = rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.CreationDate,
			&project.LastModifiedDate,
			&project.ParentProjectID,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *Repository) UpdateProject(project *models.Project) error {
	project.LastModifiedDate = time.Now()
	_, err := r.db.Exec(
		`UPDATE Projects SET Name = ?, Description = ?, LastModifiedDate = ?, ParentProjectID = ? WHERE ID = ?`,
		project.Name,
		project.Description,
		project.LastModifiedDate,
		project.ParentProjectID,
		project.ID,
	)
	return err
}

func (r *Repository) DeleteProject(id int) error {
	_, err := r.db.Exec(`DELETE FROM Projects WHERE ID = ?`, id)
	return err
}
