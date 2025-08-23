// Package handlers/medical_records.go contains medical record management HTTP handlers
package handlers

import (
	"encoding/json"
	"net/http"
	"pet-mgt/backend/internal/middleware"
	"pet-mgt/backend/internal/store"
	"time"

	"github.com/go-chi/chi/v5"
)

// MedicalRecordHandler handles medical record-related HTTP requests
type MedicalRecordHandler struct {
	db store.Database
}

// NewMedicalRecordHandler creates a new MedicalRecordHandler with database dependency
func NewMedicalRecordHandler(db store.Database) *MedicalRecordHandler {
	return &MedicalRecordHandler{
		db: db,
	}
}

// CreateMedicalRecord creates a new medical record for a pet
func (h *MedicalRecordHandler) CreateMedicalRecord(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	// Require a veterinarian profile (admin must also have a veterinarian profile to satisfy FK)
	if _, err := h.db.GetVeterinarianByID(r.Context(), user.Sub); err != nil && user.Role != "admin" {
		ErrorResponse(w, http.StatusForbidden, "Veterinarian profile required to create medical records")
		return
	}

	petID := chi.URLParam(r, "petId")
	if petID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Pet ID is required")
		return
	}

	var req struct {
		DateOfVisit          time.Time `json:"date_of_visit"`
		ReasonForVisit       string    `json:"reason_for_visit"`
		Diagnosis            string    `json:"diagnosis"`
		MedicationPrescribed []string  `json:"medication_prescribed"`
		Notes                string    `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Verify pet exists
	_, err := h.db.GetPetByID(r.Context(), petID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Pet not found")
		return
	}

	// Use current time if date not provided
	if req.DateOfVisit.IsZero() {
		req.DateOfVisit = time.Now()
	}

	record := store.NewMedicalRecord(
		petID,
		user.Sub,
		req.ReasonForVisit,
		req.Diagnosis,
		req.DateOfVisit,
		req.MedicationPrescribed,
		req.Notes,
	)

	if err := h.db.CreateMedicalRecord(r.Context(), record); err != nil {
		ErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to create medical record",
		)
		return
	}

	SuccessResponse(w, record)
}

// GetMedicalRecords retrieves all medical records for a pet
func (h *MedicalRecordHandler) GetMedicalRecords(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	petID := chi.URLParam(r, "petId")
	if petID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Pet ID is required")
		return
	}

	// Verify pet exists and check authorization
	pet, err := h.db.GetPetByID(r.Context(), petID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Pet not found")
		return
	}

	// Authorization check: clients can only access their own pet's records, vets and admins can access any
	if user.Role == "client" && pet.OwnerID != user.Sub {
		ErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	records, err := h.db.GetMedicalRecordsByPetID(r.Context(), petID)
	if err != nil {
		ErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve medical records",
		)
		return
	}

	SuccessResponse(w, records)
}

// GetMedicalRecord retrieves a specific medical record
func (h *MedicalRecordHandler) GetMedicalRecord(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	recordID := chi.URLParam(r, "id")
	if recordID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Record ID is required")
		return
	}

	record, err := h.db.GetMedicalRecordByID(r.Context(), recordID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Medical record not found")
		return
	}

	// Get the pet to check ownership
	pet, err := h.db.GetPetByID(r.Context(), record.PetID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Pet not found")
		return
	}

	// Authorization check: clients can only access their own pet's records, vets and admins can access any
	if user.Role == "client" && pet.OwnerID != user.Sub {
		ErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	SuccessResponse(w, record)
}

// UpdateMedicalRecord updates a medical record
func (h *MedicalRecordHandler) UpdateMedicalRecord(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	// Require a veterinarian profile; this also satisfies the FK constraint when updating
	if _, err := h.db.GetVeterinarianByID(r.Context(), user.Sub); err != nil {
		ErrorResponse(w, http.StatusForbidden, "Only veterinarians can update medical records")
		return
	}

	recordID := chi.URLParam(r, "id")
	if recordID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Record ID is required")
		return
	}

	record, err := h.db.GetMedicalRecordByID(r.Context(), recordID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Medical record not found")
		return
	}

	var req struct {
		DateOfVisit          time.Time `json:"date_of_visit"`
		ReasonForVisit       string    `json:"reason_for_visit"`
		Diagnosis            string    `json:"diagnosis"`
		MedicationPrescribed []string  `json:"medication_prescribed"`
		Notes                string    `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update record fields
	record.DateOfVisit = req.DateOfVisit
	record.ReasonForVisit = req.ReasonForVisit
	record.Diagnosis = req.Diagnosis
	record.MedicationPrescribed = req.MedicationPrescribed
	record.Notes = req.Notes
	record.UpdatedAt = time.Now()

	if err := h.db.UpdateMedicalRecord(r.Context(), record); err != nil {
		ErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to update medical record",
		)
		return
	}

	SuccessResponse(w, record)
}

// DeleteMedicalRecord deletes a medical record
func (h *MedicalRecordHandler) DeleteMedicalRecord(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	// Require a veterinarian profile
	if _, err := h.db.GetVeterinarianByID(r.Context(), user.Sub); err != nil {
		ErrorResponse(w, http.StatusForbidden, "Only veterinarians can delete medical records")
		return
	}

	recordID := chi.URLParam(r, "id")
	if recordID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Record ID is required")
		return
	}

	// Check if record exists
	_, err := h.db.GetMedicalRecordByID(r.Context(), recordID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Medical record not found")
		return
	}

	if err := h.db.DeleteMedicalRecord(r.Context(), recordID); err != nil {
		ErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to delete medical record",
		)
		return
	}

	MessageResponse(w, http.StatusOK, "Medical record deleted successfully")
}
