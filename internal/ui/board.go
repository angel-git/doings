package ui

import (
	"fmt"
	"strings"

	"doings/internal/config"
	"doings/internal/task"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Mode represents the current interaction mode
type Mode int

const (
	ModeNormal Mode = iota
	ModeInput
	ModeConfirm
)

// Cursor tracks the current position on the board
type Cursor struct {
	Column int // Which column (0-indexed)
	Row    int // Which task in column (0-indexed)
}

// BoardModel represents the main board view
type BoardModel struct {
	columns       []string                // Column names from config
	tasks         map[string][]*task.Task // Tasks grouped by status
	cursor        Cursor                  // Current position
	width         int                     // Terminal width
	height        int                     // Terminal height
	mode          Mode                    // Current mode
	textInput     textinput.Model         // Text input for new tasks
	confirmMsg    string                  // Confirmation message
	confirmAction func() tea.Cmd          // Action to perform on confirmation
}

// NewBoardModel creates a new board model
func NewBoardModel(columns []string, tasks []*task.Task) BoardModel {
	// Group tasks by status
	tasksByStatus := make(map[string][]*task.Task)
	for _, t := range tasks {
		tasksByStatus[t.Status] = append(tasksByStatus[t.Status], t)
	}

	// Initialize text input
	ti := textinput.New()
	ti.Placeholder = "Task title..."
	ti.CharLimit = 100
	ti.Width = 50

	return BoardModel{
		columns:   columns,
		tasks:     tasksByStatus,
		cursor:    Cursor{Column: 0, Row: 0},
		width:     80,
		height:    24,
		mode:      ModeNormal,
		textInput: ti,
	}
}

// Init initializes the board model
func (m BoardModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m BoardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle reload message
	if reload, ok := msg.(reloadMsg); ok {
		m.handleReload(reload)
		return m, nil
	}

	switch m.mode {
	case ModeInput:
		return m.updateInput(msg)
	case ModeConfirm:
		return m.updateConfirm(msg)
	default:
		return m.updateNormal(msg)
	}
}

// updateNormal handles input in normal mode
func (m BoardModel) updateNormal(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "h", "left":
			if m.cursor.Column > 0 {
				m.cursor.Column--
				m.adjustCursorRow()
			}

		case "l", "right":
			if m.cursor.Column < len(m.columns)-1 {
				m.cursor.Column++
				m.adjustCursorRow()
			}

		case "j", "down":
			columnName := m.columns[m.cursor.Column]
			tasksInColumn := m.tasks[columnName]
			if m.cursor.Row < len(tasksInColumn)-1 {
				m.cursor.Row++
			}

		case "k", "up":
			if m.cursor.Row > 0 {
				m.cursor.Row--
			}

		case "H":
			// Move task left
			return m, m.moveTaskLeft()

		case "L":
			// Move task right
			return m, m.moveTaskRight()

		case "n":
			// Create new task
			m.mode = ModeInput
			m.textInput.Reset()
			m.textInput.Focus()
			return m, textinput.Blink

		case "d":
			// Delete task
			return m, m.confirmDelete()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// updateInput handles input in text input mode
func (m BoardModel) updateInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Create task
			title := m.textInput.Value()
			if title != "" {
				cmd = m.createTask(title)
			}
			m.mode = ModeNormal
			m.textInput.Blur()
			return m, cmd

		case tea.KeyEsc:
			// Cancel
			m.mode = ModeNormal
			m.textInput.Blur()
			return m, nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// updateConfirm handles input in confirm mode
func (m BoardModel) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			// Confirm
			m.mode = ModeNormal
			if m.confirmAction != nil {
				return m, m.confirmAction()
			}
			return m, nil

		case "n", "N", "esc":
			// Cancel
			m.mode = ModeNormal
			m.confirmAction = nil
			return m, nil
		}
	}

	return m, nil
}

// adjustCursorRow adjusts the row when switching columns
func (m *BoardModel) adjustCursorRow() {
	columnName := m.columns[m.cursor.Column]
	tasksInColumn := m.tasks[columnName]

	if len(tasksInColumn) == 0 {
		m.cursor.Row = 0
		return
	}

	if m.cursor.Row >= len(tasksInColumn) {
		m.cursor.Row = len(tasksInColumn) - 1
	}
}

// moveTaskLeft moves the current task to the previous column
func (m *BoardModel) moveTaskLeft() tea.Cmd {
	if m.cursor.Column == 0 {
		return nil
	}

	columnName := m.columns[m.cursor.Column]
	tasksInColumn := m.tasks[columnName]
	if len(tasksInColumn) == 0 {
		return nil
	}

	currentTask := tasksInColumn[m.cursor.Row]
	newStatus := m.columns[m.cursor.Column-1]

	return m.moveTask(currentTask, newStatus)
}

// moveTaskRight moves the current task to the next column
func (m *BoardModel) moveTaskRight() tea.Cmd {
	if m.cursor.Column >= len(m.columns)-1 {
		return nil
	}

	columnName := m.columns[m.cursor.Column]
	tasksInColumn := m.tasks[columnName]
	if len(tasksInColumn) == 0 {
		return nil
	}

	currentTask := tasksInColumn[m.cursor.Row]
	newStatus := m.columns[m.cursor.Column+1]

	return m.moveTask(currentTask, newStatus)
}

