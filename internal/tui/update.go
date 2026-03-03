package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackroberts-gh/tuido/internal/model"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

// handleKeyPress routes key presses to mode-specific handlers
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global quit keys
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	switch m.mode {
	case modeList:
		return m.handleListMode(msg)
	case modeAdd:
		return m.handleAddMode(msg)
	case modeEditPriority:
		return m.handlePriorityMode(msg)
	case modeEditDueDate:
		return m.handleDueDateMode(msg)
	case modeHelp:
		return m.handleHelpMode(msg)
	}

	return m, nil
}

// handleListMode handles keyboard input in list view mode
func (m Model) handleListMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.clearMessages()

	visibleTasks := m.getVisibleTasks()
	maxCursor := len(visibleTasks) - 1

	switch msg.String() {
	case "q":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if maxCursor >= 0 && m.cursor < maxCursor {
			m.cursor++
		}

	case " ":
		// Toggle completion
		task := m.getCurrentTask()
		if task != nil {
			m.taskList.Toggle(task.ID)
			m.saveToStorage()
		}

	case "a":
		// Enter add mode
		m.mode = modeAdd
		m.input = ""
		m.addField = 0
		m.addCursor = 0 // Start at Low priority
		m.addPriority = model.PriorityLow
		m.addDueSelection = 0
		m.clearMessages()

	case "d":
		// Delete task
		task := m.getCurrentTask()
		if task != nil {
			m.taskList.Remove(task.ID)
			m.saveToStorage()

			// Adjust cursor if needed
			visibleTasks = m.getVisibleTasks()
			if m.cursor >= len(visibleTasks) && m.cursor > 0 {
				m.cursor--
			}
		}

	case "p":
		// Enter priority edit mode
		task := m.getCurrentTask()
		if task != nil {
			m.mode = modeEditPriority
			m.selectedID = task.ID
			m.clearMessages()
		}

	case "e":
		// Enter due date edit mode
		task := m.getCurrentTask()
		if task != nil {
			m.mode = modeEditDueDate
			m.selectedID = task.ID
			m.input = ""
			m.clearMessages()
		}

	case "t":
		// Toggle show completed
		m.showCompleted = !m.showCompleted
		// Reset cursor to 0 to avoid out of bounds
		m.cursor = 0

	case "?":
		// Show help
		m.mode = modeHelp
		m.clearMessages()
	}

	return m, nil
}

// handleAddMode handles keyboard input in add task mode
func (m Model) handleAddMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
func (m Model) handleAddTaskInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeList
		m.input = ""
		m.addField = 0
		m.clearMessages()

	case "enter":
		if strings.TrimSpace(m.input) != "" {
			m.addField = 1
			m.addCursor = 0 // Reset cursor to Low priority
			m.clearMessages()
		} else {
			m.err = fmt.Errorf("task cannot be empty")
		}

	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}

	default:
		if len(msg.String()) == 1 {
			m.input += msg.String()
		}
	}
	return m, nil
}

// handleAddPrioritySelect handles priority selection in add mode
func (m Model) handleAddPrioritySelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		m.addCursor = 0 // Reset cursor to Today

	case "1":
		m.addPriority = model.PriorityLow
		m.addField = 2
		m.addCursor = 0

	case "2":
		m.addPriority = model.PriorityMedium
		m.addField = 2
		m.addCursor = 0

	case "3":
		m.addPriority = model.PriorityHigh
		m.addField = 2
		m.addCursor = 0
	}
	return m, nil
}

// handleAddDueSelect handles due date selection in add mode
func (m Model) handleAddDueSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		if m.addCursor < 3 {
			m.addCursor++
		}

	case "enter":
		m.addDueSelection = m.addCursor
		m.addTask()

	case "1", "2", "3", "4":
		m.addDueSelection = int(msg.String()[0] - '1')
		m.addTask()
	}
	return m, nil
}

