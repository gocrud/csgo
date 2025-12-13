# 依赖注入 (DI)

[← 返回主目录](../README.md)

CSGO 的依赖注入系统提供了完整的服务容器和依赖管理功能，支持自动依赖解析、类型安全的服务注册和获取。

## 特性

- ✅ 类型安全的服务注册和解析
- ✅ 自动依赖注入（构造函数注入）
- ✅ Singleton 生命周期管理
- ✅ 泛型 API（Get、GetOr、GetAll、TryGet）
- ✅ 命名服务支持
- ✅ 实例注册
- ✅ 条件注册（TryAdd）
- ✅ 资源自动释放（IDisposable）

## 快速开始

### 1. 基本使用

```go
package main

import (
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

// 定义服务
type UserService struct {
    repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
    return &UserService{repo: repo}
}

type UserRepository struct{}

func NewUserRepository() *UserRepository {
    return &UserRepository{}
}

func main() {
    // 创建应用构建器
    builder := web.CreateBuilder()
    
    // 注册服务（自动依赖注入）
    builder.Services.Add(NewUserRepository)
    builder.Services.Add(NewUserService)  // 自动注入 UserRepository
    
    app := builder.Build()
    
    // 使用服务
    app.MapGet("/users", func() string {
        userService := di.Get[*UserService](app.Services)
        return "Hello from UserService"
    })
    
    app.Run()
}
```

### 2. 在控制器中使用

```go
type UserController struct {
    userService *UserService
    logger      logging.ILogger
}

func NewUserController(
    userService *UserService,
    loggerFactory logging.ILoggerFactory,
) *UserController {
    return &UserController{
        userService: userService,
        logger:      logging.GetLogger[UserController](loggerFactory),
    }
}

// 注册控制器
web.AddController(builder.Services, NewUserController)
```

## 核心概念

### IServiceCollection

服务集合接口，用于注册服务。所有服务注册都通过此接口完成。

```go
type IServiceCollection interface {
    // 注册单例服务
    Add(constructor interface{}) IServiceCollection
    
    // 注册单例实例
    AddInstance(instance interface{}) IServiceCollection
    
    // 注册命名服务
    AddNamed(name string, constructor interface{}) IServiceCollection
    
    // 条件注册（如果不存在才注册）
    TryAdd(constructor interface{}) IServiceCollection
    
    // 注册后台服务
    AddHostedService(constructor interface{}) IServiceCollection
}
```

### IServiceProvider

服务提供者接口，用于解析和获取服务。

```go
type IServiceProvider interface {
    // 获取服务（指针填充方式）
    Get(target interface{})
    
    // 获取命名服务
    GetNamed(target interface{}, serviceKey string)
    
    // 释放资源
    Dispose() error
}
```

### ServiceLifetime

服务生命周期枚举。当前版本支持 Singleton 生命周期。

```go
const (
    Singleton ServiceLifetime = iota  // 单例模式
)
```

## Singleton 生命周期

### 什么是 Singleton？

Singleton（单例）生命周期表示服务在整个应用程序生命周期内只创建一次实例，所有请求共享同一个实例。

### 特点

- **全局唯一**：整个应用只有一个实例
- **首次创建**：第一次请求时创建，之后复用
- **线程安全**：框架保证创建过程的线程安全
- **生命周期长**：随应用启动而创建，随应用关闭而销毁

### 适用场景

✅ **推荐使用 Singleton 的场景：**

1. **无状态服务**：不保存请求相关状态的服务
   ```go
   type EmailService struct{}
   
   func (s *EmailService) SendEmail(to, subject, body string) error {
       // 无状态操作，可以安全地共享
       return nil
   }
   ```

2. **配置对象**：应用配置、常量数据
   ```go
   type AppConfig struct {
       Port     int
       Database string
   }
   ```

3. **资源密集型对象**：数据库连接池、HTTP 客户端
   ```go
   type DatabaseConnection struct {
       pool *sql.DB
   }
   ```

4. **缓存服务**：内存缓存、分布式缓存客户端
   ```go
   type CacheService struct {
       cache map[string]interface{}
       mu    sync.RWMutex
   }
   ```

❌ **不推荐使用 Singleton 的场景：**

1. **有状态服务**：保存用户会话、请求上下文等
2. **请求相关数据**：需要按请求隔离的数据
3. **一次性使用对象**：每次使用都需要新实例的对象

### 线程安全

使用 Singleton 服务时，如果服务有可变状态，必须确保线程安全：

