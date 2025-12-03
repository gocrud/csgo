# 配置 API 参考

本文档提供 CSGO 配置系统的完整 API 参考。

## 目录

- [IConfiguration](#iconfiguration)
- [IConfigurationBuilder](#iconfigurationbuilder)
- [配置源](#配置源)
- [Options 模式](#options-模式)
- [绑定函数](#绑定函数)

---

## IConfiguration

配置访问接口。

```go
import "github.com/gocrud/csgo/configuration"
```

### Get

获取配置值。

```go
func (c *Configuration) Get(key string) string
```

**参数：**
- `key` - 配置键，使用 `:` 分隔层级

**返回值：**
- `string` - 配置值，不存在返回空字符串

**示例：**

```go
// 配置: {"database": {"host": "localhost"}}
host := config.Get("database:host")  // "localhost"
```

---

### GetSection

获取配置节。

```go
func (c *Configuration) GetSection(key string) IConfigurationSection
```

**参数：**
- `key` - 配置节键

**返回值：**
- `IConfigurationSection` - 配置节

**示例：**

```go
dbSection := config.GetSection("database")
host := dbSection.Get("host")
port := dbSection.Get("port")
```

---

### GetChildren

获取所有子节。

```go
func (c *Configuration) GetChildren() []IConfigurationSection
```

---

### Bind

将配置节绑定到结构体。

```go
func (c *Configuration) Bind(section string, target interface{}) error
```

**参数：**
- `section` - 配置节键
- `target` - 目标结构体指针

**示例：**

```go
type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Username string `json:"username"`
}

var dbConfig DatabaseConfig
err := config.Bind("database", &dbConfig)
```

---

### Set

设置配置值（运行时）。

```go
func (c *Configuration) Set(key, value string)
```

---

## IConfigurationBuilder

配置构建器接口。

### AddJsonFile

添加 JSON 配置文件。

```go
func (b *ConfigurationBuilder) AddJsonFile(path string, optional bool) *ConfigurationBuilder
```

**参数：**
- `path` - 文件路径
- `optional` - 是否可选（如果为 false，文件不存在会 panic）

**示例：**

```go
builder.Configuration.AddJsonFile("appsettings.json", false)
builder.Configuration.AddJsonFile("appsettings.Development.json", true)
```

---

### AddYamlFile

添加 YAML 配置文件。

```go
func (b *ConfigurationBuilder) AddYamlFile(path string, optional bool) *ConfigurationBuilder
```

---

### AddEnvironmentVariables

添加环境变量。

```go
func (b *ConfigurationBuilder) AddEnvironmentVariables(prefix string) *ConfigurationBuilder
```

**参数：**
- `prefix` - 环境变量前缀（空字符串表示所有环境变量）

**示例：**

```go
// 加载 MYAPP_ 前缀的环境变量
// MYAPP_DATABASE__HOST -> database:host
builder.Configuration.AddEnvironmentVariables("MYAPP_")
```

---

### AddCommandLine

添加命令行参数。

```go
func (b *ConfigurationBuilder) AddCommandLine(args []string) *ConfigurationBuilder
```

**支持格式：**
- `--key=value`
- `--key value`
- `--section.key=value` → `section:key`

**示例：**

```go
builder.Configuration.AddCommandLine(os.Args[1:])
// --database.host=localhost → database:host = localhost
```

---

### AddInMemory

添加内存配置。

```go
func (b *ConfigurationBuilder) AddInMemory(values map[string]string) *ConfigurationBuilder
```

---

## 配置源

### JsonConfigurationSource

JSON 文件配置源。

```go
type JsonConfigurationSource struct {
    Path     string
    Optional bool
}
```

---

### YamlConfigurationSource

YAML 文件配置源。

```go
type YamlConfigurationSource struct {
    Path     string
    Optional bool
}
```

---

### EnvironmentVariablesConfigurationSource

环境变量配置源。

```go
type EnvironmentVariablesConfigurationSource struct {
    Prefix string
}
```

**键名转换规则：**
- `__` → `:` （双下划线转为层级分隔符）
- `_` → `:` （单下划线也转为层级分隔符）
- 前缀会被移除

---

### CommandLineConfigurationSource

命令行配置源。

```go
type CommandLineConfigurationSource struct {
    Args []string
}
```

**键名转换规则：**
- `.` → `:` （点转为层级分隔符）
- `--` 前缀会被移除

---

### InMemoryConfigurationSource

内存配置源。

```go
type InMemoryConfigurationSource struct {
    Data map[string]string
}
```

---

## Options 模式

### Configure

注册配置选项。

```go
func Configure[T any](services di.IServiceCollection, config IConfiguration, section string)
```

**类型参数：**
- `T` - 选项类型

**示例：**

```go
type DatabaseOptions struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Database string `json:"database"`
}

configuration.Configure[DatabaseOptions](builder.Services, builder.Configuration, "database")
```

---

### ConfigureWithDefaults

注册带默认值的配置选项。

```go
func ConfigureWithDefaults[T any](services di.IServiceCollection, config IConfiguration, section string, defaults func() *T)
```

**示例：**

```go
configuration.ConfigureWithDefaults[DatabaseOptions](
    builder.Services,
    builder.Configuration,
    "database",
    func() *DatabaseOptions {
        return &DatabaseOptions{
            Host: "localhost",
            Port: 5432,
        }
    },
)
```

---

### ConfigureWithValidation

注册带验证的配置选项。

```go
func ConfigureWithValidation[T any](services di.IServiceCollection, config IConfiguration, section string, validator func(*T) error) error
```

**示例：**

```go
err := configuration.ConfigureWithValidation[DatabaseOptions](
    builder.Services,
    builder.Configuration,
    "database",
    func(opts *DatabaseOptions) error {
        if opts.Host == "" {
            return errors.New("database host is required")
        }
        return nil
    },
)
```

---

### IOptions[T]

选项访问接口。

```go
type IOptions[T any] interface {
    Value() *T
}
```

**使用方式：**

```go
// 注入 IOptions
func NewUserService(opts configuration.IOptions[DatabaseOptions]) *UserService {
    dbOpts := opts.Value()
    // 使用 dbOpts.Host, dbOpts.Port 等
}
```

---

### IOptionsMonitor[T]

动态配置监视接口（支持热更新）。

```go
type IOptionsMonitor[T any] interface {
    CurrentValue() *T
    OnChange(callback func(*T))
}
```

**使用方式：**

```go
func NewUserService(monitor configuration.IOptionsMonitor[DatabaseOptions]) *UserService {
    // 获取当前值
    current := monitor.CurrentValue()
    
    // 监听变化
    monitor.OnChange(func(newOpts *DatabaseOptions) {
        log.Printf("配置已更新: %+v", newOpts)
    })
}
```

---

## 绑定函数

### BindOptions

绑定配置到选项对象。

```go
func BindOptions[T any](config IConfiguration, section string) (*T, error)
```

**示例：**

```go
dbOpts, err := configuration.BindOptions[DatabaseOptions](config, "database")
if err != nil {
    log.Fatal(err)
}
```

---

### MustBindOptions

绑定配置，失败时 panic。

```go
func MustBindOptions[T any](config IConfiguration, section string) *T
```

**示例：**

```go
dbOpts := configuration.MustBindOptions[DatabaseOptions](config, "database")
```

---

## 完整示例

### appsettings.json

```json
{
  "server": {
    "host": "0.0.0.0",
    "port": 8080
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "database": "myapp",
    "username": "admin"
  },
  "logging": {
    "level": "info"
  }
}
```

### main.go

```go
package main

import (
    "github.com/gocrud/csgo/configuration"
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

type ServerOptions struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

type DatabaseOptions struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Database string `json:"database"`
    Username string `json:"username"`
}

func main() {
    builder := web.CreateBuilder()
    
    // 配置已自动加载 appsettings.json
    
    // 注册选项
    configuration.Configure[ServerOptions](builder.Services, builder.Configuration, "server")
    configuration.Configure[DatabaseOptions](builder.Services, builder.Configuration, "database")
    
    // 注册使用选项的服务
    builder.Services.AddSingleton(func(opts configuration.IOptions[DatabaseOptions]) *DatabaseConnection {
        dbOpts := opts.Value()
        return NewDatabaseConnection(dbOpts.Host, dbOpts.Port, dbOpts.Database)
    })
    
    app := builder.Build()
    
    app.MapGet("/config", func(c *web.HttpContext) web.IActionResult {
        serverOpts := di.GetRequiredService[configuration.IOptions[ServerOptions]](app.Services)
        return c.Ok(serverOpts.Value())
    })
    
    app.Run()
}
```

---

## 相关文档

- [配置管理指南](../guides/configuration.md)
- [依赖注入 API](di.md)
- [Web 框架 API](web.md)

