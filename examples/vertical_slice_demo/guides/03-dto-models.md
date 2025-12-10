# DTO 模型组织

## 目录

- [DTO 核心概念](#dto-核心概念)
- [DTO 组织方式](#dto-组织方式)
- [DTO 共享策略](#dto-共享策略)
- [决策树：DTO 应该放在哪里](#决策树dto-应该放在哪里)
- [DTO 命名规范](#dto-命名规范)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)

---

## DTO 核心概念

### 什么是 DTO

DTO（Data Transfer Object）是 API 层和业务层之间传输数据的对象，包括 Request、Response、ListItem 等。

**DTO 的作用：**
- 定义 API 输入格式（Request）
- 定义 API 输出格式（Response）
- 与 Domain 模型解耦
- 隐藏敏感信息
- 添加验证规则

### DTO vs Domain 模型

| 维度 | Domain 模型 | DTO |
|------|------------|-----|
| **位置** | `shared/domain/` | `features/*/models.go` |
| **用途** | 内部业务逻辑，数据库映射 | 外部数据传输，API 通信 |
| **字段** | 完整字段（包含敏感） | 只有必要字段 |
| **标签** | gorm 标签 | json 标签，binding 标签 |
| **验证** | 业务规则验证 | 输入格式验证 |

**示例对比：**

```go
// shared/domain/user.go - Domain 模型
type User struct {
    ID        int64
    Name      string
    Email     string
    Password  string    // ✅ 包含敏感字段
    Salt      string
    IsDeleted bool
    CreatedAt time.Time
}

// features/users/models.go - DTO
type UserResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    // ❌ 不包含 Password, Salt, IsDeleted
}

type CreateUserRequest struct {
    Name     string `json:"name" binding:"required,min=2,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}
```

---

## DTO 组织方式

### 方案1：与操作放在一起

**适用场景：** 每个操作的 DTO 都不同，< 3 个 DTO

**目录结构：**
```
features/users/
├── create_user.go              # ✅ Request/Response 和逻辑在一起
├── list_users.go
├── update_user.go
├── controller.go
└── service_extensions.go
```

**示例：**

```go
// features/users/create_user.go
package users

import "github.com/gocrud/csgo/web"

// ===== DTO 定义 =====
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

type CreateUserResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// ===== Handler =====
type CreateUserHandler struct {
    userRepo IUserRepository
}

func (h *CreateUserHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    // 业务逻辑...
    user := &domain.User{
        Name:  req.Name,
        Email: req.Email,
    }
    h.userRepo.Create(user)
    
    // 返回响应
    response := &CreateUserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }
    
    return c.Created(response)
}
```

**优点：**
- ✅ 功能完全内聚
- ✅ 一个文件看到所有相关代码
- ✅ 符合垂直切片理念

**缺点：**
- ❌ 多个操作可能重复定义
- ❌ DTO 不能跨操作复用

---

### 方案2：功能内 models.go ⭐推荐

**适用场景：** 功能内多个操作共享 DTO，3-10 个 DTO

**目录结构：**
```
features/orders/
├── models.go                   # ✅ 功能内共享的 DTO
├── create_order.go
├── list_orders.go
├── update_order_status.go
├── get_order_detail.go
├── controller.go
└── service_extensions.go
```

**models.go:**

```go
package orders

import "time"

// ===== 共享的响应结构 =====
type OrderResponse struct {
    ID         int64       `json:"id"`
    UserID     int64       `json:"user_id"`
    TotalPrice float64     `json:"total_price"`
    Status     string      `json:"status"`
    Items      []OrderItem `json:"items"`
    CreatedAt  time.Time   `json:"created_at"`
}

type OrderItem struct {
    ProductID   int64   `json:"product_id"`
    ProductName string  `json:"product_name"`
    Quantity    int     `json:"quantity"`
    Price       float64 `json:"price"`
}

// ===== Request DTOs =====
type CreateOrderRequest struct {
    Items []CreateOrderItem `json:"items" binding:"required,min=1"`
}

type CreateOrderItem struct {
    ProductID int64 `json:"product_id" binding:"required"`
    Quantity  int   `json:"quantity" binding:"required,gt=0"`
}

type UpdateOrderStatusRequest struct {
    Status string `json:"status" binding:"required,oneof=pending paid shipped completed cancelled"`
}

// ===== 列表项（简化版）=====
type OrderListItem struct {
    ID         int64     `json:"id"`
    TotalPrice float64   `json:"total_price"`
    Status     string    `json:"status"`
    ItemCount  int       `json:"item_count"`
    CreatedAt  time.Time `json:"created_at"`
}

// ===== 详情（组合版）=====
type OrderDetailResponse struct {
    Order    OrderResponse `json:"order"`
    User     UserInfo      `json:"user"`
    Products []ProductInfo `json:"products"`
}

type UserInfo struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

type ProductInfo struct {
    ID    int64   `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

// ===== 转换函数 =====
func toOrderResponse(order *domain.Order) *OrderResponse {
    return &OrderResponse{
        ID:         order.ID,
        Status:     order.Status,
        TotalPrice: order.TotalPrice,
        Items:      toOrderItems(order.Items),
        CreatedAt:  order.CreatedAt,
    }
}

func toOrderListItem(order *domain.Order) *OrderListItem {
    return &OrderListItem{
        ID:         order.ID,
        TotalPrice: order.TotalPrice,
        Status:     order.Status,
        ItemCount:  len(order.Items),
        CreatedAt:  order.CreatedAt,
    }
}
```

**create_order.go:**

```go
package orders

func (h *CreateOrderHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req CreateOrderRequest  // ✅ 使用 models.go 中的 DTO
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    order := createOrder(req)
    return c.Created(toOrderResponse(order))  // ✅ 使用转换函数
}
```

**list_orders.go:**

```go
package orders

func (h *ListOrdersHandler) Handle(c *web.HttpContext) web.IActionResult {
    orders, _ := h.orderRepo.List()
    
    items := make([]OrderListItem, len(orders))  // ✅ 使用 models.go 中的 DTO
    for i, order := range orders {
        items[i] = toOrderListItem(order)  // ✅ 使用转换函数
    }
    
    return c.Ok(items)
}
```

**优点：**
- ✅ 功能内 DTO 集中管理
- ✅ 易于复用和维护
- ✅ 平衡了复用和独立性
- ✅ 80% 的场景适用 ⭐⭐⭐⭐⭐

**缺点：**
- ❌ models.go 可能会变大（> 10 个 DTO 时）

---

### 方案3：requests/responses 分类

**适用场景：** DTO 超过 10 个，需要清晰分类

**目录结构：**
```
features/reports/
├── requests/                   # 请求 DTO
│   ├── generate_report.go
│   ├── export_report.go
│   └── schedule_report.go
│
├── responses/                  # 响应 DTO
│   ├── report_detail.go
│   ├── report_list.go
│   └── report_summary.go
│
├── models/                     # 内部模型（requests 和 responses 共享）
│   ├── filter.go
│   ├── pagination.go
│   └── chart.go
│
├── generate_report.go          # Handler
├── export_report.go
├── list_reports.go
├── controller.go
└── service_extensions.go
```

**requests/generate_report.go:**

```go
package requests

import "time"

type GenerateReportRequest struct {
    Name       string    `json:"name" binding:"required"`
    ReportType string    `json:"report_type" binding:"required,oneof=sales user inventory"`
    StartDate  time.Time `json:"start_date" binding:"required"`
    EndDate    time.Time `json:"end_date" binding:"required"`
    GroupBy    string    `json:"group_by" binding:"required,oneof=day week month"`
    ChartTypes []string  `json:"chart_types"`
    Filters    []Filter  `json:"filters"`
}

type Filter struct {
    Field    string      `json:"field"`
    Operator string      `json:"operator"`
    Value    interface{} `json:"value"`
}
```

**responses/report_detail.go:**

```go
package responses

import "time"

type ReportDetailResponse struct {
    ID          int64         `json:"id"`
    Name        string        `json:"name"`
    Type        string        `json:"report_type"`
    Charts      []ChartData   `json:"charts"`
    Summary     ReportSummary `json:"summary"`
    GeneratedAt time.Time     `json:"generated_at"`
}

type ChartData struct {
    Type   string      `json:"type"`
    Title  string      `json:"title"`
    Data   interface{} `json:"data"`
    Labels []string    `json:"labels"`
}

type ReportSummary struct {
    TotalOrders   int     `json:"total_orders"`
    TotalRevenue  float64 `json:"total_revenue"`
    AvgOrderValue float64 `json:"avg_order_value"`
}
```

**models/filter.go (共享结构):**

```go
package models

// Filter 在 Request 和 Response 中都使用
type Filter struct {
    Field    string      `json:"field"`
    Operator string      `json:"operator"`
    Value    interface{} `json:"value"`
}

type Pagination struct {
    Page       int `json:"page"`
    PageSize   int `json:"page_size"`
    Total      int `json:"total"`
    TotalPages int `json:"total_pages"`
}
```

**generate_report.go (Handler):**

```go
package reports

import (
    "vertical_slice_demo/apps/admin/features/reports/requests"
    "vertical_slice_demo/apps/admin/features/reports/responses"
)

func (h *GenerateReportHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req requests.GenerateReportRequest  // ✅ 使用 requests 包
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    // 业务逻辑...
    report, _ := h.generator.Generate(req)
    
    // 转换为响应
    response := &responses.ReportDetailResponse{...}  // ✅ 使用 responses 包
    return c.Ok(response)
}
```

**优点：**
- ✅ 分类清晰，易于查找
- ✅ 适合 DTO 很多的场景
- ✅ Request 和 Response 分离明确
- ✅ 可以单独为 DTO 编写文档

**缺点：**
- ❌ 目录层级增加
- ❌ 跨目录引用
- ❌ 简单功能会显得复杂

---

### 方案对比

| 方案 | 适用场景 | DTO 位置 | 复用性 | 内聚性 | 推荐度 |
|------|---------|---------|--------|--------|--------|
| 与操作一起 | 简单功能，DTO 差异大 | 每个 handler 文件中 | ❌ 低 | ⭐⭐⭐⭐⭐ 高 | ⭐⭐⭐⭐ |
| 功能内 models.go | 中等复杂，部分共享 | features/xxx/models.go | ✅ 中 | ⭐⭐⭐⭐ 高 | ⭐⭐⭐⭐⭐ |
| 分类组织 | 复杂功能，DTO > 10 | requests/ responses/ | ✅ 中 | ⭐⭐⭐ 中 | ⭐⭐⭐ |

---

## DTO 共享策略

### 1. 私有 DTO（功能内独有）

**位置：** `features/*/models.go` 或操作文件内

**适用场景：** DTO 只在当前功能内使用

**示例：**

```go
// apps/admin/features/tags/models.go
package tags

// ✅ CreateTagRequest 只在 tags 功能内使用
type CreateTagRequest struct {
    Name  string `json:"name" binding:"required"`
    Color string `json:"color" binding:"required"`
}

type TagResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Color string `json:"color"`
}
```

**何时使用：**
- ✅ 功能特有的 Request/Response
- ✅ 不需要跨功能共享
- ✅ 大部分场景（80%）

---

### 2. 端内共享 DTO（应用内共享）

**位置：** `apps/*/shared/dtos/`

**适用场景：** 多个功能需要共享，但不跨端

**示例：**

```
apps/admin/shared/dtos/
├── pagination.go        # Admin 端内多个功能共享
├── filter.go
└── sort_option.go
```

```go
// apps/admin/shared/dtos/pagination.go
package dtos

// Pagination Admin 端的分页结构（多个功能共享）
type Pagination struct {
    Page       int `json:"page" binding:"gte=1"`
    PageSize   int `json:"page_size" binding:"gte=1,lte=100"`
    Total      int `json:"total"`
    TotalPages int `json:"total_pages"`
}

type PaginationRequest struct {
    Page     int `form:"page" binding:"gte=1"`
    PageSize int `form:"page_size" binding:"gte=1,lte=100"`
}

type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Pagination Pagination  `json:"pagination"`
}
```

**使用方式：**

```go
// apps/admin/features/users/list_users.go
import "vertical_slice_demo/apps/admin/shared/dtos"

func (h *ListUsersHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req dtos.PaginationRequest  // ✅ 使用端内共享 DTO
    c.BindQuery(&req)
    
    users, total := h.userRepo.List(req.Page, req.PageSize)
    
    return c.Ok(dtos.PaginatedResponse{
        Data: users,
        Pagination: dtos.Pagination{
            Page:       req.Page,
            PageSize:   req.PageSize,
            Total:      total,
            TotalPages: (total + req.PageSize - 1) / req.PageSize,
        },
    })
}
```

**何时使用：**
- ✅ Admin 端多个功能需要相同的分页结构
- ✅ 端内统一的过滤条件格式
- ✅ 端内统一的错误响应格式
- ❌ 不是跨端共享

---

### 3. 全局共享 DTO（多端共享）

**位置：** `shared/contracts/dtos/`

**适用场景：** 多个端需要保证 API 格式一致

**示例：**

```
shared/contracts/dtos/
├── user_response.go       # Admin 和 API 端返回相同格式
├── product_response.go
└── order_response.go
```

```go
// shared/contracts/dtos/user_response.go
package dtos

// UserResponse 跨端共享的用户响应格式
type UserResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}
```

**使用方式：**

```go
// apps/admin/features/users/create_user.go
package users

import "vertical_slice_demo/shared/contracts/dtos"

// ✅ Request 独立定义（Admin 可以指定角色）
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Role     string `json:"role" binding:"required,oneof=admin user"`
}

func (h *CreateUserHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    c.MustBindJSON(&req)
    
    user := createUser(req)
    
    // ✅ Response 使用共享的
    return c.Created(dtos.UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        Role:  user.Role,
    })
}
```

```go
// apps/api/features/auth/register.go
package auth

import "vertical_slice_demo/shared/contracts/dtos"

// ✅ Request 独立定义（C端不能指定角色）
type RegisterRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    // ❌ 没有 Role 字段
}

func (h *RegisterHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req RegisterRequest
    c.MustBindJSON(&req)
    
    user := registerUser(req)
    
    // ✅ Response 使用相同的共享 DTO
    return c.Created(dtos.UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        Role:  "user", // C端默认角色
    })
}
```

**决策规则：**
- ✅ Response 倾向共享（保证一致性）
- ✅ Request 倾向独立（保持灵活性）

**何时使用：**
- ✅ 多端返回相同格式的数据
- ✅ 需要保证 API 一致性
- ❌ 不是格式差异很大的情况

---

## 决策树：DTO 应该放在哪里

```
DTO 需要被多个端使用吗？
│
├─ 是 → Response 格式相同吗？
│  │
│  ├─ 是 → shared/contracts/dtos/ ⭐
│  │  └─ 保证 API 一致性
│  │
│  └─ 否 → 各端独立
│     ├─ apps/admin/features/xxx/models.go
│     └─ apps/api/features/xxx/models.go
│
└─ 否 → 端内多个功能使用吗？
   │
   ├─ 是 → apps/*/shared/dtos/
   │  └─ 端内共享（如：Pagination）
   │
   └─ 否 → 功能内 DTO
      │
      ├─ < 3 个 DTO → 与操作放在一起
      │  └─ features/xxx/create_xxx.go
      │
      ├─ 3-10 个 DTO → 功能内 models.go ⭐
      │  └─ features/xxx/models.go
      │
      └─ > 10 个 DTO → 分类组织
         └─ features/xxx/
            ├─ requests/
            ├─ responses/
            └─ models/
```

---

## DTO 命名规范

### Request DTOs

```go
// ✅ 好的命名
type CreateUserRequest struct      // 动作 + 实体 + Request
type UpdateProductRequest struct
type SearchOrdersRequest struct

// ❌ 不好的命名
type UserCreateDTO struct          // 顺序混乱
type CreateReq struct              // 缩写不清晰
type UserData struct               // 语义不明
```

### Response DTOs

```go
// ✅ 好的命名
type UserResponse struct           // 实体 + Response（通用响应）
type UserDetailResponse struct     // 实体 + 具体用途 + Response
type OrderSummaryResponse struct

// ❌ 不好的命名
type UserDTO struct                // 过于宽泛
type GetUserResp struct            // 缩写
type UserOutput struct             // 不常见
```

### List Item DTOs

```go
// ✅ 好的命名
type UserListItem struct           // 实体 + ListItem
type ProductListItem struct
type OrderSummaryItem struct       // 或 SummaryItem

// ❌ 不好的命名
type UserList struct               // 容易和 []User 混淆
type UserInList struct             // 冗余
type UserItem struct               // 不够明确
```

---

## 最佳实践

### 1. 优先功能内聚

```go
// ✅ 好的做法：DTO 跟随功能
features/orders/
├── models.go           # 订单相关的所有 DTO
└── create_order.go

// ❌ 不好的做法：过早抽象到全局
shared/dtos/order_dto.go  # 只有一个地方用
```

### 2. 按需复用

```go
// ✅ 好的做法：只在功能内复用
// features/orders/models.go
type OrderResponse struct {...}  // 在功能内多处使用

// ❌ 不好的做法："可能会用到"就提取
shared/dtos/order_dto.go  # 只有一个功能用
```

### 3. 清晰命名

```go
// ✅ 好的命名
type CreateUserRequest struct     // 清晰的动作 + 实体 + 类型
type UserDetailResponse struct    // 实体 + 具体用途 + 类型  
type ProductListItem struct       // 实体 + ListItem

// ❌ 不好的命名
type UserDTO struct               // 过于宽泛，不知道用途
type CreateReq struct             // 缺少实体名称
type Data struct                  // 语义不明
```

### 4. 隐藏敏感信息

```go
// shared/domain/user.go - Domain 模型
type User struct {
    ID       int64
    Name     string
    Password string    // ✅ 在 Domain 中
    Salt     string
}

// features/users/models.go - DTO
type UserResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    // ❌ 绝不包含 Password、Salt
}

// 或使用 json:"-" 标签显式忽略
type UserDTO struct {
    ID       int64  `json:"id"`
    Name     string `json:"name"`
    Password string `json:"-"`  // ✅ 不会序列化到 JSON
}
```

### 5. 验证规则在 DTO 中

```go
// ✅ 好的做法：使用 binding 标签
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required,min=2,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8,max=100"`
    Age      int    `json:"age" binding:"gte=0,lte=150"`
}

// Handler 中自动验证
func (h *CreateUserHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err  // 自动返回 400 + 验证错误信息
    }
    // req 已经通过验证
}
```

### 6. 转换函数位置

**简单转换：放在 handler 中**

```go
// features/users/create_user.go
func (h *CreateUserHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    c.MustBindJSON(&req)
    
    // ✅ 简单转换直接在这里
    user := &domain.User{
        Name:  req.Name,
        Email: req.Email,
    }
    
    h.userRepo.Create(user)
    
    return c.Created(UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    })
}
```

**复杂转换：提取到 models.go**

```go
// features/orders/models.go
func toOrderResponse(order *domain.Order) *OrderResponse {
    return &OrderResponse{
        ID:         order.ID,
        Status:     order.Status,
        TotalPrice: order.TotalPrice,
        Items:      toOrderItems(order.Items),
        CreatedAt:  order.CreatedAt,
    }
}

// features/orders/create_order.go
func (h *CreateOrderHandler) Handle(c *web.HttpContext) web.IActionResult {
    order := createOrder(req)
    return c.Created(toOrderResponse(order))  // ✅ 使用转换函数
}
```

---

## 常见问题

### 问题1：DTO 应该放在哪里？

**场景：** 开发新功能时不确定 DTO 位置

**解决方案：** 根据 DTO 数量选择
- < 3 个 → 与操作放在一起
- 3-10 个 → 功能内 models.go ⭐
- \> 10 个 → requests/responses 分类

---

### 问题2：多端需要相同的 DTO 格式怎么办？

**场景：** Admin 创建用户和 API 注册用户返回格式相同

**解决方案：**
- Response 相同 → shared/contracts/dtos/
- Request 不同 → 各端独立

**示例：**

```go
// shared/contracts/dtos/user_response.go
type UserResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}

// apps/admin/features/users/models.go
type CreateUserRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Role     string `json:"role"`  // ✅ Admin 可以指定角色
}

