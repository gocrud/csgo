# 配置管理 (Configuration)

[← 返回主目录](../README.md)

CSGO 的配置系统提供了强大而灵活的配置管理功能，支持多种配置源、配置绑定、热更新等特性。

## 特性

- ✅ 多种配置源（JSON、YAML、INI、XML、环境变量、命令行）
- ✅ 配置分层和覆盖
- ✅ 配置绑定到结构体
- ✅ 配置热更新（文件监控）
- ✅ 配置节点和层级结构
- ✅ Options 模式支持
- ✅ 类型安全的配置读取
- ✅ 开发/生产环境配置

## 快速开始

### 💡 推荐方式：Configure[T] 模式

**这是最简单、最强大的配置管理方式！**

```go
package main

import (
    "github.com/gocrud/csgo/configuration"
    "github.com/gocrud/csgo/web"
)

// 1. 定义配置结构
type AppSettings struct {
    Name    string `json:"name"`
    Version string `json:"version"`
}

type ServerSettings struct {
    Port int    `json:"port"`
    Host string `json:"host"`
}

func main() {
    builder := web.CreateBuilder()
    
    // 2. 注册配置（一行代码搞定！）
    configuration.Configure[AppSettings](builder.Services, "app")
    configuration.Configure[ServerSettings](builder.Services, "server")
    
    // 3. 注册服务（自动注入配置）
    builder.Services.Add(NewMyService)
    
    app := builder.Build()
    app.Run()
}

// 4. 在服务中使用配置
type MyService struct {
    appConfig    configuration.IOptionsMonitor[AppSettings]
    serverConfig configuration.IOptionsMonitor[ServerSettings]
}

func NewMyService(
    appConfig configuration.IOptionsMonitor[AppSettings],
    serverConfig configuration.IOptionsMonitor[ServerSettings],
) *MyService {
    return &MyService{
        appConfig:    appConfig,
        serverConfig: serverConfig,
    }
}

func (s *MyService) DoWork() {
    // 自动获取最新配置（支持热更新）
    app := s.appConfig.CurrentValue()
    server := s.serverConfig.CurrentValue()
    
    fmt.Printf("App: %s v%s\n", app.Name, app.Version)
    fmt.Printf("Server: %s:%d\n", server.Host, server.Port)
}
```

**配置文件（appsettings.json）：**

```json
{
  "app": {
    "name": "MyApp",
    "version": "1.0.0"
  },
  "server": {
    "port": 8080,
    "host": "localhost"
  }
}
```

**✨ 优势：**
- ✅ 类型安全
- ✅ 自动热更新
- ✅ 依赖注入集成
- ✅ 一行代码注册
- ✅ 无需手动绑定

### 传统方式：直接读取配置

如果只需要简单读取配置值：

```go
package main

import (
    "github.com/gocrud/csgo/web"
)

func main() {
    builder := web.CreateBuilder()
    
    // 直接读取配置（自动加载 appsettings.json）
    port := builder.Configuration.GetInt("server:port", 8080)
    dbConnection := builder.Configuration.Get("database:connection")
    
    app := builder.Build()
    app.Run()
}
```

**⚠️ 注意：** 这种方式不支持热更新和依赖注入，只适合简单场景。

### 2. 配置文件

创建 `appsettings.json`：

```json
{
  "server": {
    "port": 8080,
    "host": "localhost"
  },
  "database": {
    "connection": "postgres://localhost/mydb",
    "maxConnections": 10
  },
  "logging": {
    "level": "Information"
  }
}
```

### 3. 环境特定配置

创建 `appsettings.development.json`：

```json
{
  "logging": {
    "level": "Debug"
  },
  "database": {
    "connection": "postgres://localhost/mydb_dev"
  }
}
```

框架会自动根据环境加载对应的配置文件并覆盖基础配置。

## 配置源

### JSON 配置文件

最常用的配置格式，支持层级结构。

```go
builder.Configuration.AddJsonFile("appsettings.json", false, true)
// 参数：路径、是否可选、是否监控变化
```

**appsettings.json：**

```json
{
  "app": {
    "name": "MyApp",
    "version": "1.0.0"
  },
  "server": {
    "port": 8080,
    "timeout": 30
  }
}
```

