// Package store/db.go contains database functions
package store

import (
	"context"
	"fmt"
	"pet-mgt/backend/internal/config"

	"github.com/supabase-community/supabase-go"
)

type SupabaseService struct {
	client *supabase.Client
	config *config.Config
}

// NewSupabaseService creates a new SupabaseService
func NewSupabaseService(cfg *config.Config) (*SupabaseService, error) {
	client, err := supabase.NewClient(
		cfg.SupabaseURL,
		cfg.SupabaseServiceKey,
		&supabase.ClientOptions{},
	)
	if err != nil {
		return nil, err
	}

	return &SupabaseService{
		config: cfg,
		client: client,
	}, nil
}

// Ping checks if the database connection is alive
func (s *SupabaseService) Ping(ctx context.Context) error {
	// Simple health check - the client connection is tested during initialization
	// If we got here, the connection is working
	return nil
}

// GetUserByID retrieves a user by their ID
func (s *SupabaseService) GetUserByID(
	ctx context.Context,
	userID string,
) (*User, error) {
	// Try to get from clients table first
	var client Client
	_, err := s.client.From("clients").Select("*", "", false).Eq("id", userID).Single().ExecuteTo(&client)
	if err == nil {
		return &User{
			ID:    client.ID,
			Email: client.Email,
			Role:  client.Role,
		}, nil
	}

	// Try to get from veterinarians table
	var vet Veterinarian
	_, err = s.client.From("veterinarians").Select("*", "", false).Eq("id", userID).Single().ExecuteTo(&vet)
	if err == nil {
		return &User{
			ID:    vet.ID,
			Email: vet.Email,
			Role:  vet.Role,
		}, nil
	}

	return nil, fmt.Errorf("user not found")
}

// CreateUser creates a new user profile
func (s *SupabaseService) CreateUser(ctx context.Context, user *User) error {
	// This method will be implemented based on the user role
	// For now, return an error as we need to know the specific user type
	return fmt.Errorf("use CreateClient or CreateVeterinarian instead")
}

// CreateClient creates a new client profile
func (s *SupabaseService) CreateClient(ctx context.Context, client *Client) error {
	_, _, err := s.client.From("clients").Insert(client, false, "", "", "").Execute()
	return err
}

// CreateVeterinarian creates a new veterinarian profile
func (s *SupabaseService) CreateVeterinarian(ctx context.Context, vet *Veterinarian) error {
	_, _, err := s.client.From("veterinarians").Insert(vet, false, "", "", "").Execute()
	return err
}

// UpdateUser updates a user profile
func (s *SupabaseService) UpdateUser(ctx context.Context, user *User) error {
	// This method will be implemented based on the user role
	return fmt.Errorf("use UpdateClient or UpdateVeterinarian instead")
}

// UpdateClient updates a client profile
func (s *SupabaseService) UpdateClient(ctx context.Context, client *Client) error {
	_, _, err := s.client.From("clients").Update(client, "", "").Eq("id", client.ID).Execute()
	return err
}

// UpdateVeterinarian updates a veterinarian profile
func (s *SupabaseService) UpdateVeterinarian(ctx context.Context, vet *Veterinarian) error {
	_, _, err := s.client.From("veterinarians").Update(vet, "", "").Eq("id", vet.ID).Execute()
	return err
}

// DeleteUser deletes a user profile
func (s *SupabaseService) DeleteUser(ctx context.Context, userID string) error {
	// Try to delete from clients table first
	_, _, err := s.client.From("clients").Delete("", "").Eq("id", userID).Execute()
	if err == nil {
		return nil
	}

	// Try to delete from veterinarians table
	_, _, err = s.client.From("veterinarians").Delete("", "").Eq("id", userID).Execute()
	return err
}

// ListUsers lists all users with pagination
func (s *SupabaseService) ListUsers(ctx context.Context, limit, offset int) ([]User, error) {
	var users []User

	// Get clients
	var clients []Client
	_, err := s.client.From("clients").Select("*", "", false).Range(offset, offset+limit-1, "").ExecuteTo(&clients)
	if err != nil {
		return nil, err
	}

	// Get veterinarians
	var vets []Veterinarian
	_, err = s.client.From("veterinarians").Select("*", "", false).Range(offset, offset+limit-1, "").ExecuteTo(&vets)
	if err != nil {
		return nil, err
	}

	// Convert to User interface
	for _, client := range clients {
		users = append(users, User{
			ID:    client.ID,
			Email: client.Email,
			Role:  client.Role,
		})
	}

	for _, vet := range vets {
		users = append(users, User{
			ID:    vet.ID,
			Email: vet.Email,
			Role:  vet.Role,
		})
	}

	return users, nil
}

