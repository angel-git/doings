package ui

import (
	"doings/internal/task"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DetailModel represents the detail view of a task
type DetailModel struct {
	task   *task.Task
	cursor int // Selected line in checklist (0-indexed)
	width  int
	height int
}

// NewDetailModel creates a new detail model
func NewDetailModel(t *task.Task) DetailModel {
	return DetailModel{
		task:   t,
		cursor: 0,
		width:  80,
		height: 24,
	}
}

// Init initializes the detail model
func (m DetailModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the detail model
func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.cursor < len(m.task.Checklist)-1 {
				m.cursor++
			}

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}

	return m, nil
}

// View renders the detail view
func (m DetailModel) View() string {
	var sections []string

	// Header with title and status
	header := m.renderHeader()
	sections = append(sections, header)

	// Separator
	sections = append(sections, m.renderSeparator())

	// Description
	if m.task.Description != "" {
		desc := m.renderDescription()
		sections = append(sections, desc)
		sections = append(sections, m.renderSeparator())
	}

	// Checklist
	if len(m.task.Checklist) > 0 {
		checklist := m.renderChecklist()
		sections = append(sections, checklist)
	}

	// Help bar
	help := m.renderHelp()

	// Combine all sections
	content := strings.Join(sections, "\n")

	return lipgloss.JoinVertical(lipgloss.Left, content, help)
}

// renderHeader renders the task title and status
func (m DetailModel) renderHeader() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		PaddingLeft(1)

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("62")).
		Bold(true).
		PaddingRight(1)

	title := titleStyle.Render("# " + m.task.Title)
	status := statusStyle.Render("[" + m.task.Status + "]")

	// Calculate spacing to right-align status
	titleLen := len("# " + m.task.Title)
	statusLen := len("[" + m.task.Status + "]")
	spacing := m.width - titleLen - statusLen - 4 // -4 for padding

	if spacing < 1 {
		spacing = 1
	}

	return title + strings.Repeat(" ", spacing) + status
}

// renderSeparator renders a horizontal line
func (m DetailModel) renderSeparator() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		PaddingLeft(1)

	separatorWidth := m.width - 2
	if separatorWidth < 1 {
		separatorWidth = 1
	}

	return style.Render(strings.Repeat("─", separatorWidth))
}

// renderDescription renders the task description
func (m DetailModel) renderDescription() string {
	style := lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Foreground(lipgloss.Color("252"))

	return style.Render(m.task.Description)
}

// renderChecklist renders the checklist items
func (m DetailModel) renderChecklist() string {
	var lines []string

	for i, item := range m.task.Checklist {
		lines = append(lines, m.renderChecklistItem(item, i == m.cursor))
	}

	return strings.Join(lines, "\n")
}

// renderChecklistItem renders a single checklist item
func (m DetailModel) renderChecklistItem(item task.CheckItem, selected bool) string {
	// Build the checkbox
	checkbox := "[ ]"
	if item.Checked {
		checkbox = "[x]"
	}

	// Build the line with indentation
	indent := strings.Repeat("    ", item.Indent)
	line := indent + "- " + checkbox + " " + item.Text

	// Apply styling
	style := lipgloss.NewStyle().PaddingLeft(1)

	if selected {
		style = style.
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Bold(true)
	}

	// If checked, make it dimmer
	if item.Checked && !selected {
		style = style.Foreground(lipgloss.Color("240"))
	}

	return style.Render(line)
}

// renderHelp renders the help bar
func (m DetailModel) renderHelp() string {
	help := "j/k: navigate | Space: toggle | Esc: back | s: save"
	return HelpStyle.Render(help)
}
