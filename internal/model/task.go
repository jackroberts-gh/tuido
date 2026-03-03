package model

import (
	"time"
)

// Priority represents the priority level of a task
type Priority int

const (
	PriorityLow    Priority = 0
	PriorityMedium Priority = 1
	PriorityHigh   Priority = 2
)

// String returns a string representation of the priority
func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "Low"
	case PriorityMedium:
		return "Medium"
	case PriorityHigh:
		return "High"
	default:
		return "Unknown"
	}
}

// Task represents a single TODO item
type Task struct {
	ID          string     `json:"id"`
	Text        string     `json:"text"`
	Completed   bool       `json:"completed"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// IsOverdue returns true if the task has a due date that has passed
func (t *Task) IsOverdue() bool {
	if t.DueDate == nil || t.Completed {
		return false
	}
	return time.Now().After(*t.DueDate)
}

// IsDueSoon returns true if the task is due within the next 2 days
func (t *Task) IsDueSoon() bool {
	if t.DueDate == nil || t.Completed {
		return false
	}
	now := time.Now()
	twoDaysFromNow := now.Add(48 * time.Hour)
	return t.DueDate.After(now) && t.DueDate.Before(twoDaysFromNow)
}

// DaysUntilDue returns the number of days until the task is due
// Returns 0 if no due date is set
func (t *Task) DaysUntilDue() int {
	if t.DueDate == nil {
		return 0
	}
	duration := time.Until(*t.DueDate)
	days := int(duration.Hours() / 24)
	return days
}
