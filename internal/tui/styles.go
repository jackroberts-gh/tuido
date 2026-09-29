package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Adaptive color palette - works in both light and dark mode
var (
	// Primary colors - Deeper royal purple for crisp contrast on light backdrops
	primary = lipgloss.AdaptiveColor{
		Light: "#5B21B6", // Deep Violet/Purple (WCAG AA compliant)
		Dark:  "#F0ABFC", // Subtle pink for dark mode
	}

	// Accent colors - Deeper indigo/navy blue instead of cyan for high legibility
	accent = lipgloss.AdaptiveColor{
		Light: "#234ab6", // Deep Indigo/Blue
		Dark:  "#22D3EE", // Lighter cyan for dark mode
	}

	// Status colors
	success = lipgloss.AdaptiveColor{
		Light: "#047857", // Deep Emerald
		Dark:  "#10B981", // Lighter green for dark mode
	}
	danger = lipgloss.AdaptiveColor{
		Light: "#B91C1C", // Deep Red
		Dark:  "#EF4444", // Lighter red for dark mode
	}

	// Priority colors
	priorityHigh = lipgloss.AdaptiveColor{
		Light: "#B91C1C", // Deep Red
		Dark:  "#F87171",
	}
	priorityMedium = lipgloss.AdaptiveColor{
		Light: "#B45309", // Deep Amber/Brown-Orange for readability
		Dark:  "#FBBF24",
	}
	priorityLow = lipgloss.AdaptiveColor{
		Light: "#047857", // Deep Green
		Dark:  "#34D399",
	}

	// Text colors
	text = lipgloss.AdaptiveColor{
		Light: "#0F172A", // Slate Black/Dark Navy for primary text
		Dark:  "#E5E7EB",
	}
	textMuted = lipgloss.AdaptiveColor{
		Light: "#334155", // Slate Gray
		Dark:  "#9CA3AF",
	}
	textDim = lipgloss.AdaptiveColor{
		Light: "#475569", // Darker Slate Gray (high contrast for light mode strikethrough)
		Dark:  "#9CA3AF", // Light Gray for dark mode
	}

	// Border colors
	border = lipgloss.AdaptiveColor{
		Light: "#94A3B8", // Slate border
		Dark:  "#6B7280",
	}
)

var (
	// Border styles
	borderRound = lipgloss.RoundedBorder()

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
				PaddingLeft(1).
				PaddingRight(1)

	cursorStyle = lipgloss.NewStyle().
			Foreground(primary).
			Bold(true)

	checkboxStyle = lipgloss.NewStyle().
			Foreground(text)

	checkboxSelectedStyle = lipgloss.NewStyle().
				Foreground(primary).
				Bold(true)

	checkboxCompletedStyle = lipgloss.NewStyle().
				Foreground(textDim).
				Strikethrough(true)

	// Priority badge styles
	priorityHighStyle = lipgloss.NewStyle().
				Foreground(priorityHigh)

	priorityMediumStyle = lipgloss.NewStyle().
				Foreground(priorityMedium)

	priorityLowStyle = lipgloss.NewStyle().
				Foreground(priorityLow)

	// Task text styles
	taskTextStyle = lipgloss.NewStyle().
			Foreground(text)

	completedTaskStyle = lipgloss.NewStyle().
				Foreground(textDim).
				Strikethrough(true)

	// Due date styles
	dueDateStyle = lipgloss.NewStyle().
			Foreground(textMuted)

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
			PaddingLeft(1).
			MarginTop(1).
			MarginLeft(1)

	footerKeyStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true)

	footerDescStyle = lipgloss.NewStyle().
			Foreground(textMuted).
			Faint(true)

	footerSepStyle = lipgloss.NewStyle().
			Foreground(textDim).
			Faint(true)
)
