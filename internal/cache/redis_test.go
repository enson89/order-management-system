package cache

import (
	"encoding/json"
	"github.com/enson89/order-management-system/internal/db"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisCache_SetAndGet(t *testing.T) {
	mockDB := &db.Order{
		ID:           1,
		CustomerName: "Alice",
		ProductName:  "Phone",
		Quantity:     2,
		Status:       "Pending",
	}

	// Create mock Redis Client
	client, mock := redismock.NewClientMock()

	cache := &RedisCache{Client: client}
	key := "order:1"
	ttl := 5 * time.Minute

	// Marshal order to JSON for storing in Redis
	expectedJSON, err := json.Marshal(mockDB)
	assert.NoError(t, err)

	// Set expectations for Redis SET
	mock.ExpectSet(key, string(expectedJSON), ttl).SetVal("OK")

	// Test Redis SET
	err = cache.Set(key, mockDB, ttl)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Set expectations for Redis GET
	mock.ExpectGet(key).SetVal(string(expectedJSON))

	var retrievedOrder db.Order
	err = cache.Get(key, &retrievedOrder)
	assert.NoError(t, err)
	assert.Equal(t, mockDB.CustomerName, retrievedOrder.CustomerName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisCache_Delete(t *testing.T) {
	client, mock := redismock.NewClientMock()
	cache := &RedisCache{Client: client}

	key := "order:1"

	// Set expectations for Redis DEL
	mock.ExpectDel(key).SetVal(1)

	err := cache.Delete(key)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
