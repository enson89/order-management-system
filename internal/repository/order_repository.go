package repository

import (
	"database/sql"
	"github.com/enson89/order-management-system/internal/db"
)

// OrderRepositoryInterface defines the methods for the order repository
type OrderRepositoryInterface interface {
	CreateOrder(order *db.Order) error
	GetOrderById(id int) (*db.Order, error)
}

// OrderRepository implements OrderRepositoryInterface
type OrderRepository struct {
	DB *sql.DB
}

// NewOrderRepository creates a new OrderRepository
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

// CreateOrder inserts a new order in the database
func (r *OrderRepository) CreateOrder(order *db.Order) error {
	_, err := r.DB.Exec(
		"INSERT INTO orders (customer_name, product_name, quantity, status) VALUES ($1, $2, $3, $4)",
		order.CustomerName, order.ProductName, order.Quantity, order.Status,
	)
	return err
}

// GetOrderById fetches an order by ID from the database
func (r *OrderRepository) GetOrderById(id int) (*db.Order, error) {
	order := &db.Order{}
	row := r.DB.QueryRow("SELECT id, customer_name, product_name, quantity, status, created_at FROM orders WHERE id = $1", id)
	err := row.Scan(&order.ID, &order.CustomerName, &order.ProductName, &order.Quantity, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil, err
	}
	return order, nil
}
