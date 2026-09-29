package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// borderRound is theme-independent, so it lives outside of styles.
var borderRound = lipgloss.RoundedBorder()

// palette holds the colors for the active theme. Each color is picked for the
// terminal's current background via lipgloss.LightDark.
type palette struct {
	// Primary - Deeper royal purple for crisp contrast on light backdrops
	primary color.Color
	// Accent - Deeper indigo/navy blue instead of cyan for high legibility
	accent color.Color

	// Status colors
	success color.Color
	danger  color.Color

	// Priority colors
	priorityHigh   color.Color
	priorityMedium color.Color
	priorityLow    color.Color

	// Text colors
	text      color.Color
	textMuted color.Color
	textDim   color.Color

	// Border colors
	border color.Color
}

// styles holds every style the views render with, resolved for one theme.
// Rebuild it with newStyles whenever the terminal background changes.
type styles struct {
	palette

	// Task list container
	taskListStyle lipgloss.Style

	// Task item styles
	taskStyle              lipgloss.Style
	selectedTaskStyle      lipgloss.Style
	cursorStyle            lipgloss.Style
	checkboxStyle          lipgloss.Style
	checkboxSelectedStyle  lipgloss.Style
	checkboxCompletedStyle lipgloss.Style

	// Priority badge styles
	priorityHighStyle   lipgloss.Style
	priorityMediumStyle lipgloss.Style
	priorityLowStyle    lipgloss.Style

	// Task text styles
	taskTextStyle      lipgloss.Style
	completedTaskStyle lipgloss.Style

	// Due date styles
	dueDateStyle lipgloss.Style

	// Message styles
	errorStyle   lipgloss.Style
	successStyle lipgloss.Style

	// Input dialog styles
	dialogBoxStyle lipgloss.Style
	hintStyle      lipgloss.Style

	// Help styles
	helpHeaderStyle lipgloss.Style
	helpKeyStyle    lipgloss.Style
	helpDescStyle   lipgloss.Style

	// Footer styles
	footerStyle     lipgloss.Style
	footerKeyStyle  lipgloss.Style
	footerDescStyle lipgloss.Style
	footerSepStyle  lipgloss.Style
}

// newPalette resolves the color palette for a light or dark background.
func newPalette(isDark bool) palette {
	lightDark := lipgloss.LightDark(isDark)

	return palette{
		primary: lightDark(
			lipgloss.Color("#5B21B6"), // Deep Violet/Purple (WCAG AA compliant)
			lipgloss.Color("#F0ABFC"), // Subtle pink for dark mode
		),
		accent: lightDark(
			lipgloss.Color("#234ab6"), // Deep Indigo/Blue
			lipgloss.Color("#22D3EE"), // Lighter cyan for dark mode
		),
		success: lightDark(
			lipgloss.Color("#047857"), // Deep Emerald
			lipgloss.Color("#10B981"), // Lighter green for dark mode
		),
		danger: lightDark(
			lipgloss.Color("#B91C1C"), // Deep Red
			lipgloss.Color("#EF4444"), // Lighter red for dark mode
		),
		priorityHigh: lightDark(
			lipgloss.Color("#B91C1C"), // Deep Red
			lipgloss.Color("#F87171"),
		),
		priorityMedium: lightDark(
			lipgloss.Color("#B45309"), // Deep Amber/Brown-Orange for readability
			lipgloss.Color("#FBBF24"),
		),
		priorityLow: lightDark(
			lipgloss.Color("#047857"), // Deep Green
			lipgloss.Color("#34D399"),
		),
		text: lightDark(
			lipgloss.Color("#0F172A"), // Slate Black/Dark Navy for primary text
			lipgloss.Color("#E5E7EB"),
		),
		textMuted: lightDark(
			lipgloss.Color("#334155"), // Slate Gray
			lipgloss.Color("#9CA3AF"),
		),
		textDim: lightDark(
			lipgloss.Color("#475569"), // Darker Slate Gray (high contrast for light mode strikethrough)
			lipgloss.Color("#9CA3AF"), // Light Gray for dark mode
		),
		border: lightDark(
			lipgloss.Color("#94A3B8"), // Slate border
			lipgloss.Color("#6B7280"),
		),
	}
}

// newStyles builds the style set for a light or dark terminal background.
func newStyles(isDark bool) styles {
	p := newPalette(isDark)

	return styles{
		palette: p,

		taskListStyle: lipgloss.NewStyle().
			Border(borderRound).
			BorderForeground(p.border).
			Padding(0, 1).
			MarginBottom(1),

		taskStyle: lipgloss.NewStyle().
			PaddingLeft(1),

		selectedTaskStyle: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1),

		cursorStyle: lipgloss.NewStyle().
			Foreground(p.primary).
			Bold(true),

		checkboxStyle: lipgloss.NewStyle().
			Foreground(p.text),

		checkboxSelectedStyle: lipgloss.NewStyle().
			Foreground(p.primary).
			Bold(true),

		checkboxCompletedStyle: lipgloss.NewStyle().
			Foreground(p.textDim).
			Strikethrough(true),

		priorityHighStyle: lipgloss.NewStyle().
			Foreground(p.priorityHigh),

		priorityMediumStyle: lipgloss.NewStyle().
			Foreground(p.priorityMedium),

		priorityLowStyle: lipgloss.NewStyle().
			Foreground(p.priorityLow),

		taskTextStyle: lipgloss.NewStyle().
			Foreground(p.text),

		completedTaskStyle: lipgloss.NewStyle().
			Foreground(p.textDim).
			Strikethrough(true),

		dueDateStyle: lipgloss.NewStyle().
			Foreground(p.textMuted),

		errorStyle: lipgloss.NewStyle().
			Foreground(p.danger).
			Padding(0, 1).
			MarginTop(1).
			Bold(true),

		successStyle: lipgloss.NewStyle().
			Foreground(p.success).
			Padding(0, 1).
			MarginTop(1).
			Bold(true),

		dialogBoxStyle: lipgloss.NewStyle().
			Border(borderRound).
			BorderForeground(p.primary).
			Padding(1, 2),

		hintStyle: lipgloss.NewStyle().
			Foreground(p.textMuted).
			Italic(true).
			MarginTop(1),

		helpHeaderStyle: lipgloss.NewStyle().
			Foreground(p.primary).
			Bold(true).
			MarginBottom(1),

		helpKeyStyle: lipgloss.NewStyle().
			Foreground(p.accent).
			Bold(true),

		helpDescStyle: lipgloss.NewStyle().
			Foreground(p.text),

		footerStyle: lipgloss.NewStyle().
			Foreground(p.textMuted).
			BorderStyle(lipgloss.Border{Top: "─"}).
			BorderForeground(p.border).
			BorderTop(true).
			PaddingTop(1).
			PaddingLeft(1).
			MarginTop(1).
			MarginLeft(1),

		footerKeyStyle: lipgloss.NewStyle().
			Foreground(p.accent).
			Bold(true),

		footerDescStyle: lipgloss.NewStyle().
			Foreground(p.textMuted).
			Faint(true),

		footerSepStyle: lipgloss.NewStyle().
			Foreground(p.textDim).
			Faint(true),
	}
}