```go
type CounterService struct {
    count int64
}

func (s *CounterService) Increment() int64 {
    // ❌ 不安全：并发访问会有竞态条件
    s.count++
    return s.count
}

// ✅ 安全：使用原子操作
func (s *CounterService) IncrementSafe() int64 {
    return atomic.AddInt64(&s.count, 1)
}
```

## 服务注册

### Add - 注册单例服务

使用构造函数注册服务，支持自动依赖注入。

```go
// 无依赖的服务
builder.Services.Add(NewUserRepository)

// 有依赖的服务（自动注入依赖）
builder.Services.Add(NewUserService)  // 自动注入 UserRepository

// 构造函数可以返回 (service, error)
builder.Services.Add(func() (*DatabaseConnection, error) {
    db, err := sql.Open("postgres", "...")
    if err != nil {
        return nil, err
    }
    return &DatabaseConnection{db: db}, nil
})
```

### AddInstance - 注册实例

注册已创建的实例对象。

```go
// 注册配置实例
config := &AppConfig{
    Port:     8080,
    Database: "postgres://...",
}
builder.Services.AddInstance(config)

// 注册全局变量
var globalCache = NewCacheService()
builder.Services.AddInstance(globalCache)
```

**注意**：AddInstance 适用于应用启动时就需要创建的对象，或者外部传入的对象。

### AddNamed - 注册命名服务

为同一类型注册多个实例，使用名称区分。

```go
// 注册多个数据库连接
builder.Services.AddNamed("primary", func() *Database {
    return NewDatabase("primary-connection-string")
})

builder.Services.AddNamed("secondary", func() *Database {
    return NewDatabase("secondary-connection-string")
})

builder.Services.AddNamed("cache", func() *Database {
    return NewDatabase("cache-connection-string")
})
```

### TryAdd - 条件注册

只在服务不存在时才注册，避免覆盖已有注册。

```go
// 默认实现
builder.Services.TryAdd(NewDefaultEmailService)

// 用户可以在之前注册自定义实现来覆盖默认值
// 如果已注册，TryAdd 不会覆盖
builder.Services.Add(NewCustomEmailService)
builder.Services.TryAdd(NewDefaultEmailService)  // 不会注册
```

**使用场景**：

- 框架提供默认实现
- 库提供可选的默认服务
- 避免重复注册

### AddHostedService - 注册后台服务

注册实现 `IHostedService` 接口的后台服务。

```go
type BackgroundWorker struct {
    *hosting.BackgroundService
}

func NewBackgroundWorker() *BackgroundWorker {
    worker := &BackgroundWorker{
        BackgroundService: hosting.NewBackgroundService(),
    }
    worker.SetExecuteFunc(worker.doWork)
    return worker
}

func (w *BackgroundWorker) doWork(ctx context.Context) error {
    // 后台任务逻辑
    return nil
}

// 注册
builder.Services.AddHostedService(NewBackgroundWorker)
```

## 服务解析

### Get - 泛型获取（推荐）

使用泛型 API 获取服务，类型安全且简洁。

```go
// 获取指针类型服务
userService := di.Get[*UserService](provider)

// 获取值类型服务（自动解引用）
config := di.Get[AppConfig](provider)

// 如果服务不存在会 panic
```

### GetOr - 带默认值获取

获取服务，如果不存在则返回默认值，不会 panic。

```go
// 提供默认值
defaultConfig := &AppConfig{Port: 8080}
config := di.GetOr[*AppConfig](provider, defaultConfig)

// 适用于可选服务
logger := di.GetOr[*Logger](provider, nil)
if logger != nil {
    logger.Log("Service available")
}
```

### TryGet - 安全获取

尝试获取服务，返回 (service, ok) 形式。

```go
// 检查服务是否存在
if svc, ok := di.TryGet[*UserService](provider); ok {
    // 服务存在，使用它
    svc.DoSomething()
} else {
    // 服务不存在，使用替代逻辑
    log.Println("UserService not available")
}
```

### GetNamed - 获取命名服务

根据名称获取特定的服务实例。

```go
// 获取命名服务
primaryDB := di.GetNamed[*Database](provider, "primary")
secondaryDB := di.GetNamed[*Database](provider, "secondary")
cacheDB := di.GetNamed[*Database](provider, "cache")
```

### GetAll - 获取所有实例

获取某个类型的所有已注册实例。

```go
// 获取所有插件
plugins := di.GetAll[IPlugin](provider)
for _, plugin := range plugins {
    plugin.Initialize()
}

// 获取所有通知器
notifiers := di.GetAll[*INotifier](provider)
for _, notifier := range notifiers {
    notifier.Notify("Event occurred")
}
```

