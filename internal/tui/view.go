package tui

import (
	"fmt"
	"strings"

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

	// Build the line
	line := fmt.Sprintf("%s %s %s", cursor, checkbox, taskText)

	// Apply padding
	if selected {
		return selectedTaskStyle.Render(line)
	}
	return taskStyle.Render(line)
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
			{"1-4", "select & add"},
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
	labels := []string{"Today", "Tomorrow", "This week", "Next week"}
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
