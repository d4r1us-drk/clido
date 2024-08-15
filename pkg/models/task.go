package models

import "time"

type Task struct {
	ID              int
	Name            string
	Description     string
	ProjectID       int
	TaskCompleted   bool
	DueDate         *time.Time
	CompletionDate  *time.Time
	CreationDate    time.Time
	LastUpdatedDate time.Time
	Priority        int
	ParentTaskID    *int
}
