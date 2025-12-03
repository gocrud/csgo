# CSGO 框架更新说明

> 最后更新：2024-12

本文档记录了框架的重要设计决策和与 .NET 的差异。

---

## 🔄 主要设计决策

### 1. 移除 Scoped 生命周期

**决策：** 框架只支持 **Singleton** 和 **Transient** 两种生命周期。

**原因：**
- ✅ 符合 Go 生态习惯（Gin/Echo 都是这样）
- ✅ 性能最优（无运行时作用域管理开销）
- ✅ 代码简单（无复杂的作用域生命周期）
- ✅ 避免线程安全陷阱

**影响：**

| 组件 | 生命周期 | 说明 |
|------|---------|------|
| **Controllers** | Singleton | 启动时创建一次，**必须无状态** |
| **配置 (IOptions)** | Singleton | 配置在启动时读取 |
| **业务服务** | Singleton 或 Transient | 由开发者选择 |
| **请求级服务** | Transient | 在 handler 中手动解析 |

**迁移指南：**

```go
// ❌ 旧代码（如果之前有）
services.AddScoped(NewUserService)
scope := provider.CreateScope()

// ✅ 新代码
services.AddTransient(NewUserService)  // 或 AddSingleton

// 在需要时获取
func HandleRequest(c *web.HttpContext, app *web.WebApplication) {
    userService := di.GetRequiredService[*UserService](app.Services)
}
```

---

### 2. Controllers 是单例

**设计：** 所有 Controllers 在 `app.MapControllers()` 时创建一次，整个应用生命周期复用。

#### ⚠️ 重要规则

```go
// ✅ 正确：无状态 Controller
type UserController struct {
    userService *UserService  // 依赖（不可变）
    config      *AppConfig    // 配置（不可变）
}

// ❌ 错误：有状态 Controller
type BadController struct {
    currentUser *User         // ❌ 请求状态，会被覆盖！
    requestID   string        // ❌ 线程不安全！
}
```

#### 最佳实践

1. **Controllers 只负责路由注册**
```go
func (c *UserController) MapRoutes(app *web.WebApplication) {
    app.MapGet("/users/:id", c.GetUser)
}
```

2. **请求数据从 HttpContext 获取**
```go
func (c *UserController) GetUser(ctx *web.HttpContext) web.IActionResult {
    id, _ := ctx.PathInt("id")              // ✅ 从请求获取
    userID := ctx.GetString("user_id")      // ✅ 从上下文获取
    
    user := c.userService.GetUser(id)       // ✅ 使用注入的服务
    return ctx.Ok(user)
}
```

3. **需要请求级服务时动态获取**
```go
func (c *UserController) GetUser(ctx *web.HttpContext) web.IActionResult {
    // ✅ 每次请求动态获取（如果服务是 Transient）
    logger := di.GetRequiredService[*RequestLogger](c.app.Services)
    logger.Log("Getting user...")
    
    return ctx.Ok(user)
}
```

---

### 3. 配置注入模式

**推荐：** 使用 `IOptions[T]` 或 `IOptionsMonitor[T]` 模式。

#### IOptions vs IOptionsMonitor

| 特性 | IOptions[T] | IOptionsMonitor[T] |
|------|-------------|-------------------|
| 读取时机 | 启动时一次 | 支持热更新 |
| 性能 | ⚡ 最快 | ⚠️ 有锁 |
| 推荐场景 | 静态配置（90%） | 动态配置（10%） |

#### 完整示例

**1. 定义配置结构**
```go
// config/settings.go
type DatabaseSettings struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Username string `json:"username"`
    Password string `json:"password"`
}

type AppSettings struct {
    AppName  string           `json:"appName"`
    Port     int              `json:"port"`
    Database DatabaseSettings `json:"database"`
}
```

