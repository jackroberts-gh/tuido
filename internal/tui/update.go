package tui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
	"github.com/jackroberts-gh/tuido/v2/internal/model"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, tea.ClearScreen

	case saveMsg:
		// Perform the actual save to storage
		m.performSave()
		return m, nil

	case themeTickMsg:
		// Fallback path for terminals without DEC mode 2031. If the terminal
		// has since told us it supports notifications, drop the poll loop here
		// and rely on the pushed events instead.
		if m.themePushed {
			return m, nil
		}
		return m, pollTheme()

	case tea.ModeReportMsg:
		// Reply to our mode 2031 query. A terminal that recognises the mode
		// will push theme changes, so the poll can be switched off. Terminals
		// without support answer "not recognized" (or never answer at all) and
		// keep polling.
		if dm, ok := msg.Mode.(ansi.DECMode); ok && int(dm) == themeNotificationMode {
			if !msg.Value.IsNotRecognized() {
				m.useThemeNotifications()
			}
		}
		return m, nil

	case tea.BackgroundColorMsg:
		// Reply to our own background color query (the poll, or startup)
		m.applyTheme(msg.IsDark())
		return m, nil

	case uv.DarkColorSchemeEvent:
		// Pushed by the terminal the moment it switched to a dark theme.
		// Receiving this at all proves notifications work, so stop polling
		// even if the mode report never arrived.
		m.useThemeNotifications()
		m.applyTheme(true)
		return m, nil

	case uv.LightColorSchemeEvent:
		// Pushed by the terminal the moment it switched to a light theme
		m.useThemeNotifications()
		m.applyTheme(false)
		return m, nil

	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)

	default:
		// Update spinner
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

// handleKeyPress routes key presses to mode-specific handlers
func (m Model) handleKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	// Global quit keys
	if msg.String() == "ctrl+c" {
		return m, quitCmd()
	}

	switch m.mode {
	case modeList:
		return m.handleListMode(msg)
	case modeAdd:
		return m.handleAddMode(msg)
	case modeHelp:
		return m.handleHelpMode()
	case modeDelete:
		return m.handleDeleteMode(msg)
	}

	return m, nil
}

// quitCmd disables theme change notifications before quitting, so the terminal
// is not left with mode 2031 enabled.
func quitCmd() tea.Cmd {
	return tea.Sequence(stopThemeNotificationsCmd(), tea.Quit)
}

// handleListMode handles keyboard input in list view mode
func (m Model) handleListMode(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	m.clearMessages()

	visibleTasks := m.getVisibleTasks()
	maxCursor := len(visibleTasks) - 1

	key := msg.String()

	// Handle two-key sequences like "sd" and "sp"
	if m.lastKey == "s" {
		switch key {
		case "d":
			m = m.cycleSortMode(sortDueDate, sortDueDateReverse)
			return m, nil
		case "p":
			m = m.cycleSortMode(sortPriority, sortPriorityReverse)
			return m, nil
		default:
			// Invalid sequence, clear lastKey
			m.lastKey = ""
		}
	}

	switch key {
	case "q":
		return m, quitCmd()

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
		m.lastKey = ""

	case "down", "j":
		if maxCursor >= 0 && m.cursor < maxCursor {
			m.cursor++
		}
		m.lastKey = ""

	case "right", "l":
		// Step status forward, stopping at completed
		return m.stepCurrentTask(true)

	case "left", "h":
		// Step status back, stopping at not started ("undo" progress)
		return m.stepCurrentTask(false)

	case "a":
		// Enter add mode
		m.mode = modeAdd
		m.editTaskID = "" // Clear edit task ID (we're adding, not editing)
		m.input = ""
		m.addField = 0
		m.addCursor = 0 // Start at Low priority
		m.addPriority = model.PriorityLow
		m.addDueSelection = 0
		m.clearMessages()
		m.lastKey = ""

	case "e":
		// Enter edit mode for the current task
		task := m.getCurrentTask()
		if task != nil {
			m.mode = modeAdd
			m.editTaskID = task.ID
			m.input = task.Text
			m.addField = 0
			// Set cursor position based on task priority
			switch task.Priority {
			case model.PriorityLow:
				m.addCursor = 0
			case model.PriorityMedium:
				m.addCursor = 1
			case model.PriorityHigh:
				m.addCursor = 2
			}
			m.addPriority = task.Priority
			// Calculate due date selection from task due date
			m.addDueSelection = m.calculateDueSelectionFromDate(task.DueDate)
			m.clearMessages()
		}
		m.lastKey = ""

	case "d":
		// Show delete confirmation dialog
		task := m.getCurrentTask()
		if task != nil {
			m.deleteTaskID = task.ID
			m.mode = modeDelete
			m.clearMessages()
		}
		m.lastKey = ""

	case "s":
		// Start sort sequence - wait for next key
		m.lastKey = "s"

	case "t":
		// Toggle show completed - preserve cursor position if possible
		currentTask := m.getCurrentTask()
		var currentTaskID string
		if currentTask != nil {
			currentTaskID = currentTask.ID
		}

		m.showCompleted = !m.showCompleted

		// Try to find the same task in the new visible list
		if currentTaskID != "" {
			visibleTasks := m.getVisibleTasks()
			for i, task := range visibleTasks {
				if task.ID == currentTaskID {
					m.cursor = i
					m.lastKey = ""
					return m, nil
				}
			}
		}

		// Task not found (was filtered out), default to 0
		m.cursor = 0
		m.lastKey = ""

	case "?":
		// Show help
		m.mode = modeHelp
		m.clearMessages()
		m.lastKey = ""

	default:
		m.lastKey = ""
	}

	return m, nil
}

