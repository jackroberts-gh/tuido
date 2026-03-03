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
	InProgress  bool       `json:"in_progress"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}
