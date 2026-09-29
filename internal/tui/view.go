package tui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/jackroberts-gh/tuido/v2/internal/model"
)

const (
	priorityColumnWidth = 10
	dueDateColumnWidth  = 12
)

// View renders the current view based on the mode
func (m Model) View() tea.View {
	var content string
	switch m.mode {
	case modeList:
		content = m.renderList()
	case modeAdd:
		content = m.renderAddTask()
	case modeHelp:
		content = m.renderHelp()
	case modeDelete:
		content = m.renderDeleteConfirmation()
	default:
		content = "Unknown mode"
	}

	// AltScreen replaces the v1 tea.WithAltScreen program option
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// renderList renders the main task list view
func (m Model) renderList() string {
	var b strings.Builder

	// Task list with responsive width
	listStyle := m.styles.taskListStyle
	if m.width > 0 {
		listStyle = listStyle.Width(m.width - 6)
	}

	visibleTasks := m.getVisibleTasks()
	taskContent := strings.Builder{}

	// Render table header
	taskContent.WriteString(m.renderTableHeader())
	taskContent.WriteString("\n")
	taskContent.WriteString(m.renderHeaderSeparator())

	// Render tasks as table rows
	if len(visibleTasks) == 0 {
		// Show hint when no tasks
		taskContent.WriteString("\n")
		hintText := "Hit a to add a new task"
		emptyHintStyle := lipgloss.NewStyle().
			Foreground(m.styles.textDim).
			PaddingLeft(1)
		taskContent.WriteString(emptyHintStyle.Render("  " + hintText))
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
		b.WriteString(m.styles.errorStyle.Render(fmt.Sprintf(" ✗ %v ", m.err)))
	}
	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(m.styles.successStyle.Render(fmt.Sprintf(" ✓ %s ", m.message)))
	}

	// Footer with shortcuts - main actions
	b.WriteString("\n")
	footer1 := m.buildFooter([]footerItem{
		{"↑↓", "navigate"},
		{"←→", "status"},
		{"a", "add"},
		{"e", "edit"},
		{"d", "delete"},
		{"sp", "sort priority"},
		{"sd", "sort date"},
		{"t", "filter"},
	})
	footerStyleWithWidth := m.styles.footerStyle
	if m.width > 0 {
		footerStyleWithWidth = footerStyleWithWidth.Width(m.width - 4)
	}
	b.WriteString(footerStyleWithWidth.Render(footer1))

	// Footer second line - help and quit
	b.WriteString("\n")
	footer2 := m.buildFooter([]footerItem{
		{"?", "help"},
		{"q", "quit"},
	})
	footerStyle2 := m.styles.footerStyle.BorderTop(false).PaddingTop(0).MarginTop(0)
	if m.width > 0 {
		footerStyle2 = footerStyle2.Width(m.width - 4)
	}
	b.WriteString(footerStyle2.Render(footer2))

	return b.String()
}

type footerItem struct {
	key  string
	desc string
}

func (m Model) buildFooter(items []footerItem) string {
	parts := make([]string, len(items))
	for i, item := range items {
		parts[i] = fmt.Sprintf("%s %s", m.styles.footerKeyStyle.Render(item.key), m.styles.footerDescStyle.Render(item.desc))
	}
	return strings.Join(parts, m.styles.footerSepStyle.Render(" • "))
}

// renderTableHeader renders the table header
func (m Model) renderTableHeader() string {
	// Calculate column widths
	taskWidth := m.getTaskColumnWidth()

	// Build header with proper spacing
	// Account for: cursor(1) + space(1) = 2 chars before checkbox/Task header
	// Task column includes checkbox(3) + space(1) + text = taskWidth

	// Build the full header string first, then apply styling
	header := fmt.Sprintf("  %-*s  %-*s  %-*s",
		taskWidth, "Task",
		priorityColumnWidth, "Priority",
		dueDateColumnWidth, "Due Date")

	// Use muted text color for the entire header (so Due Date isn't in accent color)
	headerStyle := lipgloss.NewStyle().
		Foreground(m.styles.textMuted).
		Bold(true).
		PaddingLeft(1)

	return headerStyle.Render(header)
}

