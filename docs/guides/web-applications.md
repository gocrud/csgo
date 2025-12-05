# Web 应用指南

本指南介绍如何使用 CSGO 框架构建 Web 应用，包括应用构建、路由系统、HttpContext、ActionResult、中间件和控制器。

## 目录

- [快速开始](#快速开始)
- [WebApplicationBuilder](#webapplicationbuilder)
- [WebApplication](#webapplication)
- [路由系统](#路由系统)
- [HttpContext](#httpcontext)
- [ActionResult](#actionresult)
- [中间件](#中间件)
- [控制器](#控制器)
- [静态文件](#静态文件)
- [CORS](#cors)
- [完整示例](#完整示例)
- [最佳实践](#最佳实践)

## 快速开始

### 最简单的 Web 应用

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gocrud/csgo/web"
)

func main() {
    // 创建应用构建器
    builder := web.CreateBuilder()
    
    // 构建应用
    app := builder.Build()
    
    // 定义路由（支持多种处理器类型）
    app.MapGet("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Hello, CSGO!"})
    })
    
    // 运行应用
    app.Run()
}
```

### 使用 HttpContext 和 ActionResult（推荐）

```go
package main

import (
    "github.com/gocrud/csgo/web"
)

func main() {
    builder := web.CreateBuilder()
    app := builder.Build()
    
    // 使用 ActionResult 模式 - 代码更清晰
    app.MapGet("/", func(c *web.HttpContext) web.IActionResult {
        return c.Ok(gin.H{"message": "Hello, CSGO!"})
    })
    
    app.MapGet("/users/:id", func(c *web.HttpContext) web.IActionResult {
        id, err := c.MustPathInt("id")
        if err != nil {
            return err  // 自动返回 400 Bad Request
        }
        
        user := getUserByID(id)
        if user == nil {
            return c.NotFound("用户不存在")
        }
        
        return c.Ok(user)
    })
    
    app.Run()
}
```

### 带服务的 Web 应用

```go
package main

import (
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

type GreetingService struct{}

func NewGreetingService() *GreetingService {
    return &GreetingService{}
}

func (s *GreetingService) Greet(name string) string {
    return "Hello, " + name + "!"
}

func main() {
    builder := web.CreateBuilder()
    
    // 注册服务
    builder.Services.AddSingleton(NewGreetingService)
    
    app := builder.Build()
    
    // 使用 ActionResult 模式
    app.MapGet("/greet/:name", func(c *web.HttpContext) web.IActionResult {
        service := di.GetRequiredService[*GreetingService](app.Services)
        greeting := service.Greet(c.Param("name"))
        return c.Ok(gin.H{"greeting": greeting})
    })
    
    app.Run()
}
```

## WebApplicationBuilder

`WebApplicationBuilder` 是 Web 应用的构建器，对应 .NET 的 `WebApplicationBuilder`。

### 创建构建器

```go
// 基本创建
builder := web.CreateBuilder()

// 带命令行参数
builder := web.CreateBuilder(os.Args...)
```

### 属性

| 属性 | 类型 | 说明 |
|------|------|------|
| `Services` | `di.IServiceCollection` | 服务集合，用于注册依赖 |
| `Configuration` | `configuration.IConfiguration` | 应用配置 |
| `Environment` | `hosting.IHostEnvironment` | 环境信息 |
| `Host` | `*ConfigureHostBuilder` | 主机配置 |
| `WebHost` | `*ConfigureWebHostBuilder` | Web 主机配置 |

### 注册服务

```go
builder := web.CreateBuilder()

// 注册单例服务
builder.Services.AddSingleton(NewDatabaseConnection)

// 注册作用域服务
builder.Services.AddScoped(NewUserService)

// 注册瞬态服务
builder.Services.AddTransient(NewEmailService)

// 注册托管服务
builder.Services.AddHostedService(NewBackgroundWorker)
```

### 配置主机

```go
builder := web.CreateBuilder()

// 配置主机服务
builder.Host.ConfigureServices(func(services di.IServiceCollection) {
    services.AddSingleton(NewMyService)
})

// 配置监听地址
builder.WebHost.UseUrls("http://localhost:5000", "https://localhost:5001")
```

### 访问配置

```go
builder := web.CreateBuilder()

// 读取配置
dbHost := builder.Configuration.Get("Database:Host")

// 检查环境
if builder.Environment.IsDevelopment() {
    // 开发环境特定配置
}
```

## WebApplication

`WebApplication` 表示已配置的 Web 应用。

### 构建应用

```go
builder := web.CreateBuilder()
// ... 配置服务
app := builder.Build()
```

### 属性

| 属性 | 类型 | 说明 |
|------|------|------|
| `Services` | `di.IServiceProvider` | 服务提供者 |

### 运行应用

```go
// 阻塞运行（推荐）
app.Run()

// 异步运行
ctx := context.Background()
app.RunAsync(ctx)

// 手动控制
app.Start(ctx)
// ... 执行其他操作
app.Stop(ctx)
```

## 路由系统

### 处理器类型

CSGO 的路由方法支持**三种处理器类型**，无需使用不同的方法名：

```go
// 1. gin.HandlerFunc - 传统方式
app.MapGet("/old", func(c *gin.Context) {
    c.JSON(200, gin.H{"style": "traditional"})
})

// 2. func(*web.HttpContext) - 使用 HttpContext
app.MapGet("/new", func(c *web.HttpContext) {
    page := c.QueryInt("page", 1)
    c.JSON(200, gin.H{"page": page})
})

// 3. func(*web.HttpContext) web.IActionResult - ActionResult 模式（推荐）
app.MapGet("/best", func(c *web.HttpContext) web.IActionResult {
    return c.Ok(gin.H{"style": "action result"})
})
```

### 基本路由

```go
app := builder.Build()

// GET 请求
app.MapGet("/users", func(c *web.HttpContext) web.IActionResult {
    return c.Ok(users)
})

// POST 请求
app.MapPost("/users", func(c *web.HttpContext) web.IActionResult {
    var user User
    if err := c.MustBindJSON(&user); err != nil {
        return err
    }
    return c.Created(user)
})

// PUT 请求
app.MapPut("/users/:id", func(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    return c.Ok(gin.H{"updated": id})
})

// DELETE 请求
app.MapDelete("/users/:id", func(c *web.HttpContext) web.IActionResult {
    return c.NoContent()
})

// PATCH 请求
app.MapPatch("/users/:id", func(c *web.HttpContext) web.IActionResult {
    return c.Ok(gin.H{"patched": true})
})
```

### 路由参数

```go
// 路径参数
app.MapGet("/users/:id", func(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    return c.Ok(gin.H{"id": id})
})

// 查询参数
app.MapGet("/search", func(c *web.HttpContext) web.IActionResult {
    query := c.Query("q")
    page := c.QueryInt("page", 1)
    size := c.QueryInt("size", 10)
    return c.Ok(gin.H{"query": query, "page": page, "size": size})
})

// 通配符
app.MapGet("/files/*filepath", func(c *web.HttpContext) web.IActionResult {
    path := c.Param("filepath")
    return web.Content(200, "File: "+path)
})
```

### 路由组

```go
// 创建路由组
api := app.MapGroup("/api")

// 在组内定义路由
api.MapGet("/users", GetUsers)
api.MapPost("/users", CreateUser)

// 嵌套路由组
v1 := api.MapGroup("/v1")
v1.MapGet("/products", GetProductsV1)

v2 := api.MapGroup("/v2")
v2.MapGet("/products", GetProductsV2)
```

### 路由元数据（OpenAPI）

```go
app.MapGet("/users", GetUsers).
    WithOpenApi(
        openapi.Name("GetUsers"),
        openapi.Summary("获取所有用户"),
        openapi.Description("返回系统中所有用户的列表"),
        openapi.Tags("Users"),
        openapi.Produces[[]User](200),
    )

app.MapPost("/users", CreateUser).
    WithOpenApi(
        openapi.Name("CreateUser"),
        openapi.Summary("创建用户"),
        openapi.Description("创建一个新用户并返回创建结果"),
        openapi.Tags("Users"),
        openapi.Accepts[CreateUserRequest]("application/json"),
        openapi.Produces[User](201),
    )
```

### 路由组元数据

```go
// 组级别配置
users := app.MapGroup("/api/users").
    WithOpenApi(
        openapi.Tags("Users"),
    )

// 组内路由自动继承标签
users.MapGet("", GetAllUsers).
    WithOpenApi(
        openapi.Name("GetAllUsers"),
        openapi.Summary("获取所有用户"),
        openapi.Produces[[]User](200),
    )

users.MapGet("/:id", GetUserByID).
    WithOpenApi(
        openapi.Name("GetUserByID"),
        openapi.Summary("根据 ID 获取用户"),
        openapi.Produces[User](200),
        openapi.ProducesProblem(404),
    )
```

## HttpContext

`HttpContext` 是对 `gin.Context` 的包装，提供更便捷的 API 和统一的响应格式。

### 创建和使用

```go
// HttpContext 会自动创建，直接在处理器中使用
app.MapGet("/example", func(c *web.HttpContext) web.IActionResult {
    // c 就是 HttpContext，嵌入了 gin.Context
    // 所有 gin.Context 的方法仍然可用
    return c.Ok(data)
})
```

### 参数获取

```go
app.MapGet("/users/:id", func(c *web.HttpContext) web.IActionResult {
    // 路径参数（自动错误处理）
    id, err := c.MustPathInt("id")
    if err != nil {
        return err  // 自动返回 400 Bad Request
    }
    
    // 路径参数（手动错误处理）
    id2, parseErr := c.PathInt("id")
    if parseErr != nil {
        return c.BadRequest("无效的 ID")
    }
    
    // 查询参数（带默认值）
    page := c.QueryInt("page", 1)
    size := c.QueryInt("size", 10)
    active := c.QueryBool("active", true)
    
    // 原始 gin.Context 方法仍然可用
    name := c.Query("name")
    header := c.GetHeader("Authorization")
    
    return c.Ok(gin.H{"id": id, "page": page})
})
```

### 参数方法一览

| 方法 | 说明 | 返回值 |
|------|------|--------|
| `PathInt(key)` | 获取路径参数转 int | `(int, error)` |
| `PathInt64(key)` | 获取路径参数转 int64 | `(int64, error)` |
| `MustPathInt(key)` | 获取路径参数，失败返回错误结果 | `(int, IActionResult)` |
| `QueryInt(key, default)` | 获取查询参数转 int | `int` |
| `QueryInt64(key, default)` | 获取查询参数转 int64 | `int64` |
| `QueryBool(key, default)` | 获取查询参数转 bool | `bool` |

### 请求绑定

```go
type CreateUserRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

app.MapPost("/users", func(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    
    // 方式 1: MustBindJSON - 失败自动返回错误
    if err := c.MustBindJSON(&req); err != nil {
        return err  // 自动返回 400 Bad Request
    }
    
    // 方式 2: BindJSON - 返回 ok 和错误结果
    ok, errResult := c.BindJSON(&req)
    if !ok {
        return errResult
    }
    
    // 方式 3: 查询参数绑定
    var query SearchQuery
    ok, errResult = c.BindQuery(&query)
    if !ok {
        return errResult
    }
    
    return c.Created(user)
})
```

### 响应方法

HttpContext 提供便捷的响应方法，返回 `IActionResult`：

```go
func handler(c *web.HttpContext) web.IActionResult {
    // 成功响应
    return c.Ok(data)              // 200 OK
    return c.Created(data)         // 201 Created
    return c.NoContent()           // 204 No Content
    
    // 错误响应
    return c.BadRequest("消息")    // 400 Bad Request
    return c.Unauthorized("消息")  // 401 Unauthorized
    return c.Forbidden("消息")     // 403 Forbidden
    return c.NotFound("消息")      // 404 Not Found
    return c.Conflict("消息")      // 409 Conflict
    return c.InternalError("消息") // 500 Internal Server Error
    
    // 自定义错误
    return c.Error(422, "VALIDATION_ERROR", "验证失败")
}
```

## ActionResult

`IActionResult` 接口类似于 .NET 的 `IActionResult`，提供统一的响应处理机制。

### 统一响应格式

所有 ActionResult 使用统一的 JSON 格式：

```json
// 成功响应
{
  "success": true,
  "data": { "id": 1, "name": "Alice" }
}

// 错误响应
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "用户不存在"
  }
}
```

### 结果类型一览

| 方法 | HTTP 状态码 | 说明 |
|------|-------------|------|
| `Ok(data)` | 200 | 成功响应 |
| `Created(data)` | 201 | 创建成功 |
| `NoContent()` | 204 | 无内容响应 |
| `BadRequest(msg)` | 400 | 请求错误 |
| `BadRequestWithCode(code, msg)` | 400 | 请求错误（自定义错误码） |
| `Unauthorized(msg)` | 401 | 未授权 |
| `Forbidden(msg)` | 403 | 禁止访问 |
| `NotFound(msg)` | 404 | 资源不存在 |
| `Conflict(msg)` | 409 | 资源冲突 |
| `InternalError(msg)` | 500 | 服务器错误 |
| `Error(code, errCode, msg)` | 自定义 | 自定义错误 |
| `Redirect(url)` | 302 | 重定向 |
| `RedirectPermanent(url)` | 301 | 永久重定向 |
| `Json(code, data)` | 自定义 | 自定义 JSON |
| `Content(code, text)` | 自定义 | 纯文本响应 |
| `File(path)` | 200 | 文件响应 |
| `FileDownload(path, name)` | 200 | 文件下载 |
| `Status(code)` | 自定义 | 仅状态码 |

### 使用静态方法

除了通过 `HttpContext` 调用，也可以直接使用静态方法：

```go
app.MapGet("/example", func(c *web.HttpContext) web.IActionResult {
    // 通过 HttpContext
    return c.Ok(data)
    
    // 或直接使用静态方法
    return web.Ok(data)
    return web.NotFound("资源不存在")
    return web.Redirect("/new-location")
    return web.Json(201, customResponse)
})
```

### 完整示例

```go
app.MapPost("/users", func(c *web.HttpContext) web.IActionResult {
    // 1. 绑定请求
    var req CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    // 2. 验证业务逻辑
    if userExists(req.Email) {
        return c.Conflict("邮箱已被注册")
    }
    
    // 3. 创建用户
    user, err := createUser(req)
    if err != nil {
        return c.InternalError("创建用户失败")
    }
    
    // 4. 返回结果
    return c.Created(user)
})

app.MapGet("/users/:id", func(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    user := getUserByID(id)
    if user == nil {
        return c.NotFound("用户不存在")
    }
    
    return c.Ok(user)
})

app.MapDelete("/users/:id", func(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    if !deleteUser(id) {
        return c.NotFound("用户不存在")
    }
    
    return c.NoContent()
})
```

## 中间件

### 使用中间件

```go
app := builder.Build()

// 使用 Gin 内置中间件
app.Use(gin.Logger())
app.Use(gin.Recovery())

// 使用自定义中间件
app.Use(AuthMiddleware())
app.Use(CorsMiddleware())
```

### 自定义中间件

```go
// 传统方式
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "未授权"})
            return
        }
        c.Next()
    }
}

// 使用 HttpContext
func LoggingMiddleware() web.HandlerFunc {
    return func(c *web.HttpContext) {
        start := time.Now()
        c.Next()
        duration := time.Since(start)
        log.Printf("%s %s %d %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
    }
}
```

### 路由组中间件

```go
// 为特定路由组添加中间件
admin := app.MapGroup("/admin", AuthMiddleware(), AdminOnlyMiddleware())
admin.MapGet("/dashboard", GetDashboard)
admin.MapGet("/users", GetAdminUsers)

// 公开路由（无中间件）
public := app.MapGroup("/public")
public.MapGet("/health", HealthCheck)
```

## 控制器

### IController 接口

```go
// 定义控制器
type UserController struct {
    userService UserService
}

func NewUserController(userService UserService) *UserController {
    return &UserController{userService: userService}
}

// 实现 web.IController 接口
func (ctrl *UserController) MapRoutes(app *web.WebApplication) {
    users := app.MapGroup("/api/users")
    users.MapGet("", ctrl.GetAll)
    users.MapGet("/:id", ctrl.GetByID)
    users.MapPost("", ctrl.Create)
}

// 使用 ActionResult
func (ctrl *UserController) GetAll(c *web.HttpContext) web.IActionResult {
    users := ctrl.userService.ListUsers()
    return c.Ok(users)
}

func (ctrl *UserController) GetByID(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    user := ctrl.userService.GetUser(id)
    if user == nil {
        return c.NotFound("用户不存在")
    }
    
    return c.Ok(user)
}

func (ctrl *UserController) Create(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    user := ctrl.userService.Create(req)
    return c.Created(user)
}
```

### 注册控制器

```go
// 使用 AddController 注册
web.AddController(builder.Services, func(sp di.IServiceProvider) *UserController {
    return NewUserController(di.GetRequiredService[UserService](sp))
})

// 构建应用后映射控制器
app := builder.Build()
app.MapControllers()
```

### 控制器注册扩展

```go
// controllers/controller_extensions.go
func AddControllers(services di.IServiceCollection) {
    web.AddController(services, func(sp di.IServiceProvider) *UserController {
        return NewUserController(di.GetRequiredService[UserService](sp))
    })
    
    web.AddController(services, func(sp di.IServiceProvider) *OrderController {
        return NewOrderController(di.GetRequiredService[OrderService](sp))
    })
}

// main.go
controllers.AddControllers(builder.Services)
app := builder.Build()
app.MapControllers()
```

## 静态文件

### 提供静态文件

```go
app := builder.Build()

// 使用 UseStaticFiles 方法
app.UseStaticFiles(func(opts *web.StaticFileOptions) {
    opts.RequestPath = "/static"
    opts.FileSystem = "./public"
})
```

### SPA 支持

```go
// 为 SPA 应用提供前端文件
app.UseStaticFiles(func(opts *web.StaticFileOptions) {
    opts.RequestPath = "/"
    opts.FileSystem = "./frontend/dist"
})

// 处理 SPA 路由（所有未匹配的路由返回 index.html）
app.MapGet("/{catchAll}", func(c *web.HttpContext) web.IActionResult {
    return web.File("./frontend/dist/index.html")
})
```

## CORS

### 配置 CORS

```go
import "github.com/gocrud/csgo/web"

builder := web.CreateBuilder()

// 添加 CORS 服务
builder.AddCors(func(opts *web.CorsOptions) {
    opts.AllowOrigins = []string{"http://localhost:3000", "https://myapp.com"}
    opts.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
    opts.AllowHeaders = []string{"Content-Type", "Authorization"}
    opts.AllowCredentials = true
})

app := builder.Build()

// 启用 CORS 中间件
app.UseCors()
```

### 手动配置

```go
import "github.com/gin-contrib/cors"

app := builder.Build()

app.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}))
```

## 完整示例

### REST API 应用

```go
package main

