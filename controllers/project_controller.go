package controllers

import (
	"errors"

	"github.com/d4r1us-drk/clido/models"
	"github.com/d4r1us-drk/clido/repository"
	"github.com/d4r1us-drk/clido/utils"
)

// Error constants for project operations.
var (
	ErrNoProjectName           = errors.New("project name is required")
	ErrParentProjectNotFound   = errors.New("parent project not found")
	ErrNoParentProjectProvided = errors.New("no parent project provided")
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
	parentProjectID, err := pc.getParentProjectID(parentProjectIdentifier)
	if err != nil && !errors.Is(err, ErrNoParentProjectProvided) {
		return err
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

	// Retrieve and apply the parent project ID (if any)
	parentProjectID, err := pc.getParentProjectID(parentProjectIdentifier)
	if err != nil && !errors.Is(err, ErrNoParentProjectProvided) {
		return err
	}
	project.ParentProjectID = parentProjectID

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

// getParentProjectID checks and retrieves the parent project ID based on the identifier (name or ID).
func (pc *ProjectController) getParentProjectID(parentProjectIdentifier string) (*int, error) {
	if parentProjectIdentifier == "" {
		// No parent project identifier provided, so no parent project ID is needed
		return nil, ErrNoParentProjectProvided
	}

	// Try to parse the parent project identifier as an integer ID
	parentID, err := utils.ParseIntOrError(parentProjectIdentifier)
	if err == nil {
		// Successfully parsed as ID, now check if the project exists by ID
		project, getProjectErr := pc.repo.GetProjectByID(parentID)
		if getProjectErr != nil || project == nil {
			return nil, ErrParentProjectNotFound
		}
		return &parentID, nil
	}

	// If parsing failed, treat it as a project name and search by name
	project, lookupErr := pc.repo.GetProjectByName(parentProjectIdentifier)
	if lookupErr != nil || project == nil {
		return nil, ErrParentProjectNotFound
	}

	// Return the found project ID
	return &project.ID, nil
}
