package main

import (
	"log"
	"os"

	// Adjust import paths based on your actual module path
	"github.com/prasanna-1807/synapse-sprint/backend/internal/config"
	"github.com/prasanna-1807/synapse-sprint/backend/internal/database"
	"github.com/prasanna-1807/synapse-sprint/backend/internal/repository" // Import repository package
)

const (
	DatabaseName = "synapse_sprint_db" // Define database name
)

func main() {
	logger := log.New(os.Stdout, "SYNAPSE-SPRINT: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("Starting Synapse Sprint backend server...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	mongoClient, err := database.ConnectDB(cfg.MongoURI)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.DisconnectDB(mongoClient)

	// Get database handle
	db := mongoClient.Database(DatabaseName) // Use constant for DB name

	// Initialize repositories (DAL)
	userRepo := repository.NewMongoUserRepository(db)
	logger.Println("User repository initialized.")
	// We will use userRepo later when initializing services

	// TODO: Initialize services (Business Logic) - Pass userRepo here
	// TODO: Initialize API handlers/router - Services will be passed here
	// TODO: Start the HTTP server on port cfg.ServerPort

	logger.Println("Server setup complete (config loaded, DB connected, repos initialized). Exiting for now.")
}
