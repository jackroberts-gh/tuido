package tui

import (
	"fmt"
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
		return m, tea.ClearScreen

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	default:
		// Update spinner
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
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
	case modeHelp:
		return m.handleHelpMode()
	}

	return m, nil
}

// handleListMode handles keyboard input in list view mode
func (m Model) handleListMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		return m, tea.Quit

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

	case " ":
		// Cycle through status: not started -> in-progress -> completed -> not started
		task := m.getCurrentTask()
		if task != nil {
			wasCompleted := task.Completed

			// Before cycling, find the next task in the list (for cursor movement when completing)
			var nextTaskID string
			if !wasCompleted && task.InProgress && m.showCompleted {
				// We're about to complete a task (in-progress -> completed)
				// Find the next uncompleted task
				visibleTasksBefore := m.getVisibleTasks()
				for i := m.cursor + 1; i < len(visibleTasksBefore); i++ {
					if !visibleTasksBefore[i].Completed {
						nextTaskID = visibleTasksBefore[i].ID
						break
					}
				}
			}

			m.taskList.CycleStatus(task.ID)

			// Check if task just became completed
			taskAfter := m.getCurrentTask()
			isNowCompleted := taskAfter != nil && taskAfter.Completed && !wasCompleted

			visibleTasks := m.getVisibleTasks()
			maxCursor := len(visibleTasks) - 1

			if isNowCompleted && m.showCompleted {
				// Task just completed - move cursor to the next uncompleted task
				if nextTaskID != "" {
					// Find where the next task ended up after re-sorting
					for i, t := range visibleTasks {
						if t.ID == nextTaskID {
							m.cursor = i
							m.saveToStorage()
							m.lastKey = ""
							return m, nil
						}
					}
				}
				// If no next task was found, go to first uncompleted task
				for i, t := range visibleTasks {
					if !t.Completed {
						m.cursor = i
						m.saveToStorage()
						m.lastKey = ""
						return m, nil
					}
				}
				// If no uncompleted tasks, go to end
				m.cursor = maxCursor
			} else {
				// Ensure cursor stays in bounds
				if m.cursor > maxCursor && maxCursor >= 0 {
					m.cursor = maxCursor
				}
			}

			m.saveToStorage()
		}
		m.lastKey = ""

	case "a":
		// Enter add mode
		m.mode = modeAdd
		m.input = ""
		m.addField = 0
		m.addCursor = 0 // Start at Low priority
		m.addPriority = model.PriorityLow
		m.addDueSelection = 0
		m.clearMessages()
		m.lastKey = ""

	case "d":
		// Check if this is the start of "sd" sequence or just delete
		// If lastKey is empty and we're pressing "d", it could be either
		// We'll treat standalone "d" as delete, and "sd" requires "s" first
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
		if m.addCursor < 4 {
			m.addCursor++
		}

	case "enter":
		m.addDueSelection = m.addCursor
		m.addTask()

	case "1", "2", "3", "4", "5":
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
	case 4: // No due date
		return nil
	}

	return &result
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
