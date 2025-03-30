package cache

import (
	"github.com/enson89/order-management-system/internal/db"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisCache(t *testing.T) {
	cache := NewRedisCache("localhost", 6379, "")

	order := &db.Order{
		ID:           1,
		CustomerName: "Alice",
		ProductName:  "Phone",
		Quantity:     2,
		Status:       "Pending",
	}

	key := "order:1"

	// Set order in cache
	err := cache.Set(key, order, 5*time.Minute)
	assert.NoError(t, err)

	// Get order from cache
	var retrievedOrder db.Order
	err = cache.Get(key, &retrievedOrder)
	assert.NoError(t, err)
	assert.Equal(t, order.CustomerName, retrievedOrder.CustomerName)

	// Delete order from cache
	err = cache.Delete(key)
	assert.NoError(t, err)

	// Try to get after delete
	err = cache.Get(key, &retrievedOrder)
	assert.Error(t, err)
}
