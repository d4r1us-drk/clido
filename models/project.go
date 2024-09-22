package models

import (
	"time"

	"gorm.io/gorm"
)

// Project represents a project entity in the to-do list system.
// It supports hierarchical relationships, where a project can have a parent project and multiple subprojects.
// Each project can also have associated tasks.
//
// Fields:
//   - ID: The unique identifier for the project.
//   - Name: The name of the project, which must be unique and non-null.
//   - Description: A description of the project (optional).
//   - CreationDate: The date and time when the project was created (automatically set).
//   - LastModifiedDate: The date and time when the project was last updated (automatically set).
//   - ParentProjectID: The ID of the parent project, if this project is a subproject (optional).
//   - ParentProject: A reference to the parent project (not serialized to JSON).
//   - SubProjects: A list of subprojects belonging to this project (not serialized to JSON).
//   - Tasks: A list of tasks associated with this project (not serialized to JSON).
type Project struct {
	ID               int       `gorm:"primaryKey"                 json:"id"`
	Name             string    `gorm:"unique;not null"            json:"name"`
	Description      string    `                                  json:"description"`
	CreationDate     time.Time `gorm:"not null"                   json:"creation_date"`
	LastModifiedDate time.Time `gorm:"not null"                   json:"last_modified_date"`
	ParentProjectID  *int      `                                  json:"parent_project_id,omitempty"`
	ParentProject    *Project  `gorm:"foreignKey:ParentProjectID" json:"-"`
	SubProjects      []Project `gorm:"foreignKey:ParentProjectID" json:"-"`
	Tasks            []Task    `gorm:"foreignKey:ProjectID"       json:"-"`
}

// BeforeCreate is a GORM hook that sets the CreationDate and LastModifiedDate fields
// to the current time before a new project is inserted into the database.
func (p *Project) BeforeCreate(_ *gorm.DB) error {
	p.CreationDate = time.Now()
	p.LastModifiedDate = time.Now()
	return nil
}

// BeforeUpdate is a GORM hook that updates the LastModifiedDate field to the current time
// before an existing project is updated in the database.
func (p *Project) BeforeUpdate(_ *gorm.DB) error {
	p.LastModifiedDate = time.Now()
	return nil
}
