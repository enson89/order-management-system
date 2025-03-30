package service

import (
	"github.com/enson89/order-management-system/internal/cache"
	"github.com/enson89/order-management-system/internal/db"
	"github.com/enson89/order-management-system/internal/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrderSuccess(t *testing.T) {
	mockDB, mockSQL, _ := sqlmock.New()
	defer mockDB.Close()

	client, mockRedis := redismock.NewClientMock()
	cache := &cache.RedisCache{Client: client}

	orderRepo := repository.NewOrderRepository(mockDB)
	orderService := NewOrderService(orderRepo, cache, "localhost:9092", "order_events")

	order := &db.Order{
		CustomerName: "Alice",
		ProductName:  "Phone",
		Quantity:     3,
		Status:       "Pending",
	}

	mockSQL.ExpectExec("INSERT INTO orders").
		WithArgs(order.CustomerName, order.ProductName, order.Quantity, order.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockRedis.ExpectSet("order:0", mock.Anything, 5*time.Minute).SetVal("OK")

	err := orderService.CreateOrder(order)
	assert.NoError(t, err)
	assert.NoError(t, mockSQL.ExpectationsWereMet())
	assert.NoError(t, mockRedis.ExpectationsWereMet())
}

func TestGetOrderById_CacheMiss_DBHit(t *testing.T) {
	mockDB, mockSQL, _ := sqlmock.New()
	defer mockDB.Close()

	client, mockRedis := redismock.NewClientMock()
	cache := &cache.RedisCache{Client: client}

	orderRepo := repository.NewOrderRepository(mockDB)
	orderService := NewOrderService(orderRepo, cache, "localhost:9092", "order_events")

	order := &db.Order{
		ID:           1,
		CustomerName: "Alice",
		ProductName:  "Phone",
		Quantity:     2,
		Status:       "Pending",
	}

	// Expect Redis GET to return a cache miss
	mockRedis.ExpectGet("order:1").RedisNil()

	// Expect DB query after cache miss
	rows := sqlmock.NewRows([]string{"id", "customer_name", "product_name", "quantity", "status", "created_at"}).
		AddRow(order.ID, order.CustomerName, order.ProductName, order.Quantity, order.Status, order.CreatedAt)

	mockSQL.ExpectQuery("SELECT id, customer_name, product_name, quantity, status, created_at FROM orders WHERE id = ?").
		WithArgs(order.ID).
		WillReturnRows(rows)

	// Expect Redis SET after DB hit
	mockRedis.ExpectSet("order:1", mock.Anything, 5*time.Minute).SetVal("OK")

	result, err := orderService.GetOrderById(1)
	assert.NoError(t, err)
	assert.Equal(t, order.CustomerName, result.CustomerName)
	assert.NoError(t, mockSQL.ExpectationsWereMet())
	assert.NoError(t, mockRedis.ExpectationsWereMet())
}
