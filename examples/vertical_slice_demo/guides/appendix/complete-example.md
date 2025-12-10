# 完整项目示例

## 目录

- [项目结构总览](#项目结构总览)
- [shared 目录详解](#shared-目录详解)
- [apps/admin 目录详解](#appsadmin-目录详解)
- [apps/api 目录详解](#appsapi-目录详解)
- [关键文件示例](#关键文件示例)

---

## 项目结构总览

```
vertical_slice_demo/
├── main.go
├── go.mod
├── go.sum
│
├── shared/                              # 跨端共享代码
│   ├── domain/                          # 共享领域模型
│   │   ├── common/                      # 公共基础模型
│   │   │   ├── base_entity.go
│   │   │   └── soft_delete.go
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── category.go
│   │   ├── order.go
│   │   └── order_item.go
│   │
│   ├── repositories/                    # 共享仓储
│   │   ├── user_repository.go
│   │   ├── product_repository.go
│   │   ├── category_repository.go
│   │   └── order_repository.go
│   │
│   ├── services/                        # 共享服务
│   │   ├── order_service.go
│   │   └── notification_service.go
│   │
│   └── contracts/                       # 共享契约
│       └── dtos/                        # 跨端共享 DTO
│           ├── user_response.go
│           ├── product_response.go
│           └── order_response.go
│
├── apps/                                # 应用端
│   ├── admin/                           # 管理端
│   │   ├── main.go
│   │   ├── internal/                    # Admin 应用内部共享
│   │   │   ├── auth/
│   │   │   │   └── permission_checker.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth_middleware.go
│   │   │   │   └── logging_middleware.go
│   │   │   └── utils/
│   │   │       └── pagination_helper.go
│   │   │
│   │   └── features/                    # 功能模块
│   │       ├── users/                   # 用户管理（模式2：按操作拆分）
│   │       │   ├── models.go
│   │       │   ├── create_user.go
│   │       │   ├── list_users.go
│   │       │   ├── update_user.go
│   │       │   ├── delete_user.go
│   │       │   ├── controller.go
│   │       │   └── service_extensions.go
│   │       │
│   │       ├── products/                # 商品管理（模式2）
│   │       │   ├── models.go
│   │       │   ├── create_product.go
│   │       │   ├── list_products.go
│   │       │   ├── update_product.go
│   │       │   ├── controller.go
│   │       │   └── service_extensions.go
│   │       │
│   │       ├── categories/              # 分类管理（模式1：单文件）
│   │       │   ├── handler.go
│   │       │   └── service_extensions.go
│   │       │
│   │       ├── orders/                  # 订单管理（模式2）
│   │       │   ├── models.go
│   │       │   ├── list_orders.go
│   │       │   ├── get_order_detail.go
│   │       │   ├── update_order_status.go
│   │       │   ├── controller.go
│   │       │   └── service_extensions.go
│   │       │
│   │       ├── logs/                    # 操作日志（模式1，私有模型）
│   │       │   ├── models.go
│   │       │   ├── handler.go
│   │       │   └── service_extensions.go
│   │       │
│   │       └── reports/                 # 报表系统（模式3：内部分层）
│   │           ├── internal/
│   │           │   ├── entity/
│   │           │   │   ├── report_entity.go
│   │           │   │   ├── template_entity.go
│   │           │   │   └── config_entity.go
│   │           │   ├── data/
│   │           │   │   ├── report_store.go
│   │           │   │   └── template_store.go
│   │           │   └── business/
│   │           │       ├── report_generator.go
│   │           │       ├── data_aggregator.go
│   │           │       └── chart_builder.go
│   │           ├── models.go
│   │           ├── generate_report.go
│   │           ├── export_report.go
│   │           ├── list_reports.go
│   │           ├── controller.go
│   │           └── service_extensions.go
│   │
│   └── api/                             # C端应用
│       ├── main.go
│       ├── internal/                    # API 应用内部共享
│       │   ├── auth/
│       │   │   └── jwt_validator.go
│       │   ├── middleware/
│       │   │   ├── auth_middleware.go
│       │   │   └── rate_limiter.go
│       │   └── utils/
│       │       └── response_helper.go
│       │
│       └── features/                    # 功能模块
│           ├── auth/                    # 认证（模式2）
│           │   ├── models.go
│           │   ├── register.go
│           │   ├── login.go
│           │   ├── logout.go
│           │   ├── controller.go
│           │   └── service_extensions.go
│           │
│           ├── products/                # 商品浏览（模式2）
│           │   ├── models.go
│           │   ├── list_products.go
│           │   ├── get_product_detail.go
│           │   ├── search_products.go
│           │   ├── controller.go
│           │   └── service_extensions.go
│           │
│           ├── cart/                    # 购物车（模式2）
│           │   ├── models.go
│           │   ├── add_to_cart.go
│           │   ├── get_cart.go
│           │   ├── update_cart_item.go
│           │   ├── remove_from_cart.go
│           │   ├── controller.go
│           │   └── service_extensions.go
│           │
│           └── orders/                  # 订单（模式2）
│               ├── models.go
│               ├── create_order.go
│               ├── list_my_orders.go
│               ├── get_order_detail.go
│               ├── cancel_order.go
│               ├── controller.go
│               └── service_extensions.go
│
└── configs/                             # 配置文件
    ├── database.yml
    ├── redis.yml
    └── app.yml
```

---

## shared 目录详解

### domain/ - 共享领域模型

**common/base_entity.go:**
```go
package common

import "time"

type BaseEntity struct {
    ID        int64     `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
```

**user.go:**
```go
package domain

import "vertical_slice_demo/shared/domain/common"

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

### repositories/ - 共享仓储

**user_repository.go:**
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
```

### contracts/dtos/ - 跨端共享 DTO

**user_response.go:**
```go
package dtos

type UserResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}
```

---

## apps/admin 目录详解

### internal/ - Admin 应用内部共享

**auth/permission_checker.go:**
```go
package auth

type PermissionChecker struct {
    adminRepo IAdminRepository
}

func (p *PermissionChecker) CheckPermission(adminID int64, resource, action string) bool {
    // Admin 端特定的权限检查
}
```

### features/users/ - 用户管理（模式2）

**models.go:**
```go
package users

type CreateUserRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Role     string `json:"role" binding:"required,oneof=admin user"`
}

type UpdateUserRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Role  string `json:"role" binding:"required,oneof=admin user"`
}

type UserListItem struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}
```

**create_user.go:**
```go
package users

import (
    "vertical_slice_demo/shared/domain"
    "vertical_slice_demo/shared/repositories"
    "vertical_slice_demo/shared/contracts/dtos"
)

type CreateUserHandler struct {
    userRepo repositories.IUserRepository
}

func (h *CreateUserHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    c.MustBindJSON(&req)
    
    user := &domain.User{
        Name:  req.Name,
        Email: req.Email,
        Role:  req.Role,
    }
    
    h.userRepo.Create(user)
    
    return c.Created(dtos.UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        Role:  user.Role,
    })
}
```

### features/categories/ - 分类管理（模式1）

**handler.go:**
```go
package categories

