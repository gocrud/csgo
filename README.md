# CSGO Framework

一个受 ASP.NET Core 启发的 Go Web 框架，提供完整的依赖注入、控制器模式和现代化开发体验。

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## ✨ 核心特性

- 🎯 **完整的依赖注入** - 类似 .NET 的服务注册（`AddSingleton`、`AddTransient`），Go 风格的指针填充解析
- 🌐 **现代 Web 框架** - 基于 Gin，提供控制器模式、路由系统和中间件支持
- 🎭 **HttpContext & ActionResult** - 类似 .NET 的请求处理模式，统一响应格式（`Ok`、`NotFound`、`BadRequest`）
- 📦 **模块化设计** - 业务模块扩展方法，清晰的代码组织和依赖管理
- 📖 **Swagger 集成** - 自动 API 文档生成，支持 OpenAPI 3.0
- ⚙️ **配置管理** - 多源配置系统（JSON、环境变量、命令行）
- 🚀 **应用托管** - Host Builder 模式，完整的应用生命周期管理
- 🔧 **开发体验** - 类型安全、IDE 友好、简洁的 API 设计

## 🚀 快速开始

### 安装

```bash
go get github.com/gocrud/csgo
```

### 第一个应用

```go
package main

import (
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

func main() {
    // 创建应用构建器
    builder := web.CreateBuilder()
    
    // 注册服务
    builder.Services.AddSingleton(NewUserService)
    
    // 构建应用
    app := builder.Build()
    
    // 使用 HttpContext + ActionResult（推荐）
    app.MapGet("/hello", func(c *web.HttpContext) web.IActionResult {
        userService := di.GetRequiredService[*UserService](app.Services)
        return c.Ok(gin.H{"message": userService.GetGreeting()})
    })
    
    // 运行应用
    app.Run()
}

type UserService struct{}

func NewUserService() *UserService {
    return &UserService{}
}

func (s *UserService) GetGreeting() string {
    return "Hello from CSGO!"
}
```

运行应用：

```bash
go run main.go
```

访问 http://localhost:8080/hello，你会看到：

```json
{"message": "Hello from CSGO!"}
```

## 📚 核心概念

### 依赖注入

CSGO 提供了完整的 DI 容器，支持两种服务生命周期：

```go
// Singleton - 全局唯一实例（推荐用于无状态服务）
services.AddSingleton(NewDatabaseConnection)
services.AddSingleton(NewUserService)

// Transient - 每次请求都创建新实例（用于有状态服务）
services.AddTransient(NewEmailService)
services.AddTransient(NewRequestLogger)

// 服务解析（指针填充方式）
var db *DatabaseConnection
provider.GetRequiredService(&db)

// 或使用泛型辅助方法（推荐）
db := di.GetRequiredService[*DatabaseConnection](provider)
```

**注意：** 框架采用简化设计，不支持 Scoped 生命周期。Controllers 是单例的，必须保持无状态。

📖 [查看完整 DI 指南](docs/guides/dependency-injection.md) | [框架变更说明](docs/FRAMEWORK_CHANGES.md)

### Web 应用

基于 Gin 构建，提供控制器模式和路由系统：

```go
builder := web.CreateBuilder()

// 添加 CORS
builder.AddCors(func(opts *CorsOptions) {
    opts.AllowOrigins = []string{"http://localhost:3000"}
})

app := builder.Build()

// 使用中间件
app.UseCors()

// 定义路由组
api := app.MapGroup("/api")
api.MapGet("/users", GetUsers)
api.MapPost("/users", CreateUser)

app.Run()
```

[查看 Web 应用指南 →](docs/guides/web-applications.md)

### 控制器模式

类似 ASP.NET Core MVC 的控制器，支持 ActionResult：

```go
type UserController struct {
    userService *UserService
}

func NewUserController(userService *UserService) *UserController {
    return &UserController{userService: userService}
}

// 使用 IController 接口
func (ctrl *UserController) MapRoutes(app *web.WebApplication) {
    users := app.MapGroup("/api/users")
    users.MapGet("/:id", ctrl.GetByID)
    users.MapPost("", ctrl.Create)
}

// 使用 HttpContext + ActionResult
func (ctrl *UserController) GetByID(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err  // 自动返回 400 Bad Request
    }
    
    user := ctrl.userService.GetUserByID(id)
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

// 注册控制器
web.AddController(builder.Services, func(sp di.IServiceProvider) *UserController {
    return NewUserController(di.GetRequiredService[*UserService](sp))
})

app := builder.Build()
app.MapControllers()  // 自动映射所有控制器
```

