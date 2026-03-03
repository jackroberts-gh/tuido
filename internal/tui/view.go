package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/jackroberts-gh/tuido/internal/model"
)

// View renders the current view based on the mode
func (m Model) View() string {
	switch m.mode {
	case modeList:
		return m.renderList()
	case modeAdd:
		return m.renderAddTask()
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
	taskContent := strings.Builder{}

	// Render table header
	taskContent.WriteString(m.renderTableHeader())
	taskContent.WriteString(m.renderHeaderSeparator())

	// Render tasks as table rows
	if len(visibleTasks) == 0 {
		// Show hint when no tasks
		taskContent.WriteString("\n")
		hintText := "Hit a to add a new task"
		hintStyle := lipgloss.NewStyle().
			Foreground(textDim).
			PaddingLeft(1)
		taskContent.WriteString(hintStyle.Render("  " + hintText))
	} else {
		for i, task := range visibleTasks {
			taskContent.WriteString("\n")
			taskContent.WriteString(m.renderTask(task, i == m.cursor))
		}
	}
	b.WriteString(listStyle.Render(taskContent.String()))

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
		{"space", "cycle status"},
		{"d", "delete"},
		{"sp", "sort priority"},
		{"sd", "sort date"},
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
		parts[i] = fmt.Sprintf("%s %s", footerKeyStyle.Render(item.key), footerDescStyle.Render(item.desc))
	}
	return strings.Join(parts, footerSepStyle.Render(" • "))
}

// renderTableHeader renders the table header
func (m Model) renderTableHeader() string {
	// Calculate column widths
	taskWidth := m.getTaskColumnWidth()
	priorityWidth := 10
	dueDateWidth := 12

	// Build header with proper spacing
	// Account for: cursor(1) + space(1) = 2 chars before checkbox/Task header
	// Task column includes checkbox(3) + space(1) + text = taskWidth

	// Build the full header string first, then apply styling
	header := fmt.Sprintf("  %-*s  %-*s  %-*s",
		taskWidth, "Task",
		priorityWidth, "Priority",
		dueDateWidth, "Due Date")

	// Use muted text color for the entire header (so Due Date isn't in accent color)
	headerStyle := lipgloss.NewStyle().
		Foreground(textMuted).
		Bold(true).
		PaddingLeft(1)

	return headerStyle.Render(header)
}

// renderHeaderSeparator renders a line separator below the header
func (m Model) renderHeaderSeparator() string {
	// Calculate total width for separator line
	taskWidth := m.getTaskColumnWidth()
	totalWidth := 2 + taskWidth + 2 + 10 + 2 + 12 // spacing included

	separator := strings.Repeat("─", totalWidth)

	separatorStyle := lipgloss.NewStyle().
		Foreground(border).
		PaddingLeft(1)

	return separatorStyle.Render(separator)
}

// getTaskColumnWidth calculates the width for the task column based on terminal width
func (m Model) getTaskColumnWidth() int {
	if m.width > 0 {
		// Reserve space for borders, padding, priority and due date columns
		// Border + padding = ~10, Priority = 10, Due Date = 12, spacing = 4
		reserved := 36
		taskWidth := m.width - reserved
		if taskWidth < 20 {
			taskWidth = 20 // Minimum task column width
		}
		// No maximum limit - let it grow with terminal width
		return taskWidth
	}
	return 40 // Default width
}

