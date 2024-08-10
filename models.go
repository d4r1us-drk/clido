package main

import "time"

type Project struct {
	ID               int
	Name             string
	Description      string
	CreationDate     time.Time
	LastModifiedDate time.Time
}

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
}
