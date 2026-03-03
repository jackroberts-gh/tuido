package model

import (
	"time"

	"github.com/google/uuid"
)

// TaskList represents a collection of tasks
type TaskList struct {
	Tasks []Task `json:"tasks"`
}

// NewTaskList creates a new empty task list
func NewTaskList() *TaskList {
	return &TaskList{
		Tasks: []Task{},
	}
}

// Add adds a new task to the list
func (tl *TaskList) Add(text string, priority Priority, dueDate *time.Time) {
	task := Task{
		ID:        uuid.New().String(),
		Text:      text,
		Completed: false,
		Priority:  priority,
		DueDate:   dueDate,
		CreatedAt: time.Now(),
	}
	tl.Tasks = append(tl.Tasks, task)
}

// Remove removes a task by ID
func (tl *TaskList) Remove(id string) bool {
	for i, task := range tl.Tasks {
		if task.ID == id {
			tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)
			return true
		}
	}
	return false
}

// Toggle toggles the completion status of a task by ID
func (tl *TaskList) Toggle(id string) bool {
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == id {
			tl.Tasks[i].Completed = !tl.Tasks[i].Completed
			if tl.Tasks[i].Completed {
				now := time.Now()
				tl.Tasks[i].CompletedAt = &now
				// Clear in-progress when completing
				tl.Tasks[i].InProgress = false
			} else {
				tl.Tasks[i].CompletedAt = nil
			}
			return true
		}
	}
	return false
}

// ToggleInProgress toggles the in-progress status of a task by ID
func (tl *TaskList) ToggleInProgress(id string) bool {
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == id {
			// Only toggle if not completed
			if !tl.Tasks[i].Completed {
				tl.Tasks[i].InProgress = !tl.Tasks[i].InProgress
				return true
			}
			return false
		}
	}
	return false
}

// CycleStatus cycles through task states: not started -> in-progress -> completed -> not started
func (tl *TaskList) CycleStatus(id string) bool {
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == id {
			if tl.Tasks[i].Completed {
				// Completed -> Not started
				tl.Tasks[i].Completed = false
				tl.Tasks[i].InProgress = false
				tl.Tasks[i].CompletedAt = nil
			} else if tl.Tasks[i].InProgress {
				// In-progress -> Completed
				tl.Tasks[i].InProgress = false
				tl.Tasks[i].Completed = true
				now := time.Now()
				tl.Tasks[i].CompletedAt = &now
			} else {
				// Not started -> In-progress
				tl.Tasks[i].InProgress = true
			}
			return true
		}
	}
	return false
}

// UpdatePriority updates the priority of a task by ID
func (tl *TaskList) UpdatePriority(id string, priority Priority) bool {
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == id {
			tl.Tasks[i].Priority = priority
			return true
		}
	}
	return false
}

// UpdateDueDate updates the due date of a task by ID
func (tl *TaskList) UpdateDueDate(id string, dueDate *time.Time) bool {
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == id {
			tl.Tasks[i].DueDate = dueDate
			return true
		}
	}
	return false
}

// GetByID retrieves a task by ID
func (tl *TaskList) GetByID(id string) *Task {
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == id {
			return &tl.Tasks[i]
		}
	}
	return nil
}

// GetByIndex retrieves a task by index (for cursor navigation)
func (tl *TaskList) GetByIndex(index int) *Task {
	if index < 0 || index >= len(tl.Tasks) {
		return nil
	}
	return &tl.Tasks[index]
}

// FilterActive returns only incomplete tasks
func (tl *TaskList) FilterActive() []Task {
	active := []Task{}
	for _, task := range tl.Tasks {
		if !task.Completed {
			active = append(active, task)
		}
	}
	return active
}

// FilterCompleted returns only completed tasks
func (tl *TaskList) FilterCompleted() []Task {
	completed := []Task{}
	for _, task := range tl.Tasks {
		if task.Completed {
			completed = append(completed, task)
		}
	}
	return completed
}

// Count returns the total number of tasks
func (tl *TaskList) Count() int {
	return len(tl.Tasks)
}

// CountActive returns the number of incomplete tasks
func (tl *TaskList) CountActive() int {
	count := 0
	for _, task := range tl.Tasks {
		if !task.Completed {
			count++
		}
	}
	return count
}

// CountCompleted returns the number of completed tasks
func (tl *TaskList) CountCompleted() int {
	count := 0
	for _, task := range tl.Tasks {
		if task.Completed {
			count++
		}
	}
	return count
}
