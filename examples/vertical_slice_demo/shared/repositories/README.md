# 仓储层（Repositories）

## 📦 概述

仓储层负责数据访问，封装了所有与数据存储相关的操作。采用 Repository 模式，提供统一的数据访问接口。

## 🎯 Repository 模式的优势

- ✅ **关注点分离**：业务逻辑不需要关心数据如何存储
- ✅ **易于测试**：可以轻松创建 Mock 实现
- ✅ **易于替换**：从内存切换到数据库不影响业务代码
- ✅ **统一接口**：所有数据访问都通过标准接口

## 📂 当前的仓储

### 1. UserRepository - 用户仓储

```go
type IUserRepository interface {
    Create(user *domain.User) error
    GetByID(id int64) (*domain.User, error)
    GetByEmail(email string) (*domain.User, error)
    List(offset, limit int) ([]*domain.User, error)
    Update(user *domain.User) error
    Delete(id int64) error
    ExistsByEmail(email string) bool
}
```

### 2. ProductRepository - 商品仓储

```go
type IProductRepository interface {
    Create(product *domain.Product) error
    GetByID(id int64) (*domain.Product, error)
    List(offset, limit int, status string) ([]*domain.Product, error)
    Update(product *domain.Product) error
    Delete(id int64) error
    UpdateStock(id int64, quantity int) error
}
```

### 3. OrderRepository - 订单仓储

```go
type IOrderRepository interface {
    Create(order *domain.Order) error
    GetByID(id int64) (*domain.Order, error)
    GetByUserID(userID int64, offset, limit int) ([]*domain.Order, error)
    Update(order *domain.Order) error
    UpdateStatus(id int64, status string) error
}
```

## 🏗️ 架构设计

```
┌─────────────────────────────────────┐
│        Business Logic               │
│   (Features / Handlers)             │
└─────────────┬───────────────────────┘
              │ 依赖接口
              ▼
┌─────────────────────────────────────┐
│    contracts/repositories/          │
│    (接口定义)                       │
│                                     │
│  - IUserRepository                  │
│  - IProductRepository               │
│  - IOrderRepository                 │
└─────────────┬───────────────────────┘
              │ 实现
              ▼
┌─────────────────────────────────────┐
│    repositories/                    │
│    (具体实现)                       │
│                                     │
│  - UserRepositoryImpl (内存版)      │
│  - ProductRepositoryImpl (内存版)   │
│  - OrderRepositoryImpl (内存版)     │
└─────────────┬───────────────────────┘
              │
              ▼
┌─────────────────────────────────────┐
│       Data Storage                  │
│   (内存 / 数据库 / 缓存)           │
└─────────────────────────────────────┘
```

## 💡 使用示例

### 在功能切片中使用

```go
// Handler 依赖仓储接口
type CreateUserHandler struct {
    userRepo repositories.IUserRepository
}

func NewCreateUserHandler(userRepo repositories.IUserRepository) *CreateUserHandler {
    return &CreateUserHandler{userRepo: userRepo}
}

func (h *CreateUserHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }

    // 业务验证
    if h.userRepo.ExistsByEmail(req.Email) {
        return c.BadRequest("邮箱已存在")
    }

    // 创建用户
    user := &domain.User{
        Name:  req.Name,
        Email: req.Email,
    }

    if err := h.userRepo.Create(user); err != nil {
        return c.InternalError("创建失败")
    }

    return c.Created(user)
}
```

## 🔄 从内存切换到数据库

当前实现使用内存存储（`map`），在实际项目中需要切换到真实数据库。

### 步骤 1：安装数据库驱动

```bash
go get gorm.io/gorm
go get gorm.io/driver/postgres
```

### 步骤 2：创建数据库实现

```go
// repositories/user_repository_gorm.go
package repositories

import (
    "gorm.io/gorm"
    "vertical_slice_demo/shared/contracts/repositories"
    "vertical_slice_demo/shared/domain"
)

type UserRepositoryGorm struct {
    db *gorm.DB
}

func NewUserRepositoryGorm(db *gorm.DB) repositories.IUserRepository {
    return &UserRepositoryGorm{db: db}
}

func (r *UserRepositoryGorm) Create(user *domain.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepositoryGorm) GetByID(id int64) (*domain.User, error) {
    var user domain.User
    err := r.db.First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// ... 其他方法实现
```