// renderTask renders a single task as a table row
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
	} else if task.InProgress {
		// Show spinner for in-progress tasks
		spinnerView := m.spinner.View()
		if selected {
			checkbox = checkboxSelectedStyle.Render("[" + spinnerView + "]")
		} else {
			checkbox = checkboxStyle.Render("[" + spinnerView + "]")
		}
	} else if selected {
		checkbox = checkboxSelectedStyle.Render("[ ]")
	} else {
		checkbox = checkboxStyle.Render("[ ]")
	}

	// Task text (truncate if too long)
	taskWidth := m.getTaskColumnWidth()
	taskText := strings.ToLower(task.Text)
	if len(taskText) > taskWidth-4 { // Reserve space for checkbox
		taskText = taskText[:taskWidth-7] + "..."
	}

	// Pad task text BEFORE styling to ensure proper alignment
	taskDisplayWidth := taskWidth - 4 // Account for checkbox
	taskText = fmt.Sprintf("%-*s", taskDisplayWidth, taskText)

	var styledTaskText string
	if task.Completed {
		styledTaskText = completedTaskStyle.Render(taskText)
	} else {
		if selected {
			styledTaskText = checkboxSelectedStyle.Render(taskText)
		} else {
			styledTaskText = taskTextStyle.Render(taskText)
		}
	}

	// Priority - pad before styling
	priorityText := ""
	switch task.Priority {
	case model.PriorityLow:
		priorityText = "low"
	case model.PriorityMedium:
		priorityText = "medium"
	case model.PriorityHigh:
		priorityText = "high"
	}
	priorityText = fmt.Sprintf("%-10s", priorityText)
	priority := m.formatPriority(task.Priority, priorityText, task.Completed)

	// Due date - pad before styling
	dueDateText := m.getDueDateText(&task)
	dueDateText = fmt.Sprintf("%-12s", dueDateText)
	dueDate := m.formatDueDateStyled(&task, dueDateText, task.Completed)

	// Build the line with properly aligned columns
	line := fmt.Sprintf("%s %s %s  %s  %s",
		cursor,
		checkbox,
		styledTaskText,
		priority,
		dueDate)

	// Apply padding
	if selected {
		return selectedTaskStyle.Render(line)
	}
	return taskStyle.Render(line)
}

// formatPriority formats the priority with appropriate styling (text is already padded)
func (m Model) formatPriority(priority model.Priority, paddedText string, completed bool) string {
	var style lipgloss.Style

	switch priority {
	case model.PriorityLow:
		style = priorityLowStyle
	case model.PriorityMedium:
		style = priorityMediumStyle
	case model.PriorityHigh:
		style = priorityHighStyle
	default:
		style = lipgloss.NewStyle()
	}

	// Apply strikethrough if completed
	if completed {
		style = style.Strikethrough(true).Foreground(textDim)
	}

	return style.Render(paddedText)
}

// getDueDateText returns the plain text for due date
func (m Model) getDueDateText(task *model.Task) string {
	if task.DueDate == nil {
		return "-"
	}

	// Calculate days until due
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	dueDay := time.Date(task.DueDate.Year(), task.DueDate.Month(), task.DueDate.Day(), 0, 0, 0, 0, task.DueDate.Location())

	daysDiff := int(dueDay.Sub(today).Hours() / 24)

	// Return relative label based on days difference
	switch {
	case daysDiff == 0:
		return "today"
	case daysDiff == 1:
		return "tomorrow"
	case daysDiff >= 2 && daysDiff <= 7:
		return "this week"
	case daysDiff >= 8 && daysDiff <= 14:
		return "next week"
	case daysDiff < 0:
		// Overdue - show how many days ago
		if daysDiff == -1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", -daysDiff)
	default:
		// Future date beyond 2 weeks - show actual date
		return task.DueDate.Format("Jan 2")
	}
}

// formatDueDateStyled formats the due date with styling (text is already padded)
func (m Model) formatDueDateStyled(task *model.Task, paddedText string, completed bool) string {
	var style lipgloss.Style

	if completed {
		// If completed, apply strikethrough and dim styling regardless of due date
		style = lipgloss.NewStyle().Foreground(textDim).Strikethrough(true)
	} else if task.DueDate == nil {
		// No due date - use dim style
		style = lipgloss.NewStyle().Foreground(textDim)
	} else {
		// Has a due date but not completed - use plain text (no color)
		style = lipgloss.NewStyle().Foreground(text)
	}

	return style.Render(paddedText)
}

