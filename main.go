package main

import (
	"fmt"
	"log"
	"os"

	"doings/internal/config"
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

	fmt.Println("Ready!")
}
