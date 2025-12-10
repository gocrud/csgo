# 配置依赖注入使用指南

## 新的配置方式（推荐）

### 1. 基本用法

```go
package admin

import (
    "github.com/gocrud/csgo/configuration"
    "github.com/gocrud/csgo/web"
    "vertical_slice_demo/configs"
)

func Bootstrap() *web.WebApplication {
    builder := web.CreateBuilder()
    
    // 构建配置
    builder.Configuration.
        AddJsonFile("configs/config.dev.json", true, false).
        AddEnvironmentVariables("APP_").
        Build()

    // 一步完成：注册配置到 DI 容器
    // 自动注册 IOptions[T]、IOptionsMonitor[T]、IOptionsSnapshot[T]
    configuration.Configure[configs.Config](builder.Services, "")
    
    // ... 其他服务注册
    
    return builder.Build()
}
```

### 2. 在服务中注入配置

```go
package myservice

import (
    "fmt"
    "github.com/gocrud/csgo/configuration"
    "vertical_slice_demo/configs"
)

// 使用 IOptions[T] - 静态配置（单例，不支持热更新）
type TestService struct {
    config configuration.IOptions[configs.Config]
}

func NewTestService(config configuration.IOptions[configs.Config]) *TestService {
    return &TestService{config: config}
}

func (s *TestService) Print() {
    cfg := s.config.Value()
    fmt.Println(cfg.Database.Database)
}

// 使用 IOptionsMonitor[T] - 支持热更新的配置（单例）
type MonitorService struct {
    config configuration.IOptionsMonitor[configs.Config]
}

func NewMonitorService(config configuration.IOptionsMonitor[configs.Config]) *MonitorService {
    svc := &MonitorService{config: config}
    
    // 监听配置变化
    config.OnChange(func(newConfig *configs.Config, name string) {
        fmt.Println("配置已更新:", newConfig.Database.Host)
    })
    
    return svc
}

func (s *MonitorService) GetCurrentConfig() *configs.Config {
    return s.config.CurrentValue()
}

// 使用 IOptionsSnapshot[T] - 请求作用域配置（瞬态）
type RequestService struct {
    config configuration.IOptionsSnapshot[configs.Config]
}

func NewRequestService(config configuration.IOptionsSnapshot[configs.Config]) *RequestService {
    return &RequestService{config: config}
}

func (s *RequestService) Process() {
    // 每次请求都会获取最新的配置快照
    cfg := s.config.Value()
    fmt.Println(cfg.Server.AdminPort)
}
```

### 3. 绑定配置子节点

```go
// 只绑定 Database 配置
configuration.Configure[configs.DatabaseConfig](builder.Services, "Database")

// 使用时
type DbService struct {
    dbConfig configuration.IOptions[configs.DatabaseConfig]
}

func NewDbService(dbConfig configuration.IOptions[configs.DatabaseConfig]) *DbService {
    return &DbService{dbConfig: dbConfig}
}

func (s *DbService) Connect() {
    cfg := s.dbConfig.Value()
    // 使用 cfg.Host, cfg.Port, etc.
}
```

### 4. 配置验证

```go
// 注册带验证的配置
configuration.ConfigureWithValidation[configs.EmailSettings](
    builder.Services,
    "Email",
    func(opts *configs.EmailSettings) error {
        if opts.SmtpHost == "" {
            return fmt.Errorf("SMTP host is required")
        }
        if opts.Port <= 0 {
            return fmt.Errorf("Port must be positive")
        }
        return nil
    },
)
```

### 5. 配置默认值

```go
// 注册带默认值的配置
configuration.ConfigureWithDefaults[configs.CacheConfig](
    builder.Services,
    "Cache",
    func() *configs.CacheConfig {
        return &configs.CacheConfig{
            Host:     "localhost",
            Port:     6379,
            Password: "",
        }
    },
)
```

## 和 .NET 的对比

| 功能 | .NET | Go (csgo) |
|------|------|-----------|
| 基本配置注册 | `services.Configure<T>(config)` | `configuration.Configure[T](services, "")` |
| 子节点配置 | `services.Configure<T>(config.GetSection("Database"))` | `configuration.Configure[T](services, "Database")` |
| 静态配置注入 | `IOptions<T>` | `configuration.IOptions[T]` |
| 热更新配置 | `IOptionsMonitor<T>` | `configuration.IOptionsMonitor[T]` |
| 请求作用域配置 | `IOptionsSnapshot<T>` | `configuration.IOptionsSnapshot[T]` |
| 带验证配置 | `services.AddOptions<T>().Validate(...)` | `configuration.ConfigureWithValidation[T](...)` |
| 带默认值配置 | 自行实现 | `configuration.ConfigureWithDefaults[T](...)` |

## 旧方式对比（已弃用）

### 旧方式（不推荐）
```go
// 需要两步操作
var appConfig configs.Config
builder.Configuration.Bind("", &appConfig)

// 然后手动注册
configuration.Configure[configs.Config](services, config, "")
```

### 新方式（推荐）
```go
// 一步到位
configuration.Configure[configs.Config](builder.Services, "")
```

## 核心优势

1. **一步到位**：`Configure` 同时注册 IOptions、IOptionsMonitor、IOptionsSnapshot
2. **参数简化**：不需要传递 config 参数，自动从 services 获取
3. **符合 .NET 习惯**：`configuration.Configure<T>(services, section)` 风格
4. **自动注册**：WebApplicationBuilder 自动注册 IConfiguration 到 DI 容器
5. **完整支持**：支持静态配置、热更新、请求作用域配置

## 注意事项

1. **IConfiguration 自动注册**：`WebApplicationBuilder.Build()` 会自动将 `IConfiguration` 注册到 DI 容器
2. **配置顺序**：必须先调用 `builder.Configuration.Build()` 再调用 `configuration.Configure[T]`
3. **热更新**：只有 `IOptionsMonitor[T]` 支持热更新，`IOptions[T]` 是静态的
4. **作用域**：
   - `IOptions[T]` - Singleton（单例，整个应用生命周期）
   - `IOptionsMonitor[T]` - Singleton（单例，但值可更新）
   - `IOptionsSnapshot[T]` - Transient（瞬态，每次解析都是新实例）