import (
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/swagger"
    "github.com/gocrud/csgo/web"
)

// 模型
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

// 服务接口
type UserService interface {
    GetAll() []User
    GetByID(id int) *User
    Create(name, email string) *User
    Delete(id int) bool
}

// 服务实现
type userService struct {
    users  map[int]*User
    nextID int
}

func NewUserService() UserService {
    return &userService{
        users: map[int]*User{
            1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
            2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
        },
        nextID: 3,
    }
}

func (s *userService) GetAll() []User {
    users := make([]User, 0, len(s.users))
    for _, u := range s.users {
        users = append(users, *u)
    }
    return users
}

func (s *userService) GetByID(id int) *User {
    return s.users[id]
}

func (s *userService) Create(name, email string) *User {
    user := &User{ID: s.nextID, Name: name, Email: email}
    s.users[s.nextID] = user
    s.nextID++
    return user
}

func (s *userService) Delete(id int) bool {
    if _, ok := s.users[id]; ok {
        delete(s.users, id)
        return true
    }
    return false
}

// 控制器
type UserController struct {
    userService UserService
}

func NewUserController(userService UserService) *UserController {
    return &UserController{userService: userService}
}

func (ctrl *UserController) MapRoutes(app *web.WebApplication) {
    users := app.MapGroup("/api/users")
    users.WithTags("Users")
    
    users.MapGet("", ctrl.GetAll).
        WithSummary("获取所有用户")
    
    users.MapGet("/:id", ctrl.GetByID).
        WithSummary("根据 ID 获取用户")
    
    users.MapPost("", ctrl.Create).
        WithSummary("创建用户")
    
    users.MapDelete("/:id", ctrl.Delete).
        WithSummary("删除用户")
}