**2. 注册配置（main.go）**
```go
import "github.com/gocrud/csgo/configuration"

func main() {
    builder := web.CreateBuilder()
    
    // 基础注册
    configuration.Configure[AppSettings](builder.Services, builder.Configuration, "App")
    
    // 带默认值
    configuration.ConfigureWithDefaults[AppSettings](
        builder.Services,
        builder.Configuration,
        "App",
        func() *AppSettings {
            return &AppSettings{Port: 8080}
        },
    )
    
    // 带验证
    err := configuration.ConfigureWithValidation[DatabaseSettings](
        builder.Services,
        builder.Configuration,
        "Database",
        func(opts *DatabaseSettings) error {
            if opts.Host == "" {
                return fmt.Errorf("database host is required")
            }
            return nil
        },
    )
    
    app := builder.Build()
    app.Run()
}
```

**3. 在 Controller 中使用**
```go
type UserController struct {
    userService *UserService
    appSettings *AppSettings  // ✅ 配置快照
}

func NewUserController(app *web.WebApplication) *UserController {
    userService := di.GetRequiredService[*UserService](app.Services)
    
    // ✅ 解析配置
    appOptions := di.GetRequiredService[configuration.IOptions[AppSettings]](app.Services)
    
    return &UserController{
        userService: userService,
        appSettings: appOptions.Value(),  // 获取配置值
    }
}

func (c *UserController) GetInfo(ctx *web.HttpContext) web.IActionResult {
    // ✅ 使用配置
    return ctx.Ok(gin.H{
        "app": c.appSettings.AppName,
        "db":  c.appSettings.Database.Host,
    })
}
```

**4. 使用 IOptionsMonitor（热更新）**
```go
type UserController struct {
    settingsMonitor configuration.IOptionsMonitor[AppSettings]
}

func NewUserController(app *web.WebApplication) *UserController {
    monitor := di.GetRequiredService[configuration.IOptionsMonitor[AppSettings]](app.Services)
    
    // 监听配置变化
    monitor.OnChange(func(newSettings *AppSettings, name string) {
        log.Printf("配置已更新: %v", newSettings)
    })
    
    return &UserController{settingsMonitor: monitor}
}

func (c *UserController) GetInfo(ctx *web.HttpContext) web.IActionResult {
    // ✅ 获取最新配置
    settings := c.settingsMonitor.CurrentValue()
    return ctx.Ok(settings)
}
```

---

### 4. 控制器注册简化

**当前设计：** `AddController` 不再注册到 DI 容器，只存储工厂函数。

#### 注册方式

```go
// controllers/controller_extensions.go
func AddControllers(services di.IServiceCollection) {
    // ✅ 方式1：简单的构造函数
    web.AddController(services, NewHealthController)
    
    // ✅ 方式2：带依赖解析
    web.AddController(services, func(sp di.IServiceProvider) *UserController {
        userService := di.GetRequiredService[*UserService](sp)
        return NewUserController(userService)
    })
    
    // ✅ 方式3：多个依赖
    web.AddController(services, func(sp di.IServiceProvider) *OrderController {
        orderService := di.GetRequiredService[*OrderService](sp)
        userService := di.GetRequiredService[*UserService](sp)
        logger := di.GetRequiredService[*Logger](sp)
        return NewOrderController(orderService, userService, logger)
    })
}
```

#### Controller 构造函数模式

```go
// 模式1：接收 WebApplication（可以动态解析服务）
func NewUserController(app *web.WebApplication) *UserController {
    userService := di.GetRequiredService[*UserService](app.Services)
    
    return &UserController{
        app:         app,
        userService: userService,
    }
}

// 模式2：接收具体依赖（依赖明确）
func NewUserController(userService *UserService, config *AppConfig) *UserController {
    return &UserController{
        userService: userService,
        config:      config,
    }
}

// 模式3：混合方式（既有固定依赖，也能动态解析）
type UserController struct {
    app         *web.WebApplication  // 用于动态解析
    userService *UserService         // 固定依赖
}

func NewUserController(app *web.WebApplication, userService *UserService) *UserController {
    return &UserController{
        app:         app,
        userService: userService,
    }
}
```

