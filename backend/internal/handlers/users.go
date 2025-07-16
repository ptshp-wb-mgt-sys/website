// Package handlers/users.go contains user management HTTP handlers
package handlers

import (
	"encoding/json"
	"net/http"
	"pet-mgt/backend/internal/middleware"
	"pet-mgt/backend/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// CreateUser creates a new user profile
func CreateUser(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
		Role    string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Only admins can create any user type, others can only create their own profile
	if user.Role != "admin" && req.Role != user.Role {
		ErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	db, ok := middleware.GetDBFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusInternalServerError, "Database not found in context")
		return
	}

	switch req.Role {
	case "client":
		client := &store.Client{
			ID:      user.Sub,
			Name:    req.Name,
			Email:   req.Email,
			Phone:   req.Phone,
			Address: req.Address,
			Role:    req.Role,
		}
		if err := db.CreateClient(r.Context(), client); err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to create client")
			return
		}
		SuccessResponse(w, client)

	case "veterinarian":
		vet := &store.Veterinarian{
			ID:    user.Sub,
			Name:  req.Name,
			Email: req.Email,
			Phone: req.Phone,
			Role:  req.Role,
		}
		if err := db.CreateVeterinarian(r.Context(), vet); err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to create veterinarian")
			return
		}
		SuccessResponse(w, vet)

	default:
		ErrorResponse(w, http.StatusBadRequest, "Invalid role")
	}
}

// GetUser retrieves a user by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	userID := chi.URLParam(r, "id")
	if userID == "" {
		ErrorResponse(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// Authorization check: users can only access their own data unless they're admin
	if user.Sub != userID && user.Role != "admin" {
		ErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	db, ok := middleware.GetDBFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusInternalServerError, "Database not found in context")
		return
	}

	// Try to get as client first
	if client, err := db.GetClientByID(r.Context(), userID); err == nil {
		SuccessResponse(w, client)
		return
	}

	// Try to get as veterinarian
	if vet, err := db.GetVeterinarianByID(r.Context(), userID); err == nil {
		SuccessResponse(w, vet)
		return
	}

	ErrorResponse(w, http.StatusNotFound, "User not found")
}

// UpdateUser updates a user profile
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	userID := chi.URLParam(r, "id")
	if userID == "" {
		ErrorResponse(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// Authorization check: users can only update their own data unless they're admin
	if user.Sub != userID && user.Role != "admin" {
		ErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	db, ok := middleware.GetDBFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusInternalServerError, "Database not found in context")
		return
	}

	// Try to update as client first
	if client, err := db.GetClientByID(r.Context(), userID); err == nil {
		client.Name = req.Name
		client.Email = req.Email
		client.Phone = req.Phone
		client.Address = req.Address

		if err := db.UpdateClient(r.Context(), client); err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to update client")
			return
		}
		SuccessResponse(w, client)
		return
	}

	// Try to update as veterinarian
	if vet, err := db.GetVeterinarianByID(r.Context(), userID); err == nil {
		vet.Name = req.Name
		vet.Email = req.Email
		vet.Phone = req.Phone

		if err := db.UpdateVeterinarian(r.Context(), vet); err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to update veterinarian")
			return
		}
		SuccessResponse(w, vet)
		return
	}

	ErrorResponse(w, http.StatusNotFound, "User not found")
}

// DeleteUser deletes a user profile
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	userID := chi.URLParam(r, "id")
	if userID == "" {
		ErrorResponse(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// Only admins can delete users
	if user.Role != "admin" {
		ErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	db, ok := middleware.GetDBFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusInternalServerError, "Database not found in context")
		return
	}
	if err := db.DeleteUser(r.Context(), userID); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	MessageResponse(w, http.StatusOK, "User deleted successfully")
}

// ListUsers lists all users with pagination
func ListUsers(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	// Only admins can list all users
	if user.Role != "admin" {
		ErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // default limit
	offset := 0 // default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	db, ok := middleware.GetDBFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusInternalServerError, "Database not found in context")
		return
	}
	users, err := db.ListUsers(r.Context(), limit, offset)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to list users")
		return
	}

	SuccessResponse(w, users)
}
