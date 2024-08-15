package models

import (
	"time"

	"github.com/d4r1us-drk/clido/pkg/utils"
)

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
	Priority        utils.Priority
	ParentTaskID    *int
}
