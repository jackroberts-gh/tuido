package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackroberts-gh/tuido/internal/model"
	"github.com/jackroberts-gh/tuido/internal/storage"
)

// viewMode represents the current view/mode of the application
type viewMode int

const (
	modeList viewMode = iota // Main task list view
	modeAdd                   // Adding new task
	modeHelp                  // Help screen
)

// sortMode represents how tasks are sorted
type sortMode int

const (
	sortNone            sortMode = iota // No sorting (default order)
	sortPriority                        // Sort by priority (high to low)
	sortPriorityReverse                 // Sort by priority (low to high)
	sortDueDate                         // Sort by due date (earliest first)
	sortDueDateReverse                  // Sort by due date (latest first)
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
	message       string    // Success/info message to display
	sortBy        sortMode  // Current sort mode
	lastKey       string    // Last key pressed (for key sequences like "sd", "sp")
	// Add task form fields
	addField        int            // Current field in add mode (0=task, 1=priority, 2=due)
	addCursor       int            // Cursor position within priority/due lists
	addPriority     model.Priority // Selected priority for new task
	addDueSelection int            // Selected due date option (0=today, 1=tomorrow, 2=this week, 3=next week)
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
		message:       "",
	}
}

// Init initializes the model (required by BubbleTea)
func (m Model) Init() tea.Cmd {
	return nil
}

// getVisibleTasks returns the tasks that should be displayed based on showCompleted and sorting
func (m Model) getVisibleTasks() []model.Task {
	var tasks []model.Task
	if m.showCompleted {
		tasks = m.taskList.Tasks
	} else {
		tasks = m.taskList.FilterActive()
	}

	// Apply sorting
	return m.sortTasks(tasks)
}

// sortTasks sorts tasks based on the current sort mode
func (m Model) sortTasks(tasks []model.Task) []model.Task {
	if m.sortBy == sortNone {
		return tasks
	}

	// Make a copy to avoid modifying the original
	sorted := make([]model.Task, len(tasks))
	copy(sorted, tasks)

	switch m.sortBy {
	case sortPriority:
		// Sort by priority: high > medium > low
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i].Priority < sorted[j].Priority {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
	case sortPriorityReverse:
		// Sort by priority: low > medium > high
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i].Priority > sorted[j].Priority {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
	case sortDueDate:
		// Sort by due date: earliest first, nil dates at the end
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				// Handle nil due dates
				if sorted[i].DueDate == nil && sorted[j].DueDate != nil {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				} else if sorted[i].DueDate != nil && sorted[j].DueDate != nil {
					if sorted[i].DueDate.After(*sorted[j].DueDate) {
						sorted[i], sorted[j] = sorted[j], sorted[i]
					}
				}
			}
		}
	case sortDueDateReverse:
		// Sort by due date: latest first, nil dates at the end
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				// Handle nil due dates
				if sorted[i].DueDate == nil && sorted[j].DueDate != nil {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				} else if sorted[i].DueDate != nil && sorted[j].DueDate != nil {
					if sorted[i].DueDate.Before(*sorted[j].DueDate) {
						sorted[i], sorted[j] = sorted[j], sorted[i]
					}
				}
			}
		}
	}

	return sorted
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
