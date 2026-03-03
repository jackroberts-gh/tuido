package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackroberts-gh/tuido/internal/model"
	"github.com/jackroberts-gh/tuido/internal/storage"
)

// viewMode represents the current view/mode of the application
type viewMode int

const (
	modeList         viewMode = iota // Main task list view
	modeAdd                           // Adding new task
	modeEditPriority                  // Editing task priority
	modeEditDueDate                   // Editing task due date
	modeHelp                          // Help screen
)

// Model represents the application state for BubbleTea
type Model struct {
	taskList      *model.TaskList
	storage       *storage.Storage
	cursor        int       // Currently selected task index
	mode          viewMode  // Current view mode
	input         string    // Text input buffer
	err           error     // Error message to display
	width         int       // Terminal width
	height        int       // Terminal height
	showCompleted bool      // Whether to show completed tasks
	selectedID    string    // ID of the selected task (for editing)
	message       string    // Success/info message to display
}

// NewModel creates a new Model with the given task list and storage
func NewModel(taskList *model.TaskList, storage *storage.Storage) Model {
	return Model{
		taskList:      taskList,
		storage:       storage,
		cursor:        0,
		mode:          modeList,
		input:         "",
		err:           nil,
		width:         80,
		height:        24,
		showCompleted: true,
		selectedID:    "",
		message:       "",
	}
}

// Init initializes the model (required by BubbleTea)
func (m Model) Init() tea.Cmd {
	return nil
}

// getVisibleTasks returns the tasks that should be displayed based on showCompleted
func (m Model) getVisibleTasks() []model.Task {
	if m.showCompleted {
		return m.taskList.Tasks
	}
	return m.taskList.FilterActive()
}

// getCurrentTask returns the currently selected task (based on cursor position)
func (m Model) getCurrentTask() *model.Task {
	visibleTasks := m.getVisibleTasks()
	if m.cursor < 0 || m.cursor >= len(visibleTasks) {
		return nil
	}
	// Find the actual task in the full list by ID
	taskID := visibleTasks[m.cursor].ID
	return m.taskList.GetByID(taskID)
}

// saveToStorage saves the task list to storage and handles errors
func (m *Model) saveToStorage() {
	if err := m.storage.Save(m.taskList); err != nil {
		m.err = err
	}
}

// clearMessages clears error and info messages
func (m *Model) clearMessages() {
	m.err = nil
	m.message = ""
}
