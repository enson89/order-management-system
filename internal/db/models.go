package db

import "time"

type Order struct {
	ID           int       `json:"id"`
	CustomerName string    `json:"customer_name"`
	ProductName  string    `json:"product_name"`
	Quantity     int       `json:"quantity"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
