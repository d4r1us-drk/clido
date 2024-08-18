package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID               int       `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"unique;not null" json:"name"`
	Description      string    `json:"description"`
	CreationDate     time.Time `gorm:"not null" json:"creation_date"`
	LastModifiedDate time.Time `gorm:"not null" json:"last_modified_date"`
	ParentProjectID  *int      `json:"parent_project_id,omitempty"`
	ParentProject    *Project  `gorm:"foreignKey:ParentProjectID" json:"-"`
	SubProjects      []Project `gorm:"foreignKey:ParentProjectID" json:"-"`
	Tasks            []Task    `gorm:"foreignKey:ProjectID" json:"-"`
}

func (p *Project) BeforeCreate(_ *gorm.DB) error {
	p.CreationDate = time.Now()
	p.LastModifiedDate = time.Now()
	return nil
}

func (p *Project) BeforeUpdate(_ *gorm.DB) error {
	p.LastModifiedDate = time.Now()
	return nil
}