// renderHeaderSeparator renders a line separator below the header
func (m Model) renderHeaderSeparator() string {
	// Calculate total width for separator line
	taskWidth := m.getTaskColumnWidth()
	totalWidth := 2 + taskWidth + 2 + priorityColumnWidth + 2 + dueDateColumnWidth // spacing included

	// A solid run of dashes has no break points, so the renderer can only wrap
	// it by pushing the whole run onto the next line - which would show up as a
	// stray line under the header. It needs the PaddingLeft(1) column below
	// plus two columns of slack to stay on one line, so cap it at that. This
	// only binds on very narrow terminals, where getTaskColumnWidth clamps the
	// task column to its 20-column minimum.
	if maxDashes := m.innerWidth() - 3; totalWidth > maxDashes {
		totalWidth = maxDashes
	}
	if totalWidth < 1 {
		return ""
	}

	separator := strings.Repeat("─", totalWidth)

	separatorStyle := lipgloss.NewStyle().
		Foreground(m.styles.border).
		PaddingLeft(1)

	return separatorStyle.Render(separator)
}

// taskColumnReserved is the space a row needs around the task column.
//
// The task list sits in a bordered box of Width(m.width-6) with Padding(0, 1),
// so its usable text area is innerWidth() = m.width-8. Two shapes bound how
// wide the task column can be:
//
//	selected row: padding(1) + cursor(1) + space(1) + checkbox(3) + space(1) +
//	  text(taskWidth-4) + gap(2) + priority(10) + gap(2) + due(12) +
//	  padding(1) = taskWidth + 30  <=  m.width-8  =>  taskWidth <= m.width-38
//
//	separator: an unbreakable run of dashes, which needs two more columns of
//	  slack than breakable text (see renderHeaderSeparator), so it costs one
//	  column more than the row above  =>  taskWidth <= m.width-39
//
// The separator is the tighter of the two, so reserve 39.
const taskColumnReserved = 39

// getTaskColumnWidth calculates the width for the task column based on terminal width
func (m Model) getTaskColumnWidth() int {
	if m.width > 0 {
		taskWidth := m.width - taskColumnReserved
		if taskWidth < 20 {
			taskWidth = 20 // Minimum task column width
		}
		// No maximum limit - let it grow with terminal width
		return taskWidth
	}
	return 40 // Default width
}

// innerWidth returns the usable text width inside the task list box, which is
// Width(m.width-6) less its horizontal padding of 2.
func (m Model) innerWidth() int {
	return m.width - 8
}

