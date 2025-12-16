# Web 框架

[← 返回主目录](../README.md)

CSGO Web 框架基于 Gin，提供了现代化的 Web 应用开发体验，包括路由、控制器、中间件、请求验证等完整功能。

## 特性

- ✅ 简洁的应用构建器（WebApplicationBuilder）
- ✅ HttpContext 和 ActionResult 模式
- ✅ 类型安全的路由系统
- ✅ 控制器模式支持
- ✅ 中间件管道
- ✅ 自动请求验证
- ✅ 统一的 API 响应格式
- ✅ CORS 支持
- ✅ 静态文件服务
- ✅ 依赖注入集成

## 快速开始

### 1. 创建第一个应用

```go
package main

import (
    "github.com/gocrud/csgo/web"
)

func main() {
    // 创建应用构建器
    builder := web.CreateBuilder()
    
    // 构建应用
    app := builder.Build()
    
    // 定义路由
    app.MapGet("/", func(c *web.HttpContext) web.IActionResult {
        return c.Ok(web.M{"message": "Hello, CSGO!"})
    })
    
    // 运行应用
    app.Run()  // 默认监听 :8080
}
```

### 2. 使用依赖注入

```go
// 定义服务
type UserService struct{}

func NewUserService() *UserService {
    return &UserService{}
}

func (s *UserService) GetUser(id int) string {
    return fmt.Sprintf("User %d", id)
}

func main() {
    builder := web.CreateBuilder()
    
    // 注册服务
    builder.Services.Add(NewUserService)
    
    app := builder.Build()
    
    // 在路由中使用服务
    app.MapGet("/users/:id", func(c *web.HttpContext) web.IActionResult {
        userService := di.Get[*UserService](c.Services)
        id := c.Params().PathInt("id").Value()
        user := userService.GetUser(id)
        return c.Ok(web.M{"user": user})
    })
    
    app.Run()
}
```

### 3. 使用控制器

```go
type UserController struct {
    userService *UserService
}

func NewUserController(userService *UserService) *UserController {
    return &UserController{userService: userService}
}

func (ctrl *UserController) GetUser(c *web.HttpContext) web.IActionResult {
    id := c.Params().PathInt("id").Value()
    user := ctrl.userService.GetUser(id)
    return c.Ok(user)
}

// 在 main 中注册
builder := web.CreateBuilder()
builder.Services.Add(NewUserService)
web.AddController(builder.Services, NewUserController)

app := builder.Build()
app.MapControllers()
app.Run()
```

## WebApplicationBuilder

### 创建构建器

```go
// 创建默认构建器
builder := web.CreateBuilder()

// 传入命令行参数
builder := web.CreateBuilder(os.Args[1:]...)
```

构建器会自动：
- 加载配置（appsettings.json）
- 设置环境（Development/Production）
- 注册基础服务（日志、配置等）
- 初始化依赖注入容器

### 配置服务

```go
builder := web.CreateBuilder()

// 注册服务
builder.Services.Add(NewUserService)
builder.Services.Add(NewOrderService)

// 注册配置选项
var dbConfig DatabaseConfig
builder.Configuration.Bind("database", &dbConfig)
builder.Services.AddInstance(&dbConfig)
```

### 配置主机

```go
builder := web.CreateBuilder()

// 配置监听地址
builder.WebHost.UseUrls("http://localhost:5000")

// 配置关闭超时
builder.WebHost.UseShutdownTimeout(30)
```

### 访问配置和环境

```go
builder := web.CreateBuilder()

// 访问配置
port := builder.Configuration.GetInt("server:port", 8080)
dbConn := builder.Configuration.Get("database:connection")

// 访问环境
if builder.Environment.IsDevelopment() {
    // 开发环境特定配置
}
```

### 构建应用

```go
app := builder.Build()  // 构建 WebApplication 实例
```

## WebApplication

### 运行应用

