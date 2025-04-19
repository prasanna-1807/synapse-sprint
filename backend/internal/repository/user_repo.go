package repository

import (
	"context"
	"errors" // Import errors package
	"log"
	"os"
	"time"

	// Adjust import paths
	"github.com/prasanna-1807/synapse-sprint/backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (primitive.ObjectID, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error)
	// Add UpdatePassword, ListUsers, etc. as needed
}

// mongoUserRepository implements UserRepository for MongoDB
type mongoUserRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
	logger     *log.Logger
}

// NewMongoUserRepository creates a new MongoDB user repository instance
func NewMongoUserRepository(db *mongo.Database) UserRepository {
	// Create a logger specific to this repository
	logger := log.New(os.Stdout, "USER_REPO: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Get the 'users' collection
	collection := db.Collection("users")

	// Optional: Create indexes (e.g., unique index on username)
	// It's good practice to ensure indexes exist on startup
	_, err := collection.Indexes().CreateOne(
		context.Background(), // Use background context for startup tasks
		mongo.IndexModel{
			Keys:    bson.D{{Key: "username", Value: 1}}, // 1 for ascending index
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		// Log the error but don't necessarily fail startup,
		// maybe the index already exists or DB permissions are missing.
		// Production apps might handle this more robustly.
		logger.Printf("Warning: Failed to create unique index on username: %v", err)
	} else {
		logger.Println("Unique index on username ensured.")
	}

	return &mongoUserRepository{
		db:         db,
		collection: collection,
		logger:     logger,
	}
}

// Create inserts a new user into the database
func (r *mongoUserRepository) Create(ctx context.Context, user *domain.User) (primitive.ObjectID, error) {
	r.logger.Printf("Attempting to create user: %s", user.Username)

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.ID = primitive.NewObjectID() // Generate ID before insert (optional but can be useful)

	// Insert document
	res, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		// Check for duplicate key error (MongoDB error code 11000)
		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) {
			for _, we := range writeErr.WriteErrors {
				if we.Code == 11000 {
					r.logger.Printf("Username '%s' already exists.", user.Username)
					return primitive.NilObjectID, errors.New("username already exists") // Return a specific error
				}
			}
		}
		// Generic error
		r.logger.Printf("Failed to insert user %s: %v", user.Username, err)
		return primitive.NilObjectID, err
	}

	// Check if InsertOne succeeded and returned an ID
	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		r.logger.Printf("InsertOne result ID is not an ObjectID for user %s", user.Username)
		return primitive.NilObjectID, errors.New("failed to get inserted ID")
	}

	r.logger.Printf("Successfully created user %s with ID %s", user.Username, insertedID.Hex())
	// return insertedID, nil // Use ID from result
	return user.ID, nil // Or return the ID we generated
}

// FindByUsername finds a user by their username
func (r *mongoUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	r.logger.Printf("Attempting to find user by username: %s", username)
	var user domain.User
	// bson.M is a shorthand for map[string]interface{} - useful for simple queries
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Printf("User not found: %s", username)
			return nil, errors.New("user not found") // Return specific, clear error
		}
		r.logger.Printf("Error finding user %s: %v", username, err)
		return nil, err // Return original error for other DB issues
	}
	r.logger.Printf("Found user: %s (ID: %s)", user.Username, user.ID.Hex())
	return &user, nil
}

// FindByID finds a user by their MongoDB ObjectID
func (r *mongoUserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	r.logger.Printf("Attempting to find user by ID: %s", id.Hex())
	var user domain.User
	// Use bson.D for ordered elements if needed, bson.M is fine here.
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Printf("User not found by ID: %s", id.Hex())
			return nil, errors.New("user not found")
		}
		r.logger.Printf("Error finding user by ID %s: %v", id.Hex(), err)
		return nil, err
	}
	r.logger.Printf("Found user by ID: %s (Username: %s)", id.Hex(), user.Username)
	return &user, nil
}

// --- Implement other UserRepository methods here ---
