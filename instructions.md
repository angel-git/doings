# Project: Terminal Task Board (v1)

## Goal

Build a **solo-developer, local-first TUI application** to manage tasks using

Markdown files stored in the repository.

The tool is **not** a team coordination system.
It is designed for **personal workflow management**, optimized for:
- terminal usage
- git version control
- minimal friction
- fast iteration

The TUI acts as an **editor and visualizer** for Markdown-based tasks.

---

## Core Principles

- Source of truth is **plain text files**
- No server, no sync, no external services
- Git is optional but first-class
- Favor usability over configurability
- v1 should be usable in < 2 days of development

---

## Tech Stack

- Language: **Go**
- TUI framework: **Bubble Tea**
- Styling: **Lip Gloss**
- Storage:
    - Tasks: Markdown files
    - Board config: TOML file

---

## File Structure

On initialization, the app creates (if it doesn't exist) a `.tasks/` directory in the current working directory with the following structure:

```
.tasks/
  config.toml
  <task-id>.md
```


### config.toml

Defines configuration of the project, ie: the columns of the board and their order.

Example:

```toml
[board]
columns = ["TODO", "DOING", "DONE"]
```

### Task Files

- One Markdown file per task
- Filename acts as the task ID, it starts with a sequence number for ordering, e.g. `001-my-first-task.md`
- File is human-editable outside the app

#### Task Markdown Format (v1)

```markdown
# Task title
status = "TODO"
---
Some description
---
- [ ] First step
- [ ] Second step
    - [ ] Subtask
- [x] Completed step
```

Rules:
- status must match a column defined in board.toml
- Bullet points represent actionable items
- Nested bullets are allowed
- No advanced metadata parsing required in v1


## TUI Layout (v1)

### Main Board View

- Columns displayed horizontally
- Tasks listed vertically under each column
- Selected task is highlighted

Keyboard interactions:
- hjkl to navigate
- n → create new task
- d → delete selected task (with confirmation)
- Enter → open task detail view
- shift+h / shift+l → move task between columns
- Moving a task updates the status field in its Markdown file.

### Task Detail View

Displays the contents of the task Markdown file.

Supported interactions:
- Navigate bullet points
- Toggle checkbox ([ ] ↔ [x]) with Space
- Basic insert/delete of bullet lines
- Save changes to disk explicitly
- Esc → return to board view
- The app edits the Markdown file directly.
- No full Markdown editor is required — only bullet lists and headers.

## Non-Goals for v1

Explicitly out of scope:
- Multi-user support
- Sync or networking
- Task dependencies
- Search or filtering
- Custom keybindings
- Plugins or extensions
- Undo/redo beyond file-level git usage

## Success Criteria for v1

The project is considered successful if:
- A user can manage daily tasks entirely from the TUI
- All task data is readable and editable without the app
- The tool feels fast and low-friction
- The user wants to use it again the next day

