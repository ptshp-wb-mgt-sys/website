// Package handlers contains QR code management handlers
package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"pet-mgt/backend/internal/middleware"
	"pet-mgt/backend/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

// QRCodeHandler handles QR code operations
type QRCodeHandler struct {
	db store.Database
}

// NewQRCodeHandler creates a new QRCodeHandler
func NewQRCodeHandler(db store.Database) *QRCodeHandler {
	return &QRCodeHandler{db: db}
}

// GenerateQRCode generates a QR code for a pet
func (h *QRCodeHandler) GenerateQRCode(w http.ResponseWriter, r *http.Request) {
	petID := chi.URLParam(r, "petId")
	if petID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Pet ID is required")
		return
	}

	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID, role := user.Sub, user.Role

	// Get pet details
	pet, err := h.db.GetPetByID(r.Context(), petID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Pet not found")
		return
	}

	// Check permissions - only pet owner or admin can generate QR codes
	if role != "admin" && pet.OwnerID != userID {
		ErrorResponse(w, http.StatusForbidden, "Forbidden: You can only generate QR codes for your own pets")
		return
	}

	// Get owner details
	owner, err := h.db.GetClientByID(r.Context(), pet.OwnerID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Owner not found")
		return
	}

	// Check if QR code already exists
	existingQR, err := h.db.GetQRCodeByPetID(r.Context(), petID)
	if err == nil && existingQR != nil {
		// Return existing QR code
		SuccessResponse(w, existingQR)
		return
	}

	// Generate unique public URL
	publicURL := fmt.Sprintf("/public/pets/%s", uuid.New().String())

	// Create encoded content
	encodedContent := store.EncodedContent{
		PetName:          pet.Name,
		PetType:          pet.Type,
		OwnerName:        owner.Name,
		OwnerPhone:       owner.Phone,
		OwnerEmail:       owner.Email,
		OwnerAddress:     owner.Address,
		PublicProfileURL: publicURL,
	}

	// Generate QR code image
	qrData := fmt.Sprintf("%s%s", getBaseURL(r), publicURL)
	qrCodeBytes, err := qrcode.Encode(qrData, qrcode.Medium, 256)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to generate QR code image")
		return
	}

	// Encode QR code as base64
	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCodeBytes)

	// Create QR code record
	qrCode := store.NewQRCode(petID, qrCodeBase64, publicURL, encodedContent)

	if err := h.db.CreateQRCode(r.Context(), qrCode); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to save QR code")
		return
	}

	SuccessResponse(w, qrCode)
}

// GetQRCode retrieves QR code for a pet
func (h *QRCodeHandler) GetQRCode(w http.ResponseWriter, r *http.Request) {
	petID := chi.URLParam(r, "petId")
	if petID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Pet ID is required")
		return
	}

	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID, role := user.Sub, user.Role

	// Get pet details for permission check
	pet, err := h.db.GetPetByID(r.Context(), petID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Pet not found")
		return
	}

	// Check permissions
	if role != "admin" && role != "veterinarian" && pet.OwnerID != userID {
		ErrorResponse(w, http.StatusForbidden, "Forbidden: Insufficient permissions")
		return
	}

	// Get QR code
	qrCode, err := h.db.GetQRCodeByPetID(r.Context(), petID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "QR code not found")
		return
	}

	SuccessResponse(w, qrCode)
}

// GetPublicPetProfile retrieves public pet profile via QR code URL (no authentication required)
func (h *QRCodeHandler) GetPublicPetProfile(w http.ResponseWriter, r *http.Request) {
	publicURL := chi.URLParam(r, "publicUrl")
	if publicURL == "" {
		ErrorResponse(w, http.StatusBadRequest, "Public URL is required")
		return
	}

	// Add leading slash if not present
	if publicURL[0] != '/' {
		publicURL = "/" + publicURL
	}

	// Get public pet profile
	profile, err := h.db.GetPublicPetProfile(r.Context(), publicURL)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Pet profile not found")
		return
	}

	SuccessResponse(w, profile)
}

// UpdateQRCode updates QR code information
func (h *QRCodeHandler) UpdateQRCode(w http.ResponseWriter, r *http.Request) {
	petID := chi.URLParam(r, "petId")
	if petID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Pet ID is required")
		return
	}

	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID, role := user.Sub, user.Role

	// Get pet details for permission check
	pet, err := h.db.GetPetByID(r.Context(), petID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Pet not found")
		return
	}

	// Check permissions - only pet owner or admin can update QR codes
	if role != "admin" && pet.OwnerID != userID {
		ErrorResponse(w, http.StatusForbidden, "Forbidden: You can only update QR codes for your own pets")
		return
	}

	// Parse request body
	var updateData struct {
		EmergencyContact string   `json:"emergency_contact,omitempty"`
		MedicalAlerts    []string `json:"medical_alerts,omitempty"`
		IsActive         *bool    `json:"is_active,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get existing QR code
	qrCode, err := h.db.GetQRCodeByPetID(r.Context(), petID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "QR code not found")
		return
	}

	// Update encoded content
	if updateData.EmergencyContact != "" {
		qrCode.EncodedContent.EmergencyContact = updateData.EmergencyContact
	}
	if updateData.MedicalAlerts != nil {
		qrCode.EncodedContent.MedicalAlerts = updateData.MedicalAlerts
	}
	if updateData.IsActive != nil {
		qrCode.IsActive = *updateData.IsActive
	}

	// Update QR code
	if err := h.db.UpdateQRCode(r.Context(), qrCode); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update QR code")
		return
	}

	SuccessResponse(w, qrCode)
}

// DeleteQRCode deactivates a QR code
func (h *QRCodeHandler) DeleteQRCode(w http.ResponseWriter, r *http.Request) {
	petID := chi.URLParam(r, "petId")
	if petID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Pet ID is required")
		return
	}

	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID, role := user.Sub, user.Role

	// Get pet details for permission check
	pet, err := h.db.GetPetByID(r.Context(), petID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Pet not found")
		return
	}

	// Check permissions - only pet owner or admin can delete QR codes
	if role != "admin" && pet.OwnerID != userID {
		ErrorResponse(w, http.StatusForbidden, "Forbidden: You can only delete QR codes for your own pets")
		return
	}

	// Get QR code
	qrCode, err := h.db.GetQRCodeByPetID(r.Context(), petID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "QR code not found")
		return
	}

	// Deactivate instead of hard delete
	qrCode.IsActive = false
	if err := h.db.UpdateQRCode(r.Context(), qrCode); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to deactivate QR code")
		return
	}

	MessageResponse(w, http.StatusOK, "QR code deactivated successfully")
}

// getBaseURL extracts the base URL from the request
func getBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// Check for forwarded protocol headers
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	return fmt.Sprintf("%s://%s", scheme, r.Host)
}
