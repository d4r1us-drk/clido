package models

import (
	"time"

	"github.com/d4r1us-drk/clido/pkg/utils"
)

type Task struct {
	ID              int            `json:"id"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	ProjectID       int            `json:"project_id"`
	TaskCompleted   bool           `json:"task_completed"`
	DueDate         *time.Time     `json:"due_date,omitempty"`
	CompletionDate  *time.Time     `json:"completion_date,omitempty"`
	CreationDate    time.Time      `json:"creation_date"`
	LastUpdatedDate time.Time      `json:"last_updated_date"`
	Priority        utils.Priority `json:"priority"`
	ParentTaskID    *int           `json:"parent_task_id,omitempty"`
}
