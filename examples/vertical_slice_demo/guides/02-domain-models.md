# 数据库表模型组织

## 目录

- [Domain 模型定义](#domain-模型定义)
- [模型组织方式](#模型组织方式)
- [决策树：模型应该放在哪里](#决策树模型应该放在哪里)
- [表模型拆分策略](#表模型拆分策略)
- [模型字段设计规范](#模型字段设计规范)
- [最佳实践](#最佳实践)

---

## Domain 模型定义

### 什么是 Domain 模型

Domain 模型（领域模型）是对业务实体的抽象，通常与数据库表一一对应，包含完整的数据结构和业务规则。

**特点：**
- 包含所有字段（包括敏感字段）
- 用于内部业务逻辑
- 与数据库表映射
- 不直接暴露给外部

**示例：**

```go
// shared/domain/user.go
package domain

import "time"

type User struct {
    ID        int64     `gorm:"primaryKey"`
    Name      string    `gorm:"size:100;not null"`
    Email     string    `gorm:"size:255;uniqueIndex;not null"`
    Password  string    `gorm:"size:255;not null"` // ✅ 包含敏感字段
    Salt      string    `gorm:"size:50;not null"`  // ✅ 内部字段
    IsDeleted bool      `gorm:"default:false"`     // ✅ 软删除标记
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
    return "users"
}
```

---

### Domain 模型的职责

**Domain 模型负责：**
- ✅ 数据持久化（与数据库表映射）
- ✅ 业务规则封装
- ✅ 实体关联关系
- ✅ 领域逻辑

**Domain 模型不负责：**
- ❌ 数据传输（那是 DTO 的职责）
- ❌ 表示层逻辑（那是 Handler 的职责）
- ❌ API 格式定义（那是 Response DTO 的职责）

---

### Domain vs DTO 的本质区别

| 维度 | Domain 模型 | DTO |
|------|------------|-----|
| **位置** | `shared/domain/` 或 `features/*/internal/entity/` | `features/*/models.go` |
| **用途** | 内部业务逻辑，数据库映射 | 外部数据传输，API 通信 |
| **字段** | 完整字段（包含敏感） | 只有必要字段 |
| **标签** | gorm 标签，可能无 JSON | 必有 JSON 标签，binding 标签 |
| **验证** | 业务规则验证 | 输入格式验证 |
| **可见性** | 内部可见 | 外部可见 |

**代码对比：**

```go
// shared/domain/user.go - Domain 模型（内部使用）
package domain

type User struct {
    ID        int64     `gorm:"primaryKey"`
    Name      string    `gorm:"size:100;not null"`
    Email     string    `gorm:"size:255;uniqueIndex;not null"`
    Password  string    `gorm:"size:255;not null"` // ✅ 包含敏感字段
    Salt      string    `gorm:"size:50;not null"`  // ✅ 内部字段
    IsDeleted bool      `gorm:"default:false"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// features/users/models.go - DTO（对外传输）
package users

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

**转换函数：**

```go
// features/users/create_user.go
package users

import "vertical_slice_demo/shared/domain"

// DTO → Domain
func toUserDomain(req CreateUserRequest) *domain.User {
    return &domain.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashPassword(req.Password), // 加密处理
    }
}

// Domain → DTO
func toUserResponse(user *domain.User) *UserResponse {
    return &UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        // 不暴露敏感字段
    }
}
```

---

## 模型组织方式

### 1. 共享模型 (shared/domain/)

**适用场景：**
- 多个端都需要使用的核心实体
- 需要跨端保证数据一致性
- 核心业务领域模型

**目录结构：**
```
shared/domain/
├── user.go              # 用户实体
├── product.go           # 商品实体
├── order.go             # 订单实体
├── order_item.go        # 订单明细
└── common/              # 公共基础模型
    ├── base_entity.go   # 基础实体
    └── soft_delete.go   # 软删除
```

**示例：**

```go
// shared/domain/user.go
package domain

import (
    "time"
    "vertical_slice_demo/shared/domain/common"
)

type User struct {
    common.BaseEntity   // 嵌入基础字段
    common.SoftDelete   // 嵌入软删除
    
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

```go
// shared/domain/product.go
package domain

import "vertical_slice_demo/shared/domain/common"

type Product struct {
    common.BaseEntity
    common.SoftDelete
    
    Name        string  `gorm:"size:200;not null"`
    Description string  `gorm:"type:text"`
    Price       float64 `gorm:"type:decimal(10,2);not null"`
    Stock       int     `gorm:"not null;default:0"`
    CategoryID  int64   `gorm:"not null;index"`
}

func (Product) TableName() string {
    return "products"
}
```

```go
// shared/domain/order.go
package domain

import "vertical_slice_demo/shared/domain/common"

type Order struct {
    common.BaseEntity
    
    UserID     int64       `gorm:"not null;index"`
    TotalPrice float64     `gorm:"type:decimal(10,2);not null"`
    Status     OrderStatus `gorm:"size:20;not null;default:'pending'"`
    Items      []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderStatus string

const (
    OrderStatusPending   OrderStatus = "pending"
    OrderStatusPaid      OrderStatus = "paid"
    OrderStatusShipped   OrderStatus = "shipped"
    OrderStatusCompleted OrderStatus = "completed"
    OrderStatusCancelled OrderStatus = "cancelled"
)

func (Order) TableName() string {
    return "orders"
}

type OrderItem struct {
    common.BaseEntity
    
    OrderID   int64   `gorm:"not null;index"`
    ProductID int64   `gorm:"not null;index"`
    Quantity  int     `gorm:"not null"`
    Price     float64 `gorm:"type:decimal(10,2);not null"`
}

func (OrderItem) TableName() string {
    return "order_items"
}
```

**何时使用：**
- ✅ Admin 端和 API 端都需要操作用户
- ✅ 多个端共享订单数据
- ✅ 核心业务实体
- ❌ 不是仅某个端特有的数据

---

### 2. 私有模型 (功能内部)

**适用场景：**
- 单端特有的数据实体
- 不需要跨端共享
- 功能特定的业务数据

**位置选择：**
- **简单功能：** `apps/*/features/*/models.go`
- **复杂功能（使用 internal）：** `apps/*/features/*/internal/entity/`

**示例1：简单功能的私有模型**

```go
// apps/admin/features/logs/models.go
package logs

import "time"

// AdminLog 管理端操作日志（不需要在 API 端使用）
type AdminLog struct {
    ID         int64     `gorm:"primaryKey"`
    AdminID    int64     `gorm:"not null;index"`
    Action     string    `gorm:"size:50;not null"`
    Resource   string    `gorm:"size:100;not null"`
    ResourceID int64     `gorm:"not null"`
    Details    string    `gorm:"type:text"`
    IPAddress  string    `gorm:"size:45"`
    CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func (AdminLog) TableName() string {
    return "admin_logs"
}
```

**示例2：复杂功能的私有模型（使用 internal）**

```go
// apps/admin/features/reports/internal/entity/report_entity.go
package entity

import "time"

// ReportEntity 报表内部实体（只在 reports 功能内使用）
type ReportEntity struct {
    ID         int64       `gorm:"primaryKey"`
    Name       string      `gorm:"size:200;not null"`
    Type       ReportType  `gorm:"size:50;not null"`
    Config     string      `gorm:"type:json"`
    CreatedBy  int64       `gorm:"not null;index"`
    CreatedAt  time.Time   `gorm:"autoCreateTime"`
}

type ReportType string

const (
    ReportTypeSales     ReportType = "sales"
    ReportTypeUser      ReportType = "user"
    ReportTypeInventory ReportType = "inventory"
)

func (ReportEntity) TableName() string {
    return "reports"
}
```

**何时使用：**
- ✅ 管理端独有的操作日志
- ✅ C端独有的用户偏好设置
- ✅ 单端特定的业务数据
- ❌ 不是需要跨端共享的

---

### 3. 公共基础模型 (shared/domain/common/)

**适用场景：**
- 通用的基础字段
- 所有实体共享的结构
- 跨领域的公共概念

**目录结构：**
```
shared/domain/common/
├── base_entity.go       # 基础实体（ID、时间戳）
├── soft_delete.go       # 软删除标记
└── audit_fields.go      # 审计字段
```

**base_entity.go:**

```go
package common

import "time"

// BaseEntity 基础实体，所有 Domain 模型都应嵌入
type BaseEntity struct {
    ID        int64     `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
```

**soft_delete.go:**

```go
package common

import "time"

// SoftDelete 软删除标记
type SoftDelete struct {
    IsDeleted bool       `gorm:"default:false;index" json:"-"`
    DeletedAt *time.Time `gorm:"index" json:"-"`
}

// Delete 标记为已删除
func (s *SoftDelete) Delete() {
    s.IsDeleted = true
    now := time.Now()
    s.DeletedAt = &now
}

// IsActive 是否未删除
func (s *SoftDelete) IsActive() bool {
    return !s.IsDeleted
}
```

**audit_fields.go:**

```go
package common

// AuditFields 审计字段
type AuditFields struct {
    CreatedBy int64 `gorm:"not null" json:"created_by"`
    UpdatedBy int64 `gorm:"not null" json:"updated_by"`
}
```

**使用方式：**

```go
// shared/domain/product.go
package domain

import "vertical_slice_demo/shared/domain/common"

type Product struct {
    common.BaseEntity    // ✅ 嵌入：ID, CreatedAt, UpdatedAt
    common.SoftDelete    // ✅ 嵌入：IsDeleted, DeletedAt
    common.AuditFields   // ✅ 嵌入：CreatedBy, UpdatedBy
    
    Name  string  `gorm:"size:200;not null"`
    Price float64 `gorm:"type:decimal(10,2);not null"`
}
```

**何时使用：**
- ✅ 所有表都有 ID、CreatedAt、UpdatedAt
- ✅ 统一的软删除策略
- ✅ 统一的审计字段
- ❌ 不是特定业务逻辑

---

## 决策树：模型应该放在哪里

```
这个实体会被多个端使用吗？
│
├─ 是 → 是核心业务实体吗？
│  │
│  ├─ 是 → shared/domain/
│  │  └─ 例如：user.go, product.go, order.go
│  │
│  └─ 否 → 是通用基础结构吗？
│     │
│     ├─ 是 → shared/domain/common/
│     │  └─ 例如：base_entity.go, soft_delete.go
│     │
│     └─ 否 → 重新评估是否真的需要共享
│
└─ 否 → 功能复杂度如何？
   │
   ├─ 简单（< 1000行）
   │  └─→ apps/*/features/*/models.go
   │     └─ 例如：admin_log.go（在 logs/models.go 中）
   │
   └─ 复杂（> 1000行，使用 internal）
      └─→ apps/*/features/*/internal/entity/
         └─ 例如：report_entity.go
```

---

## 表模型拆分策略

### 按业务领域拆分

**原则：** 一个业务领域一个模型文件

```
shared/domain/
├── user.go              # 用户领域
├── product.go           # 商品领域
├── category.go          # 分类领域
├── order.go             # 订单领域
└── order_item.go        # 订单明细（属于订单领域）
```

**反例（不推荐）：**
```
shared/domain/
└── models.go            # ❌ 所有模型都在一个文件（难以维护）
```

---

### 按聚合根拆分

**什么是聚合根：** 一组紧密关联的实体的根实体

**示例：订单聚合**

```go
// shared/domain/order.go
package domain

// Order - 聚合根
type Order struct {
    ID         int64
    UserID     int64
    TotalPrice float64
    Status     OrderStatus
    Items      []OrderItem // 值对象，依附于 Order
}

// OrderItem - 值对象，属于 Order 聚合
type OrderItem struct {
    ID        int64
    OrderID   int64
    ProductID int64
    Quantity  int
    Price     float64
}
```

**示例：商品聚合**

```go
// shared/domain/product.go
package domain

// Product - 聚合根
type Product struct {
    ID       int64
    Name     string
    Price    float64
    Variants []ProductVariant // 值对象
}

// ProductVariant - 值对象，属于 Product 聚合
type ProductVariant struct {
    ID        int64
    ProductID int64
    SKU       string
    Color     string
    Size      string
    Stock     int
}
```

---

### 关联关系处理

**一对多关系：**

```go
// shared/domain/user.go
package domain

type User struct {
    ID    int64
    Name  string
    Email string
    // ❌ 不包含 Orders 集合（避免循环依赖和性能问题）
}

// shared/domain/order.go
package domain

type Order struct {
    ID     int64
    UserID int64  // ✅ 只保存外键
    Status string
}
```

**查询时按需加载：**

```go
// shared/repositories/order_repository.go
func (r *OrderRepository) GetByUserID(userID int64) ([]*domain.Order, error) {
    var orders []*domain.Order
    err := r.db.Where("user_id = ?", userID).Find(&orders).Error
    return orders, err
}
```

**多对多关系：**

```go
// shared/domain/product.go
package domain

type Product struct {
    ID   int64
    Name string
    // ❌ 不直接包含 Tags 集合
}

// shared/domain/tag.go
package domain

type Tag struct {
    ID   int64
    Name string
    // ❌ 不直接包含 Products 集合
}

// shared/domain/product_tag.go - 关联表
package domain

type ProductTag struct {
    ID        int64
    ProductID int64 `gorm:"uniqueIndex:idx_product_tag"`
    TagID     int64 `gorm:"uniqueIndex:idx_product_tag"`
}

func (ProductTag) TableName() string {
    return "product_tags"
}
```

---

## 模型字段设计规范

### 必需字段

所有 Domain 模型都应该包含：

```go
type Entity struct {
    ID        int64     `gorm:"primaryKey"`           // ✅ 主键
    CreatedAt time.Time `gorm:"autoCreateTime"`       // ✅ 创建时间
    UpdatedAt time.Time `gorm:"autoUpdateTime"`       // ✅ 更新时间
}
```

**推荐：** 使用 `common.BaseEntity` 嵌入

---

### 可选字段

根据业务需求选择：

```go
type Entity struct {
    common.BaseEntity
    
    // 软删除
    IsDeleted bool       `gorm:"default:false;index"`
    DeletedAt *time.Time `gorm:"index"`
    
    // 审计字段
    CreatedBy int64 `gorm:"not null"`
    UpdatedBy int64 `gorm:"not null"`
}
```

**推荐：** 使用 `common.SoftDelete` 和 `common.AuditFields` 嵌入

---

### 字段命名规范

- **使用 PascalCase**（Go 约定）
  ```go
  Name      string  // ✅
  user_name string  // ❌
  ```

- **布尔字段用 Is 前缀**
  ```go
  IsDeleted bool    // ✅
  IsActive  bool    // ✅
  Deleted   bool    // ❌
  ```

- **时间字段用 At 后缀**
  ```go
  CreatedAt time.Time  // ✅
  UpdatedAt time.Time  // ✅
  CreateTime time.Time // ❌
  ```

- **外键字段用 ID 后缀**
  ```go
  UserID    int64  // ✅
  ProductID int64  // ✅
  User      int64  // ❌
  ```

---

### GORM 标签规范

```go
type User struct {
    ID    int64  `gorm:"primaryKey"`                    // 主键
    Name  string `gorm:"size:100;not null"`            // 长度限制，非空
    Email string `gorm:"size:255;uniqueIndex;not null"` // 唯一索引
    Age   int    `gorm:"default:0"`                     // 默认值
    Role  string `gorm:"size:20;not null;default:'user';index"` // 组合
}
```

**常用标签：**
- `primaryKey` - 主键
- `size:255` - 字段长度
- `not null` - 非空
- `default:'value'` - 默认值
- `uniqueIndex` - 唯一索引
- `index` - 普通索引
- `type:text` - 指定数据库类型
- `autoCreateTime` - 自动设置创建时间
- `autoUpdateTime` - 自动设置更新时间

---

### JSON 标签规范

```go
type User struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    Password  string `json:"-"`           // ✅ 不序列化敏感字段
    IsDeleted bool   `json:"-"`           // ✅ 内部字段不暴露
}
```

**注意：** Domain 模型可以不加 JSON 标签，如果需要序列化应该通过 DTO 转换。

---

## 最佳实践

### 1. 避免循环依赖

**❌ 错误示例：**
```go
// user.go
type User struct {
    ID     int64
    Orders []Order  // ❌ 引用 Order
}

// order.go
type Order struct {
    ID   int64
    User User       // ❌ 引用 User - 循环依赖！
}
```

**✅ 正确示例：**
```go
// user.go
type User struct {
    ID int64
    // ✅ 不包含 Orders 集合
}

// order.go
type Order struct {
    ID     int64
    UserID int64    // ✅ 只保存外键
}
```

---

### 2. 敏感字段处理

```go
type User struct {
    ID       int64  `gorm:"primaryKey" json:"id"`
    Name     string `gorm:"size:100;not null" json:"name"`
    Password string `gorm:"size:255;not null" json:"-"`  // ✅ json:"-"
    Salt     string `gorm:"size:50;not null" json:"-"`   // ✅ 不暴露
}
```

**原则：**
- ✅ Password, Salt 等敏感字段使用 `json:"-"`
- ✅ 通过 DTO 转换，不直接返回 Domain 模型
- ❌ 不要在 API 响应中直接使用 Domain 模型

---

### 3. 统一软删除

```go
// 使用 common.SoftDelete
type Product struct {
    common.BaseEntity
    common.SoftDelete  // ✅ 统一的软删除
    
    Name  string
    Price float64
}

// 业务代码
func (s *ProductService) Delete(id int64) error {
    product, _ := s.repo.GetByID(id)
    product.Delete()  // ✅ 调用软删除方法
    return s.repo.Update(product)
}
```

**好处：**
- ✅ 行为一致
- ✅ 可恢复数据
- ✅ 审计追踪

---

### 4. 时间戳自动管理

```go
type Entity struct {
    ID        int64     `gorm:"primaryKey"`
    CreatedAt time.Time `gorm:"autoCreateTime"`  // ✅ 自动设置
    UpdatedAt time.Time `gorm:"autoUpdateTime"`  // ✅ 自动更新
}
```

**或使用嵌入：**
```go
type Product struct {
    common.BaseEntity  // ✅ 包含自动时间戳
    Name  string
    Price float64
}
```

---

### 5. 不在模型中写业务逻辑

**❌ 不推荐：**
```go
type User struct {
    ID    int64
    Email string
}

func (u *User) SendEmail(content string) error {
    // ❌ 业务逻辑不应在模型中
    return emailService.Send(u.Email, content)
}

func (u *User) CreateOrder(items []OrderItem) (*Order, error) {
    // ❌ 业务逻辑不应在模型中
    return orderService.Create(u.ID, items)
}
```

**✅ 推荐：**
```go
// Domain 模型只包含数据结构和简单的辅助方法
type User struct {
    ID    int64
    Email string
}

func (u *User) IsEmailValid() bool {
    // ✅ 简单的验证方法可以
    return strings.Contains(u.Email, "@")
}

// 业务逻辑放在 Service 中
type UserService struct {
    userRepo  IUserRepository
    emailSvc  IEmailService
}

func (s *UserService) SendWelcomeEmail(userID int64) error {
    user, _ := s.userRepo.GetByID(userID)
    return s.emailSvc.SendWelcome(user.Email)
}
```

---

### 6. 合理使用嵌入

**✅ 推荐嵌入：**
```go
type Product struct {
    common.BaseEntity   // ✅ 基础字段
    common.SoftDelete   // ✅ 通用行为
    
    Name  string
    Price float64
}
```

**❌ 避免过度嵌入：**
```go
type Product struct {
    common.BaseEntity
    common.SoftDelete
    common.AuditFields
    common.Versioning
    common.I18n
    common.Metadata
    // ❌ 嵌入太多，结构复杂
}
```

**原则：** 只嵌入确实需要的通用结构

---

## 总结

**Domain 模型组织的核心原则：**

1. **职责单一** - Domain 模型只负责数据结构和持久化
2. **合理共享** - 跨端使用才放 shared/domain/
3. **私有优先** - 单端使用优先放功能内
4. **避免循环** - 只用外键，不直接引用关联对象
5. **统一基础** - 使用 common 包统一基础结构
6. **安全第一** - 敏感字段不暴露
7. **简单辅助** - 只包含简单验证方法，不写业务逻辑

---

**返回 [主文档](../ORGANIZATION_GUIDE.md)**