// renderAddTask renders the add task input view
func (m Model) renderAddTask() string {
	var b strings.Builder

	// Content box with form
	formContent := strings.Builder{}

	// Question 1: Task
	questionLabel := checkboxSelectedStyle.Render("Task")
	if m.addField != 0 {
		questionLabel = checkboxStyle.Render("Task")
	}
	formContent.WriteString(" ")
	formContent.WriteString(questionLabel)
	formContent.WriteString(": ")

	if m.addField == 0 {
		// Active input
		formContent.WriteString(m.input)
		formContent.WriteString("█")
	} else {
		// Completed
		formContent.WriteString(taskTextStyle.Render(m.input))
	}
	formContent.WriteString("\n\n")

	// Question 2: Priority
	questionLabel = checkboxSelectedStyle.Render("Priority")
	if m.addField != 1 {
		questionLabel = checkboxStyle.Render("Priority")
	}
	formContent.WriteString(" ")
	formContent.WriteString(questionLabel)
	formContent.WriteString(":")

	if m.addField >= 1 {
		formContent.WriteString("\n")
		if m.addField == 1 {
			// Active selection
			formContent.WriteString(m.renderPriorityList(0))
			formContent.WriteString("\n")
			formContent.WriteString(m.renderPriorityList(1))
			formContent.WriteString("\n")
			formContent.WriteString(m.renderPriorityList(2))
		} else {
			// Completed
			formContent.WriteString("   ")
			formContent.WriteString(m.formatPriorityAnswer(m.addPriority))
		}
	}
	formContent.WriteString("\n\n")

	// Question 3: Due date
	questionLabel = checkboxSelectedStyle.Render("Due date")
	if m.addField != 2 {
		questionLabel = checkboxStyle.Render("Due date")
	}
	formContent.WriteString(" ")
	formContent.WriteString(questionLabel)
	formContent.WriteString(":")

	if m.addField >= 2 {
		formContent.WriteString("\n")
		// Active selection
		formContent.WriteString(m.renderDueList(0))
		formContent.WriteString("\n")
		formContent.WriteString(m.renderDueList(1))
		formContent.WriteString("\n")
		formContent.WriteString(m.renderDueList(2))
		formContent.WriteString("\n")
		formContent.WriteString(m.renderDueList(3))
		formContent.WriteString("\n")
		formContent.WriteString(m.renderDueList(4))
	}

	// Apply box style
	boxStyle := taskListStyle
	if m.width > 0 {
		boxStyle = boxStyle.Width(m.width - 6)
	}
	b.WriteString(boxStyle.Render(formContent.String()))

	// Footer with context-sensitive shortcuts
	b.WriteString("\n")
	var footerItems []footerItem
	switch m.addField {
	case 0:
		footerItems = []footerItem{
			{"↵", "next"},
			{"esc", "cancel"},
		}
	case 1:
		footerItems = []footerItem{
			{"↑↓", "navigate"},
			{"1-3", "select"},
			{"↵", "next"},
			{"esc", "cancel"},
		}
	case 2:
		footerItems = []footerItem{
			{"↑↓", "navigate"},
			{"1-5", "select & add"},
			{"↵", "add"},
			{"esc", "cancel"},
		}
	}

	footer := m.buildFooter(footerItems)
	footerStyleWithWidth := footerStyle
	if m.width > 0 {
		footerStyleWithWidth = footerStyleWithWidth.Width(m.width - 4)
	}
	b.WriteString(footerStyleWithWidth.Render(footer))

	return b.String()
}

// formatPriorityAnswer formats the completed priority answer
func (m Model) formatPriorityAnswer(priority model.Priority) string {
	switch priority {
	case model.PriorityLow:
		return priorityLowStyle.Render(priority.String())
	case model.PriorityMedium:
		return priorityMediumStyle.Render(priority.String())
	case model.PriorityHigh:
		return priorityHighStyle.Render(priority.String())
	}
	return priority.String()
}

// renderPriorityList renders a single priority option
func (m Model) renderPriorityList(index int) string {
	priorities := []model.Priority{model.PriorityLow, model.PriorityMedium, model.PriorityHigh}
	priority := priorities[index]

	cursor := " "
	if m.addCursor == index {
		cursor = cursorStyle.Render("▶")
	}

	var style func(...string) string
	switch priority {
	case model.PriorityLow:
		style = priorityLowStyle.Render
	case model.PriorityMedium:
		style = priorityMediumStyle.Render
	case model.PriorityHigh:
		style = priorityHighStyle.Render
	}

	return fmt.Sprintf("  %s %s", cursor, style(priority.String()))
}

// renderDueList renders a single due date option
func (m Model) renderDueList(index int) string {
	labels := []string{"Today", "Tomorrow", "This week", "Next week", "No due date"}
	label := labels[index]

	cursor := " "
	if m.addCursor == index {
		cursor = cursorStyle.Render("▶")
	}

	// Highlight selected option in purple
	var styledLabel string
	if m.addCursor == index {
		styledLabel = checkboxSelectedStyle.Render(label)
	} else {
		styledLabel = dueDateStyle.Render(label)
	}

	return fmt.Sprintf("  %s %s", cursor, styledLabel)
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
		{"space", "Cycle status (→ in-progress → done → todo)"},
		{"a", "Add new task"},
		{"d", "Delete selected task"},
		{"sp", "Sort by priority"},
		{"sd", "Sort by due date"},
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
