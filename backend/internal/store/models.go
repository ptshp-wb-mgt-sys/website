// Package store/models.go contains the models for the application
package store

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Database interface defines methods for data access operations
type Database interface {
	// User operations
	GetUserByID(ctx context.Context, userID string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, userID string) error
	ListUsers(ctx context.Context, limit, offset int) ([]User, error)

	// Client/Veterinarian specific
	CreateClient(ctx context.Context, client *Client) error
	CreateVeterinarian(ctx context.Context, vet *Veterinarian) error
	UpdateClient(ctx context.Context, client *Client) error
	UpdateVeterinarian(ctx context.Context, vet *Veterinarian) error
	GetClientByID(ctx context.Context, clientID string) (*Client, error)
	GetVeterinarianByID(ctx context.Context, vetID string) (*Veterinarian, error)

	// Pet operations
	GetPetsByUserID(ctx context.Context, userID string) ([]Pet, error)
	GetPetByID(ctx context.Context, petID string) (*Pet, error)
	CreatePet(ctx context.Context, pet *Pet) error
	UpdatePet(ctx context.Context, pet *Pet) error
	DeletePet(ctx context.Context, petID string) error

	// Medical record operations
	GetMedicalRecordsByPetID(ctx context.Context, petID string) ([]MedicalRecord, error)
	GetMedicalRecordByID(ctx context.Context, recordID string) (*MedicalRecord, error)
	CreateMedicalRecord(ctx context.Context, record *MedicalRecord) error
	UpdateMedicalRecord(ctx context.Context, record *MedicalRecord) error
	DeleteMedicalRecord(ctx context.Context, recordID string) error

	// QR Code operations
	GetQRCodeByPetID(ctx context.Context, petID string) (*QRCode, error)
	GetQRCodeByPublicURL(ctx context.Context, publicURL string) (*QRCode, error)
	CreateQRCode(ctx context.Context, qrCode *QRCode) error
	UpdateQRCode(ctx context.Context, qrCode *QRCode) error
	DeleteQRCode(ctx context.Context, qrCodeID string) error
	GetPublicPetProfile(
		ctx context.Context,
		publicURL string,
	) (*PublicPetProfile, error)

	// Appointment operations
	GetAppointmentsByClientID(
		ctx context.Context,
		clientID string,
	) ([]Appointment, error)
	GetAppointmentsByVeterinarianID(
		ctx context.Context,
		vetID string,
	) ([]Appointment, error)
	GetAppointmentByID(ctx context.Context, appointmentID string) (*Appointment, error)
	CreateAppointment(ctx context.Context, appointment *Appointment) error
	UpdateAppointment(ctx context.Context, appointment *Appointment) error
	DeleteAppointment(ctx context.Context, appointmentID string) error
	GetAvailableAppointmentSlots(
		ctx context.Context,
		vetID string,
		date time.Time,
	) ([]TimeSlot, error)

	// Product operations
	GetProductsByVeterinarianID(ctx context.Context, vetID string) ([]Product, error)
	GetProductByID(ctx context.Context, productID string) (*Product, error)
	CreateProduct(ctx context.Context, product *Product) error
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, productID string) error
	ListProducts(ctx context.Context, filters ProductFilters) ([]Product, error)
	UpdateProductStock(ctx context.Context, productID string, quantity int) error

	// Order operations
	GetOrdersByClientID(ctx context.Context, clientID string) ([]Order, error)
	GetOrdersByVeterinarianID(ctx context.Context, vetID string) ([]Order, error)
	GetOrderByID(ctx context.Context, orderID string) (*Order, error)
	CreateOrder(ctx context.Context, order *Order) error
	UpdateOrderStatus(ctx context.Context, orderID string, status string) error
	GetOrderItems(ctx context.Context, orderID string) ([]OrderItem, error)
	CreateOrderItem(ctx context.Context, item *OrderItem) error

	// Health check
	Ping(ctx context.Context) error

	// Cleanup
	Close() error
}