type Category struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

type CategoryHandler struct {
    categories map[int64]*Category
    mu         sync.RWMutex
}

func (h *CategoryHandler) Create(c *web.HttpContext) web.IActionResult {
    // 所有逻辑在一起
}

func (h *CategoryHandler) List(c *web.HttpContext) web.IActionResult {
    // ...
}

func (h *CategoryHandler) MapRoutes(app *web.WebApplication) {
    g := app.MapGroup("/api/admin/categories")
    g.MapPost("", h.Create)
    g.MapGet("", h.List)
}
```

### features/reports/ - 报表系统（模式3）

**internal/entity/report_entity.go:**
```go
package entity

type ReportEntity struct {
    ID         int64
    Name       string
    Type       ReportType
    Config     string
    CreatedAt  time.Time
}
```

**internal/business/report_generator.go:**
```go
package business

type ReportGenerator struct {
    store      *data.ReportStore
    aggregator *DataAggregator
}

func (g *ReportGenerator) Generate(config entity.ReportConfig) (*entity.ReportEntity, error) {
    // 复杂的生成逻辑
}
```

**generate_report.go:**
```go
package reports

import "vertical_slice_demo/apps/admin/features/reports/internal/business"

type GenerateReportHandler struct {
    generator *business.ReportGenerator
}