// renderTask renders a single task as a table row
func (m Model) renderTask(task model.Task, selected bool) string {
	// Cursor indicator
	cursor := " "
	if selected {
		cursor = m.styles.cursorStyle.Render("▶")
	}

	// Checkbox
	var checkbox string
	if task.Completed {
		checkbox = m.styles.checkboxCompletedStyle.Render("[x]")
	} else if task.InProgress {
		// Show spinner for in-progress tasks
		spinnerView := m.spinner.View()
		if selected {
			checkbox = m.styles.checkboxSelectedStyle.Render("[" + spinnerView + "]")
		} else {
			checkbox = m.styles.checkboxStyle.Render("[" + spinnerView + "]")
		}
	} else if selected {
		checkbox = m.styles.checkboxSelectedStyle.Render("[ ]")
	} else {
		checkbox = m.styles.checkboxStyle.Render("[ ]")
	}

	// Task text (truncate if too long)
	taskWidth := m.getTaskColumnWidth()
	taskText := strings.ToLower(task.Text)

	// Add in-progress indicator if applicable
	var inProgressIndicator string
	if task.InProgress && !task.Completed {
		inProgressIndicator = " (in-progress)"
	}

	// Calculate available space for task text
	taskDisplayWidth := taskWidth - 4 // Account for checkbox
	availableSpace := taskDisplayWidth - len(inProgressIndicator)

	// Truncate task text if needed. On very narrow terminals availableSpace can
	// be tiny (or negative), so only append the ellipsis when there is room for
	// it and never slice past the start of the string.
	if len(taskText) > availableSpace {
		switch {
		case availableSpace <= 0:
			taskText = ""
		case availableSpace <= 3:
			taskText = taskText[:availableSpace]
		default:
			taskText = taskText[:availableSpace-3] + "..."
		}
	}

	var styledTaskText string
	if task.Completed {
		// Pad the combined text
		combined := fmt.Sprintf("%-*s", taskDisplayWidth, taskText+inProgressIndicator)
		styledTaskText = m.styles.completedTaskStyle.Render(combined)
	} else {
		var baseStyle lipgloss.Style
		if selected {
			baseStyle = m.styles.checkboxSelectedStyle
		} else {
			baseStyle = m.styles.taskTextStyle
		}

		if inProgressIndicator != "" {
			// Render task text, then indicator separately (not bold, italic, muted)
			indicatorStyle := lipgloss.NewStyle().Italic(true).Foreground(m.styles.textMuted)
			// Add padding after the indicator
			totalLen := len(taskText) + len(inProgressIndicator)
			padding := strings.Repeat(" ", taskDisplayWidth-totalLen)
			styledTaskText = baseStyle.Render(taskText) + indicatorStyle.Render(inProgressIndicator) + padding
		} else {
			paddedText := fmt.Sprintf("%-*s", taskDisplayWidth, taskText)
			styledTaskText = baseStyle.Render(paddedText)
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
	priorityText = fmt.Sprintf("%-*s", priorityColumnWidth, priorityText)
	priority := m.formatPriority(task.Priority, priorityText, task.Completed)

	// Due date - pad before styling
	dueDateText := m.getDueDateText(&task)
	dueDateText = fmt.Sprintf("%-*s", dueDateColumnWidth, dueDateText)
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
		return m.styles.selectedTaskStyle.Render(line)
	}
	return m.styles.taskStyle.Render(line)
}

// formatPriority formats the priority with appropriate styling (text is already padded)
func (m Model) formatPriority(priority model.Priority, paddedText string, completed bool) string {
	var style lipgloss.Style

	switch priority {
	case model.PriorityLow:
		style = m.styles.priorityLowStyle
	case model.PriorityMedium:
		style = m.styles.priorityMediumStyle
	case model.PriorityHigh:
		style = m.styles.priorityHighStyle
	default:
		style = lipgloss.NewStyle()
	}

	// Apply strikethrough if completed
	if completed {
		style = style.Strikethrough(true).Foreground(m.styles.textDim)
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
		style = lipgloss.NewStyle().Foreground(m.styles.textDim).Strikethrough(true)
	} else if task.DueDate == nil {
		// No due date - use dim style
		style = lipgloss.NewStyle().Foreground(m.styles.textDim)
	} else {
		// Has a due date but not completed - use plain text (no color)
		style = lipgloss.NewStyle().Foreground(m.styles.text)
	}

	return style.Render(paddedText)
}

// renderAddTask renders the add task input view
func (m Model) renderAddTask() string {
	var b strings.Builder

	// Content box with form
	formContent := strings.Builder{}

	// Question 1: Task
	questionLabel := m.styles.checkboxSelectedStyle.Render("Task")
	if m.addField != 0 {
		questionLabel = m.styles.checkboxStyle.Render("Task")
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
		formContent.WriteString(m.styles.taskTextStyle.Render(m.input))
	}
	formContent.WriteString("\n\n")

	// Question 2: Priority
	questionLabel = m.styles.checkboxSelectedStyle.Render("Priority")
	if m.addField != 1 {
		questionLabel = m.styles.checkboxStyle.Render("Priority")
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
	questionLabel = m.styles.checkboxSelectedStyle.Render("Due date")
	if m.addField != 2 {
		questionLabel = m.styles.checkboxStyle.Render("Due date")
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
	boxStyle := m.styles.taskListStyle
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
	footerStyleWithWidth := m.styles.footerStyle
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
		return m.styles.priorityLowStyle.Render(priority.String())
	case model.PriorityMedium:
		return m.styles.priorityMediumStyle.Render(priority.String())
	case model.PriorityHigh:
		return m.styles.priorityHighStyle.Render(priority.String())
	}
	return priority.String()
}

// renderPriorityList renders a single priority option
func (m Model) renderPriorityList(index int) string {
	priorities := []model.Priority{model.PriorityLow, model.PriorityMedium, model.PriorityHigh}
	priority := priorities[index]

	cursor := " "
	if m.addCursor == index {
		cursor = m.styles.cursorStyle.Render("▶")
	}

	var style func(...string) string
	switch priority {
	case model.PriorityLow:
		style = m.styles.priorityLowStyle.Render
	case model.PriorityMedium:
		style = m.styles.priorityMediumStyle.Render
	case model.PriorityHigh:
		style = m.styles.priorityHighStyle.Render
	}

	return fmt.Sprintf("  %s %s", cursor, style(priority.String()))
}

// renderDueList renders a single due date option
func (m Model) renderDueList(index int) string {
	labels := []string{"Today", "Tomorrow", "This week", "Next week", "No due date"}
	label := labels[index]

	cursor := " "
	if m.addCursor == index {
		cursor = m.styles.cursorStyle.Render("▶")
	}

	// Highlight selected option in purple
	var styledLabel string
	if m.addCursor == index {
		styledLabel = m.styles.checkboxSelectedStyle.Render(label)
	} else {
		styledLabel = m.styles.dueDateStyle.Render(label)
	}

	return fmt.Sprintf("  %s %s", cursor, styledLabel)
}

// renderHelp renders the help screen
func (m Model) renderHelp() string {
	var help strings.Builder

	help.WriteString(m.styles.helpHeaderStyle.Render("⌨  Keyboard Shortcuts"))
	help.WriteString("\n\n")

	shortcuts := []struct {
		key  string
		desc string
	}{
		{"↑↓ / kj", "Navigate tasks"},
		{"→ / l", "Advance status (todo to in-progress to done)"},
		{"← / h", "Undo status (done to in-progress to todo)"},
		{"a", "Add new task"},
		{"e", "Edit selected task"},
		{"d", "Delete selected task"},
		{"sp", "Sort by priority"},
		{"sd", "Sort by due date"},
		{"t", "Toggle completed tasks visibility"},
		{"?", "Show this help"},
		{"q / Ctrl+C", "Quit application"},
	}

	for _, shortcut := range shortcuts {
		// Pad the key before styling to ensure alignment
		paddedKey := fmt.Sprintf("%-12s", shortcut.key)
		line := fmt.Sprintf("%s  %s",
			m.styles.helpKeyStyle.Render(paddedKey),
			m.styles.helpDescStyle.Render(shortcut.desc))
		help.WriteString(line)
		help.WriteString("\n")
	}

	help.WriteString("\n")
	help.WriteString(m.styles.hintStyle.Render("Press any key to return"))

	boxStyle := m.styles.dialogBoxStyle
	if m.width > 0 {
		boxStyle = boxStyle.MaxWidth(m.width - 4)
	}
	return boxStyle.Render(help.String())
}

// renderDeleteConfirmation renders the delete confirmation modal
func (m Model) renderDeleteConfirmation() string {
	var dialog strings.Builder

	// Verify task exists
	task := m.taskList.GetByID(m.deleteTaskID)
	if task == nil {
		// Task not found, return to list mode
		return m.renderList()
	}

	// Question
	dialog.WriteString("Are you sure you want to delete?")
	dialog.WriteString("\n\n")

	// Options
	yesOption := fmt.Sprintf("%s %s", m.styles.footerKeyStyle.Render("Y"), m.styles.footerDescStyle.Render("Yes"))
	noOption := fmt.Sprintf("%s %s", m.styles.footerKeyStyle.Render("N"), m.styles.footerDescStyle.Render("No"))
	options := fmt.Sprintf("%s  %s", yesOption, noOption)
	dialog.WriteString(options)

	// Apply dialog box style
	boxStyle := m.styles.dialogBoxStyle.Padding(1, 2)
	if m.width > 0 {
		boxStyle = boxStyle.MaxWidth(40)
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		boxStyle.Render(dialog.String()),
	)
}
