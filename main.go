package main

import (
	"log"
	"os"
	"path/filepath"

	"doings/internal/config"
	"doings/internal/task"
	"doings/internal/ui"
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

	// Load tasks
	tasks, errs := task.ListTasks(config.TasksDir)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Printf("Warning: %v", err)
		}
	}

	// Create and run the Bubble Tea program
	p := tea.NewProgram(ui.NewBoardModel(cfg.Board.Columns, tasks))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
