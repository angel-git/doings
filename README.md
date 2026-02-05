# Doings - Terminal Task Board

A solo-developer, local-first TUI application for managing tasks using Markdown files.

## ✨ Features

### Board View
- **Kanban-style columns** with customizable statuses
- **Vim-style navigation** (hjkl or arrow keys)
- **Create tasks** quickly with 'n' key
- **Move tasks** between columns with H/L keys
- **Delete tasks** with confirmation
- **Task counts** displayed in column headers
- **Scrolling support** for columns with many tasks
- **Empty state** with helpful onboarding message

### Detail View
- **Full task display** with title, status, description, and checklist
- **Navigate checklist** with j/k keys
- **Toggle checkboxes** with Space
- **Add items** with o (below) or O (above)
- **Delete items** with x key
- **Save changes** with s key
- **Unsaved changes warning** when exiting

### User Experience
- **Startup warnings** for malformed files or invalid statuses
- **Status messages** for errors and success confirmations
- **Help screen** (? key) with all keyboard shortcuts
- **Visual polish** with consistent color palette
- **Responsive** to terminal resizing
- **Title truncation** for long task names
- **Professional styling** with Lip Gloss

## 📁 File Structure

```
doings/
├── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go       # TOML configuration handling
│   ├── task/
│   │   ├── task.go         # Task struct and markdown parsing
│   │   └── storage.go      # File operations (CRUD)
│   ├── ui/
│   │   ├── board.go        # Board view (kanban)
│   │   ├── detail.go       # Task detail view
│   │   ├── help.go         # Help screen
│   │   └── styles.go       # Lip Gloss styles
│   └── app/
│       └── app.go          # Main app model (view switching)
├── .tasks/                 # Task storage directory
│   ├── config.toml         # Board configuration
│   └── *.md                # Task files
└── plan/                   # Development plan documents
    ├── 00-overview.md
    ├── 01-project-setup.md
    ├── 02-config-file-structure.md
    ├── 03-task-parsing.md
    ├── 04-task-crud.md
    ├── 05-board-model-navigation.md
    ├── 06-board-view-rendering.md
    ├── 07-task-movement-crud.md
    ├── 08-task-detail-view.md
    ├── 09-checklist-editing.md
    └── 10-polish-error-handling.md
```

## 🚀 Quick Start

### Build
```bash
go build
```

### Run
```bash
./doings
```

On first run, the application will:
1. Create `.tasks/` directory
2. Generate default `config.toml`
3. Show empty state with instructions

## ⌨️ Keyboard Shortcuts

### Board View
- **h/j/k/l** or **arrows**: Navigate
- **n**: Create new task
- **d**: Delete task (with confirmation)
- **H/L**: Move task left/right between columns
- **Enter**: Open task detail
- **?**: Show help
- **q**: Quit

### Detail View
- **j/k**: Navigate checklist items
- **Space**: Toggle checkbox
- **o**: Add item below
- **O**: Add item above
- **x**: Delete current item
- **s**: Save changes
- **Esc**: Return to board (warns if unsaved changes)

### Help
- **?**: Toggle help screen
- **Any key**: Close help

## 📝 Task File Format

Tasks are stored as Markdown files in `.tasks/`:

```markdown
# Task Title
status = "TODO"
---
Task description goes here.
Can be multiple lines.
---
- [ ] First checklist item
- [x] Completed item
    - [ ] Nested item (4 spaces per level)
- [ ] Another item
```

### File Naming
Files use timestamp-based IDs: `1738764000-task-name.md`

### Status Values
Configured in `.tasks/config.toml`:
```toml
[board]
columns = ["TODO", "DOING", "DONE"]
```

## 🎨 Visual Design

### Color Palette
- **Primary** (Teal #62): Borders, selection
- **Accent** (Pink #205): Titles, headers
- **Success** (Green #42): Success messages
- **Warning** (Orange #208): Warnings
- **Error** (Red #196): Error messages
- **Muted** (Gray #240): Secondary info
- **Text** (Light gray #252): Main text
- **Highlight** (Yellow #230): Selected text

### Features
- Rounded borders
- Clear visual hierarchy
- Accessible contrast
- Consistent spacing
- Professional appearance

## 🛠️ Error Handling

### Graceful Degradation
- **Malformed files**: Show warning, skip file
- **Invalid status**: Show warning, include task anyway
- **Empty titles**: Validation with error message
- **No tasks**: Helpful empty state message

### User Feedback
- **Warnings**: Dismissible overlay at startup
- **Errors**: Clear messages in status bar
- **Success**: Confirmation messages (e.g., "Saved!")
- **Validation**: Immediate feedback on invalid input

## 🔄 State Management

### Persistence
- All changes immediately saved to disk
- Markdown files are the source of truth
- No database or complex state management

### Reloading
- Tasks reloaded after create/delete/move operations
- Detail view shows live task state
- Unsaved changes tracked with [*] indicator

## 📦 Dependencies

- **Bubble Tea**: TUI framework
- **Lip Gloss**: Terminal styling
- **Bubbles**: TUI components (text input)
- **BurntSushi/toml**: TOML parsing

## 🏗️ Architecture

### Model-View-Update (MVU)
Built on Bubble Tea's MVU pattern:
- **Model**: Application state
- **View**: Rendering logic
- **Update**: Message handling

### View Hierarchy
```
AppModel (view switcher)
  ├── BoardModel (kanban view)
  ├── DetailModel (task detail)
  └── HelpScreen (help overlay)
```

### Message Flow
1. User input → KeyMsg
2. AppModel routes to active view
3. View updates state
4. Commands trigger side effects
5. Results sent back as messages

## ✅ Development Steps Completed

1. ✅ **Project Setup** - Go module, dependencies
2. ✅ **Config & File Structure** - .tasks/ initialization
3. ✅ **Task Parsing** - Markdown to struct conversion
4. ✅ **Task CRUD** - Create, read, update, delete operations
5. ✅ **Board Model & Navigation** - Cursor movement, state
6. ✅ **Board View Rendering** - Kanban display with Lip Gloss
7. ✅ **Task Movement & CRUD** - H/L movement, create, delete
8. ✅ **Task Detail View** - Full task display with navigation
9. ✅ **Checklist Editing** - Toggle, add, delete, save
10. ✅ **Polish & Error Handling** - Warnings, edge cases, help

## 🎯 Use Cases

### Perfect For:
- Solo developers managing personal tasks
- Local-first workflow enthusiasts
- Developers who live in the terminal
- Quick task capture and management
- Simple kanban-style project tracking

### Not For:
- Team collaboration (no sync/sharing)
- Complex project management
- Mobile/web access
- Real-time updates across devices

## 📄 License

This project was created as a learning exercise following a step-by-step development plan.

## 🙏 Acknowledgments

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) by Charm
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) by Charm
- [Bubbles](https://github.com/charmbracelet/bubbles) by Charm

---

**Made with ❤️ for terminal lovers**
