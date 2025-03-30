package main

import (
	"github.com/enson89/order-management-system/internal/cache"
	"github.com/enson89/order-management-system/internal/config"
	"github.com/enson89/order-management-system/internal/db"
	"github.com/enson89/order-management-system/internal/handlers"
	"github.com/enson89/order-management-system/internal/repository"
	"github.com/enson89/order-management-system/internal/routes"
	"github.com/enson89/order-management-system/internal/service"
	"log"
	"net/http"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize PostgreSQL connection
	dbConn, err := db.InitDB(
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
	)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer dbConn.Close()

	// Initialize Redis
	redisClient := cache.NewRedisCache(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
	)

	// Initialize Repository and Services
	orderRepo := repository.NewOrderRepository(dbConn)
	orderService := service.NewOrderService(orderRepo, redisClient, cfg.Kafka.Brokers[0], cfg.Kafka.Topic)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Initialize routes
	router := routes.InitRoutes(orderHandler)

	// Start HTTP server
	log.Println("Server started at port:", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, router))
}
