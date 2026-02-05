package task

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Task represents a single task from a Markdown file
type Task struct {
	ID          string      // Filename without .md extension
	Title       string      // From # heading
	Status      string      // From status = "..." line
	Description string      // Text between --- separators
	Checklist   []CheckItem // Bullet points
	FilePath    string      // Full path to file
}

// CheckItem represents a checkbox item in a task
type CheckItem struct {
	Text    string
	Checked bool
	Indent  int // Nesting level (0, 1, 2...)
}

// ParseTaskFile reads and parses a task markdown file
func ParseTaskFile(path string) (*Task, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	task := &Task{
		FilePath: path,
		ID:       strings.TrimSuffix(filepath.Base(path), ".md"),
	}

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if len(lines) == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	// Parse the file in sections
	state := "title"
	var descLines []string
	var checklistLines []string
	descStarted := false

	for i, line := range lines {
		switch state {
		case "title":
			if strings.HasPrefix(line, "# ") {
				task.Title = parseTitle(line)
				state = "status"
			} else if line != "" {
				return nil, fmt.Errorf("expected title on line %d, got: %s", i+1, line)
			}

		case "status":
			if strings.HasPrefix(line, "status = ") {
				task.Status = parseStatus(line)
				state = "description"
			} else if line != "" {
				return nil, fmt.Errorf("expected status on line %d, got: %s", i+1, line)
			}

		case "description":
			if line == "---" {
				if !descStarted {
					descStarted = true
				} else {
					// End of description
					task.Description = strings.TrimSpace(strings.Join(descLines, "\n"))
					state = "checklist"
				}
			} else if descStarted {
				descLines = append(descLines, line)
			}

		case "checklist":
			if strings.TrimSpace(line) != "" {
				checklistLines = append(checklistLines, line)
			}
		}
	}

	// Parse checklist items
	task.Checklist = parseChecklist(checklistLines)

	return task, nil
}

// parseTitle extracts the title from a markdown heading
func parseTitle(line string) string {
	return strings.TrimSpace(strings.TrimPrefix(line, "#"))
}

// parseStatus extracts the status value from status = "VALUE" line
func parseStatus(line string) string {
	// Remove 'status = ' prefix and trim quotes
	value := strings.TrimPrefix(line, "status = ")
	value = strings.Trim(value, `"`)
	return value
}

// parseChecklist converts markdown checklist lines into CheckItem structs
func parseChecklist(lines []string) []CheckItem {
	var items []CheckItem

	for _, line := range lines {
		if !strings.Contains(line, "[ ]") && !strings.Contains(line, "[x]") {
			continue
		}

		// Calculate indent level (number of leading spaces / 4)
		indent := 0
		trimmed := strings.TrimLeft(line, " ")
		leadingSpaces := len(line) - len(trimmed)
		indent = leadingSpaces / 4

		// Check if item is checked
		checked := strings.Contains(line, "[x]")

		// Extract text after the checkbox
		text := trimmed
		if strings.HasPrefix(text, "- [ ] ") {
			text = strings.TrimPrefix(text, "- [ ] ")
		} else if strings.HasPrefix(text, "- [x] ") {
			text = strings.TrimPrefix(text, "- [x] ")
		}

		items = append(items, CheckItem{
			Text:    text,
			Checked: checked,
			Indent:  indent,
		})
	}

	return items
}

// ToMarkdown converts the Task back to markdown format
func (t *Task) ToMarkdown() string {
	var b strings.Builder

	// Title
	b.WriteString("# ")
	b.WriteString(t.Title)
	b.WriteString("\n")

	// Status
	b.WriteString("status = \"")
	b.WriteString(t.Status)
	b.WriteString("\"\n")

	// Description
	b.WriteString("---\n")
	if t.Description != "" {
		b.WriteString(t.Description)
		b.WriteString("\n")
	}
	b.WriteString("---\n")

	// Checklist
	for _, item := range t.Checklist {
		// Add indentation
		for i := 0; i < item.Indent; i++ {
			b.WriteString("    ")
		}

		// Add checkbox
		if item.Checked {
			b.WriteString("- [x] ")
		} else {
			b.WriteString("- [ ] ")
		}

		// Add text
		b.WriteString(item.Text)
		b.WriteString("\n")
	}

	return b.String()
}