```go
app := builder.Build()

// 方式 1：使用默认地址运行
app.Run()  // 默认 :8080

// 方式 2：指定地址运行
app.Run("http://localhost:5000")

// 方式 3：使用 Context 运行
ctx := context.Background()
app.RunWithContext(ctx)

// 方式 4：手动控制生命周期
ctx := context.Background()
app.Start(ctx)
// ... 做其他事情
app.Stop(ctx)
```

### 访问服务

```go
app := builder.Build()

// 从应用的服务容器解析服务
userService := di.Get[*UserService](app.Services)
config := di.Get[*AppConfig](app.Services)
```

## 路由系统

### 基本路由

```go
app := builder.Build()

// GET 请求
app.MapGet("/hello", func(c *web.HttpContext) web.IActionResult {
    return c.Ok(web.M{"message": "Hello"})
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
    // 更新逻辑
    return c.Ok(nil)
})

// DELETE 请求
app.MapDelete("/users/:id", func(c *web.HttpContext) web.IActionResult {
    // 删除逻辑
    return c.NoContent()
})

// PATCH 请求
app.MapPatch("/users/:id", func(c *web.HttpContext) web.IActionResult {
    // 部分更新逻辑
    return c.Ok(nil)
})
```

### 路径参数

```go
// 定义路径参数
app.MapGet("/users/:id", func(c *web.HttpContext) web.IActionResult {
    // 获取路径参数
    id := c.RawCtx().Param("id")
    
    // 或使用参数验证器
    idInt := c.Params().PathInt("id").Value()
    
    return c.Ok(web.M{"id": idInt})
})

// 多个路径参数
app.MapGet("/users/:userId/orders/:orderId", 
    func(c *web.HttpContext) web.IActionResult {
        userId := c.RawCtx().Param("userId")
        orderId := c.RawCtx().Param("orderId")
        return c.Ok(web.M{"userId": userId, "orderId": orderId})
    })
```

### 查询参数

```go
app.MapGet("/search", func(c *web.HttpContext) web.IActionResult {
    // 获取单个查询参数
    keyword := c.RawCtx().Query("keyword")
    
    // 获取带默认值的查询参数
    page := c.RawCtx().DefaultQuery("page", "1")
    
    // 绑定到结构体
    var query SearchQuery
    if ok, err := c.BindQuery(&query); !ok {
        return err
    }
    
    return c.Ok(query)
})

type SearchQuery struct {
    Keyword string `form:"keyword"`
    Page    int    `form:"page"`
    Size    int    `form:"size"`
}
```

### 路由组

```go
app := builder.Build()

// 创建 API 路由组
api := app.MapGroup("/api")
{
    // /api/users
    api.MapGet("/users", getUsers)
    api.MapPost("/users", createUser)
    
    // /api/orders
    api.MapGet("/orders", getOrders)
}

// 创建带版本的路由组
v1 := app.MapGroup("/api/v1")
{
    v1.MapGet("/users", getUsersV1)
}

v2 := app.MapGroup("/api/v2")
{
    v2.MapGet("/users", getUsersV2)
}

// 嵌套路由组
api := app.MapGroup("/api")
users := api.MapGroup("/users")
{
    users.MapGet("", listUsers)          // GET /api/users
    users.MapGet("/:id", getUser)        // GET /api/users/:id
    users.MapPost("", createUser)        // POST /api/users
    users.MapPut("/:id", updateUser)     // PUT /api/users/:id
    users.MapDelete("/:id", deleteUser)  // DELETE /api/users/:id
}
```

### 路由组中间件

```go
// 为路由组添加中间件
api := app.MapGroup("/api", authMiddleware, loggingMiddleware)
{
    api.MapGet("/users", getUsers)  // 会应用中间件
}

// 或者
api := app.MapGroup("/api")
api.Use(authMiddleware)  // 添加中间件到组
api.MapGet("/users", getUsers)
```

## HttpContext

### 获取请求信息

