package tui

import (
	"sort"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/jackroberts-gh/tuido/v2/internal/model"
	"github.com/jackroberts-gh/tuido/v2/internal/storage"
)

const (
	saveDebounceDuration = 500 * time.Millisecond

	// themePollInterval is how often the terminal is re-queried for its
	// background color, so the palette follows light/dark mode switches.
	//
	// This is only used on terminals that do not support DEC private mode 2031.
	// Terminals that do support it push a color scheme event the moment the
	// theme changes (see themeNotificationsCmd), and the poll is switched off
	// entirely for them.
	themePollInterval = 500 * time.Millisecond

	// themeNotificationMode is the DEC private mode number for light/dark
	// change notifications (ansi.ModeLightDark).
	themeNotificationMode = 2031
)

// saveMsg is sent when the debounced save timer expires
type saveMsg struct{}

// themeTickMsg drives the background color polling loop
type themeTickMsg time.Time

// pollTheme re-queries the terminal background color and schedules the next poll
func pollTheme() tea.Cmd {
	return tea.Batch(tea.RequestBackgroundColor, themeTick())
}

// themeTick schedules the next theme poll
func themeTick() tea.Cmd {
	return tea.Tick(themePollInterval, func(t time.Time) tea.Msg {
		return themeTickMsg(t)
	})
}

// themeNotificationsCmd asks the terminal to report light/dark mode changes as
// they happen (DEC private mode 2031), then asks whether that mode is actually
// supported. Supporting terminals answer the query and push a color scheme
// event on every theme switch, which lets us stop polling entirely; the rest
// never answer and stay on the poll.
func themeNotificationsCmd() tea.Cmd {
	return tea.Batch(
		tea.Raw(ansi.SetModeLightDark),
		tea.Raw(ansi.RequestModeLightDark),
	)
}

// stopThemeNotificationsCmd turns the mode back off, so the terminal is left as
// we found it on quit.
func stopThemeNotificationsCmd() tea.Cmd {
	return tea.Raw(ansi.ResetModeLightDark)
}

// viewMode represents the current view/mode of the application
type viewMode int

const (
	modeList   viewMode = iota // Main task list view
	modeAdd                    // Adding new task
	modeHelp                   // Help screen
	modeDelete                 // Delete confirmation dialog
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
	savePending   bool          // Whether a save is scheduled but not yet executed
	styles        styles        // Styles for the active (light/dark) terminal theme
	isDark        bool          // Whether the terminal background is dark
	themePushed   bool          // Terminal pushes theme changes (DEC mode 2031), so no polling
	// Delete confirmation
	deleteTaskID string // ID of task to delete (when in modeDelete)
	// Add/Edit task form fields
	editTaskID      string         // ID of task being edited (empty if adding new task)
	addField        int            // Current field in add mode (0=task, 1=priority, 2=due)
	addCursor       int            // Cursor position within priority/due lists
	addPriority     model.Priority // Selected priority for new task
	addDueSelection int            // Selected due date option (0=today, 1=tomorrow, 2=this week, 3=next week)
}

// NewModel creates a new Model with the given task list and storage
func NewModel(taskList *model.TaskList, storage *storage.Storage) Model {
	s := spinner.New()
	s.Spinner = spinner.MiniDot

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
		// Assume a dark background until the terminal tells us otherwise.
		styles: newStyles(true),
		isDark: true,
	}
}

// Init initializes the model (required by BubbleTea)
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, themeNotificationsCmd(), pollTheme())
}

// StartInAddMode configures the model to start in add task mode
func (m Model) StartInAddMode() Model {
	m.mode = modeAdd
	m.input = ""
	m.addField = 0
	m.addCursor = 0
	m.addPriority = model.PriorityLow
	m.addDueSelection = 0
	m.editTaskID = ""
	return m
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

// scheduleSave returns a command that will trigger a save after a debounce delay
func (m *Model) scheduleSave() tea.Cmd {
	m.savePending = true
	return tea.Tick(saveDebounceDuration, func(t time.Time) tea.Msg {
		return saveMsg{}
	})
}

// performSave executes the actual save to storage
func (m *Model) performSave() {
	m.savePending = false
	if err := m.storage.Save(m.taskList); err != nil {
		m.err = err
	}
}

// applyTheme switches the palette to match the terminal background, rebuilding
// the styles only when the light/dark mode actually changed. Both the polled
// reply and the pushed color scheme events land here.
func (m *Model) applyTheme(isDark bool) {
	if isDark == m.isDark {
		return
	}
	m.isDark = isDark
	m.styles = newStyles(isDark)
}

// useThemeNotifications records that the terminal supports DEC mode 2031 and
// will push theme changes, which makes the background color poll redundant.
// Once this is set, themeTickMsg stops rescheduling itself.
func (m *Model) useThemeNotifications() {
	m.themePushed = true
}

// clearMessages clears error and info messages
func (m *Model) clearMessages() {
	m.err = nil
	m.message = ""
}