### 指针填充方式（可选）

除了泛型 API，也可以使用传统的指针填充方式：

```go
// 获取服务
var userService *UserService
provider.Get(&userService)

// 获取命名服务
var primaryDB *Database
provider.GetNamed(&primaryDB, "primary")
```

## 依赖自动解析

### 构造函数注入

DI 容器会自动分析构造函数参数，并注入所需的依赖。

```go
// 定义服务层级
type Repository struct{}

type Service struct {
    repo *Repository
}

type Controller struct {
    service *Service
    config  *Config
}

// 注册（按依赖顺序或任意顺序都可以）
builder.Services.Add(func() *Repository {
    return &Repository{}
})

builder.Services.Add(func(repo *Repository) *Service {
    return &Service{repo: repo}  // 自动注入 Repository
})

builder.Services.Add(func(svc *Service, cfg *Config) *Controller {
    return &Controller{
        service: svc,   // 自动注入 Service
        config:  cfg,   // 自动注入 Config
    }
})
```

### 依赖图

DI 容器会自动构建依赖图并按正确顺序创建实例：

```
Repository (无依赖)
    ↓
Service (依赖 Repository)
    ↓
Controller (依赖 Service, Config)
```

### 循环依赖检测

框架会自动检测循环依赖并报错：

```go
// ❌ 循环依赖会导致 panic
type ServiceA struct {
    b *ServiceB
}

type ServiceB struct {
    a *ServiceA
}

builder.Services.Add(func(b *ServiceB) *ServiceA {
    return &ServiceA{b: b}
})

builder.Services.Add(func(a *ServiceA) *ServiceB {
    return &ServiceB{a: a}
})

// Build 时会检测到循环依赖并报错
app := builder.Build()  // panic: circular dependency detected
```

**解决方案**：

1. 重新设计服务架构，消除循环依赖
2. 使用事件/消息模式解耦
3. 延迟解析依赖（在需要时从 provider 获取）

### 可选依赖

使用 `GetOr` 或 `TryGet` 实现可选依赖：

```go
type Service struct {
    required *RequiredService
    optional *OptionalService
}

func NewService(
    provider di.IServiceProvider,
    required *RequiredService,
) *Service {
    // 可选依赖，如果不存在使用 nil
    optional := di.GetOr[*OptionalService](provider, nil)
    
    return &Service{
        required: required,
        optional: optional,
    }
}
```

## 资源管理

### IDisposable 接口

实现 `IDisposable` 接口的服务会在应用关闭时自动释放资源。

```go
type DatabaseConnection struct {
    db *sql.DB
}

func (d *DatabaseConnection) Dispose() error {
    if d.db != nil {
        return d.db.Close()
    }
    return nil
}

// 注册
builder.Services.Add(NewDatabaseConnection)

// 应用关闭时自动调用 Dispose
app := builder.Build()
defer app.Dispose()  // 自动释放所有 IDisposable 服务
```

### 释放顺序

服务按注册顺序的逆序释放（LIFO），确保依赖关系正确：

```
注册顺序：Repository → Service → Controller
释放顺序：Controller → Service → Repository
```

## 最佳实践

### 1. 优先使用构造函数注入

```go
// ✅ 推荐：依赖通过构造函数注入
type UserService struct {
    repo   *UserRepository
    logger logging.ILogger
}

func NewUserService(repo *UserRepository, logger logging.ILogger) *UserService {
    return &UserService{
        repo:   repo,
        logger: logger,
    }
}

// ❌ 不推荐：在方法中手动获取依赖
type UserService struct {
    provider di.IServiceProvider
}

func (s *UserService) GetUser(id int) {
    repo := di.Get[*UserRepository](s.provider)  // 不推荐
}
```

### 2. 使用接口定义依赖

```go
// 定义接口
type IUserRepository interface {
    FindByID(id int) (*User, error)
    Save(user *User) error
}

// 实现接口
type UserRepository struct{}

func (r *UserRepository) FindByID(id int) (*User, error) {
    return nil, nil
}

func (r *UserRepository) Save(user *User) error {
    return nil
}

// 依赖接口而非具体实现
type UserService struct {
    repo IUserRepository  // 使用接口
}
```

### 3. 服务注册集中管理

创建模块化的服务注册函数：

