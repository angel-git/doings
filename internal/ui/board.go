package ui

import (
	"strings"

	"doings/internal/task"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Cursor tracks the current position on the board
type Cursor struct {
	Column int // Which column (0-indexed)
	Row    int // Which task in column (0-indexed)
}

// BoardModel represents the main board view
type BoardModel struct {
	columns []string                // Column names from config
	tasks   map[string][]*task.Task // Tasks grouped by status
	cursor  Cursor                  // Current position
	width   int                     // Terminal width
	height  int                     // Terminal height
}

// NewBoardModel creates a new board model
func NewBoardModel(columns []string, tasks []*task.Task) BoardModel {
	// Group tasks by status
	tasksByStatus := make(map[string][]*task.Task)
	for _, t := range tasks {
		tasksByStatus[t.Status] = append(tasksByStatus[t.Status], t)
	}

	return BoardModel{
		columns: columns,
		tasks:   tasksByStatus,
		cursor:  Cursor{Column: 0, Row: 0},
		width:   80,
		height:  24,
	}
}

// Init initializes the board model
func (m BoardModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m BoardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "h", "left":
			// Move left
			if m.cursor.Column > 0 {
				m.cursor.Column--
				m.adjustCursorRow()
			}

		case "l", "right":
			// Move right
			if m.cursor.Column < len(m.columns)-1 {
				m.cursor.Column++
				m.adjustCursorRow()
			}

		case "j", "down":
			// Move down
			columnName := m.columns[m.cursor.Column]
			tasksInColumn := m.tasks[columnName]
			if m.cursor.Row < len(tasksInColumn)-1 {
				m.cursor.Row++
			}

		case "k", "up":
			// Move up
			if m.cursor.Row > 0 {
				m.cursor.Row--
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// adjustCursorRow adjusts the row when switching columns
func (m *BoardModel) adjustCursorRow() {
	columnName := m.columns[m.cursor.Column]
	tasksInColumn := m.tasks[columnName]

	// If column is empty, set row to 0
	if len(tasksInColumn) == 0 {
		m.cursor.Row = 0
		return
	}

	// If current row is beyond the tasks in this column, adjust
	if m.cursor.Row >= len(tasksInColumn) {
		m.cursor.Row = len(tasksInColumn) - 1
	}
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
	help := m.renderHelp()

	return lipgloss.JoinVertical(lipgloss.Left, board, help)
}
