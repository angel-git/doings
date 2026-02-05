package ui

import (
	"doings/internal/task"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DetailMode represents the current mode in detail view
type DetailMode int

const (
	DetailModeNormal DetailMode = iota
	DetailModeInput
	DetailModeConfirm
)

// DetailModel represents the detail view of a task
type DetailModel struct {
	task        *task.Task
	cursor      int // Selected line in checklist (0-indexed)
	width       int
	height      int
	modified    bool            // Has unsaved changes
	mode        DetailMode      // Current interaction mode
	textInput   textinput.Model // Text input for new items
	message     string          // Temporary status message
	confirmMsg  string          // Confirmation prompt
	confirmYes  func()          // Action on 'y'
	confirmNo   func()          // Action on 'n'
	insertBelow bool            // For 'o' vs 'O'
}

// NewDetailModel creates a new detail model
func NewDetailModel(t *task.Task) DetailModel {
	ti := textinput.New()
	ti.Placeholder = "Item text..."
	ti.CharLimit = 200
	ti.Width = 50

	return DetailModel{
		task:      t,
		cursor:    0,
		width:     80,
		height:    24,
		modified:  false,
		mode:      DetailModeNormal,
		textInput: ti,
	}
}

// Init initializes the detail model
func (m DetailModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the detail model
func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd) {
	// Handle save result message
	if saveMsg, ok := msg.(saveResultMsg); ok {
		m.handleSaveResult(saveMsg)
		return m, nil
	}

	switch m.mode {
	case DetailModeInput:
		return m.updateInput(msg)
	case DetailModeConfirm:
		return m.updateConfirm(msg)
	default:
		return m.updateNormal(msg)
	}
}

// updateNormal handles updates in normal mode
func (m DetailModel) updateNormal(msg tea.Msg) (DetailModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if len(m.task.Checklist) > 0 && m.cursor < len(m.task.Checklist)-1 {
				m.cursor++
			}

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case " ":
			// Toggle checkbox
			m.toggleCheckbox()

		case "o":
			// Add item below
			m.insertBelow = true
			m.startInput()
			return m, textinput.Blink

		case "O":
			// Add item above
			m.insertBelow = false
			m.startInput()
			return m, textinput.Blink

		case "x":
			// Delete item
			m.deleteItem()

		case "s":
			// Save changes
			return m, m.saveTask()
		}
	}

	return m, nil
}

// updateInput handles input mode
func (m DetailModel) updateInput(msg tea.Msg) (DetailModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Add the item
			text := m.textInput.Value()
			if text != "" {
				m.addItem(text)
			}
			m.mode = DetailModeNormal
			m.textInput.Blur()
			m.textInput.Reset()
			return m, nil

		case tea.KeyEsc:
			// Cancel
			m.mode = DetailModeNormal
			m.textInput.Blur()
			m.textInput.Reset()
			return m, nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// updateConfirm handles confirm mode
func (m DetailModel) updateConfirm(msg tea.Msg) (DetailModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			// Confirm yes
			m.mode = DetailModeNormal
			if m.confirmYes != nil {
				m.confirmYes()
				m.confirmYes = nil
				m.confirmNo = nil
			}
			return m, nil

		case "n", "N":
			// Confirm no
			m.mode = DetailModeNormal
			if m.confirmNo != nil {
				m.confirmNo()
				m.confirmYes = nil
				m.confirmNo = nil
			}
			return m, nil

		case "c", "C", "esc":
			// Cancel
			m.mode = DetailModeNormal
			m.confirmYes = nil
			m.confirmNo = nil
			return m, nil
		}
	}

	return m, nil
}

// toggleCheckbox toggles the checkbox of the current item
func (m *DetailModel) toggleCheckbox() {
	if len(m.task.Checklist) == 0 || m.cursor >= len(m.task.Checklist) {
		return
	}

	m.task.Checklist[m.cursor].Checked = !m.task.Checklist[m.cursor].Checked
	m.modified = true
	m.message = ""
}

// startInput starts input mode for adding a new item
func (m *DetailModel) startInput() {
	m.mode = DetailModeInput
	m.textInput.Reset()
	m.textInput.Focus()
	m.message = ""
}

