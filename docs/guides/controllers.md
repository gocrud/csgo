# 控制器指南

本指南介绍如何使用 ASP.NET Core 风格的 Controller 模式来组织和构建 API。

## 目录

- [什么是控制器模式](#什么是控制器模式)
- [快速开始](#快速开始)
- [IController 接口](#icontroller-接口)
- [使用 HttpContext 和 ActionResult](#使用-httpcontext-和-actionresult)
- [控制器注册](#控制器注册)
- [路由定义](#路由定义)
- [项目结构](#项目结构)
- [最佳实践](#最佳实践)
- [与 .NET 对比](#与-net-对比)

## 什么是控制器模式

Controller 模式是一种将相关的 HTTP 请求处理逻辑组织在一起的设计模式，源自 MVC（Model-View-Controller）架构。

### 优势

| 优势 | 说明 |
|------|------|
| **清晰的代码组织** | 按功能模块分组 |
| **关注点分离** | Controller、Service、Model 分层 |
| **依赖注入** | 构造函数注入模式 |
| **易于测试** | 可以轻松 mock 依赖 |
| **团队协作** | 不同开发者负责不同控制器 |

## 快速开始

### 1. 定义控制器

```go
package controllers

import (
    "github.com/gocrud/csgo/web"
)

// UserController 处理用户相关的 HTTP 请求
type UserController struct {
    userService UserService
}

// NewUserController 创建控制器（构造函数注入）
func NewUserController(userService UserService) *UserController {
    return &UserController{userService: userService}
}

// MapRoutes 实现 web.IController 接口
func (ctrl *UserController) MapRoutes(app *web.WebApplication) {
    users := app.MapGroup("/api/users")
    users.WithTags("Users")
    
    users.MapGet("", ctrl.GetAll)
    users.MapGet("/:id", ctrl.GetByID)
    users.MapPost("", ctrl.Create)
    users.MapPut("/:id", ctrl.Update)
    users.MapDelete("/:id", ctrl.Delete)
}

// GetAll 处理 GET /api/users - 使用 ActionResult
func (ctrl *UserController) GetAll(c *web.HttpContext) web.IActionResult {
    users := ctrl.userService.ListUsers()
    return c.Ok(users)
}

// GetByID 处理 GET /api/users/:id
func (ctrl *UserController) GetByID(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    user := ctrl.userService.GetUser(id)
    if user == nil {
        return c.NotFound("用户不存在")
    }
    
    return c.Ok(user)
}
```

### 2. 注册控制器

```go
package controllers

import (
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

// AddControllers 注册所有控制器
func AddControllers(services di.IServiceCollection) {
    // 使用 web.AddController 自动发现注册
    web.AddController(services, func(sp di.IServiceProvider) *UserController {
        userService := di.GetRequiredService[UserService](sp)
        return NewUserController(userService)
    })
    
    web.AddController(services, func(sp di.IServiceProvider) *OrderController {
        orderService := di.GetRequiredService[OrderService](sp)
        return NewOrderController(orderService)
    })
}
```

### 3. 使用控制器

```go
package main

import (
    "myapp/controllers"
    "myapp/services"
    "github.com/gocrud/csgo/web"
)

func main() {
    builder := web.CreateBuilder()
    
    // 注册服务
    services.AddServices(builder.Services)
    
    // 注册控制器
    controllers.AddControllers(builder.Services)
    
    // 构建应用
    app := builder.Build()
    
    // 自动映射所有控制器路由
    app.MapControllers()
    
    app.Run()
}
```

## IController 接口

CSGO 提供 `IController` 接口用于定义控制器：

```go
// IController 定义控制器接口
type IController interface {
    // MapRoutes 注册控制器的路由
    MapRoutes(app *WebApplication)
}
```

### 实现接口

```go
type ProductController struct {
    productService ProductService
}

// 确保实现了 IController 接口
var _ web.IController = (*ProductController)(nil)

func (ctrl *ProductController) MapRoutes(app *web.WebApplication) {
    products := app.MapGroup("/api/products")
    products.MapGet("", ctrl.List)
    products.MapPost("", ctrl.Create)
}
```

### ControllerBase（可选）

CSGO 提供可选的 `ControllerBase` 基类：

```go
type MyController struct {
    web.ControllerBase  // 嵌入基类
    myService MyService
}

func NewMyController(sp di.IServiceProvider, myService MyService) *MyController {
    return &MyController{
        ControllerBase: web.NewControllerBase(sp),
        myService:      myService,
    }
}
```

## 使用 HttpContext 和 ActionResult

CSGO 提供 `HttpContext` 和 `IActionResult` 来简化请求处理，类似 .NET 的 Controller 模式。

### HttpContext

`HttpContext` 包装了 `gin.Context`，提供更便捷的 API：

```go
func (ctrl *UserController) GetByID(c *web.HttpContext) web.IActionResult {
    // 路径参数（自动错误处理）
    id, err := c.MustPathInt("id")
    if err != nil {
        return err  // 自动返回 400 Bad Request
    }
    
    // 查询参数（带默认值）
    page := c.QueryInt("page", 1)
    size := c.QueryInt("size", 10)
    
    // 原始 gin.Context 方法仍然可用
    name := c.Query("name")
    
    return c.Ok(data)
}
```

### ActionResult

`IActionResult` 提供统一的响应处理：

```go
func (ctrl *UserController) Create(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    
    // 绑定请求（自动错误处理）
    if err := c.MustBindJSON(&req); err != nil {
        return err  // 自动返回 400 Bad Request
    }
    
    // 业务逻辑
    if ctrl.userService.EmailExists(req.Email) {
        return c.Conflict("邮箱已被注册")
    }
    
    user := ctrl.userService.Create(req)
    return c.Created(user)  // 201 Created
}
```

### 响应方法一览

| 方法 | HTTP 状态码 | 说明 |
|------|-------------|------|
| `c.Ok(data)` | 200 | 成功响应 |
| `c.Created(data)` | 201 | 创建成功 |
| `c.NoContent()` | 204 | 无内容 |
| `c.BadRequest(msg)` | 400 | 请求错误 |
| `c.Unauthorized(msg)` | 401 | 未授权 |
| `c.Forbidden(msg)` | 403 | 禁止访问 |
| `c.NotFound(msg)` | 404 | 未找到 |
| `c.Conflict(msg)` | 409 | 资源冲突 |
| `c.InternalError(msg)` | 500 | 服务器错误 |

### 统一响应格式

```json
// 成功
{
  "success": true,
  "data": { "id": 1, "name": "Alice" }
}

// 错误
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "用户不存在"
  }
}
```

### 完整控制器示例

```go
type UserController struct {
    userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
    return &UserController{userService: userService}
}

func (ctrl *UserController) MapRoutes(app *web.WebApplication) {
    users := app.MapGroup("/api/users")
    users.WithTags("Users")
    
    users.MapGet("", ctrl.GetAll).WithSummary("获取所有用户")
    users.MapGet("/:id", ctrl.GetByID).WithSummary("根据 ID 获取用户")
    users.MapPost("", ctrl.Create).WithSummary("创建用户")
    users.MapPut("/:id", ctrl.Update).WithSummary("更新用户")
    users.MapDelete("/:id", ctrl.Delete).WithSummary("删除用户")
}

func (ctrl *UserController) GetAll(c *web.HttpContext) web.IActionResult {
    page := c.QueryInt("page", 1)
    size := c.QueryInt("size", 10)
    
    users := ctrl.userService.ListUsers(page, size)
    return c.Ok(users)
}

func (ctrl *UserController) GetByID(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    user := ctrl.userService.GetUser(id)
    if user == nil {
        return c.NotFound("用户不存在")
    }
    
    return c.Ok(user)
}

func (ctrl *UserController) Create(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    if ctrl.userService.EmailExists(req.Email) {
        return c.Conflict("邮箱已被注册")
    }
    
    user := ctrl.userService.Create(req.Name, req.Email)
    return c.Created(user)
}

func (ctrl *UserController) Update(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    var req UpdateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    user := ctrl.userService.Update(id, req)
    if user == nil {
        return c.NotFound("用户不存在")
    }
    
    return c.Ok(user)
}

func (ctrl *UserController) Delete(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    if !ctrl.userService.Delete(id) {
        return c.NotFound("用户不存在")
    }
    
    return c.NoContent()
}
```

## 控制器注册

### 方式 1: AddController（推荐）

使用泛型函数自动注册，支持 `app.MapControllers()` 自动发现：

```go
// 注册单个控制器
web.AddController(services, func(sp di.IServiceProvider) *UserController {
    return NewUserController(di.GetRequiredService[UserService](sp))
})

// 自动映射所有注册的控制器
app.MapControllers()
```

### 方式 2: AddControllerInstance

适合需要更灵活控制的场景：

```go
web.AddControllerInstance(services, func(sp di.IServiceProvider) web.IController {
    userService := di.GetRequiredService[UserService](sp)
    return NewUserController(userService)
})
```

### 方式 3: 手动注册

不使用自动发现，手动调用 `MapRoutes`：

```go
app := builder.Build()

// 手动创建和注册控制器
userCtrl := controllers.NewUserController(app.Services)
userCtrl.MapRoutes(app)

orderCtrl := controllers.NewOrderController(app.Services)
orderCtrl.MapRoutes(app)
```

## 路由定义

### RESTful 路由

```go
func (ctrl *UserController) MapRoutes(app *web.WebApplication) {
    users := app.MapGroup("/api/users")
    users.WithTags("Users")
    
    // GET    /api/users          - 列表
    users.MapGet("", ctrl.List)
    
    // GET    /api/users/:id      - 详情
    users.MapGet("/:id", ctrl.GetByID)
    
    // POST   /api/users          - 创建
    users.MapPost("", ctrl.Create)
    
    // PUT    /api/users/:id      - 更新
    users.MapPut("/:id", ctrl.Update)
    
    // DELETE /api/users/:id      - 删除
    users.MapDelete("/:id", ctrl.Delete)
    
    // PATCH  /api/users/:id      - 部分更新
    users.MapPatch("/:id", ctrl.Patch)
}
```

### 嵌套路由

```go
func (ctrl *OrderController) MapRoutes(app *web.WebApplication) {
    orders := app.MapGroup("/api/orders")
    
    orders.MapGet("", ctrl.List)
    orders.MapGet("/:id", ctrl.GetByID)
    orders.MapPost("", ctrl.Create)
    
    // 嵌套资源: /api/orders/:orderId/items
    items := orders.MapGroup("/:orderId/items")
    items.MapGet("", ctrl.ListItems)
    items.MapPost("", ctrl.AddItem)
}
```

### Swagger 文档

为路由添加 API 文档：

```go
func (ctrl *UserController) MapRoutes(app *web.WebApplication) {
    users := app.MapGroup("/api/users")
    users.WithTags("Users")
    
    users.MapGet("", ctrl.List).
        WithSummary("获取所有用户").
        WithDescription("返回系统中所有用户的列表")
    
    users.MapGet("/:id", ctrl.GetByID).
        WithSummary("根据 ID 获取用户").
        WithDescription("根据用户 ID 返回单个用户信息")
    
    users.MapPost("", ctrl.Create).
        WithSummary("创建用户").
        WithDescription("创建一个新用户")
}
```

## 项目结构

### 推荐结构

```
myapp/
├── main.go                           # 主程序
├── controllers/                      # 控制器层
│   ├── user_controller.go
│   ├── order_controller.go
│   ├── product_controller.go
│   └── controller_extensions.go      # 控制器注册
├── services/                         # 服务层（业务逻辑）
│   ├── user_service.go
│   ├── order_service.go
│   └── service_extensions.go
├── models/                           # 数据模型
│   ├── user.go
│   └── order.go
└── repositories/                     # 数据访问层（可选）
    ├── user_repository.go
    └── order_repository.go
```

### 控制器扩展文件

```go
// controllers/controller_extensions.go
package controllers

import (
    "myapp/services"
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

// AddControllers 注册所有控制器到 DI 容器
func AddControllers(svc di.IServiceCollection) {
    web.AddController(svc, func(sp di.IServiceProvider) *UserController {
        return NewUserController(di.GetRequiredService[services.UserService](sp))
    })
    
    web.AddController(svc, func(sp di.IServiceProvider) *OrderController {
        return NewOrderController(di.GetRequiredService[services.OrderService](sp))
    })
    
    web.AddController(svc, func(sp di.IServiceProvider) *ProductController {
        return NewProductController(di.GetRequiredService[services.ProductService](sp))
    })
}
```

## 最佳实践

### 1. 单一职责原则

每个控制器只负责一个资源或功能模块：

```go
// ✅ 推荐
type UserController struct { ... }
type OrderController struct { ... }
type ProductController struct { ... }

// ❌ 不推荐
type ApiController struct { ... }  // 太宽泛
```

### 2. 使用 ActionResult

统一使用 `IActionResult` 返回响应：

```go
// ✅ 推荐：使用 ActionResult
func (ctrl *UserController) GetByID(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    user := ctrl.userService.GetUser(id)
    if user == nil {
        return c.NotFound("用户不存在")
    }
    
    return c.Ok(user)
}

// ❌ 不推荐：直接操作 gin.Context
func (ctrl *UserController) GetByID(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    user := ctrl.userService.GetUser(id)
    if user == nil {
        c.JSON(404, gin.H{"error": "用户不存在"})
        return
    }
    c.JSON(200, user)
}
```

### 3. 依赖注入

通过构造函数注入依赖：

```go
// ✅ 推荐
func NewUserController(
    userService services.UserService,
    logger logging.ILogger,
) *UserController {
    return &UserController{
        userService: userService,
        logger:      logger,
    }
}

// ❌ 不推荐
func NewUserController() *UserController {
    return &UserController{
        userService: services.NewUserService(),  // 硬编码依赖
    }
}
```

### 4. 请求验证

使用 `MustBindJSON` 自动处理验证错误：

```go
type CreateUserRequest struct {
    Name  string `json:"name" binding:"required,min=2,max=50"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=150"`
}

func (ctrl *UserController) Create(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err  // 自动返回 400 Bad Request
    }
    
    user := ctrl.userService.Create(req)
    return c.Created(user)
}
```

### 5. 错误处理

善用 `MustPathInt` 等方法简化错误处理：

```go
func (ctrl *UserController) GetByID(c *web.HttpContext) web.IActionResult {
    // 一行代码搞定参数获取和错误处理
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    user := ctrl.userService.GetUser(id)
    if user == nil {
        return c.NotFound("用户不存在")
    }
    
    return c.Ok(user)
}
```

## 与 .NET 对比

### .NET Controller

```csharp
[ApiController]
[Route("api/[controller]")]
public class UserController : ControllerBase
{
    private readonly IUserService _userService;
    
    public UserController(IUserService userService)
    {
        _userService = userService;
    }
    
    [HttpGet]
    public IActionResult GetUsers() => Ok(_userService.ListUsers());
    
    [HttpGet("{id}")]
    public IActionResult GetUser(int id)
    {
        var user = _userService.GetUser(id);
        return user == null ? NotFound() : Ok(user);
    }
    
    [HttpPost]
    public IActionResult CreateUser([FromBody] User user)
    {
        _userService.CreateUser(user);
        return CreatedAtAction(nameof(GetUser), new { id = user.Id }, user);
    }
}

// Program.cs
builder.Services.AddControllers();
var app = builder.Build();
app.MapControllers();
```

### CSGO Controller

```go
type UserController struct {
    userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
    return &UserController{userService: userService}
}

func (ctrl *UserController) MapRoutes(app *web.WebApplication) {
    users := app.MapGroup("/api/users")
    users.MapGet("", ctrl.GetAll)
    users.MapGet("/:id", ctrl.GetByID)
    users.MapPost("", ctrl.Create)
}

func (ctrl *UserController) GetAll(c *web.HttpContext) web.IActionResult {
    return c.Ok(ctrl.userService.ListUsers())
}

func (ctrl *UserController) GetByID(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    user := ctrl.userService.GetUser(id)
    if user == nil {
        return c.NotFound("用户不存在")
    }
    
    return c.Ok(user)
}

func (ctrl *UserController) Create(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    user := ctrl.userService.Create(req)
    return c.Created(user)
}

// main.go
web.AddController(builder.Services, NewUserController)
app := builder.Build()
app.MapControllers()
```

### 对比总结

| 特性 | .NET | CSGO | 一致性 |
|------|------|-----|--------|
| 控制器类 | `class UserController` | `type UserController struct` | ✅ |
| 构造函数注入 | ✅ 自动 | ✅ 手动（构造函数） | ✅ |
| 路由注册 | `[HttpGet]` 特性 | `MapGet()` 显式注册 | ✅ |
| HTTP 方法 | `[HttpGet]`, `[HttpPost]` | `MapGet`, `MapPost` | ✅ |
| IActionResult | `IActionResult` | `web.IActionResult` | ✅ |
| HttpContext | `HttpContext` | `web.HttpContext` | ✅ |
| Ok/NotFound | `Ok()`, `NotFound()` | `c.Ok()`, `c.NotFound()` | ✅ |
| 参数绑定 | ✅ 自动 | ✅ `MustBindJSON` | ✅ |
| 控制器注册 | `AddControllers()` | `AddController()` | ✅ |
| 路由映射 | `MapControllers()` | `MapControllers()` | ✅ |

## 完整示例

查看 `examples/controller_api_demo` 获取完整的可运行示例：

```bash
cd examples/controller_api_demo
go run main.go
```

访问：
- **API**: http://localhost:8080
- **Swagger**: http://localhost:8080/swagger

## 相关资源

- [Web 应用指南](web-applications.md) - Web 应用完整指南
- [依赖注入](dependency-injection.md) - DI 系统详解
- [API 文档](api-documentation.md) - Swagger 集成
- [业务模块](business-modules.md) - 模块化设计

---

**下一步**: 查看 [Web 应用指南](web-applications.md) 了解完整的 Web 应用开发流程。