// stepCurrentTask moves the selected task one step through its status, forward
// when advance is true and back otherwise. Neither direction wraps, so the ends
// of the progression are stable: pressing right on a completed task or left on
// a not-started task does nothing.
//
// The cursor follows the task rather than moving on to the next one (as the
// forward cycle does), because stepping is used to adjust a specific task - the
// task itself may shift position if completing or uncompleting re-sorts it.
func (m Model) stepCurrentTask(advance bool) (tea.Model, tea.Cmd) {
	task := m.getCurrentTask()
	if task == nil {
		m.lastKey = ""
		return m, nil
	}

	taskID := task.ID
	var changed bool
	if advance {
		changed = m.taskList.AdvanceStatus(taskID)
	} else {
		changed = m.taskList.RegressStatus(taskID)
	}
	m.lastKey = ""
	if !changed {
		// Already at the end of the progression - nothing to save.
		return m, nil
	}

	visibleTasks := m.getVisibleTasks()
	maxCursor := len(visibleTasks) - 1

	// Follow the task to its new position, if it is still visible.
	found := false
	for i, t := range visibleTasks {
		if t.ID == taskID {
			m.cursor = i
			found = true
			break
		}
	}
	// Filtered out (completed while completed tasks are hidden) - keep the
	// cursor in range.
	if !found && m.cursor > maxCursor && maxCursor >= 0 {
		m.cursor = maxCursor
	}

	return m, m.scheduleSave()
}

// handleAddMode handles keyboard input in add task mode
func (m Model) handleAddMode(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch m.addField {
	case 0:
		// Task input field
		return m.handleAddTaskInput(msg)
	case 1:
		// Priority selection
		return m.handleAddPrioritySelect(msg)
	case 2:
		// Due date selection
		return m.handleAddDueSelect(msg)
	}
	return m, nil
}

// handleAddTaskInput handles task input in add mode
func (m Model) handleAddTaskInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeList
		m.input = ""
		m.addField = 0
		m.clearMessages()

	case "enter":
		if strings.TrimSpace(m.input) != "" {
			m.addField = 1
			// Set cursor based on current priority (for edit mode)
			switch m.addPriority {
			case model.PriorityLow:
				m.addCursor = 0
			case model.PriorityMedium:
				m.addCursor = 1
			case model.PriorityHigh:
				m.addCursor = 2
			}
			m.clearMessages()
		} else {
			m.err = fmt.Errorf("task cannot be empty")
		}

	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}

	default:
		// Text holds the actual characters typed (empty for non-text keys)
		m.input += msg.Text
	}
	return m, nil
}

// handleAddPrioritySelect handles priority selection in add mode
func (m Model) handleAddPrioritySelect(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeList
		m.input = ""
		m.addField = 0
		m.clearMessages()

	case "up", "k":
		if m.addCursor > 0 {
			m.addCursor--
		}

	case "down", "j":
		if m.addCursor < 2 {
			m.addCursor++
		}

	case "enter":
		// Update priority based on cursor position
		switch m.addCursor {
		case 0:
			m.addPriority = model.PriorityLow
		case 1:
			m.addPriority = model.PriorityMedium
		case 2:
			m.addPriority = model.PriorityHigh
		}
		m.addField = 2
		m.addCursor = m.addDueSelection // Set cursor to current due date selection

	case "1":
		m.addPriority = model.PriorityLow
		m.addField = 2
		m.addCursor = m.addDueSelection

	case "2":
		m.addPriority = model.PriorityMedium
		m.addField = 2
		m.addCursor = m.addDueSelection

	case "3":
		m.addPriority = model.PriorityHigh
		m.addField = 2
		m.addCursor = m.addDueSelection
	}
	return m, nil
}

