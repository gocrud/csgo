# 代码组织指南

## 🎯 核心原则

**不要过早抽象！** 只有在确定需要跨端共享时，才将代码提取到 `shared/`。

## 📊 决策树

```
你的功能需要被多个端使用吗？
├─ 是 → 放在 shared/
│  ├─ 数据模型 → shared/domain/
│  ├─ 数据访问 → shared/repositories/
│  └─ 业务服务 → shared/services/
│
└─ 否 → 放在对应端的 features/
   │
   ├─ 逻辑简单？（单表 CRUD，< 200 行）
   │  └─ 是 → 单文件实现
   │     └─ features/xxx/
   │         ├── handler.go           # 所有逻辑
   │         └── service_extensions.go
   │
   ├─ 中等复杂？（多个操作，200-1000 行）
   │  └─ 是 → 按操作拆分
   │     └─ features/xxx/
   │         ├── models.go            # 数据模型
   │         ├── create_xxx.go        # 创建操作
   │         ├── list_xxx.go          # 列表操作
   │         ├── update_xxx.go        # 更新操作
   │         ├── controller.go        # 路由
   │         └── service_extensions.go
   │
   └─ 逻辑复杂？（复杂业务，> 1000 行）
      └─ 是 → 内部分层
          └─ features/xxx/
              ├── models/              # 数据模型
              ├── data/                # 内部数据访问
              ├── business/            # 内部业务逻辑
              ├── create_xxx.go        # Handler
              ├── list_xxx.go
              ├── controller.go
              └── service_extensions.go
```

## 📁 三种组织模式详解

### 模式 1️⃣：单文件实现（Simple）

**适用场景：**
- 简单的 CRUD 操作
- 代码量 < 200 行
- 业务逻辑简单

**示例：分类管理**

```
apps/admin/features/categories/
├── handler.go               # 所有功能
└── service_extensions.go    # DI 注册
```

**handler.go 内容：**

```go
package categories

import (
    "sync"
    "github.com/gocrud/csgo/web"
)

// ===== 数据模型 =====
type Category struct {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

// ===== 数据访问 =====
type CategoryHandler struct {
    categories map[int64]*Category
    mu         sync.RWMutex
    nextID     int64
}

func NewCategoryHandler() *CategoryHandler {
    return &CategoryHandler{
        categories: make(map[int64]*Category),
        nextID:     1,
    }
}

// ===== HTTP Handlers =====
func (h *CategoryHandler) Create(c *web.HttpContext) web.IActionResult {
    var req Category
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    h.mu.Lock()
    defer h.mu.Unlock()
    
    req.ID = h.nextID
    h.nextID++
    h.categories[req.ID] = &req
    
    return c.Created(req)
}

func (h *CategoryHandler) List(c *web.HttpContext) web.IActionResult {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    list := make([]*Category, 0, len(h.categories))
    for _, cat := range h.categories {
        list = append(list, cat)
    }
    
    return c.Ok(list)
}

func (h *CategoryHandler) MapRoutes(app *web.WebApplication) {
    g := app.MapGroup("/api/admin/categories")
    g.MapPost("", h.Create)
    g.MapGet("", h.List)
}
```

**优点：**
- ✅ 简单直接，所有代码在一起
- ✅ 改动方便，不需要跨文件
- ✅ 容易理解和维护

**缺点：**
- ❌ 文件变大后不好维护
- ❌ 多人协作容易冲突

---

### 模式 2️⃣：按操作拆分（Recommended）

**适用场景：**
- 多个操作
- 代码量 200-1000 行
- 需要团队协作

**示例：标签管理**

```
apps/admin/features/tags/
├── models.go                # 数据模型
├── store.go                 # 内部数据访问
├── create_tag.go            # 创建操作
├── list_tags.go             # 列表操作
├── update_tag.go            # 更新操作
├── delete_tag.go            # 删除操作
├── controller.go            # 路由映射
└── service_extensions.go    # DI 注册
```

**models.go:**

```go
package tags

type Tag struct {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Color       string `json:"color"`
    Description string `json:"description"`
}

type CreateTagRequest struct {
    Name        string `json:"name" binding:"required"`
    Color       string `json:"color" binding:"required"`
    Description string `json:"description"`
}

type UpdateTagRequest struct {
    Name        string `json:"name" binding:"required"`
    Color       string `json:"color" binding:"required"`
    Description string `json:"description"`
}
```

**store.go:**

```go
package tags

import "sync"

// TagStore 内部数据访问层（不暴露到外部）
type TagStore struct {
    tags   map[int64]*Tag
    mu     sync.RWMutex
    nextID int64
}

func NewTagStore() *TagStore {
    return &TagStore{
        tags:   make(map[int64]*Tag),
        nextID: 1,
    }
}

func (s *TagStore) Create(tag *Tag) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    tag.ID = s.nextID
    s.nextID++
    s.tags[tag.ID] = tag
    return nil
}

func (s *TagStore) GetByID(id int64) (*Tag, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    tag, exists := s.tags[id]
    if !exists {
        return nil, errors.New("tag not found")
    }
    return tag, nil
}

func (s *TagStore) List() []*Tag {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    list := make([]*Tag, 0, len(s.tags))
    for _, tag := range s.tags {
        list = append(list, tag)
    }
    return list
}

// ... 其他方法
```

**create_tag.go:**

```go
package tags

import "github.com/gocrud/csgo/web"

type CreateTagHandler struct {
    store *TagStore
}

func NewCreateTagHandler(store *TagStore) *CreateTagHandler {
    return &CreateTagHandler{store: store}
}

func (h *CreateTagHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req CreateTagRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    // 业务验证
    if req.Name == "" {
        return c.BadRequest("标签名称不能为空")
    }
    
    // 创建标签
    tag := &Tag{
        Name:        req.Name,
        Color:       req.Color,
        Description: req.Description,
    }
    
    if err := h.store.Create(tag); err != nil {
        return c.InternalError("创建失败")
    }
    
    return c.Created(tag)
}
```

**list_tags.go:**

```go
package tags

import "github.com/gocrud/csgo/web"

type ListTagsHandler struct {
    store *TagStore
}

func NewListTagsHandler(store *TagStore) *ListTagsHandler {
    return &ListTagsHandler{store: store}
}

func (h *ListTagsHandler) Handle(c *web.HttpContext) web.IActionResult {
    tags := h.store.List()
    return c.Ok(tags)
}
```

**controller.go:**

```go
package tags

import "github.com/gocrud/csgo/web"

type TagController struct {
    createHandler *CreateTagHandler
    listHandler   *ListTagsHandler
    updateHandler *UpdateTagHandler
    deleteHandler *DeleteTagHandler
}

func NewTagController(
    createHandler *CreateTagHandler,
    listHandler *ListTagsHandler,
    updateHandler *UpdateTagHandler,
    deleteHandler *DeleteTagHandler,
) *TagController {
    return &TagController{
        createHandler: createHandler,
        listHandler:   listHandler,
        updateHandler: updateHandler,
        deleteHandler: deleteHandler,
    }
}

func (ctrl *TagController) MapRoutes(app *web.WebApplication) {
    tags := app.MapGroup("/api/admin/tags")
    tags.MapPost("", ctrl.createHandler.Handle)
    tags.MapGet("", ctrl.listHandler.Handle)
    tags.MapPut("/:id", ctrl.updateHandler.Handle)
    tags.MapDelete("/:id", ctrl.deleteHandler.Handle)
}
```

**service_extensions.go:**

