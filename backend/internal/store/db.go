// Package store/db.go contains database functions
package store

import (
	"context"
	"fmt"
	"pet-mgt/backend/internal/config"
	"time"

	"github.com/supabase-community/supabase-go"
)

type SupabaseService struct {
	client *supabase.Client
	config *config.Config
}

// NewSupabaseService creates a new SupabaseService
func NewSupabaseService(cfg *config.Config) (*SupabaseService, error) {
	client, err := supabase.NewClient(
		cfg.SupabaseURL,
		cfg.SupabaseServiceKey,
		&supabase.ClientOptions{},
	)
	if err != nil {
		return nil, err
	}

	return &SupabaseService{
		config: cfg,
		client: client,
	}, nil
}

// Ping checks if the database connection is alive
func (s *SupabaseService) Ping(ctx context.Context) error {
	// Simple health check - the client connection is tested during initialization
	// If we got here, the connection is working
	return nil
}

// GetUserByID retrieves a user by their ID
func (s *SupabaseService) GetUserByID(
	ctx context.Context,
	userID string,
) (*User, error) {
	// Try to get from clients table first
	var client Client
	_, err := s.client.From("clients").
		Select("*", "", false).
		Eq("id", userID).
		Single().
		ExecuteTo(&client)
	if err == nil {
		return &User{
			ID:    client.ID,
			Email: client.Email,
			Role:  client.Role,
		}, nil
	}

	// Try to get from veterinarians table
	var vet Veterinarian
	_, err = s.client.From("veterinarians").
		Select("*", "", false).
		Eq("id", userID).
		Single().
		ExecuteTo(&vet)
	if err == nil {
		return &User{
			ID:    vet.ID,
			Email: vet.Email,
			Role:  vet.Role,
		}, nil
	}

	return nil, fmt.Errorf("user not found")
}

// CreateUser creates a new user profile
func (s *SupabaseService) CreateUser(ctx context.Context, user *User) error {
	// This method will be implemented based on the user role
	// NOTE: For now, return an error as we need to know the specific user type
	return fmt.Errorf("use CreateClient or CreateVeterinarian instead")
}

// CreateClient creates a new client profile
func (s *SupabaseService) CreateClient(ctx context.Context, client *Client) error {
	_, _, err := s.client.From("clients").Insert(client, false, "", "", "").Execute()
	return err
}

// CreateVeterinarian creates a new veterinarian profile
func (s *SupabaseService) CreateVeterinarian(
	ctx context.Context,
	vet *Veterinarian,
) error {
	_, _, err := s.client.From("veterinarians").Insert(vet, false, "", "", "").Execute()
	return err
}

// UpdateUser updates a user profile
func (s *SupabaseService) UpdateUser(ctx context.Context, user *User) error {
	// This method will be implemented based on the user role
	return fmt.Errorf("use UpdateClient or UpdateVeterinarian instead")
}

// UpdateClient updates a client profile
func (s *SupabaseService) UpdateClient(ctx context.Context, client *Client) error {
	_, _, err := s.client.From("clients").
		Update(client, "", "").
		Eq("id", client.ID).
		Execute()
	return err
}

// UpdateVeterinarian updates a veterinarian profile
func (s *SupabaseService) UpdateVeterinarian(
	ctx context.Context,
	vet *Veterinarian,
) error {
	_, _, err := s.client.From("veterinarians").
		Update(vet, "", "").
		Eq("id", vet.ID).
		Execute()
	return err
}

// DeleteUser deletes a user profile
func (s *SupabaseService) DeleteUser(ctx context.Context, userID string) error {
	// Try to delete from clients table first
	_, _, err := s.client.From("clients").Delete("", "").Eq("id", userID).Execute()
	if err == nil {
		return nil
	}

	// Try to delete from veterinarians table
	_, _, err = s.client.From("veterinarians").Delete("", "").Eq("id", userID).Execute()
	return err
}

