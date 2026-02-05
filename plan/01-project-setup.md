# Step 01: Project Setup

## Goal
Initialize the Go project with proper module structure and dependencies.

## Tasks

### 1.1 Initialize Go Module
```bash
go mod init doings
```

### 1.2 Install Dependencies
```bash
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss
go get github.com/BurntSushi/toml
```

### 1.3 Create Directory Structure
```
doings/
├── main.go
├── internal/
│   ├── config/
│   ├── task/
│   ├── ui/
│   └── app/
```

### 1.4 Create Minimal main.go
A simple entry point that:
- Prints "Terminal Task Board" 
- Exits cleanly

## Testing This Step

```bash
go build -o doings
./doings
```

**Expected output**: 
- Builds without errors
- Prints startup message
- Exits cleanly

## Files to Create
- `main.go` - Entry point
- `internal/` directories (empty for now)

## Definition of Done
- [x] `go build` succeeds
- [x] Running `./doings` shows a message and exits
- [x] All dependencies are in `go.sum`

## Status
**COMPLETED** - All tasks finished successfully on Thu Feb 05 2026