// addItem adds a new checklist item
func (m *DetailModel) addItem(text string) {
	// Determine indent level
	indent := 0
	if len(m.task.Checklist) > 0 && m.cursor < len(m.task.Checklist) {
		indent = m.task.Checklist[m.cursor].Indent
	}

	newItem := task.CheckItem{
		Text:    text,
		Checked: false,
		Indent:  indent,
	}

	// Insert at appropriate position
	var insertPos int
	if m.insertBelow {
		insertPos = m.cursor + 1
	} else {
		insertPos = m.cursor
	}

	// Handle empty checklist
	if len(m.task.Checklist) == 0 {
		m.task.Checklist = []task.CheckItem{newItem}
		m.cursor = 0
	} else {
		// Insert the item
		m.task.Checklist = append(m.task.Checklist[:insertPos],
			append([]task.CheckItem{newItem}, m.task.Checklist[insertPos:]...)...)
		m.cursor = insertPos
	}

	m.modified = true
	m.message = ""
}

// deleteItem deletes the current checklist item
func (m *DetailModel) deleteItem() {
	if len(m.task.Checklist) == 0 || m.cursor >= len(m.task.Checklist) {
		return
	}

	// Remove the item
	m.task.Checklist = append(m.task.Checklist[:m.cursor], m.task.Checklist[m.cursor+1:]...)

	// Adjust cursor
	if m.cursor >= len(m.task.Checklist) && m.cursor > 0 {
		m.cursor--
	}

	m.modified = true
	m.message = ""
}

// saveTask saves the task to disk
func (m *DetailModel) saveTask() tea.Cmd {
	return func() tea.Msg {
		if err := task.SaveTask(m.task); err != nil {
			return saveResultMsg{err: err}
		}
		return saveResultMsg{err: nil}
	}
}

// saveResultMsg is returned after saving
type saveResultMsg struct {
	err error
}

// HandleSaveResult handles the save result message
func (m *DetailModel) handleSaveResult(msg saveResultMsg) {
	if msg.err != nil {
		m.message = "Error saving: " + msg.err.Error()
	} else {
		m.modified = false
		m.message = "Saved!"
	}
}

// HasUnsavedChanges returns whether there are unsaved changes
func (m DetailModel) HasUnsavedChanges() bool {
	return m.modified
}

// IsNormalMode returns whether the detail view is in normal mode
func (m DetailModel) IsNormalMode() bool {
	return m.mode == DetailModeNormal
}

// SetConfirmUnsaved sets up the unsaved changes confirmation
func (m *DetailModel) SetConfirmUnsaved(onSave func(), onDiscard func()) {
	m.mode = DetailModeConfirm
	m.confirmMsg = "Unsaved changes. Save? (y/n/c)"
	m.confirmYes = func() {
		if err := task.SaveTask(m.task); err == nil {
			m.modified = false
		}
		if onSave != nil {
			onSave()
		}
	}
	m.confirmNo = func() {
		if onDiscard != nil {
			onDiscard()
		}
	}
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

	modifiedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("208")).
		Bold(true)

	title := titleStyle.Render("# " + m.task.Title)
	status := statusStyle.Render("[" + m.task.Status + "]")

	// Add modified indicator
	modifiedIndicator := ""
	if m.modified {
		modifiedIndicator = modifiedStyle.Render(" [*]")
	}

	// Calculate spacing to right-align status
	titleLen := len("# " + m.task.Title)
	statusLen := len("[" + m.task.Status + "]")
	modLen := 0
	if m.modified {
		modLen = 4 // " [*]"
	}
	spacing := m.width - titleLen - statusLen - modLen - 4 // -4 for padding

	if spacing < 1 {
		spacing = 1
	}

	return title + strings.Repeat(" ", spacing) + status + modifiedIndicator
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
	var help string

	switch m.mode {
	case DetailModeInput:
		help = "Enter new item text (Enter to confirm, Esc to cancel):\n" + m.textInput.View()
	case DetailModeConfirm:
		help = m.confirmMsg
	default:
		help = "j/k: navigate | Space: toggle | o/O: add | x: delete | s: save | Esc: back"
		if m.message != "" {
			help = m.message + " | " + help
		}
	}

	return HelpStyle.Render(help)
}
