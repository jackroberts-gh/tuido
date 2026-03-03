package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jackroberts-gh/tuido/internal/model"
)

func TestSaveAndLoad(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_tasks.json")

	// Create storage with test file path
	storage := &Storage{filePath: testFile}

	// Create a task list with some tasks
	tl := model.NewTaskList()
	tl.Add("Task 1", model.PriorityHigh, nil)
	tl.Add("Task 2", model.PriorityLow, nil)

	// Save the task list
	err := storage.Save(tl)
	if err != nil {
		t.Fatalf("Failed to save tasks: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Task file was not created")
	}

	// Load the task list
	loadedTL, err := storage.Load()
	if err != nil {
		t.Fatalf("Failed to load tasks: %v", err)
	}

	// Verify the loaded data matches
	if len(loadedTL.Tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(loadedTL.Tasks))
	}

	if loadedTL.Tasks[0].Text != "Task 1" {
		t.Errorf("Expected 'Task 1', got '%s'", loadedTL.Tasks[0].Text)
	}

	if loadedTL.Tasks[0].Priority != model.PriorityHigh {
		t.Errorf("Expected PriorityHigh, got %v", loadedTL.Tasks[0].Priority)
	}

	if loadedTL.Tasks[1].Text != "Task 2" {
		t.Errorf("Expected 'Task 2', got '%s'", loadedTL.Tasks[1].Text)
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	// Create storage with non-existent file path
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "nonexistent.json")

	storage := &Storage{filePath: testFile}

	// Load should return empty task list without error
	tl, err := storage.Load()
	if err != nil {
		t.Fatalf("Load should not error on non-existent file: %v", err)
	}

	if tl == nil {
		t.Fatal("Task list should not be nil")
	}

	if len(tl.Tasks) != 0 {
		t.Errorf("Expected empty task list, got %d tasks", len(tl.Tasks))
	}
}

func TestSaveAtomicity(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_tasks.json")

	storage := &Storage{filePath: testFile}

	// Create and save initial task list
	tl := model.NewTaskList()
	tl.Add("Task 1", model.PriorityMedium, nil)

	err := storage.Save(tl)
	if err != nil {
		t.Fatalf("Failed to save initial tasks: %v", err)
	}

	// Verify temp file doesn't exist after successful save
	tempFile := testFile + ".tmp"
	if _, err := os.Stat(tempFile); err == nil {
		t.Error("Temporary file should be cleaned up after save")
	}

	// Verify main file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Main task file should exist")
	}
}

func TestLoadEmptyJSON(t *testing.T) {
	// Create a temporary directory with an empty JSON file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "empty.json")

	// Write valid but empty JSON
	err := os.WriteFile(testFile, []byte(`{"tasks":[]}`), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	storage := &Storage{filePath: testFile}

	// Load should succeed and return empty task list
	tl, err := storage.Load()
	if err != nil {
		t.Fatalf("Load should not error on empty JSON: %v", err)
	}

	if len(tl.Tasks) != 0 {
		t.Errorf("Expected empty task list, got %d tasks", len(tl.Tasks))
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	// Create a temporary directory with invalid JSON
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "invalid.json")

	// Write invalid JSON
	err := os.WriteFile(testFile, []byte(`{invalid json`), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	storage := &Storage{filePath: testFile}

	// Load should return error
	_, err = storage.Load()
	if err == nil {
		t.Error("Load should error on invalid JSON")
	}
}

func TestSavePreservesTaskData(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_tasks.json")

	storage := &Storage{filePath: testFile}

	// Create a task with all fields populated
	tl := model.NewTaskList()
	tl.Add("Complete task", model.PriorityHigh, nil)
	taskID := tl.Tasks[0].ID

	// Set task to in-progress and then complete
	tl.ToggleInProgress(taskID)
	tl.Toggle(taskID)

	// Save
	err := storage.Save(tl)
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Load
	loadedTL, err := storage.Load()
	if err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	// Verify all data is preserved
	if len(loadedTL.Tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(loadedTL.Tasks))
	}

	task := loadedTL.Tasks[0]
	if task.ID != taskID {
		t.Error("Task ID should be preserved")
	}
	if task.Text != "Complete task" {
		t.Error("Task text should be preserved")
	}
	if !task.Completed {
		t.Error("Completed status should be preserved")
	}
	if task.CompletedAt == nil {
		t.Error("CompletedAt should be preserved")
	}
}
