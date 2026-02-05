package app

import (
	"doings/internal/task"
	"doings/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

// View represents which view is currently active
type View int

const (
	ViewBoard View = iota
	ViewDetail
)

// AppModel is the top-level model that manages view switching
type AppModel struct {
	view    View
	board   ui.BoardModel
	detail  ui.DetailModel
	columns []string
}

// NewAppModel creates a new app model
func NewAppModel(columns []string, tasks []*task.Task) AppModel {
	return AppModel{
		view:    ViewBoard,
		board:   ui.NewBoardModel(columns, tasks),
		columns: columns,
	}
}

// Init initializes the app model
func (m AppModel) Init() tea.Cmd {
	return m.board.Init()
}

// Update handles messages and routes them to the appropriate view
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.view {
	case ViewBoard:
		return m.updateBoard(msg)
	case ViewDetail:
		return m.updateDetail(msg)
	default:
		return m, nil
	}
}

// updateBoard handles updates when board view is active
func (m AppModel) updateBoard(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Check for Enter key to switch to detail view
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "enter" {
			// Get the currently selected task
			selectedTask := m.board.GetSelectedTask()
			if selectedTask != nil {
				m.detail = ui.NewDetailModel(selectedTask)
				m.view = ViewDetail
				return m, nil
			}
		}
	}

	// Pass window size messages to detail view too
	if windowMsg, ok := msg.(tea.WindowSizeMsg); ok {
		updatedDetail, _ := m.detail.Update(windowMsg)
		m.detail = updatedDetail
	}

	// Update board model
	updatedBoard, cmd := m.board.Update(msg)
	m.board = updatedBoard.(ui.BoardModel)
	return m, cmd
}

// updateDetail handles updates when detail view is active
func (m AppModel) updateDetail(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Check for Esc key to return to board view
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "esc" {
			m.view = ViewBoard
			return m, nil
		}
	}

	// Pass window size messages to board view too
	if windowMsg, ok := msg.(tea.WindowSizeMsg); ok {
		updatedBoard, _ := m.board.Update(windowMsg)
		m.board = updatedBoard.(ui.BoardModel)
	}

	// Update detail model
	updatedDetail, cmd := m.detail.Update(msg)
	m.detail = updatedDetail
	return m, cmd
}

// View renders the current view
func (m AppModel) View() string {
	switch m.view {
	case ViewBoard:
		return m.board.View()
	case ViewDetail:
		return m.detail.View()
	default:
		return ""
	}
}
