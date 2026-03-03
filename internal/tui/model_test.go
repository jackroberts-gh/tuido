package tui

import (
	"testing"
	"time"

	"github.com/jackroberts-gh/tuido/internal/model"
)

func TestApplySortCriteriaPriority(t *testing.T) {
	m := Model{sortBy: sortPriority}

	tasks := []model.Task{
		{ID: "1", Text: "Low", Priority: model.PriorityLow},
		{ID: "2", Text: "High", Priority: model.PriorityHigh},
		{ID: "3", Text: "Medium", Priority: model.PriorityMedium},
	}

	sorted := m.applySortCriteria(tasks)

	// Should be sorted: High, Medium, Low
	if sorted[0].Priority != model.PriorityHigh {
		t.Errorf("First task should be High priority, got %v", sorted[0].Priority)
	}
	if sorted[1].Priority != model.PriorityMedium {
		t.Errorf("Second task should be Medium priority, got %v", sorted[1].Priority)
	}
	if sorted[2].Priority != model.PriorityLow {
		t.Errorf("Third task should be Low priority, got %v", sorted[2].Priority)
	}
}

func TestApplySortCriteriaPriorityReverse(t *testing.T) {
	m := Model{sortBy: sortPriorityReverse}

	tasks := []model.Task{
		{ID: "1", Text: "Low", Priority: model.PriorityLow},
		{ID: "2", Text: "High", Priority: model.PriorityHigh},
		{ID: "3", Text: "Medium", Priority: model.PriorityMedium},
	}

	sorted := m.applySortCriteria(tasks)

	// Should be sorted: Low, Medium, High
	if sorted[0].Priority != model.PriorityLow {
		t.Errorf("First task should be Low priority, got %v", sorted[0].Priority)
	}
	if sorted[1].Priority != model.PriorityMedium {
		t.Errorf("Second task should be Medium priority, got %v", sorted[1].Priority)
	}
	if sorted[2].Priority != model.PriorityHigh {
		t.Errorf("Third task should be High priority, got %v", sorted[2].Priority)
	}
}

func TestApplySortCriteriaDueDate(t *testing.T) {
	m := Model{sortBy: sortDueDate}

	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	nextWeek := now.Add(7 * 24 * time.Hour)

	tasks := []model.Task{
		{ID: "1", Text: "Next week", DueDate: &nextWeek},
		{ID: "2", Text: "No due date", DueDate: nil},
		{ID: "3", Text: "Tomorrow", DueDate: &tomorrow},
	}

	sorted := m.applySortCriteria(tasks)

	// Should be sorted: Tomorrow, Next week, No due date (nil last)
	if sorted[0].ID != "3" {
		t.Errorf("First task should be 'Tomorrow', got %s", sorted[0].Text)
	}
	if sorted[1].ID != "1" {
		t.Errorf("Second task should be 'Next week', got %s", sorted[1].Text)
	}
	if sorted[2].ID != "2" {
		t.Errorf("Third task should be 'No due date', got %s", sorted[2].Text)
	}
}

func TestApplySortCriteriaDueDateReverse(t *testing.T) {
	m := Model{sortBy: sortDueDateReverse}

	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	nextWeek := now.Add(7 * 24 * time.Hour)

	tasks := []model.Task{
		{ID: "1", Text: "Next week", DueDate: &nextWeek},
		{ID: "2", Text: "No due date", DueDate: nil},
		{ID: "3", Text: "Tomorrow", DueDate: &tomorrow},
	}

	sorted := m.applySortCriteria(tasks)

	// Should be sorted: No due date (nil first), Next week, Tomorrow
	if sorted[0].ID != "2" {
		t.Errorf("First task should be 'No due date', got %s", sorted[0].Text)
	}
	if sorted[1].ID != "1" {
		t.Errorf("Second task should be 'Next week', got %s", sorted[1].Text)
	}
	if sorted[2].ID != "3" {
		t.Errorf("Third task should be 'Tomorrow', got %s", sorted[2].Text)
	}
}

func TestApplySortCriteriaEmpty(t *testing.T) {
	m := Model{sortBy: sortPriority}
	tasks := []model.Task{}

	sorted := m.applySortCriteria(tasks)

	if len(sorted) != 0 {
		t.Errorf("Expected empty slice, got %d tasks", len(sorted))
	}
}

func TestSortTasksCompletedFirst(t *testing.T) {
	m := Model{sortBy: sortPriority}

	tasks := []model.Task{
		{ID: "1", Text: "Active Low", Priority: model.PriorityLow, Completed: false},
		{ID: "2", Text: "Completed High", Priority: model.PriorityHigh, Completed: true},
		{ID: "3", Text: "Active High", Priority: model.PriorityHigh, Completed: false},
	}

	sorted := m.sortTasks(tasks)

	// Completed tasks should appear first (unsorted), then active tasks sorted by priority
	if !sorted[0].Completed {
		t.Error("First task should be completed")
	}
	if sorted[0].ID != "2" {
		t.Errorf("First task should be 'Completed High', got %s", sorted[0].Text)
	}

	// Remaining tasks should be sorted by priority (High, then Low)
	if sorted[1].Priority != model.PriorityHigh {
		t.Errorf("Second task should be High priority, got %v", sorted[1].Priority)
	}
	if sorted[2].Priority != model.PriorityLow {
		t.Errorf("Third task should be Low priority, got %v", sorted[2].Priority)
	}
}

func TestGetVisibleTasksShowCompleted(t *testing.T) {
	tl := model.NewTaskList()
	tl.Add("Task 1", model.PriorityLow, nil)
	tl.Add("Task 2", model.PriorityMedium, nil)
	tl.Toggle(tl.Tasks[0].ID) // Complete first task

	m := Model{
		taskList:      tl,
		showCompleted: true,
		sortBy:        sortNone,
	}

	visible := m.getVisibleTasks()
	if len(visible) != 2 {
		t.Errorf("Expected 2 visible tasks with showCompleted=true, got %d", len(visible))
	}
}

func TestGetVisibleTasksHideCompleted(t *testing.T) {
	tl := model.NewTaskList()
	tl.Add("Task 1", model.PriorityLow, nil)
	tl.Add("Task 2", model.PriorityMedium, nil)
	tl.Toggle(tl.Tasks[0].ID) // Complete first task

	m := Model{
		taskList:      tl,
		showCompleted: false,
		sortBy:        sortNone,
	}

	visible := m.getVisibleTasks()
	if len(visible) != 1 {
		t.Errorf("Expected 1 visible task with showCompleted=false, got %d", len(visible))
	}
	if visible[0].Completed {
		t.Error("Visible task should not be completed")
	}
}