```go
func handler(c *web.HttpContext) web.IActionResult {
    // 获取原始 gin.Context
    ginCtx := c.RawCtx()
    
    // 获取请求 Context
    ctx := c.Context()
    
    // 获取请求方法
    method := ginCtx.Request.Method
    
    // 获取请求路径
    path := ginCtx.Request.URL.Path
    
    // 获取请求头
    userAgent := ginCtx.GetHeader("User-Agent")
    
    // 获取 Cookie
    token, err := ginCtx.Cookie("token")
    
    return c.Ok(nil)
}
```

### 请求体绑定

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func createUser(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    
    // 绑定 JSON（返回错误则自动返回 400）
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    // 或者使用两个返回值的方式
    if ok, err := c.BindJSON(&req); !ok {
        return err
    }
    
    // 使用请求数据
    user := createUserFromRequest(req)
    return c.Created(user)
}
```

### 请求验证

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

// 注册验证器
func init() {
    validator := validation.NewValidator[CreateUserRequest]()
    validator.Field(func(r *CreateUserRequest) string { return r.Name }).
        NotEmpty().
        MinLength(2)
    validator.Field(func(r *CreateUserRequest) string { return r.Email }).
        NotEmpty().
        EmailAddress()
    validation.RegisterValidator[CreateUserRequest](validator)
}

// 使用自动验证
func createUser(c *web.HttpContext) web.IActionResult {
    // 自动绑定并验证
    req, err := web.BindAndValidate[CreateUserRequest](c)
    if err != nil {
        return err  // 自动返回验证错误
    }
    
    // 验证通过，处理业务逻辑
    user := createUserFromRequest(*req)
    return c.Created(user)
}
```

### 访问服务

```go
func handler(c *web.HttpContext) web.IActionResult {
    // 从 HttpContext 访问服务容器
    userService := di.Get[*UserService](c.Services)
    
    // 使用服务
    users := userService.GetAllUsers()
    
    return c.Ok(users)
}
```

## ActionResult

### 成功响应

```go
// 200 OK
return c.Ok(data)
return c.Ok(web.M{"message": "Success"})

// 201 Created
return c.Created(user)

// 204 No Content
return c.NoContent()
```

**响应格式：**

```json
{
  "success": true,
  "data": { /* 你的数据 */ }
}
```

### 错误响应

```go
// 400 Bad Request
return c.BadRequest("无效的请求参数")

// 401 Unauthorized
return c.Unauthorized("未授权访问")

// 403 Forbidden
return c.Forbidden("没有访问权限")

// 404 Not Found
return c.NotFound("资源不存在")

// 409 Conflict
return c.Conflict("资源冲突")

// 500 Internal Server Error
return c.InternalError("服务器内部错误")

// 自定义错误
return c.Error(418, "I_AM_TEAPOT", "我是一个茶壶")
```

**错误响应格式：**

```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "资源不存在"
  }
}
```

### 验证错误响应

```go
// 验证失败时自动返回
req, err := web.BindAndValidate[CreateUserRequest](c)
if err != nil {
    return err  // 自动格式化验证错误
}
```

**验证错误响应格式：**

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION.FAILED",
    "message": "验证失败",
    "fields": [
      {
        "field": "name",
        "message": "不能为空",
        "code": "VALIDATION.REQUIRED"
      },
      {
        "field": "email",
        "message": "邮箱格式不正确",
        "code": "VALIDATION.EMAIL"
      }
    ]
  }
}
```

### 业务错误响应

**推荐方式：使用 FromError（简洁）**

```go
import "github.com/gocrud/csgo/errors"

func getUser(c *web.HttpContext) web.IActionResult {
    user, err := userService.GetUser(id)
    if err != nil {
        // FromError 自动识别错误类型并返回对应的响应
        // BizError -> 自动映射状态码，ValidationErrors -> 400，普通 error -> 500
        return c.FromError(err, "获取用户失败")
    }
    return c.Ok(user)
}

