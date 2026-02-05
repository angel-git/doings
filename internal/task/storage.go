package task

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CreateTask creates a new task file with the given title and status
func CreateTask(tasksDir string, title string, status string) (*Task, error) {
	id := generateTaskID(title)
	filePath := filepath.Join(tasksDir, id+".md")

	task := &Task{
		ID:          id,
		Title:       title,
		Status:      status,
		Description: "",
		Checklist:   []CheckItem{},
		FilePath:    filePath,
	}

	if err := SaveTask(task); err != nil {
		return nil, err
	}

	return task, nil
}

// ListTasks reads all task files from the given directory
func ListTasks(tasksDir string) ([]*Task, []error) {
	var tasks []*Task
	var errors []error

	entries, err := os.ReadDir(tasksDir)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to read tasks directory: %w", err)}
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(tasksDir, entry.Name())
		task, err := ParseTaskFile(filePath)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to parse %s: %w", entry.Name(), err))
			continue
		}

		tasks = append(tasks, task)
	}

	return tasks, errors
}

// SaveTask writes the task to its file
func SaveTask(task *Task) error {
	content := task.ToMarkdown()
	if err := os.WriteFile(task.FilePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}
	return nil
}

// DeleteTask removes the task file from the filesystem
func DeleteTask(task *Task) error {
	if err := os.Remove(task.FilePath); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

// generateTaskID creates a unique ID based on timestamp and slugified title
func generateTaskID(title string) string {
	timestamp := time.Now().Unix()
	slug := slugify(title)
	return fmt.Sprintf("%d-%s", timestamp, slug)
}

// slugify converts a string to a URL-friendly slug
func slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)
	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	// Remove non-alphanumeric characters except hyphens
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	// Remove consecutive hyphens
	slug := result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")
	return slug
}
