// Package handlers contains product management handlers
package handlers

import (
	"encoding/json"
	"net/http"
	"pet-mgt/backend/internal/middleware"
	"pet-mgt/backend/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// ProductHandler handles product operations
type ProductHandler struct {
	db store.Database
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(db store.Database) *ProductHandler {
	return &ProductHandler{db: db}
}

// CreateProduct creates a new product (veterinarians only)
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Only veterinarians and admins can create products
	if user.Role != "veterinarian" && user.Role != "admin" {
		ErrorResponse(w, http.StatusForbidden, "Only veterinarians can create products")
		return
	}

	// Parse request body
	var req struct {
		Name                   string                  `json:"name"`
		Description            string                  `json:"description"`
		Category               string                  `json:"category"`
		Price                  float64                 `json:"price"`
		StockQuantity          int                     `json:"stock_quantity"`
		SKU                    string                  `json:"sku,omitempty"`
		Brand                  string                  `json:"brand,omitempty"`
		Weight                 float64                 `json:"weight,omitempty"`
		Dimensions             store.ProductDimensions `json:"dimensions"`
		IsPrescriptionRequired bool                    `json:"is_prescription_required"`
		Images                 []string                `json:"images,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Name == "" || req.Category == "" || req.Price <= 0 {
		ErrorResponse(
			w,
			http.StatusBadRequest,
			"Missing or invalid required fields (name, category, price)",
		)
		return
	}

	// Create product
	veterinarianID := user.Sub
	if user.Role == "admin" && r.URL.Query().Get("veterinarian_id") != "" {
		veterinarianID = r.URL.Query().Get("veterinarian_id")
	}

	product := store.NewProduct(
		veterinarianID,
		req.Name,
		req.Description,
		req.Category,
		req.Price,
	)
	product.StockQuantity = req.StockQuantity
	product.SKU = req.SKU
	product.Brand = req.Brand
	product.Weight = req.Weight
	product.Dimensions = req.Dimensions
	product.IsPrescriptionRequired = req.IsPrescriptionRequired
	product.Images = req.Images

	if err := h.db.CreateProduct(r.Context(), product); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	SuccessResponse(w, product)
}

// GetProducts lists products with filtering
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	filters := store.ProductFilters{
		Category:       r.URL.Query().Get("category"),
		Brand:          r.URL.Query().Get("brand"),
		VeterinarianID: r.URL.Query().Get("veterinarian_id"),
		Search:         r.URL.Query().Get("search"),
		Limit:          10,
		Offset:         0,
	}

	// Parse numeric filters
	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		if price, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters.MinPrice = price
		}
	}

	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		if price, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters.MaxPrice = price
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	// Get products
	products, err := h.db.ListProducts(r.Context(), filters)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	SuccessResponse(w, products)
}

// GetProduct retrieves a specific product
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	if productID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	// Get product
	product, err := h.db.GetProductByID(r.Context(), productID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	SuccessResponse(w, product)
}

// UpdateProduct updates a product (veterinarian or admin only)
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	if productID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get existing product
	product, err := h.db.GetProductByID(r.Context(), productID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	// Check permissions - only product owner or admin can update
	if user.Role != "admin" && product.VeterinarianID != user.Sub {
		ErrorResponse(w, http.StatusForbidden, "You can only update your own products")
		return
	}

	// Parse request body
	var updateData struct {
		Name                   string                   `json:"name,omitempty"`
		Description            string                   `json:"description,omitempty"`
		Category               string                   `json:"category,omitempty"`
		Price                  *float64                 `json:"price,omitempty"`
		StockQuantity          *int                     `json:"stock_quantity,omitempty"`
		SKU                    string                   `json:"sku,omitempty"`
		Brand                  string                   `json:"brand,omitempty"`
		Weight                 *float64                 `json:"weight,omitempty"`
		Dimensions             *store.ProductDimensions `json:"dimensions,omitempty"`
		IsPrescriptionRequired *bool                    `json:"is_prescription_required,omitempty"`
		IsActive               *bool                    `json:"is_active,omitempty"`
		Images                 []string                 `json:"images,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update fields
	if updateData.Name != "" {
		product.Name = updateData.Name
	}
	if updateData.Description != "" {
		product.Description = updateData.Description
	}
	if updateData.Category != "" {
		product.Category = updateData.Category
	}
	if updateData.Price != nil {
		if *updateData.Price <= 0 {
			ErrorResponse(w, http.StatusBadRequest, "Price must be greater than 0")
			return
		}
		product.Price = *updateData.Price
	}
	if updateData.StockQuantity != nil {
		product.StockQuantity = *updateData.StockQuantity
	}
	if updateData.SKU != "" {
		product.SKU = updateData.SKU
	}
	if updateData.Brand != "" {
		product.Brand = updateData.Brand
	}
	if updateData.Weight != nil {
		product.Weight = *updateData.Weight
	}
	if updateData.Dimensions != nil {
		product.Dimensions = *updateData.Dimensions
	}
	if updateData.IsPrescriptionRequired != nil {
		product.IsPrescriptionRequired = *updateData.IsPrescriptionRequired
	}
	if updateData.IsActive != nil {
		product.IsActive = *updateData.IsActive
	}
	if updateData.Images != nil {
		product.Images = updateData.Images
	}

	// Update product
	if err := h.db.UpdateProduct(r.Context(), product); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	SuccessResponse(w, product)
}

// DeleteProduct deactivates a product (veterinarian or admin only)
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	if productID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get product
	product, err := h.db.GetProductByID(r.Context(), productID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	// Check permissions - only product owner or admin can delete
	if user.Role != "admin" && product.VeterinarianID != user.Sub {
		ErrorResponse(w, http.StatusForbidden, "You can only delete your own products")
		return
	}

	// Deactivate product instead of hard delete
	product.IsActive = false
	if err := h.db.UpdateProduct(r.Context(), product); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to deactivate product")
		return
	}

	MessageResponse(w, http.StatusOK, "Product deactivated successfully")
}

// GetVeterinarianProducts retrieves products for a specific veterinarian
func (h *ProductHandler) GetVeterinarianProducts(
	w http.ResponseWriter,
	r *http.Request,
) {
	vetID := chi.URLParam(r, "vetId")
	if vetID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Veterinarian ID is required")
		return
	}

	// Get current user from context for permission check
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Check permissions - veterinarians can only see their own products, others can see all
	if user.Role == "veterinarian" && user.Sub != vetID {
		ErrorResponse(w, http.StatusForbidden, "You can only view your own products")
		return
	}

	// Get products
	products, err := h.db.GetProductsByVeterinarianID(r.Context(), vetID)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	SuccessResponse(w, products)
}

// UpdateProductStock updates product stock quantity
func (h *ProductHandler) UpdateProductStock(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	if productID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get product for permission check
	product, err := h.db.GetProductByID(r.Context(), productID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	// Check permissions - only product owner or admin can update stock
	if user.Role != "admin" && product.VeterinarianID != user.Sub {
		ErrorResponse(
			w,
			http.StatusForbidden,
			"You can only update stock for your own products",
		)
		return
	}

	// Parse request body
	var req struct {
		Quantity int `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Quantity < 0 {
		ErrorResponse(w, http.StatusBadRequest, "Stock quantity cannot be negative")
		return
	}

	// Update stock
	if err := h.db.UpdateProductStock(r.Context(), productID, req.Quantity); err != nil {
		ErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to update product stock",
		)
		return
	}

	MessageResponse(w, http.StatusOK, "Stock updated successfully")
}
