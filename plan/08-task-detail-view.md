# Step 08: Task Detail View

## Goal
Implement the detail view for viewing and navigating task contents.

## Tasks

### 8.1 Create Detail Model (`internal/ui/detail.go`)

```go
type DetailModel struct {
    task          *task.Task
    cursor        int           // Selected line in checklist
    width         int
    height        int
    modified      bool          // Has unsaved changes
}
```

### 8.2 Implement View

Layout:
```
┌─────────────────────────────────────────┐
│ # Task Title                     [TODO] │
├─────────────────────────────────────────┤
│ This is the task description that can   │
│ span multiple lines.                    │
├─────────────────────────────────────────┤
│ - [ ] First step                        │
│ - [x] Completed step                    │
│     - [ ] Nested subtask  <-- selected  │
│ - [ ] Another step                      │
└─────────────────────────────────────────┘

j/k: navigate | Space: toggle | Esc: back | s: save
```

### 8.3 Implement Navigation

Key bindings:
- `j` / `down` - Move to next checklist item
- `k` / `up` - Move to previous checklist item
- `Esc` - Return to board view

### 8.4 Switch Between Views

In main app model:
```go
type AppModel struct {
    board   BoardModel
    detail  DetailModel
    view    View  // ViewBoard or ViewDetail
}

type View int
const (
    ViewBoard View = iota
    ViewDetail
)
```

Or use Bubble Tea's built-in model switching pattern.

### 8.5 Open Detail View

Key binding: `Enter` (from board view)

Flow:
1. Get selected task
2. Create DetailModel with task
3. Switch to detail view

## Testing This Step

```bash
./doings
```

**Test actions**:
1. Select a task, press Enter - detail view opens
2. See task title, status, description, checklist
3. Press `j`/`k` - navigate checklist
4. Press Esc - return to board

**Note**: Editing will be added in the next step

## Files to Create/Modify
- `internal/ui/detail.go` - Detail view model
- `internal/ui/board.go` - Handle Enter key
- `internal/app/app.go` - (optional) App-level model for view switching
- `main.go` - Update to use app model

## Definition of Done
- [ ] Enter opens detail view for selected task
- [ ] Task title, status, description displayed
- [ ] Checklist items displayed with correct indentation
- [ ] `j`/`k` navigates checklist
- [ ] Selected checklist item is highlighted
- [ ] Esc returns to board view
