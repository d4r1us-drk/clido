package models

import (
	"time"

	"github.com/d4r1us-drk/clido/pkg/utils"
	"gorm.io/gorm"
)

// Task represents a task entity in the to-do list system.
// Each task belongs to a project and can have a parent task (for subtasks) and multiple subtasks.
// The task tracks its completion status, priority, and relevant timestamps.
//
// Fields:
//   - ID: The unique identifier for the task.
//   - Name: The name of the task, which is required.
//   - Description: A description of the task (optional).
//   - ProjectID: The ID of the project to which the task belongs (required).
//   - Project: A reference to the project this task belongs to (not serialized to JSON).
//   - TaskCompleted: A boolean indicating whether the task is completed (required).
//   - DueDate: The due date for the task (optional).
//   - CompletionDate: The date when the task was completed (optional, set when TaskCompleted is true).
//   - CreationDate: The date and time when the task was created (automatically set).
//   - LastUpdatedDate: The date and time when the task was last updated (automatically set).
//   - Priority: The priority level of the task, represented by an integer (1: High, 2: Medium, 3: Low, 4: None).
//   - ParentTaskID: The ID of the parent task, if this task is a subtask (optional).
//   - ParentTask: A reference to the parent task (not serialized to JSON).
//   - SubTasks: A list of subtasks belonging to this task (not serialized to JSON).
type Task struct {
	ID              int            `gorm:"primaryKey"              json:"id"`
	Name            string         `gorm:"not null"                json:"name"`
	Description     string         `                               json:"description"`
	ProjectID       int            `gorm:"not null"                json:"project_id"`
	Project         Project        `gorm:"foreignKey:ProjectID"    json:"-"`  // Foreign key for the project (ignored in JSON)
	TaskCompleted   bool           `gorm:"not null"                json:"task_completed"` // Tracks completion status
	DueDate         *time.Time     `                               json:"due_date,omitempty"` // Optional due date
	CompletionDate  *time.Time     `                               json:"completion_date,omitempty"` // Optional completion date
	CreationDate    time.Time      `gorm:"not null"                json:"creation_date"`
	LastUpdatedDate time.Time      `gorm:"not null"                json:"last_updated_date"`
	Priority        utils.Priority `gorm:"not null;default:4"      json:"priority"`  // Default priority set to "None"
	ParentTaskID    *int           `                               json:"parent_task_id,omitempty"` // Optional parent task
	ParentTask      *Task          `gorm:"foreignKey:ParentTaskID" json:"-"`  // Foreign key for the parent task (ignored in JSON)
	SubTasks        []Task         `gorm:"foreignKey:ParentTaskID" json:"-"`  // List of subtasks (ignored in JSON)
}

// BeforeCreate is a GORM hook that sets the CreationDate and LastUpdatedDate fields
// to the current time before a new task is inserted into the database.
func (t *Task) BeforeCreate(_ *gorm.DB) error {
	t.CreationDate = time.Now()
	t.LastUpdatedDate = time.Now()
	return nil
}

// BeforeUpdate is a GORM hook that updates the LastUpdatedDate field to the current time
// before an existing task is updated in the database.
func (t *Task) BeforeUpdate(_ *gorm.DB) error {
	t.LastUpdatedDate = time.Now()
	return nil
}
