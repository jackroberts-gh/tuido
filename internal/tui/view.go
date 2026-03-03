package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/jackroberts-gh/tuido/internal/model"
)

// View renders the current view based on the mode
func (m Model) View() string {
	switch m.mode {
	case modeList:
		return m.renderList()
	case modeAdd:
		return m.renderAddTask()
	case modeEditPriority:
		return m.renderPrioritySelector()
	case modeEditDueDate:
		return m.renderDueDateInput()
	case modeHelp:
		return m.renderHelp()
	default:
		return "Unknown mode"
	}
}

// renderList renders the main task list view
func (m Model) renderList() string {
	var b strings.Builder

	// Task list with responsive width
	listStyle := taskListStyle
	if m.width > 0 {
		listStyle = listStyle.Width(m.width - 6)
	}

	visibleTasks := m.getVisibleTasks()
	if len(visibleTasks) == 0 {
		emptyMsg := "No tasks yet! Press 'a' to add your first task."
		b.WriteString(listStyle.Render(emptyStateStyle.Render(emptyMsg)))
	} else {
		taskContent := strings.Builder{}
		for i, task := range visibleTasks {
			taskContent.WriteString(m.renderTask(task, i == m.cursor))
			if i < len(visibleTasks)-1 {
				taskContent.WriteString("\n")
			}
		}
		b.WriteString(listStyle.Render(taskContent.String()))
	}

	// Messages
	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(errorStyle.Render(fmt.Sprintf(" ✗ %v ", m.err)))
	}
	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(successStyle.Render(fmt.Sprintf(" ✓ %s ", m.message)))
	}

	// Footer with shortcuts
	b.WriteString("\n")
	footer := m.buildFooter([]footerItem{
		{"a", "add"},
		{"space", "toggle"},
		{"d", "delete"},
		{"p", "priority"},
		{"e", "date"},
		{"t", "filter"},
		{"?", "help"},
		{"q", "quit"},
	})
	footerStyleWithWidth := footerStyle
	if m.width > 0 {
		footerStyleWithWidth = footerStyleWithWidth.Width(m.width - 4)
	}
	b.WriteString(footerStyleWithWidth.Render(footer))

	return b.String()
}

type footerItem struct {
	key  string
	desc string
}

func (m Model) buildFooter(items []footerItem) string {
	parts := make([]string, len(items))
	for i, item := range items {
		parts[i] = fmt.Sprintf("%s %s", footerKeyStyle.Render(item.key), item.desc)
	}
	return strings.Join(parts, footerSepStyle.Render(" • "))
}

// renderTask renders a single task
func (m Model) renderTask(task model.Task, selected bool) string {
	// Cursor indicator
	cursor := " "
	if selected {
		cursor = cursorStyle.Render("▶")
	}

	// Checkbox
	var checkbox string
	if task.Completed {
		checkbox = checkboxCompletedStyle.Render("[x]")
	} else if selected {
		checkbox = checkboxSelectedStyle.Render("[ ]")
	} else {
		checkbox = checkboxStyle.Render("[ ]")
	}

	// Task text
	var taskText string
	if task.Completed {
		taskText = completedTaskStyle.Render(task.Text)
	} else {
		if selected {
			taskText = checkboxSelectedStyle.Render(task.Text)
		} else {
			taskText = taskTextStyle.Render(task.Text)
		}
	}

	// Due date
	dueDateStr := ""
	if task.DueDate != nil {
		if task.Completed {
			dueDateStr = " " + completedTaskStyle.Render(m.formatDueDateText(&task))
		} else if selected {
			dueDateStr = " " + checkboxSelectedStyle.Render(m.formatDueDateText(&task))
		} else {
			dueDateStr = " " + m.formatDueDate(&task)
		}
	}

	// Build the line
	line := fmt.Sprintf("%s %s %s%s", cursor, checkbox, taskText, dueDateStr)

	// Apply padding
	if selected {
		return selectedTaskStyle.Render(line)
	}
	return taskStyle.Render(line)
}

// formatDueDateText returns the due date text without styling
func (m Model) formatDueDateText(task *model.Task) string {
	if task.DueDate == nil {
		return ""
	}

	now := time.Now()
	dueDate := *task.DueDate

	// Calculate difference
	diff := dueDate.Sub(now)
	days := int(diff.Hours() / 24)

	if days < 0 {
		return fmt.Sprintf("⚠ overdue %dd", -days)
	} else if days == 0 {
		return "⚠ due today"
	} else if days == 1 {
		return "⏰ tomorrow"
	} else if days <= 7 {
		return fmt.Sprintf("⏰ %dd", days)
	} else {
		return fmt.Sprintf("📅 %s", dueDate.Format("Jan 02"))
	}
}