// 服务层
func (s *UserService) GetUser(id int) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    if user == nil {
        // 使用业务错误构建器
        return nil, errors.Business("USER").NotFound("用户不存在")
    }
    return user, nil
}
```

**传统方式：手动类型判断（仍然支持）**

```go
func getUser(c *web.HttpContext) web.IActionResult {
    user, err := userService.GetUser(id)
    if err != nil {
        // 手动判断错误类型
        if bizErr, ok := err.(*errors.BizError); ok {
            return c.BizError(bizErr)
        }
        return c.InternalError("服务器错误")
    }
    return c.Ok(user)
}
```

**自定义错误处理器**

```go
// 在应用启动时注册
func init() {
    // 注册数据库错误处理器
    web.RegisterErrorHandler(
        func(err error) bool {
            return errors.Is(err, sql.ErrNoRows)
        },
        func(err error, msg ...string) web.IActionResult {
            return web.Error(404, "NOT_FOUND", "记录不存在")
        },
    )
}

// 控制器中使用
func getUser(c *web.HttpContext) web.IActionResult {
    user, err := repo.FindByID(id)  // 可能返回 sql.ErrNoRows
    if err != nil {
        return c.FromError(err, "用户不存在")  // 自动应用处理器
    }
    return c.Ok(user)
}
```

### 其他响应类型

```go
// 重定向
return web.Redirect("/new-url")
return web.RedirectPermanent("/new-url")

// 纯文本
return web.Content(200, "Plain text response")

// 自定义 JSON（不使用标准格式）
return web.Json(200, web.M{"custom": "format"})

// 文件下载
return web.File("/path/to/file.pdf")
return web.FileDownload("/path/to/file.pdf", "download.pdf")

// 图片响应（二进制流）
return web.PNG(imageData)
return web.JPEG(imageData)
return web.WebP(imageData)
return web.BinaryImage(imageData, "image/gif")

// 图片响应（Base64编码的JSON）
return web.Base64Image(imageData, "image/png")

// 仅状态码
return web.Status(204)
```

### 图片响应详解

框架提供了专门的图片响应方法：

```go
// 方式1：二进制图片流（直接返回图片数据）
func getAvatar(c *web.HttpContext) web.IActionResult {
    imageData, _ := loadImageFromDB()
    return web.PNG(imageData)  // 返回PNG格式
}

// 方式2：Base64编码（包含在JSON中）
func getThumbnail(c *web.HttpContext) web.IActionResult {
    imageData, _ := loadThumbnailFromDB()
    return web.Base64Image(imageData, "image/png")
}
// 响应格式：{"success": true, "data": {"image": "base64...", "contentType": "image/png"}}

// 所有支持的图片格式
web.PNG(imageData)           // image/png
web.JPEG(imageData)          // image/jpeg
web.WebP(imageData)          // image/webp
web.BinaryImage(data, type)  // 自定义类型
```

**OpenAPI 配置：**

```go
import "github.com/gocrud/csgo/openapi"

// 二进制图片
app.MapGet("/api/images/logo", getLogoPNG).
    WithOpenApi(
        openapi.OptSummary("获取Logo"),
        openapi.OptBinaryImageResponse("image/png"),
    )

// Base64 图片
app.MapGet("/api/images/thumbnail/:id", getThumbnail).
    WithOpenApi(
        openapi.OptSummary("获取缩略图"),
        openapi.OptBase64ImageResponse(),
        openapi.OptPath[int]("id", "图片ID"),
    )
```

## 控制器模式

### 定义控制器

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

// 实现 IController 接口
func (ctrl *UserController) MapRoutes(app *web.WebApplication) {
    users := app.MapGroup("/api/users")
    users.MapGet("", ctrl.List)
    users.MapGet("/:id", ctrl.Get)
    users.MapPost("", ctrl.Create)
    users.MapPut("/:id", ctrl.Update)
    users.MapDelete("/:id", ctrl.Delete)
}

// 控制器方法
func (ctrl *UserController) List(c *web.HttpContext) web.IActionResult {
    users := ctrl.userService.GetAll()
    return c.Ok(users)
}

func (ctrl *UserController) Get(c *web.HttpContext) web.IActionResult {
    id := c.Params().PathInt("id").Value()
    user, err := ctrl.userService.GetByID(id)
    if err != nil {
        if bizErr, ok := err.(*errors.BizError); ok {
            return c.BizError(bizErr)
        }
        return c.InternalError("服务器错误")
    }
    return c.Ok(user)
}

func (ctrl *UserController) Create(c *web.HttpContext) web.IActionResult {
    req, err := web.BindAndValidate[CreateUserRequest](c)
    if err != nil {
        return err
    }
    
    user, err := ctrl.userService.Create(req)
    if err != nil {
        return c.handleError(err)
    }
    
    return c.Created(user)
}

// 统一错误处理
func (ctrl *UserController) handleError(err error) web.IActionResult {
    if bizErr, ok := err.(*errors.BizError); ok {
        return web.BizError(bizErr)
    }
    ctrl.logger.LogError(err, "Unexpected error")
    return web.InternalError("服务器错误")
}
```

