# Step 03: Task Parsing

## Goal
Parse Markdown task files into Go structs.

## Tasks

### 3.1 Define Task Struct (`internal/task/task.go`)

```go
type Task struct {
    ID          string       // Filename without .md extension
    Title       string       // From # heading
    Status      string       // From status = "..." line
    Description string       // Text between --- separators
    Checklist   []CheckItem  // Bullet points
    FilePath    string       // Full path to file
}

type CheckItem struct {
    Text     string
    Checked  bool
    Indent   int        // Nesting level (0, 1, 2...)
}
```

### 3.2 Implement Parser

Parse the markdown format:
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

Functions:
- `ParseTaskFile(path string) (*Task, error)`
- `parseTitle(line string) string`
- `parseStatus(line string) string`
- `parseChecklist(lines []string) []CheckItem`

### 3.3 Handle Parsing Errors

If a file cannot be parsed:
- Return error with filename and reason
- Caller will log warning and skip

## Testing This Step

Create a test file manually:
```bash
mkdir -p .tasks
cat > .tasks/1738764000-test-task.md << 'EOF'
# My Test Task
status = "TODO"
---
This is a test description
---
- [ ] First item
- [x] Completed item
    - [ ] Nested item
EOF
```

Add temporary code in main.go to:
1. Parse the file
2. Print the parsed struct

**Expected output**:
- Task struct with all fields populated correctly
- Nested checklist items have correct indent levels

## Files to Create
- `internal/task/task.go` - Task struct and parser

## Definition of Done
- [ ] Can parse valid task markdown files
- [ ] Title, status, description extracted correctly
- [ ] Checklist items parsed with correct checked state
- [ ] Nested items have correct indent level
- [ ] Invalid files return descriptive errors
