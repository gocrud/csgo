# CSGO - C# 风格的 Go Web 框架

<div align="center">
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Documentation](https://img.shields.io/badge/docs-latest-brightgreen)](docs/)

**CSGO** 是一个受 .NET/ASP.NET Core 启发的现代化 Go Web 框架，提供优雅的 API 和完整的企业级特性。

## ✨ 特性

- 🚀 **简洁优雅** - 受 ASP.NET Core 启发的 API 设计，上手即用
- 💉 **依赖注入** - 内置强大的 DI 容器，支持自动依赖解析
- 🎯 **类型安全** - 充分利用 Go 泛型，提供类型安全的 API
- 📝 **请求验证** - FluentValidation 风格的验证系统
- 🔧 **配置管理** - 灵活的配置系统，支持多种配置源
- 📊 **日志系统** - 结构化日志，支持多种输出格式
- 🛠️ **中间件** - 强大的中间件管道，灵活扩展请求处理流程
- 🎮 **控制器** - 可选的控制器模式，更好地组织代码
- 🔒 **错误处理** - 统一的错误处理和业务错误管理
- 🌐 **CORS 支持** - 开箱即用的跨域资源共享
- 📚 **API 文档** - 集成 Swagger/OpenAPI 支持
- 🔄 **后台服务** - 托管服务支持，轻松实现后台任务
- ⚡ **高性能** - 基于 Gin 框架，性能出色

## 🚀 快速开始

### 安装

```bash
go get github.com/gocrud/csgo
```

### Hello World

```go
package main

import "github.com/gocrud/csgo/web"

func main() {
    // 创建应用构建器
    builder := web.CreateBuilder()
    
    // 构建应用
    app := builder.Build()
    
    // 定义路由
    app.MapGet("/", func(c *web.HttpContext) web.IActionResult {
        return c.Ok(web.M{"message": "Hello, CSGO!"})
    })
    
    // 运行应用
    app.Run()  // 默认监听 :8080
}
```

运行应用：

```bash
go run main.go
```

访问 http://localhost:8080/ 查看结果。

### 完整示例

```go
package main

import (
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

// 定义服务
type UserService struct{}

func NewUserService() *UserService {
    return &UserService{}
}

func (s *UserService) GetUser(id int) string {
    return fmt.Sprintf("User %d", id)
}

// 定义请求模型
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    builder := web.CreateBuilder()
    
    // 注册服务
    builder.Services.Add(NewUserService)
    
    app := builder.Build()
    
    // 获取用户
    app.MapGet("/users/:id", func(c *web.HttpContext) web.IActionResult {
        userService := di.Get[*UserService](c.Services)
        id := c.Params().PathInt("id").Value()
        user := userService.GetUser(id)
        return c.Ok(web.M{"user": user})
    })
    
    // 创建用户
    app.MapPost("/users", func(c *web.HttpContext) web.IActionResult {
        var req CreateUserRequest
        if err := c.MustBindJSON(&req); err != nil {
            return err
        }
        return c.Created(web.M{"message": "User created", "name": req.Name})
    })
    
    app.Run()
}
```

## 📚 文档

完整文档请查看 [docs](docs/) 目录：

### 🎓 入门教程

- **[快速入门](docs/00-getting-started/)** - 30 分钟上手 CSGO
  - [安装配置](docs/00-getting-started/installation.md)
  - [第一个应用](docs/00-getting-started/hello-world.md)
  - [核心概念](docs/00-getting-started/concepts.md)

### 📖 核心基础

- **[基础知识](docs/01-fundamentals/)** - 深入理解核心概念
  - [Web 应用基础](docs/01-fundamentals/web-basics.md)
  - [路由系统](docs/01-fundamentals/routing.md)
  - [依赖注入](docs/01-fundamentals/dependency-injection.md)
  - [配置管理](docs/01-fundamentals/configuration.md)
  - [HttpContext](docs/01-fundamentals/http-context.md)

### 🔨 构建 API

- **[API 开发](docs/02-building-apis/)** - 构建生产级 API
  - [控制器模式](docs/02-building-apis/controllers.md)
  - [请求验证](docs/02-building-apis/validation.md)
  - [错误处理](docs/02-building-apis/error-handling.md)
  - [API 文档](docs/02-building-apis/api-docs.md)
  - [最佳实践](docs/02-building-apis/best-practices.md)

### 🚀 高级特性

- **[进阶主题](docs/03-advanced-features/)** - 掌握高级功能
  - [中间件](docs/03-advanced-features/middleware.md)
  - [后台服务](docs/03-advanced-features/background-services.md)
  - [日志系统](docs/03-advanced-features/logging.md)
  - [性能优化](docs/03-advanced-features/performance.md)
  - [单元测试](docs/03-advanced-features/testing.md)

### 🔧 模块文档

- **[Web 框架](web/README.md)** - Web 应用开发
- **[依赖注入 (DI)](di/README.md)** - 服务容器和依赖注入
- **[配置系统](configuration/README.md)** - 配置管理
- **[验证系统](validation/README.md)** - FluentValidation 风格的验证
- **[日志系统](logging/README.md)** - 结构化日志
- **[错误处理](errors/README.md)** - 业务错误和错误码
- **[主机托管](hosting/README.md)** - 应用生命周期和后台服务
- **[Swagger](swagger/README.md)** - API 文档生成

## 🏗️ 项目结构

推荐的项目结构：

```
myapp/
├── main.go                 # 应用入口
├── appsettings.json        # 配置文件
├── appsettings.Development.json
├── go.mod
├── go.sum
├── controllers/            # 控制器
│   ├── user_controller.go
│   └── product_controller.go
├── services/              # 业务服务
│   ├── user_service.go
│   └── product_service.go
├── models/                # 数据模型
│   ├── user.go
│   └── product.go
├── repositories/          # 数据访问层
│   ├── user_repository.go
│   └── product_repository.go
└── validators/            # 验证器
    ├── user_validator.go
    └── product_validator.go
```

## 🌟 核心概念

### WebApplicationBuilder

`WebApplicationBuilder` 是应用程序的构建器，负责配置和初始化：

```go
builder := web.CreateBuilder()

// 配置服务
builder.Services.Add(NewUserService)

// 访问配置
port := builder.Configuration.GetInt("server:port", 8080)

// 访问环境
if builder.Environment.IsDevelopment() {
    // 开发环境特定配置
}

// 构建应用
app := builder.Build()
```

### 依赖注入

内置的 DI 容器支持自动依赖解析：

```go
// 注册服务
builder.Services.Add(NewDatabase)
builder.Services.Add(NewUserRepository)  // 自动注入 Database
builder.Services.Add(NewUserService)     // 自动注入 UserRepository

// 使用服务
app.MapGet("/users", func(c *web.HttpContext) web.IActionResult {
    userService := di.Get[*UserService](c.Services)
    users := userService.GetAll()
    return c.Ok(users)
})
```

### HttpContext 和 ActionResult

统一的请求处理和响应格式：

```go
func handler(c *web.HttpContext) web.IActionResult {
    // 访问请求
    id := c.RawCtx().Param("id")
    
    // 绑定 JSON
    var req Request
    if err := c.MustBindJSON(&req); err != nil {
        return err  // 自动返回 400 错误
    }
    
    // 访问服务
    service := di.Get[*Service](c.Services)
    
    // 返回响应
    return c.Ok(data)           // 200 OK
    return c.Created(data)      // 201 Created
    return c.NoContent()        // 204 No Content
    return c.BadRequest("...")  // 400 Bad Request
    return c.NotFound("...")    // 404 Not Found
}
```

### 请求验证

FluentValidation 风格的验证系统：

```go
// 定义验证器
func NewCreateUserValidator() *validation.AbstractValidator[CreateUserRequest] {
    v := validation.NewValidator[CreateUserRequest]()
    
    v.Field(func(r *CreateUserRequest) string { return r.Name }).
        NotEmpty().
        MinLength(2).
        MaxLength(50)
    
    v.Field(func(r *CreateUserRequest) string { return r.Email }).
        NotEmpty().
        EmailAddress()
    
    return v
}

// 注册验证器
func init() {
    validation.RegisterValidator[CreateUserRequest](NewCreateUserValidator())
}

// 使用验证
func createUser(c *web.HttpContext) web.IActionResult {
    req, err := web.BindAndValidate[CreateUserRequest](c)
    if err != nil {
        return err  // 自动返回验证错误
    }
    // 验证通过，处理业务逻辑
    return c.Created(user)
}
```

### 控制器模式

可选的控制器模式，更好地组织代码：

```go
type UserController struct {
    userService *UserService
}

func NewUserController(userService *UserService) *UserController {
    return &UserController{userService: userService}
}

func (ctrl *UserController) MapRoutes(app *web.WebApplication) {
    users := app.MapGroup("/api/users")
    users.MapGet("", ctrl.List)
    users.MapGet("/:id", ctrl.Get)
    users.MapPost("", ctrl.Create)
    users.MapPut("/:id", ctrl.Update)
    users.MapDelete("/:id", ctrl.Delete)
}

func (ctrl *UserController) List(c *web.HttpContext) web.IActionResult {
    users := ctrl.userService.GetAll()
    return c.Ok(users)
}

// 注册控制器
web.AddController(builder.Services, NewUserController)
app.MapControllers()
```

### 后台服务

轻松实现后台任务：

```go
type EmailWorker struct {
    *hosting.BackgroundService
    emailService *EmailService
}

func NewEmailWorker(emailService *EmailService) *EmailWorker {
    worker := &EmailWorker{
        BackgroundService: hosting.NewBackgroundService(),
        emailService:      emailService,
    }
    worker.SetExecuteFunc(worker.execute)
    return worker
}

func (w *EmailWorker) execute(ctx context.Context) error {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            w.emailService.ProcessQueue()
        case <-w.StoppingToken():
            return nil
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

// 注册后台服务
builder.Services.AddHostedService(NewEmailWorker)
```

## 🎯 设计原则

CSGO 遵循以下设计原则：

1. **约定优于配置** - 提供合理的默认值，减少配置工作
2. **类型安全** - 使用 Go 泛型提供类型安全的 API
3. **依赖注入** - 松耦合、可测试的代码
4. **清晰的职责分离** - Controller → Service → Repository
5. **统一的错误处理** - 一致的错误响应格式
6. **开发者体验** - 简洁、直观、易于使用的 API

## 🤝 与 .NET 的对比

| 功能 | .NET/ASP.NET Core | CSGO |
|------|-------------------|------|
| 应用构建器 | `WebApplication.CreateBuilder()` | `web.CreateBuilder()` |
| 依赖注入 | `services.AddSingleton<T>()` | `services.Add(NewT)` |
| 路由 | `app.MapGet("/api/users", ...)` | `app.MapGet("/api/users", ...)` |
| 控制器 | `[ApiController]` | `web.AddController()` |
| 请求验证 | `FluentValidation` | `validation.NewValidator[T]()` |
| 后台服务 | `IHostedService` | `hosting.IHostedService` |
| 配置 | `IConfiguration` | `configuration.IConfiguration` |
| 日志 | `ILogger<T>` | `logging.ILogger` |

## 📦 依赖

CSGO 基于以下优秀的开源项目：

- [Gin](https://github.com/gin-gonic/gin) - 高性能的 HTTP Web 框架
- [Zerolog](https://github.com/rs/zerolog) - 零分配的 JSON 日志库

## 🗺️ 路线图

- [x] Web 框架基础
- [x] 依赖注入
- [x] 配置管理
- [x] 请求验证
- [x] 日志系统
- [x] 错误处理
- [x] 后台服务
- [x] Swagger 集成
- [ ] 数据库集成 (GORM)
- [ ] 认证授权 (JWT)
- [ ] 缓存支持 (Redis)
- [ ] 消息队列支持
- [ ] 健康检查
- [ ] 限流和熔断
- [ ] 分布式追踪

## 🤝 贡献

欢迎贡献！请查看 [贡献指南](CONTRIBUTING.md) 了解详情。

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 🙏 鸣谢

感谢所有贡献者和以下项目的启发：

- [ASP.NET Core](https://github.com/dotnet/aspnetcore) - 现代 Web 框架的典范
- [Gin](https://github.com/gin-gonic/gin) - 高性能的 Go Web 框架
- [Echo](https://github.com/labstack/echo) - 极简的 Go Web 框架

## 📮 联系方式

- 问题反馈：[GitHub Issues](https://github.com/gocrud/csgo/issues)
- 讨论交流：[GitHub Discussions](https://github.com/gocrud/csgo/discussions)

---

<div align="center">

**[快速开始](docs/00-getting-started/)** | **[完整文档](docs/)** | **[示例项目](examples/)**

</div>

