# CSGO 框架快速参考

> 一页纸速查手册 📄

---

## 📦 安装

```bash
go get github.com/gocrud/csgo
```

---

## 🚀 基础应用

```go
package main

import (
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

func main() {
    // 1. 创建 Builder
    builder := web.CreateBuilder()
    
    // 2. 注册服务
    builder.Services.AddSingleton(NewUserService)
    
    // 3. 构建应用
    app := builder.Build()
    
    // 4. 定义路由
    app.MapGet("/hello", func(c *web.HttpContext) web.IActionResult {
        svc := di.GetRequiredService[*UserService](app.Services)
        return c.Ok(svc.GetGreeting())
    })
    
    // 5. 运行
    app.Run()
}
```

---

## 💉 依赖注入

### 注册服务

```go
// Singleton - 全局唯一（推荐用于无状态服务）
services.AddSingleton(NewDatabaseConnection)

// Transient - 每次创建新实例
services.AddTransient(NewEmailService)

// 命名服务
services.AddKeyedSingleton("primary", NewPrimaryDb)
services.AddKeyedTransient("logger", NewLogger)
```

### 解析服务

```go
// ✅ 推荐：泛型辅助函数
userService := di.GetRequiredService[*UserService](provider)
cache := di.GetRequiredService[*Cache](provider)

// ✅ 可选：指针填充
var userService *UserService
provider.GetRequiredService(&userService)

// 可选服务
userService, err := di.GetService[*UserService](provider)
if err != nil {
    // 服务不存在
}

// 命名服务
primaryDb := di.GetRequiredKeyedService[*Database](provider, "primary")
```

---

## 🎮 Controller 模式

### 定义 Controller

```go
// controllers/user_controller.go
package controllers

import (
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

type UserController struct {
    userService *UserService
}

func NewUserController(app *web.WebApplication) *UserController {
    return &UserController{
        userService: di.GetRequiredService[*UserService](app.Services),
    }
}

// 实现 IController 接口
func (c *UserController) MapRoutes(app *web.WebApplication) {
    users := app.MapGroup("/api/users")
    users.MapGet("", c.GetAll)
    users.MapGet("/:id", c.GetByID)
    users.MapPost("", c.Create)
    users.MapPut("/:id", c.Update)
    users.MapDelete("/:id", c.Delete)
}

// 处理器使用 HttpContext + ActionResult
func (c *UserController) GetByID(ctx *web.HttpContext) web.IActionResult {
    id, err := ctx.PathInt("id")
    if err != nil {
        return ctx.BadRequest("无效的ID")
    }
    
    user := c.userService.GetUser(id)
    if user == nil {
        return ctx.NotFound("用户不存在")
    }
    
    return ctx.Ok(user)
}

func (c *UserController) Create(ctx *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    if err := ctx.MustBindJSON(&req); err != nil {
        return err  // 自动返回 400
    }
    
    user := c.userService.Create(req)
    return ctx.Created(user)
}
```

### 注册 Controller

```go
// controllers/controller_extensions.go
func AddControllers(services di.IServiceCollection) {
    web.AddController(services, NewUserController)
    web.AddController(services, NewOrderController)
}

// main.go
func main() {
    builder := web.CreateBuilder()
    
    services.AddServices(builder.Services)
    controllers.AddControllers(builder.Services)
    
    app := builder.Build()
    app.MapControllers()  // 自动映射所有控制器
    app.Run()
}
```

---

## ⚙️ 配置管理

### 定义配置

```go
// config/settings.go
type AppSettings struct {
    AppName  string `json:"appName"`
    Port     int    `json:"port"`
}
```

### 注册配置

```go
// main.go
import "github.com/gocrud/csgo/configuration"

configuration.Configure[AppSettings](
    builder.Services,
    builder.Configuration,
    "App",  // appsettings.json 中的节点名
)
```

### 使用配置

```go
// 在 Controller 中
type UserController struct {
    settings *AppSettings
}

func NewUserController(app *web.WebApplication) *UserController {
    opts := di.GetRequiredService[configuration.IOptions[AppSettings]](app.Services)
    return &UserController{
        settings: opts.Value(),
    }
}

func (c *UserController) GetInfo(ctx *web.HttpContext) web.IActionResult {
    return ctx.Ok(gin.H{
        "app":  c.settings.AppName,
        "port": c.settings.Port,
    })
}
```

### appsettings.json

```json
{
  "App": {
    "appName": "My API",
    "port": 8080
  }
}
```

---

## 🌐 路由和响应

### HttpContext 方法

```go
// 路径参数
id, err := ctx.PathInt("id")
userId := ctx.Param("userId")

// 查询参数
page := ctx.QueryInt("page", 1)      // 带默认值
name := ctx.Query("name")

// 请求体绑定
var req CreateUserRequest
if err := ctx.MustBindJSON(&req); err != nil {
    return err  // 自动返回 400
}
```

### ActionResult 响应

