# Step 05: Board Model & Navigation

## Goal
Create the Bubble Tea model for the board with navigation between columns and tasks.

## Tasks

### 5.1 Create Board Model (`internal/ui/board.go`)

```go
type BoardModel struct {
    columns      []string        // Column names from config
    tasks        map[string][]*task.Task  // Tasks grouped by status
    cursor       Cursor          // Current position
    width        int             // Terminal width
    height       int             // Terminal height
}

type Cursor struct {
    Column int  // Which column (0-indexed)
    Row    int  // Which task in column (0-indexed)
}
```

### 5.2 Implement Bubble Tea Interface

```go
func (m BoardModel) Init() tea.Cmd
func (m BoardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m BoardModel) View() string
```

### 5.3 Implement Navigation (Update method)

Key bindings:
- `h` / `left` - Move cursor left (previous column)
- `l` / `right` - Move cursor right (next column)
- `j` / `down` - Move cursor down (next task)
- `k` / `up` - Move cursor up (previous task)
- `q` - Quit application

Navigation rules:
- Cursor wraps at edges or stops (your choice - stopping is simpler)
- When moving to empty column, cursor row becomes 0
- When column has fewer tasks than current row, adjust row

### 5.4 Wire Up in main.go

```go
func main() {
    // ... initialization ...
    
    p := tea.NewProgram(ui.NewBoardModel(config, tasks))
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## Testing This Step

```bash
./doings
```

**Test actions**:
1. Press `h`, `j`, `k`, `l` - cursor should move
2. Press `q` - app should quit
3. Navigation should not crash on empty columns

**Note**: View will be minimal (just showing cursor position is fine for now)

## Files to Create/Modify
- `internal/ui/board.go` - Board model
- `main.go` - Wire up Bubble Tea program

## Definition of Done
- [ ] App starts with Bubble Tea
- [ ] `hjkl` keys move cursor
- [ ] Cursor position is tracked correctly
- [ ] `q` quits the app
- [ ] Navigation handles empty columns gracefully
