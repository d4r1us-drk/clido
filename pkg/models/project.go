package models

import "time"

type Project struct {
	ID               int
	Name             string
	Description      string
	CreationDate     time.Time
	LastModifiedDate time.Time
	ParentProjectId  *int
}
