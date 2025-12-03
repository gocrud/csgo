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

### AddScoped

注册作用域服务，每个请求/作用域创建一个实例。

```go
func (sc *ServiceCollection) AddScoped(factory interface{})
```

**参数：**
- `factory` - 工厂函数，返回服务实例

**示例：**

```go
services.AddScoped(func(db *DatabaseConnection) *UserRepository {
    return NewUserRepository(db)
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

## IServiceProvider

服务提供者接口，用于解析已注册的服务。

### GetService

尝试获取服务，如果未注册返回 false。

```go
func (sp *ServiceProvider) GetService(target interface{}) bool
```

**参数：**
- `target` - 指向服务指针的指针（用于接收服务实例）

**返回值：**
- `bool` - 是否成功获取服务

**示例：**

```go
var userService *UserService
if provider.GetService(&userService) {
    // 服务存在
    userService.DoSomething()
} else {
    // 服务不存在
}
```

---

### GetRequiredService

获取必需的服务，如果未注册则 panic。

```go
func (sp *ServiceProvider) GetRequiredService(target interface{})
```

**参数：**
- `target` - 指向服务指针的指针

**示例：**

```go
var userService *UserService
provider.GetRequiredService(&userService)
// 如果服务不存在，会 panic
```

---

### CreateScope

创建新的服务作用域。

```go
func (sp *ServiceProvider) CreateScope() IServiceScope
```

**返回值：**
- `IServiceScope` - 新的服务作用域

**示例：**

```go
scope := provider.CreateScope()
defer scope.Dispose()

var scopedService *MyScopedService
scope.ServiceProvider().GetRequiredService(&scopedService)
```

---

## 服务生命周期

| 生命周期 | 说明 | 使用场景 |
|---------|------|---------|
| `Singleton` | 全局唯一实例 | 数据库连接、配置、缓存 |
| `Scoped` | 每作用域一个实例 | 请求上下文、工作单元 |
| `Transient` | 每次请求新实例 | 无状态服务、轻量级操作 |

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
userService := di.GetRequiredService[*UserService](provider)
```

---

### GetService[T]

泛型方式尝试获取服务。

```go
func GetService[T any](provider IServiceProvider) (T, bool)
```

**类型参数：**
- `T` - 服务类型

**返回值：**
- `T` - 服务实例（如果不存在则为零值）
- `bool` - 是否成功获取

**示例：**

```go
userService, ok := di.GetService[*UserService](provider)
if ok {
    userService.DoSomething()
}
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
    
    builder.Services.AddScoped(func(db *DatabaseConnection) *UserRepository {
        return &UserRepository{db: db}
    })
    
    builder.Services.AddScoped(func(repo *UserRepository) *UserService {
        return &UserService{repo: repo}
    })
    
    app := builder.Build()
    
    app.MapGet("/users", func(c *web.HttpContext) web.IActionResult {
        // 使用泛型辅助函数（推荐）
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

## 相关文档

- [依赖注入指南](../guides/dependency-injection.md)
- [Web 应用 API](web.md)
- [配置 API](configuration.md)