```go
// 成功响应
ctx.Ok(data)              // 200 OK
ctx.Created(data)         // 201 Created
ctx.NoContent()           // 204 No Content

// 错误响应
ctx.BadRequest("错误")     // 400
ctx.Unauthorized("未授权") // 401
ctx.Forbidden("禁止")      // 403
ctx.NotFound("未找到")     // 404
ctx.Conflict("冲突")       // 409
ctx.InternalError("错误")  // 500
```

### 路由定义

```go
// 基础路由
app.MapGet("/users", GetUsers)
app.MapPost("/users", CreateUser)
app.MapPut("/users/:id", UpdateUser)
app.MapDelete("/users/:id", DeleteUser)

// 路由组
api := app.MapGroup("/api")
api.MapGet("/users", GetUsers)
api.MapGet("/orders", GetOrders)

// 嵌套组
v1 := api.MapGroup("/v1")
v1.MapGet("/users", GetUsersV1)
```

---

## 📖 Swagger 集成

```go
import (
    "github.com/gocrud/csgo/openapi"
    "github.com/gocrud/csgo/swagger"
)

func main() {
    builder := web.CreateBuilder()
    
    // 配置 Swagger
    swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
        opts.Title = "My API"
        opts.Version = "v1"
        opts.Description = "API 文档"
    })
    
    app := builder.Build()
    
    // 启用 Swagger
    swagger.UseSwagger(app)
    swagger.UseSwaggerUI(app)
    
    // 为路由添加文档
    app.MapGet("/users", GetUsers).
        WithSummary("获取所有用户").
        WithDescription("返回系统中所有用户")
    
    app.Run()
}
```

访问 Swagger UI: http://localhost:8080/swagger

---

## 🎯 Controller 最佳实践

### ✅ 推荐

```go
type UserController struct {
    // ✅ 依赖服务（不可变）
    userService *UserService
    
    // ✅ 配置（不可变）
    settings *AppSettings
    
    // ✅ WebApplication（动态解析）
    app *web.WebApplication
}

func (c *UserController) GetUser(ctx *web.HttpContext) web.IActionResult {
    // ✅ 从请求获取数据
    id, _ := ctx.PathInt("id")
    
    // ✅ 使用注入的服务
    user := c.userService.GetUser(id)
    
    return ctx.Ok(user)
}
```

### ❌ 避免

```go
type BadController struct {
    // ❌ 请求状态（会被覆盖）
    currentUser *User
    
    // ❌ 请求ID（线程不安全）
    requestID string
}
```

**原因：** Controllers 是单例的，会被多个请求并发访问！

---

## 📊 生命周期对比

| 生命周期 | 创建时机 | 适用场景 |
|---------|---------|---------|
| **Singleton** | 应用启动时 | 数据库连接、缓存、配置、无状态服务 |
| **Transient** | 每次请求时 | 有状态服务、轻量级操作、请求日志 |

**注意：** 框架不支持 Scoped 生命周期。

---

## 🔗 常用导入

```go
import (
    "github.com/gocrud/csgo/configuration"
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/hosting"
    "github.com/gocrud/csgo/openapi"
    "github.com/gocrud/csgo/swagger"
    "github.com/gocrud/csgo/web"
)
```

---

## 📚 完整示例

```go
package main

import (
    "myapp/config"
    "myapp/controllers"
    "myapp/services"
    
    "github.com/gocrud/csgo/configuration"
    "github.com/gocrud/csgo/swagger"
    "github.com/gocrud/csgo/web"
)

func main() {
    // 1. 创建 Builder
    builder := web.CreateBuilder()
    
    // 2. 配置
    configuration.Configure[config.AppSettings](
        builder.Services,
        builder.Configuration,
        "App",
    )
    
    // 3. 注册服务
    services.AddServices(builder.Services)
    
    // 4. 注册控制器
    controllers.AddControllers(builder.Services)
    
    // 5. Swagger
    swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
        opts.Title = "My API"
        opts.Version = "v1"
    })
    
    // 6. 构建应用
    app := builder.Build()
    
    // 7. 中间件
    swagger.UseSwagger(app)
    swagger.UseSwaggerUI(app)
    
    // 8. 映射控制器
    app.MapControllers()
    
    // 9. 运行
    app.Run()
}
```

---

## 🆘 常见问题

### Q: 为什么没有 Scoped？
A: 为了简化设计和提升性能。Controllers 是单例的，符合 Go 生态习惯。

### Q: Controller 如何处理请求状态？
A: 从 `HttpContext` 参数获取，不要存储在 Controller 字段中。

### Q: 如何实现请求级别的服务？
A: 注册为 Transient，在 handler 中动态获取：
```go
logger := di.GetRequiredService[*RequestLogger](app.Services)
```

### Q: 配置如何热更新？
A: 使用 `IOptionsMonitor[T]` 而不是 `IOptions[T]`。

---

## 📖 更多文档

- [完整 API 参考](api/)
- [详细指南](guides/)
- [框架变更说明](FRAMEWORK_CHANGES.md)
- [示例代码](../examples/)

---

**Happy Coding! 🎉**

