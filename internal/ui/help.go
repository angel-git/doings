package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// RenderHelpScreen renders the help screen overlay
func RenderHelpScreen(width, height int) string {
	helpContent := `
Keyboard Shortcuts

Navigation
  h/←              Move left
  j/↓              Move down
  k/↑              Move up
  l/→              Move right
  gg               Jump to first task/item
  G                Jump to last task/item
  Enter            Open task detail

Board Actions
  i/n              Create new task
  dd               Delete task
  H                Move task left
  L                Move task right
  q                Quit application

Detail View
  j/k              Navigate checklist
  gg/G             Jump to first/last item
  Space            Toggle checkbox
  o                Add item below
  O                Add item above
  x                Delete item
  J                Move item down
  K                Move item up
  e                Edit task description
  s                Save changes
  Esc              Return to board

Help
  ?                Show/hide this help
`

	helpStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(AccentColor).
		Padding(1, 2).
		Width(50)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(AccentColor).
		Align(lipgloss.Center)

	title := titleStyle.Render("HELP")
	content := helpStyle.Render(helpContent)

	// Center the help box
	help := lipgloss.JoinVertical(lipgloss.Center, title, content)

	// Add background overlay
	overlayStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center)

	return overlayStyle.Render(help)
}