```go
package tags

import (
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

func AddTagFeature(services di.IServiceCollection) {
    // 注册内部存储（Singleton）
    services.AddSingleton(NewTagStore)
    
    // 注册 Handlers
    services.AddSingleton(NewCreateTagHandler)
    services.AddSingleton(NewListTagsHandler)
    services.AddSingleton(NewUpdateTagHandler)
    services.AddSingleton(NewDeleteTagHandler)
    
    // 注册 Controller
    web.AddController(services, NewTagController)
}
```

**优点：**
- ✅ 结构清晰，职责分明
- ✅ 每个操作一个文件，容易定位
- ✅ 团队协作友好，不同人改不同文件
- ✅ 内部 Store 可以复用

**缺点：**
- ❌ 文件数量多一些

---

### 模式 3️⃣：内部分层（Complex）

**适用场景：**
- 复杂业务逻辑
- 代码量 > 1000 行
- 需要内部复用

**示例：报表系统**

```
apps/admin/features/reports/
├── models/                          # 数据模型层
│   ├── report.go
│   ├── report_config.go
│   └── report_template.go
│
├── data/                            # 数据访问层（内部）
│   ├── report_store.go
│   └── template_store.go
│
├── business/                        # 业务逻辑层（内部）
│   ├── report_generator.go         # 报表生成
│   ├── report_exporter.go          # 报表导出
│   ├── data_aggregator.go          # 数据聚合
│   └── chart_builder.go            # 图表构建
│
├── generate_report.go               # Handler：生成报表
├── export_report.go                 # Handler：导出报表
├── schedule_report.go               # Handler：定时报表
├── list_reports.go                  # Handler：报表列表
├── controller.go                    # 路由
└── service_extensions.go            # DI 注册
```

**models/report.go:**

```go
package models

type Report struct {
    ID         int64
    Name       string
    Type       ReportType
    DataSource DataSource
    Charts     []Chart
    CreatedAt  time.Time
}

type ReportType string

const (
    ReportTypeSales     ReportType = "sales"
    ReportTypeUser      ReportType = "user"
    ReportTypeInventory ReportType = "inventory"
)

type Chart struct {
    Type   ChartType
    Title  string
    Data   interface{}
}
```

**data/report_store.go:**

```go
package data

// ReportStore 内部数据访问（不暴露）
type ReportStore struct {
    orderRepo   repositories.IOrderRepository   // 使用共享仓储
    productRepo repositories.IProductRepository
    reports     map[int64]*models.Report
    mu          sync.RWMutex
}

func NewReportStore(
    orderRepo repositories.IOrderRepository,
    productRepo repositories.IProductRepository,
) *ReportStore {
    return &ReportStore{
        orderRepo:   orderRepo,
        productRepo: productRepo,
        reports:     make(map[int64]*models.Report),
    }
}

func (s *ReportStore) GetOrdersForPeriod(start, end time.Time) ([]*domain.Order, error) {
    // 从共享仓储获取数据
    return s.orderRepo.GetByPeriod(start, end)
}

func (s *ReportStore) SaveReport(report *models.Report) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.reports[report.ID] = report
    return nil
}
```

**business/report_generator.go:**

```go
package business

// ReportGenerator 报表生成器（内部业务逻辑）
type ReportGenerator struct {
    store      *data.ReportStore
    aggregator *DataAggregator
    chartBuilder *ChartBuilder
}

func NewReportGenerator(
    store *data.ReportStore,
    aggregator *DataAggregator,
    chartBuilder *ChartBuilder,
) *ReportGenerator {
    return &ReportGenerator{
        store:      store,
        aggregator: aggregator,
        chartBuilder: chartBuilder,
    }
}

func (g *ReportGenerator) Generate(config models.ReportConfig) (*models.Report, error) {
    // 1. 获取数据
    data, err := g.store.GetOrdersForPeriod(config.StartDate, config.EndDate)
    if err != nil {
        return nil, err
    }
    
    // 2. 聚合数据
    aggregated := g.aggregator.Aggregate(data, config.GroupBy)
    
    // 3. 构建图表
    charts := g.chartBuilder.BuildCharts(aggregated, config.ChartTypes)
    
    // 4. 生成报表
    report := &models.Report{
        Name:   config.Name,
        Type:   config.Type,
        Charts: charts,
    }
    
    // 5. 保存报表
    if err := g.store.SaveReport(report); err != nil {
        return nil, err
    }
    
    return report, nil
}
```

**business/data_aggregator.go:**

```go
package business

// DataAggregator 数据聚合器（内部逻辑）
type DataAggregator struct{}

func NewDataAggregator() *DataAggregator {
    return &DataAggregator{}
}

func (a *DataAggregator) Aggregate(orders []*domain.Order, groupBy string) map[string]interface{} {
    // 复杂的聚合逻辑
    result := make(map[string]interface{})
    
    switch groupBy {
    case "day":
        result = a.aggregateByDay(orders)
    case "month":
        result = a.aggregateByMonth(orders)
    case "product":
        result = a.aggregateByProduct(orders)
    }
    
    return result
}

func (a *DataAggregator) aggregateByDay(orders []*domain.Order) map[string]interface{} {
    // 按天聚合
    return nil
}
```

**generate_report.go (Handler):**

```go
package reports

type GenerateReportHandler struct {
    generator *business.ReportGenerator
}

func NewGenerateReportHandler(generator *business.ReportGenerator) *GenerateReportHandler {
    return &GenerateReportHandler{generator: generator}
}

func (h *GenerateReportHandler) Handle(c *web.HttpContext) web.IActionResult {
    var config models.ReportConfig
    if err := c.MustBindJSON(&config); err != nil {
        return err
    }
    
    // 调用业务层生成报表
    report, err := h.generator.Generate(config)
    if err != nil {
        return c.InternalError("生成报表失败")
    }
    
    return c.Ok(report)
}
```

**service_extensions.go:**

```go
package reports

func AddReportFeature(services di.IServiceCollection) {
    // 注册数据层
    services.AddSingleton(data.NewReportStore)
    services.AddSingleton(data.NewTemplateStore)
    
    // 注册业务层
    services.AddSingleton(business.NewReportGenerator)
    services.AddSingleton(business.NewReportExporter)
    services.AddSingleton(business.NewDataAggregator)
    services.AddSingleton(business.NewChartBuilder)
    
    // 注册 Handlers
    services.AddSingleton(NewGenerateReportHandler)
    services.AddSingleton(NewExportReportHandler)
    services.AddSingleton(NewScheduleReportHandler)
    services.AddSingleton(NewListReportsHandler)
    
    // 注册 Controller
    web.AddController(services, NewReportController)
}
```

**优点：**
- ✅ 适合复杂业务
- ✅ 内部逻辑可复用
- ✅ 职责分层清晰
- ✅ 容易测试每一层

**缺点：**
- ❌ 结构复杂
- ❌ 需要更多的规划

---

## 📦 数据传输模型（DTO）组织

### 核心理念

DTO（Data Transfer Object）是 API 层和业务层之间传输数据的对象，包括 Request、Response、ListItem 等。

**关键问题：** DTO 应该放在哪里？如何组织？

### 方案 1️⃣：与操作放在一起（推荐垂直切片）

**适用场景：** 每个操作的 DTO 都不同，追求功能完全独立

**目录结构：**

```
features/users/
├── create_user.go              # ✅ Request/Response 和逻辑在一起
├── list_users.go
├── update_user.go
├── controller.go
└── service_extensions.go
```

**示例代码：**

```go
// features/users/create_user.go
package users

import "github.com/gocrud/csgo/web"

// ===== DTO 定义 =====
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Role     string `json:"role" binding:"required,oneof=admin user"`
}

type CreateUserResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}

// ===== Handler =====
type CreateUserHandler struct {
    userRepo repositories.IUserRepository
}

func (h *CreateUserHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    // 业务逻辑...
    user := &domain.User{...}
    h.userRepo.Create(user)
    
    // 返回响应
    response := &CreateUserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        Role:  user.Role,
    }
    
    return c.Created(response)
}
```

