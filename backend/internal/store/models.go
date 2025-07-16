// Package store/models.go contains the models for the application
package store

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Database interface defines methods for data access operations
type Database interface {
	// User operations
	GetUserByID(ctx context.Context, userID string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, userID string) error
	ListUsers(ctx context.Context, limit, offset int) ([]User, error)

	// Client/Veterinarian specific
	CreateClient(ctx context.Context, client *Client) error
	CreateVeterinarian(ctx context.Context, vet *Veterinarian) error
	UpdateClient(ctx context.Context, client *Client) error
	UpdateVeterinarian(ctx context.Context, vet *Veterinarian) error
	GetClientByID(ctx context.Context, clientID string) (*Client, error)
	GetVeterinarianByID(ctx context.Context, vetID string) (*Veterinarian, error)

	// Pet operations
	GetPetsByUserID(ctx context.Context, userID string) ([]Pet, error)
	GetPetByID(ctx context.Context, petID string) (*Pet, error)
	CreatePet(ctx context.Context, pet *Pet) error
	UpdatePet(ctx context.Context, pet *Pet) error
	DeletePet(ctx context.Context, petID string) error

	// Medical record operations
	GetMedicalRecordsByPetID(ctx context.Context, petID string) ([]MedicalRecord, error)
	GetMedicalRecordByID(ctx context.Context, recordID string) (*MedicalRecord, error)
	CreateMedicalRecord(ctx context.Context, record *MedicalRecord) error
	UpdateMedicalRecord(ctx context.Context, record *MedicalRecord) error
	DeleteMedicalRecord(ctx context.Context, recordID string) error

	// Health check
	Ping(ctx context.Context) error

	// Cleanup
	Close() error
}

// User represents a user in the system (base interface for Client and Veterinarian)
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Client represents a pet owner
type Client struct {
	ID      string `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Email   string `json:"email" db:"email"`
	Phone   string `json:"phone" db:"phone"`
	Address string `json:"address" db:"address"`
	Role    string `json:"role" db:"role"`
}

// Veterinarian represents a veterinarian user
type Veterinarian struct {
	ID             string         `json:"id" db:"id"`
	Name           string         `json:"name" db:"name"`
	Email          string         `json:"email" db:"email"`
	Phone          string         `json:"phone" db:"phone"`
	ClinicAddress  string         `json:"clinic_address" db:"clinic_address"`
	AvailableHours []WorkingHours `json:"available_hours" db:"available_hours"`
	Role           string         `json:"role" db:"role"`
}

// WorkingHours represents available working hours for a veterinarian
type WorkingHours struct {
	DayOfWeek string `json:"day_of_week" db:"day_of_week"`
	Start     string `json:"start" db:"start"`
	End       string `json:"end" db:"end"`
}

// Pet represents a pet in the system
type Pet struct {
	ID          string    `json:"id" db:"id"`
	OwnerID     string    `json:"owner_id" db:"owner_id"`
	Name        string    `json:"name" db:"name"`
	Type        string    `json:"type" db:"type"`
	Breed       string    `json:"breed" db:"breed"`
	DateOfBirth time.Time `json:"date_of_birth" db:"date_of_birth"`
	Weight      float64   `json:"weight" db:"weight"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// MedicalRecord represents a veterinary visit record
type MedicalRecord struct {
	ID                   string    `json:"id" db:"id"`
	PetID                string    `json:"pet_id" db:"pet_id"`
	VeterinarianID       string    `json:"veterinarian_id" db:"veterinarian_id"`
	DateOfVisit          time.Time `json:"date_of_visit" db:"date_of_visit"`
	ReasonForVisit       string    `json:"reason_for_visit" db:"reason_for_visit"`
	Diagnosis            string    `json:"diagnosis" db:"diagnosis"`
	MedicationPrescribed []string  `json:"medication_prescribed" db:"medication_prescribed"`
	Notes                string    `json:"notes" db:"notes"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// NewPet creates a new Pet with generated ID and timestamps
func NewPet(ownerID, name, petType, breed string, dateOfBirth time.Time, weight float64) *Pet {
	now := time.Now()
	return &Pet{
		ID:          uuid.New().String(),
		OwnerID:     ownerID,
		Name:        name,
		Type:        petType,
		Breed:       breed,
		DateOfBirth: dateOfBirth,
		Weight:      weight,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewMedicalRecord creates a new MedicalRecord with generated ID and timestamps
func NewMedicalRecord(petID, veterinarianID, reasonForVisit, diagnosis string, dateOfVisit time.Time, medicationPrescribed []string, notes string) *MedicalRecord {
	now := time.Now()
	return &MedicalRecord{
		ID:                   uuid.New().String(),
		PetID:                petID,
		VeterinarianID:       veterinarianID,
		DateOfVisit:          dateOfVisit,
		ReasonForVisit:       reasonForVisit,
		Diagnosis:            diagnosis,
		MedicationPrescribed: medicationPrescribed,
		Notes:                notes,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}