### YAML 配置文件

支持 YAML 格式配置。

```go
builder.Configuration.AddYamlFile("config.yaml", false, true)
```

**config.yaml：**

```yaml
app:
  name: MyApp
  version: 1.0.0
server:
  port: 8080
  timeout: 30
```

### INI 配置文件

支持传统的 INI 格式。

```go
builder.Configuration.AddIniFile("config.ini", false, true)
```

**config.ini：**

```ini
[app]
name=MyApp
version=1.0.0

[server]
port=8080
timeout=30
```

### XML 配置文件

支持 XML 格式配置。

```go
builder.Configuration.AddXmlFile("config.xml", false, true)
```

**config.xml：**

```xml
<configuration>
  <app>
    <name>MyApp</name>
    <version>1.0.0</version>
  </app>
  <server>
    <port>8080</port>
    <timeout>30</timeout>
  </server>
</configuration>
```

### 环境变量

从环境变量读取配置。

```go
// 读取所有环境变量
builder.Configuration.AddEnvironmentVariables("")

// 只读取特定前缀的环境变量
builder.Configuration.AddEnvironmentVariables("MYAPP_")
```

**使用示例：**

```bash
# 设置环境变量
export MYAPP_SERVER__PORT=9000
export MYAPP_DATABASE__CONNECTION="postgres://..."

# 在Go中读取
# server__port -> server:port
port := builder.Configuration.GetInt("server:port", 8080)  // 9000
```

**注意**：环境变量使用双下划线（`__`）表示层级，在配置中转换为冒号（`:`）。

### 命令行参数

从命令行参数读取配置。

```go
builder.Configuration.AddCommandLine(os.Args[1:])
```

**使用示例：**

```bash
# 运行应用时传递参数
./myapp --server:port=9000 --database:connection="postgres://..."
```

### 内存配置

从内存中的 map 读取配置。

```go
config := map[string]string{
    "server:port": "8080",
    "server:host": "localhost",
}
builder.Configuration.AddInMemoryCollection(config)
```

### Key-Per-File 配置

适用于容器环境（如 Docker、Kubernetes），每个文件代表一个配置键。

```go
// 读取 /etc/secrets 目录下的所有文件
builder.Configuration.AddKeyPerFile("/etc/secrets", false)
```

**目录结构：**

```
/etc/secrets/
  ├── database/
  │   ├── username
  │   └── password
  └── api/
      └── key
```

**生成的配置：**

```
database:username = <username文件内容>
database:password = <password文件内容>
api:key = <key文件内容>
```

## 配置层级和覆盖

### 配置优先级

配置源按添加顺序应用，后添加的会覆盖先添加的：

```go
builder := web.CreateBuilder()

// 1. 基础配置
builder.Configuration.AddJsonFile("appsettings.json", false, true)

// 2. 环境配置（覆盖基础配置）
env := builder.Environment.GetEnvironmentName()
builder.Configuration.AddJsonFile(
    fmt.Sprintf("appsettings.%s.json", env), 
    true,  // 可选
    true,
)

// 3. 环境变量（覆盖文件配置）
builder.Configuration.AddEnvironmentVariables("")

// 4. 命令行（最高优先级）
builder.Configuration.AddCommandLine(os.Args[1:])
```

**优先级顺序（从低到高）：**

```
appsettings.json
  ↓ 覆盖
appsettings.{Environment}.json
  ↓ 覆盖
环境变量
  ↓ 覆盖
命令行参数
```

### 配置键的层级

配置使用冒号（`:`）分隔层级：

```json
{
  "database": {
    "primary": {
      "connection": "...",
      "timeout": 30
    }
  }
}
```

读取时使用冒号分隔：

```go
connection := config.Get("database:primary:connection")
timeout := config.GetInt("database:primary:timeout", 30)
```

## 配置绑定

### 绑定到结构体

将配置节点绑定到 Go 结构体，实现强类型配置。

