package repositories

import (
	"fmt"
	"sync"
	"time"

	"vertical_slice_demo/shared/contracts/repositories"
	"vertical_slice_demo/shared/domain"
	"vertical_slice_demo/shared/infrastructure/database"
)

// UserRepository 用户仓储实现（内存版）
type UserRepository struct {
	db    *database.DB
	users map[int64]*domain.User
	mu    sync.RWMutex
	nextID int64
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *database.DB) repositories.IUserRepository {
	return &UserRepository{
		db:     db,
		users:  make(map[int64]*domain.User),
		nextID: 1,
	}
}

// Create 创建用户
func (r *UserRepository) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	r.nextID++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	r.users[user.ID] = user
	return nil
}

// GetByID 根据 ID 获取用户
func (r *UserRepository) GetByID(id int64) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// List 获取用户列表
func (r *UserRepository) List(offset, limit int) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	start := offset
	if start > len(users) {
		start = len(users)
	}
	end := start + limit
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], nil
}

// Update 更新用户
func (r *UserRepository) Update(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.ID]; !ok {
		return fmt.Errorf("user not found")
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	return nil
}

// Delete 删除用户
func (r *UserRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return fmt.Errorf("user not found")
	}

	delete(r.users, id)
	return nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(email string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return true
		}
	}
	return false
}

