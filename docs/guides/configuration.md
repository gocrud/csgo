# 配置管理指南

本指南介绍 CSGO 框架的配置管理系统，包括多源配置、类型安全的选项模式和动态配置更新。

## 目录

- [快速开始](#快速开始)
- [配置源](#配置源)
- [读取配置](#读取配置)
- [配置绑定](#配置绑定)
- [配置 API 速查](#配置-api-速查)
- [选项模式](#选项模式)
- [动态配置](#动态配置)
- [配置分节](#配置分节)
- [最佳实践](#最佳实践)

## 快速开始

### 基本示例

```go
package main

import (
    "fmt"
    "github.com/gocrud/csgo/configuration"
)

func main() {
    // 创建配置构建器
    builder := configuration.NewConfigurationBuilder()
    
    // 添加配置源
    builder.
        AddJsonFile("appsettings.json", false, false).
        AddEnvironmentVariables("APP_")
    
    // 构建配置
    config := builder.Build()
    
    // 读取配置值
    dbHost := config.Get("Database:Host")
    fmt.Println("Database Host:", dbHost)
}
```

### 与 Web 应用集成

Web 应用会自动配置配置系统：

```go
package main

import "github.com/gocrud/csgo/web"

func main() {
    // CreateBuilder 会自动设置配置
    builder := web.CreateBuilder()
    
    // 配置已经可用
    dbHost := builder.Configuration.Get("Database:Host")
    
    app := builder.Build()
    app.Run()
}
```

## 配置源

CSGO 支持多种配置源，按优先级顺序应用（后面的覆盖前面的）。

### JSON 文件

```go
builder := configuration.NewConfigurationBuilder()

// 必需的配置文件
builder.AddJsonFile("appsettings.json", false, false)

// 可选的环境特定配置
builder.AddJsonFile("appsettings.Development.json", true, false)

// 支持热重载
builder.AddJsonFile("appsettings.json", false, true)
```

**appsettings.json 示例：**

```json
{
  "Database": {
    "Host": "localhost",
    "Port": 5432,
    "Name": "mydb"
  },
  "Logging": {
    "Level": "Information"
  },
  "Features": {
    "EnableCaching": true,
    "CacheExpiration": 3600
  }
}
```

### YAML 文件

```go
builder.AddYamlFile("config.yaml", false, false)
```

**config.yaml 示例：**

```yaml
database:
  host: localhost
  port: 5432
  name: mydb

logging:
  level: Information
```

### 环境变量

```go
// 加载所有环境变量
builder.AddEnvironmentVariables("")

// 仅加载特定前缀的变量
builder.AddEnvironmentVariables("APP_")
```

**环境变量格式：**

```bash
# 使用 : 或 __ 作为层级分隔符
export Database:Host=localhost
export Database__Port=5432

# 使用前缀
export APP_Database__Host=localhost
export APP_Database__Port=5432
```

### 命令行参数

```go
import "os"

builder.AddCommandLine(os.Args[1:])
```

**命令行格式：**

```bash
# 两种格式都支持
./myapp --Database:Host=localhost --Database:Port=5432
./myapp --Database:Host localhost --Database:Port 5432
```

### 内存集合

用于测试或默认值：

```go
defaults := map[string]string{
    "Database:Host": "localhost",
    "Database:Port": "5432",
    "Timeout": "30",
}

builder.AddInMemoryCollection(defaults)
```

### 默认配置顺序

Web 应用使用以下默认顺序：

```go
builder := web.CreateBuilder(args...)

// 内部配置顺序：
// 1. appsettings.json
// 2. appsettings.{Environment}.json
// 3. 环境变量
// 4. 命令行参数
```

## 读取配置

### 简单值

```go
// 获取字符串值
host := config.Get("Database:Host")

// 配置键使用 : 分隔层级
logLevel := config.Get("Logging:Level")

// 不存在的键返回空字符串
value := config.Get("NonExistent") // ""
```

### 配置分节

```go
// 获取配置分节
dbSection := config.GetSection("Database")

// 读取分节中的值
host := dbSection.Get("Host")
port := dbSection.Get("Port")

// 分节支持嵌套
loggingSection := config.GetSection("Logging")
consoleSection := loggingSection.GetSection("Console")
level := consoleSection.Get("Level")
```

### 遍历子节点

```go
// 获取所有子分节
children := config.GetChildren()

for _, child := range children {
    fmt.Printf("Key: %s, Path: %s, Value: %s\n", 
        child.Key(), child.Path(), child.Value())
}
```

## 配置绑定

将配置绑定到强类型结构体：

### 定义配置结构

```go
type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Name     string `json:"name"`
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoggingConfig struct {
    Level      string `json:"level"`
    OutputPath string `json:"outputPath"`
}

type AppConfig struct {
    Database DatabaseConfig `json:"database"`
    Logging  LoggingConfig  `json:"logging"`
}
```

### 绑定配置

```go
// 方式 1: 直接使用 Bind
var appConfig AppConfig
if err := config.Bind("", &appConfig); err != nil {
    panic(err)
}

// 方式 2: 使用泛型辅助函数 BindOptions（直接返回值）
dbConfig, err := configuration.BindOptions[DatabaseConfig](config, "Database")
if err != nil {
    panic(err)
}

// 方式 3: 使用 MustBindOptions（失败时 panic，直接返回值）
cacheConfig := configuration.MustBindOptions[CacheConfig](config, "Cache")

// 使用绑定的配置
fmt.Printf("Connecting to %s:%d\n", dbConfig.Host, dbConfig.Port)
```

## 配置 API 速查

CSGO 提供多种配置 API，根据场景选择最合适的方式：

| API | 用途 | 说明 |
|-----|------|------|
| `Configure[T]()` | 注册配置到 DI | 注册 IOptions 和 IOptionsMonitor，支持热重载 |
| `ConfigureWithDefaults[T]()` | 带默认值注册 | 先应用默认值，再用配置覆盖 |
| `ConfigureWithValidation[T]()` | 带验证注册 | 启动时验证配置，热重载时也验证 |
| `BindOptions[T]()` | 手动绑定 | 直接返回 `*T`，无需传指针 |
| `MustBindOptions[T]()` | 手动绑定(panic) | 失败时 panic，直接返回 `*T` |

```go
// 推荐：使用 Configure 系列函数
configuration.Configure[DatabaseOptions](services, config, "Database")
configuration.ConfigureWithDefaults[CacheOptions](services, config, "Cache", defaultFunc)
configuration.ConfigureWithValidation[EmailOptions](services, config, "Email", validator)

// 手动绑定（不使用 DI 时）- 直接返回值
opts := configuration.MustBindOptions[AppSettings](config, "App")
```

## 选项模式

选项模式是 .NET 风格的类型安全配置方式。

### 定义选项类

```go
type DatabaseOptions struct {
    Host            string
    Port            int
    Name            string
    Username        string
    Password        string
    ConnectionRetry int
    MaxConnections  int
}

type CacheOptions struct {
    EnableCaching   bool
    CacheExpiration int
    Provider        string
}
```

### 注册选项

```go
import "github.com/gocrud/csgo/configuration"

// 方式 1: 使用 Configure 函数（推荐）
configuration.Configure[DatabaseOptions](builder.Services, builder.Configuration, "Database")

// 方式 2: 手动注册（更灵活）
builder.Services.AddSingleton(func() configuration.IOptions[DatabaseOptions] {
    var opts DatabaseOptions
    builder.Configuration.Bind("Database", &opts)
    return configuration.NewOptions(&opts)
})
```

`Configure[T]()` 会同时注册 `IOptions[T]` 和 `IOptionsMonitor[T]`，并自动支持配置热重载。

### 使用选项

```go
type UserService struct {
    dbOptions configuration.IOptions[DatabaseOptions]
}

func NewUserService(
    dbOptions configuration.IOptions[DatabaseOptions],
) *UserService {
    return &UserService{
        dbOptions: dbOptions,
    }
}

func (s *UserService) Connect() error {
    opts := s.dbOptions.Value()
    
    connectionString := fmt.Sprintf(
        "host=%s port=%d dbname=%s user=%s password=%s",
        opts.Host, opts.Port, opts.Name,
        opts.Username, opts.Password,
    )
    
    // 使用连接字符串连接数据库
    return nil
}
```

### 选项验证

```go
type EmailOptions struct {
    SmtpHost     string
    SmtpPort     int
    SenderEmail  string
    SenderName   string
}

// 定义验证函数
func validateEmailOptions(opts *EmailOptions) error {
    if opts.SmtpHost == "" {
        return fmt.Errorf("SMTP host is required")
    }
    if opts.SmtpPort <= 0 || opts.SmtpPort > 65535 {
        return fmt.Errorf("invalid SMTP port: %d", opts.SmtpPort)
    }
    if opts.SenderEmail == "" {
        return fmt.Errorf("sender email is required")
    }
    return nil
}

// 使用 ConfigureWithValidation（推荐）
err := configuration.ConfigureWithValidation[EmailOptions](
    builder.Services, 
    builder.Configuration, 
    "Email", 
    validateEmailOptions,
)
if err != nil {
    panic(fmt.Sprintf("Invalid email configuration: %v", err))
}
```

`ConfigureWithValidation[T]()` 会在注册时验证配置，并在配置热重载时再次验证（无效配置不会被应用）。

## 动态配置

### 配置更改监听

```go
// 创建选项监视器
monitor := configuration.NewOptionsMonitor(&DatabaseOptions{
    Host: "localhost",
    Port: 5432,
})

// 注册更改监听器
monitor.OnChange(func(opts *DatabaseOptions, name string) {
    fmt.Printf("Database configuration changed: %s:%d\n", 
        opts.Host, opts.Port)
    
    // 重新连接数据库或更新连接池
})

// 获取当前值
currentOpts := monitor.CurrentValue()

// 更新配置（触发监听器）
monitor.Set(&DatabaseOptions{
    Host: "newhost",
    Port: 3306,
})
```

### 热重载支持

```go
// 配置文件热重载
builder := configuration.NewConfigurationBuilder()
builder.AddJsonFile("appsettings.json", false, true) // reloadOnChange = true

config := builder.Build()

// 注册配置更改回调
config.OnChange(func() {
    fmt.Println("Configuration reloaded!")
    
    // 重新读取配置
    dbHost := config.Get("Database:Host")
    // 应用新配置...
})
```

## 配置分节

### 层级配置

```json
{
  "Services": {
    "UserService": {
      "Endpoint": "https://api.example.com/users",
      "Timeout": 30,
      "Retry": {
        "MaxAttempts": 3,
        "BackoffMs": 1000
      }
    },
    "OrderService": {
      "Endpoint": "https://api.example.com/orders",
      "Timeout": 60
    }
  }
}
```

```go
// 获取服务配置分节
servicesSection := config.GetSection("Services")

// 获取特定服务配置
userServiceSection := servicesSection.GetSection("UserService")
endpoint := userServiceSection.Get("Endpoint")

// 获取嵌套配置
retrySection := userServiceSection.GetSection("Retry")
maxAttempts := retrySection.Get("MaxAttempts")

// 或使用完整路径
maxAttempts := config.Get("Services:UserService:Retry:MaxAttempts")
```

### 数组配置

```json
{
  "AllowedHosts": ["localhost", "example.com", "*.myapp.com"],
  "Servers": [
    { "Name": "Server1", "Url": "http://server1" },
    { "Name": "Server2", "Url": "http://server2" }
  ]
}
```

```go
type ServerConfig struct {
    Name string `json:"name"`
    Url  string `json:"url"`
}

// 绑定数组配置
var servers []ServerConfig
config.Bind("Servers", &servers)

for _, server := range servers {
    fmt.Printf("Server: %s at %s\n", server.Name, server.Url)
}
```

## 环境特定配置

### 按环境加载配置

```go
import "github.com/gocrud/csgo/hosting"

// 获取当前环境
env := hosting.NewEnvironment()

// 构建配置
builder := configuration.NewConfigurationBuilder()
builder.AddJsonFile("appsettings.json", false, false)
builder.AddJsonFile(
    fmt.Sprintf("appsettings.%s.json", env.Name()), 
    true, // 可选
    false,
)

config := builder.Build()
```

**文件结构：**

```
project/
├── appsettings.json              # 基础配置
├── appsettings.Development.json  # 开发环境
├── appsettings.Staging.json      # 预发布环境
└── appsettings.Production.json   # 生产环境
```

### 环境变量控制

```bash
# 设置环境
export CSGO_ENVIRONMENT=Production

# 或通过命令行
./myapp --environment=Production
```

## 最佳实践

### 1. 配置层次结构

遵循清晰的配置层次结构：

```json
{
  "App": {
    "Name": "MyApp",
    "Version": "1.0.0"
  },
  "Database": {
    "Primary": { "Host": "...", "Port": 5432 },
    "Cache": { "Host": "...", "Port": 6379 }
  },
  "Services": {
    "Authentication": { "Endpoint": "...", "Timeout": 30 },
    "Payment": { "Endpoint": "...", "Timeout": 60 }
  },
  "Features": {
    "EnableCaching": true,
    "EnableMetrics": false
  }
}
```

### 2. 使用强类型选项

始终使用选项模式而不是直接读取配置字符串：

**❌ 不推荐：**

```go
func (s *Service) Process() {
    host := s.config.Get("Database:Host")
    portStr := s.config.Get("Database:Port")
    port, _ := strconv.Atoi(portStr) // 可能失败
}
```

**✅ 推荐：**

```go
type Service struct {
    dbOptions configuration.IOptions[DatabaseOptions]
}

func (s *Service) Process() {
    opts := s.dbOptions.Value()
    // opts.Host 和 opts.Port 已经是正确类型
}
```

### 3. 敏感信息处理

不要在配置文件中存储敏感信息：

**❌ 不推荐：**

```json
{
  "Database": {
    "Password": "MySecretPassword123"
  }
}
```

**✅ 推荐：**

```bash
# 使用环境变量
export Database__Password=MySecretPassword123

# 或使用密钥管理服务
export Database__Password=$(aws ssm get-parameter --name /myapp/db-password --query Parameter.Value)
```

### 4. 配置验证

在应用启动时验证配置：

```go
func ValidateConfiguration(config configuration.IConfiguration) error {
    required := []string{
        "Database:Host",
        "Database:Port",
        "Database:Name",
    }
    
    for _, key := range required {
        if config.Get(key) == "" {
            return fmt.Errorf("required configuration missing: %s", key)
        }
    }
    
    return nil
}

func main() {
    builder := web.CreateBuilder()
    
    if err := ValidateConfiguration(builder.Configuration); err != nil {
        panic(fmt.Sprintf("Invalid configuration: %v", err))
    }
    
    app := builder.Build()
    app.Run()
}
```

### 5. 默认值

为配置提供合理的默认值：

```go
type AppOptions struct {
    Timeout         int
    MaxRetries      int
    EnableCaching   bool
}

// 方式 1: 使用 ConfigureWithDefaults（推荐）
configuration.ConfigureWithDefaults[AppOptions](
    builder.Services, 
    builder.Configuration, 
    "App",
    func() *AppOptions {
        return &AppOptions{
            Timeout:       30,
            MaxRetries:    3,
            EnableCaching: true,
        }
    },
)

// 方式 2: 手动处理默认值
func LoadAppOptions(config configuration.IConfiguration) *AppOptions {
    opts := &AppOptions{
        Timeout:       30,
        MaxRetries:    3,
        EnableCaching: true,
    }
    config.Bind("App", opts)
    return opts
}
```

`ConfigureWithDefaults[T]()` 会先应用默认值，然后用配置文件中的值覆盖。

### 6. 配置文档化

在代码中注释配置选项：

```go
type ApiOptions struct {
    // 服务端点 URL
    // 示例: "https://api.example.com"
    // 必需: 是
    Endpoint string `json:"endpoint"`
    
    // 请求超时时间（秒）
    // 默认: 30
    // 范围: 1-300
    Timeout int `json:"timeout"`
    
    // 是否启用请求重试
    // 默认: true
    EnableRetry bool `json:"enableRetry"`
    
    // 最大重试次数（仅在 EnableRetry=true 时有效）
    // 默认: 3
    // 范围: 1-10
    MaxRetries int `json:"maxRetries"`
}
```

## 完整示例

### 多层配置应用

```go
package main

import (
    "fmt"
    "github.com/gocrud/csgo/configuration"
    "github.com/gocrud/csgo/web"
)

// 定义配置结构
type DatabaseOptions struct {
    Host           string `json:"host"`
    Port           int    `json:"port"`
    Name           string `json:"name"`
    MaxConnections int    `json:"maxConnections"`
}

type CacheOptions struct {
    Enabled    bool   `json:"enabled"`
    Provider   string `json:"provider"`
    Expiration int    `json:"expiration"`
}

type AppOptions struct {
    Name    string `json:"name"`
    Version string `json:"version"`
}

// 服务使用配置
type DataService struct {
    dbOpts    configuration.IOptions[DatabaseOptions]
    cacheOpts configuration.IOptions[CacheOptions]
}

func NewDataService(
    dbOpts configuration.IOptions[DatabaseOptions],
    cacheOpts configuration.IOptions[CacheOptions],
) *DataService {
    return &DataService{
        dbOpts:    dbOpts,
        cacheOpts: cacheOpts,
    }
}

func (s *DataService) Initialize() error {
    db := s.dbOpts.Value()
    fmt.Printf("Connecting to database: %s:%d/%s\n", 
        db.Host, db.Port, db.Name)
    
    if s.cacheOpts.Value().Enabled {
        cache := s.cacheOpts.Value()
        fmt.Printf("Enabling cache: provider=%s, expiration=%ds\n",
            cache.Provider, cache.Expiration)
    }
    
    return nil
}

func main() {
    // 创建应用
    builder := web.CreateBuilder()
    
    // 注册选项（使用 Configure 函数，自动支持热重载）
    configuration.Configure[DatabaseOptions](builder.Services, builder.Configuration, "Database")
    configuration.Configure[CacheOptions](builder.Services, builder.Configuration, "Cache")
    
    // 注册服务
    builder.Services.AddSingleton(NewDataService)
    
    app := builder.Build()
    
    // 初始化服务
    var dataService *DataService
    app.Services.GetRequiredService(&dataService)
    dataService.Initialize()
    
    app.Run()
}
```

**appsettings.json：**

```json
{
  "App": {
    "Name": "MyApp",
    "Version": "1.0.0"
  },
  "Database": {
    "Host": "localhost",
    "Port": 5432,
    "Name": "mydb",
    "MaxConnections": 100
  },
  "Cache": {
    "Enabled": true,
    "Provider": "Redis",
    "Expiration": 3600
  }
}
```

## 与 .NET 对比

| .NET | CSGO | 说明 |
|------|-----|------|
| `IConfiguration` | `configuration.IConfiguration` | 配置接口 |
| `IConfigurationBuilder` | `configuration.IConfigurationBuilder` | 配置构建器 |
| `AddJsonFile()` | `AddJsonFile()` | 添加 JSON 配置源 |
| `AddEnvironmentVariables()` | `AddEnvironmentVariables()` | 添加环境变量 |
| `GetSection()` | `GetSection()` | 获取配置分节 |
| `Bind()` | `Bind()` | 绑定到对象 |
| `IOptions<T>` | `configuration.IOptions[T]` | 选项接口 |
| `IOptionsMonitor<T>` | `configuration.IOptionsMonitor[T]` | 选项监视器 |

## 相关资源

- [快速开始](../getting-started.md) - 开始使用 CSGO
- [应用托管](hosting.md) - 应用生命周期管理
- [依赖注入](dependency-injection.md) - DI 系统详解
- [最佳实践](../best-practices.md) - 推荐模式

---

**下一步**: 查看 [应用托管指南](hosting.md) 了解如何管理应用生命周期。

