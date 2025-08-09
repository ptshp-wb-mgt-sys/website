// Package handlers/pets.go contains pet management HTTP handlers
package handlers

import (
	"encoding/json"
	"net/http"
	"pet-mgt/backend/internal/middleware"
	"pet-mgt/backend/internal/store"
	"time"

	"github.com/go-chi/chi/v5"
)

// PetHandler handles pet-related HTTP requests
type PetHandler struct {
	db store.Database
}

// NewPetHandler creates a new PetHandler with database dependency
func NewPetHandler(db store.Database) *PetHandler {
	return &PetHandler{
		db: db,
	}
}

// CreatePet creates a new pet
func (h *PetHandler) CreatePet(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	var req struct {
		Name        string    `json:"name"`
		Type        string    `json:"type"`
		Breed       string    `json:"breed"`
		DateOfBirth time.Time `json:"date_of_birth"`
		Weight      float64   `json:"weight"`
		OwnerID     string    `json:"owner_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Authorization check: clients can only create pets for themselves, admins can create for anyone
	if user.Role == "client" && req.OwnerID != user.Sub {
		ErrorResponse(
			w,
			http.StatusForbidden,
			"Clients can only create pets for themselves",
		)
		return
	}

	if user.Role == "veterinarian" {
		ErrorResponse(w, http.StatusForbidden, "Veterinarians cannot create pets")
		return
	}

	// If no owner_id specified, use the authenticated user's ID
	if req.OwnerID == "" {
		req.OwnerID = user.Sub
	}

	pet := store.NewPet(
		req.OwnerID,
		req.Name,
		req.Type,
		req.Breed,
		req.DateOfBirth,
		req.Weight,
	)

	if err := h.db.CreatePet(r.Context(), pet); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to create pet")
		return
	}

	SuccessResponse(w, pet)
}

// GetPet retrieves a pet by ID
func (h *PetHandler) GetPet(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	petID := chi.URLParam(r, "id")
	if petID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Pet ID is required")
		return
	}

	pet, err := h.db.GetPetByID(r.Context(), petID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Pet not found")
		return
	}

	// Authorization check: clients can only access their own pets, vets and admins can access any
	if user.Role == "client" && pet.OwnerID != user.Sub {
		ErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	SuccessResponse(w, pet)
}

// UpdatePet updates a pet
func (h *PetHandler) UpdatePet(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	petID := chi.URLParam(r, "id")
	if petID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Pet ID is required")
		return
	}

	pet, err := h.db.GetPetByID(r.Context(), petID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Pet not found")
		return
	}

	// Authorization check: clients can only update their own pets, admins can update any
	if user.Role == "client" && pet.OwnerID != user.Sub {
		ErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	if user.Role == "veterinarian" {
		ErrorResponse(
			w,
			http.StatusForbidden,
			"Veterinarians cannot update pet details",
		)
		return
	}

	var req struct {
		Name        string    `json:"name"`
		Type        string    `json:"type"`
		Breed       string    `json:"breed"`
		DateOfBirth time.Time `json:"date_of_birth"`
		Weight      float64   `json:"weight"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update pet fields
	pet.Name = req.Name
	pet.Type = req.Type
	pet.Breed = req.Breed
	pet.DateOfBirth = req.DateOfBirth.Format("2006-01-02")
	pet.Weight = req.Weight
	pet.UpdatedAt = time.Now()

	if err := h.db.UpdatePet(r.Context(), pet); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update pet")
		return
	}

	SuccessResponse(w, pet)
}

// DeletePet deletes a pet
func (h *PetHandler) DeletePet(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	petID := chi.URLParam(r, "id")
	if petID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Pet ID is required")
		return
	}

	pet, err := h.db.GetPetByID(r.Context(), petID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Pet not found")
		return
	}

	// Authorization check: clients can only delete their own pets, admins can delete any
	if user.Role == "client" && pet.OwnerID != user.Sub {
		ErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	if user.Role == "veterinarian" {
		ErrorResponse(w, http.StatusForbidden, "Veterinarians cannot delete pets")
		return
	}

	if err := h.db.DeletePet(r.Context(), petID); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to delete pet")
		return
	}

	MessageResponse(w, http.StatusOK, "Pet deleted successfully")
}

// GetPetsByClient retrieves all pets for a specific client
func (h *PetHandler) GetPetsByClient(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	clientID := chi.URLParam(r, "clientId")
	if clientID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Client ID is required")
		return
	}

	// Authorization check: clients can only access their own pets, vets and admins can access any
	if user.Role == "client" && clientID != user.Sub {
		ErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	pets, err := h.db.GetPetsByUserID(r.Context(), clientID)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve pets")
		return
	}

	SuccessResponse(w, pets)
}
