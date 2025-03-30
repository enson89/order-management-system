package handlers

import (
	"encoding/json"
	"github.com/enson89/order-management-system/internal/db"
	"github.com/enson89/order-management-system/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	Service *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{Service: svc}
}

// HealthCheck godoc
// @Summary Health Check
// @Description Returns the health status of the application
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *OrderHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Set the header to JSON
	w.Header().Set("Content-Type", "application/json")
	// Write HTTP 200 OK status
	w.WriteHeader(http.StatusOK)
	// Return a JSON response with a health status message
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// @Summary Create Order
// @Description Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Param order body db.Order true "Order Request"
// @Success 201 {object} db.Order
// @Failure 400 {string} string "Invalid request"
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order db.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := h.Service.CreateOrder(&order)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// @Summary Get Order
// @Description Get an order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} db.Order
// @Failure 404 {string} string "Order not found"
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	order, err := h.Service.GetOrderById(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(order)
}
