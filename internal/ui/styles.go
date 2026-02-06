package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette using ANSI colors for light/dark theme compatibility
	// These colors automatically adapt to the terminal's color scheme
	PrimaryColor   = lipgloss.AdaptiveColor{Light: "4", Dark: "12"} // Blue (adapts to theme)
	AccentColor    = lipgloss.AdaptiveColor{Light: "6", Dark: "14"} // Cyan (adapts to theme)
	SuccessColor   = lipgloss.AdaptiveColor{Light: "2", Dark: "10"} // Green (adapts to theme)
	WarningColor   = lipgloss.AdaptiveColor{Light: "3", Dark: "11"} // Yellow (adapts to theme)
	ErrorColor     = lipgloss.AdaptiveColor{Light: "1", Dark: "9"}  // Red (adapts to theme)
	MutedColor     = lipgloss.AdaptiveColor{Light: "8", Dark: "8"}  // Gray (adapts to theme)
	TextColor      = lipgloss.AdaptiveColor{Light: "", Dark: ""}    // Default terminal foreground
	HighlightColor = lipgloss.AdaptiveColor{Light: "5", Dark: "13"} // Magenta (adapts to theme)

	// Column styles
	ColumnStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(0, 1)

	ColumnTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(MutedColor).
				Align(lipgloss.Center)

	FocusedColumnTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(AccentColor).
				Align(lipgloss.Center)

	// Task styles
	TaskStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(TextColor)

	SelectedTaskStyle = lipgloss.NewStyle().
				Padding(0, 1).
				Background(HighlightColor).
				Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "0"}).
				Bold(true)

	// Help/status bar
	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Padding(1, 0)

	// Error/Warning styles
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor)
)
