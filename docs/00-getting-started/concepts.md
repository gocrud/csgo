# 核心概念

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

本章将介绍 CSGO 框架的核心概念和设计思想。

## WebApplicationBuilder

### 作用

`WebApplicationBuilder` 是应用程序的构建器，负责配置和初始化应用：

```go
builder := web.CreateBuilder()
```

它会自动：
- 加载配置文件（appsettings.json）
- 设置环境（Development/Production）
- 初始化依赖注入容器
- 配置日志系统

### 主要属性

```go
builder := web.CreateBuilder()

// 服务容器
builder.Services         // 用于注册服务

// 配置管理
builder.Configuration    // 读取配置

// 环境信息
builder.Environment      // 环境判断

// Host 配置
builder.Host            // 配置主机

// WebHost 配置
builder.WebHost         // 配置 Web 服务器
```

### 构建应用

```go
builder := web.CreateBuilder()

// 配置服务和应用
// ...

// 构建应用
app := builder.Build()
```

## 依赖注入（DI）

### 什么是依赖注入？

依赖注入是一种设计模式，让对象的依赖关系由外部管理，而不是对象自己创建。

### 为什么需要 DI？

- ✅ **松耦合**：服务之间解耦
- ✅ **可测试**：方便单元测试
- ✅ **可维护**：集中管理依赖关系
- ✅ **可重用**：服务可以在多处使用

### 使用 DI

```go
// 1. 定义服务
type UserService struct {
    db *Database
}

func NewUserService(db *Database) *UserService {
    return &UserService{db: db}
}

// 2. 注册服务
builder := web.CreateBuilder()
builder.Services.Add(NewDatabase)
builder.Services.Add(NewUserService)  // 自动注入 Database

// 3. 使用服务
app := builder.Build()
app.MapGet("/users", func(c *web.HttpContext) web.IActionResult {
    // 从容器获取服务
    userService := di.Get[*UserService](c.Services)
    users := userService.GetAll()
    return c.Ok(users)
})
```

### Singleton 生命周期

CSGO 当前支持 Singleton（单例）生命周期：

```go
// 注册单例服务
builder.Services.Add(NewUserService)

// 整个应用共享同一个实例
```

特点：
- 全局唯一实例
- 首次使用时创建
- 应用关闭时销毁

## 路由系统

### 定义路由

CSGO 提供了简洁的路由定义方式：

```go
app := builder.Build()

// HTTP 方法
app.MapGet("/users", getUsers)       // GET
app.MapPost("/users", createUser)    // POST
app.MapPut("/users/:id", updateUser) // PUT
app.MapDelete("/users/:id", deleteUser) // DELETE
app.MapPatch("/users/:id", patchUser) // PATCH
```

### 路径参数

```go
// 定义路径参数
app.MapGet("/users/:id", func(c *web.HttpContext) web.IActionResult {
    id := c.RawCtx().Param("id")
    return c.Ok(web.M{"id": id})
})

// 多个参数
app.MapGet("/users/:userId/posts/:postId", func(c *web.HttpContext) web.IActionResult {
    userId := c.RawCtx().Param("userId")
    postId := c.RawCtx().Param("postId")
    return c.Ok(web.M{"userId": userId, "postId": postId})
})
```

### 查询参数

```go
app.MapGet("/search", func(c *web.HttpContext) web.IActionResult {
    keyword := c.RawCtx().Query("keyword")
    page := c.RawCtx().DefaultQuery("page", "1")
    return c.Ok(web.M{"keyword": keyword, "page": page})
})
// GET /search?keyword=golang&page=2
```

### 路由组

```go
api := app.MapGroup("/api")
{
    v1 := api.MapGroup("/v1")
    {
        users := v1.MapGroup("/users")
        {
            users.MapGet("", listUsers)        // GET /api/v1/users
            users.MapGet("/:id", getUser)      // GET /api/v1/users/:id
            users.MapPost("", createUser)      // POST /api/v1/users
        }
    }
}
```

## HttpContext 和 ActionResult

### HttpContext

`HttpContext` 封装了 HTTP 请求和响应：

```go
func handler(c *web.HttpContext) web.IActionResult {
    // 获取原始 gin.Context
    ginCtx := c.RawCtx()
    
    // 获取请求 Context
    ctx := c.Context()
    
    // 访问服务容器
    service := di.Get[*Service](c.Services)
    
    // 绑定请求体
    var req Request
    c.MustBindJSON(&req)
    
    return c.Ok(data)
}
```