```go
// users/services.go
package users

func AddUserServices(services di.IServiceCollection) {
    services.Add(NewUserRepository)
    services.Add(NewUserService)
    services.Add(NewUserValidator)
}

// main.go
package main

func main() {
    builder := web.CreateBuilder()
    
    // 模块化注册
    users.AddUserServices(builder.Services)
    orders.AddOrderServices(builder.Services)
    auth.AddAuthServices(builder.Services)
    
    app := builder.Build()
    app.Run()
}
```

### 4. 使用泛型 API

```go
// ✅ 推荐：类型安全的泛型 API
userService := di.Get[*UserService](provider)

// ⚠️ 可选：指针填充方式
var userService *UserService
provider.Get(&userService)
```

### 5. 合理使用命名服务

```go
// ✅ 好的使用场景：多个相同类型的不同实例
services.AddNamed("mysql", NewMySQLDatabase)
services.AddNamed("postgres", NewPostgresDatabase)
services.AddNamed("redis", NewRedisDatabase)

// ❌ 不好的使用：应该使用不同的类型
services.AddNamed("user", NewUserService)
services.AddNamed("order", NewOrderService)
// 应该定义不同的服务类型
```

### 6. 避免服务定位器反模式

```go
// ❌ 服务定位器反模式（不推荐）
type Service struct {
    provider di.IServiceProvider  // 持有 provider
}

func (s *Service) DoWork() {
    // 在方法中解析依赖
    repo := di.Get[*Repository](s.provider)
    logger := di.Get[*Logger](s.provider)
}

// ✅ 依赖注入（推荐）
type Service struct {
    repo   *Repository
    logger *Logger
}

func NewService(repo *Repository, logger *Logger) *Service {
    return &Service{
        repo:   repo,
        logger: logger,
    }
}
```

### 7. 单例服务注意线程安全

```go
// ✅ 线程安全的单例服务
type CacheService struct {
    data map[string]interface{}
    mu   sync.RWMutex
}

func (c *CacheService) Get(key string) interface{} {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.data[key]
}

func (c *CacheService) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}
```

## API 参考

### 泛型 API

```go
// 获取服务（不存在会 panic）
func Get[T any](provider IServiceProvider) T

// 获取服务（不存在返回默认值）
func GetOr[T any](provider IServiceProvider, defaultValue T) T

// 尝试获取服务
func TryGet[T any](provider IServiceProvider) (T, bool)

// 获取命名服务
func GetNamed[T any](provider IServiceProvider, name string) T

// 获取所有服务
func GetAll[T any](provider IServiceProvider) []T
```

### 服务注册

```go
// 注册单例服务
Add(constructor interface{}) IServiceCollection

// 注册实例
AddInstance(instance interface{}) IServiceCollection

// 注册命名服务
AddNamed(name string, constructor interface{}) IServiceCollection

// 条件注册
TryAdd(constructor interface{}) IServiceCollection

// 注册后台服务
AddHostedService(constructor interface{}) IServiceCollection
```

### 服务解析

```go
// 获取服务（指针填充）
Get(target interface{})

// 获取命名服务（指针填充）
GetNamed(target interface{}, serviceKey string)

// 释放资源
Dispose() error
```

## 常见问题

### 如何注册接口？

Go 的反射机制不能直接注册接口类型，需要注册具体实现：

```go
// 定义接口
type IUserService interface {
    GetUser(id int) (*User, error)
}

// 实现接口
type UserService struct {}

func (s *UserService) GetUser(id int) (*User, error) {
    return nil, nil
}

// ✅ 注册具体实现
services.Add(func() *UserService {
    return &UserService{}
})

// 使用时转换为接口
userService := di.Get[*UserService](provider)
var svc IUserService = userService  // 赋值给接口变量
```

### 如何处理初始化失败？

构造函数可以返回 error：

```go
services.Add(func() (*Database, error) {
    db, err := sql.Open("postgres", "...")
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    return &Database{db: db}, nil
})

// Build 时如果初始化失败会 panic
app := builder.Build()  // 可能 panic
```

### 服务什么时候创建？

Singleton 服务在首次请求时创建（懒加载），之后复用同一实例。

```go
services.Add(NewUserService)  // 注册，不创建
provider := services.Build()   // 构建，不创建
svc := di.Get[*UserService](provider)  // 首次获取时创建
svc2 := di.Get[*UserService](provider) // 复用已创建的实例
```

### 可以运行时动态注册服务吗？

不可以。所有服务必须在调用 `Build()` 之前注册完成。Build 后服务注册就被冻结了。

```go
services.Add(NewService1)
provider := services.Build()  // Build 后不能再注册

services.Add(NewService2)  // ❌ 无效，Build 后的注册被忽略
```

---

[← 返回主目录](../README.md)

