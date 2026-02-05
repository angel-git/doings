# Step 09: Checklist Editing

## Goal
Implement checklist interactions: toggling checkboxes, adding/removing items.

## Tasks

### 9.1 Toggle Checkbox

Key binding: `Space`

Implementation:
1. Get current checklist item
2. Toggle `Checked` state
3. Mark task as modified

```go
func (m *DetailModel) toggleCheckbox() {
    item := &m.task.Checklist[m.cursor]
    item.Checked = !item.Checked
    m.modified = true
}
```

### 9.2 Add New Checklist Item

Key bindings:
- `o` - Add item below current (like vim)
- `O` - Add item above current

Flow:
1. Enter input mode
2. User types item text
3. Enter confirms, Esc cancels
4. Insert item at appropriate position
5. New item inherits indent level of current item

### 9.3 Delete Checklist Item

Key binding: `dd` (vim-style) or just `x`

For v1, let's use `x` for simplicity:
1. Remove current item
2. Adjust cursor if needed
3. Mark as modified

### 9.4 Save Changes

Key binding: `s`

Implementation:
1. Serialize task to markdown
2. Write to file
3. Clear modified flag
4. Show brief "Saved!" message

### 9.5 Unsaved Changes Warning

When pressing Esc with unsaved changes:
1. Show "Unsaved changes. Save? (y/n/c)"
2. `y` - save and exit
3. `n` - discard and exit
4. `c` - cancel (stay in detail view)

### 9.6 Update Detail Model

```go
type DetailModel struct {
    // ... existing fields ...
    mode        DetailMode
    textInput   textinput.Model
    message     string        // Temporary status message
}

type DetailMode int
const (
    DetailModeNormal DetailMode = iota
    DetailModeInput
    DetailModeConfirm
)
```

## Testing This Step

```bash
./doings
```

**Test actions**:
1. Open task detail
2. Press Space on checkbox - toggles [ ] ↔ [x]
3. Press `o`, type text, Enter - new item added
4. Press `x` - item deleted
5. Press `s` - changes saved
6. Make changes, press Esc - warning shown
7. Check `.tasks/*.md` file reflects changes

## Files to Modify
- `internal/ui/detail.go` - Add editing functionality
- `internal/task/task.go` - Ensure ToMarkdown handles all cases

## Definition of Done
- [ ] Space toggles checkbox state
- [ ] `o` adds new item below
- [ ] `x` deletes current item
- [ ] `s` saves changes to file
- [ ] Modified indicator shown when unsaved
- [ ] Esc with changes shows confirmation
- [ ] Saved file matches expected markdown format
