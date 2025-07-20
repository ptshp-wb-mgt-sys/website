// Package handlers contains HTTP handlers for the application
package handlers

import "pet-mgt/backend/internal/store"

// Handlers contains all HTTP handlers for the application
type Handlers struct {
	User          *UserHandler
	Pet           *PetHandler
	MedicalRecord *MedicalRecordHandler
}

// NewHandlers creates a new Handlers instance with all handler dependencies
func NewHandlers(db store.Database) *Handlers {
	return &Handlers{
		User:          NewUserHandler(db),
		Pet:           NewPetHandler(db),
		MedicalRecord: NewMedicalRecordHandler(db),
	}
}
