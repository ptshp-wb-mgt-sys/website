// Package handlers contains appointment management handlers
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"pet-mgt/backend/internal/middleware"
	"pet-mgt/backend/internal/store"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

// AppointmentHandler handles appointment operations
type AppointmentHandler struct {
	db store.Database
}

// NewAppointmentHandler creates a new AppointmentHandler
func NewAppointmentHandler(db store.Database) *AppointmentHandler {
	return &AppointmentHandler{db: db}
}

// CreateAppointment books a new appointment
func (h *AppointmentHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse request body
	var req struct {
		VeterinarianID  string    `json:"veterinarian_id"`
		PetID           string    `json:"pet_id"`
		AppointmentDate time.Time `json:"appointment_date"`
		DurationMinutes int       `json:"duration_minutes"`
		Reason          string    `json:"reason"`
		Notes           string    `json:"notes,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.VeterinarianID == "" || req.PetID == "" || req.Reason == "" {
		ErrorResponse(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	if req.DurationMinutes <= 0 {
		req.DurationMinutes = 30 // Default duration
	}

	// Verify user is client and owns the pet (unless admin)
	role := deriveRole(r.Context(), h.db, user)
	if role != "admin" {
		if role != "client" {
			ErrorResponse(w, http.StatusForbidden, "Only clients can book appointments")
			return
		}

		// Verify pet ownership
		pet, err := h.db.GetPetByID(r.Context(), req.PetID)
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "Pet not found")
			return
		}
		if pet.OwnerID != user.Sub {
			ErrorResponse(
				w,
				http.StatusForbidden,
				"You can only book appointments for your own pets",
			)
			return
		}
	}

	// Verify veterinarian exists
	_, err := h.db.GetVeterinarianByID(r.Context(), req.VeterinarianID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Veterinarian not found")
		return
	}

	// Create appointment
	appointment := store.NewAppointment(
		user.Sub,
		req.VeterinarianID,
		req.PetID,
		req.AppointmentDate,
		req.DurationMinutes,
		req.Reason,
	)
	appointment.Notes = req.Notes

	// Validate against availability and conflicts
	dateOnly := req.AppointmentDate.UTC().Format("2006-01-02")
	date, _ := time.ParseInLocation("2006-01-02", dateOnly, time.UTC)
	availableSlots, err := h.db.GetAvailableAppointmentSlots(r.Context(), req.VeterinarianID, date)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to validate availability")
		return
	}

	// Check that requested time aligns with an available slot start
	requestedStart := req.AppointmentDate
	matchFound := false
	for _, s := range availableSlots {
		if s.Available && s.StartTime.Equal(requestedStart) {
			matchFound = true
			break
		}
	}
	if !matchFound {
		ErrorResponse(w, http.StatusConflict, "Selected time is not available")
		return
	}

	if err := h.db.CreateAppointment(r.Context(), appointment); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to create appointment")
		return
	}

	SuccessResponse(w, appointment)
}

// GetAppointments retrieves appointments for the current user
func (h *AppointmentHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var appointments []store.Appointment
	var err error

	role := deriveRole(r.Context(), h.db, user)

	switch role {
	case "client":
		appointments, err = h.db.GetAppointmentsByClientID(r.Context(), user.Sub)
	case "veterinarian":
		appointments, err = h.db.GetAppointmentsByVeterinarianID(r.Context(), user.Sub)
	case "admin":
		// For admin, check query params to see if they want specific user's appointments
		clientID := r.URL.Query().Get("client_id")
		vetID := r.URL.Query().Get("veterinarian_id")

		if clientID != "" {
			appointments, err = h.db.GetAppointmentsByClientID(r.Context(), clientID)
		} else if vetID != "" {
			appointments, err = h.db.GetAppointmentsByVeterinarianID(r.Context(), vetID)
		} else {
			// NOTE: Return all appointments would need a new method - for now return empty
			appointments = []store.Appointment{}
		}
	default:
		ErrorResponse(w, http.StatusForbidden, "Invalid user role")
		return
	}

	if err != nil {
		ErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve appointments",
		)
		return
	}

	SuccessResponse(w, appointments)
}

// GetAppointment retrieves a specific appointment
func (h *AppointmentHandler) GetAppointment(w http.ResponseWriter, r *http.Request) {
	appointmentID := chi.URLParam(r, "id")
	if appointmentID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Appointment ID is required")
		return
	}

	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get appointment
	appointment, err := h.db.GetAppointmentByID(r.Context(), appointmentID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Appointment not found")
		return
	}

	// Check permissions
	role := deriveRole(r.Context(), h.db, user)
	if role != "admin" {
		if role == "client" && appointment.ClientID != user.Sub {
			ErrorResponse(
				w,
				http.StatusForbidden,
				"You can only view your own appointments",
			)
			return
		}
		if role == "veterinarian" && appointment.VeterinarianID != user.Sub {
			ErrorResponse(
				w,
				http.StatusForbidden,
				"You can only view your own appointments",
			)
			return
		}
	}

	SuccessResponse(w, appointment)
}

// UpdateAppointment updates an appointment
func (h *AppointmentHandler) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	appointmentID := chi.URLParam(r, "id")
	if appointmentID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Appointment ID is required")
		return
	}

	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get existing appointment
	appointment, err := h.db.GetAppointmentByID(r.Context(), appointmentID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Appointment not found")
		return
	}

	// Check permissions
	canUpdate := false
	role := deriveRole(r.Context(), h.db, user)
	if role == "admin" {
		canUpdate = true
	} else if role == "client" && appointment.ClientID == user.Sub {
		canUpdate = true
	} else if role == "veterinarian" && appointment.VeterinarianID == user.Sub {
		canUpdate = true
	}

	if !canUpdate {
		ErrorResponse(
			w,
			http.StatusForbidden,
			"You can only update your own appointments",
		)
		return
	}

	// Parse request body
	var updateData struct {
		AppointmentDate *time.Time `json:"appointment_date,omitempty"`
		DurationMinutes *int       `json:"duration_minutes,omitempty"`
		Reason          string     `json:"reason,omitempty"`
		Status          string     `json:"status,omitempty"`
		Notes           string     `json:"notes,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update fields
	if updateData.AppointmentDate != nil {
		appointment.AppointmentDate = *updateData.AppointmentDate
	}
	if updateData.DurationMinutes != nil {
		appointment.DurationMinutes = *updateData.DurationMinutes
	}
	if updateData.Reason != "" {
		appointment.Reason = updateData.Reason
	}
	if updateData.Status != "" {
		// Validate status
		validStatuses := map[string]bool{
			"scheduled":   true,
			"completed":   true,
			"cancelled":   true,
			"rescheduled": true,
		}
		if !validStatuses[updateData.Status] {
			ErrorResponse(w, http.StatusBadRequest, "Invalid status")
			return
		}
		appointment.Status = updateData.Status
	}
	if updateData.Notes != "" {
		appointment.Notes = updateData.Notes
	}

	appointment.UpdatedAt = time.Now()

	// Update appointment
	if err := h.db.UpdateAppointment(r.Context(), appointment); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update appointment")
		return
	}

	SuccessResponse(w, appointment)
}

