package service

import (
	"encoding/json"
	"fmt"
	"github.com/enson89/order-management-system/internal/cache"
	"github.com/enson89/order-management-system/internal/db"
	"github.com/enson89/order-management-system/internal/kafka"
	"github.com/enson89/order-management-system/internal/repository"
	"log"
	"time"
)

const cacheTTL = 5 * time.Minute

type OrderService struct {
	Repo  *repository.OrderRepository
	Cache *cache.RedisCache
	Kafka string
	Topic string
}

func NewOrderService(repo *repository.OrderRepository, cache *cache.RedisCache, kafkaBroker, topic string) *OrderService {
	return &OrderService{
		Repo:  repo,
		Cache: cache,
		Kafka: kafkaBroker,
		Topic: topic,
	}
}

// CreateOrder creates an order and publishes a Kafka message
func (s *OrderService) CreateOrder(order *db.Order) error {
	// Save to DB
	err := s.Repo.CreateOrder(order)
	if err != nil {
		return err
	}

	// Cache the order
	cacheKey := fmt.Sprintf("order:%d", order.ID)
	s.Cache.Set(cacheKey, order, cacheTTL)

	// Publish to Kafka
	orderJSON, _ := json.Marshal(order)
	err = kafka.ProduceMessage(s.Kafka, s.Topic, string(orderJSON))
	if err != nil {
		log.Printf("Failed to send Kafka message: %v\n", err)
	}
	return nil
}

// GetOrderById retrieves an order by ID with Redis caching
func (s *OrderService) GetOrderById(id int) (*db.Order, error) {
	cacheKey := fmt.Sprintf("order:%d", id)

	// Try to get from Redis cache
	var order db.Order
	err := s.Cache.Get(cacheKey, &order)
	if err == nil {
		log.Printf("Cache hit for order %d\n", id)
		return &order, nil
	}

	// If cache miss, fetch from DB
	log.Printf("Cache miss for order %d. Fetching from DB...\n", id)
	orderPtr, err := s.Repo.GetOrderById(id)
	if err != nil {
		return nil, err
	}

	// Cache the result for future use
	s.Cache.Set(cacheKey, orderPtr, cacheTTL)
	return orderPtr, nil
}
