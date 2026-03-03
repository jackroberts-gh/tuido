package storage

import (
	"encoding/json"
	"os"

	"github.com/jackroberts-gh/tuido/internal/config"
	"github.com/jackroberts-gh/tuido/internal/model"
)

// Storage handles reading and writing tasks to disk
type Storage struct {
	filePath string
}

// NewStorage creates a new Storage instance and ensures the data directory exists
func NewStorage() (*Storage, error) {
	dataDir, err := config.GetDataDir()
	if err != nil {
		return nil, err
	}

	// Create the .tuido directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	filePath, err := config.GetTasksFilePath()
	if err != nil {
		return nil, err
	}

	return &Storage{
		filePath: filePath,
	}, nil
}

// Load reads the tasks from the JSON file
// Returns an empty TaskList if the file doesn't exist
func (s *Storage) Load() (*model.TaskList, error) {
	// Check if file exists
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		// File doesn't exist, return empty task list
		return model.NewTaskList(), nil
	}

	// Read file
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var taskList model.TaskList
	if err := json.Unmarshal(data, &taskList); err != nil {
		return nil, err
	}

	// Initialize the Tasks slice if it's nil
	if taskList.Tasks == nil {
		taskList.Tasks = []model.Task{}
	}

	return &taskList, nil
}

// Save writes the tasks to the JSON file atomically
func (s *Storage) Save(taskList *model.TaskList) error {
	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(taskList, "", "  ")
	if err != nil {
		return err
	}

	// Write to a temporary file first
	tempFile := s.filePath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return err
	}

	// Atomically rename temp file to actual file
	if err := os.Rename(tempFile, s.filePath); err != nil {
		// Clean up temp file on error
		os.Remove(tempFile)
		return err
	}

	return nil
}
