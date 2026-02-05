# Step 06: Board View Rendering

## Goal
Render the board with columns and tasks using Lip Gloss styling.

## Tasks

### 6.1 Create Styles (`internal/ui/styles.go`)

```go
var (
    // Column styles
    ColumnStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(0, 1)
    
    ColumnTitleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("205"))
    
    // Task styles
    TaskStyle = lipgloss.NewStyle().
        Padding(0, 1)
    
    SelectedTaskStyle = lipgloss.NewStyle().
        Padding(0, 1).
        Background(lipgloss.Color("62")).
        Foreground(lipgloss.Color("230"))
    
    // Status bar
    HelpStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("241"))
)
```

### 6.2 Implement View Method

Layout:
```
┌─────────┐ ┌─────────┐ ┌─────────┐
│  TODO   │ │  DOING  │ │  DONE   │
├─────────┤ ├─────────┤ ├─────────┤
│ Task 1  │ │ Task 3  │ │ Task 5  │
│ Task 2  │ │         │ │ Task 6  │
│         │ │         │ │         │
└─────────┘ └─────────┘ └─────────┘

hjkl: navigate | n: new | d: delete | Enter: open | H/L: move | q: quit
```

### 6.3 Handle Terminal Resizing

```go
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
```

Calculate column widths based on terminal width:
- Equal width for each column
- Account for borders and padding

### 6.4 Render Components

Functions:
```go
func (m BoardModel) renderColumn(colIndex int) string
func (m BoardModel) renderTask(task *Task, selected bool) string
func (m BoardModel) renderHelp() string
```

## Testing This Step

```bash
./doings
```

**Visual checks**:
1. Columns displayed horizontally
2. Tasks appear under correct columns
3. Selected task is highlighted
4. Resizing terminal adjusts layout
5. Help bar shows available commands

## Files to Create/Modify
- `internal/ui/styles.go` - Lip Gloss styles
- `internal/ui/board.go` - Update View method

## Definition of Done
- [ ] Board displays columns horizontally
- [ ] Tasks appear in their respective columns
- [ ] Selected task is visually distinct
- [ ] Terminal resize is handled
- [ ] Help bar shows key bindings
- [ ] Layout looks clean and readable
