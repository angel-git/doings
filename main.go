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

	// Test CRUD operations
	fmt.Println("\n--- Testing CRUD Operations ---")

	// 1. Create a new task
	fmt.Println("\n1. Creating new task...")
	newTask, err := task.CreateTask(config.TasksDir, "Test CRUD Task", "TODO")
	if err != nil {
		log.Printf("Error creating task: %v", err)
	} else {
		fmt.Printf("   Created: %s (ID: %s)\n", newTask.Title, newTask.ID)
	}

	// 2. List all tasks
	fmt.Println("\n2. Listing all tasks...")
	tasks, errs := task.ListTasks(config.TasksDir)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Printf("   Warning: %v", err)
		}
	}
	fmt.Printf("   Found %d task(s):\n", len(tasks))
	for _, t := range tasks {
		fmt.Printf("   - %s [%s]\n", t.Title, t.Status)
	}

	// 3. Update a task
	if len(tasks) > 0 {
		fmt.Println("\n3. Updating first task...")
		firstTask := tasks[0]
		firstTask.Status = "DOING"
		firstTask.Description = "Updated description via CRUD test"
		firstTask.Checklist = append(firstTask.Checklist, task.CheckItem{
			Text:    "New checklist item",
			Checked: false,
			Indent:  0,
		})
		if err := task.SaveTask(firstTask); err != nil {
			log.Printf("   Error saving task: %v", err)
		} else {
			fmt.Printf("   Updated: %s\n", firstTask.Title)
		}
	}

	// 4. Delete the newly created task
	if newTask != nil {
		fmt.Println("\n4. Deleting newly created task...")
		if err := task.DeleteTask(newTask); err != nil {
			log.Printf("   Error deleting task: %v", err)
		} else {
			fmt.Printf("   Deleted: %s\n", newTask.Title)
		}
	}

	// List tasks again to verify deletion
	fmt.Println("\n5. Listing tasks after deletion...")
	tasks, _ = task.ListTasks(config.TasksDir)
	fmt.Printf("   Found %d task(s):\n", len(tasks))
	for _, t := range tasks {
		fmt.Printf("   - %s [%s]\n", t.Title, t.Status)
	}

	fmt.Println("\nReady!")
}
