// Package handlers_test contains tests for HTTP handlers
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pet-mgt/backend/internal/middleware"
	"pet-mgt/backend/internal/store"
	"testing"
	"time"
)

// MockDatabase implements the Database interface for testing
type MockDatabase struct {
	users map[string]*store.User
	pets  map[string]*store.Pet
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		users: make(map[string]*store.User),
		pets:  make(map[string]*store.Pet),
	}
}

func (m *MockDatabase) GetUserByID(
	ctx context.Context,
	userID string,
) (*store.User, error) {
	if user, exists := m.users[userID]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *MockDatabase) CreateUser(ctx context.Context, user *store.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockDatabase) UpdateUser(ctx context.Context, user *store.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockDatabase) DeleteUser(ctx context.Context, userID string) error {
	delete(m.users, userID)
	return nil
}

func (m *MockDatabase) ListUsers(
	ctx context.Context,
	limit, offset int,
) ([]store.User, error) {
	var users []store.User
	for _, user := range m.users {
		users = append(users, *user)
	}
	return users, nil
}

func (m *MockDatabase) GetPetsByUserID(
	ctx context.Context,
	userID string,
) ([]store.Pet, error) {
	var pets []store.Pet
	for _, pet := range m.pets {
		if pet.OwnerID == userID {
			pets = append(pets, *pet)
		}
	}
	return pets, nil
}

func (m *MockDatabase) GetPetByID(
	ctx context.Context,
	petID string,
) (*store.Pet, error) {
	if pet, exists := m.pets[petID]; exists {
		return pet, nil
	}
	return nil, nil
}

func (m *MockDatabase) CreatePet(ctx context.Context, pet *store.Pet) error {
	m.pets[pet.ID] = pet
	return nil
}

func (m *MockDatabase) UpdatePet(ctx context.Context, pet *store.Pet) error {
	m.pets[pet.ID] = pet
	return nil
}

func (m *MockDatabase) DeletePet(ctx context.Context, petID string) error {
	delete(m.pets, petID)
	return nil
}

func (m *MockDatabase) GetMedicalRecordsByPetID(
	ctx context.Context,
	petID string,
) ([]store.MedicalRecord, error) {
	return []store.MedicalRecord{}, nil
}

func (m *MockDatabase) GetMedicalRecordByID(
	ctx context.Context,
	recordID string,
) (*store.MedicalRecord, error) {
	return nil, nil
}

func (m *MockDatabase) CreateMedicalRecord(
	ctx context.Context,
	record *store.MedicalRecord,
) error {
	return nil
}

func (m *MockDatabase) UpdateMedicalRecord(
	ctx context.Context,
	record *store.MedicalRecord,
) error {
	return nil
}

func (m *MockDatabase) DeleteMedicalRecord(ctx context.Context, recordID string) error {
	return nil
}

func (m *MockDatabase) Ping(ctx context.Context) error {
	return nil
}

func (m *MockDatabase) Close() error {
	return nil
}

// Add these methods to MockDatabase:
func (m *MockDatabase) CreateClient(ctx context.Context, client *store.Client) error {
	// For test, just store as a User
	m.users[client.ID] = &store.User{
		ID:    client.ID,
		Email: client.Email,
		Role:  client.Role,
	}
	return nil
}

func (m *MockDatabase) CreateVeterinarian(
	ctx context.Context,
	vet *store.Veterinarian,
) error {
	m.users[vet.ID] = &store.User{ID: vet.ID, Email: vet.Email, Role: vet.Role}
	return nil
}

func (m *MockDatabase) UpdateClient(ctx context.Context, client *store.Client) error {
	m.users[client.ID] = &store.User{
		ID:    client.ID,
		Email: client.Email,
		Role:  client.Role,
	}
	return nil
}

func (m *MockDatabase) UpdateVeterinarian(
	ctx context.Context,
	vet *store.Veterinarian,
) error {
	m.users[vet.ID] = &store.User{ID: vet.ID, Email: vet.Email, Role: vet.Role}
	return nil
}

func (m *MockDatabase) GetClientByID(
	ctx context.Context,
	clientID string,
) (*store.Client, error) {
	if u, ok := m.users[clientID]; ok {
		return &store.Client{ID: u.ID, Email: u.Email, Role: u.Role}, nil
	}
	return nil, nil
}

func (m *MockDatabase) GetVeterinarianByID(
	ctx context.Context,
	vetID string,
) (*store.Veterinarian, error) {
	if u, ok := m.users[vetID]; ok {
		return &store.Veterinarian{ID: u.ID, Email: u.Email, Role: u.Role}, nil
	}
	return nil, nil
}

// QR Code operations (stub implementations for testing)
func (m *MockDatabase) GetQRCodeByPetID(
	ctx context.Context,
	petID string,
) (*store.QRCode, error) {
	return nil, nil
}

func (m *MockDatabase) GetQRCodeByPublicURL(
	ctx context.Context,
	publicURL string,
) (*store.QRCode, error) {
	return nil, nil
}

func (m *MockDatabase) CreateQRCode(
	ctx context.Context,
	qrCode *store.QRCode,
) error {
	return nil
}

func (m *MockDatabase) UpdateQRCode(
	ctx context.Context,
	qrCode *store.QRCode,
) error {
	return nil
}

