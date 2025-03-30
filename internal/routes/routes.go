package routes

import (
	"github.com/enson89/order-management-system/internal/handlers"

	"github.com/gorilla/mux"
)

func InitRoutes(handler *handlers.OrderHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/orders", handler.CreateOrder).Methods("POST")
	router.HandleFunc("/orders/{id}", handler.GetOrderById).Methods("GET")
	return router
}