### 注册控制器

```go
func main() {
    builder := web.CreateBuilder()
    
    // 注册服务
    builder.Services.Add(NewUserService)
    
    // 注册控制器
    web.AddController(builder.Services, NewUserController)
    
    app := builder.Build()
    
    // 自动映射所有控制器路由
    app.MapControllers()
    
    app.Run()
}
```

## 中间件

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
// 日志中间件
func loggingMiddleware(c *gin.Context) {
    start := time.Now()
    
    // 处理请求
    c.Next()
    
    // 请求完成后
    latency := time.Since(start)
    status := c.Writer.Status()
    
    fmt.Printf("[%s] %s %d %v\n", 
        c.Request.Method,
        c.Request.URL.Path,
        status,
        latency,
    )
}

// 认证中间件
func authMiddleware(c *gin.Context) {
    token := c.GetHeader("Authorization")
    
    if token == "" {
        c.JSON(401, web.M{"error": "Unauthorized"})
        c.Abort()  // 停止后续处理
        return
    }
    
    // 验证 token
    user, err := validateToken(token)
    if err != nil {
        c.JSON(401, web.M{"error": "Invalid token"})
        c.Abort()
        return
    }
    
    // 设置用户信息到上下文
    c.Set("user", user)
    c.Next()
}

// 使用
app.Use(loggingMiddleware)
app.Use(authMiddleware)
```

### 路由级中间件

```go
// 只应用到特定路由
app.MapGet("/admin/users", authMiddleware, getAdminUsers)

// 应用到路由组
admin := app.MapGroup("/admin", authMiddleware)
{
    admin.MapGet("/users", getAdminUsers)
    admin.MapPost("/users", createAdminUser)
}
```

### 恢复中间件

Gin 默认包含恢复中间件，捕获 panic 并返回 500：

```go
// 已自动启用，无需手动添加
// 如果 panic，会自动返回 500 错误
```

## CORS 配置

### 启用 CORS

```go
builder := web.CreateBuilder()

// 添加 CORS 支持
builder.AddCors(func(opts *web.CorsOptions) {
    opts.AllowOrigins = []string{"http://localhost:3000"}
    opts.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
    opts.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
    opts.AllowCredentials = true
    opts.MaxAge = 12 * time.Hour
})

app := builder.Build()

// 使用 CORS 中间件
app.UseCors()

app.MapGet("/api/users", getUsers)
app.Run()
```

### 开发环境 CORS

```go
builder := web.CreateBuilder()

if builder.Environment.IsDevelopment() {
    // 开发环境允许所有源
    builder.AddCors(func(opts *web.CorsOptions) {
        opts.AllowAllOrigins = true
        opts.AllowMethods = []string{"*"}
        opts.AllowHeaders = []string{"*"}
    })
}

app := builder.Build()
app.UseCors()
```

## 静态文件

### 提供静态文件

```go
app := builder.Build()

// 提供静态文件目录
app.ServeStaticFiles("/static", "./public")
// 访问：http://localhost:8080/static/image.jpg -> ./public/image.jpg

// 提供单个文件
app.ServeStaticFile("/favicon.ico", "./assets/favicon.ico")