---

## 📋 API 变化清单

### 移除的 API

| API | 替代方案 |
|-----|---------|
| `AddScoped()` | 使用 `AddTransient()` 或 `AddSingleton()` |
| `TryAddScoped()` | 使用 `TryAddTransient()` |
| `AddKeyedScoped()` | 使用 `AddKeyedTransient()` |
| `CreateScope()` | 不需要，Controllers 是单例 |
| `IServiceScope` | 已移除 |
| `IServiceScopeFactory` | 已移除 |
| `GetServiceScopeFactory()` | 已移除 |
| `WithValidateScopes()` | 已移除 |

### 当前支持的生命周期

```go
// ✅ Singleton - 全局单例
services.AddSingleton(NewDatabaseConnection)
services.AddKeyedSingleton("primary", NewPrimaryDb)

// ✅ Transient - 每次创建新实例
services.AddTransient(NewEmailService)
services.AddKeyedTransient("logger", NewLogger)
```

---

## 🔄 与 .NET 的差异

| 特性 | .NET Core | CSGO | 原因 |
|------|-----------|------|------|
| **生命周期** | Singleton/Scoped/Transient | Singleton/Transient | Go 生态习惯 |
| **Controllers** | Scoped（每请求） | Singleton（全局） | 性能和简单性 |
| **HttpContext** | 自动注入 | 通过参数传递 | Go 语言限制 |
| **参数绑定** | 自动（FromBody 等） | 手动（BindJSON） | Go 语言限制 |
| **作用域管理** | 自动（每请求创建） | 手动（需要时获取） | 简化设计 |

---

## ✅ 最佳实践总结

### 1. 服务注册

```go
// ✅ 无状态服务 → Singleton
services.AddSingleton(NewDatabaseConnection)
services.AddSingleton(NewCache)

// ✅ 有状态/轻量级服务 → Transient
services.AddTransient(NewEmailService)
services.AddTransient(NewRequestLogger)

// ✅ 配置 → IOptions
configuration.Configure[AppSettings](services, config, "App")
```

### 2. Controller 设计

```go
type UserController struct {
    // ✅ 允许：服务依赖
    userService *UserService
    
    // ✅ 允许：配置
    settings *AppSettings
    
    // ✅ 允许：WebApplication（动态解析）
    app *web.WebApplication
    
    // ❌ 禁止：请求状态
    // currentUser *User  // ❌ 不要这样！
}
```

### 3. 依赖注入

```go
// ✅ 推荐：泛型辅助函数
userService := di.GetRequiredService[*UserService](provider)

// ✅ 可选：指针填充
var userService *UserService
provider.GetRequiredService(&userService)
```

### 4. 配置使用

```go
// ✅ 静态配置 → IOptions
opts := di.GetRequiredService[configuration.IOptions[AppSettings]](sp)
settings := opts.Value()

// ✅ 动态配置 → IOptionsMonitor
monitor := di.GetRequiredService[configuration.IOptionsMonitor[AppSettings]](sp)
settings := monitor.CurrentValue()
```

---

## 📚 相关文档

- [依赖注入 API 参考](api/di.md)
- [控制器指南](guides/controllers.md)
- [配置 API 参考](api/configuration.md)
- [Web 应用指南](guides/web-applications.md)

---

## 🎯 总结

CSGO 框架采用**简化设计**，权衡了以下因素：

| 方面 | 取舍 |
|------|------|
| **性能** | ✅ 优先 - 无运行时作用域开销 |
| **简单性** | ✅ 优先 - 代码和概念更简单 |
| **Go 习惯** | ✅ 符合 - Gin/Echo 都这样 |
| **.NET 一致性** | ⚠️ 折衷 - 核心概念一致，生命周期简化 |

这种设计让框架更快、更简单、更符合 Go 生态，同时保持了与 .NET 的核心理念一致。