// DeleteAppointment cancels an appointment
func (h *AppointmentHandler) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	appointmentID := chi.URLParam(r, "id")
	if appointmentID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Appointment ID is required")
		return
	}

	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get appointment
	appointment, err := h.db.GetAppointmentByID(r.Context(), appointmentID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Appointment not found")
		return
	}

	// Check permissions
	canDelete := false
	role := deriveRole(r.Context(), h.db, user)
	if role == "admin" {
		canDelete = true
	} else if role == "client" && appointment.ClientID == user.Sub {
		canDelete = true
	} else if role == "veterinarian" && appointment.VeterinarianID == user.Sub {
		canDelete = true
	}

	if !canDelete {
		ErrorResponse(
			w,
			http.StatusForbidden,
			"You can only cancel your own appointments",
		)
		return
	}

	// Cancel appointment (soft delete by updating status)
	appointment.Status = "cancelled"
	appointment.UpdatedAt = time.Now()

	if err := h.db.UpdateAppointment(r.Context(), appointment); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to cancel appointment")
		return
	}

	MessageResponse(w, http.StatusOK, "Appointment cancelled successfully")
}

// GetAvailableSlots retrieves available appointment slots for a veterinarian
func (h *AppointmentHandler) GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	vetID := chi.URLParam(r, "vetId")
	if vetID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Veterinarian ID is required")
		return
	}

	// Parse date from query parameter
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "Date parameter is required")
		return
	}

	date, err := time.ParseInLocation("2006-01-02", dateStr, time.UTC)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD")
		return
	}

	// Get available slots
	slots, err := h.db.GetAvailableAppointmentSlots(r.Context(), vetID, date)
	if err != nil {
		ErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve available slots",
		)
		return
	}

	SuccessResponse(w, slots)
}

