# internal 目录使用指南

## 目录

- [什么是 internal 目录](#什么是-internal-目录)
- [为什么使用 internal](#为什么使用-internal)
- [不同层次的 internal](#不同层次的-internal)
- [决策树：是否使用 internal](#决策树是否使用-internal)
- [internal 目录结构规范](#internal-目录结构规范)
- [最佳实践](#最佳实践)
- [常见错误](#常见错误)

---

## 什么是 internal 目录

Go 语言的特殊目录，`internal/` 包只能被其父目录及子目录导入，无法被外部包访问。

**Go 语言规则：**
```
vertical_slice_demo/
├── apps/
│   └── admin/
│       └── features/
│           └── reports/
│               ├── internal/
│               │   └── entity/
│               │       └── report.go      # ✅ 可被 reports/ 下的代码导入
│               │                          # ❌ 不能被 reports/ 外的代码导入
│               └── handler.go
```

**访问规则：**
- ✅ `apps/admin/features/reports/handler.go` 可以导入 `reports/internal/entity`
- ✅ `apps/admin/features/reports/generate_report.go` 可以导入 `reports/internal/entity`
- ❌ `apps/admin/features/users/handler.go` 不能导入 `reports/internal/entity`
- ❌ `apps/api/features/xxx/` 不能导入 `admin/features/reports/internal/`

---

## 为什么使用 internal

### 好处

1. **封装内部实现**
   - 防止外部代码直接访问内部细节
   - 明确对外接口边界

2. **提高代码安全性**
   - 内部实现可以自由重构
   - 不用担心破坏外部依赖

3. **防止误用**
   - 避免其他模块绕过对外接口直接访问内部

4. **清晰的架构边界**
   - 一眼看出哪些是对外接口（外层）
   - 哪些是内部实现（internal）

### 示例：没有 internal 的问题

```go
// features/reports/entity/report_entity.go
package entity

// 外部可以直接访问
type ReportEntity struct {...}

// features/users/handler.go
import "vertical_slice_demo/apps/admin/features/reports/entity"

// ❌ 其他功能可以直接使用 reports 的内部实体
func (h *UserHandler) SomeMethod() {
    report := &entity.ReportEntity{}  // 不应该这样做
}
```

### 示例：使用 internal 的好处

```go
// features/reports/internal/entity/report_entity.go
package entity

// 被 internal 保护
type ReportEntity struct {...}

// features/users/handler.go
import "vertical_slice_demo/apps/admin/features/reports/internal/entity"
// ❌ 编译错误：use of internal package not allowed
```

---

## 不同层次的 internal

### 1. 功能级 internal（推荐）

**位置：** `apps/*/features/*/internal/`

**适用场景：** 功能内部的实现细节

**目录结构：**
```
apps/admin/features/reports/
├── internal/                        # 功能内部实现
│   ├── entity/                      # 内部实体
│   │   ├── report_entity.go
│   │   └── template_entity.go
│   ├── data/                        # 内部数据访问
│   │   └── report_store.go
│   └── business/                    # 内部业务逻辑
│       └── report_generator.go
│
├── models.go                        # 对外 DTO
├── generate_report.go               # 对外 Handler
└── controller.go
```

**何时使用：**
- ✅ 复杂功能需要内部分层（模式3）
- ✅ 内部实体不希望被其他功能访问
- ✅ 内部 Store/Business 只在本功能内使用

**示例：**

```go
// internal/entity/report_entity.go
package entity

// ReportEntity 内部报表实体（只在 reports 功能内使用）
type ReportEntity struct {
    ID         int64
    Name       string
    Type       ReportType
    Config     string
    CreatedAt  time.Time
}
```

```go
// internal/data/report_store.go
package data

import "vertical_slice_demo/apps/admin/features/reports/internal/entity"

// ReportStore 内部数据访问（不暴露给外部）
type ReportStore struct {
    reports map[int64]*entity.ReportEntity
    mu      sync.RWMutex
}

func (s *ReportStore) Create(report *entity.ReportEntity) error {
    // ...
}
```

```go
// internal/business/report_generator.go
package business

import (
    "vertical_slice_demo/apps/admin/features/reports/internal/entity"
    "vertical_slice_demo/apps/admin/features/reports/internal/data"
)

// ReportGenerator 报表生成器（内部业务逻辑）
type ReportGenerator struct {
    store *data.ReportStore
}

func (g *ReportGenerator) Generate(config entity.ReportConfig) (*entity.ReportEntity, error) {
    // 复杂的生成逻辑
}
```

```go
// generate_report.go (对外 Handler)
package reports

import (
    "vertical_slice_demo/apps/admin/features/reports/internal/business"
    "vertical_slice_demo/apps/admin/features/reports/internal/entity"
)

type GenerateReportHandler struct {
    generator *business.ReportGenerator  // ✅ 可以使用 internal 包
}

func (h *GenerateReportHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req GenerateReportRequest  // 对外 DTO
    c.MustBindJSON(&req)
    
    // 转换 DTO 到内部实体
    config := entity.ReportConfig{
        Name: req.Name,
        Type: entity.ReportType(req.ReportType),
    }
    
    // 调用内部业务逻辑
    report, _ := h.generator.Generate(config)
    
    // 转换内部实体到 DTO
    response := toReportResponse(report)
    
    return c.Ok(response)
}
```

---

### 2. 应用级 internal

**位置：** `apps/*/internal/`

**适用场景：** 应用内多个功能共享，但不跨应用

**目录结构：**
```
apps/admin/
├── internal/                        # Admin 应用内部共享
│   ├── models/                      # 应用内共享的内部实体
│   │   └── admin_session.go
│   ├── stores/                      # 应用内共享的 Store
│   │   └── cache_store.go
│   ├── utils/                       # 应用内工具函数
│   │   └── permission_checker.go
│   └── middleware/                  # 应用内中间件
│       └── auth_middleware.go
│
└── features/
    ├── users/
    ├── products/
    └── orders/
```

**何时使用：**
- ✅ Admin 端多个功能需要共享，但 API 端不需要
- ✅ 应用内部的工具函数
- ✅ 应用特定的缓存/会话管理
- ✅ 应用特定的中间件

**示例：**

```go
// apps/admin/internal/utils/permission_checker.go
package utils

// PermissionChecker Admin 端权限检查（内部工具）
type PermissionChecker struct {
    adminRepo IAdminRepository
}

func NewPermissionChecker(adminRepo IAdminRepository) *PermissionChecker {
    return &PermissionChecker{adminRepo: adminRepo}
}

func (p *PermissionChecker) CheckPermission(adminID int64, resource string, action string) bool {
    // Admin 端特定的权限检查逻辑
    admin, _ := p.adminRepo.GetByID(adminID)
    return admin.HasPermission(resource, action)
}
```

```go
// apps/admin/internal/middleware/auth_middleware.go
package middleware

import "vertical_slice_demo/apps/admin/internal/utils"

// AdminAuthMiddleware Admin 端认证中间件
type AdminAuthMiddleware struct {
    checker *utils.PermissionChecker  // ✅ 使用应用内部工具
}

func (m *AdminAuthMiddleware) Handle(c *web.HttpContext) {
    adminID := c.GetUserID()
    if !m.checker.CheckPermission(adminID, c.Path(), "access") {
        c.AbortWithStatus(403)
        return
    }
    c.Next()
}
```

```go
// apps/admin/features/users/handler.go
package users

import "vertical_slice_demo/apps/admin/internal/utils"

type CreateUserHandler struct {
    userRepo  IUserRepository
    checker   *utils.PermissionChecker  // ✅ 多个功能共享
}
```

---

### 3. 全局 shared（不使用 internal）

**位置：** `shared/domain/`, `shared/repositories/`, `shared/services/`

**适用场景：** 明确需要跨应用共享

**何时使用：**
- ✅ 多个应用端共享的核心实体（User, Product）
- ✅ 跨端的数据访问层
- ✅ 跨端的业务服务

**说明：** shared 目录下通常不使用 internal，因为其目的就是共享

```
shared/
├── domain/                          # 共享领域模型（不用 internal）
│   ├── user.go
│   └── product.go
│
├── repositories/                    # 共享仓储（不用 internal）
│   └── user_repository.go
│
└── services/                        # 共享服务（不用 internal）
    └── order_service.go
```

**为什么 shared 不用 internal：**
- shared 的目的就是让多个应用使用
- 如果需要隐藏细节，说明不应该放在 shared

---

## 决策树：是否使用 internal

```
这段代码需要被外部访问吗？
│
├─ 否（只在功能内部使用）
│  │
│  ├─ 是 DTO（Request/Response）？
│  │  └─→ 不要放 internal，放外层
│  │     └─ models.go
│  │
│  ├─ 是对外 Handler？
│  │  └─→ 不要放 internal，放外层
│  │     └─ create_xxx.go
│  │
│  └─ 是内部实现（Entity/Store/Business）？
│     └─→ features/*/internal/
│        ├─ internal/entity/
│        ├─ internal/data/
│        └─ internal/business/
│
├─ 需要在应用内共享（但不跨应用）
│  │
│  ├─ 是工具/辅助代码？
│  │  └─→ apps/*/internal/utils/
│  │
│  ├─ 是应用特定逻辑？
│  │  └─→ apps/*/internal/
│  │
│  └─ 是核心业务实体？
│     └─→ 重新考虑，可能应该在 shared/
│
└─ 需要跨应用共享
   └─→ shared/（不使用 internal）
```

---

## internal 目录结构规范

### 功能级 internal 推荐结构

```
features/xxx/internal/
├── entity/              # 内部实体/领域对象
│   ├── xxx_entity.go
│   └── config.go
│
├── data/                # 数据访问层
│   ├── xxx_store.go
│   └── cache.go
│
└── business/            # 业务逻辑层
    ├── xxx_service.go
    ├── calculator.go
    └── validator.go
```

**entity/ 目录：**
- 内部实体定义
- 与 Domain 模型不同，这些是功能特定的
- 不对外暴露

**data/ 目录：**
- 内部数据访问
- 可能是内存存储、缓存等
- 不对外暴露

**business/ 目录：**
- 内部业务逻辑
- 复杂计算、验证等
- 可以在多个 Handler 间复用
- 不对外暴露

---

### 应用级 internal 推荐结构

```
apps/*/internal/
├── models/              # 应用内共享实体
│   └── session.go
│
├── stores/              # 应用内共享 Store
│   └── cache_store.go
│
├── utils/               # 应用内工具
│   ├── permission_checker.go
│   └── id_generator.go
│
└── middleware/          # 应用内中间件
    ├── auth.go
    └── logging.go
```

---

## 最佳实践

### 1. 优先使用功能级 internal

```go
// ✅ 推荐：大部分情况下使用功能级 internal
apps/admin/features/reports/
└── internal/
    ├── entity/
    ├── data/
    └── business/

// ❌ 避免：不要过早使用应用级 internal
apps/admin/internal/
└── report_stuff.go  // 只有一个功能用，应该在功能内
```

**原则：** 默认用功能级 internal，确实需要应用内共享时才用应用级

---

### 2. 谨慎使用应用级 internal

```go
// ✅ 好的使用：真正的应用级共享
apps/admin/internal/utils/
└── permission_checker.go  // Admin 端多个功能都要检查权限

// ❌ 不好的使用：应该放功能内或 shared
apps/admin/internal/
└── user_validator.go  // 如果是核心业务，应该在 shared
```

**提醒：** 应用级 internal 容易变成"另一个 shared"，要谨慎使用

---

### 3. internal 中可以有 internal

```go
features/reports/internal/data/
├── internal/            # ✅ data 层的内部实现
│   └── cache.go
└── report_store.go      # data 层的对外接口
```

**使用场景：** 当 internal 内部也需要进一步封装时

---

### 4. 对外接口放在 internal 外

```go
// ✅ 正确：对外接口在外层
features/reports/
├── internal/
│   ├── entity/
│   ├── data/
│   └── business/
│
├── models.go            # ✅ DTO 在外面
├── generate_report.go   # ✅ Handler 在外面
└── controller.go        # ✅ Controller 在外面
```

**原则：**
- models.go（DTO）必须在外面
- handler.go（对外接口）必须在外面
- controller.go（路由）必须在外面

---

### 5. 转换函数的位置

**Entity → DTO 转换：**

```go
// 方式1：在 handler 中
// features/reports/generate_report.go
func (h *GenerateReportHandler) Handle(c *web.HttpContext) web.IActionResult {
    report, _ := h.generator.Generate(config)
    
    // ✅ 在 handler 中转换
    response := ReportResponse{
        ID:   report.ID,
        Name: report.Name,
    }
    
    return c.Ok(response)
}

// 方式2：在 models.go 中
// features/reports/models.go
func toReportResponse(report *entity.ReportEntity) *ReportResponse {
    return &ReportResponse{
        ID:   report.ID,
        Name: report.Name,
    }
}
```

**Domain → Entity 转换：**

```go
// ✅ 在 internal/entity/ 中
// features/reports/internal/entity/converters.go
func FromDomainOrder(order *domain.Order) *OrderEntity {
    return &OrderEntity{
        ID:         order.ID,
        TotalPrice: order.TotalPrice,
    }
}
```

---

## 常见错误

### ❌ 错误1：所有东西都放 internal

```go
// ❌ 错误示例
features/reports/
└── internal/
    ├── models.go        # ❌ DTO 应该在外面
    ├── handler.go       # ❌ Handler 是对外接口
    └── entity.go
```

**正确做法：**

```go
// ✅ 正确示例
features/reports/
├── internal/            # 只有内部实现
│   └── entity/
│       └── report_entity.go
├── models.go            # ✅ DTO 在外面
└── generate_report.go   # ✅ Handler 在外面
```

---

### ❌ 错误2：不该用 internal 却用了

```go
// ❌ 错误示例：简单功能不需要 internal
features/simple_crud/
└── internal/
    ├── handler.go       # ❌ 简单功能直接暴露即可
    └── models.go
```

**正确做法：**

```go
// ✅ 正确示例：简单功能不用 internal
features/simple_crud/
├── handler.go           # ✅ 直接暴露
└── service_extensions.go
```

**原则：** 只有复杂功能（模式3）才需要 internal

---

### ❌ 错误3：应用级 internal 变成垃圾堆

```go
// ❌ 错误示例：什么都往里放
apps/admin/internal/
├── user_stuff.go        # 应该在 features/users/
├── product_util.go      # 应该在 features/products/
├── random_helper.go     # 不知道干什么的
└── legacy_code.go       # 遗留代码
```

**正确做法：**

```go
// ✅ 正确示例：只放真正需要应用内共享的
apps/admin/internal/
├── utils/
│   └── permission_checker.go  # ✅ 多个功能共享
└── middleware/
    └── auth_middleware.go     # ✅ 应用级中间件
```

---

### ❌ 错误4：过度嵌套 internal

```go
// ❌ 错误示例：过度嵌套
features/reports/
└── internal/
    └── data/
        └── internal/
            └── cache/
                └── internal/
                    └── storage.go  # 太深了
```

**正确做法：**

```go
// ✅ 正确示例：合理的层次
features/reports/
└── internal/
    ├── entity/
    ├── data/
    │   └── internal/       # 最多一层嵌套
    │       └── cache.go
    └── business/
```

---

## 完整示例：带 internal 的报表功能

```
apps/admin/features/reports/
├── internal/                               # 内部实现
│   ├── entity/                             # 内部实体
│   │   ├── report_entity.go
│   │   ├── report_config.go
│   │   └── chart_config.go
│   │
│   ├── data/                               # 数据访问
│   │   ├── report_store.go
│   │   └── internal/                       # data 层的内部实现
│   │       └── cache.go
│   │
│   └── business/                           # 业务逻辑
│       ├── report_generator.go
│       ├── data_aggregator.go
│       └── chart_builder.go
│
├── models.go                               # 对外 DTO
├── generate_report.go                      # 对外 Handler
├── export_report.go
├── list_reports.go
├── controller.go
└── service_extensions.go
```

**访问规则：**
- ✅ `generate_report.go` 可以导入 `internal/business`
- ✅ `internal/business` 可以导入 `internal/data`
- ✅ `internal/data` 可以导入 `internal/entity`
- ❌ `apps/admin/features/users/` 不能导入 `reports/internal/`
- ❌ `apps/api/features/` 不能导入 `admin/features/reports/internal/`

---

## 总结

**internal 使用的核心原则：**

1. **封装内部实现** - 使用 internal 保护实现细节
2. **功能级优先** - 大部分情况用 features/*/internal/
3. **对外接口在外层** - DTO、Handler、Controller 不放 internal
4. **谨慎应用级** - 确实需要应用内共享才用 apps/*/internal/
5. **shared 不用 internal** - shared 的目的就是共享
6. **合理层次** - 避免过度嵌套
7. **简单功能不需要** - 只有复杂功能（模式3）才用 internal

---

**返回 [主文档](../ORGANIZATION_GUIDE.md)**
