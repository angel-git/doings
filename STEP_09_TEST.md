# Step 9 Implementation Testing Guide

## What was implemented

Step 9 (Checklist Editing) has been successfully implemented with the following features:

### Files Modified
1. **internal/ui/detail.go** - Complete editing functionality
   - Added DetailMode enum (Normal, Input, Confirm)
   - Added mode, textInput, message, confirmMsg, and confirmation callbacks
   - Implemented toggle checkbox with Space key
   - Implemented add item with o (below) and O (above) keys
   - Implemented delete item with x key
   - Implemented save with s key
   - Added unsaved changes tracking with [*] indicator
   - Added confirmation dialog for unsaved changes on Esc
   - Updated view to show status messages and different modes

2. **internal/app/app.go** - Updated view switching logic
   - Added shouldExitDetail flag for proper confirmation flow
   - Updated updateDetail to check for unsaved changes before exit
   - Proper handling of confirmation dialogs

### Key Features Implemented
✅ **Space** toggles checkbox state ([ ] ↔ [x])  
✅ **o** adds new item below current position  
✅ **O** adds new item above current position  
✅ **x** deletes current checklist item  
✅ **s** saves changes to file  
✅ **[*]** indicator shown in header when modified  
✅ Status messages shown in help bar (e.g., "Saved!")  
✅ Esc with unsaved changes shows confirmation  
✅ Confirmation dialog: y=save & exit, n=discard & exit, c=cancel  

## How to Test

Run the application:
```bash
./doings
```

### Test Checklist

**Build and Basic Navigation:**
- [x] Application builds without errors
- [ ] Open a task with Enter
- [ ] Navigate checklist with j/k

**Toggle Checkboxes:**
- [ ] Press Space on unchecked item - becomes [x]
- [ ] Press Space on checked item - becomes [ ]
- [ ] Header shows [*] modified indicator
- [ ] Checked items appear dimmed

**Add Items:**
- [ ] Press 'o' on an item
- [ ] Enter text and press Enter
- [ ] New item appears below current item
- [ ] Press 'O' on an item
- [ ] Enter text and press Enter
- [ ] New item appears above current item
- [ ] New items inherit indent level
- [ ] Press Esc in input mode cancels

**Delete Items:**
- [ ] Press 'x' on an item
- [ ] Item is removed from list
- [ ] Cursor adjusts appropriately
- [ ] [*] indicator appears

**Save Changes:**
- [ ] Make changes (toggle, add, or delete)
- [ ] Press 's' to save
- [ ] "Saved!" message appears
- [ ] [*] indicator disappears
- [ ] Check .tasks/*.md file - changes are persisted

**Unsaved Changes Warning:**
- [ ] Make changes (don't save)
- [ ] Press Esc
- [ ] Confirmation prompt appears: "Unsaved changes. Save? (y/n/c)"
- [ ] Press 'y' - changes saved and return to board
- [ ] Make changes again, press Esc
- [ ] Press 'n' - changes discarded and return to board
- [ ] Make changes again, press Esc
- [ ] Press 'c' - stay in detail view with changes

**Edge Cases:**
- [ ] Delete all items - cursor handles empty list
- [ ] Add item to empty checklist
- [ ] Navigate with no items (shouldn't crash)
- [ ] Very long item text wraps properly

### Test with Sample Task

The test task at `.tasks/1738764000-test-task.md` has a good checklist for testing.

### Expected Key Bindings (Detail View)

**Navigation:**
- j/down: Move to next checklist item
- k/up: Move to previous checklist item

**Editing:**
- Space: Toggle checkbox
- o: Add item below
- O: Add item above
- x: Delete current item
- s: Save changes

**View Control:**
- Esc: Return to board (or show unsaved changes warning)

## Definition of Done

All items from plan/09-checklist-editing.md have been implemented:

- [x] Space toggles checkbox state
- [x] 'o' adds new item below
- [x] 'x' deletes current item
- [x] 's' saves changes to file
- [x] Modified indicator shown when unsaved ([*])
- [x] Esc with changes shows confirmation
- [x] Saved file matches expected markdown format

## Implementation Notes

### Indentation
- New items inherit the indent level of the current item
- Indent level is preserved when saved to markdown (4 spaces per level)

### Modes
The detail view has three modes:
1. **Normal Mode**: Navigate and edit
2. **Input Mode**: Enter text for new items
3. **Confirm Mode**: Handle unsaved changes dialog

### Save Flow
1. User presses 's' in normal mode
2. SaveTask command is sent
3. saveResultMsg is returned
4. If successful: modified=false, message="Saved!"
5. If error: message shows error

### Unsaved Changes Flow
1. User presses Esc with modified=true
2. Confirmation dialog is shown
3. User presses y: save and exit
4. User presses n: discard and exit
5. User presses c/Esc: stay in detail view

### File Format
Changes are persisted using the existing ToMarkdown() method which maintains:
- Title with # prefix
- status = "VALUE" line
- Description between --- separators
- Checklist with proper indentation and checkbox states
