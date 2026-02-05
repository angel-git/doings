# Step 10: Polish & Error Handling

## Goal
Add finishing touches, error handling, and improve user experience.

## Tasks

### 10.1 Startup Warnings

Display warnings for:
- Malformed task files (couldn't parse)
- Tasks with invalid status (not matching any column)
- Missing config values

Show warnings at startup, then continue:
```
Warning: Could not parse .tasks/1234-broken.md: missing status line
Warning: Task "Old Task" has unknown status "ARCHIVED"

Press any key to continue...
```

Or show in a status bar area.

### 10.2 Empty State

When no tasks exist:
```
┌─────────┐ ┌─────────┐ ┌─────────┐
│  TODO   │ │  DOING  │ │  DONE   │
├─────────┤ ├─────────┤ ├─────────┤
│         │ │         │ │         │
│  No     │ │         │ │         │
│  tasks  │ │         │ │         │
│         │ │         │ │         │
└─────────┘ └─────────┘ └─────────┘

Press 'n' to create your first task
```

### 10.3 Status Bar Improvements

Show contextual information:
- Current mode (Normal, Input, etc.)
- Error/success messages (fade after 3 seconds)
- Modified indicator in detail view

### 10.4 Edge Case Handling

- Empty task title when creating (show error, don't create)
- Very long task titles (truncate in board view)
- Many tasks in column (scrolling or truncation)
- Very long description/checklist (scrolling in detail view)

### 10.5 Visual Polish

- Consistent colors across views
- Clear visual hierarchy
- Readable contrast
- Nice borders and spacing

### 10.6 Help Screen (Optional)

Key binding: `?`

Show full list of keybindings:
```
┌─────── Keyboard Shortcuts ───────┐
│                                  │
│  Navigation                      │
│    h/j/k/l  Move cursor          │
│    Enter    Open task            │
│    Esc      Go back              │
│                                  │
│  Actions                         │
│    n        New task             │
│    d        Delete task          │
│    H/L      Move task            │
│                                  │
│  Detail View                     │
│    Space    Toggle checkbox      │
│    o        Add item             │
│    x        Delete item          │
│    s        Save                 │
│                                  │
│         Press ? to close         │
└──────────────────────────────────┘
```

## Testing This Step

**Error handling tests**:
1. Create malformed `.md` file manually - app shows warning
2. Create task with unknown status - shows warning
3. Delete all tasks - empty state shown

**Edge case tests**:
1. Create task with empty name - error shown
2. Create task with very long name - truncated in view
3. Create many tasks - verify layout handles it

**Polish tests**:
1. Overall visual appearance is clean
2. Colors are consistent
3. No visual glitches on resize

## Files to Modify
- `internal/ui/board.go` - Warnings, empty state
- `internal/ui/detail.go` - Status messages
- `internal/ui/styles.go` - Color refinements
- `main.go` - Startup warning logic

## Definition of Done
- [ ] Malformed files show warnings but don't crash
- [ ] Empty state is handled gracefully
- [ ] Long content is handled (truncate/scroll)
- [ ] Visual appearance is polished
- [ ] Error messages are helpful
- [ ] App feels complete and usable
