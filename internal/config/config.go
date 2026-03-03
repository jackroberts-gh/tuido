package config

import (
	"os"
	"path/filepath"
)

// GetDataDir returns the path to the .tuido data directory in the user's home directory
func GetDataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".tuido"), nil
}

// GetTasksFilePath returns the full path to the tasks.json file
func GetTasksFilePath() (string, error) {
	dataDir, err := GetDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "tasks.json"), nil
}
