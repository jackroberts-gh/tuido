package model

import (
	"testing"
	"time"
)

func TestCycleStatus(t *testing.T) {
	tl := NewTaskList()
	tl.Add("Test task", PriorityMedium, nil)

	if len(tl.Tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tl.Tasks))
	}

	taskID := tl.Tasks[0].ID

	// Test: Not started -> In-progress
	if !tl.CycleStatus(taskID) {
		t.Fatal("CycleStatus should return true for valid ID")
	}
	task := tl.GetByID(taskID)
	if task == nil {
		t.Fatal("Task should not be nil")
	}
	if !task.InProgress {
		t.Error("Task should be in progress")
	}
	if task.Completed {
		t.Error("Task should not be completed")
	}

	// Test: In-progress -> Completed
	if !tl.CycleStatus(taskID) {
		t.Fatal("CycleStatus should return true for valid ID")
	}
	task = tl.GetByID(taskID)
	if task == nil {
		t.Fatal("Task should not be nil")
	}
	if task.InProgress {
		t.Error("Task should not be in progress")
	}
	if !task.Completed {
		t.Error("Task should be completed")
	}
	if task.CompletedAt == nil {
		t.Error("CompletedAt should be set")
	}

	// Test: Completed -> Not started
	if !tl.CycleStatus(taskID) {
		t.Fatal("CycleStatus should return true for valid ID")
	}
	task = tl.GetByID(taskID)
	if task == nil {
		t.Fatal("Task should not be nil")
	}
	if task.InProgress {
		t.Error("Task should not be in progress")
	}
	if task.Completed {
		t.Error("Task should not be completed")
	}
	if task.CompletedAt != nil {
		t.Error("CompletedAt should be nil")
	}
}

func TestCycleStatusInvalidID(t *testing.T) {
	tl := NewTaskList()
	if tl.CycleStatus("invalid-id") {
		t.Error("CycleStatus should return false for invalid ID")
	}
}

func TestToggle(t *testing.T) {
	tl := NewTaskList()
	tl.Add("Test task", PriorityLow, nil)
	taskID := tl.Tasks[0].ID

	// Test: Not completed -> Completed
	if !tl.Toggle(taskID) {
		t.Fatal("Toggle should return true for valid ID")
	}
	task := tl.GetByID(taskID)
	if !task.Completed {
		t.Error("Task should be completed")
	}
	if task.CompletedAt == nil {
		t.Error("CompletedAt should be set")
	}
	if task.InProgress {
		t.Error("InProgress should be cleared when completing")
	}

	// Test: Completed -> Not completed
	if !tl.Toggle(taskID) {
		t.Fatal("Toggle should return true for valid ID")
	}
	task = tl.GetByID(taskID)
	if task.Completed {
		t.Error("Task should not be completed")
	}
	if task.CompletedAt != nil {
		t.Error("CompletedAt should be nil")
	}
}

func TestAddRemove(t *testing.T) {
	tl := NewTaskList()

	// Test Add
	tl.Add("Task 1", PriorityHigh, nil)
	if len(tl.Tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tl.Tasks))
	}
	if tl.Tasks[0].Text != "Task 1" {
		t.Errorf("Expected 'Task 1', got '%s'", tl.Tasks[0].Text)
	}
	if tl.Tasks[0].Priority != PriorityHigh {
		t.Errorf("Expected PriorityHigh, got %v", tl.Tasks[0].Priority)
	}

	taskID := tl.Tasks[0].ID

	// Test Remove
	if !tl.Remove(taskID) {
		t.Fatal("Remove should return true for valid ID")
	}
	if len(tl.Tasks) != 0 {
		t.Fatalf("Expected 0 tasks, got %d", len(tl.Tasks))
	}

	// Test Remove invalid ID
	if tl.Remove("invalid-id") {
		t.Error("Remove should return false for invalid ID")
	}
}

func TestFilterActive(t *testing.T) {
	tl := NewTaskList()
	tl.Add("Task 1", PriorityLow, nil)
	tl.Add("Task 2", PriorityMedium, nil)
	tl.Add("Task 3", PriorityHigh, nil)

	// Complete one task
	tl.Toggle(tl.Tasks[1].ID)

	active := tl.FilterActive()
	if len(active) != 2 {
		t.Errorf("Expected 2 active tasks, got %d", len(active))
	}
}

func TestFilterCompleted(t *testing.T) {
	tl := NewTaskList()
	tl.Add("Task 1", PriorityLow, nil)
	tl.Add("Task 2", PriorityMedium, nil)
	tl.Add("Task 3", PriorityHigh, nil)

	// Complete two tasks
	tl.Toggle(tl.Tasks[0].ID)
	tl.Toggle(tl.Tasks[2].ID)

	completed := tl.FilterCompleted()
	if len(completed) != 2 {
		t.Errorf("Expected 2 completed tasks, got %d", len(completed))
	}
}

func TestCount(t *testing.T) {
	tl := NewTaskList()

	if tl.Count() != 0 {
		t.Error("Empty list should have count 0")
	}

	tl.Add("Task 1", PriorityLow, nil)
	tl.Add("Task 2", PriorityMedium, nil)

	if tl.Count() != 2 {
		t.Errorf("Expected count 2, got %d", tl.Count())
	}
}

func TestCountActiveCompleted(t *testing.T) {
	tl := NewTaskList()
	tl.Add("Task 1", PriorityLow, nil)
	tl.Add("Task 2", PriorityMedium, nil)
	tl.Add("Task 3", PriorityHigh, nil)

	if tl.CountActive() != 3 {
		t.Errorf("Expected 3 active tasks, got %d", tl.CountActive())
	}
	if tl.CountCompleted() != 0 {
		t.Errorf("Expected 0 completed tasks, got %d", tl.CountCompleted())
	}

	// Complete one task
	tl.Toggle(tl.Tasks[1].ID)

	if tl.CountActive() != 2 {
		t.Errorf("Expected 2 active tasks, got %d", tl.CountActive())
	}
	if tl.CountCompleted() != 1 {
		t.Errorf("Expected 1 completed task, got %d", tl.CountCompleted())
	}
}

func TestUpdatePriority(t *testing.T) {
	tl := NewTaskList()
	tl.Add("Task 1", PriorityLow, nil)
	taskID := tl.Tasks[0].ID

	if !tl.UpdatePriority(taskID, PriorityHigh) {
		t.Fatal("UpdatePriority should return true for valid ID")
	}

	task := tl.GetByID(taskID)
	if task.Priority != PriorityHigh {
		t.Errorf("Expected PriorityHigh, got %v", task.Priority)
	}

	if tl.UpdatePriority("invalid-id", PriorityMedium) {
		t.Error("UpdatePriority should return false for invalid ID")
	}
}

func TestUpdateDueDate(t *testing.T) {
	tl := NewTaskList()
	tl.Add("Task 1", PriorityLow, nil)
	taskID := tl.Tasks[0].ID

	now := time.Now()
	if !tl.UpdateDueDate(taskID, &now) {
		t.Fatal("UpdateDueDate should return true for valid ID")
	}

	task := tl.GetByID(taskID)
	if task.DueDate == nil {
		t.Fatal("DueDate should not be nil")
	}
	if !task.DueDate.Equal(now) {
		t.Error("DueDate should match the set time")
	}

	if tl.UpdateDueDate("invalid-id", &now) {
		t.Error("UpdateDueDate should return false for invalid ID")
	}
}
