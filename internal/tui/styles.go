package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Adaptive color palette - works in both light and dark mode
var (
	// Primary colors
	primary = lipgloss.AdaptiveColor{
		Light: "#7D56F4", // Darker purple for light mode
		Dark:  "#A78BFA", // Lighter purple for dark mode
	}

	// Accent colors
	accent = lipgloss.AdaptiveColor{
		Light: "#0891B2", // Darker cyan for light mode
		Dark:  "#22D3EE", // Lighter cyan for dark mode
	}

	// Status colors
	success = lipgloss.AdaptiveColor{
		Light: "#059669", // Darker green for light mode
		Dark:  "#10B981", // Lighter green for dark mode
	}
	warning = lipgloss.AdaptiveColor{
		Light: "#D97706", // Darker amber for light mode
		Dark:  "#F59E0B", // Lighter amber for dark mode
	}
	danger = lipgloss.AdaptiveColor{
		Light: "#DC2626", // Darker red for light mode
		Dark:  "#EF4444", // Lighter red for dark mode
	}

	// Priority colors
	priorityHigh = lipgloss.AdaptiveColor{
		Light: "#DC2626", // Darker red for light mode
		Dark:  "#F87171", // Lighter red for dark mode
	}
	priorityMedium = lipgloss.AdaptiveColor{
		Light: "#D97706", // Darker yellow for light mode
		Dark:  "#FBBF24", // Lighter yellow for dark mode
	}
	priorityLow = lipgloss.AdaptiveColor{
		Light: "#059669", // Darker green for light mode
		Dark:  "#34D399", // Lighter green for dark mode
	}

	// Text colors
	text = lipgloss.AdaptiveColor{
		Light: "#000000", // Black text for light mode
		Dark:  "#E5E7EB", // Light text for dark mode
	}
	textMuted = lipgloss.AdaptiveColor{
		Light: "#374151", // Darker gray for light mode
		Dark:  "#9CA3AF", // Light gray for dark mode
	}
	textDim = lipgloss.AdaptiveColor{
		Light: "#6B7280", // Medium gray for light mode
		Dark:  "#6B7280", // Darker gray for dark mode
	}

	// Border colors
	border = lipgloss.AdaptiveColor{
		Light: "#6B7280", // Dark gray border for light mode
		Dark:  "#6B7280", // Dark gray border for dark mode
	}
)

var (
	// Border styles
	borderNormal = lipgloss.NormalBorder()
	borderRound  = lipgloss.RoundedBorder()

	// Title bar style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primary).
			Padding(0, 1).
			MarginBottom(1)

	// App container
	appStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Stats badges
	statsContainerStyle = lipgloss.NewStyle().
				MarginBottom(1)

	statBadgeStyle = lipgloss.NewStyle().
			Foreground(text).
			Padding(0, 1).
			MarginRight(1)

	statNumberStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	// Task list container
	taskListStyle = lipgloss.NewStyle().
			Border(borderRound).
			BorderForeground(border).
			Padding(0, 1).
			MarginBottom(1)

	// Task item styles
	taskStyle = lipgloss.NewStyle().
			PaddingLeft(1)

	selectedTaskStyle = lipgloss.NewStyle().
				Foreground(primary).
				Bold(true).
				PaddingLeft(1).
				PaddingRight(1)

	cursorStyle = lipgloss.NewStyle().
			Foreground(primary).
			Bold(true)

	checkboxStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	checkboxCompletedStyle = lipgloss.NewStyle().
				Foreground(success).
				Bold(true)

	// Priority badge styles
	priorityHighStyle = lipgloss.NewStyle().
				Foreground(priorityHigh).
				Bold(true)

	priorityMediumStyle = lipgloss.NewStyle().
				Foreground(priorityMedium).
				Bold(true)

	priorityLowStyle = lipgloss.NewStyle().
				Foreground(priorityLow)

	// Task text styles
	taskTextStyle = lipgloss.NewStyle().
			Foreground(text)

	completedTaskStyle = lipgloss.NewStyle().
				Foreground(textDim).
				Strikethrough(true)

	// Due date styles
	overdueStyle = lipgloss.NewStyle().
			Foreground(danger).
			Bold(true)

	dueSoonStyle = lipgloss.NewStyle().
			Foreground(warning)

	dueDateStyle = lipgloss.NewStyle().
			Foreground(textMuted)

	// Empty state
	emptyStateStyle = lipgloss.NewStyle().
			Foreground(textMuted).
			Italic(true).
			Padding(2, 0)

	// Message styles
	errorStyle = lipgloss.NewStyle().
			Foreground(danger).
			Padding(0, 1).
			MarginTop(1).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(success).
			Padding(0, 1).
			MarginTop(1).
			Bold(true)

	// Input dialog styles
	dialogBoxStyle = lipgloss.NewStyle().
			Border(borderRound).
			BorderForeground(primary).
			Padding(1, 2)

	dialogTitleStyle = lipgloss.NewStyle().
				Foreground(primary).
				Bold(true).
				MarginBottom(1)

	promptStyle = lipgloss.NewStyle().
			Foreground(accent).
			MarginBottom(1)

	inputBoxStyle = lipgloss.NewStyle().
			Border(borderNormal).
			BorderForeground(accent).
			Padding(0, 1).
			Foreground(text)

	inputCursorStyle = lipgloss.NewStyle().
				Foreground(accent)

	hintStyle = lipgloss.NewStyle().
			Foreground(textMuted).
			Italic(true).
			MarginTop(1)

	// Help styles
	helpHeaderStyle = lipgloss.NewStyle().
			Foreground(primary).
			Bold(true).
			MarginBottom(1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(text)

	// Footer styles
	footerStyle = lipgloss.NewStyle().
			Foreground(textMuted).
			BorderStyle(lipgloss.Border{Top: "─"}).
			BorderForeground(border).
			BorderTop(true).
			PaddingTop(1).
			MarginTop(1)

	footerKeyStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	footerSepStyle = lipgloss.NewStyle().
			Foreground(textDim)
)
