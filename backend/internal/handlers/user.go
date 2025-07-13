// Package handlers/user.go contains user-related HTTP handlers
package handlers

import (
	"net/http"
	"pet-mgt/backend/internal/middleware"
)

// GetUserProfile returns the authenticated user's profile information
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get user from JWT token (already verified by middleware)
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	// Return user info
	JSONResponse(w, http.StatusOK, map[string]any{
		"user_id": user.Sub,
		"email":   user.Email,
		"role":    user.Role,
	})
}