func (m *MockDatabase) DeleteQRCode(
	ctx context.Context,
	qrCodeID string,
) error {
	return nil
}

func (m *MockDatabase) GetPublicPetProfile(
	ctx context.Context,
	publicURL string,
) (*store.PublicPetProfile, error) {
	return nil, nil
}

// Appointment operations (stub implementations for testing)
func (m *MockDatabase) GetAppointmentsByClientID(
	ctx context.Context,
	clientID string,
) ([]store.Appointment, error) {
	return []store.Appointment{}, nil
}

func (m *MockDatabase) GetAppointmentsByVeterinarianID(
	ctx context.Context,
	vetID string,
) ([]store.Appointment, error) {
	return []store.Appointment{}, nil
}

func (m *MockDatabase) GetAppointmentByID(
	ctx context.Context,
	appointmentID string,
) (*store.Appointment, error) {
	return nil, nil
}

func (m *MockDatabase) CreateAppointment(
	ctx context.Context,
	appointment *store.Appointment,
) error {
	return nil
}

func (m *MockDatabase) UpdateAppointment(
	ctx context.Context,
	appointment *store.Appointment,
) error {
	return nil
}

func (m *MockDatabase) DeleteAppointment(
	ctx context.Context,
	appointmentID string,
) error {
	return nil
}

func (m *MockDatabase) GetAvailableAppointmentSlots(
	ctx context.Context,
	vetID string,
	date time.Time,
) ([]store.TimeSlot, error) {
	return []store.TimeSlot{}, nil
}

// Product operations (stub implementations for testing)
func (m *MockDatabase) GetProductsByVeterinarianID(
	ctx context.Context,
	vetID string,
) ([]store.Product, error) {
	return []store.Product{}, nil
}

func (m *MockDatabase) GetProductByID(
	ctx context.Context,
	productID string,
) (*store.Product, error) {
	return nil, nil
}

func (m *MockDatabase) CreateProduct(
	ctx context.Context,
	product *store.Product,
) error {
	return nil
}

func (m *MockDatabase) UpdateProduct(
	ctx context.Context,
	product *store.Product,
) error {
	return nil
}

func (m *MockDatabase) DeleteProduct(
	ctx context.Context,
	productID string,
) error {
	return nil
}

func (m *MockDatabase) ListProducts(
	ctx context.Context,
	filters store.ProductFilters,
) ([]store.Product, error) {
	return []store.Product{}, nil
}

func (m *MockDatabase) UpdateProductStock(
	ctx context.Context,
	productID string,
	quantity int,
) error {
	return nil
}

// Order operations (stub implementations for testing)
func (m *MockDatabase) GetOrdersByClientID(
	ctx context.Context,
	clientID string,
) ([]store.Order, error) {
	return []store.Order{}, nil
}

func (m *MockDatabase) GetOrdersByVeterinarianID(
	ctx context.Context,
	vetID string,
) ([]store.Order, error) {
	return []store.Order{}, nil
}

func (m *MockDatabase) GetOrderByID(
	ctx context.Context,
	orderID string,
) (*store.Order, error) {
	return nil, nil
}

func (m *MockDatabase) CreateOrder(
	ctx context.Context,
	order *store.Order,
) error {
	return nil
}

func (m *MockDatabase) UpdateOrderStatus(
	ctx context.Context,
	orderID string,
	status string,
) error {
	return nil
}

func (m *MockDatabase) GetOrderItems(
	ctx context.Context,
	orderID string,
) ([]store.OrderItem, error) {
	return []store.OrderItem{}, nil
}

func (m *MockDatabase) CreateOrderItem(
	ctx context.Context,
	item *store.OrderItem,
) error {
	return nil
}

// Helper function to create a request with context
func createRequestWithContext(
	method, path string,
	body any,
	user *middleware.UserClaims,
) *http.Request {
	var req *http.Request
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
	return req.WithContext(ctx)
}

// TestGetUserProfile tests the GetUserProfile handler
func TestGetUserProfile(t *testing.T) {
	user := &middleware.UserClaims{
		Sub:   "test-user-id",
		Email: "test@example.com",
		Role:  "client",
	}

	// Create handler with mock database
	mockDB := NewMockDatabase()
	userHandler := NewUserHandler(mockDB)

	req := createRequestWithContext("GET", "/api/v1/profile", nil, user)
	w := httptest.NewRecorder()

	userHandler.GetUserProfile(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if response["user_id"] != user.Sub {
		t.Errorf("Expected user_id %s, got %s", user.Sub, response["user_id"])
	}
}

// TestCreateUser tests the CreateUser handler
func TestCreateUser(t *testing.T) {
	user := &middleware.UserClaims{
		Sub:   "test-user-id",
		Email: "test@example.com",
		Role:  "client",
	}

	// Create handler with mock database
	mockDB := NewMockDatabase()
	userHandler := NewUserHandler(mockDB)

	// Create test request
	newUser := map[string]any{
		"name":    "John Doe",
		"email":   "john@example.com",
		"phone":   "123-456-7890",
		"address": "123 Main St",
		"role":    "client",
	}

	req := createRequestWithContext("POST", "/api/v1/users", newUser, user)
	w := httptest.NewRecorder()

	userHandler.CreateUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if !response["success"].(bool) {
		t.Errorf("Expected success to be true")
	}
}
