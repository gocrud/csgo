# 依赖注入 API 参考

本文档提供 CSGO 依赖注入系统的完整 API 参考。

## 目录

- [IServiceCollection](#iservicecollection)
- [IServiceProvider](#iserviceprovider)
- [服务生命周期](#服务生命周期)
- [泛型辅助函数](#泛型辅助函数)

---

## IServiceCollection

服务注册接口，用于在应用启动时注册服务。

```go
import "github.com/gocrud/csgo/di"
```

### AddSingleton

注册单例服务，整个应用生命周期只创建一个实例。

```go
func (sc *ServiceCollection) AddSingleton(factory interface{})
```

**参数：**
- `factory` - 工厂函数，返回服务实例

**示例：**

```go
// 无依赖
services.AddSingleton(func() *DatabaseConnection {
    return NewDatabaseConnection()
})

// 有依赖
services.AddSingleton(func(config *Config) *DatabaseConnection {
    return NewDatabaseConnection(config.ConnectionString)
})
```

---

### AddTransient

注册瞬态服务，每次请求都创建新实例。

```go
func (sc *ServiceCollection) AddTransient(factory interface{})
```

**参数：**
- `factory` - 工厂函数，返回服务实例

**示例：**

```go
services.AddTransient(func() *EmailService {
    return NewEmailService()
})

// 有依赖
services.AddTransient(func(logger *Logger) *EmailService {
    return NewEmailService(logger)
})
```

---

### AddHostedService

注册托管服务，随应用启动和停止。

```go
func (sc *ServiceCollection) AddHostedService(factory interface{})
```

**参数：**
- `factory` - 工厂函数，返回 `hosting.IHostedService`

**示例：**

```go
services.AddHostedService(func() hosting.IHostedService {
    return NewBackgroundWorker()
})
```

---

### AddKeyedSingleton / AddKeyedTransient

注册命名服务，允许同一类型的多个实现。

```go
func (sc *ServiceCollection) AddKeyedSingleton(key string, factory interface{})
func (sc *ServiceCollection) AddKeyedTransient(key string, factory interface{})
```

**示例：**

```go
// 注册多个数据库连接
services.AddKeyedSingleton("primary", func() *Database {
    return NewDatabase("postgres://primary")
})
services.AddKeyedSingleton("secondary", func() *Database {
    return NewDatabase("postgres://secondary")
})

// 获取特定的实现
var primary *Database
provider.GetKeyedService(&primary, "primary")
```

---

## IServiceProvider

服务提供者接口，用于解析已注册的服务。

### GetService

尝试获取服务，如果未注册返回错误。

```go
func (sp *ServiceProvider) GetService(target interface{}) error
```

**参数：**
- `target` - 指向服务的指针

**返回值：**
- `error` - 错误信息（如果服务未注册）

**示例：**

```go
var userService *UserService
err := provider.GetService(&userService)
if err != nil {
    // 服务不存在，处理错误
}
```

---

### GetRequiredService

获取必需的服务，如果未注册则 panic。

```go
func (sp *ServiceProvider) GetRequiredService(target interface{})
```

**参数：**
- `target` - 指向服务的指针

**示例：**

```go
var userService *UserService
provider.GetRequiredService(&userService)
// 如果服务不存在，会 panic
```

---

### TryGetService

尝试获取服务，返回布尔值表示成功与否。

```go
func (sp *ServiceProvider) TryGetService(target interface{}) bool
```

**返回值：**
- `bool` - 是否成功获取服务

**示例：**

```go
var optionalService *OptionalService
if provider.TryGetService(&optionalService) {
    // 服务存在
    optionalService.DoSomething()
}
```

---

### GetKeyedService

获取命名服务。

```go
func (sp *ServiceProvider) GetKeyedService(target interface{}, key string) error
```

**示例：**

```go
var primaryDb *Database
provider.GetKeyedService(&primaryDb, "primary")
```

---

## 服务生命周期

CSGO 框架支持两种服务生命周期：

| 生命周期 | 说明 | 使用场景 | 性能 |
|---------|------|---------|------|
| **Singleton** | 全局唯一实例，应用启动时创建 | 数据库连接池、配置、缓存、无状态服务 | ⚡ 最快 |
| **Transient** | 每次请求创建新实例 | 有状态的请求处理、轻量级操作 | ⭐ 适中 |

### 重要说明

⚠️ **框架不支持 Scoped 生命周期**

与 ASP.NET Core 不同，CSGO 框架采用更简单的设计：
- **Controllers 是单例的** - 在应用启动时创建一次，整个生命周期复用
- **Controllers 必须是无状态的** - 不要在 Controller 字段中存储请求相关数据
- **业务逻辑在服务层** - 服务可以是 Singleton 或 Transient

这种设计：
- ✅ 符合 Go 生态习惯（Gin/Echo 都是这样）
- ✅ 性能更好（无运行时开销）
- ✅ 代码更简单（无复杂的作用域管理）

---

## 泛型辅助函数

### GetRequiredService[T]

泛型方式获取服务（推荐）。

```go
func GetRequiredService[T any](provider IServiceProvider) T
```

**类型参数：**
- `T` - 服务类型

**返回值：**
- `T` - 服务实例

**示例：**

```go
// ✅ 推荐：一行代码获取服务
userService := di.GetRequiredService[*UserService](provider)
```

---

### GetService[T]

泛型方式尝试获取服务。

```go
func GetService[T any](provider IServiceProvider) (T, error)
```

**类型参数：**
- `T` - 服务类型

**返回值：**
- `T` - 服务实例（如果不存在则为零值）
- `error` - 错误信息

**示例：**

```go
userService, err := di.GetService[*UserService](provider)
if err != nil {
    // 处理错误
}
```

---

### GetKeyedService[T]

泛型方式获取命名服务。

```go
func GetKeyedService[T any](provider IServiceProvider, key string) (T, error)
```

**示例：**

```go
primaryDb := di.GetRequiredKeyedService[*Database](provider, "primary")
```

---

## 完整示例

```go
package main

import (
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

// 服务定义
type Config struct {
    DatabaseURL string
}

type DatabaseConnection struct {
    url string
}

type UserRepository struct {
    db *DatabaseConnection
}

type UserService struct {
    repo *UserRepository
}

func main() {
    builder := web.CreateBuilder()
    
    // 注册服务链
    builder.Services.AddSingleton(func() *Config {
        return &Config{DatabaseURL: "postgres://localhost/mydb"}
    })
    
    builder.Services.AddSingleton(func(config *Config) *DatabaseConnection {
        return &DatabaseConnection{url: config.DatabaseURL}
    })
    
    builder.Services.AddTransient(func(db *DatabaseConnection) *UserRepository {
        return &UserRepository{db: db}
    })
    
    builder.Services.AddTransient(func(repo *UserRepository) *UserService {
        return &UserService{repo: repo}
    })
    
    app := builder.Build()
    
    app.MapGet("/users", func(c *web.HttpContext) web.IActionResult {
        // ✅ 推荐：使用泛型辅助函数
        userService := di.GetRequiredService[*UserService](app.Services)
        
        // 或使用指针填充方式
        // var userService *UserService
        // app.Services.GetRequiredService(&userService)
        
        return c.Ok(userService.ListUsers())
    })
    
    app.Run()
}
```

---

## Controller 中的依赖注入

Controllers 是单例的，在应用启动时创建一次：

```go
type UserController struct {
    userService *UserService  // 在构造函数中注入
}

func NewUserController(app *web.WebApplication) *UserController {
    // 从 DI 容器解析服务
    userService := di.GetRequiredService[*UserService](app.Services)
    
    return &UserController{
        userService: userService,
    }
}

func (c *UserController) GetUser(ctx *web.HttpContext) web.IActionResult {
    // 使用注入的服务
    id, _ := ctx.PathInt("id")
    user := c.userService.GetUser(id)
    return ctx.Ok(user)
}
```

---

## 配置注入

使用 IOptions 模式注入配置（推荐）：

```go
import "github.com/gocrud/csgo/configuration"

// 定义配置
type DatabaseOptions struct {
    Host string
    Port int
}

// 注册配置
configuration.Configure[DatabaseOptions](builder.Services, builder.Configuration, "Database")

// Controller 中使用
type UserController struct {
    dbOptions configuration.IOptions[DatabaseOptions]
}

func NewUserController(app *web.WebApplication) *UserController {
    dbOptions := di.GetRequiredService[configuration.IOptions[DatabaseOptions]](app.Services)
    
    return &UserController{
        dbOptions: dbOptions,
    }
}

func (c *UserController) Connect(ctx *web.HttpContext) web.IActionResult {
    opts := c.dbOptions.Value()
    // 使用 opts.Host 和 opts.Port
    return ctx.Ok(fmt.Sprintf("Connected to %s:%d", opts.Host, opts.Port))
}
```

---

## 相关文档

- [依赖注入指南](../guides/dependency-injection.md)
- [控制器指南](../guides/controllers.md)
- [配置 API](configuration.md)
- [Web 应用 API](web.md)