```go
// 定义配置结构
type ServerConfig struct {
    Port    int    `json:"port"`
    Host    string `json:"host"`
    Timeout int    `json:"timeout"`
}

type DatabaseConfig struct {
    Connection     string `json:"connection"`
    MaxConnections int    `json:"maxConnections"`
}

// 绑定配置
var serverConfig ServerConfig
builder.Configuration.Bind("server", &serverConfig)

var dbConfig DatabaseConfig
builder.Configuration.Bind("database", &dbConfig)

// 使用配置
fmt.Printf("Server: %s:%d\n", serverConfig.Host, serverConfig.Port)
fmt.Printf("Database: %s\n", dbConfig.Connection)
```

### 嵌套结构绑定

支持嵌套结构的配置绑定。

```go
type AppConfig struct {
    Server struct {
        Port int    `json:"port"`
        Host string `json:"host"`
    } `json:"server"`
    
    Database struct {
        Connection string `json:"connection"`
        Pool       struct {
            MinSize int `json:"minSize"`
            MaxSize int `json:"maxSize"`
        } `json:"pool"`
    } `json:"database"`
}

var config AppConfig
builder.Configuration.Bind("", &config)  // 绑定根节点

// 访问配置
port := config.Server.Port
maxSize := config.Database.Pool.MaxSize
```

### 数组配置绑定

支持数组配置的绑定。

**配置文件：**

```json
{
  "servers": [
    {
      "name": "server1",
      "host": "192.168.1.1",
      "port": 8080
    },
    {
      "name": "server2",
      "host": "192.168.1.2",
      "port": 8081
    }
  ]
}
```

**绑定代码：**

```go
type Server struct {
    Name string `json:"name"`
    Host string `json:"host"`
    Port int    `json:"port"`
}

type Config struct {
    Servers []Server `json:"servers"`
}

var config Config
builder.Configuration.Bind("", &config)

// 使用
for _, server := range config.Servers {
    fmt.Printf("Server: %s at %s:%d\n", server.Name, server.Host, server.Port)
}
```

## Options 模式（推荐）

Options 模式是一种优雅的配置管理方式，提供类型安全、热更新支持和依赖注入集成。

### 方式 1：使用 Configure[T]（推荐）

**最简单、最优雅的方式**，自动处理配置绑定、热更新和依赖注入。

#### 1. 定义配置结构

```go
package config

type AppSettings struct {
    Name    string `json:"name"`
    Version string `json:"version"`
    Debug   bool   `json:"debug"`
}

type DatabaseSettings struct {
    Host           string `json:"host"`
    Port           int    `json:"port"`
    MaxConnections int    `json:"maxConnections"`
}
```

#### 2. 注册配置

```go
import "github.com/gocrud/csgo/configuration"

func main() {
    builder := web.CreateBuilder()
    
    // 方式 A：基本用法
    configuration.Configure[AppSettings](builder.Services, "app")
    configuration.Configure[DatabaseSettings](builder.Services, "database")
    
    // 方式 B：带默认值
    configuration.ConfigureWithDefaults[AppSettings](builder.Services, "app", func() *AppSettings {
        return &AppSettings{
            Name:    "MyApp",
            Version: "1.0.0",
            Debug:   false,
        }
    })
    
    // 方式 C：带验证
    configuration.ConfigureWithValidation[DatabaseSettings](builder.Services, "database", 
        func(opts *DatabaseSettings) error {
            if opts.Host == "" {
                return fmt.Errorf("database host is required")
            }
            if opts.Port <= 0 {
                return fmt.Errorf("invalid database port")
            }
            return nil
        },
    )
    
    app := builder.Build()
    app.Run()
}
```

#### 3. 配置文件

```json
{
  "app": {
    "name": "MyApp",
    "version": "1.0.0",
    "debug": true
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "maxConnections": 100
  }
}
```

#### 4. 在服务中使用配置

**风格 1：使用 IOptionsMonitor（推荐，支持热更新）**

