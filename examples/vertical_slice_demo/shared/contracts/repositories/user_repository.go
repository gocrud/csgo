package repositories

import "vertical_slice_demo/shared/domain"

// IUserRepository 用户仓储接口
type IUserRepository interface {
	// Create 创建用户
	Create(user *domain.User) error
	// GetByID 根据 ID 获取用户
	GetByID(id int64) (*domain.User, error)
	// GetByEmail 根据邮箱获取用户
	GetByEmail(email string) (*domain.User, error)
	// List 获取用户列表
	List(offset, limit int) ([]*domain.User, error)
	// Update 更新用户
	Update(user *domain.User) error
	// Delete 删除用户
	Delete(id int64) error
	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(email string) bool
}

