package repositories

import (
	"fmt"
	"sync"
	"time"

	"vertical_slice_demo/shared/contracts/repositories"
	"vertical_slice_demo/shared/domain"
	"vertical_slice_demo/shared/infrastructure/database"
)

// OrderRepository 订单仓储实现（内存版）
type OrderRepository struct {
	db     *database.DB
	orders map[int64]*domain.Order
	mu     sync.RWMutex
	nextID int64
}

// NewOrderRepository 创建订单仓储
func NewOrderRepository(db *database.DB) repositories.IOrderRepository {
	return &OrderRepository{
		db:     db,
		orders: make(map[int64]*domain.Order),
		nextID: 1,
	}
}

// Create 创建订单
func (r *OrderRepository) Create(order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order.ID = r.nextID
	r.nextID++
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	r.orders[order.ID] = order
	return nil
}

// GetByID 根据 ID 获取订单
func (r *OrderRepository) GetByID(id int64) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[id]
	if !ok {
		return nil, fmt.Errorf("order not found")
	}
	return order, nil
}

// GetByUserID 根据用户 ID 获取订单列表
func (r *OrderRepository) GetByUserID(userID int64, offset, limit int) ([]*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]*domain.Order, 0)
	for _, order := range r.orders {
		if order.UserID == userID {
			orders = append(orders, order)
		}
	}

	start := offset
	if start > len(orders) {
		start = len(orders)
	}
	end := start + limit
	if end > len(orders) {
		end = len(orders)
	}

	return orders[start:end], nil
}

// Update 更新订单
func (r *OrderRepository) Update(order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.ID]; !ok {
		return fmt.Errorf("order not found")
	}

	order.UpdatedAt = time.Now()
	r.orders[order.ID] = order
	return nil
}

// UpdateStatus 更新订单状态
func (r *OrderRepository) UpdateStatus(id int64, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[id]
	if !ok {
		return fmt.Errorf("order not found")
	}

	order.Status = status
	order.UpdatedAt = time.Now()
	return nil
}

