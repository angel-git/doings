# Step 04: Task CRUD Operations

## Goal
Implement file operations for creating, reading, updating, and deleting tasks.

## Tasks

### 4.1 Create Storage Package (`internal/task/storage.go`)

Functions:

#### List Tasks
```go
func ListTasks(tasksDir string) ([]*Task, []error)
```
- Read all `.md` files in `.tasks/`
- Parse each file
- Return valid tasks and list of errors (for warnings)

#### Create Task
```go
func CreateTask(tasksDir string, title string, status string) (*Task, error)
```
- Generate ID: `{timestamp}-{slugified-title}.md`
- Create file with template content
- Return parsed task

#### Update Task
```go
func SaveTask(task *Task) error
```
- Serialize task back to markdown format
- Write to file (overwrite)

#### Delete Task
```go
func DeleteTask(task *Task) error
```
- Remove the file from filesystem

### 4.2 Implement ID Generation

```go
func generateTaskID(title string) string {
    timestamp := time.Now().Unix()
    slug := slugify(title) // lowercase, replace spaces with hyphens
    return fmt.Sprintf("%d-%s", timestamp, slug)
}
```

### 4.3 Implement Task Serialization

Convert Task struct back to markdown:
```go
func (t *Task) ToMarkdown() string
```

## Testing This Step

```bash
./doings
```

Add temporary code to:
1. Create a new task
2. List all tasks
3. Update a task's status
4. Delete the task

Verify files are created/modified/deleted correctly.

## Files to Create/Modify
- `internal/task/storage.go` - File operations
- `internal/task/task.go` - Add ToMarkdown method

## Definition of Done
- [x] Can create new task files with correct format
- [x] Can list all tasks from `.tasks/` directory
- [x] Can update and save task changes
- [x] Can delete task files
- [x] ID generation uses timestamp + slugified title
- [x] Malformed files are reported but don't crash

## Status
**COMPLETED** - All tasks finished successfully on Thu Feb 05 2026
