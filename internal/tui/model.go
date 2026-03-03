package tui

import (
	"sort"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackroberts-gh/tuido/internal/model"
	"github.com/jackroberts-gh/tuido/internal/storage"
)

// viewMode represents the current view/mode of the application
type viewMode int

const (
	modeList viewMode = iota // Main task list view
	modeAdd                  // Adding new task
	modeHelp                 // Help screen
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
	cursor        int           // Currently selected task index
	mode          viewMode      // Current view mode
	input         string        // Text input buffer
	err           error         // Error message to display
	width         int           // Terminal width
	height        int           // Terminal height
	showCompleted bool          // Whether to show completed tasks
	message       string        // Success/info message to display
	sortBy        sortMode      // Current sort mode
	lastKey       string        // Last key pressed (for key sequences like "sd", "sp")
	spinner       spinner.Model // Spinner for in-progress tasks
	// Add task form fields
	addField        int            // Current field in add mode (0=task, 1=priority, 2=due)
	addCursor       int            // Cursor position within priority/due lists
	addPriority     model.Priority // Selected priority for new task
	addDueSelection int            // Selected due date option (0=today, 1=tomorrow, 2=this week, 3=next week)
}

// NewModel creates a new Model with the given task list and storage
func NewModel(taskList *model.TaskList, storage *storage.Storage) Model {
	s := spinner.New()
	s.Spinner = spinner.Line

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
		spinner:       s,
	}
}

// Init initializes the model (required by BubbleTea)
func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
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
// Completed tasks always appear first in completion order (unsorted)
func (m Model) sortTasks(tasks []model.Task) []model.Task {
	// Separate completed and uncompleted tasks
	var completed, uncompleted []model.Task
	for _, task := range tasks {
		if task.Completed {
			completed = append(completed, task)
		} else {
			uncompleted = append(uncompleted, task)
		}
	}

	// Only sort uncompleted tasks - completed tasks stay in completion order
	if m.sortBy != sortNone {
		uncompleted = m.applySortCriteria(uncompleted)
	}

	// Return completed tasks first (in original order), then sorted uncompleted
	result := make([]model.Task, 0, len(tasks))
	result = append(result, completed...)
	result = append(result, uncompleted...)
	return result
}

// applySortCriteria applies the current sort criteria to a list of tasks
func (m Model) applySortCriteria(tasks []model.Task) []model.Task {
	if len(tasks) == 0 {
		return tasks
	}

	// Make a copy to avoid modifying the original
	sorted := make([]model.Task, len(tasks))
	copy(sorted, tasks)

	switch m.sortBy {
	case sortPriority:
		// Sort by priority: high > medium > low
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Priority > sorted[j].Priority
		})
	case sortPriorityReverse:
		// Sort by priority: low > medium > high
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Priority < sorted[j].Priority
		})
	case sortDueDate:
		// Sort by due date: earliest first, nil dates at the end (furthest in future)
		sort.Slice(sorted, func(i, j int) bool {
			// nil dates go last
			if sorted[i].DueDate == nil {
				return false
			}
			if sorted[j].DueDate == nil {
				return true
			}
			return sorted[i].DueDate.Before(*sorted[j].DueDate)
		})
	case sortDueDateReverse:
		// Sort by due date: latest first, nil dates at the beginning (furthest in future)
		sort.Slice(sorted, func(i, j int) bool {
			// nil dates go first
			if sorted[i].DueDate == nil {
				return true
			}
			if sorted[j].DueDate == nil {
				return false
			}
			return sorted[i].DueDate.After(*sorted[j].DueDate)
		})
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
