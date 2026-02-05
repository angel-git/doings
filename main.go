package main

import (
	"fmt"
	"log"
	"os"

	"doings/internal/config"
	"doings/internal/task"
)

func main() {
	fmt.Println("Terminal Task Board")

	// Initialize .tasks/ directory and config
	created, err := config.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
		os.Exit(1)
	}

	if created {
		fmt.Printf("Created default config in %s/%s\n", config.TasksDir, config.ConfigFile)
	}

	// Temporary test: Parse the test task file
	fmt.Println("\n--- Testing Task Parser ---")
	testTask, err := task.ParseTaskFile(".tasks/1738764000-test-task.md")
	if err != nil {
		log.Printf("Error parsing task: %v", err)
	} else {
		fmt.Printf("ID: %s\n", testTask.ID)
		fmt.Printf("Title: %s\n", testTask.Title)
		fmt.Printf("Status: %s\n", testTask.Status)
		fmt.Printf("Description: %s\n", testTask.Description)
		fmt.Printf("Checklist items:\n")
		for _, item := range testTask.Checklist {
			indent := ""
			for i := 0; i < item.Indent; i++ {
				indent += "  "
			}
			checked := "[ ]"
			if item.Checked {
				checked = "[x]"
			}
			fmt.Printf("  %s%s %s\n", indent, checked, item.Text)
		}
	}

	fmt.Println("\nReady!")
}
