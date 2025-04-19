package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ConnectDB establishes a connection to MongoDB
func ConnectDB(mongoURI string) (*mongo.Client, error) {
	logger := log.New(os.Stdout, "DATABASE: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("Connecting to MongoDB...")

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Use context with timeout for connection attempt
	// context.TODO() is often used as a temporary context when unsure;
	// a more specific context (e.g., with cancelation) might be used in production handlers.
	// For startup, a timeout context is good practice.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure context resources are released

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Printf("Failed to create client: %v\n", err)
		return nil, err
	}

	// Ping the primary server to verify the connection
	// Use another short timeout context for the ping
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	err = client.Ping(pingCtx, readpref.Primary())
	if err != nil {
		logger.Printf("Failed to connect to MongoDB (ping failed): %v\n", err)
		// Disconnect if ping fails after successful Connect call
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer disconnectCancel()
		_ = client.Disconnect(disconnectCtx) // Attempt cleanup, ignore error for now
		return nil, err
	}

	logger.Println("Successfully connected to MongoDB!")
	return client, nil
}

// DisconnectDB disconnects the MongoDB client
// We'll use this later for graceful shutdown
func DisconnectDB(client *mongo.Client) {
	logger := log.New(os.Stdout, "DATABASE: ", log.Ldate|log.Ltime|log.Lshortfile)
	if client == nil {
		return
	}
	logger.Println("Disconnecting from MongoDB...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		logger.Printf("Error disconnecting from MongoDB: %v\n", err)
	} else {
		logger.Println("Successfully disconnected from MongoDB.")
	}
}
