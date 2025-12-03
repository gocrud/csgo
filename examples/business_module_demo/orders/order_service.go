package orders

import (
	"fmt"
	"time"
)

// Order represents an order entity.
type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// OrderService defines the order service interface.
type OrderService interface {
	GetOrder(id int) (*Order, error)
	ListOrdersByUser(userID int) ([]*Order, error)
	CreateOrder(order *Order) error
}

// orderService is the default implementation.
type orderService struct {
	orders map[int]*Order
}

// NewOrderService creates a new OrderService.
func NewOrderService() OrderService {
	return &orderService{
		orders: map[int]*Order{
			1: {
				ID:         1,
				UserID:     1,
				TotalPrice: 99.99,
				Status:     "completed",
				CreatedAt:  time.Now().Add(-24 * time.Hour),
			},
		},
	}
}

func (s *orderService) GetOrder(id int) (*Order, error) {
	order, ok := s.orders[id]
	if !ok {
		return nil, fmt.Errorf("order not found: %d", id)
	}
	return order, nil
}

func (s *orderService) ListOrdersByUser(userID int) ([]*Order, error) {
	orders := make([]*Order, 0)
	for _, o := range s.orders {
		if o.UserID == userID {
			orders = append(orders, o)
		}
	}
	return orders, nil
}

func (s *orderService) CreateOrder(order *Order) error {
	if _, exists := s.orders[order.ID]; exists {
		return fmt.Errorf("order already exists: %d", order.ID)
	}
	order.CreatedAt = time.Now()
	s.orders[order.ID] = order
	return nil
}

