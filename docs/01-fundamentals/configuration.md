# 配置管理

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

本章讲解 CSGO 的配置管理系统。

## 完整文档

关于配置管理的完整详细文档，请查看：

👉 **[配置管理完整文档](../../configuration/README.md)**

## 核心内容概览

### 1. 配置源
- JSON 文件
- YAML 文件
- 环境变量
- 命令行参数

### 2. 配置绑定
- 绑定到结构体
- 嵌套配置
- 数组配置

### 3. Configure[T] 模式（推荐）
- 类型安全的配置
- 自动依赖注入
- 热更新支持
- 一行代码注册

### 4. Options 模式
- 定义 Options
- 使用 Options
- 默认值

### 5. 配置热更新
- 启用文件监控
- 监听配置变化
- 实时更新

## 快速示例

### 推荐方式：Configure[T]

```go
import "github.com/gocrud/csgo/configuration"

// 定义配置结构
type ServerConfig struct {
    Port int    `json:"port"`
    Host string `json:"host"`
}

type DatabaseConfig struct {
    Connection string `json:"connection"`
}

func main() {
    builder := web.CreateBuilder()
    
    // 注册配置（推荐方式）
    configuration.Configure[ServerConfig](builder.Services, "server")
    configuration.Configure[DatabaseConfig](builder.Services, "database")
    
    // 注册服务（自动注入配置）
    builder.Services.Add(NewMyService)
    
    app := builder.Build()
    app.Run()
}

// 在服务中使用配置
type MyService struct {
    serverConfig configuration.IOptionsMonitor[ServerConfig]
}

func NewMyService(serverConfig configuration.IOptionsMonitor[ServerConfig]) *MyService {
    return &MyService{serverConfig: serverConfig}
}

func (s *MyService) DoWork() {
    // 自动获取最新配置（支持热更新）
    config := s.serverConfig.CurrentValue()
    fmt.Printf("Server: %s:%d\n", config.Host, config.Port)
}
```

### 传统方式：直接读取

```go
builder := web.CreateBuilder()

// 直接读取配置
port := builder.Configuration.GetInt("server:port", 8080)
dbConn := builder.Configuration.Get("database:connection")

// 手动绑定到结构体
type ServerConfig struct {
    Port int    `json:"port"`
    Host string `json:"host"`
}

var config ServerConfig
builder.Configuration.Bind("server", &config)

app := builder.Build()
app.Run()
```

## 配置文件示例

**appsettings.json：**

```json
{
  "server": {
    "port": 8080,
    "host": "localhost"
  },
  "database": {
    "connection": "postgres://localhost/mydb"
  }
}
```

## 下一步

继续学习：[HTTP 上下文](http-context.md) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

