# Logging Package

日志包提供与 .NET `Microsoft.Extensions.Logging` 完全一致的日志抽象，底层使用 zerolog 实现。

## 特性

- ✅ 与 .NET ILogger 接口完全一致
- ✅ 支持依赖注入
- ✅ 支持从 appsettings.json 配置
- ✅ 支持控制台和文件输出
- ✅ 可替换日志引擎（当前使用 zerolog）
- ✅ 支持泛型 ILogger<T>

## 快速开始

### 1. 基础使用（自动配置）

```go
package main

import (
    "github.com/gocrud/csgo/web"
    "github.com/gocrud/csgo/logging"
)

func main() {
    // 创建应用（日志已自动注册）
    builder := web.CreateBuilder()
    
    // 注册服务
    builder.Services.AddSingleton(NewUserService)
    
    app := builder.Build()
    app.Run()
}

// 在服务中使用日志
type UserService struct {
    logger logging.ILogger
}

func NewUserService(factory logging.ILoggerFactory) *UserService {
    return &UserService{
        logger: factory.CreateLogger("UserService"),
    }
}

func (s *UserService) GetUser(id int) {
    s.logger.LogInformation("Getting user with id: %d", id)
    // ...
}
```

### 2. 使用泛型 ILogger<T>

```go
func NewUserService(factory logging.ILoggerFactory) *UserService {
    return &UserService{
        logger: logging.GetLogger[UserService](factory),
    }
}
```

### 3. 配置日志（appsettings.json）

```json
{
  "Logging": {
    "LogLevel": {
      "Default": "Information"
    },
    "Console": {
      "Enabled": true
    },
    "File": {
      "Enabled": true,
      "Path": "logs/app.log"
    }
  }
}
```

### 4. 程序化配置

```go
builder := web.CreateBuilder()

// 自定义日志配置
logging.AddZerolog(builder.Services, func(opts *logging.LoggingOptions) {
    opts.MinLevel = logging.LogLevelDebug
    opts.Console.Enabled = true
    opts.Console.UseConsoleWriter = true  // 开发环境使用可读格式
    opts.File.Enabled = true
    opts.File.Path = "logs/myapp.log"
})
```

## 日志级别

```go
logging.LogLevelTrace       // 最详细的消息
logging.LogLevelDebug       // 调试信息
logging.LogLevelInformation // 常规信息
logging.LogLevelWarning     // 警告
logging.LogLevelError       // 错误
logging.LogLevelCritical    // 严重错误
logging.LogLevelNone        // 禁用日志
```

## API 参考

### ILogger 接口

```go
type ILogger interface {
    Log(level LogLevel, message string, args ...interface{})
    LogTrace(message string, args ...interface{})
    LogDebug(message string, args ...interface{})
    LogInformation(message string, args ...interface{})
    LogWarning(message string, args ...interface{})
    LogError(err error, message string, args ...interface{})
    LogCritical(err error, message string, args ...interface{})
    IsEnabled(level LogLevel) bool
}
```

### 使用示例

```go
// 信息日志
logger.LogInformation("User %s logged in", username)

// 警告日志
logger.LogWarning("Cache miss for key: %s", key)

// 错误日志
logger.LogError(err, "Failed to save user %d", userID)

// 检查日志级别
if logger.IsEnabled(logging.LogLevelDebug) {
    logger.LogDebug("Expensive debug info: %v", expensiveOperation())
}
```

## 与 .NET 的对比

| .NET | CSGO | 说明 |
|------|------|------|
| `ILogger` | `logging.ILogger` | 完全一致 |
| `ILogger<T>` | `logging.GetLogger[T]()` | 泛型支持 |
| `ILoggerFactory` | `logging.ILoggerFactory` | 完全一致 |
| `LogLevel` | `logging.LogLevel` | 完全一致 |
| `services.AddLogging()` | 自动注册 | 自动集成 |

## 更换日志引擎

当前使用 zerolog 作为底层实现。如需更换：

1. 实现 `ILogger` 和 `ILoggerFactory` 接口
2. 在 `AddLogging()` 中注册新的实现

```go
// 示例：使用自定义日志引擎
services.AddSingleton(func() logging.ILoggerFactory {
    return NewCustomLoggerFactory()
})
```

## 最佳实践

1. **使用依赖注入**：通过构造函数注入 `ILoggerFactory`
2. **使用泛型**：使用 `GetLogger[T]()` 自动设置类别名
3. **检查级别**：对于昂贵的日志操作，先检查 `IsEnabled()`
4. **结构化日志**：使用格式化字符串而不是字符串拼接
5. **错误日志**：使用 `LogError()` 并传入 error 对象

## 配置示例

### 开发环境（appsettings.Development.json）

```json
{
  "Logging": {
    "LogLevel": {
      "Default": "Debug"
    },
    "Console": {
      "Enabled": true
    }
  }
}
```

### 生产环境（appsettings.Production.json）

```json
{
  "Logging": {
    "LogLevel": {
      "Default": "Warning"
    },
    "Console": {
      "Enabled": true
    },
    "File": {
      "Enabled": true,
      "Path": "/var/log/myapp/app.log"
    }
  }
}
```

