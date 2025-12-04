# 配置管理指南

## 📦 框架配置系统

本示例使用 csgo 框架的配置管理系统，支持多种配置源和优先级覆盖。

## 🎯 配置加载顺序

配置按以下顺序加载（后加载的会覆盖先加载的）：

```
1. config.dev.json      (基础配置文件)
2. 环境变量            (APP_ 前缀)
3. 命令行参数          (--key value)
```

## 📂 使用示例

### 1. 在 Bootstrap 中加载配置

```go
import "github.com/gocrud/csgo/configuration"

func Bootstrap() *web.WebApplication {
    builder := web.CreateBuilder()
    
    // 构建配置
    config := configuration.NewConfigurationBuilder().
        AddJsonFile("configs/config.dev.json", true, false).
        AddEnvironmentVariables("APP_").
        Build()
    
    // 绑定到结构体
    var appConfig configs.Config
    config.Bind("", &appConfig)
    
    // 注册到 DI 容器
    builder.Services.AddSingleton(func() configuration.IConfiguration {
        return config
    })
    builder.Services.AddSingleton(func() *configs.Config {
        return &appConfig
    })
    
    // ...
}
```

### 2. 读取配置值

#### 方式 A：通过 IConfiguration 接口

```go
type MyHandler struct {
    config configuration.IConfiguration
}

func (h *MyHandler) Handle(c *web.HttpContext) {
    // 读取单个值
    port := h.config.Get("server:admin_port")
    
    // 读取节点
    serverSection := h.config.GetSection("server")
    adminPort := serverSection.Get("admin_port")
}
```

#### 方式 B：通过绑定的结构体（推荐）

```go
type MyHandler struct {
    config *configs.Config
}

func (h *MyHandler) Handle(c *web.HttpContext) {
    // 直接访问结构体字段
    port := h.config.Server.AdminPort
    dbHost := h.config.Database.Host
}
```

### 3. 配置文件格式

**JSON 格式** (`config.dev.json`):

```json
{
  "server": {
    "admin_port": ":8081",
    "api_port": ":8080"
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "postgres",
    "password": "password"
  }
}
```

**转换为扁平化键值：**

```
server:admin_port = :8081
server:api_port = :8080
database:host = localhost
database:port = 5432
```

### 4. 环境变量覆盖

使用 `APP_` 前缀，支持两种格式：

**格式 1：双下划线（推荐）**

```bash
export APP_Database__Host=192.168.1.100
export APP_Database__Port=3306
```

**格式 2：单下划线**

```bash
export APP_Server_AdminPort=:9091
```

**优先级：** 环境变量 > JSON 文件

### 5. 命令行参数

```bash
# 格式 1: --key=value
./bin/admin --server:admin_port=:9091 --database:host=192.168.1.100

# 格式 2: --key value
./bin/admin --server:admin_port :9091

# 也支持点号
./bin/admin --server.admin_port=:9091
```

**优先级：** 命令行 > 环境变量 > JSON 文件

## 🔧 配置结构定义

### 基本结构

```go
// configs/config.go
package configs

type Config struct {
    Server   ServerConfig   `json:"server"`
    Database DatabaseConfig `json:"database"`
    Cache    CacheConfig    `json:"cache"`
}

type ServerConfig struct {
    AdminPort string `json:"admin_port"`
    ApiPort   string `json:"api_port"`
}

type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    User     string `json:"user"`
    Password string `json:"password"`
    Database string `json:"database"`
}
```

**关键点：**
- 使用 `json` tag 映射配置键名
- 框架会自动将 `server:admin_port` 映射到 `Server.AdminPort`

### 嵌套配置

```go
type Config struct {
    Redis RedisConfig `json:"redis"`
}

type RedisConfig struct {
    Master  RedisNode   `json:"master"`
    Slaves  []RedisNode `json:"slaves"`
    Cluster bool        `json:"cluster"`
}

type RedisNode struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}
```

**JSON 配置：**

```json
{
  "redis": {
    "master": {
      "host": "localhost",
      "port": 6379
    },
    "slaves": [
      {"host": "slave1", "port": 6379},
      {"host": "slave2", "port": 6379}
    ],
    "cluster": false
  }
}
```

**扁平化格式：**