```go
type UserService struct {
    dbConfig configuration.IOptionsMonitor[DatabaseSettings]
    logger   logging.ILogger
}

func NewUserService(
    dbConfig configuration.IOptionsMonitor[DatabaseSettings],
    logger logging.ILogger,
) *UserService {
    return &UserService{
        dbConfig: dbConfig,
        logger:   logger,
    }
}

func (s *UserService) Connect() error {
    // 获取当前配置（自动获取最新值）
    config := s.dbConfig.CurrentValue()
    
    connStr := fmt.Sprintf("%s:%d", config.Host, config.Port)
    s.logger.LogInformation("Connecting to database: %s", connStr)
    
    // 连接数据库...
    return nil
}

func (s *UserService) WatchConfigChanges() {
    // 监听配置变化
    s.dbConfig.OnChange(func(newConfig *DatabaseSettings, name string) {
        s.logger.LogInformation("Database config changed: %s:%d", 
            newConfig.Host, newConfig.Port)
        // 重新连接或更新连接池配置
    })
}
```

**风格 2：直接注入配置值（快照）**

```go
type EmailService struct {
    appConfig AppSettings
    logger    logging.ILogger
}

func NewEmailService(
    appConfig AppSettings,  // 直接注入配置值
    logger logging.ILogger,
) *EmailService {
    return &EmailService{
        appConfig: appConfig,
        logger:    logger,
    }
}

func (s *EmailService) SendWelcome(email string) error {
    subject := fmt.Sprintf("Welcome to %s", s.appConfig.Name)
    // 发送邮件...
    return nil
}
```

#### 5. 注册服务

```go
func main() {
    builder := web.CreateBuilder()
    
    // 注册配置
    configuration.Configure[AppSettings](builder.Services, "app")
    configuration.Configure[DatabaseSettings](builder.Services, "database")
    
    // 注册服务（自动注入配置）
    builder.Services.Add(NewUserService)    // 自动注入 IOptionsMonitor[DatabaseSettings]
    builder.Services.Add(NewEmailService)   // 自动注入 AppSettings
    
    app := builder.Build()
    app.Run()
}
```

### 方式 2：手动绑定（传统方式）

如果你需要更多控制，可以手动绑定配置。

#### 定义 Options

```go
type ServiceOptions struct {
    Timeout     time.Duration `json:"timeout"`
    MaxRetries  int          `json:"maxRetries"`
    EnableCache bool         `json:"enableCache"`
}

// 提供默认值
func NewServiceOptions() *ServiceOptions {
    return &ServiceOptions{
        Timeout:     30 * time.Second,
        MaxRetries:  3,
        EnableCache: true,
    }
}
```

#### 注册服务时绑定配置

```go
// 服务接受 Options
type MyService struct {
    options *ServiceOptions
}

func NewMyService(options *ServiceOptions) *MyService {
    return &MyService{options: options}
}

// 手动注册
builder.Services.Add(func(config configuration.IConfiguration) *MyService {
    options := NewServiceOptions()
    
    // 从配置绑定（覆盖默认值）
    if err := config.Bind("myService", options); err != nil {
        log.Printf("Failed to bind config: %v", err)
    }
    
    return NewMyService(options)
})
```

#### 配置文件

```json
{
  "myService": {
    "timeout": "60s",
    "maxRetries": 5,
    "enableCache": false
  }
}
```

### Configure vs 手动绑定

| 特性 | Configure[T] | 手动绑定 |
|------|-------------|---------|
| **代码简洁性** | ✅ 一行代码搞定 | ❌ 需要自己写工厂函数 |
| **类型安全** | ✅ 完全类型安全 | ✅ 类型安全 |
| **热更新支持** | ✅ 自动支持 | ❌ 需要手动实现 |
| **依赖注入** | ✅ 自动注册 | ❌ 需要手动注册 |
| **两种注入风格** | ✅ 支持 IOptionsMonitor[T] 和 T | ❌ 只支持一种 |
| **默认值** | ✅ ConfigureWithDefaults | ✅ 手动提供 |
| **验证** | ✅ ConfigureWithValidation | ❌ 需要手动验证 |

**推荐使用 `Configure[T]` 方式！**

## 配置热更新

### 启用文件监控

在添加配置源时启用 `reloadOnChange` 参数：

```go
// 启用热更新
builder.Configuration.AddJsonFile("appsettings.json", false, true)
//                                                            ↑
//                                                       reloadOnChange
```

### 方式 1：使用 IOptionsMonitor（推荐）

`Configure[T]` 自动支持配置热更新，无需额外代码！

