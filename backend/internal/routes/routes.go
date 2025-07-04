// Package routes contains the routes for the application
package routes

import (
	"pet-mgt/backend/internal/config"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// SetupRouter sets up the chi router
func SetupRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()
	setupGloballMiddleware(r)

	r.Use(middleware.Heartbeat("/ping"))

	r.Route("/api/v1", func(r chi.Router) {
		protectedRoutes(r)
	})

	return r
}

// protectedRoutes sets up the routes that require authentication
func protectedRoutes(r chi.Router) {
	// add protected routes here
}

// setupGloballMiddleware sets up the middleware for the router
func setupGloballMiddleware(r *chi.Mux) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
}