// apps/api/features/auth/models.go
type RegisterRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
    // ❌ 没有 Role（C端不能指定）
}
```

---

### 问题3：如何区分 Domain 模型和 DTO？

**关键区别：**

| 维度 | Domain 模型 | DTO |
|------|------------|-----|
| 位置 | `shared/domain/` | `features/*/models.go` |
| 用途 | 内部业务逻辑 | 外部数据传输 |
| 字段 | 完整（包含敏感） | 必要字段 |
| 标签 | gorm 标签 | json + binding 标签 |

**示例：**

```go
// shared/domain/user.go - Domain 模型（内部使用）
type User struct {
    ID        int64     `gorm:"primaryKey"`
    Name      string    `gorm:"size:100;not null"`
    Email     string    `gorm:"size:255;uniqueIndex;not null"`
    Password  string    `gorm:"size:255;not null"` // ✅ 包含敏感字段
    Salt      string    `gorm:"size:50;not null"`
    IsDeleted bool      `gorm:"default:false"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// features/users/models.go - DTO（对外传输）
type UserResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    // ❌ 不包含 Password、Salt、IsDeleted
}
```

---

## 总结

**DTO 组织的核心原则：**

1. **80% 场景用 models.go** - 功能内集中管理 DTO
2. **Response 倾向共享** - 保证 API 一致性
3. **Request 倾向独立** - 保持各端灵活性
4. **隐藏敏感信息** - 永远不暴露 Password、Salt
5. **添加验证规则** - 使用 binding 标签
6. **清晰命名** - CreateXxxRequest, XxxResponse
7. **按需复用** - 不要过早抽象

---

**返回 [主文档](../ORGANIZATION_GUIDE.md)**
