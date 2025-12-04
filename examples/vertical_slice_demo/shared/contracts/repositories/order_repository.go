package repositories

import "vertical_slice_demo/shared/domain"

// IOrderRepository 订单仓储接口
type IOrderRepository interface {
	// Create 创建订单
	Create(order *domain.Order) error
	// GetByID 根据 ID 获取订单
	GetByID(id int64) (*domain.Order, error)
	// GetByUserID 根据用户 ID 获取订单列表
	GetByUserID(userID int64, offset, limit int) ([]*domain.Order, error)
	// Update 更新订单
	Update(order *domain.Order) error
	// UpdateStatus 更新订单状态
	UpdateStatus(id int64, status string) error
}