// moveTask moves a task to a new status
func (m *BoardModel) moveTask(t *task.Task, newStatus string) tea.Cmd {
	t.Status = newStatus
	if err := task.SaveTask(t); err != nil {
		return nil // TODO: show error
	}

	return m.reloadTasks()
}

// createTask creates a new task
func (m *BoardModel) createTask(title string) tea.Cmd {
	// Create task in first column (TODO)
	status := m.columns[0]
	_, err := task.CreateTask(config.TasksDir, title, status)
	if err != nil {
		return nil // TODO: show error
	}

	// Move cursor to first column
	m.cursor.Column = 0
	m.cursor.Row = 0

	return m.reloadTasks()
}

// confirmDelete shows delete confirmation
func (m *BoardModel) confirmDelete() tea.Cmd {
	columnName := m.columns[m.cursor.Column]
	tasksInColumn := m.tasks[columnName]
	if len(tasksInColumn) == 0 {
		return nil
	}

	currentTask := tasksInColumn[m.cursor.Row]
	m.confirmMsg = fmt.Sprintf("Delete '%s'? (y/n)", currentTask.Title)
	m.mode = ModeConfirm

	m.confirmAction = func() tea.Cmd {
		if err := task.DeleteTask(currentTask); err != nil {
			return nil // TODO: show error
		}

		// Adjust cursor after deletion
		if m.cursor.Row > 0 {
			m.cursor.Row--
		}

		return m.reloadTasks()
	}

	return nil
}

// reloadTasks reloads all tasks from disk
func (m *BoardModel) reloadTasks() tea.Cmd {
	return func() tea.Msg {
		tasks, _ := task.ListTasks(config.TasksDir)
		return reloadMsg{tasks: tasks}
	}
}

// reloadMsg is sent when tasks need to be reloaded
type reloadMsg struct {
	tasks []*task.Task
}

// handleReload handles the reload message
func (m *BoardModel) handleReload(msg reloadMsg) {
	// Regroup tasks by status
	m.tasks = make(map[string][]*task.Task)
	for _, t := range msg.tasks {
		m.tasks[t.Status] = append(m.tasks[t.Status], t)
	}

	// Adjust cursor if needed
	m.adjustCursorRow()
}

// GetSelectedTask returns the currently selected task, or nil if none
func (m BoardModel) GetSelectedTask() *task.Task {
	columnName := m.columns[m.cursor.Column]
	tasksInColumn := m.tasks[columnName]

	if len(tasksInColumn) == 0 {
		return nil
	}

	if m.cursor.Row >= len(tasksInColumn) {
		return nil
	}

	return tasksInColumn[m.cursor.Row]
}

// renderTask renders a single task with appropriate styling
func (m BoardModel) renderTask(t *task.Task, selected bool) string {
	style := TaskStyle
	if selected {
		style = SelectedTaskStyle
	}
	return style.Render(t.Title)
}

// renderColumn renders a single column with its tasks
func (m BoardModel) renderColumn(colIndex int) string {
	columnName := m.columns[colIndex]
	tasksInColumn := m.tasks[columnName]

	// Calculate column width based on terminal width
	columnWidth := m.getColumnWidth()

	// Title
	title := ColumnTitleStyle.Width(columnWidth).Render(columnName)

	// Tasks
	var taskLines []string
	if len(tasksInColumn) == 0 {
		taskLines = append(taskLines, TaskStyle.Width(columnWidth).Render("(empty)"))
	} else {
		for i, t := range tasksInColumn {
			selected := colIndex == m.cursor.Column && i == m.cursor.Row
			taskLine := m.renderTask(t, selected)
			taskLines = append(taskLines, lipgloss.NewStyle().Width(columnWidth).Render(taskLine))
		}
	}

	// Add empty lines to fill column height
	minHeight := 10
	for len(taskLines) < minHeight {
		taskLines = append(taskLines, strings.Repeat(" ", columnWidth))
	}

	// Combine title and tasks
	content := title + "\n" + strings.Join(taskLines, "\n")

	// Apply column style
	return ColumnStyle.Width(columnWidth).Render(content)
}

// getColumnWidth calculates the width for each column
func (m BoardModel) getColumnWidth() int {
	if len(m.columns) == 0 {
		return 20
	}

	// Account for borders and padding (roughly 4 chars per column border)
	borderOverhead := len(m.columns) * 4
	availableWidth := m.width - borderOverhead - 2

	if availableWidth < len(m.columns)*15 {
		return 15 // Minimum width
	}

	return availableWidth / len(m.columns)
}

// renderHelp renders the help/status bar
func (m BoardModel) renderHelp() string {
	help := "hjkl/arrows: navigate | n: new | d: delete | Enter: open | H/L: move | q: quit"
	return HelpStyle.Render(help)
}

// View renders the board
func (m BoardModel) View() string {
	// Render all columns
	var columns []string
	for i := range m.columns {
		columns = append(columns, m.renderColumn(i))
	}

	// Place columns side by side
	board := lipgloss.JoinHorizontal(lipgloss.Top, columns...)

	// Add help bar at bottom
	var help string
	switch m.mode {
	case ModeInput:
		help = HelpStyle.Render("Enter task title (Enter to confirm, Esc to cancel):\n" + m.textInput.View())
	case ModeConfirm:
		help = HelpStyle.Render(m.confirmMsg)
	default:
		help = m.renderHelp()
	}

	return lipgloss.JoinVertical(lipgloss.Left, board, help)
}
