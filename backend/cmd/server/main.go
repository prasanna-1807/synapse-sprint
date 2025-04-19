package main

import (
	"log"
	"os"

	// Import the internal config package
	// Use the module path you defined in go.mod
	"github.com/prasanna-1807/synapse-sprint/backend/internal/config"
)

func main() {
	logger := log.New(os.Stdout, "SYNAPSE-SPRINT: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("Starting Synapse Sprint backend server...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err) // Use Fatalf to exit if config fails
	}

	// Use cfg.ServerPort and cfg.MongoURI below...

	// TODO: Initialize database connection (MongoDB) using cfg.MongoURI
	// TODO: Initialize repositories (DAL)
	// TODO: Initialize services (Business Logic)
	// TODO: Initialize API handlers/router
	// TODO: Start the HTTP server on port cfg.ServerPort

	logger.Println("Server setup placeholder complete (config loaded). Exiting for now.")
}