```go
// features/users/list_users.go
package users

// ===== DTO 定义 =====
type UserListItem struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}

type ListUsersResponse struct {
    Users []UserListItem `json:"users"`
    Total int            `json:"total"`
}

// ===== Handler =====
type ListUsersHandler struct {
    userRepo repositories.IUserRepository
}

func (h *ListUsersHandler) Handle(c *web.HttpContext) web.IActionResult {
    users, _ := h.userRepo.List(offset, limit)
    
    items := make([]UserListItem, len(users))
    for i, user := range users {
        items[i] = UserListItem{
            ID:    user.ID,
            Name:  user.Name,
            Email: user.Email,
            Role:  user.Role,
        }
    }
    
    response := &ListUsersResponse{
        Users: items,
        Total: len(items),
    }
    
    return c.Ok(response)
}
```

**优点：**
- ✅ 功能完全内聚，改动范围最小
- ✅ 一个文件看到所有相关代码（DTO + 逻辑）
- ✅ 容易找到 DTO 定义
- ✅ 删除功能时 DTO 一起删除
- ✅ 符合垂直切片理念

**缺点：**
- ❌ 多个操作可能重复定义相似结构
- ❌ DTO 不能跨操作复用
- ❌ 文件可能偏长

**何时使用：**
- 每个操作的 DTO 都不同或差异大
- 追求功能完全独立
- 团队规模小（< 5 人）
- 代码审查重视功能完整性

---

### 方案 2️⃣：功能内共享 models.go（平衡方案 ⭐推荐）

**适用场景：** 功能内多个操作共享 DTO，需要一定程度的复用

**目录结构：**

```
features/orders/
├── models.go                   # ✅ 功能内共享的 DTO
├── create_order.go             # 使用 models 中的 DTO
├── list_orders.go
├── update_order_status.go
├── get_order_detail.go
├── controller.go
└── service_extensions.go
```

**示例代码：**

```go
// features/orders/models.go
package orders

import "time"

// ===== 共享的实体结构 =====
type Order struct {
    ID         int64       `json:"id"`
    UserID     int64       `json:"user_id"`
    TotalPrice float64     `json:"total_price"`
    Status     string      `json:"status"`
    Items      []OrderItem `json:"items"`
    CreatedAt  time.Time   `json:"created_at"`
}

type OrderItem struct {
    ProductID int64   `json:"product_id"`
    Quantity  int     `json:"quantity"`
    Price     float64 `json:"price"`
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

// ===== Response DTOs =====
type OrderResponse struct {
    ID         int64       `json:"id"`
    TotalPrice float64     `json:"total_price"`
    Status     string      `json:"status"`
    Items      []OrderItem `json:"items"`
    CreatedAt  time.Time   `json:"created_at"`
}

type OrderListItem struct {
    ID         int64     `json:"id"`
    TotalPrice float64   `json:"total_price"`
    Status     string    `json:"status"`
    ItemCount  int       `json:"item_count"`
    CreatedAt  time.Time `json:"created_at"`
}

type OrderDetailResponse struct {
    Order    OrderResponse  `json:"order"`
    User     UserInfo       `json:"user"`
    Products []ProductInfo  `json:"products"`
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
```

```go
// features/orders/create_order.go
package orders

type CreateOrderHandler struct {
    orderRepo   repositories.IOrderRepository
    productRepo repositories.IProductRepository
}

func (h *CreateOrderHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req CreateOrderRequest  // ✅ 使用 models.go 中的
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    // 业务逻辑...
    
    response := &OrderResponse{...}  // ✅ 使用 models.go 中的
    return c.Created(response)
}
```

```go
// features/orders/list_orders.go
package orders

type ListOrdersHandler struct {
    orderRepo repositories.IOrderRepository
}

func (h *ListOrdersHandler) Handle(c *web.HttpContext) web.IActionResult {
    orders, _ := h.orderRepo.GetByUserID(userID, offset, limit)
    
    // 转换为列表项
    items := make([]OrderListItem, len(orders))  // ✅ 使用 models.go 中的
    for i, order := range orders {
        items[i] = OrderListItem{
            ID:         order.ID,
            TotalPrice: order.TotalPrice,
            Status:     order.Status,
            ItemCount:  len(order.Items),
            CreatedAt:  order.CreatedAt,
        }
    }
    
    return c.Ok(items)
}
```

**优点：**
- ✅ 功能内 DTO 集中管理，避免重复
- ✅ 依然保持功能内聚
- ✅ 易于复用和维护
- ✅ 统一查看所有 DTO
- ✅ 平衡了复用和独立性

**缺点：**
- ❌ models.go 可能会变大（> 10 个 DTO 时）
- ❌ 需要判断哪些放 models，哪些放操作文件

**何时使用：**
- 功能内多个操作共享 DTO（推荐）
- 需要一定程度的复用
- 80% 的场景适用 ⭐⭐⭐⭐⭐

---

### 方案 3️⃣：分类组织（DTO 较多时）

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
├── models/                     # 内部模型
│   ├── chart.go
│   └── data_source.go
│
├── generate_report.go          # Handler
├── export_report.go
├── list_reports.go
├── controller.go
└── service_extensions.go
```

**示例代码：**

```go
// features/reports/requests/generate_report.go
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

```go
// features/reports/responses/report_detail.go
package responses

type ReportDetailResponse struct {
    ID          int64          `json:"id"`
    Name        string         `json:"name"`
    Type        string         `json:"report_type"`
    Charts      []ChartData    `json:"charts"`
    Summary     ReportSummary  `json:"summary"`
    GeneratedAt time.Time      `json:"generated_at"`
}

type ChartData struct {
    Type   string      `json:"type"`
    Title  string      `json:"title"`
    Data   interface{} `json:"data"`
    Labels []string    `json:"labels"`
}

type ReportSummary struct {
    TotalOrders  int     `json:"total_orders"`
    TotalRevenue float64 `json:"total_revenue"`
    AvgOrderValue float64 `json:"avg_order_value"`
}
```

```go
// features/reports/responses/report_list.go
package responses

type ReportListResponse struct {
    Reports []ReportListItem `json:"reports"`
    Total   int              `json:"total"`
}

type ReportListItem struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Type        string    `json:"type"`
    Status      string    `json:"status"`
    GeneratedAt time.Time `json:"generated_at"`
}
```

```go
// features/reports/generate_report.go
package reports

import (
    "vertical_slice_demo/apps/admin/features/reports/requests"
    "vertical_slice_demo/apps/admin/features/reports/responses"
)

type GenerateReportHandler struct {
    generator *business.ReportGenerator
}

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
- ❌ 可能过度设计
- ❌ 简单功能会显得复杂

**何时使用：**
- DTO 超过 10 个
- Request 和 Response 结构复杂
- 需要为 API 生成文档
- 团队较大（> 10 人）

---

### 方案 4️⃣：全局共享（❌ 不推荐）

**适用场景：** 传统分层架构（违反垂直切片原则）

**目录结构：**

```
shared/dtos/                    # ❌ 违反垂直切片原则
├── user_dto.go
├── product_dto.go
└── order_dto.go

features/users/
└── create_user.go             # 引用 shared/dtos
```

**示例代码：**

```go
// shared/dtos/user_dto.go
package dtos

type UserDTO struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

```go
// features/users/create_user.go
package users

import "vertical_slice_demo/shared/dtos"

type CreateUserHandler struct {
    userRepo repositories.IUserRepository
}

func (h *CreateUserHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req dtos.CreateUserRequest  // ❌ 依赖全局 DTO
    // ...
}
```

