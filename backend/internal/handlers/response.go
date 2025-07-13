// Package handlers/response.go contains HTTP response utilities
package handlers

import (
	"encoding/json"
	"net/http"
)

// JSONResponse sends a JSON response with the given status and data
func JSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ErrorResponse sends a standardized error response
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, map[string]string{"error": message})
}

// SuccessResponse sends a standardized success response
func SuccessResponse(w http.ResponseWriter, data any) {
	JSONResponse(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    data,
	})
}

// MessageResponse sends a simple message response
func MessageResponse(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, map[string]string{"message": message})
}
