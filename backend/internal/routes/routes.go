// Package routes contains the routes for the application
package routes

import (
	"pet-mgt/backend/internal/config"
	"pet-mgt/backend/internal/handlers"
	"pet-mgt/backend/internal/middleware"
	"pet-mgt/backend/internal/store"
	"time"

	"github.com/go-chi/chi/v5"
	chiMw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

// SetupRouter sets up the chi router
func SetupRouter(cfg *config.Config, db store.Database) *chi.Mux {
	r := chi.NewRouter()
	setupGlobalMiddleware(cfg, r)
	r.Use(chiMw.Heartbeat("/ping"))

	// Initialize all handlers with database dependency
	h := handlers.NewHandlers(db)

	r.Route("/api/v1", func(r chi.Router) {
		publicRoutes(r)
		protectedRoutes(r, h, cfg)
	})

	return r
}

// publicRoutes sets up the public routes
func publicRoutes(r chi.Router) {
	r.Use(httprate.LimitByIP(60, 1*time.Minute))
	// some public handlers handlers
}

// protectedRoutes sets up the protected routes
func protectedRoutes(r chi.Router, h *handlers.Handlers, cfg *config.Config) {
	r.Use(httprate.LimitByIP(100, 1*time.Minute))
	r.Use(middleware.JWTAuth(cfg))

	// User routes
	r.Get("/profile", h.User.GetUserProfile)
	r.Post("/users", h.User.CreateUser)
	r.Get("/users", h.User.ListUsers)
	r.Get("/users/{id}", h.User.GetUser)
	r.Put("/users/{id}", h.User.UpdateUser)
	r.Delete("/users/{id}", h.User.DeleteUser)

	// Pet routes
	r.Post("/pets", h.Pet.CreatePet)
	r.Get("/pets/{id}", h.Pet.GetPet)
	r.Put("/pets/{id}", h.Pet.UpdatePet)
	r.Delete("/pets/{id}", h.Pet.DeletePet)
	r.Get("/clients/{clientId}/pets", h.Pet.GetPetsByClient)

	// Medical record routes
	r.Post("/pets/{petId}/medical-records", h.MedicalRecord.CreateMedicalRecord)
	r.Get("/pets/{petId}/medical-records", h.MedicalRecord.GetMedicalRecords)
	r.Get("/medical-records/{id}", h.MedicalRecord.GetMedicalRecord)
	r.Put("/medical-records/{id}", h.MedicalRecord.UpdateMedicalRecord)
	r.Delete("/medical-records/{id}", h.MedicalRecord.DeleteMedicalRecord)
}

// setupGlobalMiddleware sets up the middleware for the router
func setupGlobalMiddleware(cfg *config.Config, r *chi.Mux) {
	r.Use(chiMw.Logger)
	r.Use(chiMw.Recoverer)
	r.Use(chiMw.Timeout(60 * time.Second))
	r.Use(middleware.CORS(cfg))
}