```go
// 1. 定义配置
type CacheSettings struct {
    Timeout int `json:"timeout"`  // 秒
    MaxSize int `json:"maxSize"`
}

// 2. 注册配置（自动支持热更新）
configuration.Configure[CacheSettings](builder.Services, "cache")

// 3. 在服务中使用（自动获取最新配置）
type CacheService struct {
    config configuration.IOptionsMonitor[CacheSettings]
    logger logging.ILogger
}

func NewCacheService(
    config configuration.IOptionsMonitor[CacheSettings],
    logger logging.ILogger,
) *CacheService {
    svc := &CacheService{
        config: config,
        logger: logger,
    }
    
    // 可选：监听配置变化事件
    config.OnChange(func(newConfig *CacheSettings, name string) {
        svc.logger.LogInformation("Cache config changed: timeout=%d, maxSize=%d",
            newConfig.Timeout, newConfig.MaxSize)
        // 这里可以执行配置变化后的逻辑
        // 例如：重新初始化缓存、调整连接池大小等
    })
    
    return svc
}

func (s *CacheService) Get(key string) (interface{}, error) {
    // CurrentValue() 总是返回最新配置，无需手动刷新！
    config := s.config.CurrentValue()
    
    timeout := time.Duration(config.Timeout) * time.Second
    maxSize := config.MaxSize
    
    // 使用最新配置执行业务逻辑...
    s.logger.LogDebug("Using cache with timeout=%v, maxSize=%d", timeout, maxSize)
    
    return nil, nil
}

func (s *CacheService) Set(key string, value interface{}) error {
    // 每次调用都自动使用最新配置
    config := s.config.CurrentValue()
    
    if config.MaxSize <= 0 {
        return fmt.Errorf("cache is disabled")
    }
    
    // 设置缓存...
    return nil
}
```

**配置文件（appsettings.json）：**

```json
{
  "cache": {
    "timeout": 60,
    "maxSize": 1000
  }
}
```

**修改配置文件后，服务会自动使用新配置，无需重启应用！**

### 方式 2：监听 IConfiguration 变化（传统方式）

如果不使用 `Configure[T]`，可以手动监听配置变化：

```go
builder.Configuration.OnChange(func() {
    // 配置发生变化时调用
    fmt.Println("Configuration changed!")
    
    // 重新读取配置
    port := builder.Configuration.GetInt("server:port", 8080)
    fmt.Printf("New port: %d\n", port)
})
```

### 方式 3：手动实现热更新（不推荐）

需要手动管理配置刷新和线程安全：

```go
type CacheService struct {
    config configuration.IConfiguration
    mu     sync.RWMutex
    
    timeout time.Duration
    maxSize int
}

func NewCacheService(config configuration.IConfiguration) *CacheService {
    svc := &CacheService{
        config: config,
    }
    
    // 初始加载配置
    svc.reloadConfig()
    
    // 监听配置变化（手动刷新）
    config.OnChange(func() {
        svc.reloadConfig()
    })
    
    return svc
}

func (s *CacheService) reloadConfig() {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.timeout = time.Duration(s.config.GetInt("cache:timeout", 60)) * time.Second
    s.maxSize = s.config.GetInt("cache:maxSize", 1000)
    
    fmt.Printf("Cache config reloaded: timeout=%v, maxSize=%d\n", 
        s.timeout, s.maxSize)
}

func (s *CacheService) GetTimeout() time.Duration {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.timeout
}
```

### 热更新对比

| 方式 | IOptionsMonitor[T] | 手动监听 IConfiguration |
|------|-------------------|------------------------|
| **代码复杂度** | ✅ 简单 | ❌ 复杂（需要手动刷新） |
| **线程安全** | ✅ 自动处理 | ❌ 需要手动加锁 |
| **类型安全** | ✅ 完全类型安全 | ⚠️ 运行时转换 |
| **自动更新** | ✅ 自动获取最新值 | ❌ 需要手动刷新字段 |
| **变化通知** | ✅ OnChange 回调 | ✅ OnChange 回调 |

**强烈推荐使用 `IOptionsMonitor[T]` 方式！**

## 配置节点

### 获取配置节点

使用 `GetSection` 获取配置的子节点：