### 步骤 3：更新 DI 注册

```go
// repositories/service_extensions.go
func AddRepositories(services di.IServiceCollection) {
    // 原来：内存版
    // services.AddSingleton(NewUserRepository)
    
    // 现在：数据库版
    services.AddSingleton(NewUserRepositoryGorm)
    services.AddSingleton(NewProductRepositoryGorm)
    services.AddSingleton(NewOrderRepositoryGorm)
}
```

业务代码无需修改，因为依赖的是接口！

## 🧪 单元测试

### 创建 Mock 实现

```go
// mocks/user_repository_mock.go
type MockUserRepository struct {
    users map[int64]*domain.User
}

func NewMockUserRepository() *MockUserRepository {
    return &MockUserRepository{
        users: make(map[int64]*domain.User),
    }
}

func (m *MockUserRepository) Create(user *domain.User) error {
    user.ID = int64(len(m.users) + 1)
    m.users[user.ID] = user
    return nil
}

// ... 其他方法
```

### 在测试中使用

```go
func TestCreateUserHandler(t *testing.T) {
    // 使用 Mock 仓储
    mockRepo := NewMockUserRepository()
    handler := NewCreateUserHandler(mockRepo)
    
    // 测试逻辑
    // ...
}
```

## 📝 最佳实践

### 1. 接口定义在 contracts 中

```go
// ✅ 好的做法
// shared/contracts/repositories/user_repository.go
package repositories

type IUserRepository interface {
    Create(user *domain.User) error
    // ...
}
```

### 2. 实现在 repositories 中

```go
// ✅ 好的做法
// shared/repositories/user_repository_impl.go
package repositories

func NewUserRepository(db *database.DB) contracts.IUserRepository {
    return &UserRepository{db: db}
}
```

### 3. 返回接口类型

```go
// ✅ 好的做法
func NewUserRepository() repositories.IUserRepository {
    return &UserRepository{}
}

// ❌ 不好的做法
func NewUserRepository() *UserRepository {
    return &UserRepository{}
}
```

### 4. 使用领域模型

```go
// ✅ 好的做法
func (r *UserRepository) Create(user *domain.User) error

// ❌ 不好的做法
func (r *UserRepository) Create(name, email string) error
```

### 5. 清晰的方法命名

```go
// ✅ 好的做法
GetByID(id int64) (*User, error)
GetByEmail(email string) (*User, error)
ExistsByEmail(email string) bool

// ❌ 不好的做法
Get(param interface{}) (*User, error)
Check(field string, value string) bool
```

## 🔧 高级功能

### 分页查询

```go
type PageRequest struct {
    Offset int
    Limit  int
}

type PageResult struct {
    Items      []*User
    Total      int
    HasMore    bool
}

func (r *UserRepository) ListWithPage(req PageRequest) (*PageResult, error) {
    // 实现分页逻辑
}
```

### 事务支持

```go
type IUnitOfWork interface {
    BeginTransaction() error
    Commit() error
    Rollback() error
}

func (r *UserRepository) CreateWithTransaction(user *User, uow IUnitOfWork) error {
    // 在事务中执行
}
```

### 查询构建器

```go
type UserQueryBuilder struct {
    filters map[string]interface{}
}

func (b *UserQueryBuilder) WhereEmail(email string) *UserQueryBuilder {
    b.filters["email"] = email
    return b
}

func (b *UserQueryBuilder) WhereRole(role string) *UserQueryBuilder {
    b.filters["role"] = role
    return b
}

func (b *UserQueryBuilder) Execute() ([]*User, error) {
    // 执行查询
}
```

## 📚 相关文档

- [领域模型](../domain/)
- [共享服务](../services/)
- [依赖注入指南](../../../docs/guides/dependency-injection.md)

---

**Repository 模式让数据访问更清晰！** 🎉