func (ctrl *UserController) GetAll(c *web.HttpContext) web.IActionResult {
    users := ctrl.userService.GetAll()
    return c.Ok(users)
}

func (ctrl *UserController) GetByID(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    user := ctrl.userService.GetByID(id)
    if user == nil {
        return c.NotFound("用户不存在")
    }
    
    return c.Ok(user)
}

func (ctrl *UserController) Create(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    user := ctrl.userService.Create(req.Name, req.Email)
    return c.Created(user)
}

func (ctrl *UserController) Delete(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    if !ctrl.userService.Delete(id) {
        return c.NotFound("用户不存在")
    }
    
    return c.NoContent()
}

func main() {
    // 1. 创建构建器
    builder := web.CreateBuilder()
    
    // 2. 注册服务
    builder.Services.AddSingleton(NewUserService)
    
    // 3. 注册控制器
    web.AddController(builder.Services, func(sp di.IServiceProvider) *UserController {
        return NewUserController(di.GetRequiredService[UserService](sp))
    })
    
    // 4. 配置 Swagger
    swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
        opts.Title = "User API"
        opts.Version = "v1"
        opts.Description = "用户管理 API - 使用 HttpContext 和 ActionResult"
    })
    
    // 5. 构建应用
    app := builder.Build()
    
    // 6. 配置中间件
    swagger.UseSwagger(app)
    swagger.UseSwaggerUI(app)
    
    // 7. 映射控制器
    app.MapControllers()
    
    // 8. 根路由
    app.MapGet("/", func(c *web.HttpContext) web.IActionResult {
        return c.Ok(gin.H{
            "message": "User API",
            "version": "v1",
            "docs":    "/swagger",
        })
    })
    
    // 9. 运行
    println("Server: http://localhost:8080")
    println("Swagger: http://localhost:8080/swagger")
    app.Run()
}
```

## 最佳实践

### 1. 项目结构

```
myapp/
├── main.go
├── controllers/
│   ├── user_controller.go
│   └── extensions.go
├── services/
│   ├── user_service.go
│   └── extensions.go
├── models/
│   └── user.go
├── middleware/
│   └── auth.go
└── config/
    └── appsettings.json
