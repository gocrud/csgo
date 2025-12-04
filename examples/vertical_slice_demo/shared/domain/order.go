package domain

import "time"

// Order 订单实体
type Order struct {
	ID         int64       `json:"id"`
	UserID     int64       `json:"user_id"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"` // pending, paid, shipped, completed, cancelled
	Items      []OrderItem `json:"items"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// OrderItem 订单项
type OrderItem struct {
	ID        int64   `json:"id"`
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