// SetAvailability allows a veterinarian to set or update their available working hours
func (h *AppointmentHandler) SetAvailability(w http.ResponseWriter, r *http.Request) {
	vetID := chi.URLParam(r, "id")
	if vetID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Veterinarian ID is required")
		return
	}

	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Only the veterinarian themselves or admin can update availability
	role := deriveRole(r.Context(), h.db, user)
	if role != "admin" && !(role == "veterinarian" && user.Sub == vetID) {
		ErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	var req struct {
		AvailableHours []store.WorkingHours `json:"available_hours"`
		ClinicAddress  string               `json:"clinic_address,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	vet, err := h.db.GetVeterinarianByID(r.Context(), vetID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Veterinarian not found")
		return
	}

	// Basic normalization: ensure day_of_week is capitalized consistently (e.g., Mon, Tue)
	normalized := make([]store.WorkingHours, 0, len(req.AvailableHours))
	for _, wh := range req.AvailableHours {
		normalized = append(normalized, store.WorkingHours{
			DayOfWeek: normalizeDay(wh.DayOfWeek),
			Start:     wh.Start,
			End:       wh.End,
		})
	}

	vet.AvailableHours = normalized
	if req.ClinicAddress != "" {
		vet.ClinicAddress = req.ClinicAddress
	}

	if err := h.db.UpdateVeterinarian(r.Context(), vet); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update availability")
		return
	}

	SuccessResponse(w, vet)
}

// normalizeDay converts various inputs to a canonical 3-letter day (e.g., Monday -> Mon)
func normalizeDay(input string) string {
	switch strings.ToLower(input) {
	case "monday", "mon":
		return "Mon"
	case "tuesday", "tue", "tues":
		return "Tue"
	case "wednesday", "wed":
		return "Wed"
	case "thursday", "thu", "thur", "thurs":
		return "Thu"
	case "friday", "fri":
		return "Fri"
	case "saturday", "sat":
		return "Sat"
	case "sunday", "sun":
		return "Sun"
	default:
		return input
	}
}

// deriveRole maps generic or missing roles from the JWT to concrete application roles
// by querying the database. If the user's role is already a concrete role, it is returned as-is.
func deriveRole(ctx context.Context, db store.Database, user *middleware.UserClaims) string {
	role := user.Role
	if role == "authenticated" || role == "user" || role == "" {
		if c, err := db.GetClientByID(ctx, user.Sub); err == nil && c != nil && c.Role != "" {
			return c.Role
		}
		if v, err := db.GetVeterinarianByID(ctx, user.Sub); err == nil && v != nil && v.Role != "" {
			return v.Role
		}
	}
	return role
}

// ListVeterinarians returns all veterinarians for appointment booking
func (h *AppointmentHandler) ListVeterinarians(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

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

	// NOTE: For now, we'll use the existing ListUsers method and filter by role
	// In a real implementation, you'd want a dedicated method for listing veterinarians
	users, err := h.db.ListUsers(r.Context(), limit, offset)
	if err != nil {
		ErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve veterinarians",
		)
		return
	}

	// Filter to only veterinarians and convert to proper format
	var veterinarians []map[string]any
	for _, user := range users {
		if user.Role == "veterinarian" {
			// Get full veterinarian details
			vet, err := h.db.GetVeterinarianByID(r.Context(), user.ID)
			if err == nil {
				veterinarians = append(veterinarians, map[string]any{
					"id":              vet.ID,
					"name":            vet.Name,
					"email":           vet.Email,
					"phone":           vet.Phone,
					"clinic_address":  vet.ClinicAddress,
					"available_hours": vet.AvailableHours,
				})
			}
		}
	}

	SuccessResponse(w, veterinarians)
}
