package repository

import (
	"database/sql"
	"github.com/enson89/order-management-system/internal/db"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) CreateOrder(order *db.Order) error {
	_, err := r.DB.Exec(
		"INSERT INTO orders (customer_name, product_name, quantity, status) VALUES ($1, $2, $3, $4)",
		order.CustomerName, order.ProductName, order.Quantity, order.Status,
	)
	return err
}

func (r *OrderRepository) GetOrderById(id int) (*db.Order, error) {
	order := &db.Order{}
	row := r.DB.QueryRow("SELECT id, customer_name, product_name, quantity, status, created_at FROM orders WHERE id = $1", id)
	err := row.Scan(&order.ID, &order.CustomerName, &order.ProductName, &order.Quantity, &order.Status, &order.CreatedAt)
	return order, err
}