// ListUsers lists all users with pagination
func (s *SupabaseService) ListUsers(
	ctx context.Context,
	limit, offset int,
) ([]User, error) {
	var users []User

	// Get clients
	var clients []Client
	_, err := s.client.From("clients").
		Select("*", "", false).
		Range(offset, offset+limit-1, "").
		ExecuteTo(&clients)
	if err != nil {
		return nil, err
	}

	// Get veterinarians
	var vets []Veterinarian
	_, err = s.client.From("veterinarians").
		Select("*", "", false).
		Range(offset, offset+limit-1, "").
		ExecuteTo(&vets)
	if err != nil {
		return nil, err
	}

	// Convert to User interface
	for _, client := range clients {
		users = append(users, User{
			ID:    client.ID,
			Email: client.Email,
			Role:  client.Role,
		})
	}

	for _, vet := range vets {
		users = append(users, User{
			ID:    vet.ID,
			Email: vet.Email,
			Role:  vet.Role,
		})
	}

	return users, nil
}

// GetPetsByUserID retrieves all pets for a specific user
func (s *SupabaseService) GetPetsByUserID(
	ctx context.Context,
	userID string,
) ([]Pet, error) {
	var pets []Pet

	_, err := s.client.From("pets").
		Select("*", "", false).
		Eq("owner_id", userID).
		ExecuteTo(&pets)
	if err != nil {
		return nil, err
	}

	return pets, err
}

// GetPetByID retrieves a pet by ID
func (s *SupabaseService) GetPetByID(ctx context.Context, petID string) (*Pet, error) {
	var pet Pet
	_, err := s.client.From("pets").
		Select("*", "", false).
		Eq("id", petID).
		Single().
		ExecuteTo(&pet)
	if err != nil {
		return nil, err
	}
	return &pet, nil
}

// CreatePet creates a new pet
func (s *SupabaseService) CreatePet(ctx context.Context, pet *Pet) error {
	_, _, err := s.client.From("pets").Insert(pet, false, "", "", "").Execute()
	return err
}

// UpdatePet updates an existing pet
func (s *SupabaseService) UpdatePet(ctx context.Context, pet *Pet) error {
	_, _, err := s.client.From("pets").Update(pet, "", "").Eq("id", pet.ID).Execute()
	return err
}

// DeletePet deletes a pet by ID
func (s *SupabaseService) DeletePet(ctx context.Context, petID string) error {
	_, _, err := s.client.From("pets").Delete("", "").Eq("id", petID).Execute()
	return err
}

// GetMedicalRecordsByPetID retrieves all medical records for a pet
func (s *SupabaseService) GetMedicalRecordsByPetID(
	ctx context.Context,
	petID string,
) ([]MedicalRecord, error) {
	var records []MedicalRecord
	_, err := s.client.From("medical_records").
		Select("*", "", false).
		Eq("pet_id", petID).
		ExecuteTo(&records)
	return records, err
}

