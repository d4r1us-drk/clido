package controllers

import (
	"errors"

	"github.com/d4r1us-drk/clido/models"
	"github.com/d4r1us-drk/clido/repository"
	"github.com/d4r1us-drk/clido/utils"
)

// Error constants for project operations.
var (
	ErrNoProjectName         = errors.New("project name is required")
	ErrParentProjectNotFound = errors.New("parent project not found")
)

// ProjectController manages the project-related business logic.
type ProjectController struct {
	repo *repository.Repository
}

// NewProjectController creates and returns a new instance of ProjectController.
func NewProjectController(repo *repository.Repository) *ProjectController {
	return &ProjectController{repo: repo}
}

// CreateProject handles the creation of a new project.
func (pc *ProjectController) CreateProject(
	name, description, parentProjectIdentifier string,
) error {
	// Validate project name
	if name == "" {
		return ErrNoProjectName
	}

	// Retrieve the parent project ID (if any)
	var parentProjectID *int
	if parentProjectIdentifier != "" {
		parentID, projectErr := utils.ParseIntOrError(parentProjectIdentifier)
		if projectErr != nil {
			project, lookupErr := pc.repo.GetProjectByName(parentProjectIdentifier)
			if lookupErr != nil {
				return ErrParentProjectNotFound
			}
			parentID = project.ID
		}
		parentProjectID = &parentID
	}

	// Create a new project
	project := models.Project{
		Name:            name,
		Description:     description,
		ParentProjectID: parentProjectID,
	}

	// Store the project in the repository
	return pc.repo.CreateProject(&project)
}

// EditProject handles updating an existing project by its ID.
func (pc *ProjectController) EditProject(
	id int,
	name, description, parentProjectIdentifier string,
) error {
	// Retrieve the existing project
	project, getProjectErr := pc.repo.GetProjectByID(id)
	if getProjectErr != nil {
		return getProjectErr
	}

	// Apply updates
	if name != "" {
		project.Name = name
	}
	if description != "" {
		project.Description = description
	}
	if parentProjectIdentifier != "" {
		parentID, projectErr := utils.ParseIntOrError(parentProjectIdentifier)
		if projectErr != nil {
			projectByName, lookupErr := pc.repo.GetProjectByName(parentProjectIdentifier)
			if lookupErr != nil {
				return ErrParentProjectNotFound
			}
			parentID = projectByName.ID
		}
		project.ParentProjectID = &parentID
	}

	// Update the project in the repository
	return pc.repo.UpdateProject(project)
}

// ListProjects returns all projects stored in the repository.
func (pc *ProjectController) ListProjects() ([]*models.Project, error) {
	return pc.repo.GetAllProjects()
}

// GetProjectByID returns a project by its ID.
func (pc *ProjectController) GetProjectByID(id int) (*models.Project, error) {
	return pc.repo.GetProjectByID(id)
}

// GetProjectByName returns a project by its name.
func (pc *ProjectController) GetProjectByName(name string) (*models.Project, error) {
	return pc.repo.GetProjectByName(name)
}

// ListSubprojects returns subprojects for a specific parent project.
func (pc *ProjectController) ListSubprojects(parentID int) ([]*models.Project, error) {
	return pc.repo.GetSubprojects(parentID)
}

// RemoveProject handles the recursive removal of a project and all its subprojects.
func (pc *ProjectController) RemoveProject(id int) error {
	// Retrieve all subprojects of the project
	subprojects, getSubprojectsErr := pc.repo.GetSubprojects(id)
	if getSubprojectsErr != nil {
		return getSubprojectsErr
	}

	// Recursively remove subprojects
	for _, subproject := range subprojects {
		if removeErr := pc.RemoveProject(subproject.ID); removeErr != nil {
			return removeErr
		}
	}

	// Remove the parent project
	return pc.repo.DeleteProject(id)
}
