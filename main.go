package main

import (
	"log"
	"os"
	"path/filepath"

	"doings/internal/app"
	"doings/internal/config"
	"doings/internal/task"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize .tasks/ directory and config
	_, err := config.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
		os.Exit(1)
	}

	// Load config
	configPath := filepath.Join(config.TasksDir, config.ConfigFile)
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		os.Exit(1)
	}

	// Load tasks and collect warnings
	tasks, errs := task.ListTasks(config.TasksDir)
	var warnings []string

	// Add parse errors
	for _, err := range errs {
		warnings = append(warnings, err.Error())
	}

	// Check for tasks with invalid status
	validStatuses := make(map[string]bool)
	for _, col := range cfg.Board.Columns {
		validStatuses[col] = true
	}

	for _, t := range tasks {
		if !validStatuses[t.Status] {
			warnings = append(warnings,
				"Task \""+t.Title+"\" has unknown status \""+t.Status+"\"")
		}
	}

	// Create and run the Bubble Tea program with alternate screen
	p := tea.NewProgram(
		app.NewAppModel(cfg.Board.Columns, tasks, warnings),
		tea.WithAltScreen(), // Use alternate screen buffer
	)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