// User represents a user in the system (base interface for Client and Veterinarian)
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Client represents a pet owner
type Client struct {
	ID      string `json:"id"      db:"id"`
	Name    string `json:"name"    db:"name"`
	Email   string `json:"email"   db:"email"`
	Phone   string `json:"phone"   db:"phone"`
	Address string `json:"address" db:"address"`
	Role    string `json:"role"    db:"role"`
}

// Veterinarian represents a veterinarian user
type Veterinarian struct {
	ID             string         `json:"id"              db:"id"`
	Name           string         `json:"name"            db:"name"`
	Email          string         `json:"email"           db:"email"`
	Phone          string         `json:"phone"           db:"phone"`
	ClinicAddress  string         `json:"clinic_address"  db:"clinic_address"`
	AvailableHours []WorkingHours `json:"available_hours" db:"available_hours"`
	Role           string         `json:"role"            db:"role"`
}

// WorkingHours represents available working hours for a veterinarian
type WorkingHours struct {
	DayOfWeek string `json:"day_of_week" db:"day_of_week"`
	Start     string `json:"start"       db:"start"`
	End       string `json:"end"         db:"end"`
}

// Pet represents a pet in the system
type Pet struct {
	ID          string    `json:"id"            db:"id"`
	OwnerID     string    `json:"owner_id"      db:"owner_id"`
	Name        string    `json:"name"          db:"name"`
	Type        string    `json:"type"          db:"type"`
	Breed       string    `json:"breed"         db:"breed"`
	DateOfBirth time.Time `json:"date_of_birth" db:"date_of_birth"`
	Weight      float64   `json:"weight"        db:"weight"`
	CreatedAt   time.Time `json:"created_at"    db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"    db:"updated_at"`
}

// MedicalRecord represents a veterinary visit record
type MedicalRecord struct {
	ID                   string    `json:"id"                    db:"id"`
	PetID                string    `json:"pet_id"                db:"pet_id"`
	VeterinarianID       string    `json:"veterinarian_id"       db:"veterinarian_id"`
	DateOfVisit          time.Time `json:"date_of_visit"         db:"date_of_visit"`
	ReasonForVisit       string    `json:"reason_for_visit"      db:"reason_for_visit"`
	Diagnosis            string    `json:"diagnosis"             db:"diagnosis"`
	MedicationPrescribed []string  `json:"medication_prescribed" db:"medication_prescribed"`
	Notes                string    `json:"notes"                 db:"notes"`
	CreatedAt            time.Time `json:"created_at"            db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"            db:"updated_at"`
}

// QRCode represents a QR code for pet identification
type QRCode struct {
	ID             string         `json:"id"              db:"id"`
	PetID          string         `json:"pet_id"          db:"pet_id"`
	QRCodeData     string         `json:"qr_code_data"    db:"qr_code_data"`
	PublicURL      string         `json:"public_url"      db:"public_url"`
	EncodedContent EncodedContent `json:"encoded_content" db:"encoded_content"`
	IsActive       bool           `json:"is_active"       db:"is_active"`
	CreatedAt      time.Time      `json:"created_at"      db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"      db:"updated_at"`
}

// EncodedContent represents the JSON content stored in QR codes
type EncodedContent struct {
	PetName          string   `json:"pet_name"`
	PetType          string   `json:"pet_type"`
	OwnerName        string   `json:"owner_name"`
	OwnerPhone       string   `json:"owner_phone"`
	OwnerEmail       string   `json:"owner_email"`
	OwnerAddress     string   `json:"owner_address"`
	EmergencyContact string   `json:"emergency_contact,omitempty"`
	MedicalAlerts    []string `json:"medical_alerts,omitempty"`
	PublicProfileURL string   `json:"public_profile_url"`
}

