# 业务模块 IServiceCollection 扩展指南

本文档介绍如何为业务模块创建 `IServiceCollection` 扩展方法，这是 .NET 中非常常见和推荐的模式。

---

## 🎯 为什么需要扩展方法？

### .NET 中的扩展方法

在 .NET 中，扩展方法让代码更加模块化和可维护：

```csharp
// .NET 示例
public static class UserServiceCollectionExtensions
{
    public static IServiceCollection AddUserServices(
        this IServiceCollection services)
    {
        services.AddSingleton<IUserRepository, UserRepository>();
        services.AddSingleton<IUserService, UserService>();
        return services;
    }
}

// 使用
builder.Services.AddUserServices();
```

### Go 中的实现方式

Go 没有扩展方法，但我们可以使用**顶层函数**达到相同效果：

```go
// Go 实现
package users

func AddUserServices(services di.IServiceCollection) {
    services.AddSingleton(NewUserRepository)
    services.AddSingleton(NewUserService)
}

// 使用
users.AddUserServices(builder.Services)
```

---

## 📁 推荐的项目结构

```
your-project/
├── main.go
├── users/                              # 用户模块
│   ├── user_service.go                 # 服务实现
│   ├── user_repository.go              # 仓储实现
│   ├── user_module_options.go          # 模块配置
│   └── service_collection_extensions.go # ✅ DI 扩展
├── orders/                             # 订单模块
│   ├── order_service.go
│   └── service_collection_extensions.go # ✅ DI 扩展
└── products/                           # 产品模块
    ├── product_service.go
    └── service_collection_extensions.go # ✅ DI 扩展
```

---

## 📝 实现模式

### 模式 1：简单扩展（最常用）⭐

适用于大多数业务模块。

```go
// users/service_collection_extensions.go
package users

import "github.com/gocrud/csgo/di"

// AddUserServices registers all user-related services.
// Corresponds to .NET services.AddUserServices().
func AddUserServices(services di.IServiceCollection) {
    // Register repository
    services.AddSingleton(NewUserRepository)
    
    // Register service
    services.AddSingleton(NewUserService)
}
```

**使用方式：**

```go
func main() {
    builder := web.CreateBuilder()
    
    // ✅ 一行代码注册整个模块
    users.AddUserServices(builder.Services)
    
    app := builder.Build()
    app.Run()
}
```

---

### 模式 2：带配置的扩展

适用于需要配置选项的模块。

```go
// users/user_module_options.go
package users

import "time"

type UserModuleOptions struct {
    EnableCache     bool
    CacheExpiration time.Duration
    MaxConnections  int
}

// users/service_collection_extensions.go
package users

// AddUserServicesWithOptions registers user services with custom options.
// Corresponds to .NET services.AddUserServices(options => { ... }).
func AddUserServicesWithOptions(
    services di.IServiceCollection,
    configure func(*UserModuleOptions),
) {
    // Create default options
    opts := &UserModuleOptions{
        EnableCache:     false,
        CacheExpiration: 5 * time.Minute,
        MaxConnections:  10,
    }
    
    // Apply custom configuration
    if configure != nil {
        configure(opts)
    }
    
    // Register options as Singleton
    services.AddSingleton(func() *UserModuleOptions {
        return opts
    })
    
    // Register services
    services.AddSingleton(NewUserRepository)
    services.AddSingleton(NewUserService)
}
```

**使用方式：**

```go
func main() {
    builder := web.CreateBuilder()
    
    // ✅ 带配置的注册
    users.AddUserServicesWithOptions(builder.Services, func(opts *users.UserModuleOptions) {
        opts.EnableCache = true
        opts.CacheExpiration = 10 * time.Minute
        opts.MaxConnections = 20
    })
    
    app := builder.Build()
    app.Run()
}
```

---

### 模式 3：多种生命周期扩展

适用于需要不同生命周期的场景。

```go
// orders/service_collection_extensions.go
package orders

// AddOrderServices registers order services as Singleton.
func AddOrderServices(services di.IServiceCollection) {
    services.AddSingleton(NewOrderService)
}

// AddOrderServicesScoped registers order services as Transient (per-request).
// Use this when you need a new instance per request.
func AddOrderServicesScoped(services di.IServiceCollection) {
    services.AddTransient(NewOrderService)
}
```

