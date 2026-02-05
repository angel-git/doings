# Terminal Task Board - Development Plan Overview

## Project Summary
A solo-developer, local-first TUI application to manage tasks using Markdown files.

## Tech Stack
- **Language**: Go
- **TUI Framework**: Bubble Tea
- **Styling**: Lip Gloss
- **Storage**: Markdown files + TOML config

## Key Decisions
- **Module name**: `doings`
- **Task IDs**: Timestamp-based (e.g., `1738764000-my-task-name.md`)
- **Status format**: Plain text in markdown (not frontmatter)
- **Auto-initialization**: Create `.tasks/` and default config if missing
- **Error handling**: Warn and skip malformed markdown files
- **Testing**: Manual testing only for v1

## Development Steps

| Step | Description | Testable Outcome |
|------|-------------|------------------|
| 01 | Project Setup | `go build` succeeds, app runs |
| 02 | Config & File Structure | `.tasks/` created with `config.toml` |
| 03 | Task Parsing | Parse markdown files into task structs |
| 04 | Task CRUD Operations | Create, read, update, delete task files |
| 05 | Board Model & Navigation | Navigate between columns and tasks |
| 06 | Board View Rendering | Display columns and tasks in TUI |
| 07 | Task Movement | Move tasks between columns (H/L keys) |
| 08 | Task Detail View | Open and display task details |
| 09 | Checklist Editing | Toggle checkboxes, add/remove bullets |
| 10 | Polish & Error Handling | Warnings, confirmations, edge cases |

## File Structure (Final)
```
doings/
├── go.mod
├── go.sum
├── main.go
├── internal/
│   ├── config/
│   │   └── config.go       # TOML config handling
│   ├── task/
│   │   ├── task.go         # Task struct and parsing
│   │   └── storage.go      # File operations
│   ├── ui/
│   │   ├── board.go        # Board view (main)
│   │   ├── detail.go       # Task detail view
│   │   └── styles.go       # Lip Gloss styles
│   └── app/
│       └── app.go          # Main app model
└── plan/
    └── *.md                # Step-by-step plans
```