// PublicPetProfile represents a public pet profile accessible via QR code
type PublicPetProfile struct {
	PetName          string   `json:"pet_name"`
	PetType          string   `json:"pet_type"`
	OwnerName        string   `json:"owner_name"`
	OwnerPhone       string   `json:"owner_phone"`
	OwnerEmail       string   `json:"owner_email"`
	OwnerAddress     string   `json:"owner_address"`
	EmergencyContact string   `json:"emergency_contact,omitempty"`
	MedicalAlerts    []string `json:"medical_alerts,omitempty"`
}

// Appointment represents a scheduled appointment
type Appointment struct {
	ID              string    `json:"id"               db:"id"`
	ClientID        string    `json:"client_id"        db:"client_id"`
	VeterinarianID  string    `json:"veterinarian_id"  db:"veterinarian_id"`
	PetID           string    `json:"pet_id"           db:"pet_id"`
	AppointmentDate time.Time `json:"appointment_date" db:"appointment_date"`
	DurationMinutes int       `json:"duration_minutes" db:"duration_minutes"`
	Reason          string    `json:"reason"           db:"reason"`
	Status          string    `json:"status"           db:"status"`
	Notes           string    `json:"notes"            db:"notes"`
	CreatedAt       time.Time `json:"created_at"       db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"       db:"updated_at"`
}

// TimeSlot represents an available appointment time slot
type TimeSlot struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Available bool      `json:"available"`
}

// Product represents a product in the e-commerce system
type Product struct {
	ID                     string            `json:"id"                       db:"id"`
	VeterinarianID         string            `json:"veterinarian_id"          db:"veterinarian_id"`
	Name                   string            `json:"name"                     db:"name"`
	Description            string            `json:"description"              db:"description"`
	Category               string            `json:"category"                 db:"category"`
	Price                  float64           `json:"price"                    db:"price"`
	StockQuantity          int               `json:"stock_quantity"           db:"stock_quantity"`
	SKU                    string            `json:"sku"                      db:"sku"`
	Brand                  string            `json:"brand"                    db:"brand"`
	Weight                 float64           `json:"weight"                   db:"weight"`
	Dimensions             ProductDimensions `json:"dimensions"               db:"dimensions"`
	IsPrescriptionRequired bool              `json:"is_prescription_required" db:"is_prescription_required"`
	IsActive               bool              `json:"is_active"                db:"is_active"`
	Images                 []string          `json:"images"                   db:"images"`
	CreatedAt              time.Time         `json:"created_at"               db:"created_at"`
	UpdatedAt              time.Time         `json:"updated_at"               db:"updated_at"`
}

// ProductDimensions represents product dimensions
type ProductDimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"`
}

// ProductFilters represents filters for product listing
type ProductFilters struct {
	Category       string  `json:"category,omitempty"`
	MinPrice       float64 `json:"min_price,omitempty"`
	MaxPrice       float64 `json:"max_price,omitempty"`
	Brand          string  `json:"brand,omitempty"`
	VeterinarianID string  `json:"veterinarian_id,omitempty"`
	Search         string  `json:"search,omitempty"`
	Limit          int     `json:"limit,omitempty"`
	Offset         int     `json:"offset,omitempty"`
}

// Order represents a purchase order
type Order struct {
	ID              string    `json:"id"               db:"id"`
	ClientID        string    `json:"client_id"        db:"client_id"`
	VeterinarianID  string    `json:"veterinarian_id"  db:"veterinarian_id"`
	TotalAmount     float64   `json:"total_amount"     db:"total_amount"`
	Status          string    `json:"status"           db:"status"`
	PaymentStatus   string    `json:"payment_status"   db:"payment_status"`
	PaymentMethod   string    `json:"payment_method"   db:"payment_method"`
	ShippingAddress string    `json:"shipping_address" db:"shipping_address"`
	DeliveryMethod  string    `json:"delivery_method"  db:"delivery_method"`
	Notes           string    `json:"notes"            db:"notes"`
	CreatedAt       time.Time `json:"created_at"       db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"       db:"updated_at"`
}

