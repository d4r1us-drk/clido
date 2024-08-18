package models

import (
	"time"

	"github.com/d4r1us-drk/clido/pkg/utils"
	"gorm.io/gorm"
)

type Task struct {
	ID              int            `gorm:"primaryKey"              json:"id"`
	Name            string         `gorm:"not null"                json:"name"`
	Description     string         `                               json:"description"`
	ProjectID       int            `gorm:"not null"                json:"project_id"`
	Project         Project        `gorm:"foreignKey:ProjectID"    json:"-"`
	TaskCompleted   bool           `gorm:"not null"                json:"task_completed"`
	DueDate         *time.Time     `                               json:"due_date,omitempty"`
	CompletionDate  *time.Time     `                               json:"completion_date,omitempty"`
	CreationDate    time.Time      `gorm:"not null"                json:"creation_date"`
	LastUpdatedDate time.Time      `gorm:"not null"                json:"last_updated_date"`
	Priority        utils.Priority `gorm:"not null;default:4"      json:"priority"`
	ParentTaskID    *int           `                               json:"parent_task_id,omitempty"`
	ParentTask      *Task          `gorm:"foreignKey:ParentTaskID" json:"-"`
	SubTasks        []Task         `gorm:"foreignKey:ParentTaskID" json:"-"`
}

func (t *Task) BeforeCreate(_ *gorm.DB) error {
	t.CreationDate = time.Now()
	t.LastUpdatedDate = time.Now()
	return nil
}

func (t *Task) BeforeUpdate(_ *gorm.DB) error {
	t.LastUpdatedDate = time.Now()
	return nil
}