**优点：**
- ✅ 全局统一，完全复用
- ✅ 适合传统分层架构

**缺点：**
- ❌ 违反功能内聚原则
- ❌ 改 DTO 影响多个功能
- ❌ 不符合垂直切片理念
- ❌ 耦合度高，难以独立演进
- ❌ 删除功能时 DTO 残留

**何时使用：**
- 传统分层架构项目
- **垂直切片架构中不推荐使用**

---

## 📊 DTO 组织方案对比

| 方案 | 适用场景 | DTO 位置 | 复用性 | 内聚性 | 推荐度 |
|------|---------|---------|--------|--------|--------|
| 与操作一起 | 简单功能，DTO 差异大 | 每个 handler 文件中 | ❌ 低 | ⭐⭐⭐⭐⭐ 高 | ⭐⭐⭐⭐ |
| 功能内 models.go | 中等复杂，部分共享 | features/xxx/models.go | ✅ 中 | ⭐⭐⭐⭐ 高 | ⭐⭐⭐⭐⭐ |
| 分类组织 | 复杂功能，DTO > 10 | requests/ responses/ | ✅ 中 | ⭐⭐⭐ 中 | ⭐⭐⭐ |
| 全局共享 | 传统架构 | shared/dtos/ | ⭐⭐⭐ 高 | ❌ 低 | ❌ 不推荐 |

---

## 🎯 DTO 命名规范

### 1. Request DTOs

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

### 2. Response DTOs

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

### 3. List Item DTOs

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

## 💡 实战指南

### 指南 1：何时复用 vs 重新定义

**场景 A：格式完全相同 → 复用**

```go
// models.go
type UserResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// create_user.go
func (h *CreateUserHandler) Handle(c *web.HttpContext) {
    return c.Created(UserResponse{...})  // 复用
}

// update_user.go
func (h *UpdateUserHandler) Handle(c *web.HttpContext) {
    return c.Ok(UserResponse{...})  // 复用
}
```

**场景 B：格式不同 → 各自定义**

```go
// models.go

// 列表：只需要部分字段
type UserListItem struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

// 详情：需要完整信息
type UserDetailResponse struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Phone     string    `json:"phone"`
    Profile   *Profile  `json:"profile"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 指南 2：Domain 模型 vs DTO

**关键区别：**

```go
// shared/domain/user.go - 领域模型（内部使用，数据库映射）
type User struct {
    ID        int64
    Name      string
    Email     string
    Password  string    // ✅ 包含敏感字段
    Salt      string    // ✅ 内部字段
    IsDeleted bool      // ✅ 软删除标记
    CreatedAt time.Time
    UpdatedAt time.Time
}

// features/users/models.go - DTO（外部传输，API 响应）
type UserResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    // ❌ 没有 Password！
    // ❌ 没有 Salt！
    // ❌ 没有 IsDeleted！
}

// 转换函数
func toUserResponse(user *domain.User) *UserResponse {
    return &UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }
}
```

**原则：**
- Domain 模型：完整的数据结构，包含所有字段
- DTO：面向 API，只暴露必要字段
- 永远不要在 DTO 中暴露敏感信息

### 指南 3：转换函数的位置

**方式 A：在操作文件中（推荐简单场景）**

```go
// features/users/list_users.go
package users

func (h *ListUsersHandler) Handle(c *web.HttpContext) web.IActionResult {
    users, _ := h.userRepo.List(offset, limit)
    
    // ✅ 转换逻辑在这里
    items := make([]UserListItem, len(users))
    for i, user := range users {
        items[i] = UserListItem{
            ID:    user.ID,
            Name:  user.Name,
            Email: user.Email,
            Role:  user.Role,
        }
    }
    
    return c.Ok(items)
}
```

**方式 B：在 models.go 中（推荐复用场景）**

```go
// features/users/models.go
package users

// DTO 定义
type UserResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// ✅ 转换函数也在 models.go
func toUserResponse(user *domain.User) *UserResponse {
    return &UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }
}

func toUserListItem(user *domain.User) *UserListItem {
    return &UserListItem{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        Role:  user.Role,
    }
}
```

```go
// features/users/create_user.go
func (h *CreateUserHandler) Handle(c *web.HttpContext) {
    user := &domain.User{...}
    h.userRepo.Create(user)
    
    return c.Created(toUserResponse(user))  // ✅ 使用转换函数
}
```

**方式 C：独立的 mappers 文件（复杂场景）**

```go
// features/reports/mappers.go
package reports

// 复杂的转换逻辑
func toReportDetailResponse(report *models.Report, stats *Statistics) *responses.ReportDetailResponse {
    // 复杂的组装逻辑
    charts := make([]responses.ChartData, len(report.Charts))
    for i, chart := range report.Charts {
        charts[i] = responses.ChartData{
            Type:   chart.Type,
            Title:  chart.Title,
            Data:   transformChartData(chart.Data),
            Labels: generateLabels(chart),
        }
    }
    
    return &responses.ReportDetailResponse{
        ID:      report.ID,
        Name:    report.Name,
        Charts:  charts,
        Summary: buildSummary(stats),
    }
}
```

---

## 🔍 DTO 组织决策树

```
你的功能有多少个 DTO？
│
├─ < 3 个 DTO
│  └─→ 与操作放在一起
│     └─ features/xxx/create_xxx.go
│        └─ type Request struct 直接定义
│
├─ 3-10 个 DTO
│  │
│  ├─ DTO 之间有复用？
│  │  ├─ 是 → 功能内 models.go ⭐推荐
│  │  │  └─ features/xxx/models.go
│  │  │
│  │  └─ 否 → 与操作放在一起
│  │     └─ features/xxx/create_xxx.go
│  │
│  └─→ 功能内 models.go（默认选择）
│
└─ > 10 个 DTO
   └─→ 分类组织
      └─ features/xxx/
         ├── requests/
         ├── responses/
         └── models/
```

---

## ✅ DTO 最佳实践

### 1. 优先功能内聚

```go
// ✅ 好的做法：DTO 跟随功能
features/orders/
├── models.go           # 订单相关的所有 DTO
└── create_order.go

// ❌ 不好的做法：全局 DTO
shared/dtos/order_dto.go  # 过早抽象
```

### 2. 按需复用

```go
// ✅ 好的做法：只在功能内复用
// features/orders/models.go
type OrderResponse struct {...}  // 在功能内多处使用

// ❌ 不好的做法："可能会用到"
shared/dtos/order_dto.go  // 只有一个地方用
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
type UserModel struct             // 容易和 Domain 模型混淆
```

### 4. 隐藏敏感信息

```go
// shared/domain/user.go - Domain 模型
type User struct {
    ID       int64
    Name     string
    Email    string
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

// 使用 json:"-" 标签显式忽略
type UserDTO struct {
    ID       int64  `json:"id"`
    Name     string `json:"name"`
    Password string `json:"-"`  // ✅ 不会序列化到 JSON
}
```

### 5. 转换分离

```go
// ✅ 好的做法：专门的转换函数
func toUserResponse(user *domain.User) *UserResponse {
    return &UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
    }
}

// ❌ 不好的做法：在 Handler 中内联转换
func (h *Handler) Handle(c *web.HttpContext) {
    user := h.getUser()
    // 内联转换，不易复用
    response := map[string]interface{}{
        "id":    user.ID,
        "name":  user.Name,
        "email": user.Email,
    }
}
```

### 6. 验证规则在 DTO 中

```go
// ✅ 好的做法：使用 binding 标签
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required,min=2,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8,max=100"`
    Age      int    `json:"age" binding:"gte=0,lte=150"`
}

