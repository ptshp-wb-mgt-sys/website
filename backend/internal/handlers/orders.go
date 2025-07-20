// Package handlers contains order management handlers
package handlers

import (
	"encoding/json"
	"net/http"
	"pet-mgt/backend/internal/middleware"
	"pet-mgt/backend/internal/store"

	"github.com/go-chi/chi/v5"
)

// OrderHandler handles order operations
type OrderHandler struct {
	db store.Database
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(db store.Database) *OrderHandler {
	return &OrderHandler{db: db}
}

// CreateOrder creates a new order (clients only)
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Only clients can create orders (unless admin)
	if user.Role != "client" && user.Role != "admin" {
		ErrorResponse(w, http.StatusForbidden, "Only clients can create orders")
		return
	}

	// Parse request body
	var req struct {
		VeterinarianID string `json:"veterinarian_id"`
		Items          []struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
		} `json:"items"`
		PaymentMethod   string `json:"payment_method,omitempty"`
		ShippingAddress string `json:"shipping_address,omitempty"`
		DeliveryMethod  string `json:"delivery_method,omitempty"`
		Notes           string `json:"notes,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.VeterinarianID == "" || len(req.Items) == 0 {
		ErrorResponse(
			w,
			http.StatusBadRequest,
			"Missing required fields (veterinarian_id, items)",
		)
		return
	}

	// Verify veterinarian exists
	_, err := h.db.GetVeterinarianByID(r.Context(), req.VeterinarianID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Veterinarian not found")
		return
	}

	// Calculate total amount and validate products
	var totalAmount float64
	var orderItems []store.OrderItem

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			ErrorResponse(
				w,
				http.StatusBadRequest,
				"Item quantity must be greater than 0",
			)
			return
		}

		// Get product details
		product, err := h.db.GetProductByID(r.Context(), item.ProductID)
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "Product not found: "+item.ProductID)
			return
		}

		// Check if product belongs to the specified veterinarian
		if product.VeterinarianID != req.VeterinarianID {
			ErrorResponse(
				w,
				http.StatusBadRequest,
				"All products must be from the same veterinarian",
			)
			return
		}

		// Check stock availability
		if product.StockQuantity < item.Quantity {
			ErrorResponse(
				w,
				http.StatusBadRequest,
				"Insufficient stock for product: "+product.Name,
			)
			return
		}

		// Create order item
		orderItem := store.NewOrderItem(
			"",
			item.ProductID,
			item.Quantity,
			product.Price,
		)
		orderItems = append(orderItems, *orderItem)
		totalAmount += orderItem.TotalPrice
	}

	// Create order
	clientID := user.Sub
	if user.Role == "admin" && r.URL.Query().Get("client_id") != "" {
		clientID = r.URL.Query().Get("client_id")
	}

	order := store.NewOrder(clientID, req.VeterinarianID, totalAmount)
	if req.PaymentMethod != "" {
		order.PaymentMethod = req.PaymentMethod
	}
	if req.ShippingAddress != "" {
		order.ShippingAddress = req.ShippingAddress
	}
	if req.DeliveryMethod != "" {
		order.DeliveryMethod = req.DeliveryMethod
	}
	if req.Notes != "" {
		order.Notes = req.Notes
	}

	// Create order in database
	if err := h.db.CreateOrder(r.Context(), order); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	// Create order items
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
		if err := h.db.CreateOrderItem(r.Context(), &orderItems[i]); err != nil {
			ErrorResponse(
				w,
				http.StatusInternalServerError,
				"Failed to create order items",
			)
			return
		}

		// Update product stock
		product, _ := h.db.GetProductByID(r.Context(), orderItems[i].ProductID)
		newStock := product.StockQuantity - orderItems[i].Quantity
		h.db.UpdateProductStock(r.Context(), orderItems[i].ProductID, newStock)
	}

	// Return order with items
	response := map[string]any{
		"order": order,
		"items": orderItems,
	}

	SuccessResponse(w, response)
}

// GetOrders retrieves orders for the current user
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var orders []store.Order
	var err error

	switch user.Role {
	case "client":
		orders, err = h.db.GetOrdersByClientID(r.Context(), user.Sub)
	case "veterinarian":
		orders, err = h.db.GetOrdersByVeterinarianID(r.Context(), user.Sub)
	case "admin":
		// For admin, check query params
		clientID := r.URL.Query().Get("client_id")
		vetID := r.URL.Query().Get("veterinarian_id")

		if clientID != "" {
			orders, err = h.db.GetOrdersByClientID(r.Context(), clientID)
		} else if vetID != "" {
			orders, err = h.db.GetOrdersByVeterinarianID(r.Context(), vetID)
		} else {
			// Return empty for now - would need a new method for all orders
			orders = []store.Order{}
		}
	default:
		ErrorResponse(w, http.StatusForbidden, "Invalid user role")
		return
	}

	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve orders")
		return
	}

	SuccessResponse(w, orders)
}

// GetOrder retrieves a specific order
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get order
	order, err := h.db.GetOrderByID(r.Context(), orderID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Order not found")
		return
	}

	// Check permissions
	if user.Role != "admin" {
		if user.Role == "client" && order.ClientID != user.Sub {
			ErrorResponse(w, http.StatusForbidden, "You can only view your own orders")
			return
		}
		if user.Role == "veterinarian" && order.VeterinarianID != user.Sub {
			ErrorResponse(w, http.StatusForbidden, "You can only view your own orders")
			return
		}
	}

	// Get order items
	items, err := h.db.GetOrderItems(r.Context(), orderID)
	if err != nil {
		ErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve order items",
		)
		return
	}

	// Return order with items
	response := map[string]any{
		"order": order,
		"items": items,
	}

	SuccessResponse(w, response)
}

// UpdateOrderStatus updates the status of an order
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get order
	order, err := h.db.GetOrderByID(r.Context(), orderID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Order not found")
		return
	}

	// Check permissions - only veterinarian who owns the order or admin can update status
	if user.Role != "admin" &&
		(user.Role != "veterinarian" || order.VeterinarianID != user.Sub) {
		ErrorResponse(
			w,
			http.StatusForbidden,
			"You can only update status for your own orders",
		)
		return
	}

	// Parse request body
	var req struct {
		Status        string `json:"status,omitempty"`
		PaymentStatus string `json:"payment_status,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate status values
	if req.Status != "" {
		validStatuses := map[string]bool{
			"pending":    true,
			"confirmed":  true,
			"processing": true,
			"shipped":    true,
			"delivered":  true,
			"cancelled":  true,
		}
		if !validStatuses[req.Status] {
			ErrorResponse(w, http.StatusBadRequest, "Invalid order status")
			return
		}
	}

	if req.PaymentStatus != "" {
		validPaymentStatuses := map[string]bool{
			"pending":  true,
			"paid":     true,
			"failed":   true,
			"refunded": true,
		}
		if !validPaymentStatuses[req.PaymentStatus] {
			ErrorResponse(w, http.StatusBadRequest, "Invalid payment status")
			return
		}
	}

	// Update status
	statusToUpdate := req.Status
	if statusToUpdate == "" {
		statusToUpdate = order.Status
	}

	if err := h.db.UpdateOrderStatus(r.Context(), orderID, statusToUpdate); err != nil {
		ErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to update order status",
		)
		return
	}

	// If payment status is provided, we'd need a separate method to update it
	// For now, just return success for the order status update

	MessageResponse(w, http.StatusOK, "Order status updated successfully")
}