func (h *GenerateReportHandler) Handle(c *web.HttpContext) web.IActionResult {
    // 使用内部业务逻辑
}
```

---

## apps/api 目录详解

### internal/ - API 应用内部共享

**auth/jwt_validator.go:**
```go
package auth

type JWTValidator struct {
    secret string
}

func (v *JWTValidator) ValidateToken(token string) (*Claims, error) {
    // JWT 验证逻辑
}
```

### features/auth/ - 认证（模式2）

**models.go:**
```go
package auth

type RegisterRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Token string              `json:"token"`
    User  dtos.UserResponse   `json:"user"`
}
```

**register.go:**
```go
package auth

import (
    "vertical_slice_demo/shared/domain"
    "vertical_slice_demo/shared/repositories"
    "vertical_slice_demo/shared/contracts/dtos"
)

type RegisterHandler struct {
    userRepo repositories.IUserRepository
}

func (h *RegisterHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req RegisterRequest
    c.MustBindJSON(&req)
    
    user := &domain.User{
        Name:  req.Name,
        Email: req.Email,
        Role:  "user", // C端默认角色
    }
    
    h.userRepo.Create(user)
    
    return c.Created(dtos.UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        Role:  user.Role,
    })
}
```

### features/orders/ - 订单（模式2）

**models.go:**
```go
package orders

type CreateOrderRequest struct {
    Items []CreateOrderItem `json:"items" binding:"required,min=1"`
}

type CreateOrderItem struct {
    ProductID int64 `json:"product_id" binding:"required"`
    Quantity  int   `json:"quantity" binding:"required,gt=0"`
}
```

**create_order.go:**
```go
package orders

import (
    "vertical_slice_demo/shared/domain"
    "vertical_slice_demo/shared/services"
)

type CreateOrderHandler struct {
    orderService *services.OrderService
}

func (h *CreateOrderHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req CreateOrderRequest
    c.MustBindJSON(&req)
    
    userID := c.GetUserID()
    order, _ := h.orderService.CreateOrder(userID, req.Items)
    
    return c.Created(toOrderResponse(order))
}
```

---

## 关键文件示例

### main.go（应用入口）

```go
package main

import (
    "vertical_slice_demo/apps/admin"
    "vertical_slice_demo/apps/api"
)

func main() {
    // 启动 Admin 端
    go admin.Start(":8080")
    
    // 启动 API 端
    api.Start(":8081")
}
```

### apps/admin/main.go

```go
package admin

import (
    "vertical_slice_demo/apps/admin/features/users"
    "vertical_slice_demo/apps/admin/features/products"
    "vertical_slice_demo/apps/admin/features/categories"
    "vertical_slice_demo/apps/admin/features/orders"
    "vertical_slice_demo/apps/admin/features/reports"
)

func Start(addr string) {
    app := web.NewApplication()
    
    // 注册 features
    users.AddUserFeature(app.Services)
    products.AddProductFeature(app.Services)
    categories.AddCategoryFeature(app.Services)
    orders.AddOrderFeature(app.Services)
    reports.AddReportFeature(app.Services)
    
    app.Run(addr)
}
```

### service_extensions.go 示例

```go
package users

import (
    "vertical_slice_demo/shared/repositories"
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

func AddUserFeature(services di.IServiceCollection) {
    // 注册仓储（如果还没注册）
    if !services.Has(repositories.IUserRepository) {
        services.AddSingleton(repositories.NewUserRepository)
    }
    
    // 注册 Handlers
    services.AddSingleton(NewCreateUserHandler)
    services.AddSingleton(NewListUsersHandler)
    services.AddSingleton(NewUpdateUserHandler)
    services.AddSingleton(NewDeleteUserHandler)
    
    // 注册 Controller
    web.AddController(services, NewUserController)
}
```

---

## 总结

这个完整示例展示了：

1. **清晰的目录结构** - shared、apps、features 分离
2. **三种组织模式** - 简单、中等、复杂功能的不同组织方式
3. **合理的共享策略** - 功能内、应用级、全局三个层次
4. **internal 的使用** - 封装复杂功能的内部实现
5. **跨端代码复用** - shared 目录下的 Domain、Repository、DTO

---

**返回 [主文档](../../ORGANIZATION_GUIDE.md)**