// SPA 应用支持
app.ServeStaticFiles("/", "./dist")
app.ServeSPA("./dist/index.html")  // 所有未匹配路由返回 index.html
```

## 最佳实践

### 1. 使用 HttpContext 和 ActionResult

```go
// ✅ 推荐：使用 HttpContext 和 ActionResult
func handler(c *web.HttpContext) web.IActionResult {
    return c.Ok(data)
}

// ❌ 不推荐：直接使用 gin.Context
func handler(c *gin.Context) {
    c.JSON(200, web.M{"data": data})
}
```

### 2. 统一响应格式

```go
// ✅ 使用 ActionResult，自动格式化响应
return c.Ok(user)
return c.BadRequest("Invalid input")
return c.NotFound("User not found")

// ❌ 手动构建响应
c.RawCtx().JSON(200, web.M{"success": true, "data": user})
```

### 3. 使用控制器组织代码

```go
// ✅ 推荐：使用控制器
type UserController struct {
    service *UserService
}

func (ctrl *UserController) MapRoutes(app *web.WebApplication) {
    users := app.MapGroup("/api/users")
    users.MapGet("", ctrl.List)
    users.MapPost("", ctrl.Create)
}

// ❌ 不推荐：所有路由在 main 中定义
func main() {
    app.MapGet("/api/users", func(...) {...})
    app.MapPost("/api/users", func(...) {...})
    // 大量路由定义...
}
```

### 4. 验证器复用

```go
// ✅ 定义并注册验证器
func init() {
    validation.RegisterValidator[CreateUserRequest](NewCreateUserValidator())
}

// 在多个地方使用
func createUser(c *web.HttpContext) web.IActionResult {
    req, err := web.BindAndValidate[CreateUserRequest](c)
    if err != nil {
        return err
    }
    // ...
}
```

### 5. 错误处理分层

```go
// ✅ 推荐：服务层抛出业务错误，控制器使用 FromError
// 服务层：抛出业务错误
func (s *UserService) GetUser(id int) (*User, error) {
    if user == nil {
        return nil, errors.Business("USER").NotFound("用户不存在")
    }
    return user, nil
}

// 控制器层：使用 FromError 自动处理
func (ctrl *UserController) GetUser(c *web.HttpContext) web.IActionResult {
    user, err := ctrl.service.GetUser(id)
    if err != nil {
        return c.FromError(err, "获取用户失败")  // 一行搞定！
    }
    return c.Ok(user)
}

// ❌ 不推荐：手动类型判断（样板代码多）
func (ctrl *UserController) GetUser(c *web.HttpContext) web.IActionResult {
    user, err := ctrl.service.GetUser(id)
    if err != nil {
        if bizErr, ok := err.(*errors.BizError); ok {
            return c.BizError(bizErr)
        }
        return c.InternalError("服务器错误")
    }
    return c.Ok(user)
}
```

### 6. 使用路由组组织 API

```go
// ✅ 推荐：使用路由组
api := app.MapGroup("/api")
v1 := api.MapGroup("/v1")
{
    users := v1.MapGroup("/users")
    users.MapGet("", listUsers)
    users.MapPost("", createUser)
    
    orders := v1.MapGroup("/orders")
    orders.MapGet("", listOrders)
}

// ❌ 不推荐：扁平化路由
app.MapGet("/api/v1/users", listUsers)
app.MapPost("/api/v1/users", createUser)
app.MapGet("/api/v1/orders", listOrders)
```

### 7. 中间件顺序

```go
app := builder.Build()

// 正确的中间件顺序
app.Use(recoveryMiddleware)      // 1. 异常恢复
app.Use(loggingMiddleware)        // 2. 日志记录
app.Use(corsMiddleware)           // 3. CORS
app.Use(authMiddleware)           // 4. 认证
app.Use(rateLimitMiddleware)      // 5. 限流

// 定义路由
app.MapGet("/api/users", getUsers)
```

## API 参考

### WebApplicationBuilder

```go
// 创建构建器
CreateBuilder(args ...string) *WebApplicationBuilder

