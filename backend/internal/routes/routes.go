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

	r.Route("/api/v1", func(r chi.Router) {
		publicRoutes(r)
		protectedRoutes(r, db, cfg)
	})

	return r
}

// publicRoutes sets up the public routes
func publicRoutes(r chi.Router) {
	r.Use(httprate.LimitByIP(60, 1*time.Minute))
	// some public handlers handlers
}

// protectedRoutes sets up the protected routes
func protectedRoutes(r chi.Router, db store.Database, cfg *config.Config) {
	r.Use(httprate.LimitByIP(100, 1*time.Minute))
	r.Use(middleware.JWTAuth(cfg))
	r.Use(middleware.InjectDB(db))

	// User routes
	r.Get("/profile", handlers.GetUserProfile)
	r.Post("/users", handlers.CreateUser)
	r.Get("/users", handlers.ListUsers)
	r.Get("/users/{id}", handlers.GetUser)
	r.Put("/users/{id}", handlers.UpdateUser)
	r.Delete("/users/{id}", handlers.DeleteUser)

	// Pet routes
	r.Post("/pets", handlers.CreatePet)
	r.Get("/pets/{id}", handlers.GetPet)
	r.Put("/pets/{id}", handlers.UpdatePet)
	r.Delete("/pets/{id}", handlers.DeletePet)
	r.Get("/clients/{clientId}/pets", handlers.GetPetsByClient)

	// Medical record routes
	r.Post("/pets/{petId}/medical-records", handlers.CreateMedicalRecord)
	r.Get("/pets/{petId}/medical-records", handlers.GetMedicalRecords)
	r.Get("/medical-records/{id}", handlers.GetMedicalRecord)
	r.Put("/medical-records/{id}", handlers.UpdateMedicalRecord)
	r.Delete("/medical-records/{id}", handlers.DeleteMedicalRecord)
}

// setupGlobalMiddleware sets up the middleware for the router
func setupGlobalMiddleware(cfg *config.Config, r *chi.Mux) {
	r.Use(middleware.CORS(cfg))
	r.Use(chiMw.Logger)
	r.Use(chiMw.Recoverer)
	r.Use(chiMw.Timeout(30 * time.Second))
}
