package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(host string, port int, password string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       0,
	})
	return &RedisCache{Client: rdb}
}

// Set sets a key-value pair with a TTL in Redis
func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Client.Set(ctx, key, string(data), ttl).Err()
}

// Get retrieves a value by key from Redis
func (c *RedisCache) Get(key string, dest interface{}) error {
	val, err := c.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("key not found")
	} else if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// Delete removes a key from Redis
func (c *RedisCache) Delete(key string) error {
	return c.Client.Del(ctx, key).Err()
}
