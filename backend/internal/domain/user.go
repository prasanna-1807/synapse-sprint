package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive" // MongoDB driver types
)

// Role type for user roles
type Role string

const (
	RoleStudent Role = "student"
	RoleParent  Role = "parent" // Or Teacher
	RoleAdmin   Role = "admin"
)

// User represents a user in the system
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"` // MongoDB default ID field
	Username     string             `bson:"username"`
	PasswordHash string             `bson:"passwordHash"` // Store hashed passwords only!
	Role         Role               `bson:"role"`
	CreatedAt    time.Time          `bson:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"`
	// Add other fields if needed, e.g., FirstName, LastName for Parent/Admin
}
