package repositories

import "vertical_slice_demo/shared/domain"

// IProductRepository 商品仓储接口
type IProductRepository interface {
	// Create 创建商品
	Create(product *domain.Product) error
	// GetByID 根据 ID 获取商品
	GetByID(id int64) (*domain.Product, error)
	// List 获取商品列表
	List(offset, limit int, status string) ([]*domain.Product, error)
	// Update 更新商品
	Update(product *domain.Product) error
	// Delete 删除商品
	Delete(id int64) error
	// UpdateStock 更新库存
	UpdateStock(id int64, quantity int) error
}

