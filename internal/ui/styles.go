package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette
	PrimaryColor   = lipgloss.Color("62")  // Teal/purple
	AccentColor    = lipgloss.Color("205") // Pink
	SuccessColor   = lipgloss.Color("42")  // Green
	WarningColor   = lipgloss.Color("208") // Orange
	ErrorColor     = lipgloss.Color("196") // Red
	MutedColor     = lipgloss.Color("240") // Gray
	TextColor      = lipgloss.Color("252") // Light gray
	HighlightColor = lipgloss.Color("230") // Yellow

	// Column styles
	ColumnStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(0, 1)

	ColumnTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(AccentColor).
				Align(lipgloss.Center)

	FocusedColumnTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(HighlightColor).
				Align(lipgloss.Center)

	// Task styles
	TaskStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(TextColor)

	SelectedTaskStyle = lipgloss.NewStyle().
				Padding(0, 1).
				Background(PrimaryColor).
				Foreground(HighlightColor).
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