// 访问属性
builder.Services      // IServiceCollection
builder.Configuration // IConfigurationManager
builder.Environment   // IHostEnvironment
builder.WebHost       // WebHost 配置

// 构建应用
builder.Build() *WebApplication
```

### WebApplication

```go
// 运行应用
Run(urls ...string) error
RunWithContext(ctx context.Context) error
Start(ctx context.Context) error
Stop(ctx context.Context) error

// 路由
MapGet(pattern string, handlers ...Handler) IEndpointConventionBuilder
MapPost(pattern string, handlers ...Handler) IEndpointConventionBuilder
MapPut(pattern string, handlers ...Handler) IEndpointConventionBuilder
MapDelete(pattern string, handlers ...Handler) IEndpointConventionBuilder
MapPatch(pattern string, handlers ...Handler) IEndpointConventionBuilder
MapGroup(prefix string, handlers ...Handler) *RouteGroupBuilder

// 控制器
MapControllers()

// 中间件
Use(middleware ...gin.HandlerFunc)

// 静态文件
ServeStaticFiles(prefix, root string)
ServeStaticFile(path, filepath string)
```

### HttpContext

```go
// 访问原始上下文
RawCtx() *gin.Context
Context() context.Context

// 服务容器
Services di.IServiceProvider

// 成功响应
Ok(data interface{}) IActionResult
Created(data interface{}) IActionResult
NoContent() IActionResult

// 错误响应
BadRequest(message string) IActionResult
Unauthorized(message string) IActionResult
Forbidden(message string) IActionResult
NotFound(message string) IActionResult
InternalError(message string) IActionResult

// 绑定
BindJSON(target interface{}) (bool, IActionResult)
MustBindJSON(target interface{}) IActionResult
BindQuery(target interface{}) (bool, IActionResult)

// 验证
BindAndValidate[T any](c *HttpContext) (*T, IActionResult)

// 图片响应
PNG(imageData []byte) IActionResult
JPEG(imageData []byte) IActionResult
WebP(imageData []byte) IActionResult
BinaryImage(imageData []byte, contentType string) IActionResult
Base64Image(imageData []byte, contentType string) IActionResult
```

## 常见问题

### 如何获取请求头？

```go
func handler(c *web.HttpContext) web.IActionResult {
    token := c.RawCtx().GetHeader("Authorization")
    userAgent := c.RawCtx().GetHeader("User-Agent")
    return c.Ok(nil)
}
```

### 如何设置响应头？

```go
func handler(c *web.HttpContext) web.IActionResult {
    c.RawCtx().Header("X-Custom-Header", "value")
    return c.Ok(data)
}
```

### 如何处理文件上传？

```go
func uploadFile(c *web.HttpContext) web.IActionResult {
    file, err := c.RawCtx().FormFile("file")
    if err != nil {
        return c.BadRequest("No file uploaded")
    }
    
    // 保存文件
    dst := fmt.Sprintf("./uploads/%s", file.Filename)
    if err := c.RawCtx().SaveUploadedFile(file, dst); err != nil {
        return c.InternalError("Failed to save file")
    }
    
    return c.Ok(web.M{"filename": file.Filename})
}
```

### 如何在中间件中传递数据？

```go
// 在中间件中设置
func authMiddleware(c *gin.Context) {
    user := getUserFromToken(c)
    c.Set("user", user)
    c.Next()
}

// 在处理器中获取
func handler(c *web.HttpContext) web.IActionResult {
    user, exists := c.RawCtx().Get("user")
    if !exists {
        return c.Unauthorized("Not authenticated")
    }
    return c.Ok(user)
}
```

### HttpContext 和 gin.Context 的关系？

HttpContext 包装了 gin.Context，提供了统一的 API。你可以通过 `c.RawCtx()` 访问原始的 gin.Context：

```go
func handler(c *web.HttpContext) web.IActionResult {
    ginCtx := c.RawCtx()  // 获取原始 gin.Context
    // 使用 gin.Context 的所有方法
    ginCtx.JSON(200, data)
    return nil
}
```

---

[← 返回主目录](../README.md)