```go
// 获取 server 节点
serverSection := config.GetSection("server")

// 从节点读取值
port := serverSection.GetInt("port", 8080)
host := serverSection.GetString("host", "localhost")

// 获取节点路径
fmt.Println(serverSection.Path())  // "server"
fmt.Println(serverSection.Key())   // "server"
```

### 遍历子节点

```go
// 获取所有子节点
children := config.GetChildren()
for _, child := range children {
    fmt.Printf("Key: %s, Value: %s\n", child.Key(), child.Value())
}

// 获取特定节点的子节点
serverSection := config.GetSection("server")
serverChildren := serverSection.GetChildren()
```

### 检查配置是否存在

```go
if config.Exists("server:port") {
    port := config.GetInt("server:port", 8080)
    fmt.Printf("Port: %d\n", port)
} else {
    fmt.Println("Port not configured")
}
```

## 类型安全的读取

### 基本类型

```go
// 字符串
host := config.GetString("server:host", "localhost")

// 整数
port := config.GetInt("server:port", 8080)
maxConn := config.GetInt64("database:maxConnections", 100)

// 布尔值
enabled := config.GetBool("features:newUI", false)

// 浮点数
timeout := config.GetFloat64("server:timeout", 30.0)
```

### Get vs GetXxx

- `Get(key)`: 返回字符串，如果不存在返回空字符串
- `GetString(key, default)`: 返回字符串，如果不存在返回默认值
- `GetInt(key, default)`: 返回整数，如果不存在或转换失败返回默认值
- 其他类型同理

```go
// Get - 不存在返回空字符串
host := config.Get("server:host")  // ""

// GetString - 不存在返回默认值
host := config.GetString("server:host", "localhost")  // "localhost"
```

## 开发实践

### 项目配置结构

推荐的配置文件组织：

```
项目根目录/
  ├── appsettings.json              # 基础配置
  ├── appsettings.Development.json  # 开发环境配置
  ├── appsettings.Production.json   # 生产环境配置
  ├── appsettings.Staging.json      # 预发布环境配置
  └── main.go
```

### 配置结构设计

将配置组织成逻辑模块：

```json
{
  "app": {
    "name": "MyApp",
    "version": "1.0.0"
  },
  "server": {
    "port": 8080,
    "host": "localhost",
    "timeout": 30
  },
  "database": {
    "primary": {
      "connection": "postgres://...",
      "maxConnections": 10
    },
    "cache": {
      "connection": "redis://...",
      "ttl": 3600
    }
  },
  "logging": {
    "level": "Information",
    "console": {
      "enabled": true
    },
    "file": {
      "enabled": true,
      "path": "logs/app.log"
    }
  },
  "features": {
    "enableNewUI": false,
    "enableBeta": false
  }
}
```

### 敏感配置管理

**开发环境**：使用配置文件或环境变量

```bash
export DATABASE_PASSWORD="dev_password"
```

**生产环境**：使用密钥管理服务

```go
// 使用 Key-Per-File 读取 Kubernetes Secrets
builder.Configuration.AddKeyPerFile("/etc/secrets", false)

// 或使用环境变量
builder.Configuration.AddEnvironmentVariables("APP_")
```

**不要在配置文件中存储敏感信息**：

```json
{
  "database": {
    "connection": "postgres://user:PASSWORD@host/db"  // ❌ 不要这样
  }
}
```

使用环境变量或密钥文件：

```json
{
  "database": {
    "host": "localhost",
    "database": "mydb",
    "username": "user"
    // password 从环境变量或密钥文件读取
  }
}
```

```go
// 组合配置
dbHost := config.Get("database:host")
dbUser := config.Get("database:username")
dbPass := config.Get("database:password")  // 从环境变量读取

connStr := fmt.Sprintf("postgres://%s:%s@%s/%s", 
    dbUser, dbPass, dbHost, config.Get("database:database"))
```

## 最佳实践

### 1. 使用 Configure[T] 模式（强烈推荐）