### ActionResult

`ActionResult` 表示操作的结果，提供统一的响应格式：

```go
// 成功响应
return c.Ok(data)           // 200
return c.Created(data)      // 201
return c.NoContent()        // 204

// 错误响应
return c.BadRequest("...")   // 400
return c.Unauthorized("...") // 401
return c.NotFound("...")     // 404
return c.InternalError("...") // 500
```

**响应格式：**

```json
{
  "success": true,
  "data": { /* 数据 */ }
}
```

或错误响应：

```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "资源不存在"
  }
}
```

## 中间件管道

### 什么是中间件？

中间件是处理请求的组件链，每个中间件可以：
- 在请求到达处理器前执行代码
- 调用下一个中间件或处理器
- 在响应返回前执行代码

### 使用中间件

```go
app := builder.Build()

// 全局中间件
app.Use(loggingMiddleware)
app.Use(authMiddleware)

// 定义路由
app.MapGet("/api/users", getUsers)
```

### 自定义中间件

```go
func loggingMiddleware(c *gin.Context) {
    start := time.Now()
    
    // 处理请求
    c.Next()
    
    // 请求完成
    latency := time.Since(start)
    fmt.Printf("[%s] %s %v\n", 
        c.Request.Method,
        c.Request.URL.Path,
        latency,
    )
}

app.Use(loggingMiddleware)
```

### 中间件顺序

中间件按注册顺序执行：

```go
app.Use(middleware1)  // 1. 执行
app.Use(middleware2)  // 2. 执行
app.Use(middleware3)  // 3. 执行

app.MapGet("/", handler)  // 4. 执行

// 响应返回时逆序执行
// 3 -> 2 -> 1
```

## 配置系统

### 读取配置

```go
builder := web.CreateBuilder()

// 读取配置值
port := builder.Configuration.GetInt("server:port", 8080)
dbConn := builder.Configuration.Get("database:connection")

// 绑定到结构体
var config AppConfig
builder.Configuration.Bind("", &config)
```

### 配置文件

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

### 环境特定配置

**appsettings.Development.json：**

```json
{
  "database": {
    "connection": "postgres://localhost/mydb_dev"
  }
}
```

框架自动根据环境加载对应配置并覆盖基础配置。

## 应用生命周期

### 启动流程

```
1. CreateBuilder()     // 创建构建器
2. 配置服务和应用
3. Build()             // 构建应用
4. 定义路由
5. Run()               // 启动应用
```

### 关闭流程

```
1. 接收关闭信号（Ctrl+C）
2. 停止接受新请求
3. 等待现有请求完成
4. 停止后台服务
5. 清理资源
6. 退出应用
```

## 设计原则

### 1. 约定优于配置

框架提供合理的默认值，减少配置工作：

```go
builder := web.CreateBuilder()  // 自动加载配置
app := builder.Build()
app.Run()  // 默认监听 :8080
```

### 2. 类型安全

使用 Go 泛型提供类型安全的 API：

```go
// 类型安全的服务解析
userService := di.Get[*UserService](provider)

// 类型安全的请求验证
req, err := web.BindAndValidate[CreateUserRequest](c)
```

### 3. 清晰的职责分离

```
Controller  -> 处理 HTTP 请求
Service     -> 业务逻辑
Repository  -> 数据访问
```

### 4. 可测试性

依赖注入使代码容易测试：

```go
// 测试时注入 mock 服务
services.Add(func() *Database {
    return NewMockDatabase()
})
```

## 小结

本章介绍了 CSGO 框架的核心概念：

- ✅ **WebApplicationBuilder** - 应用构建器
- ✅ **依赖注入** - 服务管理
- ✅ **路由系统** - 请求映射
- ✅ **HttpContext & ActionResult** - 请求处理
- ✅ **中间件管道** - 请求拦截
- ✅ **配置系统** - 配置管理

## 下一步

恭喜完成快速入门！🎉

现在你已经了解了 CSGO 的核心概念，可以开始深入学习更多特性：

- [阶段 1：核心基础](../01-fundamentals/) - 深入学习 DI、路由、配置
- [阶段 2：构建 API](../02-building-apis/) - 学习构建生产级 API
- [阶段 3：高级特性](../03-advanced-features/) - 中间件、日志、测试

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

