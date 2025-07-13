// Package store/models.go contains the models for the application
package store

import (
	"context"
)

// Database interface defines methods for data access operations
type Database interface {
	// User operations
	GetUserByID(ctx context.Context, userID string) (*User, error)

	// Pet operations (for future implementation)
	GetPetsByUserID(ctx context.Context, userID string) ([]Pet, error)
	CreatePet(ctx context.Context, pet *Pet) error
	UpdatePet(ctx context.Context, pet *Pet) error
	DeletePet(ctx context.Context, petID string) error

	// Health check
	Ping(ctx context.Context) error

	// Cleanup
	Close() error
}

// User represents a user in the system
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Pet represents a pet in the system
type Pet struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Species     string `json:"species"`
	Breed       string `json:"breed"`
	Age         int    `json:"age"`
	Description string `json:"description"`
}
