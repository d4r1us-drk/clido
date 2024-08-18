package models

import "time"

type Project struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	CreationDate     time.Time `json:"creation_date"`
	LastModifiedDate time.Time `json:"last_modified_date"`
	ParentProjectID  *int      `json:"parent_project_id,omitempty"`
}