// OrderItem represents an individual item in an order
type OrderItem struct {
	ID         string    `json:"id"          db:"id"`
	OrderID    string    `json:"order_id"    db:"order_id"`
	ProductID  string    `json:"product_id"  db:"product_id"`
	Quantity   int       `json:"quantity"    db:"quantity"`
	UnitPrice  float64   `json:"unit_price"  db:"unit_price"`
	TotalPrice float64   `json:"total_price" db:"total_price"`
	CreatedAt  time.Time `json:"created_at"  db:"created_at"`
}

// NewPet creates a new Pet with generated ID and timestamps
func NewPet(
	ownerID, name, petType, breed string,
	dateOfBirth time.Time,
	weight float64,
) *Pet {
	now := time.Now()
	return &Pet{
		ID:          uuid.New().String(),
		OwnerID:     ownerID,
		Name:        name,
		Type:        petType,
		Breed:       breed,
		DateOfBirth: dateOfBirth,
		Weight:      weight,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewMedicalRecord creates a new MedicalRecord with generated ID and timestamps
func NewMedicalRecord(
	petID, veterinarianID, reasonForVisit, diagnosis string,
	dateOfVisit time.Time,
	medicationPrescribed []string,
	notes string,
) *MedicalRecord {
	now := time.Now()
	return &MedicalRecord{
		ID:                   uuid.New().String(),
		PetID:                petID,
		VeterinarianID:       veterinarianID,
		DateOfVisit:          dateOfVisit,
		ReasonForVisit:       reasonForVisit,
		Diagnosis:            diagnosis,
		MedicationPrescribed: medicationPrescribed,
		Notes:                notes,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

// NewQRCode creates a new QRCode with generated ID and timestamps
func NewQRCode(
	petID, qrCodeData, publicURL string,
	encodedContent EncodedContent,
) *QRCode {
	now := time.Now()
	return &QRCode{
		ID:             uuid.New().String(),
		PetID:          petID,
		QRCodeData:     qrCodeData,
		PublicURL:      publicURL,
		EncodedContent: encodedContent,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// NewAppointment creates a new Appointment with generated ID and timestamps
func NewAppointment(
	clientID, veterinarianID, petID string,
	appointmentDate time.Time,
	durationMinutes int,
	reason string,
) *Appointment {
	now := time.Now()
	return &Appointment{
		ID:              uuid.New().String(),
		ClientID:        clientID,
		VeterinarianID:  veterinarianID,
		PetID:           petID,
		AppointmentDate: appointmentDate,
		DurationMinutes: durationMinutes,
		Reason:          reason,
		Status:          "scheduled",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// NewProduct creates a new Product with generated ID and timestamps
func NewProduct(
	veterinarianID, name, description, category string,
	price float64,
) *Product {
	now := time.Now()
	return &Product{
		ID:                     uuid.New().String(),
		VeterinarianID:         veterinarianID,
		Name:                   name,
		Description:            description,
		Category:               category,
		Price:                  price,
		StockQuantity:          0,
		IsPrescriptionRequired: false,
		IsActive:               true,
		CreatedAt:              now,
		UpdatedAt:              now,
	}
}

// NewOrder creates a new Order with generated ID and timestamps
func NewOrder(clientID, veterinarianID string, totalAmount float64) *Order {
	now := time.Now()
	return &Order{
		ID:             uuid.New().String(),
		ClientID:       clientID,
		VeterinarianID: veterinarianID,
		TotalAmount:    totalAmount,
		Status:         "pending",
		PaymentStatus:  "pending",
		DeliveryMethod: "pickup",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// NewOrderItem creates a new OrderItem with generated ID and timestamps
func NewOrderItem(
	orderID, productID string,
	quantity int,
	unitPrice float64,
) *OrderItem {
	now := time.Now()
	totalPrice := float64(quantity) * unitPrice
	return &OrderItem{
		ID:         uuid.New().String(),
		OrderID:    orderID,
		ProductID:  productID,
		Quantity:   quantity,
		UnitPrice:  unitPrice,
		TotalPrice: totalPrice,
		CreatedAt:  now,
	}
}