```go
// ✅ 最佳实践：使用 Configure[T] 模式
import "github.com/gocrud/csgo/configuration"

// 定义配置结构
type ServerConfig struct {
    Port    int    `json:"port"`
    Host    string `json:"host"`
    Timeout int    `json:"timeout"`
}

type DatabaseConfig struct {
    Connection     string `json:"connection"`
    MaxConnections int    `json:"maxConnections"`
}

// 注册配置
configuration.Configure[ServerConfig](builder.Services, "server")
configuration.Configure[DatabaseConfig](builder.Services, "database")

// 在服务中使用
type MyService struct {
    serverConfig configuration.IOptionsMonitor[ServerConfig]
}

func NewMyService(serverConfig configuration.IOptionsMonitor[ServerConfig]) *MyService {
    return &MyService{serverConfig: serverConfig}
}

// ⚠️ 备选方案：手动绑定
var serverConfig ServerConfig
builder.Configuration.Bind("server", &serverConfig)
builder.Services.AddInstance(&serverConfig)

// ❌ 不推荐：到处使用字符串键
port := builder.Configuration.GetInt("server:port", 8080)
```

**为什么推荐 Configure[T]？**
- ✅ 类型安全
- ✅ 自动热更新
- ✅ 依赖注入集成
- ✅ 代码更简洁
- ✅ 支持配置验证

### 2. 提供默认值

```go
// ✅ 始终提供合理的默认值
port := config.GetInt("server:port", 8080)
timeout := config.GetInt("server:timeout", 30)

// ❌ 不提供默认值可能导致零值
port := config.GetInt("server:port", 0)  // 0 不是有效端口
```

### 3. 验证配置

在应用启动时验证关键配置：

```go
func validateConfig(config *AppConfig) error {
    if config.Server.Port < 1 || config.Server.Port > 65535 {
        return fmt.Errorf("invalid port: %d", config.Server.Port)
    }
    
    if config.Database.Connection == "" {
        return fmt.Errorf("database connection is required")
    }
    
    return nil
}

// 在 main 中使用
var config AppConfig
builder.Configuration.Bind("", &config)

if err := validateConfig(&config); err != nil {
    log.Fatalf("Invalid configuration: %v", err)
}
```

### 4. 集中配置管理

**方式 A：使用 Configure[T]（推荐）**

```go
// config/config.go
package config

// 定义各模块配置
type ServerConfig struct {
    Port    int    `json:"port"`
    Host    string `json:"host"`
    Timeout int    `json:"timeout"`
}

type DatabaseConfig struct {
    Connection     string `json:"connection"`
    MaxConnections int    `json:"maxConnections"`
}

type LoggingConfig struct {
    Level  string `json:"level"`
    Format string `json:"format"`
}

// 注册所有配置
func RegisterConfigs(services di.IServiceCollection) {
    // 使用 Configure 自动注册
    configuration.Configure[ServerConfig](services, "server")
    
    // 带默认值
    configuration.ConfigureWithDefaults[DatabaseConfig](services, "database", func() *DatabaseConfig {
        return &DatabaseConfig{
            Connection:     "localhost:5432",
            MaxConnections: 10,
        }
    })
    
    // 带验证
    configuration.ConfigureWithValidation[LoggingConfig](services, "logging",
        func(cfg *LoggingConfig) error {
            validLevels := []string{"Debug", "Information", "Warning", "Error"}
            for _, level := range validLevels {
                if cfg.Level == level {
                    return nil
                }
            }
            return fmt.Errorf("invalid log level: %s", cfg.Level)
        },
    )
}

// main.go
func main() {
    builder := web.CreateBuilder()
    
    // 注册所有配置
    config.RegisterConfigs(builder.Services)
    
    // 注册服务（自动注入配置）
    builder.Services.Add(NewUserService)  // 自动注入 DatabaseConfig
    builder.Services.Add(NewWebServer)    // 自动注入 ServerConfig
    
    app := builder.Build()
    app.Run()
}
```

**方式 B：手动加载和验证（传统方式）**