// GetMedicalRecordByID retrieves a specific medical record
func (s *SupabaseService) GetMedicalRecordByID(
	ctx context.Context,
	recordID string,
) (*MedicalRecord, error) {
	var record MedicalRecord
	_, err := s.client.From("medical_records").
		Select("*", "", false).
		Eq("id", recordID).
		Single().
		ExecuteTo(&record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// CreateMedicalRecord creates a new medical record
func (s *SupabaseService) CreateMedicalRecord(
	ctx context.Context,
	record *MedicalRecord,
) error {
	_, _, err := s.client.From("medical_records").
		Insert(record, false, "", "", "").
		Execute()
	return err
}

// UpdateMedicalRecord updates an existing medical record
func (s *SupabaseService) UpdateMedicalRecord(
	ctx context.Context,
	record *MedicalRecord,
) error {
	_, _, err := s.client.From("medical_records").
		Update(record, "", "").
		Eq("id", record.ID).
		Execute()
	return err
}

// DeleteMedicalRecord deletes a medical record by ID
func (s *SupabaseService) DeleteMedicalRecord(
	ctx context.Context,
	recordID string,
) error {
	_, _, err := s.client.From("medical_records").
		Delete("", "").
		Eq("id", recordID).
		Execute()
	return err
}

// GetClientByID retrieves a client by ID
func (s *SupabaseService) GetClientByID(
	ctx context.Context,
	clientID string,
) (*Client, error) {
	var client Client
	_, err := s.client.From("clients").
		Select("*", "", false).
		Eq("id", clientID).
		Single().
		ExecuteTo(&client)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// GetVeterinarianByID retrieves a veterinarian by ID
func (s *SupabaseService) GetVeterinarianByID(
	ctx context.Context,
	vetID string,
) (*Veterinarian, error) {
	var vet Veterinarian
	_, err := s.client.From("veterinarians").
		Select("*", "", false).
		Eq("id", vetID).
		Single().
		ExecuteTo(&vet)
	if err != nil {
		return nil, err
	}
	return &vet, nil
}

// Close performs cleanup operations
func (s *SupabaseService) Close() error {
	// Supabase client uses HTTP connections, no explicit cleanup needed
	// Future: close any connection pools or background workers here
	return nil
}

// QR Code operations

// GetQRCodeByPetID retrieves QR code by pet ID
func (s *SupabaseService) GetQRCodeByPetID(
	ctx context.Context,
	petID string,
) (*QRCode, error) {
	var qrCode QRCode
	_, err := s.client.From("qr_codes").
		Select("*", "", false).
		Eq("pet_id", petID).
		Eq("is_active", "true").
		Single().
		ExecuteTo(&qrCode)
	if err != nil {
		return nil, err
	}
	return &qrCode, nil
}

// GetQRCodeByPublicURL retrieves QR code by public URL
func (s *SupabaseService) GetQRCodeByPublicURL(
	ctx context.Context,
	publicURL string,
) (*QRCode, error) {
	var qrCode QRCode
	_, err := s.client.From("qr_codes").
		Select("*", "", false).
		Eq("public_url", publicURL).
		Eq("is_active", "true").
		Single().
		ExecuteTo(&qrCode)
	if err != nil {
		return nil, err
	}
	return &qrCode, nil
}

// CreateQRCode creates a new QR code
func (s *SupabaseService) CreateQRCode(ctx context.Context, qrCode *QRCode) error {
	_, _, err := s.client.From("qr_codes").Insert(qrCode, false, "", "", "").Execute()
	return err
}

// UpdateQRCode updates an existing QR code
func (s *SupabaseService) UpdateQRCode(ctx context.Context, qrCode *QRCode) error {
	_, _, err := s.client.From("qr_codes").
		Update(qrCode, "", "").
		Eq("id", qrCode.ID).
		Execute()
	return err
}

// DeleteQRCode deletes a QR code by ID
func (s *SupabaseService) DeleteQRCode(ctx context.Context, qrCodeID string) error {
	_, _, err := s.client.From("qr_codes").Delete("", "").Eq("id", qrCodeID).Execute()
	return err
}

// GetPublicPetProfile retrieves public pet profile via QR code URL
func (s *SupabaseService) GetPublicPetProfile(
	ctx context.Context,
	publicURL string,
) (*PublicPetProfile, error) {
	var qrCode QRCode
	_, err := s.client.From("qr_codes").
		Select("*", "", false).
		Eq("public_url", publicURL).
		Eq("is_active", "true").
		Single().
		ExecuteTo(&qrCode)
	if err != nil {
		return nil, err
	}

	// Convert to public profile
	// Fetch pet details
	var pet Pet
	_, err = s.client.From("pets").
		Select("*", "", false).
		Eq("id", qrCode.PetID).
		Single().
		ExecuteTo(&pet)
	if err != nil {
		return nil, err
	}

	// Fetch medical records and project to public-friendly shape
	var records []MedicalRecord
	_, err = s.client.From("medical_records").
		Select("*", "", false).
		Eq("pet_id", qrCode.PetID).
		ExecuteTo(&records)
	if err != nil {
		return nil, err
	}
	publicRecords := make([]PublicMedicalRecord, 0, len(records))
	for _, r := range records {
		publicRecords = append(publicRecords, PublicMedicalRecord{
			DateOfVisit:          r.DateOfVisit.Format("2006-01-02"),
			ReasonForVisit:       r.ReasonForVisit,
			Diagnosis:            r.Diagnosis,
			MedicationPrescribed: r.MedicationPrescribed,
		})
	}

	// Build enriched public profile
	profile := &PublicPetProfile{
		PetName:          qrCode.EncodedContent.PetName,
		PetType:          qrCode.EncodedContent.PetType,
		Breed:            pet.Breed,
		DateOfBirth:      pet.DateOfBirth,
		Weight:           pet.Weight,
		OwnerName:        qrCode.EncodedContent.OwnerName,
		OwnerPhone:       qrCode.EncodedContent.OwnerPhone,
		OwnerEmail:       qrCode.EncodedContent.OwnerEmail,
		OwnerAddress:     qrCode.EncodedContent.OwnerAddress,
		EmergencyContact: qrCode.EncodedContent.EmergencyContact,
		MedicalAlerts:    qrCode.EncodedContent.MedicalAlerts,
		MedicalRecords:   publicRecords,
	}

	return profile, nil
}

// Appointment operations

// GetAppointmentsByClientID retrieves appointments for a client
func (s *SupabaseService) GetAppointmentsByClientID(
	ctx context.Context,
	clientID string,
) ([]Appointment, error) {
	var appointments []Appointment
	_, err := s.client.From("appointments").
		Select("*", "", false).
		Eq("client_id", clientID).
		ExecuteTo(&appointments)
	return appointments, err
}

// GetAppointmentsByVeterinarianID retrieves appointments for a veterinarian
func (s *SupabaseService) GetAppointmentsByVeterinarianID(
	ctx context.Context,
	vetID string,
) ([]Appointment, error) {
	var appointments []Appointment
	_, err := s.client.From("appointments").
		Select("*", "", false).
		Eq("veterinarian_id", vetID).
		ExecuteTo(&appointments)
	return appointments, err
}

// GetAppointmentByID retrieves a specific appointment
func (s *SupabaseService) GetAppointmentByID(
	ctx context.Context,
	appointmentID string,
) (*Appointment, error) {
	var appointment Appointment
	_, err := s.client.From("appointments").
		Select("*", "", false).
		Eq("id", appointmentID).
		Single().
		ExecuteTo(&appointment)
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

// CreateAppointment creates a new appointment
func (s *SupabaseService) CreateAppointment(
	ctx context.Context,
	appointment *Appointment,
) error {
	_, _, err := s.client.From("appointments").
		Insert(appointment, false, "", "", "").
		Execute()
	return err
}

// UpdateAppointment updates an existing appointment
func (s *SupabaseService) UpdateAppointment(
	ctx context.Context,
	appointment *Appointment,
) error {
	_, _, err := s.client.From("appointments").
		Update(appointment, "", "").
		Eq("id", appointment.ID).
		Execute()
	return err
}

// DeleteAppointment deletes an appointment by ID
func (s *SupabaseService) DeleteAppointment(
	ctx context.Context,
	appointmentID string,
) error {
	_, _, err := s.client.From("appointments").
		Delete("", "").
		Eq("id", appointmentID).
		Execute()
	return err
}

// GetAvailableAppointmentSlots retrieves available time slots for a veterinarian
func (s *SupabaseService) GetAvailableAppointmentSlots(
	ctx context.Context,
	vetID string,
	date time.Time,
) ([]TimeSlot, error) {
	// // Interpret working hours in clinic timezone (default Asia/Manila)
	// loc, _ := time.LoadLocation("Asia/Manila")

	// Fetch veterinarian working hours
	vet, err := s.GetVeterinarianByID(ctx, vetID)
	if err != nil {
		return nil, err
	}

	// Determine windows for the requested day
	weekdayKey := weekdayToKey(date.Weekday())
	type window struct {
		start time.Time
		end   time.Time
	}
	var windows []window
	for _, wh := range vet.AvailableHours {
		if normalizeDayKey(wh.DayOfWeek) != weekdayKey {
			continue
		}
		startParsed, err1 := time.Parse("15:04", wh.Start)
		endParsed, err2 := time.Parse("15:04", wh.End)
		if err1 != nil || err2 != nil {
			continue
		}
		ws := time.Date(date.Year(), date.Month(), date.Day(), startParsed.Hour(), startParsed.Minute(), 0, 0, date.Location())
		we := time.Date(date.Year(), date.Month(), date.Day(), endParsed.Hour(), endParsed.Minute(), 0, 0, date.Location())
		if we.After(ws) {
			windows = append(windows, window{start: ws, end: we})
		}
	}

	if len(windows) == 0 {
		return []TimeSlot{}, nil
	}

	// Fetch existing appointments for that day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var appts []Appointment
	_, err = s.client.From("appointments").
		Select("*", "", false).
		Eq("veterinarian_id", vetID).
		Gte("appointment_date", startOfDay.Format(time.RFC3339)).
		Lt("appointment_date", endOfDay.Format(time.RFC3339)).
		ExecuteTo(&appts)
	if err != nil {
		return nil, err
	}

	// Ignore cancelled appointments
	activeAppts := make([]Appointment, 0, len(appts))
	for _, a := range appts {
		if a.Status != "cancelled" {
			activeAppts = append(activeAppts, a)
		}
	}

	// Build 30-minute slots within the windows and mark as available if no overlap
	const slotMinutes = 30
	var slots []TimeSlot
	for _, win := range windows {
		for ts := win.start; ts.Add(time.Duration(slotMinutes)*time.Minute).Before(win.end) || ts.Add(time.Duration(slotMinutes)*time.Minute).Equal(win.end); ts = ts.Add(time.Duration(slotMinutes) * time.Minute) {
			te := ts.Add(time.Duration(slotMinutes) * time.Minute)
			available := true
			for _, a := range activeAppts {
				as := a.AppointmentDate
				ae := as.Add(time.Duration(a.DurationMinutes) * time.Minute)
				if intervalsOverlap(ts, te, as, ae) {
					available = false
					break
				}
			}
			slots = append(slots, TimeSlot{StartTime: ts, EndTime: te, Available: available})
		}
	}

	return slots, nil
}

// intervalsOverlap returns true if [s1,e1) overlaps [s2,e2)
func intervalsOverlap(s1, e1, s2, e2 time.Time) bool {
	return s1.Before(e2) && s2.Before(e1)
}

// weekdayToKey returns canonical three-letter weekday key
func weekdayToKey(w time.Weekday) string {
	switch w {
	case time.Monday:
		return "Mon"
	case time.Tuesday:
		return "Tue"
	case time.Wednesday:
		return "Wed"
	case time.Thursday:
		return "Thu"
	case time.Friday:
		return "Fri"
	case time.Saturday:
		return "Sat"
	default:
		return "Sun"
	}
}

// normalizeDayKey maps arbitrary inputs to canonical three-letter weekday key
func normalizeDayKey(s string) string {
	switch s {
	case "Mon", "Monday", "monday", "mon":
		return "Mon"
	case "Tue", "Tues", "Tuesday", "tuesday", "tue", "tues":
		return "Tue"
	case "Wed", "Wednesday", "wednesday", "wed":
		return "Wed"
	case "Thu", "Thur", "Thurs", "Thursday", "thursday", "thu", "thur", "thurs":
		return "Thu"
	case "Fri", "Friday", "friday", "fri":
		return "Fri"
	case "Sat", "Saturday", "saturday", "sat":
		return "Sat"
	case "Sun", "Sunday", "sunday", "sun":
		return "Sun"
	default:
		return s
	}
}

// Product operations

// GetProductsByVeterinarianID retrieves products for a veterinarian
func (s *SupabaseService) GetProductsByVeterinarianID(
	ctx context.Context,
	vetID string,
) ([]Product, error) {
	var products []Product
	_, err := s.client.From("products").
		Select("*", "", false).
		Eq("veterinarian_id", vetID).
		Eq("is_active", "true").
		ExecuteTo(&products)
	return products, err
}

// GetProductByID retrieves a specific product
func (s *SupabaseService) GetProductByID(
	ctx context.Context,
	productID string,
) (*Product, error) {
	var product Product
	_, err := s.client.From("products").
		Select("*", "", false).
		Eq("id", productID).
		Single().
		ExecuteTo(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// CreateProduct creates a new product
func (s *SupabaseService) CreateProduct(ctx context.Context, product *Product) error {
	_, _, err := s.client.From("products").Insert(product, false, "", "", "").Execute()
	return err
}

// UpdateProduct updates an existing product
func (s *SupabaseService) UpdateProduct(ctx context.Context, product *Product) error {
	_, _, err := s.client.From("products").
		Update(product, "", "").
		Eq("id", product.ID).
		Execute()
	return err
}

// DeleteProduct deletes a product by ID
func (s *SupabaseService) DeleteProduct(ctx context.Context, productID string) error {
	_, _, err := s.client.From("products").Delete("", "").Eq("id", productID).Execute()
	return err
}

// ListProducts retrieves products with filtering
func (s *SupabaseService) ListProducts(
	ctx context.Context,
	filters ProductFilters,
) ([]Product, error) {
	var products []Product

	query := s.client.From("products").Select("*", "", false).Eq("is_active", "true")

	if filters.Category != "" {
		query = query.Eq("category", filters.Category)
	}
	if filters.VeterinarianID != "" {
		query = query.Eq("veterinarian_id", filters.VeterinarianID)
	}

	// Use Range for pagination like other methods in this file
	if filters.Limit > 0 && filters.Offset >= 0 {
		query = query.Range(filters.Offset, filters.Offset+filters.Limit-1, "")
	}

	_, err := query.ExecuteTo(&products)
	return products, err
}

// UpdateProductStock updates product stock quantity
func (s *SupabaseService) UpdateProductStock(
	ctx context.Context,
	productID string,
	quantity int,
) error {
	updateData := map[string]any{
		"stock_quantity": quantity,
	}
	_, _, err := s.client.From("products").
		Update(updateData, "", "").
		Eq("id", productID).
		Execute()
	return err
}

// Order operations

// GetOrdersByClientID retrieves orders for a client
func (s *SupabaseService) GetOrdersByClientID(
	ctx context.Context,
	clientID string,
) ([]Order, error) {
	var orders []Order
	_, err := s.client.From("orders").
		Select("*", "", false).
		Eq("client_id", clientID).
		ExecuteTo(&orders)
	return orders, err
}

// GetOrdersByVeterinarianID retrieves orders for a veterinarian
func (s *SupabaseService) GetOrdersByVeterinarianID(
	ctx context.Context,
	vetID string,
) ([]Order, error) {
	var orders []Order
	_, err := s.client.From("orders").
		Select("*", "", false).
		Eq("veterinarian_id", vetID).
		ExecuteTo(&orders)
	return orders, err
}

// GetOrderByID retrieves a specific order
func (s *SupabaseService) GetOrderByID(
	ctx context.Context,
	orderID string,
) (*Order, error) {
	var order Order
	_, err := s.client.From("orders").
		Select("*", "", false).
		Eq("id", orderID).
		Single().
		ExecuteTo(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// CreateOrder creates a new order
func (s *SupabaseService) CreateOrder(ctx context.Context, order *Order) error {
	_, _, err := s.client.From("orders").Insert(order, false, "", "", "").Execute()
	return err
}

// UpdateOrderStatus updates the status of an order
func (s *SupabaseService) UpdateOrderStatus(
	ctx context.Context,
	orderID string,
	status string,
) error {
	updateData := map[string]any{
		"status": status,
	}
	_, _, err := s.client.From("orders").
		Update(updateData, "", "").
		Eq("id", orderID).
		Execute()
	return err
}

// GetOrderItems retrieves items for an order
func (s *SupabaseService) GetOrderItems(
	ctx context.Context,
	orderID string,
) ([]OrderItem, error) {
	var items []OrderItem
	_, err := s.client.From("order_items").
		Select("*", "", false).
		Eq("order_id", orderID).
		ExecuteTo(&items)
	return items, err
}

// CreateOrderItem creates a new order item
func (s *SupabaseService) CreateOrderItem(ctx context.Context, item *OrderItem) error {
	_, _, err := s.client.From("order_items").Insert(item, false, "", "", "").Execute()
	return err
}