// Handler 中自动验证
func (h *CreateUserHandler) Handle(c *web.HttpContext) {
    var req CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err  // 自动返回 400 + 验证错误信息
    }
    // req 已经通过验证
}
```

---

## 🎨 实战案例分析

### 案例 1：简单 CRUD（用方案 1）

**场景：** 分类管理，只有基本的 CRUD

```go
// features/categories/handler.go - 所有代码在一个文件

package categories

// DTO 直接定义在这里
type Category struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

type CategoryHandler struct {
    categories map[int64]*Category
}

func (h *CategoryHandler) Create(c *web.HttpContext) {
    var req Category  // 直接用 Category 作为 Request
    c.MustBindJSON(&req)
    // ...
}

func (h *CategoryHandler) List(c *web.HttpContext) {
    // 直接返回 []Category
}
```

**评价：** ✅ 简单直接，适合简单场景

---

### 案例 2：订单管理（用方案 2，推荐）

**场景：** 订单功能，多个操作共享 DTO

```go
// features/orders/models.go - 集中管理 DTO

package orders

// 共享的响应结构
type OrderResponse struct {
    ID         int64       `json:"id"`
    Status     string      `json:"status"`
    TotalPrice float64     `json:"total_price"`
    Items      []OrderItem `json:"items"`
    CreatedAt  time.Time   `json:"created_at"`
}

type OrderItem struct {
    ProductID   int64   `json:"product_id"`
    ProductName string  `json:"product_name"`
    Quantity    int     `json:"quantity"`
    Price       float64 `json:"price"`
}

// 创建订单请求（独有）
type CreateOrderRequest struct {
    Items []CreateOrderItem `json:"items" binding:"required,min=1"`
}

type CreateOrderItem struct {
    ProductID int64 `json:"product_id" binding:"required"`
    Quantity  int   `json:"quantity" binding:"required,gt=0"`
}

// 列表项（简化版）
type OrderListItem struct {
    ID         int64     `json:"id"`
    Status     string    `json:"status"`
    TotalPrice float64   `json:"total_price"`
    ItemCount  int       `json:"item_count"`
    CreatedAt  time.Time `json:"created_at"`
}

// 详情（组合版）
type OrderDetailResponse struct {
    Order    OrderResponse `json:"order"`
    User     UserInfo      `json:"user"`
    Products []ProductInfo `json:"products"`
}

// 转换函数
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
        Status:     order.Status,
        TotalPrice: order.TotalPrice,
        ItemCount:  len(order.Items),
        CreatedAt:  order.CreatedAt,
    }
}
```

```go
// features/orders/create_order.go
func (h *CreateOrderHandler) Handle(c *web.HttpContext) {
    var req CreateOrderRequest  // 使用 models.go
    c.MustBindJSON(&req)
    
    order := createOrder(req)
    return c.Created(toOrderResponse(order))  // 使用转换函数
}

// features/orders/list_orders.go  
func (h *ListOrdersHandler) Handle(c *web.HttpContext) {
    orders, _ := h.orderRepo.List()
    
    items := make([]OrderListItem, len(orders))  // 使用 models.go
    for i, order := range orders {
        items[i] = toOrderListItem(order)  // 使用转换函数
    }
    
    return c.Ok(items)
}
```

**评价：** ✅ 结构清晰，复用合理，最推荐的方案

---

### 案例 3：报表系统（用方案 3）

**场景：** 复杂报表功能，DTO 超过 15 个

```
features/reports/
├── requests/
│   ├── generate_sales_report.go      # 销售报表请求
│   ├── generate_user_report.go       # 用户报表请求
│   └── export_report.go              # 导出请求
│
├── responses/
│   ├── sales_report_detail.go        # 销售报表详情响应
│   ├── user_report_detail.go         # 用户报表详情响应
│   ├── report_list.go                # 报表列表响应
│   └── export_result.go              # 导出结果响应
│
├── models/                            # 内部模型
│   ├── chart.go
│   ├── filter.go
│   └── aggregation.go
│
└── handlers...
```

**评价：** ✅ 适合复杂场景，分类清晰

---

## 📐 基于复杂度的选择

### 简单功能（< 3 个 DTO）

```
✅ 选择：与操作放在一起

features/tags/
└── handler.go
    ├── type Tag struct
    └── func Create/List/Update
```

### 中等功能（3-10 个 DTO）

```
✅ 选择：功能内 models.go ⭐推荐

features/orders/
├── models.go
│   ├── OrderResponse
│   ├── CreateOrderRequest
│   ├── OrderListItem
│   └── 转换函数
└── handlers...
```

### 复杂功能（> 10 个 DTO）

```
✅ 选择：分类组织

features/reports/
├── requests/
├── responses/
├── models/
└── handlers...
```

---

## 🎯 我的推荐

**80% 的场景用这个：**

```
features/xxx/
├── models.go              # ⭐ DTO 集中管理
│   ├── Request DTOs
│   ├── Response DTOs
│   ├── List Item DTOs
│   └── 转换函数
│
├── create_xxx.go         # Handler
├── list_xxx.go
└── service_extensions.go
```

**理由：**
1. ✅ 平衡了复用和内聚
2. ✅ DTO 集中，易于查找
3. ✅ 功能依然独立
4. ✅ 适用大多数场景
5. ✅ 团队友好

**特殊情况：**
- 非常简单（< 3 个 DTO）→ 直接放在 handler 中
- 非常复杂（> 10 个 DTO）→ 使用 requests/responses 分类

---

## ❓ DTO 组织常见问题与解决方案

### 问题 1：DTO 应该放在哪里？

**场景描述：**
开发新功能时，不确定 Request/Response DTO 应该定义在哪个位置。

**解决方案：**
根据功能复杂度和 DTO 数量选择合适的组织方式：

- **< 3 个 DTO** → 直接放在操作文件中
- **3-10 个 DTO** → 创建功能内 models.go
- **> 10 个 DTO** → 使用 requests/responses 分类

**代码示例：**

```go
// 方案 A：简单场景，直接在操作文件中定义
// features/tags/create_tag.go
package tags

type CreateTagRequest struct {
    Name  string `json:"name" binding:"required"`
    Color string `json:"color" binding:"required"`
}

type TagResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Color string `json:"color"`
}

func (h *CreateTagHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req CreateTagRequest
    // ... 处理逻辑
    return c.Created(TagResponse{...})
}
```

```go
// 方案 B：中等复杂，使用 models.go
// features/orders/models.go
package orders

type CreateOrderRequest struct {
    Items []OrderItem `json:"items"`
}

type OrderResponse struct {
    ID         int64       `json:"id"`
    TotalPrice float64     `json:"total_price"`
    Items      []OrderItem `json:"items"`
}

type OrderListItem struct {
    ID         int64   `json:"id"`
    TotalPrice float64 `json:"total_price"`
    ItemCount  int     `json:"item_count"`
}
```

```go
// 方案 C：复杂场景，使用 requests/responses 分类
// features/reports/requests/generate_report.go
package requests

type GenerateReportRequest struct {
    Name      string    `json:"name"`
    StartDate time.Time `json:"start_date"`
    Filters   []Filter  `json:"filters"`
}

// features/reports/responses/report_detail.go
package responses

type ReportDetailResponse struct {
    ID     int64   `json:"id"`
    Name   string  `json:"name"`
    Charts []Chart `json:"charts"`
}
```

---

### 问题 2：多端需要相同的 DTO 格式怎么办？

**场景描述：**
管理端创建用户和 C 端注册用户，返回的用户信息格式完全相同，是否需要共享 DTO？

**解决方案：**
- **Response 格式相同** → 共享到 `shared/contracts/dtos/`
- **Request 格式不同** → 各端独立定义
- **Response 格式不同** → 各端独立定义

**代码示例：**

```go
// shared/contracts/dtos/user_response.go
package dtos