```go
// config/config.go
package config

type Config struct {
    Server   ServerConfig   `json:"server"`
    Database DatabaseConfig `json:"database"`
    Logging  LoggingConfig  `json:"logging"`
}

func Load(configuration configuration.IConfiguration) (*Config, error) {
    var cfg Config
    if err := configuration.Bind("", &cfg); err != nil {
        return nil, err
    }
    
    if err := validate(&cfg); err != nil {
        return nil, err
    }
    
    return &cfg, nil
}

func validate(cfg *Config) error {
    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        return fmt.Errorf("invalid port: %d", cfg.Server.Port)
    }
    // ... 更多验证
    return nil
}

// main.go
func main() {
    builder := web.CreateBuilder()
    
    // 手动加载配置
    appConfig, err := config.Load(builder.Configuration)
    if err != nil {
        log.Fatal(err)
    }
    
    // 注册配置为单例
    builder.Services.AddInstance(appConfig)
    
    app := builder.Build()
    app.Run()
}
```

### 5. 使用环境隔离

```go
// 自动加载环境特定配置
env := builder.Environment.GetEnvironmentName()

builder.Configuration.
    AddJsonFile("appsettings.json", false, true).
    AddJsonFile(fmt.Sprintf("appsettings.%s.json", env), true, true).
    AddEnvironmentVariables("").
    AddCommandLine(os.Args[1:])
```

### 6. 配置注释和文档

在配置文件中添加注释说明：

```json
{
  "server": {
    "port": 8080,           // HTTP 监听端口
    "host": "localhost",     // 监听地址，生产环境使用 0.0.0.0
    "timeout": 30            // 请求超时时间（秒）
  }
}
```

或创建配置模板文件：

```
appsettings.json         # 实际配置（不提交到Git）
appsettings.example.json # 配置模板（提交到Git）
```

## API 参考

### IConfiguration

```go
// 读取配置值
Get(key string) string
GetString(key string, defaultValue string) string
GetInt(key string, defaultValue int) int
GetInt64(key string, defaultValue int64) int64
GetBool(key string, defaultValue bool) bool
GetFloat64(key string, defaultValue float64) float64

// 配置节点
GetSection(key string) IConfigurationSection
GetRequiredSection(key string) IConfigurationSection
GetChildren() []IConfigurationSection

// 配置绑定
Bind(section string, target interface{}) error
BindWithOptions(section string, target interface{}, options *BinderOptions) error

// 其他
Exists(key string) bool
Set(key string, value string)
OnChange(callback func())
```

### IConfigurationBuilder

```go
// 添加配置源
AddJsonFile(path string, optional bool, reloadOnChange bool) IConfigurationBuilder
AddYamlFile(path string, optional bool, reloadOnChange bool) IConfigurationBuilder
AddIniFile(path string, optional bool, reloadOnChange bool) IConfigurationBuilder
AddXmlFile(path string, optional bool, reloadOnChange bool) IConfigurationBuilder
AddEnvironmentVariables(prefix string) IConfigurationBuilder
AddCommandLine(args []string) IConfigurationBuilder
AddInMemoryCollection(data map[string]string) IConfigurationBuilder
AddKeyPerFile(directoryPath string, optional bool) IConfigurationBuilder

// 其他
SetBasePath(basePath string) IConfigurationBuilder
Build() IConfiguration
```

### IConfigurationRoot

```go
// 重新加载配置
Reload()

// 调试
GetDebugView() string

// 获取提供者
Providers() []IConfigurationProvider
```

## 常见问题

### 配置文件找不到？

确保配置文件在正确的位置，或使用 `SetBasePath` 设置基础路径：

```go
builder.Configuration.
    SetBasePath("./config").
    AddJsonFile("appsettings.json", false, true)
```

### 配置值没有更新？

1. 检查配置源的优先级，后添加的会覆盖先添加的
2. 确保启用了 `reloadOnChange` 参数
3. 检查文件权限和监控是否正常

### 如何读取复杂的嵌套配置？

使用配置绑定到结构体，而不是逐个读取：

```go
// ✅ 推荐
var config MyConfig
builder.Configuration.Bind("mySection", &config)

// ❌ 不推荐
val1 := config.Get("mySection:level1:level2:value1")
val2 := config.Get("mySection:level1:level2:value2")
// ...
```

### 环境变量的层级分隔符？

环境变量使用双下划线（`__`）表示层级：

```bash
export APP__SERVER__PORT=9000
# 对应配置键: server:port
```

---

[← 返回主目录](../README.md)

