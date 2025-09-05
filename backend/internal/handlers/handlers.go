// Package handlers contains HTTP handlers for the application
package handlers

import (
	"pet-mgt/backend/internal/config"
	"pet-mgt/backend/internal/store"
)

// Handlers contains all HTTP handlers for the application
type Handlers struct {
	User          *UserHandler
	Pet           *PetHandler
	MedicalRecord *MedicalRecordHandler
	QRCode        *QRCodeHandler
	Appointment   *AppointmentHandler
	Product       *ProductHandler
	Order         *OrderHandler
}

// NewHandlers creates a new Handlers instance with all handler dependencies
func NewHandlers(cfg *config.Config, db store.Database) *Handlers {
	return &Handlers{
		User:          NewUserHandler(db),
		Pet:           NewPetHandler(db),
		MedicalRecord: NewMedicalRecordHandler(db),
		QRCode:        NewQRCodeHandler(db, cfg.FrontendURL),
		Appointment:   NewAppointmentHandler(db),
		Product:       NewProductHandler(db),
		Order:         NewOrderHandler(db),
	}
}