// addTask completes the add flow and creates the task
func (m *Model) addTask() {
	dueDate := m.calculateDueDateFromSelection()
	m.taskList.Add(m.input, m.addPriority, dueDate)
	m.saveToStorage()
	m.mode = modeList
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
	}

	return &result
}

// handlePriorityMode handles keyboard input in priority selection mode
func (m Model) handlePriorityMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel priority change
		m.mode = modeList
		m.selectedID = ""
		m.clearMessages()

	case "1":
		// Set low priority
		if m.taskList.UpdatePriority(m.selectedID, model.PriorityLow) {
			m.saveToStorage()
		}
		m.mode = modeList
		m.selectedID = ""

	case "2":
		// Set medium priority
		if m.taskList.UpdatePriority(m.selectedID, model.PriorityMedium) {
			m.saveToStorage()
		}
		m.mode = modeList
		m.selectedID = ""

	case "3":
		// Set high priority
		if m.taskList.UpdatePriority(m.selectedID, model.PriorityHigh) {
			m.saveToStorage()
		}
		m.mode = modeList
		m.selectedID = ""
	}

	return m, nil
}

// handleDueDateMode handles keyboard input in due date edit mode
func (m Model) handleDueDateMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel due date edit
		m.mode = modeList
		m.selectedID = ""
		m.input = ""
		m.clearMessages()

	case "enter":
		// Parse and set due date
		if strings.TrimSpace(m.input) == "" {
			// Empty input means remove due date
			m.taskList.UpdateDueDate(m.selectedID, nil)
			m.saveToStorage()
			m.mode = modeList
			m.selectedID = ""
			m.input = ""
		} else {
			dueDate, err := parseDueDate(m.input)
			if err != nil {
				m.err = err
			} else {
				m.taskList.UpdateDueDate(m.selectedID, dueDate)
				m.saveToStorage()
				m.mode = modeList
				m.selectedID = ""
				m.input = ""
			}
		}

	case "backspace":
		// Remove last character
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}

	default:
		// Add character to input
		if len(msg.String()) == 1 {
			m.input += msg.String()
		}
	}

	return m, nil
}

// handleHelpMode handles keyboard input in help mode
func (m Model) handleHelpMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Any key returns to list mode
	m.mode = modeList
	return m, nil
}

// parseDueDate parses a due date string in various formats
func parseDueDate(input string) (*time.Time, error) {
	input = strings.TrimSpace(strings.ToLower(input))
	now := time.Now()

	// Relative dates (e.g., "3d", "1w", "2m")
	if len(input) >= 2 {
		numStr := input[:len(input)-1]
		unit := input[len(input)-1:]

		if num, err := strconv.Atoi(numStr); err == nil {
			var duration time.Duration
			switch unit {
			case "d":
				duration = time.Duration(num) * 24 * time.Hour
			case "w":
				duration = time.Duration(num) * 7 * 24 * time.Hour
			case "m":
				// Approximate month as 30 days
				duration = time.Duration(num*30) * 24 * time.Hour
			default:
				goto tryOtherFormats
			}
			result := now.Add(duration)
			return &result, nil
		}
	}

tryOtherFormats:
	// Natural language
	switch input {
	case "today":
		result := now
		return &result, nil
	case "tomorrow":
		result := now.Add(24 * time.Hour)
		return &result, nil
	case "next week":
		result := now.Add(7 * 24 * time.Hour)
		return &result, nil
	}

	// Absolute date formats
	formats := []string{
		"2006-01-02",
		"01/02/2006",
		"Jan 02, 2006",
		"Jan 02",
		"01/02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, input); err == nil {
			// For formats without year, use current year
			if !strings.Contains(format, "2006") {
				t = time.Date(now.Year(), t.Month(), t.Day(), 0, 0, 0, 0, now.Location())
				// If the date is in the past, assume next year
				if t.Before(now) {
					t = t.AddDate(1, 0, 0)
				}
			}
			return &t, nil
		}
	}

	return nil, fmt.Errorf("invalid date format. Try: tomorrow, 3d, 1w, 2026-03-15")
}
