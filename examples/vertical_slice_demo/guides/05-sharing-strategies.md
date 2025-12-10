# 共享策略完整指南

## 目录

- [共享层次概览](#共享层次概览)
- [功能级共享](#功能级共享)
- [应用级共享](#应用级共享)
- [全局共享](#全局共享)
- [决策树：代码应该共享到哪里](#决策树代码应该共享到哪里)
- [对比表：不同共享层次](#对比表不同共享层次)
- [最佳实践](#最佳实践)
- [实战案例](#实战案例)

---

## 共享层次概览

```
┌─────────────────────────────────────────┐
│  功能私有（不共享）                        │
│  features/xxx/                          │
│  - 只在当前功能使用                        │
└─────────────────────────────────────────┘
         ↓ 需要在功能内部共享
┌─────────────────────────────────────────┐
│  功能级共享                               │
│  features/xxx/internal/                 │
│  - 功能内部复用，外部不可见                 │
└─────────────────────────────────────────┘
         ↓ 需要在应用内共享
┌─────────────────────────────────────────┐
│  应用级共享                               │
│  apps/*/internal/                       │
│  - 应用内多个功能共享，但不跨应用             │
└─────────────────────────────────────────┘
         ↓ 需要跨应用共享
┌─────────────────────────────────────────┐
│  全局共享                                 │
│  shared/                                │
│  - 多个应用端共享                          │
└─────────────────────────────────────────┘
```

**核心原则：** 代码应该在尽可能小的范围内共享，逐级提升。

---

## 功能级共享

### 适用场景

**问题：** 当前功能内的多个 Handler 需要共享 Store 或 Business 逻辑，但不希望暴露给其他功能

**解决方案：** 使用 `features/xxx/internal/`

### 示例场景

报表功能有多个操作（生成、导出、调度），需要共享数据聚合器和报表生成器：

```
apps/admin/features/reports/
├── internal/
│   ├── entity/
│   │   └── report_entity.go         # 内部实体
│   ├── data/
│   │   └── report_store.go          # 共享 Store
│   └── business/
│       ├── data_aggregator.go       # 共享业务逻辑
│       └── report_generator.go
│
├── models.go                        # 对外 DTO
├── generate_report.go               # 使用 internal/business
├── export_report.go                 # 使用 internal/business
├── schedule_report.go               # 使用 internal/business
├── controller.go
└── service_extensions.go
```

### 代码示例

**internal/entity/report_entity.go:**

```go
package entity

// ReportEntity 内部报表实体（只在 reports 功能内使用）
type ReportEntity struct {
    ID         int64
    Name       string
    Type       ReportType
    DataSource DataSourceConfig
    Charts     []ChartConfig
    CreatedAt  time.Time
}

type ReportType string

const (
    ReportTypeSales     ReportType = "sales"
    ReportTypeUser      ReportType = "user"
    ReportTypeInventory ReportType = "inventory"
)
```

**internal/data/report_store.go:**

```go
package data

import "vertical_slice_demo/apps/admin/features/reports/internal/entity"

// ReportStore 内部数据存储（功能内多个 Handler 共享）
type ReportStore struct {
    reports map[int64]*entity.ReportEntity
    mu      sync.RWMutex
}

func NewReportStore() *ReportStore {
    return &ReportStore{
        reports: make(map[int64]*entity.ReportEntity),
    }
}

func (s *ReportStore) Create(report *entity.ReportEntity) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.reports[report.ID] = report
    return nil
}
```

**internal/business/data_aggregator.go:**

```go
package business

// DataAggregator 数据聚合器（功能内多个 Handler 共享）
type DataAggregator struct {
    orderRepo repositories.IOrderRepository
}

func NewDataAggregator(orderRepo repositories.IOrderRepository) *DataAggregator {
    return &DataAggregator{orderRepo: orderRepo}
}

func (a *DataAggregator) AggregateByDay(startDate, endDate time.Time) map[string]interface{} {
    // 复杂的聚合逻辑
    orders, _ := a.orderRepo.GetByPeriod(startDate, endDate)
    return a.groupByDay(orders)
}
```

**generate_report.go (Handler):**

```go
package reports

import "vertical_slice_demo/apps/admin/features/reports/internal/business"

type GenerateReportHandler struct {
    aggregator *business.DataAggregator  // ✅ 使用功能内部共享的业务逻辑
    generator  *business.ReportGenerator
}

func (h *GenerateReportHandler) Handle(c *web.HttpContext) web.IActionResult {
    // 使用共享的聚合器
    data := h.aggregator.AggregateByDay(req.StartDate, req.EndDate)
    // 使用共享的生成器
    report, _ := h.generator.Generate(data)
    return c.Ok(toReportResponse(report))
}
```

**export_report.go (Handler):**

```go
package reports

import "vertical_slice_demo/apps/admin/features/reports/internal/business"

type ExportReportHandler struct {
    aggregator *business.DataAggregator  // ✅ 复用相同的聚合器
    exporter   *business.ReportExporter
}

func (h *ExportReportHandler) Handle(c *web.HttpContext) web.IActionResult {
    // 复用聚合器
    data := h.aggregator.AggregateByDay(req.StartDate, req.EndDate)
    // 导出
    file, _ := h.exporter.ExportToExcel(data)
    return c.File(file)
}
```

### 好处

- ✅ 代码复用：多个 Handler 共享 Store 和 Business 逻辑
- ✅ 封装内部：使用 internal 防止其他功能误用
- ✅ 职责清晰：对外接口（Handler）和内部实现分离
- ✅ 易于测试：可以单独测试 internal 中的逻辑

---

## 应用级共享

### 适用场景

**问题：** Admin 端多个功能需要共享某些代码，但 API 端不需要

**解决方案：** 使用 `apps/admin/internal/`

### 何时使用

- ✅ 应用特定的权限检查
- ✅ 应用特定的缓存管理
- ✅ 应用特定的会话处理
- ✅ 应用特定的中间件
- ❌ 不是核心业务实体（那应该在 shared/domain）

### 示例场景

Admin 端的多个功能都需要检查管理员权限：

```
apps/admin/
├── internal/
│   ├── auth/
│   │   └── permission_checker.go    # Admin 端权限检查
│   ├── cache/
│   │   └── admin_cache.go           # Admin 端缓存
│   ├── models/
│   │   └── admin_session.go         # Admin 端会话
│   └── middleware/
│       └── auth_middleware.go       # Admin 端认证中间件
│
└── features/
    ├── users/                       # 使用 admin/internal/auth
    ├── products/                    # 使用 admin/internal/auth
    └── orders/                      # 使用 admin/internal/auth
```

### 代码示例

**apps/admin/internal/auth/permission_checker.go:**

```go
package auth

// PermissionChecker Admin 端权限检查（应用内共享）
type PermissionChecker struct {
    adminRepo IAdminRepository
    cache     *cache.AdminCache
}

func NewPermissionChecker(adminRepo IAdminRepository, cache *cache.AdminCache) *PermissionChecker {
    return &PermissionChecker{
        adminRepo: adminRepo,
        cache:     cache,
    }
}

func (p *PermissionChecker) CheckPermission(adminID int64, resource string, action string) bool {
    // 1. 从缓存获取权限
    if perm, ok := p.cache.GetPermission(adminID, resource); ok {
        return perm.HasAction(action)
    }
    
    // 2. 从数据库获取
    admin, _ := p.adminRepo.GetByID(adminID)
    return admin.HasPermission(resource, action)
}

func (p *PermissionChecker) RequirePermission(adminID int64, resource string, action string) error {
    if !p.CheckPermission(adminID, resource, action) {
        return errors.New("permission denied")
    }
    return nil
}
```

**apps/admin/internal/middleware/auth_middleware.go:**

```go
package middleware

import "vertical_slice_demo/apps/admin/internal/auth"

// AuthMiddleware Admin 端认证中间件
type AuthMiddleware struct {
    checker *auth.PermissionChecker
}

func NewAuthMiddleware(checker *auth.PermissionChecker) *AuthMiddleware {
    return &AuthMiddleware{checker: checker}
}

func (m *AuthMiddleware) RequirePermission(resource string, action string) web.HandlerFunc {
    return func(c *web.HttpContext) {
        adminID := c.GetUserID()
        
        if err := m.checker.RequirePermission(adminID, resource, action); err != nil {
            c.AbortWithStatus(403)
            return
        }
        
        c.Next()
    }
}
```

**apps/admin/features/users/create_user.go:**

```go
package users

import "vertical_slice_demo/apps/admin/internal/auth"

type CreateUserHandler struct {
    userRepo IUserRepository
    checker  *auth.PermissionChecker  // ✅ 使用应用级共享的权限检查
}

func (h *CreateUserHandler) Handle(c *web.HttpContext) web.IActionResult {
    adminID := c.GetUserID()
    
    // ✅ 检查权限
    if err := h.checker.RequirePermission(adminID, "users", "create"); err != nil {
        return c.Forbidden("无权限创建用户")
    }
    
    // 业务逻辑...
    return c.Created(user)
}
```

**apps/admin/features/products/create_product.go:**

```go
package products

import "vertical_slice_demo/apps/admin/internal/auth"

type CreateProductHandler struct {
    productRepo IProductRepository
    checker     *auth.PermissionChecker  // ✅ 复用相同的权限检查
}

func (h *CreateProductHandler) Handle(c *web.HttpContext) web.IActionResult {
    adminID := c.GetUserID()
    
    // ✅ 检查权限
    if err := h.checker.RequirePermission(adminID, "products", "create"); err != nil {
        return c.Forbidden("无权限创建商品")
    }
    
    // 业务逻辑...
    return c.Created(product)
}
```

### 好处

- ✅ 应用内复用：多个功能共享应用特定逻辑
- ✅ 不跨应用：API 端无法访问 Admin 端的内部代码
- ✅ 统一管理：应用特定的工具集中管理
- ✅ 易于维护：修改权限逻辑只需改一处

### 注意事项

⚠️ **避免应用级 internal 变成垃圾堆：**
- 只放真正需要应用内共享的代码
- 不要把功能特定的代码放这里
- 不要把核心业务实体放这里（应该在 shared/domain）

---

## 全局共享

### 适用场景

**问题：** 多个应用端都需要使用同一份代码

**解决方案：** 使用 `shared/`

### 何时使用

- ✅ 核心业务实体（User, Product, Order）
- ✅ 跨端数据访问（UserRepository）
- ✅ 跨端业务服务（OrderService）
- ✅ 跨端 DTO（UserResponse）

### 目录结构

```
shared/
├── domain/                          # 共享领域模型
│   ├── user.go
│   ├── product.go
│   ├── order.go
│   └── common/
│       ├── base_entity.go
│       └── soft_delete.go
│
├── repositories/                    # 共享仓储
│   ├── user_repository.go
│   ├── product_repository.go
│   └── order_repository.go
│
├── services/                        # 共享服务
│   ├── order_service.go
│   └── notification_service.go
│
└── contracts/
    └── dtos/                        # 跨端共享 DTO
        ├── user_response.go
        └── product_response.go
```

### 代码示例

**shared/domain/user.go:**

```go
package domain

import "vertical_slice_demo/shared/domain/common"

// User 核心用户实体（多端共享）
type User struct {
    common.BaseEntity
    common.SoftDelete
    
    Name     string `gorm:"size:100;not null"`
    Email    string `gorm:"size:255;uniqueIndex;not null"`
    Password string `gorm:"size:255;not null"`
    Salt     string `gorm:"size:50;not null"`
    Role     string `gorm:"size:20;not null;default:'user'"`
}

func (User) TableName() string {
    return "users"
}
```

**shared/repositories/user_repository.go:**

```go
package repositories

import "vertical_slice_demo/shared/domain"

type IUserRepository interface {
    Create(user *domain.User) error
    GetByID(id int64) (*domain.User, error)
    GetByEmail(email string) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id int64) error
}

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id int64) (*domain.User, error) {
    var user domain.User
    err := r.db.First(&user, id).Error
    return &user, err
}
```

**shared/contracts/dtos/user_response.go:**

```go
package dtos

// UserResponse 跨端共享的用户响应（Admin 和 API 端都用）
type UserResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}
```

**apps/admin/features/users/create_user.go:**

```go
package users

import (
    "vertical_slice_demo/shared/domain"
    "vertical_slice_demo/shared/repositories"
    "vertical_slice_demo/shared/contracts/dtos"
)

type CreateUserHandler struct {
    userRepo repositories.IUserRepository  // ✅ 使用共享仓储
}

func (h *CreateUserHandler) Handle(c *web.HttpContext) web.IActionResult {
    user := &domain.User{...}  // ✅ 使用共享 Domain
    h.userRepo.Create(user)
    
    // ✅ 使用共享 DTO
    return c.Created(dtos.UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        Role:  user.Role,
    })
}
```

**apps/api/features/auth/register.go:**

```go
package auth

import (
    "vertical_slice_demo/shared/domain"
    "vertical_slice_demo/shared/repositories"
    "vertical_slice_demo/shared/contracts/dtos"
)

type RegisterHandler struct {
    userRepo repositories.IUserRepository  // ✅ 复用相同仓储
}

func (h *RegisterHandler) Handle(c *web.HttpContext) web.IActionResult {
    user := &domain.User{...}  // ✅ 复用相同 Domain
    h.userRepo.Create(user)
    
    // ✅ 复用相同 DTO，保证 API 格式一致
    return c.Created(dtos.UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        Role:  "user",
    })
}
```

### 好处

- ✅ 代码复用：避免重复实现
- ✅ 数据一致性：多端操作同一份数据
- ✅ API 一致性：使用共享 DTO 保证格式统一
- ✅ 易于维护：修改一处，所有端生效

### 注意事项

⚠️ **避免过早抽象到 shared：**
- 至少有 2 个端真正需要时才放 shared
- 不要"可能会用到"就提前抽象
- 第一次实现时放在功能内，第二个端需要时再提取

---

## 决策树：代码应该共享到哪里

```
这段代码需要被共享吗？
│
├─ 不需要共享
│  └─→ 保留在 features/xxx/ 内
│
└─ 需要共享
   │
   ├─ 只在当前功能内共享？
   │  └─→ features/xxx/internal/
   │     └─ 示例：报表功能的 ReportGenerator
   │
   ├─ 在当前应用内共享，但不跨应用？
   │  │
   │  ├─ 是工具/辅助代码？
   │  │  └─→ apps/*/internal/utils/
   │  │     └─ 示例：Admin 端的 PermissionChecker
   │  │
   │  ├─ 是中间件？
   │  │  └─→ apps/*/internal/middleware/
   │  │     └─ 示例：Admin 端的 AuthMiddleware
   │  │
   │  ├─ 是应用特定逻辑？
   │  │  └─→ apps/*/internal/
   │  │     └─ 示例：Admin 端的缓存管理
   │  │
   │  └─ 是核心业务实体？
   │     └─→ 重新考虑，可能应该在 shared/domain
   │
   └─ 需要跨应用共享？
      │
      ├─ 是核心业务实体？
      │  └─→ shared/domain/
      │     └─ 示例：User, Product, Order
      │
      ├─ 是数据访问？
      │  └─→ shared/repositories/
      │     └─ 示例：UserRepository
      │
      ├─ 是业务服务？
      │  └─→ shared/services/
      │     └─ 示例：OrderService
      │
      └─ 是 DTO？
         └─→ shared/contracts/dtos/
            └─ 示例：UserResponse
```

---

## 对比表：不同共享层次

| 层次 | 位置 | 可见性 | 适用场景 | 示例 | 提取时机 |
|------|------|--------|---------|------|---------|
| **功能私有** | `features/xxx/` | 仅当前功能 | 功能特有代码 | create_order.go | 默认 |
| **功能级共享** | `features/xxx/internal/` | 仅当前功能（内部） | 功能内部复用 | ReportGenerator | 功能内有复用需求 |
| **应用级共享** | `apps/*/internal/` | 仅当前应用 | 应用特定逻辑 | PermissionChecker | 应用内2+功能需要 |
| **全局共享** | `shared/` | 所有应用 | 核心业务代码 | User, UserRepository | 2+应用需要 |

---

## 最佳实践

### 1. 默认不共享

```go
// ✅ 新功能默认放在功能内
features/tags/
├── create_tag.go    # 默认不共享
├── list_tags.go
└── models.go

// ❌ 不要一开始就想着共享
shared/tag_service.go  // 过早抽象
```

**原则：** 先实现功能，确定需要共享时再提取

---

### 2. 逐级提升

```go
// 第一步：功能内
features/reports/
└── generate_report.go

// 第二步：功能内有复用 → 提取到 internal
features/reports/
├── internal/
│   └── business/
│       └── report_generator.go
└── generate_report.go

// 第三步：应用内多个功能需要 → 提取到 apps/*/internal
apps/admin/internal/
└── reporting/
    └── report_generator.go

// 第四步：多个应用需要 → 提取到 shared
shared/services/
└── report_service.go
```

**原则：** 不要跳级，逐步提升共享范围

---

### 3. 避免过早抽象到 shared

```go
// ❌ 错误：只有一个地方用就放 shared
shared/domain/admin_log.go     // 只有 Admin 用
shared/services/admin_service.go

// ✅ 正确：先放在功能内
apps/admin/features/logs/
└── models/
    └── admin_log.go

// 第二个端也需要时再考虑提取
```

**原则：** 至少有 2 个端真正需要时才放 shared

---

### 4. 应用级 internal 谨慎使用

```go
// ✅ 好的使用：真正的应用级共享
apps/admin/internal/
├── auth/
│   └── permission_checker.go   # 多个功能都需要
└── middleware/
    └── auth_middleware.go      # 应用级中间件

// ❌ 不好的使用：应该放功能内
apps/admin/internal/
└── user_validator.go           # 只有 users 功能用
```

**原则：** 确认真的需要应用内共享

---

### 5. shared 只放纯粹共享的代码

```go
// ✅ 好的 shared：核心业务
shared/
├── domain/
│   ├── user.go           # 核心实体
│   └── product.go
├── repositories/
│   └── user_repository.go
└── contracts/
    └── dtos/
        └── user_response.go

// ❌ 不好的 shared：应用特定逻辑
shared/
├── admin_permission.go   # 应该在 apps/admin/internal
└── api_rate_limiter.go   # 应该在 apps/api/internal
```

**原则：** shared 只放跨应用的核心业务代码

---

## 实战案例

### 案例1：报表功能的共享演进

**阶段1：功能内实现**

```go
features/reports/
└── generate_report.go    // 所有逻辑在一起
```

**阶段2：功能内有复用 → internal**

```go
features/reports/
├── internal/
│   ├── entity/
│   │   └── report_entity.go
│   ├── data/
│   │   └── report_store.go
│   └── business/
│       ├── report_generator.go    # 多个 Handler 共享
│       └── data_aggregator.go     # 多个 Handler 共享
│
├── generate_report.go
├── export_report.go
└── schedule_report.go
```

**阶段3：其他功能也需要 → 应用级 internal（如果只是 Admin 端需要）**

```go
apps/admin/internal/
└── reporting/
    └── report_generator.go    # Admin 端多个功能共享

features/reports/
└── generate_report.go         # 使用应用级共享

features/dashboard/
└── dashboard_handler.go       # 也使用报表生成
```

**阶段4：API 端也需要 → shared（跨应用）**

```go
shared/services/
└── report_service.go          # 跨应用共享

apps/admin/features/reports/
└── generate_report.go         # 使用 shared

apps/api/features/reports/
└── get_report.go              # 也使用 shared
```

---

### 案例2：权限检查的共享

**Admin 端权限检查：**

```go
// apps/admin/internal/auth/permission_checker.go
// Admin 端特定的权限逻辑（不跨应用）

apps/admin/features/users/
└── create_user.go             # 使用 admin/internal/auth

apps/admin/features/products/
└── create_product.go          # 使用 admin/internal/auth
```

**API 端权限检查：**

```go
// apps/api/internal/auth/permission_checker.go
// API 端特定的权限逻辑（不同于 Admin）

apps/api/features/orders/
└── create_order.go            # 使用 api/internal/auth
```

**为什么不共享：** Admin 和 API 的权限逻辑完全不同，不应该共享

---

### 案例3：跨端用户服务

```go
// shared/services/user_service.go
// 核心用户服务（Admin 和 API 都需要）

type UserService struct {
    userRepo repositories.IUserRepository
}

func (s *UserService) CreateUser(name, email, password string) (*domain.User, error) {
    // 通用的用户创建逻辑
}

// apps/admin/features/users/create_user.go
// Admin 端使用
func (h *CreateUserHandler) Handle(c *web.HttpContext) web.IActionResult {
    user, _ := h.userService.CreateUser(req.Name, req.Email, req.Password)
    return c.Created(toUserResponse(user))
}

// apps/api/features/auth/register.go
// API 端也使用
func (h *RegisterHandler) Handle(c *web.HttpContext) web.IActionResult {
    user, _ := h.userService.CreateUser(req.Name, req.Email, req.Password)
    return c.Created(toUserResponse(user))
}
```

---

## 总结

**共享策略的核心原则：**

1. **默认不共享** - 新代码先放功能内
2. **逐级提升** - features → internal → apps/internal → shared
3. **不要跳级** - 按需逐步提升范围
4. **功能级优先** - 大部分用 features/*/internal/
5. **应用级谨慎** - 确认需要应用内共享
6. **全局需谨慎** - 至少2个端真正需要
7. **避免过早抽象** - 不要"可能会用到"就提取

**判断共享层次的关键问题：**
- 只在功能内用 → features/xxx/internal/
- 应用内多处用 → apps/*/internal/
- 多个应用用 → shared/

---

**返回 [主文档](../ORGANIZATION_GUIDE.md)**