// handleAddDueSelect handles due date selection in add mode
func (m Model) handleAddDueSelect(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeList
		m.input = ""
		m.addField = 0
		m.clearMessages()

	case "up", "k":
		if m.addCursor > 0 {
			m.addCursor--
		}

	case "down", "j":
		if m.addCursor < 4 {
			m.addCursor++
		}

	case "enter":
		m.addDueSelection = m.addCursor
		m.addTask()
		return m, m.scheduleSave()

	case "1", "2", "3", "4", "5":
		m.addDueSelection = int(msg.String()[0] - '1')
		m.addTask()
		return m, m.scheduleSave()
	}
	return m, nil
}

// addTask completes the add/edit flow and creates or updates the task
func (m *Model) addTask() {
	dueDate := m.calculateDueDateFromSelection()

	if m.editTaskID != "" {
		// Editing existing task
		task := m.taskList.GetByID(m.editTaskID)
		if task != nil {
			task.Text = m.input
			task.Priority = m.addPriority
			task.DueDate = dueDate
		}
	} else {
		// Adding new task
		m.taskList.Add(m.input, m.addPriority, dueDate)
	}

	m.mode = modeList
	m.editTaskID = ""
	m.input = ""
	m.addField = 0
	m.addCursor = 0
}

// calculateDueDateFromSelection calculates the due date based on the selection
func (m Model) calculateDueDateFromSelection() *time.Time {
	now := time.Now()
	var result time.Time

	switch m.addDueSelection {
	case 0: // Today
		result = now
	case 1: // Tomorrow
		result = now.Add(24 * time.Hour)
	case 2: // This week (7 days)
		result = now.Add(7 * 24 * time.Hour)
	case 3: // Next week (14 days)
		result = now.Add(14 * 24 * time.Hour)
	case 4: // No due date
		return nil
	}

	return &result
}

// calculateDueSelectionFromDate converts a due date back to a selection index
func (m Model) calculateDueSelectionFromDate(dueDate *time.Time) int {
	if dueDate == nil {
		return 4 // No due date
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	dueDay := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, dueDate.Location())

	daysDiff := int(dueDay.Sub(today).Hours() / 24)

	switch {
	case daysDiff == 0:
		return 0 // Today
	case daysDiff == 1:
		return 1 // Tomorrow
	case daysDiff >= 2 && daysDiff <= 7:
		return 2 // This week
	case daysDiff >= 8 && daysDiff <= 14:
		return 3 // Next week
	default:
		// For dates outside the range, default to "This week"
		return 2
	}
}

// cycleSortMode cycles through sort modes and preserves cursor position
func (m Model) cycleSortMode(primarySort, reverseSort sortMode) Model {
	// Save current task ID before sorting
	currentTask := m.getCurrentTask()
	var currentTaskID string
	if currentTask != nil {
		currentTaskID = currentTask.ID
	}

	// Cycle through sort: none -> primary -> reverse -> none
	switch m.sortBy {
	case sortNone:
		m.sortBy = primarySort
	case primarySort:
		m.sortBy = reverseSort
	case reverseSort:
		m.sortBy = sortNone
	default:
		// If in a different sort mode, switch to primary
		m.sortBy = primarySort
	}

	// Try to restore cursor to same task
	if currentTaskID != "" {
		visibleTasks := m.getVisibleTasks()
		for i, task := range visibleTasks {
			if task.ID == currentTaskID {
				m.cursor = i
				m.lastKey = ""
				return m
			}
		}
	}

	m.cursor = 0
	m.lastKey = ""
	return m
}

// handleHelpMode handles keyboard input in help mode
func (m Model) handleHelpMode() (tea.Model, tea.Cmd) {
	// Any key returns to list mode
	m.mode = modeList
	return m, nil
}

// handleDeleteMode handles keyboard input in delete confirmation mode
func (m Model) handleDeleteMode(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		// Confirm deletion
		m.taskList.Remove(m.deleteTaskID)

		// Adjust cursor if needed
		visibleTasks := m.getVisibleTasks()
		if m.cursor >= len(visibleTasks) && m.cursor > 0 {
			m.cursor--
		}

		// Return to list mode
		m.mode = modeList
		m.deleteTaskID = ""
		return m, m.scheduleSave()

	case "n", "N", "esc":
		// Cancel deletion
		m.mode = modeList
		m.deleteTaskID = ""
		return m, nil

	default:
		// Any other key cancels
		m.mode = modeList
		m.deleteTaskID = ""
		return m, nil
	}
}
