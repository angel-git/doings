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
	ViewHelp
)

// AppModel is the top-level model that manages view switching
type AppModel struct {
	view     View
	board    ui.BoardModel
	detail   ui.DetailModel
	columns  []string
	warnings []string // Startup warnings
	showHelp bool     // Show help overlay
	width    int      // Terminal width
	height   int      // Terminal height
}

// NewAppModel creates a new app model
func NewAppModel(columns []string, tasks []*task.Task, warnings []string) AppModel {
	return AppModel{
		view:     ViewBoard,
		board:    ui.NewBoardModel(columns, tasks, warnings),
		columns:  columns,
		warnings: warnings,
	}
}

// Init initializes the app model
func (m AppModel) Init() tea.Cmd {
	return m.board.Init()
}

// Update handles messages and routes them to the appropriate view
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle window size globally
	if windowMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = windowMsg.Width
		m.height = windowMsg.Height
	}

	// Handle help toggle globally
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "?" {
			m.showHelp = !m.showHelp
			return m, nil
		}
	}

	// If help is showing, any other key closes it
	if m.showHelp {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() != "?" {
				m.showHelp = false
				return m, nil
			}
		}
	}

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
	// Check for Enter key to switch to detail view (only in normal mode)
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "enter" && m.board.IsNormalMode() {
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
	// Pass window size messages to board view too
	if windowMsg, ok := msg.(tea.WindowSizeMsg); ok {
		updatedBoard, _ := m.board.Update(windowMsg)
		m.board = updatedBoard.(ui.BoardModel)
	}

	// Check for Esc key to return to board view (before detail processes it)
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "esc" {
			// Only handle Esc if we're in normal mode (not in confirm/input mode)
			if m.detail.IsNormalMode() {
				// Check for unsaved changes
				if m.detail.HasUnsavedChanges() {
					// Show confirmation dialog
					m.detail.ShowUnsavedConfirmation()
					return m, nil
				}
				// No unsaved changes, go back
				m.view = ViewBoard
				return m, nil
			}
		}
	}

	// Update detail model
	updatedDetail, cmd := m.detail.Update(msg)
	m.detail = updatedDetail

	// Handle confirmation result
	if confirmMsg, ok := msg.(ui.ConfirmResultMsg); ok {
		switch confirmMsg.Action {
		case "save", "discard":
			// Exit to board on save or discard
			m.view = ViewBoard
		case "cancel":
			// Stay in detail view
		}
		return m, cmd
	}

	return m, cmd
}

// View renders the current view
func (m AppModel) View() string {
	var baseView string

	switch m.view {
	case ViewBoard:
		baseView = m.board.View()
	case ViewDetail:
		baseView = m.detail.View()
	default:
		baseView = ""
	}

	// Overlay help screen if showing
	if m.showHelp {
		return ui.RenderHelpScreen(m.width, m.height)
	}

	return baseView
}