// ✅ 共享的响应 DTO（保证 API 一致性）
type UserResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}
```

```go
// apps/admin/features/users/models.go
package users

import "vertical_slice_demo/shared/contracts/dtos"

// ✅ Request 独立（管理端可以指定角色）
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
// apps/api/features/auth/models.go
package auth

import "vertical_slice_demo/shared/contracts/dtos"

// ✅ Request 独立（C端不能指定角色）
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

---

### 问题 3：同一功能的多个操作有相同的子结构，如何组织？

**场景描述：**
报表功能的多个 Request 都需要 `Filter` 和 `DateRange` 结构，是否需要提取？

**解决方案：**
- **使用 < 3 次** → 不提取，直接定义
- **使用 3-5 次** → 提取到 `requests/common.go`
- **使用 > 5 次** → 提取到 `models/` 目录

**代码示例：**

```go
// 方案 A：提取到 requests/common.go（推荐）
// features/reports/requests/common.go
package requests

import "time"

// ✅ 请求相关的共享子结构
type Filter struct {
    Field    string      `json:"field"`
    Operator string      `json:"operator"`
    Value    interface{} `json:"value"`
}

type DateRange struct {
    StartDate time.Time `json:"start_date"`
    EndDate   time.Time `json:"end_date"`
}
```

```go
// features/reports/requests/generate_sales_report.go
package requests

type GenerateSalesReportRequest struct {
    Name      string    `json:"name"`
    DateRange DateRange `json:"date_range"`  // ✅ 复用 common.go
    Filters   []Filter  `json:"filters"`     // ✅ 复用 common.go
    GroupBy   string    `json:"group_by"`
}
```

```go
// features/reports/requests/generate_user_report.go
package requests

type GenerateUserReportRequest struct {
    Name      string    `json:"name"`
    DateRange DateRange `json:"date_range"`  // ✅ 复用
    Filters   []Filter  `json:"filters"`     // ✅ 复用
    UserType  string    `json:"user_type"`
}
```

```go
// 方案 B：提取到 models/ 目录（共享结构很多时）
// features/reports/models/filter.go
package models

type Filter struct {
    Field    string      `json:"field"`
    Operator string      `json:"operator"`
    Value    interface{} `json:"value"`
}
```

```go
// features/reports/requests/generate_report.go
package requests

import "vertical_slice_demo/apps/admin/features/reports/models"

type GenerateReportRequest struct {
    Name    string         `json:"name"`
    Filters []models.Filter `json:"filters"`  // ✅ 引用 models
}
```

---

### 问题 4：requests 和 responses 目录之间需要共享 DTO 怎么办？

**场景描述：**
`Filter`、`Pagination` 等结构既在 Request 中使用（提交过滤条件），也在 Response 中使用（返回当前过滤条件），应该放在哪里？

**解决方案：**
- **没有跨 requests/responses 的共享** → 各自的 common.go
- **有共享且 < 3 个** → 功能级 common.go
- **有共享且 >= 3 个** → models/ 目录（推荐）

**代码示例：**

```go
// 方案 A：提取到 models/ 目录（推荐）
// features/reports/models/filter.go
package models

// ✅ requests 和 responses 都可以用
type Filter struct {
    Field    string      `json:"field"`
    Operator string      `json:"operator"`
    Value    interface{} `json:"value"`
}
```

```go
// features/reports/models/pagination.go
package models

type Pagination struct {
    Page       int `json:"page"`
    PageSize   int `json:"page_size"`
    Total      int `json:"total"`
    TotalPages int `json:"total_pages"`
}
```

```go
// features/reports/requests/generate_report.go
package requests

import "vertical_slice_demo/apps/admin/features/reports/models"

type GenerateReportRequest struct {
    Name    string         `json:"name"`
    Filters []models.Filter `json:"filters"`  // ✅ 使用 models
}
```

```go
// features/reports/responses/report_detail.go
package responses

import "vertical_slice_demo/apps/admin/features/reports/models"

type ReportDetailResponse struct {
    ID             int64          `json:"id"`
    Name           string         `json:"name"`
    AppliedFilters []models.Filter `json:"applied_filters"`  // ✅ 使用 models
}
```

```go
// features/reports/requests/list_reports.go
package requests

import "vertical_slice_demo/apps/admin/features/reports/models"

type ListReportsRequest struct {
    Pagination models.Pagination `json:"pagination"`  // ✅ 使用 models
}
```

```go
// features/reports/responses/report_list.go
package responses

import "vertical_slice_demo/apps/admin/features/reports/models"

type ReportListResponse struct {
    Reports    []ReportItem       `json:"reports"`
    Pagination models.Pagination  `json:"pagination"`  // ✅ 使用 models
}
```

**依赖关系图：**

```
models/
├── filter.go
├── pagination.go
└── sort_option.go
    ↑
    ├─────────────┬─────────────┐
    │             │             │
requests/     responses/    handlers/
```

---

### 问题 5：如何区分 Domain 模型和 DTO？

**场景描述：**
User 既有 Domain 模型（`shared/domain/user.go`），又有 DTO（`UserResponse`），两者有什么区别？

**解决方案：**
- **Domain 模型**：完整的数据结构，包含所有字段（包括敏感字段），用于内部逻辑
- **DTO**：面向 API 的传输对象，只包含必要字段，隐藏敏感信息

**代码示例：**

```go
// shared/domain/user.go - Domain 模型（内部使用）
package domain

type User struct {
    ID        int64
    Name      string
    Email     string
    Password  string    // ✅ 包含敏感字段
    Salt      string    // ✅ 内部字段
    IsDeleted bool      // ✅ 软删除标记
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

```go
// shared/contracts/dtos/user_response.go - DTO（对外传输）
package dtos

type UserResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    // ❌ 不包含 Password、Salt、IsDeleted
}
```

```go
// 转换函数（在功能内定义）
// features/users/models.go
package users

import (
    "vertical_slice_demo/shared/contracts/dtos"
    "vertical_slice_demo/shared/domain"
)