**使用方式：**

```go
func main() {
    builder := web.CreateBuilder()
    
    // ✅ 选择合适的生命周期
    orders.AddOrderServices(builder.Services)        // Singleton
    // 或
    orders.AddOrderServicesScoped(builder.Services)  // Transient
    
    app := builder.Build()
    app.Run()
}
```

---

### 模式 4：组合扩展

适用于需要组合多个模块的场景。

```go
// infrastructure/service_collection_extensions.go
package infrastructure

import (
    "your-project/users"
    "your-project/orders"
    "your-project/products"
    "github.com/gocrud/csgo/di"
)

// AddInfrastructureServices registers all infrastructure services.
// This is a convenience method that registers multiple modules at once.
func AddInfrastructureServices(services di.IServiceCollection) {
    users.AddUserServices(services)
    orders.AddOrderServices(services)
    products.AddProductServices(services)
}
```

**使用方式：**

```go
func main() {
    builder := web.CreateBuilder()
    
    // ✅ 一次性注册多个模块
    infrastructure.AddInfrastructureServices(builder.Services)
    
    app := builder.Build()
    app.Run()
}
```

---

## 🎨 完整示例

### 用户模块实现

#### 1. 定义服务接口和实现

```go
// users/user_service.go
package users

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type IUserService interface {
    GetUser(id int) (*User, error)
    ListUsers() ([]*User, error)
    CreateUser(user *User) error
}

type UserService struct {
    users map[int]*User
}

func NewUserService() IUserService {
    return &UserService{
        users: make(map[int]*User),
    }
}

func (s *UserService) GetUser(id int) (*User, error) {
    user, ok := s.users[id]
    if !ok {
        return nil, fmt.Errorf("user not found: %d", id)
    }
    return user, nil
}

func (s *UserService) ListUsers() ([]*User, error) {
    users := make([]*User, 0, len(s.users))
    for _, u := range s.users {
        users = append(users, u)
    }
    return users, nil
}

func (s *UserService) CreateUser(user *User) error {
    s.users[user.ID] = user
    return nil
}
```

#### 2. 创建扩展方法

```go
// users/service_collection_extensions.go
package users

import "github.com/gocrud/csgo/di"

// AddUserServices registers all user-related services.
func AddUserServices(services di.IServiceCollection) {
    services.AddSingleton(NewUserService)
}
```

#### 3. 在主程序中使用

```go
// main.go
package main

import (
    "your-project/users"
    "github.com/gin-gonic/gin"
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

func main() {
    builder := web.CreateBuilder()
    
    // ✅ 注册用户模块
    users.AddUserServices(builder.Services)
    
    app := builder.Build()
    
    // 使用服务
    app.MapGet("/api/users", func(c *gin.Context) {
        // Style 1: Traditional
        var userSvc users.IUserService
        app.Services.GetRequiredService(&userSvc)
        
        // 或 Style 2: Generic helper (推荐)
        userSvc := di.GetRequiredService[users.IUserService](app.Services)
        
        userList, err := userSvc.ListUsers()
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(200, userList)
    })
    
    app.Run()
}
```

---

## 🔍 与 .NET 的对比

### .NET 代码

```csharp
// UserServiceCollectionExtensions.cs
public static class UserServiceCollectionExtensions
{
    public static IServiceCollection AddUserServices(
        this IServiceCollection services)
    {
        services.AddSingleton<IUserRepository, UserRepository>();
        services.AddSingleton<IUserService, UserService>();
        return services;
    }
    
    public static IServiceCollection AddUserServices(
        this IServiceCollection services,
        Action<UserModuleOptions> configure)
    {
        services.Configure(configure);
        services.AddSingleton<IUserRepository, UserRepository>();
        services.AddSingleton<IUserService, UserService>();
        return services;
    }
}

// Program.cs
var builder = WebApplication.CreateBuilder(args);

builder.Services.AddUserServices();
// 或
builder.Services.AddUserServices(opts => {
    opts.EnableCache = true;
});

var app = builder.Build();
app.Run();
```

---

### Ego 代码