// formatDueDate formats the due date with appropriate styling
func (m Model) formatDueDate(task *model.Task) string {
	if task.DueDate == nil {
		return ""
	}

	now := time.Now()
	dueDate := *task.DueDate

	// Calculate difference
	diff := dueDate.Sub(now)
	days := int(diff.Hours() / 24)

	var dateStr string
	if days < 0 {
		// Overdue
		dateStr = fmt.Sprintf("⚠ overdue %dd", -days)
		return overdueStyle.Render(dateStr)
	} else if days == 0 {
		// Due today
		dateStr = "⚠ due today"
		return overdueStyle.Render(dateStr)
	} else if days == 1 {
		// Due tomorrow
		dateStr = "⏰ tomorrow"
		return dueSoonStyle.Render(dateStr)
	} else if days <= 7 {
		// Due within a week
		dateStr = fmt.Sprintf("⏰ %dd", days)
		return dueSoonStyle.Render(dateStr)
	} else {
		// Due later
		dateStr = fmt.Sprintf("📅 %s", dueDate.Format("Jan 02"))
		return dueDateStyle.Render(dateStr)
	}
}

// renderAddTask renders the add task input view
func (m Model) renderAddTask() string {
	var dialog strings.Builder

	dialog.WriteString(dialogTitleStyle.Render("Add New Task"))
	dialog.WriteString("\n\n")

	dialog.WriteString(promptStyle.Render("Task description:"))
	dialog.WriteString("\n")

	inputText := m.input + "█"
	dialog.WriteString(inputBoxStyle.Render(inputText))
	dialog.WriteString("\n\n")

	dialog.WriteString(hintStyle.Render("↵ Enter to add • Esc to cancel"))

	boxStyle := dialogBoxStyle
	if m.width > 0 {
		boxStyle = boxStyle.MaxWidth(m.width - 4)
	}
	return boxStyle.Render(dialog.String())
}

// renderPrioritySelector renders the priority selection view
func (m Model) renderPrioritySelector() string {
	var dialog strings.Builder

	dialog.WriteString(dialogTitleStyle.Render("Select Priority"))
	dialog.WriteString("\n\n")

	dialog.WriteString(priorityLowStyle.Render("  1 ● Low"))
	dialog.WriteString("\n\n")
	dialog.WriteString(priorityMediumStyle.Render("  2 ● Medium"))
	dialog.WriteString("\n\n")
	dialog.WriteString(priorityHighStyle.Render("  3 ● High"))
	dialog.WriteString("\n\n")

	dialog.WriteString(hintStyle.Render("Press 1, 2, or 3 • Esc to cancel"))

	boxStyle := dialogBoxStyle
	if m.width > 0 {
		boxStyle = boxStyle.MaxWidth(m.width - 4)
	}
	return boxStyle.Render(dialog.String())
}

// renderDueDateInput renders the due date input view
func (m Model) renderDueDateInput() string {
	var dialog strings.Builder

	dialog.WriteString(dialogTitleStyle.Render("Set Due Date"))
	dialog.WriteString("\n\n")

	dialog.WriteString(promptStyle.Render("Due date:"))
	dialog.WriteString("\n")

	inputText := m.input + "█"
	dialog.WriteString(inputBoxStyle.Render(inputText))
	dialog.WriteString("\n\n")

	dialog.WriteString(hintStyle.Render("Examples: tomorrow, 3d, 1w, 2026-03-15"))
	dialog.WriteString("\n")
	dialog.WriteString(hintStyle.Render("↵ Enter to set • Esc to cancel • Empty to clear"))

	boxStyle := dialogBoxStyle
	if m.width > 0 {
		boxStyle = boxStyle.MaxWidth(m.width - 4)
	}
	return boxStyle.Render(dialog.String())
}

// renderHelp renders the help screen
func (m Model) renderHelp() string {
	var help strings.Builder

	help.WriteString(helpHeaderStyle.Render("⌨  Keyboard Shortcuts"))
	help.WriteString("\n\n")

	shortcuts := []struct {
		key  string
		desc string
	}{
		{"↑ / k", "Move cursor up"},
		{"↓ / j", "Move cursor down"},
		{"space", "Toggle task completion"},
		{"a", "Add new task"},
		{"d", "Delete selected task"},
		{"p", "Change priority"},
		{"e", "Edit due date"},
		{"t", "Toggle completed tasks visibility"},
		{"?", "Show this help"},
		{"q / Ctrl+C", "Quit application"},
	}

	for _, shortcut := range shortcuts {
		line := fmt.Sprintf("%-12s  %s",
			helpKeyStyle.Render(shortcut.key),
			helpDescStyle.Render(shortcut.desc))
		help.WriteString(line)
		help.WriteString("\n")
	}

	help.WriteString("\n")
	help.WriteString(hintStyle.Render("Press any key to return"))

	boxStyle := dialogBoxStyle
	if m.width > 0 {
		boxStyle = boxStyle.MaxWidth(m.width - 4)
	}
	return boxStyle.Render(help.String())
}