[查看控制器指南 →](docs/guides/controllers.md)

### 业务模块

创建可复用的业务模块扩展：

```go
// 模块定义
package users

// AddUserServices 注册用户模块的所有服务
func AddUserServices(services di.IServiceCollection) {
    services.AddTransient(NewUserService)
    services.AddTransient(NewUserRepository)
    services.AddSingleton(NewUserCache)
}

// 在主程序中使用
builder := web.CreateBuilder()
users.AddUserServices(builder.Services)
orders.AddOrderServices(builder.Services)
```

[查看业务模块指南 →](docs/guides/business-modules.md)

### API 文档

自动生成 Swagger 文档：

```go
import "github.com/gocrud/csgo/swagger"

// 添加 Swagger
swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
    opts.Title = "My API"
    opts.Version = "v1"
    opts.Description = "API Documentation"
})

app := builder.Build()

// 启用 Swagger UI
swagger.UseSwagger(app)
swagger.UseSwaggerUI(app)

// 访问 http://localhost:8080/swagger
```

[查看 API 文档指南 →](docs/guides/api-documentation.md)

## 📖 完整文档

### 快速入门
- [快速开始](docs/getting-started.md) - 安装和第一个应用
- **[快速参考](docs/QUICK_REFERENCE.md)** - 一页纸速查手册 📄
- **[框架变更说明](docs/FRAMEWORK_CHANGES.md)** - 设计决策和最佳实践 🔄

### 用户指南
- [Web 应用](docs/guides/web-applications.md) - Web 应用完整指南
- [控制器](docs/guides/controllers.md) - 控制器模式
- [依赖注入](docs/guides/dependency-injection.md) - DI 系统
- [配置管理](docs/guides/configuration.md) - 配置系统
- [应用托管](docs/guides/hosting.md) - 生命周期管理
- [业务模块](docs/guides/business-modules.md) - 模块化设计
- [API 文档](docs/guides/api-documentation.md) - Swagger 集成

### 参考资料
- [API 参考](docs/api/) - 完整的 API 文档
- [最佳实践](docs/best-practices.md) - 推荐的代码组织和模式
- [与 .NET 对比](docs/comparison-with-dotnet.md) - API 对照和迁移指南

## 💡 示例

查看 [examples/](examples/) 目录获取完整的示例代码：

- [complete_di_demo](examples/complete_di_demo/) - DI 功能完整演示
- [business_module_demo](examples/business_module_demo/) - 业务模块设计示例
- [controller_api_demo](examples/controller_api_demo/) - 控制器模式示例
- [service_resolution_demo](examples/service_resolution_demo/) - 服务解析示例

## 🔄 与 .NET 的关系

CSGO 深受 ASP.NET Core 启发，但针对 Go 语言特性进行了优化：

| .NET | CSGO | 说明 |
|------|-----|------|
| `IServiceCollection` | `di.IServiceCollection` | 服务注册接口 |
| `AddSingleton<T>()` | `AddSingleton(factory)` | 注册单例服务 |
| `GetService<T>()` | `GetService(&target)` | 指针填充方式解析 |
| `WebApplicationBuilder` | `web.CreateBuilder()` | Web 应用构建器 |
| `app.MapGet()` | `app.MapGet()` | 路由定义 |
| `HttpContext` | `web.HttpContext` | HTTP 上下文 |
| `IActionResult` | `web.IActionResult` | 操作结果接口 |
| `Ok()` / `NotFound()` | `c.Ok()` / `c.NotFound()` | 响应辅助方法 |
| `IHostedService` | `IHostedService` | 后台服务 |

**关键差异**：
- **服务解析**：CSGO 使用 Go 惯用的指针填充方式（类似 `json.Unmarshal`），而不是泛型返回
- **类型安全**：编译时类型检查，无需类型断言
- **性能优化**：针对 Go 的 runtime 特性优化（如 `sync.Pool`、unsafe 指针等）

[查看详细对比 →](docs/comparison-with-dotnet.md)

## 🤝 贡献

欢迎贡献代码、报告问题或提出建议！

## 📄 许可证

MIT License

---

**Star ⭐ 如果这个项目对你有帮助！**