```

### 2. 使用 ActionResult 模式

```go
// ✅ 推荐：使用 ActionResult
func (ctrl *UserController) GetByID(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    user := ctrl.userService.GetUser(id)
    if user == nil {
        return c.NotFound("用户不存在")
    }
    
    return c.Ok(user)
}

// ❌ 不推荐：手动处理响应
func (ctrl *UserController) GetByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(400, gin.H{"error": "无效的 ID"})
        return
    }
    
    user := ctrl.userService.GetUser(id)
    if user == nil {
        c.JSON(404, gin.H{"error": "用户不存在"})
        return
    }
    
    c.JSON(200, user)
}
```

### 3. 服务注册扩展

```go
// services/extensions.go
func AddServices(services di.IServiceCollection) {
    services.AddSingleton(NewUserService)
    services.AddSingleton(NewOrderService)
    services.AddScoped(NewUnitOfWork)
}

// controllers/extensions.go
func AddControllers(services di.IServiceCollection) {
    web.AddController(services, NewUserController)
    web.AddController(services, NewOrderController)
}

// main.go
services.AddServices(builder.Services)
controllers.AddControllers(builder.Services)
```

### 4. 环境配置

```go
builder := web.CreateBuilder()

if builder.Environment.IsDevelopment() {
    // 开发环境配置
    builder.Services.AddSingleton(NewMockEmailService)
} else {
    // 生产环境配置
    builder.Services.AddSingleton(NewRealEmailService)
}
```

### 5. 错误处理中间件

```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            c.JSON(500, web.ApiResponse{
                Success: false,
                Error:   &web.ApiError{Code: "INTERNAL_ERROR", Message: err.Error()},
            })
        }
    }
}