// CancelOrder cancels an order (clients and veterinarians can cancel their own orders)
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		ErrorResponse(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	// Get current user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		ErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get order
	order, err := h.db.GetOrderByID(r.Context(), orderID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Order not found")
		return
	}

	// Check permissions
	canCancel := false
	if user.Role == "admin" {
		canCancel = true
	} else if user.Role == "client" && order.ClientID == user.Sub {
		canCancel = true
	} else if user.Role == "veterinarian" && order.VeterinarianID == user.Sub {
		canCancel = true
	}

	if !canCancel {
		ErrorResponse(w, http.StatusForbidden, "You can only cancel your own orders")
		return
	}

	// Check if order can be cancelled
	if order.Status == "shipped" || order.Status == "delivered" ||
		order.Status == "cancelled" {
		ErrorResponse(
			w,
			http.StatusBadRequest,
			"Order cannot be cancelled in current status",
		)
		return
	}

	// Cancel order
	if err := h.db.UpdateOrderStatus(r.Context(), orderID, "cancelled"); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to cancel order")
		return
	}

	// Restore product stock
	items, err := h.db.GetOrderItems(r.Context(), orderID)
	if err == nil {
		for _, item := range items {
			product, err := h.db.GetProductByID(r.Context(), item.ProductID)
			if err == nil {
				newStock := product.StockQuantity + item.Quantity
				h.db.UpdateProductStock(r.Context(), item.ProductID, newStock)
			}
		}
	}

	MessageResponse(w, http.StatusOK, "Order cancelled successfully")
}