```
redis:master:host = localhost
redis:master:port = 6379
redis:slaves:0:host = slave1
redis:slaves:0:port = 6379
redis:slaves:1:host = slave2
redis:slaves:1:port = 6379
redis:cluster = false
```

## 🎨 高级用法

### 1. 配置热更新

```go
config := configuration.NewConfigurationBuilder().
    AddJsonFile("config.json", true, true).  // reloadOnChange = true
    Build()

// 注册变更回调
config.OnChange(func() {
    fmt.Println("配置已更新！")
    // 重新绑定配置
    var newConfig configs.Config
    config.Bind("", &newConfig)
})
```

### 2. 多环境配置

```go
env := os.Getenv("APP_ENV")
if env == "" {
    env = "dev"
}

config := configuration.NewConfigurationBuilder().
    AddJsonFile("configs/config.json", true, false).           // 基础配置
    AddJsonFile(fmt.Sprintf("configs/config.%s.json", env), true, false).  // 环境特定配置
    AddEnvironmentVariables("APP_").
    Build()
```

**配置文件：**
- `config.json` - 通用配置
- `config.dev.json` - 开发环境
- `config.prod.json` - 生产环境

### 3. 配置验证

```go
var appConfig configs.Config
if err := config.Bind("", &appConfig); err != nil {
    panic(fmt.Sprintf("配置绑定失败: %v", err))
}

// 验证配置
if appConfig.Database.Host == "" {
    panic("数据库主机地址未配置")
}
if appConfig.Database.Port == 0 {
    panic("数据库端口未配置")
}
```

### 4. 配置分组

```go
// 只绑定某个节点
var dbConfig configs.DatabaseConfig
config.Bind("database", &dbConfig)

// 只绑定服务器配置
var serverConfig configs.ServerConfig
config.Bind("server", &serverConfig)
```

## 📝 实际示例

### 示例 1：数据库连接

```go
// shared/infrastructure/database/db.go
type DB struct {
    config *configs.DatabaseConfig
}

func NewDB(appConfig *configs.Config) *DB {
    return &DB{
        config: &appConfig.Database,
    }
}

func (db *DB) Connect() error {
    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
        db.config.Host,
        db.config.Port,
        db.config.User,
        db.config.Password,
        db.config.Database,
    )
    // 连接数据库...
}
```

### 示例 2：启动端口

```go
// cmd/admin/main.go
func main() {
    app := admin.Bootstrap()
    
    // 从 DI 容器获取配置
    var config *configs.Config
    app.Services.GetRequiredService(&config)
    
    // 使用配置的端口
    app.Run(config.Server.AdminPort)
}
```

### 示例 3：环境变量覆盖

```bash
# 开发环境使用默认配置
./bin/api

# 生产环境通过环境变量覆盖
export APP_Database__Host=prod-db.example.com
export APP_Database__Password=prod_secret_password
export APP_Cache__Host=prod-redis.example.com
./bin/api

# 或者通过命令行
./bin/api \
  --database:host=prod-db.example.com \
  --database:password=prod_secret_password \
  --cache:host=prod-redis.example.com
```

## 🔐 敏感信息处理

### 方式 1：环境变量

```bash
# 不要在配置文件中硬编码密码
export APP_Database__Password=secret_password
export APP_Redis__Password=redis_secret
```

### 方式 2：配置文件加密

```go
// 自定义配置源，解密敏感信息
type EncryptedConfigSource struct {
    innerSource configuration.IConfigurationSource
    decryptor   func(string) string
}

func (s *EncryptedConfigSource) Load() map[string]string {
    data := s.innerSource.Load()
    
    // 解密特定字段
    if encrypted, ok := data["database:password"]; ok {
        data["database:password"] = s.decryptor(encrypted)
    }
    
    return data
}
```

## ✅ 最佳实践

1. **结构化配置**：使用结构体而不是字符串键
2. **环境分离**：开发、测试、生产使用不同配置文件
3. **敏感信息**：通过环境变量传递，不要提交到代码库
4. **配置验证**：启动时验证必需的配置项
5. **默认值**：为可选配置提供合理的默认值
6. **文档化**：注释说明每个配置项的作用

## 📚 相关文档

- [csgo 配置管理文档](../../../docs/guides/configuration.md)
- [配置结构定义](config.go)

---

**使用框架配置系统，让你的应用更灵活！** 🚀

