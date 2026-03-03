package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackroberts-gh/tuido/internal/storage"
	"github.com/jackroberts-gh/tuido/internal/tui"
)

func main() {
	// Initialize storage
	store, err := storage.NewStorage()
	if err != nil {
		fmt.Printf("Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	// Load existing tasks
	taskList, err := store.Load()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	// Create BubbleTea model
	m := tui.NewModel(taskList, store)

	// Check for command-line arguments
	if len(os.Args) > 1 && os.Args[1] == "add" {
		m = m.StartInAddMode()
	}

	// Run TUI with alternate screen buffer
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
