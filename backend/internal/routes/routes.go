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
		// Public routes with their own middleware
		r.Group(func(r chi.Router) {
			r.Use(httprate.LimitByIP(60, 1*time.Minute))
			publicRoutes(r, h)
		})

		// Protected routes with their own middleware
		r.Group(func(r chi.Router) {
			r.Use(httprate.LimitByIP(100, 1*time.Minute))
			r.Use(middleware.JWTAuth(cfg, db))
			protectedRoutes(r, h)
		})
	})

	return r
}

// publicRoutes sets up the public routes
func publicRoutes(r chi.Router, h *handlers.Handlers) {
	// Middleware is now handled in the caller

	// Public QR code access (no authentication required)
	r.Get("/public/pets/{publicUrl}", h.QRCode.GetPublicPetProfile)
	r.Get(
		"/pets/public/{publicUrl}",
		h.QRCode.GetPublicPetProfile,
	) // Alternative route format
}

// protectedRoutes sets up the protected routes
func protectedRoutes(r chi.Router, h *handlers.Handlers) {
	// Middleware is now handled in the caller

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

	// QR Code routes
	r.Post("/pets/{petId}/qr-code", h.QRCode.GenerateQRCode)
	r.Get("/pets/{petId}/qr-code", h.QRCode.GetQRCode)
	r.Put("/pets/{petId}/qr-code", h.QRCode.UpdateQRCode)
	r.Delete("/pets/{petId}/qr-code", h.QRCode.DeleteQRCode)

	// Medical record routes
	r.Post("/pets/{petId}/medical-records", h.MedicalRecord.CreateMedicalRecord)
	r.Get("/pets/{petId}/medical-records", h.MedicalRecord.GetMedicalRecords)
	r.Get("/medical-records/{id}", h.MedicalRecord.GetMedicalRecord)
	r.Put("/medical-records/{id}", h.MedicalRecord.UpdateMedicalRecord)
	r.Delete("/medical-records/{id}", h.MedicalRecord.DeleteMedicalRecord)

	// Appointment routes
	r.Post("/appointments", h.Appointment.CreateAppointment)
	r.Get("/appointments", h.Appointment.GetAppointments)
	r.Get("/appointments/{id}", h.Appointment.GetAppointment)
	r.Put("/appointments/{id}", h.Appointment.UpdateAppointment)
	r.Delete("/appointments/{id}", h.Appointment.DeleteAppointment)

	// Veterinarian and appointment availability routes
	r.Get("/veterinarians", h.Appointment.ListVeterinarians)
	r.Get("/veterinarians/{vetId}/availability", h.Appointment.GetAvailableSlots)
	// Vet availability management (vet self or admin)
	r.Post("/veterinarians/{id}/availability", h.Appointment.SetAvailability)

	// Product routes
	r.Post("/products", h.Product.CreateProduct)
	r.Get("/products", h.Product.GetProducts)
	r.Get("/products/{id}", h.Product.GetProduct)
	r.Put("/products/{id}", h.Product.UpdateProduct)
	r.Delete("/products/{id}", h.Product.DeleteProduct)
	r.Get("/veterinarians/{vetId}/products", h.Product.GetVeterinarianProducts)
	r.Put("/products/{id}/stock", h.Product.UpdateProductStock)
	r.Post("/products/checkout", h.Product.CheckoutProducts)

	// Order routes
	r.Post("/orders", h.Order.CreateOrder)
	r.Get("/orders", h.Order.GetOrders)
	r.Get("/orders/{id}", h.Order.GetOrder)
	r.Put("/orders/{id}/status", h.Order.UpdateOrderStatus)
	r.Delete("/orders/{id}", h.Order.CancelOrder)
}

// setupGlobalMiddleware sets up the middleware for the router
func setupGlobalMiddleware(cfg *config.Config, r *chi.Mux) {
	r.Use(chiMw.Logger)
	r.Use(chiMw.Recoverer)
	r.Use(chiMw.Timeout(60 * time.Second))
	r.Use(middleware.CORS(cfg))
}