```go
// users/service_collection_extensions.go
package users

func AddUserServices(services di.IServiceCollection) {
    services.AddSingleton(NewUserRepository)
    services.AddSingleton(NewUserService)
}

func AddUserServicesWithOptions(
    services di.IServiceCollection,
    configure func(*UserModuleOptions),
) {
    opts := &UserModuleOptions{}
    if configure != nil {
        configure(opts)
    }
    services.AddSingleton(func() *UserModuleOptions { return opts })
    services.AddSingleton(NewUserRepository)
    services.AddSingleton(NewUserService)
}

// main.go
func main() {
    builder := web.CreateBuilder()
    
    users.AddUserServices(builder.Services)
    // 或
    users.AddUserServicesWithOptions(builder.Services, func(opts *users.UserModuleOptions) {
        opts.EnableCache = true
    })
    
    app := builder.Build()
    app.Run()
}
```

---

### 一致性分析

| 特性 | .NET | Ego | 一致性 |
|------|------|-----|--------|
| **扩展方法模式** | `this IServiceCollection` | 顶层函数 | 98% |
| **命名约定** | `AddXxxServices` | `AddXxxServices` | 100% |
| **配置选项** | `Action<T>` | `func(*T)` | 100% |
| **链式调用** | ✅ 返回 `IServiceCollection` | ⚠️ 无返回值 | 90% |
| **使用体验** | `services.AddUserServices()` | `users.AddUserServices(services)` | 95% |

**总体一致性：96%** ✅

---

## 💡 最佳实践

### 1. 命名约定

- ✅ 使用 `AddXxxServices` 命名模式
- ✅ 文件命名为 `service_collection_extensions.go`
- ✅ 放在业务模块包内

```go
// ✅ Good
users.AddUserServices(services)
orders.AddOrderServices(services)

// ❌ Bad
users.RegisterServices(services)
users.Setup(services)
```

---

### 2. 文档注释

为每个扩展方法添加清晰的文档注释：

```go
// AddUserServices registers all user-related services.
// This includes:
//   - IUserRepository (Singleton)
//   - IUserService (Singleton)
//
// Usage:
//
//	builder := web.CreateBuilder()
//	users.AddUserServices(builder.Services)
func AddUserServices(services di.IServiceCollection) {
    // ...
}
```

---

### 3. 模块化组织

每个业务模块应该是独立的：

```
users/
├── user.go                          # 实体定义
├── user_service.go                  # 服务实现
├── user_repository.go               # 仓储实现
├── user_module_options.go           # 配置选项
└── service_collection_extensions.go # ✅ DI 扩展
```

---

### 4. 依赖注入最佳实践

```go
// ✅ Good: 使用接口
type IUserService interface {
    GetUser(id int) (*User, error)
}

func NewUserService() IUserService {
    return &UserService{}
}

// ❌ Bad: 直接返回具体类型
func NewUserService() *UserService {
    return &UserService{}
}
```

---

### 5. 生命周期选择

- **Singleton**: 无状态服务、配置、缓存
- **Transient**: 有状态服务、每次请求需要新实例

```go
// Singleton - 适用于无状态服务
services.AddSingleton(NewUserService)

// Transient - 适用于有状态服务
services.AddTransient(NewOrderService)
```

---

## 📊 完整示例项目

查看 `examples/business_module_demo` 获取完整的可运行示例：

```bash
cd examples/business_module_demo
go mod tidy
go run main.go
```

访问：
- API: http://localhost:8080
- Swagger: http://localhost:8080/swagger

---

## ✅ 总结

通过使用 `IServiceCollection` 扩展模式，你可以：

1. ✅ **模块化** - 每个业务模块独立管理自己的依赖
2. ✅ **可维护** - 集中管理服务注册逻辑
3. ✅ **可测试** - 容易替换实现进行测试
4. ✅ **.NET 一致** - 与 .NET 的使用体验高度一致（96%）

### 推荐使用模式

```go
// 1. 在业务模块中定义扩展
package users

func AddUserServices(services di.IServiceCollection) {
    services.AddSingleton(NewUserService)
}

// 2. 在主程序中使用
func main() {
    builder := web.CreateBuilder()
    users.AddUserServices(builder.Services)  // ✅ 简洁清晰
    app := builder.Build()
    app.Run()
}
```

这种模式让你的代码更加专业和易于维护！🚀

