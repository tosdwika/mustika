package models

import (
	"time"
)

type Order struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CustomerID  uint      `json:"customer_id"`
	OrderDate   time.Time `json:"order_date" gorm:"type:datetime;not null"`
	ProductName string    `json:"product_name"`
	Status      string    `json:"status"`
	Total       float64   `json:"total"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