// ✅ Domain → DTO 转换，隐藏敏感信息
func toUserResponse(user *domain.User) *dtos.UserResponse {
    return &dtos.UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        // 不暴露 Password、Salt 等
    }
}
```

**关键区别：**

| 维度 | Domain 模型 | DTO |
|------|------------|-----|
| 位置 | `shared/domain/` | `features/*/models.go` 或 `shared/contracts/dtos/` |
| 用途 | 内部业务逻辑 | 外部数据传输 |
| 字段 | 完整字段（包含敏感） | 只有必要字段 |
| 标签 | 可能无 JSON 标签 | 必有 JSON 标签 |
| 验证 | 业务规则验证 | 输入格式验证 |

---

### 问题 6：什么时候应该复用 DTO，什么时候应该重新定义？

**场景描述：**
创建用户和更新用户都返回用户信息，Response 是复用还是各自定义？

**解决方案：**
- **格式完全相同** → 复用
- **格式有差异** → 各自定义
- **现在相同但可能变化** → 先复用，需要时再拆分

**代码示例：**

```go
// features/users/models.go

// 场景 A：格式完全相同 → 复用
type UserResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}

// create_user.go 和 update_user.go 都使用 UserResponse

// 场景 B：格式不同 → 各自定义
type UserListItem struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
    Role string `json:"role"`
    // 列表不需要 Email
}

type UserDetailResponse struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    Phone     string    `json:"phone"`
    Profile   *Profile  `json:"profile"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    // 详情需要更多信息
}
```

**决策规则：**

```
两个操作的 Response 格式是否完全相同？
├─ 是 → 复用
│  └─ type UserResponse struct {...}
│
├─ 否 → 各自定义
│  ├─ type UserListItem struct {...}
│  └─ type UserDetailResponse struct {...}
│
└─ 大部分相同，少量差异 → 使用嵌入
   ├─ type BaseUserResponse struct {...}
   └─ type UserDetailResponse struct {
          BaseUserResponse
          ExtraField string
      }
```

---

## 🗺️ DTO 组织完整决策指南

### 决策树 1：基本组织方式

```
开始：需要定义 DTO
│
├─ 功能需要被多个端使用吗？
│  │
│  ├─ 是 → Response 格式完全相同吗？
│  │  │
│  │  ├─ 是 → shared/contracts/dtos/xxx_response.go
│  │  │     └─ Response 共享，Request 各端独立
│  │  │
│  │  └─ 否 → 各端 features/xxx/models.go
│  │        └─ 格式不同，各自维护
│  │
│  └─ 否 → 继续判断复杂度
│     │
│     ├─ DTO 数量
│     │  ├─ < 3 个 → 与操作放在一起
│     │  │   └─ features/xxx/create_xxx.go
│     │  │       └─ type Request/Response 直接定义
│     │  │
│     │  ├─ 3-10 个 → 功能内 models.go ⭐
│     │  │   └─ features/xxx/models.go
│     │  │       ├─ Request DTOs
│     │  │       ├─ Response DTOs
│     │  │       └─ 转换函数
│     │  │
│     │  └─ > 10 个 → 分类组织
│     │      └─ features/xxx/
│     │          ├─ requests/
│     │          │   ├─ create_xxx.go
│     │          │   └─ update_xxx.go
│     │          └─ responses/
│     │              ├─ xxx_detail.go
│     │              └─ xxx_list.go
│     │
│     └─ 需要清晰的 API 文档吗？
│        ├─ 是 → 使用 requests/responses 分类
│        └─ 否 → 使用 models.go
```

### 决策树 2：功能内共享（requests/responses 目录场景）

```
有 requests/ 和 responses/ 目录
│
├─ requests 内有重复的子结构吗？
│  │
│  ├─ 是 → 使用次数
│  │  ├─ 1-2 次 → 不提取，直接定义
│  │  ├─ 3-5 次 → requests/common.go
│  │  └─ > 5 次 → 考虑 models/ 目录
│  │
│  └─ 否 → 各文件独立定义
│
├─ responses 内有重复的子结构吗？
│  │
│  ├─ 是 → 使用次数
│  │  ├─ 1-2 次 → 不提取，直接定义
│  │  ├─ 3-5 次 → responses/common.go
│  │  └─ > 5 次 → 考虑 models/ 目录
│  │
│  └─ 否 → 各文件独立定义
│
└─ requests 和 responses 有共同的子结构吗？
   │
   ├─ 是 → 共享结构数量
   │  ├─ < 3 个 → features/xxx/common.go
   │  │   └─ 功能级共享
   │  │
   │  └─ >= 3 个 → features/xxx/models/ ⭐
   │      └─ 独立目录管理
   │          ├─ filter.go
   │          ├─ pagination.go
   │          └─ sort_option.go
   │
   └─ 否 → 各自的 common.go
      ├─ requests/common.go
      └─ responses/common.go
```

### 决策树 3：跨端共享

```
这个 DTO 需要跨端使用吗？
│
├─ 是 → Request 还是 Response？
│  │
│  ├─ Request → 各端需求相同吗？
│  │  ├─ 相同 → shared/contracts/dtos/xxx_request.go
│  │  │   └─ 少见，谨慎共享
│  │  │
│  │  └─ 不同 → 各端独立 ⭐
│  │      ├─ apps/admin/features/xxx/models.go
│  │      └─ apps/api/features/xxx/models.go
│  │
│  └─ Response → 格式完全相同吗？
│     │
│     ├─ 相同 → shared/contracts/dtos/xxx_response.go ⭐
│     │   └─ 保证 API 一致性
│     │
│     └─ 不同 → 各端独立
│         ├─ apps/admin/features/xxx/models.go
│         │   └─ type AdminUserResponse（更多字段）
│         └─ apps/api/features/xxx/models.go
│             └─ type UserProfileResponse（基础字段）
│
└─ 否 → 放在各自 features/ 下
   └─ 按功能内组织方式处理
```

---

## 📊 DTO 组织方案全面对比

### 表 1：基本组织方式对比

| 方案 | DTO 位置 | 适用场景 | DTO 数量 | 复用性 | 内聚性 | 查找难度 | 推荐度 |
|------|---------|---------|----------|--------|--------|----------|--------|
| 与操作一起 | `create_xxx.go` 中直接定义 | 简单 CRUD | < 3 个 | ❌ 低 | ⭐⭐⭐⭐⭐ 高 | ✅ 易 | ⭐⭐⭐⭐ |
| 功能内 models.go | `features/xxx/models.go` | 中等复杂 | 3-10 个 | ⭐⭐⭐ 中 | ⭐⭐⭐⭐ 高 | ✅ 易 | ⭐⭐⭐⭐⭐ |
| 分类组织 | `requests/` + `responses/` | 复杂功能 | > 10 个 | ⭐⭐⭐ 中 | ⭐⭐⭐ 中 | ⭐⭐ 中 | ⭐⭐⭐ |
| 全局共享 | `shared/dtos/` | 传统架构 | 任意 | ⭐⭐⭐⭐ 高 | ❌ 低 | ⭐⭐⭐ 中 | ❌ 不推荐 |

**推荐组合：**
- 80% 场景：功能内 models.go ⭐⭐⭐⭐⭐
- 简单场景：与操作一起
- 复杂场景：分类组织

---

### 表 2：跨端共享方案对比

| 方案 | 位置 | Request | Response | 适用场景 | 耦合度 | 灵活性 | 推荐度 |
|------|------|---------|----------|---------|--------|--------|--------|
| 完全不共享 | 各端 `features/` | 独立 | 独立 | 格式差异大 | ✅ 低 | ⭐⭐⭐⭐⭐ 高 | ⭐⭐⭐ |
| Response 共享 | `shared/contracts/dtos/` | 独立 | 共享 | API 需要一致性 | ⭐⭐ 中 | ⭐⭐⭐ 中 | ⭐⭐⭐⭐⭐ |
| 全部共享 | `shared/contracts/dtos/` | 共享 | 共享 | 格式完全相同 | ⭐⭐⭐ 高 | ⭐⭐ 低 | ⭐⭐ |
| 全局 dtos | `shared/dtos/` | 共享 | 共享 | 传统分层 | ⭐⭐⭐⭐ 高 | ❌ 低 | ❌ 不推荐 |

**推荐策略：**
- ✅ Response 倾向共享（保证一致性）
- ✅ Request 倾向独立（保持灵活性）

---

### 表 3：功能内共享方案对比（requests/responses 场景）

| 方案 | 位置 | 适用场景 | 共享结构数量 | 依赖关系 | 查找难度 | 推荐度 |
|------|------|---------|-------------|----------|----------|--------|
| 不共享 | 各文件内定义 | 使用 1-2 次 | < 3 个 | 无依赖 | ✅ 易 | ⭐⭐⭐ |
| requests/common.go | `requests/common.go` | 只 requests 内共享 | 3-5 个 | requests 内部 | ✅ 易 | ⭐⭐⭐⭐ |
| responses/common.go | `responses/common.go` | 只 responses 内共享 | 3-5 个 | responses 内部 | ✅ 易 | ⭐⭐⭐⭐ |
| models/ 目录 | `models/*.go` | requests + responses 共享 | >= 3 个 | 清晰 | ⭐⭐ 中 | ⭐⭐⭐⭐⭐ |
| 功能级 common.go | `features/xxx/common.go` | 跨目录少量共享 | < 3 个 | 子包 → 父包 | ✅ 易 | ⭐⭐⭐⭐ |

**推荐策略：**
- ✅ 有跨 requests/responses 共享 → models/ 目录
- ✅ 只在 requests 内 → requests/common.go
- ✅ 只在 responses 内 → responses/common.go

---

### 表 4：requests 和 responses 跨目录共享对比

| 方案 | 共享结构位置 | 依赖关系 | 适用场景 | 语义清晰度 | 维护成本 | 推荐度 |
|------|-------------|---------|---------|-----------|----------|--------|
| 重复定义 | requests/common.go<br>responses/common.go | 无（重复） | ❌ 不推荐 | ⭐⭐ | ❌ 高 | ❌ |
| 互相引用 | responses/common.go | requests → responses | ❌ 不推荐 | ❌ 低 | ⭐⭐ 中 | ❌ |
| models/ 目录 | models/*.go | requests → models<br>responses → models | >= 3 个共享结构 | ⭐⭐⭐⭐⭐ 高 | ✅ 低 | ⭐⭐⭐⭐⭐ |
| 功能级 common.go | features/xxx/common.go | requests → parent<br>responses → parent | < 3 个共享结构 | ⭐⭐⭐⭐ 高 | ✅ 低 | ⭐⭐⭐⭐ |

**推荐策略：**
- ✅ 默认选择：models/ 目录（依赖清晰）
- ✅ 共享少时：功能级 common.go
- ❌ 避免：重复定义、互相引用

---

### 表 5：典型共享结构分类

| 结构类型 | 常见名称 | 在 Request 中 | 在 Response 中 | 是否需要跨 requests/responses 共享 | 推荐位置 |
|---------|---------|--------------|---------------|--------------------------------|---------|
| 分页信息 | Pagination | ✅ 请求第几页 | ✅ 返回分页元数据 | ✅ 是 | models/pagination.go |
| 过滤条件 | Filter | ✅ 提交过滤 | ✅ 返回当前过滤 | ✅ 是 | models/filter.go |
| 排序选项 | SortOption | ✅ 请求排序 | ✅ 返回当前排序 | ✅ 是 | models/sort_option.go |
| 时间范围 | DateRange | ✅ 查询范围 | ✅ 返回应用范围 | ✅ 是 | models/date_range.go |
| 验证规则 | ValidationRule | ✅ 提交验证 | ❌ 不在响应中 | ❌ 否 | requests/common.go |
| 统计信息 | Statistics | ❌ 不在请求中 | ✅ 返回统计 | ❌ 否 | responses/common.go |
| 图表数据 | Chart | ❌ 不在请求中 | ✅ 返回图表 | ❌ 否 | responses/common.go |
| 元数据 | Metadata | ❌ 不在请求中 | ✅ 返回元数据 | ❌ 否 | responses/common.go |

---

## 📋 DTO 组织检查清单

### 新建功能时

- [ ] 确定 DTO 数量（< 3, 3-10, > 10）
- [ ] 选择基本组织方式（操作文件、models.go、requests/responses）
- [ ] 检查是否需要跨端共享
- [ ] 确定 Request 和 Response 是否需要独立定义
- [ ] 检查是否有重复的子结构

### 功能演进时

- [ ] DTO 超过 3 个 → 考虑提取到 models.go
- [ ] DTO 超过 10 个 → 考虑使用 requests/responses 分类
- [ ] 有重复子结构 → 考虑提取到 common.go 或 models/
- [ ] 第二个端也需要 → 考虑提取到 shared/contracts/dtos/
- [ ] Response 格式相同 → 优先共享

### 代码审查时

- [ ] DTO 放置位置是否合理
- [ ] 是否有不必要的重复定义
- [ ] 是否过早抽象到 shared
- [ ] 命名是否清晰（CreateXxxRequest, XxxResponse）
- [ ] 是否隐藏了敏感信息（Password, Salt）
- [ ] 转换函数是否放在合适的位置

---

## 🎯 快速参考

### 80% 场景的默认选择

```
features/xxx/
├── models.go              # ⭐ DTO 集中管理
│   ├── Request DTOs
│   ├── Response DTOs
│   ├── List Item DTOs
│   └── 转换函数
│
├── create_xxx.go         # Handler
├── list_xxx.go
├── update_xxx.go
├── controller.go
└── service_extensions.go
```

### 跨端共享的默认选择

```
shared/contracts/dtos/
├── user_response.go       # ⭐ Response 共享
├── product_response.go
└── order_response.go

apps/admin/features/users/models.go
└── CreateUserRequest      # ⭐ Request 独立

apps/api/features/auth/models.go
└── RegisterRequest        # ⭐ Request 独立
```

### requests/responses 共享的默认选择

```
features/reports/
├── models/                # ⭐ 跨目录共享结构
│   ├── filter.go
│   ├── pagination.go
│   └── sort_option.go
│
├── requests/
│   ├── generate_report.go
│   └── list_reports.go
│
└── responses/
    ├── report_detail.go
    └── report_list.go
```

---

## 🔄 何时重构？

### 从模式 1 → 模式 2

**触发条件：**
- 文件超过 200 行
- 有 3 个以上的操作
- 团队有 2 人以上协作

**重构步骤：**
1. 提取数据模型到 `models.go`
2. 提取数据访问到 `store.go`
3. 拆分每个操作到独立文件
4. 创建 `controller.go` 统一路由

### 从模式 2 → 模式 3

**触发条件：**
- 文件总数超过 1000 行
- 有复杂的业务逻辑需要复用
- 需要单独测试业务层

**重构步骤：**
1. 创建 `models/` 目录，移动所有模型
2. 创建 `data/` 目录，移动数据访问
3. 创建 `business/` 目录，提取业务逻辑
4. Handler 只保留 HTTP 处理逻辑

### 从 features → shared

**触发条件：**
- 第二个端也需要这个功能
- 明确的跨端复用需求

**重构步骤：**
1. 移动数据模型到 `shared/domain/`
2. 定义接口到 `shared/contracts/`
3. 移动实现到 `shared/repositories/` 或 `shared/services/`
4. 更新所有引用

---

## ⚠️ 常见错误

### ❌ 错误 1：过早抽象

```
# BAD: 一开始就创建
shared/domain/admin_log.go        # 只有管理端用
shared/repositories/admin_log_repository.go

# GOOD: 先放在功能内
apps/admin/features/logs/log_entry.go
apps/admin/features/logs/handler.go
```

### ❌ 错误 2：过度分层

```
# BAD: 为了分层而分层
features/simple_crud/
├── models/
├── dtos/
├── mappers/
├── validators/
├── services/
├── repositories/
└── handlers/

# GOOD: 简单的东西简单做
features/simple_crud/
├── handler.go
└── service_extensions.go
```

### ❌ 错误 3：命名混乱

```
# BAD
features/user_manage/user_stuff.go
features/product/do_something.go

# GOOD
features/users/create_user.go
features/products/update_product.go
```

---

## ✅ 最佳实践总结

1. **从简单开始**：先单文件，不行再拆分
2. **功能内聚**：一个功能的所有代码在一个目录
3. **避免过早抽象**：只在需要时才提取到 shared
4. **清晰命名**：文件名清楚表达功能
5. **内部私有**：内部的 Store、Business 不暴露
6. **混合使用**：可以同时使用 shared 和 private
7. **及时重构**：代码超过阈值就重构

---

## 📚 相关文档

- [ARCHITECTURE.md](ARCHITECTURE.md) - 架构设计
- [shared/services/README.md](shared/services/README.md) - 共享服务
- [shared/repositories/README.md](shared/repositories/README.md) - 仓储层

---

**记住：架构是演进的，不是设计的！** 🚀

