package config

import (
	"log"
	"os"
	"strconv"
)

// Config holds the application configuration values
type Config struct {
	ServerPort string
	MongoURI   string
	// Add other configuration fields as needed, e.g., JWT secret, log level
}

// LoadConfig loads configuration from environment variables
// It provides default values for some settings if environment variables are not set.
func LoadConfig() (*Config, error) {
	logger := log.New(os.Stdout, "CONFIG: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("Loading configuration...")

	cfg := &Config{
		// Default values
		ServerPort: "8080",                      // Default port
		MongoURI:   "mongodb://localhost:27017", // Default local MongoDB URI
	}

	// Read Server Port from environment variable
	if portEnv, exists := os.LookupEnv("SERVER_PORT"); exists {
		// Basic validation: check if it's a number (more robust checks could be added)
		if _, err := strconv.Atoi(portEnv); err == nil {
			cfg.ServerPort = portEnv
		} else {
			logger.Printf("Warning: SERVER_PORT environment variable ('%s') is not a valid number. Using default '%s'.\n", portEnv, cfg.ServerPort)
		}
	}

	// Read MongoDB URI from environment variable
	if mongoURIEnv, exists := os.LookupEnv("MONGODB_URI"); exists {
		cfg.MongoURI = mongoURIEnv
	}

	logger.Println("Configuration loaded successfully.")
	logger.Printf("Server Port: %s", cfg.ServerPort)
	logger.Printf("MongoDB URI: %s", cfg.MongoURI) // Be cautious logging sensitive data like full URIs in production

	return cfg, nil // Returning nil error for simplicity now; add validation later
}
