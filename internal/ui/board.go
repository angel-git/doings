package ui

import (
	"fmt"
	"strings"

	"doings/internal/task"
	tea "github.com/charmbracelet/bubbletea"
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

// View renders the board
func (m BoardModel) View() string {
	var b strings.Builder

	b.WriteString("Terminal Task Board\n")
	b.WriteString("===================\n\n")

	// Show columns
	b.WriteString("Columns: ")
	for i, col := range m.columns {
		if i == m.cursor.Column {
			b.WriteString("[" + col + "] ")
		} else {
			b.WriteString(col + " ")
		}
	}
	b.WriteString("\n\n")

	// Show current column tasks
	currentColumn := m.columns[m.cursor.Column]
	tasksInColumn := m.tasks[currentColumn]

	b.WriteString(fmt.Sprintf("Column: %s (%d tasks)\n", currentColumn, len(tasksInColumn)))
	b.WriteString("---\n")

	if len(tasksInColumn) == 0 {
		b.WriteString("  (empty)\n")
	} else {
		for i, t := range tasksInColumn {
			if i == m.cursor.Row {
				b.WriteString(fmt.Sprintf("> %s\n", t.Title))
			} else {
				b.WriteString(fmt.Sprintf("  %s\n", t.Title))
			}
		}
	}

	b.WriteString("\n")
	b.WriteString("Controls: hjkl/arrows to navigate | q to quit\n")

	return b.String()
}
