package repository

import (
	"github.com/d4r1us-drk/clido/pkg/models"
)

func (r *Repository) CreateProject(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *Repository) GetProjectByID(id int) (*models.Project, error) {
	var project models.Project
	err := r.db.First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *Repository) GetProjectByName(name string) (*models.Project, error) {
	var project models.Project
	err := r.db.Where("name = ?", name).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *Repository) GetAllProjects() ([]*models.Project, error) {
	var projects []*models.Project
	err := r.db.Find(&projects).Error
	return projects, err
}

func (r *Repository) GetSubprojects(parentProjectID int) ([]*models.Project, error) {
	var projects []*models.Project
	err := r.db.Where("parent_project_id = ?", parentProjectID).Find(&projects).Error
	return projects, err
}

func (r *Repository) UpdateProject(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *Repository) DeleteProject(id int) error {
	return r.db.Delete(&models.Project{}, id).Error
}

func (r *Repository) GetNextProjectID() (int, error) {
	var maxID int
	err := r.db.Model(&models.Project{}).Select("COALESCE(MAX(id), 0) + 1").Scan(&maxID).Error
	return maxID, err
}
