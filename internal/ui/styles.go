package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Column styles
	ColumnStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	ColumnTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205")).
				Align(lipgloss.Center)

	// Task styles
	TaskStyle = lipgloss.NewStyle().
			Padding(0, 1)

	SelectedTaskStyle = lipgloss.NewStyle().
				Padding(0, 1).
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("230")).
				Bold(true)

	// Help/status bar
	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(1, 0)
)