app.Use(ErrorHandler())
```

## 与 .NET 对比

| .NET | CSGO | 说明 |
|------|-----|------|
| `WebApplication.CreateBuilder()` | `web.CreateBuilder()` | 创建构建器 |
| `builder.Services` | `builder.Services` | 服务集合 |
| `builder.Configuration` | `builder.Configuration` | 配置 |
| `builder.Build()` | `builder.Build()` | 构建应用 |
| `app.MapGet()` | `app.MapGet()` | 注册 GET 路由 |
| `app.MapPost()` | `app.MapPost()` | 注册 POST 路由 |
| `app.MapGroup()` | `app.MapGroup()` | 路由组 |
| `app.MapControllers()` | `app.MapControllers()` | 映射控制器 |
| `app.UseMiddleware()` | `app.Use()` | 使用中间件 |
| `app.Run()` | `app.Run()` | 运行应用 |
| `HttpContext` | `web.HttpContext` | HTTP 上下文 |
| `IActionResult` | `web.IActionResult` | 操作结果 |
| `Ok()` | `c.Ok()` / `web.Ok()` | 200 响应 |
| `NotFound()` | `c.NotFound()` / `web.NotFound()` | 404 响应 |
| `BadRequest()` | `c.BadRequest()` / `web.BadRequest()` | 400 响应 |

## 相关资源

- [控制器指南](controllers.md) - 控制器模式详解
- [依赖注入](dependency-injection.md) - DI 系统
- [配置管理](configuration.md) - 配置系统
- [应用托管](hosting.md) - 应用生命周期
- [API 文档](api-documentation.md) - Swagger 集成

---

**下一步**: 查看 [控制器指南](controllers.md) 了解如何组织 API 代码。
