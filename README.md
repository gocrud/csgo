# CSGO

### C# 风格的现代化 Go Web 框架

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue?style=flat-square)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/gocrud/csgo.svg)](https://pkg.go.dev/github.com/gocrud/csgo)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)
[![Documentation](https://img.shields.io/badge/docs-latest-brightgreen?style=flat-square)](docs/)

**受 ASP.NET Core 启发 • 优雅的 API 设计 • 企业级开发体验**

[快速开始](#-快速开始) • [文档](docs/) • [示例](#-完整示例) • [更新日志](docs/CHANGELOG.md)

---

## ✨ 核心特性

### 🚀 开发效率
- **依赖注入** - 自动依赖解析，松耦合设计
- **类型安全** - 利用 Go 泛型，编译时检查
- **开箱即用** - 最小配置即可启动

### 🛠️ 企业级功能
- **配置管理** - 多源配置，环境感知
- **请求验证** - FluentValidation 风格
- **后台服务** - 托管服务支持

### 🎯 Web 开发
- **路由系统** - 灵活的 RESTful 路由
- **中间件** - 强大的中间件管道
- **控制器** - 可选的 MVC 模式
- **错误处理** - 统一的异常管理

### 📚 工具支持
- **Swagger/OpenAPI** - 自动生成 API 文档
- **CORS 支持** - 开箱即用的跨域支持
- **参数绑定** - 自动解析请求参数
- **高性能** - 基于 Gin，性能卓越

---

## 🚀 快速开始

### 1️⃣ 安装

```bash
go get -u github.com/gocrud/csgo@latest
```

### 2️⃣ Hello World

创建 `main.go`：

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

### 3️⃣ 运行

```bash
go run main.go
```

访问 http://localhost:8080/ 查看结果 🎉

---

## 📝 完整示例

一个带依赖注入和 RESTful API 的完整示例：

```go
package main

import (
    "fmt"
    
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

// 1. 定义服务层
type UserService struct{}

func NewUserService() *UserService {
    return &UserService{}
}

func (s *UserService) GetUser(id int) string {
    return fmt.Sprintf("User #%d", id)
}

func (s *UserService) CreateUser(name, email string) map[string]any {
    return map[string]any{
        "id":    1,
        "name":  name,
        "email": email,
    }
}

// 2. 定义请求模型
type CreateUserRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

func main() {
    // 创建应用构建器
    builder := web.CreateBuilder()
    
    // 注册服务到 DI 容器
    builder.Services.Add(NewUserService)
    
    // 构建应用
    app := builder.Build()
    
    // 定义 API 路由组
    api := app.MapGroup("/api")
    {
        users := api.MapGroup("/users")
        
        // GET /api/users/:id - 获取用户
        users.MapGet("/:id", func(c *web.HttpContext) web.IActionResult {
            // 从 DI 容器获取服务
            userService := di.Get[*UserService](c.Services)
            
            // 获取路径参数
            id := c.Params().PathInt("id").Value()
            
            // 调用服务
            user := userService.GetUser(id)
            
            // 返回响应
            return c.Ok(web.M{"user": user})
        })
        
        // POST /api/users - 创建用户
        users.MapPost("", func(c *web.HttpContext) web.IActionResult {
            var req CreateUserRequest
            
            // 绑定并验证请求体
            if err := c.MustBindJSON(&req); err != nil {
                return err  // 自动返回 400 错误
            }
            
            // 从 DI 容器获取服务
            userService := di.Get[*UserService](c.Services)
            
            // 创建用户
            user := userService.CreateUser(req.Name, req.Email)
            
            // 返回 201 Created
            return c.Created(user)
        })
    }
    
    // 运行应用
    app.Run()
}
```

**测试 API**：

```bash
# 获取用户
curl http://localhost:8080/api/users/1

# 创建用户
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"张三","email":"zhangsan@example.com"}'
```

---

## 📚 文档导航

### 🎓 学习路径

#### [00-快速入门](docs/00-getting-started/) 
- [安装配置](docs/00-getting-started/installation.md)
- [第一个应用](docs/00-getting-started/hello-world.md)
- [核心概念](docs/00-getting-started/concepts.md)

#### [01-基础知识](docs/01-fundamentals/) 
- [Web 应用基础](docs/01-fundamentals/web-basics.md)
- [路由系统](docs/01-fundamentals/routing.md)
- [依赖注入](docs/01-fundamentals/dependency-injection.md)
- [配置管理](docs/01-fundamentals/configuration.md)
- [HttpContext](docs/01-fundamentals/http-context.md)
- [📦 项目实战](docs/01-fundamentals/project-simple-api.md)

#### [02-API 开发](docs/02-building-apis/) 
- [控制器模式](docs/02-building-apis/controllers.md)
- [请求验证](docs/02-building-apis/validation.md)
- [错误处理](docs/02-building-apis/error-handling.md)
- [API 文档](docs/02-building-apis/api-docs.md)
- [最佳实践](docs/02-building-apis/best-practices.md)
- [📦 项目实战](docs/02-building-apis/project-crud-api.md)

#### [03-高级特性](docs/03-advanced-features/) 
- [中间件](docs/03-advanced-features/middleware.md)
- [后台服务](docs/03-advanced-features/background-services.md)
- [日志系统](docs/03-advanced-features/logging.md)
- [性能优化](docs/03-advanced-features/performance.md)
- [单元测试](docs/03-advanced-features/testing.md)
- [📦 项目实战](docs/03-advanced-features/project-complete-app.md)

### 🔧 模块参考手册

| 模块 | 说明 | 文档 |
|------|------|------|
| 🌐 **Web** | Web 应用开发核心 | [web/README.md](web/README.md) |
| 💉 **DI** | 依赖注入容器 | [di/README.md](di/README.md) |
| ⚙️ **Configuration** | 配置管理系统 | [configuration/README.md](configuration/README.md) |
| 🏠 **Hosting** | 应用托管和后台服务 | [hosting/README.md](hosting/README.md) |
| 📚 **Swagger** | API 文档生成 | [swagger/README.md](swagger/README.md) |

## 

> 💡 **提示**：可以根据项目规模选择扁平或分层结构，CSGO 两种风格都支持。

---

## 🌟 核心概念速览

### 📦 WebApplicationBuilder - 应用构建器

`WebApplicationBuilder` 是应用程序的构建器，负责配置和初始化：

```go
builder := web.CreateBuilder()

// 注册服务到 DI 容器
builder.Services.Add(NewUserService)

// 读取配置
port := builder.Configuration.GetInt("server:port", 8080)

// 环境判断
if builder.Environment.IsDevelopment() {
    // 开发环境专属配置
}

// 构建应用
app := builder.Build()
```

### 💉 依赖注入 - 自动依赖解析

内置的 DI 容器支持构造函数自动注入：

```go
builder.Services.Add(NewDatabase)
builder.Services.Add(NewUserRepository)  // 自动注入 *Database
builder.Services.Add(NewUserService)     // 自动注入 *UserRepository

// 在处理器中使用
app.MapGet("/users", func(c *web.HttpContext) web.IActionResult {
    service := di.Get[*UserService](c.Services)
    return c.Ok(service.GetAll())
})
```

### 🎯 HttpContext & ActionResult - 统一请求处理

类型安全的请求处理和响应：

```go
func handler(c *web.HttpContext) web.IActionResult {
    // 获取路径参数
    id := c.Params().PathInt("id").Value()
    
    // 绑定请求体
    var req CreateRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err  // 自动返回 400 Bad Request
    }
    
    // 获取服务
    service := di.Get[*MyService](c.Services)
    
    // 统一的响应格式
    return c.Ok(data)              // 200 OK
    return c.Created(data)         // 201 Created
    return c.NoContent()           // 204 No Content
    return c.BadRequest("error")   // 400 Bad Request
    return c.NotFound("not found") // 404 Not Found
}
```

### 🎮 控制器模式 - MVC 风格开发

可选的控制器模式，更好地组织代码：

```go
type UserController struct {
    userService *UserService
}

func NewUserController(userService *UserService) *UserController {
    return &UserController{userService: userService}
}

// 实现 IController 接口
func (ctrl *UserController) MapRoutes(app *web.WebApplication) {
    users := app.MapGroup("/api/users")
    users.MapGet("", ctrl.List)
    users.MapGet("/:id", ctrl.Get)
    users.MapPost("", ctrl.Create)
}

func (ctrl *UserController) List(c *web.HttpContext) web.IActionResult {
    return c.Ok(ctrl.userService.GetAll())
}

// 注册控制器
web.AddController(builder.Services, NewUserController)
app.MapControllers()  // 自动注册所有控制器路由
```

### 🔄 后台服务 - 托管服务支持

轻松实现后台任务和定时任务：

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

### ⚙️ 配置管理 - 多源配置系统

支持 JSON 文件、环境变量、命令行参数：

```go
// appsettings.json
{
  "server": {
    "port": 8080,
    "host": "localhost"
  },
  "database": {
    "host": "localhost",
    "port": 5432
  }
}

// 读取配置
port := builder.Configuration.GetInt("server:port", 8080)
host := builder.Configuration.GetString("server:host", "localhost")

// 绑定到结构体
type ServerConfig struct {
    Port int    `json:"port"`
    Host string `json:"host"`
}

var cfg ServerConfig
builder.Configuration.Bind("server", &cfg)
```

---

## 🎯 设计哲学

CSGO 遵循以下核心设计原则：

### 🎨 开发体验优先
- 简洁直观的 API
- 最小化样板代码
- 约定优于配置
- IDE 友好

### 🛡️ 类型安全
- 充分利用 Go 泛型
- 编译时类型检查
- 减少运行时错误
- 自动补全支持

### 🏗️ 架构清晰
- 依赖注入为核心
- 明确的职责分离
- 松耦合易测试
- 渐进式学习曲线

---

## 🤝 与 ASP.NET Core 的对比

熟悉 .NET？快速上手 CSGO：

| 功能 | ASP.NET Core | CSGO | 说明 |
|------|-------------|------|------|
| **应用构建** | `WebApplication.CreateBuilder()` | `web.CreateBuilder()` | ✅ 几乎相同 |
| **依赖注入** | `services.AddSingleton<IService, Service>()` | `services.Add(NewService)` | 🔄 构造函数风格 |
| **路由定义** | `app.MapGet("/api/users", handler)` | `app.MapGet("/api/users", handler)` | ✅ 完全一致 |
| **控制器** | `[ApiController]` + 特性 | `web.AddController()` + 接口 | 🔄 接口风格 |
| **后台服务** | `IHostedService` | `hosting.IHostedService` | ✅ 概念相同 |
| **配置系统** | `IConfiguration` | `configuration.IConfiguration` | ✅ 相似 API |
| **响应类型** | `IActionResult` | `web.IActionResult` | ✅ 相同模式 |

---

## 🗺️ 开发路线图

### ✅ 已完成

- [x] Web 框架核心
- [x] 依赖注入系统
- [x] 配置管理
- [x] 请求验证
- [x] 错误处理
- [x] 后台服务
- [x] Swagger/OpenAPI 集成
- [x] 中间件管道
- [x] 控制器模式

### 🚧 开发中

- [ ] 更完善的示例项目
- [ ] 性能测试和基准
- [ ] 更多的中间件（认证、限流等）

### 📋 计划中

- [ ] 认证授权方案 (JWT)
- [ ] 指标收集 (Prometheus)
- [ ] 分布式追踪 (OpenTelemetry)

> 💡 有想法？欢迎在 [Discussions](https://github.com/gocrud/csgo/discussions) 中分享！

---

## 🤝 参与贡献

我们欢迎各种形式的贡献！

### 如何贡献

1. 🐛 **报告 Bug** - 在 [Issues](https://github.com/gocrud/csgo/issues) 中提交
2. 💡 **提出建议** - 在 [Discussions](https://github.com/gocrud/csgo/discussions) 中讨论
3. 📝 **改进文档** - 帮助完善文档
4. 🔧 **提交代码** - Fork 项目并提交 Pull Request

### 贡献者

感谢所有为 CSGO 做出贡献的开发者！

[![Contributors](https://contrib.rocks/image?repo=gocrud/csgo)](https://github.com/gocrud/csgo/graphs/contributors)

---

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

---

## 💬 社区与支持

### 💡 提问讨论
**[GitHub Discussions](https://github.com/gocrud/csgo/discussions)**  
适合：功能讨论、使用问题、最佳实践

### 🐛 问题反馈
**[GitHub Issues](https://github.com/gocrud/csgo/issues)**  
适合：Bug 报告、功能请求

### 📚 文档
**[在线文档](docs/)**  
适合：学习教程、API 参考

---

## ⭐ Star History

如果 CSGO 对你有帮助，请给个 Star 支持一下！

[![Star History Chart](https://api.star-history.com/svg?repos=gocrud/csgo&type=Date)](https://star-history.com/#gocrud/csgo&Date)

---

**[快速开始](#-快速开始)** • **[查看文档](docs/)** • **[参与贡献](#-参与贡献)**

Made with ❤️ by the CSGO community
