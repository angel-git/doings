# Step 07: Task Movement & Creation/Deletion

## Goal
Implement moving tasks between columns, creating new tasks, and deleting tasks.

## Tasks

### 7.1 Move Task Between Columns

Key bindings:
- `H` (Shift+h) - Move task to previous column
- `L` (Shift+l) - Move task to next column

Implementation:
1. Get current task
2. Determine target column
3. Update task.Status to new column name
4. Save task to file
5. Reload tasks (or update in-memory state)

```go
func (m *BoardModel) moveTaskLeft() tea.Cmd
func (m *BoardModel) moveTaskRight() tea.Cmd
```

### 7.2 Create New Task

Key binding: `n`

Flow:
1. Enter input mode (show text input)
2. User types task title
3. Press Enter to confirm (or Esc to cancel)
4. Create task with status = first column (e.g., "TODO")
5. Add task to board
6. Select the new task

Use Bubble Tea's text input component:
```go
import "github.com/charmbracelet/bubbles/textinput"
```

### 7.3 Delete Task

Key binding: `d`

Flow:
1. Show confirmation: "Delete 'Task Title'? (y/n)"
2. If `y`, delete file and remove from board
3. If `n` or Esc, cancel
4. Adjust cursor if needed

### 7.4 Board State Management

Add fields to BoardModel:
```go
type BoardModel struct {
    // ... existing fields ...
    mode        Mode           // Normal, Input, Confirm
    textInput   textinput.Model
    confirmMsg  string
    confirmAction func() tea.Cmd
}

type Mode int
const (
    ModeNormal Mode = iota
    ModeInput
    ModeConfirm
)
```

## Testing This Step

```bash
./doings
```

**Test actions**:
1. Press `n`, type "New Task", press Enter - task created in TODO
2. Press `L` on a task - moves to next column
3. Press `H` - moves back
4. Press `d`, then `y` - task deleted
5. Press `d`, then `n` - deletion cancelled

**File checks**:
- New task file appears in `.tasks/`
- Status changes are persisted
- Deleted task file is removed

## Files to Modify
- `internal/ui/board.go` - Add modes and handlers

## Definition of Done
- [x] `H`/`L` moves task between columns
- [x] Status change is saved to file
- [x] `n` creates new task with text input
- [x] `d` deletes task with confirmation
- [x] Cancelling input/confirm returns to normal mode
- [x] Edge cases handled (empty columns, first/last column)

## Status
**COMPLETED** - All tasks finished successfully on Thu Feb 05 2026
