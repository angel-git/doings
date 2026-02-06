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
	DetailModeEditDescription
	DetailModeConfirm
)

// DetailModel represents the detail view of a task
type DetailModel struct {
	task             *task.Task
	cursor           int // Selected line in checklist (0-indexed)
	width            int
	height           int
	modified         bool            // Has unsaved changes
	mode             DetailMode      // Current interaction mode
	textInput        textinput.Model // Text input for new items
	descriptionInput textinput.Model // Text input for description
	message          string          // Temporary status message
	confirmMsg       string          // Confirmation prompt
	insertBelow      bool            // For 'o' vs 'O'
	lastKeyPress     string          // Track last key for multi-key sequences like "gg"
}

// NewDetailModel creates a new detail model
func NewDetailModel(t *task.Task) DetailModel {
	ti := textinput.New()
	ti.Placeholder = "Item text..."
	ti.CharLimit = 200
	ti.Width = 50

	descInput := textinput.New()
	descInput.Placeholder = "Task description..."
	descInput.CharLimit = 1000
	descInput.Width = 70

	return DetailModel{
		task:             t,
		cursor:           0,
		width:            80,
		height:           24,
		modified:         false,
		mode:             DetailModeNormal,
		textInput:        ti,
		descriptionInput: descInput,
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
	case DetailModeEditDescription:
		return m.updateEditDescription(msg)
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
		key := msg.String()

		// Handle multi-key sequences
		switch key {
		case "g":
			if m.lastKeyPress == "g" {
				// gg - Jump to first item
				m.cursor = 0
				m.lastKeyPress = ""
				return m, nil
			}
			m.lastKeyPress = "g"
			return m, nil
		}

		// Reset multi-key sequence if a different key is pressed
		if key != m.lastKeyPress {
			m.lastKeyPress = ""
		}

		switch key {
		case "j", "down":
			if len(m.task.Checklist) > 0 && m.cursor < len(m.task.Checklist)-1 {
				m.cursor++
			}

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "G":
			// Jump to last item
			if len(m.task.Checklist) > 0 {
				m.cursor = len(m.task.Checklist) - 1
			}

		case "J":
			// Move item down
			m.moveItemDown()

		case "K":
			// Move item up
			m.moveItemUp()

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

		case "e":
			// Edit description
			m.startEditDescription()
			return m, textinput.Blink

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

// updateEditDescription handles description editing mode
func (m DetailModel) updateEditDescription(msg tea.Msg) (DetailModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Save the description
			m.task.Description = m.descriptionInput.Value()
			m.modified = true
			m.mode = DetailModeNormal
			m.descriptionInput.Blur()
			return m, nil

		case tea.KeyEsc:
			// Cancel
			m.mode = DetailModeNormal
			m.descriptionInput.Blur()
			return m, nil
		}
	}

	m.descriptionInput, cmd = m.descriptionInput.Update(msg)
	return m, cmd
}

// updateConfirm handles confirm mode
func (m DetailModel) updateConfirm(msg tea.Msg) (DetailModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			// Confirm yes - save and exit
			m.mode = DetailModeNormal
			if err := task.SaveTask(m.task); err == nil {
				m.modified = false
			}
			return m, func() tea.Msg {
				return ConfirmResultMsg{Action: "save"}
			}

		case "n", "N":
			// Confirm no - discard and exit
			m.mode = DetailModeNormal
			return m, func() tea.Msg {
				return ConfirmResultMsg{Action: "discard"}
			}

		case "c", "C", "esc":
			// Cancel - stay in detail view
			m.mode = DetailModeNormal
			return m, func() tea.Msg {
				return ConfirmResultMsg{Action: "cancel"}
			}
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

// startEditDescription starts description editing mode
func (m *DetailModel) startEditDescription() {
	m.mode = DetailModeEditDescription
	m.descriptionInput.SetValue(m.task.Description)
	m.descriptionInput.Focus()
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

// moveItemDown moves the current checklist item down
func (m *DetailModel) moveItemDown() {
	if len(m.task.Checklist) == 0 || m.cursor >= len(m.task.Checklist)-1 {
		return
	}

	// Swap current item with the one below
	m.task.Checklist[m.cursor], m.task.Checklist[m.cursor+1] = m.task.Checklist[m.cursor+1], m.task.Checklist[m.cursor]
	m.cursor++
	m.modified = true
	m.message = ""
}

// moveItemUp moves the current checklist item up
func (m *DetailModel) moveItemUp() {
	if len(m.task.Checklist) == 0 || m.cursor == 0 {
		return
	}

	// Swap current item with the one above
	m.task.Checklist[m.cursor], m.task.Checklist[m.cursor-1] = m.task.Checklist[m.cursor-1], m.task.Checklist[m.cursor]
	m.cursor--
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

// ConfirmResultMsg is returned after confirmation dialog
type ConfirmResultMsg struct {
	Action string // "save", "discard", or "cancel"
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

// IsConfirmMode returns whether the detail view is in confirm mode
func (m DetailModel) IsConfirmMode() bool {
	return m.mode == DetailModeConfirm
}

// ShowUnsavedConfirmation shows the unsaved changes confirmation
func (m *DetailModel) ShowUnsavedConfirmation() {
	m.mode = DetailModeConfirm
	m.confirmMsg = "Unsaved changes. Save? (y/n/c)"
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
		Foreground(HighlightColor).
		PaddingLeft(1)

	statusStyle := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		PaddingRight(1)

	modifiedStyle := lipgloss.NewStyle().
		Foreground(WarningColor).
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
		Foreground(MutedColor).
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
		Foreground(TextColor)

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
			Background(HighlightColor).
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "0"}).
			Bold(true)
	}

	// If checked, make it dimmer
	if item.Checked && !selected {
		style = style.Foreground(MutedColor)
	}

	return style.Render(line)
}

// renderHelp renders the help bar
func (m DetailModel) renderHelp() string {
	var help string

	switch m.mode {
	case DetailModeInput:
		help = "Enter new item text (Enter to confirm, Esc to cancel):\n" + m.textInput.View()
	case DetailModeEditDescription:
		help = "Edit description (Enter to save, Esc to cancel):\n" + m.descriptionInput.View()
	case DetailModeConfirm:
		help = m.confirmMsg
	default:
		help = "j/k: navigate | gg/G: first/last | Space: toggle | J/K: move up/down | o/O: add | x: delete | e: edit description | s: save | Esc: back"
		if m.message != "" {
			help = m.message + " | " + help
		}
	}

	return HelpStyle.Render(help)
}