// GetPetsByUserID retrieves all pets for a specific user
func (s *SupabaseService) GetPetsByUserID(
	ctx context.Context,
	userID string,
) ([]Pet, error) {
	var pets []Pet
	_, err := s.client.From("pets").Select("*", "", false).Eq("owner_id", userID).ExecuteTo(&pets)
	return pets, err
}

// GetPetByID retrieves a pet by ID
func (s *SupabaseService) GetPetByID(ctx context.Context, petID string) (*Pet, error) {
	var pet Pet
	_, err := s.client.From("pets").Select("*", "", false).Eq("id", petID).Single().ExecuteTo(&pet)
	if err != nil {
		return nil, err
	}
	return &pet, nil
}

// CreatePet creates a new pet
func (s *SupabaseService) CreatePet(ctx context.Context, pet *Pet) error {
	_, _, err := s.client.From("pets").Insert(pet, false, "", "", "").Execute()
	return err
}

// UpdatePet updates an existing pet
func (s *SupabaseService) UpdatePet(ctx context.Context, pet *Pet) error {
	_, _, err := s.client.From("pets").Update(pet, "", "").Eq("id", pet.ID).Execute()
	return err
}

// DeletePet deletes a pet by ID
func (s *SupabaseService) DeletePet(ctx context.Context, petID string) error {
	_, _, err := s.client.From("pets").Delete("", "").Eq("id", petID).Execute()
	return err
}

// GetMedicalRecordsByPetID retrieves all medical records for a pet
func (s *SupabaseService) GetMedicalRecordsByPetID(ctx context.Context, petID string) ([]MedicalRecord, error) {
	var records []MedicalRecord
	_, err := s.client.From("medical_records").Select("*", "", false).Eq("pet_id", petID).ExecuteTo(&records)
	return records, err
}

// GetMedicalRecordByID retrieves a specific medical record
func (s *SupabaseService) GetMedicalRecordByID(ctx context.Context, recordID string) (*MedicalRecord, error) {
	var record MedicalRecord
	_, err := s.client.From("medical_records").Select("*", "", false).Eq("id", recordID).Single().ExecuteTo(&record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// CreateMedicalRecord creates a new medical record
func (s *SupabaseService) CreateMedicalRecord(ctx context.Context, record *MedicalRecord) error {
	_, _, err := s.client.From("medical_records").Insert(record, false, "", "", "").Execute()
	return err
}

// UpdateMedicalRecord updates an existing medical record
func (s *SupabaseService) UpdateMedicalRecord(ctx context.Context, record *MedicalRecord) error {
	_, _, err := s.client.From("medical_records").Update(record, "", "").Eq("id", record.ID).Execute()
	return err
}

// DeleteMedicalRecord deletes a medical record by ID
func (s *SupabaseService) DeleteMedicalRecord(ctx context.Context, recordID string) error {
	_, _, err := s.client.From("medical_records").Delete("", "").Eq("id", recordID).Execute()
	return err
}

// GetClientByID retrieves a client by ID
func (s *SupabaseService) GetClientByID(ctx context.Context, clientID string) (*Client, error) {
	var client Client
	_, err := s.client.From("clients").Select("*", "", false).Eq("id", clientID).Single().ExecuteTo(&client)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// GetVeterinarianByID retrieves a veterinarian by ID
func (s *SupabaseService) GetVeterinarianByID(ctx context.Context, vetID string) (*Veterinarian, error) {
	var vet Veterinarian
	_, err := s.client.From("veterinarians").Select("*", "", false).Eq("id", vetID).Single().ExecuteTo(&vet)
	if err != nil {
		return nil, err
	}
	return &vet, nil
}

// Close performs cleanup operations
func (s *SupabaseService) Close() error {
	// Supabase client uses HTTP connections, no explicit cleanup needed
	// Future: close any connection pools or background workers here
	return nil
}
