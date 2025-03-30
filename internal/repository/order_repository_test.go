package repository

import (
	"github.com/enson89/order-management-system/internal/db"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	repo := NewOrderRepository(mockDB)

	order := &db.Order{
		CustomerName: "John Doe",
		ProductName:  "Laptop",
		Quantity:     2,
		Status:       "Pending",
	}

	mock.ExpectExec("INSERT INTO orders").
		WithArgs(order.CustomerName, order.ProductName, order.Quantity, order.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateOrder(order)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderById(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	repo := NewOrderRepository(mockDB)

	order := &db.Order{
		ID:           1,
		CustomerName: "Alice",
		ProductName:  "Phone",
		Quantity:     3,
		Status:       "Completed",
	}

	rows := sqlmock.NewRows([]string{"id", "customer_name", "product_name", "quantity", "status", "created_at"}).
		AddRow(order.ID, order.CustomerName, order.ProductName, order.Quantity, order.Status, order.CreatedAt)

	mock.ExpectQuery("SELECT id, customer_name, product_name, quantity, status, created_at FROM orders WHERE id = ?").
		WithArgs(order.ID).
		WillReturnRows(rows)

	result, err := repo.GetOrderById(order.ID)
	assert.NoError(t, err)
	assert.Equal(t, order.CustomerName, result.CustomerName)
	assert.NoError(t, mock.ExpectationsWereMet())
}
