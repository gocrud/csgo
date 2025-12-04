package repositories

import (
	"fmt"
	"sync"
	"time"

	"vertical_slice_demo/shared/contracts/repositories"
	"vertical_slice_demo/shared/domain"
	"vertical_slice_demo/shared/infrastructure/database"
)

// ProductRepository 商品仓储实现（内存版）
type ProductRepository struct {
	db       *database.DB
	products map[int64]*domain.Product
	mu       sync.RWMutex
	nextID   int64
}

// NewProductRepository 创建商品仓储
func NewProductRepository(db *database.DB) repositories.IProductRepository {
	return &ProductRepository{
		db:       db,
		products: make(map[int64]*domain.Product),
		nextID:   1,
	}
}

// Create 创建商品
func (r *ProductRepository) Create(product *domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	product.ID = r.nextID
	r.nextID++
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	r.products[product.ID] = product
	return nil
}

// GetByID 根据 ID 获取商品
func (r *ProductRepository) GetByID(id int64) (*domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, ok := r.products[id]
	if !ok {
		return nil, fmt.Errorf("product not found")
	}
	return product, nil
}

// List 获取商品列表
func (r *ProductRepository) List(offset, limit int, status string) ([]*domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*domain.Product, 0)
	for _, product := range r.products {
		if status == "" || product.Status == status {
			products = append(products, product)
		}
	}

	start := offset
	if start > len(products) {
		start = len(products)
	}
	end := start + limit
	if end > len(products) {
		end = len(products)
	}

	return products[start:end], nil
}

// Update 更新商品
func (r *ProductRepository) Update(product *domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[product.ID]; !ok {
		return fmt.Errorf("product not found")
	}

	product.UpdatedAt = time.Now()
	r.products[product.ID] = product
	return nil
}

// Delete 删除商品
func (r *ProductRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[id]; !ok {
		return fmt.Errorf("product not found")
	}

	delete(r.products, id)
	return nil
}

// UpdateStock 更新库存
func (r *ProductRepository) UpdateStock(id int64, quantity int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	product, ok := r.products[id]
	if !ok {
		return fmt.Errorf("product not found")
	}

	product.Stock += quantity
	product.UpdatedAt = time.Now()
	return nil
}

