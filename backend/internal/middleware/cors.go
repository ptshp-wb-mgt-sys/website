// Package middleware/cors.go contains the cors middleware configuration
package middleware

import (
	"net/http"
	"pet-mgt/backend/internal/config"

	"github.com/go-chi/cors"
)

// CORS sets up the cors middleware with standard options
func CORS(cfg *config.Config) func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: []string{
			cfg.FrontendURL,
			"http://localhost:3000",
			"http://localhost:5173",
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		ExposedHeaders: []string{
			"Link",
		},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
