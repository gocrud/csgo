# 项目结构组织

## 目录

- [核心原则](#核心原则)
- [标准目录结构说明](#标准目录结构说明)
- [决策树：shared vs features](#决策树shared-vs-features)
- [三种复杂度组织模式](#三种复杂度组织模式)
  - [模式1：单文件实现](#模式1单文件实现)
  - [模式2：按操作拆分](#模式2按操作拆分)
  - [模式3：内部分层](#模式3内部分层)
- [决策树：如何选择组织模式](#决策树如何选择组织模式)
- [文件组织规范](#文件组织规范)
- [何时重构](#何时重构)

---

## 核心原则

**不要过早抽象！** 只有在确定需要跨端共享时，才将代码提取到 `shared/`。

---

## 标准目录结构说明

### apps/ 目录

```
apps/
├── admin/              # 管理端应用
│   ├── features/       # 功能模块
│   └── internal/       # 应用内部共享（可选）
└── api/                # C端应用
    ├── features/
    └── internal/
```

**用途：** 每个独立的应用端，包含该端特有的功能和逻辑。

### shared/ 目录

```
shared/
├── domain/             # 共享领域模型
│   ├── user.go
│   ├── product.go
│   └── common/         # 公共基础模型
│       └── base_entity.go
├── repositories/       # 共享数据访问
├── services/          # 共享业务服务
└── contracts/         # 共享契约
    └── dtos/          # 跨端共享的 DTO
```

**用途：** 多个端需要共享的代码。

### features/ 目录

```
apps/admin/features/
├── users/             # 用户管理功能
├── products/          # 商品管理功能
└── orders/            # 订单管理功能
```

**用途：** 按功能垂直切分的模块。

---

## 决策树：shared vs features

```
你的功能需要被多个端使用吗？
├─ 是 → 放在 shared/
│  ├─ 数据模型 → shared/domain/
│  ├─ 数据访问 → shared/repositories/
│  └─ 业务服务 → shared/services/
│
└─ 否 → 放在对应端的 features/
   └─ apps/*/features/xxx/
```

---

## 三种复杂度组织模式

### 模式1：单文件实现

**适用场景：**
- 简单的 CRUD 操作
- 代码量 < 200 行
- 业务逻辑简单

**目录结构：**
```
apps/admin/features/categories/
├── handler.go               # 所有功能
└── service_extensions.go    # DI 注册
```

**handler.go 内容示例：**

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

### 模式2：按操作拆分

**适用场景：**
- 多个操作
- 代码量 200-1000 行
- 需要团队协作

**目录结构：**
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

### 模式3：内部分层

**适用场景：**
- 复杂业务逻辑
- 代码量 > 1000 行
- 需要内部复用
- 需要明确内部/外部边界

**目录结构（推荐使用 internal）：**
```
apps/admin/features/reports/
├── internal/                        # 内部实现（不暴露给外部）
│   ├── entity/                      # 内部实体/领域对象
│   │   ├── report_entity.go
│   │   ├── template_entity.go
│   │   └── config_entity.go
│   │
│   ├── data/                        # 数据访问层
│   │   ├── report_store.go
│   │   └── template_store.go
│   │
│   └── business/                    # 业务逻辑层
│       ├── report_generator.go
│       ├── report_exporter.go
│       ├── data_aggregator.go
│       └── chart_builder.go
│
├── models.go                        # 对外 DTO（Request/Response）
├── generate_report.go               # Handler（对外接口）
├── export_report.go
├── schedule_report.go
├── list_reports.go
├── controller.go
└── service_extensions.go
```

**internal 目录说明：**
- ✅ Go 语言特性：`internal/` 包只能被其父目录及子目录导入
- ✅ 封装内部实现：外部无法直接访问 internal 下的代码
- ✅ 明确边界：清晰区分对外接口和内部实现
- ✅ 防止误用：避免其他模块直接依赖内部实现

**internal/entity/report_entity.go:**

```go
package entity

import "time"

// ReportEntity 内部报表实体（不对外暴露）
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

type ChartConfig struct {
    Type   string
    Title  string
    Config map[string]interface{}
}
```

**internal/data/report_store.go:**

```go
package data

import (
    "sync"
    "vertical_slice_demo/apps/admin/features/reports/internal/entity"
)

// ReportStore 内部数据访问（不暴露）
type ReportStore struct {
    reports map[int64]*entity.ReportEntity
    mu      sync.RWMutex
    nextID  int64
}

func NewReportStore() *ReportStore {
    return &ReportStore{
        reports: make(map[int64]*entity.ReportEntity),
        nextID:  1,
    }
}

func (s *ReportStore) Create(report *entity.ReportEntity) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    report.ID = s.nextID
    s.nextID++
    s.reports[report.ID] = report
    return nil
}

func (s *ReportStore) GetByID(id int64) (*entity.ReportEntity, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    report, exists := s.reports[id]
    if !exists {
        return nil, errors.New("report not found")
    }
    return report, nil
}
```

**internal/business/report_generator.go:**

```go
package business

import (
    "vertical_slice_demo/apps/admin/features/reports/internal/entity"
    "vertical_slice_demo/apps/admin/features/reports/internal/data"
)

// ReportGenerator 报表生成器（内部业务逻辑）
type ReportGenerator struct {
    store      *data.ReportStore
    aggregator *DataAggregator
    builder    *ChartBuilder
}

func NewReportGenerator(
    store *data.ReportStore,
    aggregator *DataAggregator,
    builder *ChartBuilder,
) *ReportGenerator {
    return &ReportGenerator{
        store:      store,
        aggregator: aggregator,
        builder:    builder,
    }
}

func (g *ReportGenerator) Generate(config entity.ReportConfig) (*entity.ReportEntity, error) {
    // 1. 聚合数据
    aggregated := g.aggregator.Aggregate(config)
    
    // 2. 构建图表
    charts := g.builder.BuildCharts(aggregated, config.ChartTypes)
    
    // 3. 生成报表
    report := &entity.ReportEntity{
        Name:   config.Name,
        Type:   config.Type,
        Charts: charts,
    }
    
    // 4. 保存报表
    if err := g.store.Create(report); err != nil {
        return nil, err
    }
    
    return report, nil
}
```

**models.go (对外 DTO):**

```go
package reports

import "time"

// 对外的 Request DTO
type GenerateReportRequest struct {
    Name       string    `json:"name" binding:"required"`
    ReportType string    `json:"report_type" binding:"required"`
    StartDate  time.Time `json:"start_date" binding:"required"`
    EndDate    time.Time `json:"end_date" binding:"required"`
    GroupBy    string    `json:"group_by" binding:"required"`
}

// 对外的 Response DTO
type ReportResponse struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Type        string    `json:"type"`
    Charts      []Chart   `json:"charts"`
    GeneratedAt time.Time `json:"generated_at"`
}

type Chart struct {
    Type  string      `json:"type"`
    Title string      `json:"title"`
    Data  interface{} `json:"data"`
}
```

**generate_report.go (Handler):**

```go
package reports

import (
    "github.com/gocrud/csgo/web"
    "vertical_slice_demo/apps/admin/features/reports/internal/business"
    "vertical_slice_demo/apps/admin/features/reports/internal/entity"
)

type GenerateReportHandler struct {
    generator *business.ReportGenerator
}

func NewGenerateReportHandler(generator *business.ReportGenerator) *GenerateReportHandler {
    return &GenerateReportHandler{generator: generator}
}

func (h *GenerateReportHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req GenerateReportRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    // 转换 DTO 到内部实体配置
    config := entity.ReportConfig{
        Name:      req.Name,
        Type:      entity.ReportType(req.ReportType),
        StartDate: req.StartDate,
        EndDate:   req.EndDate,
        GroupBy:   req.GroupBy,
    }
    
    // 调用业务层生成报表
    report, err := h.generator.Generate(config)
    if err != nil {
        return c.InternalError("生成报表失败")
    }
    
    // 转换内部实体到 Response DTO
    response := toReportResponse(report)
    
    return c.Ok(response)
}

// 转换函数：内部实体 → DTO
func toReportResponse(report *entity.ReportEntity) *ReportResponse {
    charts := make([]Chart, len(report.Charts))
    for i, c := range report.Charts {
        charts[i] = Chart{
            Type:  c.Type,
            Title: c.Title,
            Data:  c.Config,
        }
    }
    
    return &ReportResponse{
        ID:          report.ID,
        Name:        report.Name,
        Type:        string(report.Type),
        Charts:      charts,
        GeneratedAt: report.CreatedAt,
    }
}
```

**service_extensions.go:**

```go
package reports

import (
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
    "vertical_slice_demo/apps/admin/features/reports/internal/data"
    "vertical_slice_demo/apps/admin/features/reports/internal/business"
)

func AddReportFeature(services di.IServiceCollection) {
    // 注册内部数据层
    services.AddSingleton(data.NewReportStore)
    
    // 注册内部业务层
    services.AddSingleton(business.NewReportGenerator)
    services.AddSingleton(business.NewDataAggregator)
    services.AddSingleton(business.NewChartBuilder)
    
    // 注册 Handlers（对外接口）
    services.AddSingleton(NewGenerateReportHandler)
    services.AddSingleton(NewExportReportHandler)
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
- ✅ 使用 internal 封装内部实现
- ✅ 防止外部模块误用内部代码

**缺点：**
- ❌ 结构复杂
- ❌ 需要更多的规划

---

## 决策树：如何选择组织模式

```
你的功能代码量有多少？
│
├─ < 200 行
│  └─→ 模式1：单文件实现
│     └─ features/xxx/handler.go
│
├─ 200-1000 行
│  └─→ 模式2：按操作拆分 ⭐推荐
│     └─ features/xxx/
│        ├── models.go
│        ├── create_xxx.go
│        └── list_xxx.go
│
└─ > 1000 行
   └─→ 模式3：内部分层
      └─ features/xxx/
         ├── internal/
         │   ├── entity/
         │   ├── data/
         │   └── business/
         ├── models.go
         └── handlers...
```

---

## 文件组织规范

### handler.go vs controller.go

- `handler.go` - 单文件模式，包含所有逻辑
- `controller.go` - 拆分模式，只负责路由映射

### service_extensions.go

- 每个 feature 必须有
- 负责 DI 注册
- 统一命名：`Add{Feature}Feature`

### 命名规范

- **文件名：** 小写+下划线 `create_user.go`
- **Handler：** 动词+名词+Handler `CreateUserHandler`
- **Controller：** 名词+Controller `UserController`
- **Store：** 名词+Store `UserStore`（内部使用）

---

## 何时重构

### 从模式1 → 模式2

**触发条件：**
- 文件超过 200 行
- 有 3 个以上的操作
- 团队有 2 人以上协作

**重构步骤：**
1. 提取数据模型到 `models.go`
2. 提取数据访问到 `store.go`
3. 拆分每个操作到独立文件
4. 创建 `controller.go` 统一路由

---

### 从模式2 → 模式3

**触发条件：**
- 文件总数超过 1000 行
- 有复杂的业务逻辑需要复用
- 需要单独测试业务层
- 需要明确封装内部实现

**重构步骤：**
1. 创建 `internal/` 目录
2. 创建 `internal/entity/` 目录，移动或创建内部实体
3. 创建 `internal/data/` 目录，移动数据访问（原 store.go）
4. 创建 `internal/business/` 目录，提取业务逻辑
5. Handler 只保留 HTTP 处理逻辑和 DTO 转换
6. models.go 保留在外层，作为对外 DTO

---

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

## 常见错误

### ❌ 错误1：过早抽象

```
# BAD: 一开始就创建
shared/domain/admin_log.go        # 只有管理端用
shared/repositories/admin_log_repository.go

# GOOD: 先放在功能内
apps/admin/features/logs/models.go
apps/admin/features/logs/handler.go
```

### ❌ 错误2：过度分层

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

### ❌ 错误3：命名混乱

```
# BAD
features/user_manage/user_stuff.go
features/product/do_something.go

# GOOD
features/users/create_user.go
features/products/update_product.go
```

### ❌ 错误4：internal 使用不当

```
# BAD: 简单功能不需要 internal
features/simple_crud/
└── internal/
    └── handler.go      # ❌ Handler 应该在外面

# GOOD
features/simple_crud/
└── handler.go          # ✅ 简单功能直接暴露
```

---

## 最佳实践总结

1. **从简单开始**：先单文件，不行再拆分
2. **功能内聚**：一个功能的所有代码在一个目录
3. **避免过早抽象**：只在需要时才提取到 shared
4. **清晰命名**：文件名清楚表达功能
5. **使用 internal 封装**：复杂功能使用 internal 保护内部实现
6. **及时重构**：代码超过阈值就重构

---

**返回 [主文档](../ORGANIZATION_GUIDE.md)**
