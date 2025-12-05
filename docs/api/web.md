# Web 框架 API 参考

本文档提供 CSGO Web 框架的完整 API 参考。

## 目录

- [WebApplicationBuilder](#webapplicationbuilder)
- [WebApplication](#webapplication)
- [HttpContext](#httpcontext)
- [IActionResult](#iactionresult)
- [路由](#路由)
- [控制器](#控制器)

---

## WebApplicationBuilder

Web 应用构建器，用于配置服务和中间件。

```go
import "github.com/gocrud/csgo/web"
```

### CreateBuilder

创建 Web 应用构建器。

```go
func CreateBuilder(args ...string) *WebApplicationBuilder
```

**参数：**
- `args` - 可选的命令行参数

**返回值：**
- `*WebApplicationBuilder` - 构建器实例

**示例：**

```go
builder := web.CreateBuilder()
// 或带命令行参数
builder := web.CreateBuilder(os.Args...)
```

---

### 属性

| 属性 | 类型 | 说明 |
|------|------|------|
| `Services` | `di.IServiceCollection` | 服务集合 |
| `Configuration` | `configuration.IConfiguration` | 配置对象 |
| `Environment` | `hosting.IHostEnvironment` | 环境信息 |
| `Host` | `*ConfigureHostBuilder` | 主机配置 |
| `WebHost` | `*ConfigureWebHostBuilder` | Web 主机配置 |

---

### Build

构建 Web 应用。

```go
func (b *WebApplicationBuilder) Build() *WebApplication
```

**返回值：**
- `*WebApplication` - 已配置的 Web 应用

---

## WebApplication

已配置的 Web 应用实例。

### 属性

| 属性 | 类型 | 说明 |
|------|------|------|
| `Services` | `di.IServiceProvider` | 服务提供者 |

---

### Run

阻塞运行应用。

```go
func (app *WebApplication) Run(urls ...string) error
```

---

### Start / Stop

手动控制应用生命周期。

```go
func (app *WebApplication) Start(ctx context.Context) error
func (app *WebApplication) Stop(ctx context.Context) error
```

---

### Use

添加中间件。

```go
func (app *WebApplication) Use(middleware ...gin.HandlerFunc)
```

---

### MapGet / MapPost / MapPut / MapDelete / MapPatch

注册路由。支持三种处理器类型：

```go
func (app *WebApplication) MapGet(pattern string, handlers ...Handler) IEndpointConventionBuilder
func (app *WebApplication) MapPost(pattern string, handlers ...Handler) IEndpointConventionBuilder
func (app *WebApplication) MapPut(pattern string, handlers ...Handler) IEndpointConventionBuilder
func (app *WebApplication) MapDelete(pattern string, handlers ...Handler) IEndpointConventionBuilder
func (app *WebApplication) MapPatch(pattern string, handlers ...Handler) IEndpointConventionBuilder
```

**处理器类型：**

```go
// 1. gin.HandlerFunc
app.MapGet("/", func(c *gin.Context) {
    c.JSON(200, gin.H{"msg": "hello"})
})

// 2. func(*HttpContext)
app.MapGet("/", func(c *web.HttpContext) {
    c.JSON(200, gin.H{"msg": "hello"})
})

// 3. func(*HttpContext) IActionResult（推荐）
app.MapGet("/", func(c *web.HttpContext) web.IActionResult {
    return c.Ok(gin.H{"msg": "hello"})
})
```

---

### MapGroup

创建路由组。

```go
func (app *WebApplication) MapGroup(prefix string, handlers ...Handler) *RouteGroupBuilder
```

**示例：**

```go
api := app.MapGroup("/api")
api.MapGet("/users", GetUsers)
api.MapPost("/users", CreateUser)
```

---

### MapControllers

自动映射所有注册的控制器。

```go
func (app *WebApplication) MapControllers() *WebApplication
```

---

## HttpContext

HTTP 请求上下文，包装 `gin.Context`。

### 参数获取

| 方法 | 签名 | 说明 |
|------|------|------|
| `PathInt` | `PathInt(key string) (int, error)` | 获取路径参数转 int |
| `PathInt64` | `PathInt64(key string) (int64, error)` | 获取路径参数转 int64 |
| `MustPathInt` | `MustPathInt(key string) (int, IActionResult)` | 获取路径参数，失败返回错误结果 |
| `QueryInt` | `QueryInt(key string, defaultValue int) int` | 获取查询参数转 int |
| `QueryInt64` | `QueryInt64(key string, defaultValue int64) int64` | 获取查询参数转 int64 |
| `QueryBool` | `QueryBool(key string, defaultValue bool) bool` | 获取查询参数转 bool |

**示例：**

```go
func GetUser(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err  // 自动返回 400
    }
    
    page := c.QueryInt("page", 1)
    size := c.QueryInt("size", 10)
    
    return c.Ok(gin.H{"id": id, "page": page, "size": size})
}
```

---

### 请求绑定

| 方法 | 签名 | 说明 |
|------|------|------|
| `BindJSON` | `BindJSON(target interface{}) (bool, IActionResult)` | 绑定 JSON，返回成功标志和错误结果 |
| `MustBindJSON` | `MustBindJSON(target interface{}) IActionResult` | 绑定 JSON，失败返回错误结果 |
| `BindQuery` | `BindQuery(target interface{}) (bool, IActionResult)` | 绑定查询参数 |

**示例：**

```go
func CreateUser(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err  // 自动返回 400
    }
    
    // 处理请求...
    return c.Created(user)
}
```

---

### 响应方法

| 方法 | 签名 | HTTP 状态码 |
|------|------|-------------|
| `Ok` | `Ok(data interface{}) IActionResult` | 200 |
| `Created` | `Created(data interface{}) IActionResult` | 201 |
| `NoContent` | `NoContent() IActionResult` | 204 |
| `BadRequest` | `BadRequest(message string) IActionResult` | 400 |
| `BadRequestWithCode` | `BadRequestWithCode(code, message string) IActionResult` | 400 |
| `Unauthorized` | `Unauthorized(message string) IActionResult` | 401 |
| `Forbidden` | `Forbidden(message string) IActionResult` | 403 |
| `NotFound` | `NotFound(message string) IActionResult` | 404 |
| `Conflict` | `Conflict(message string) IActionResult` | 409 |
| `InternalError` | `InternalError(message string) IActionResult` | 500 |
| `Error` | `Error(statusCode int, code, message string) IActionResult` | 自定义 |

---

## IActionResult

操作结果接口，定义统一的响应格式。

### 响应格式

```json
// 成功
{
  "success": true,
  "data": { ... }
}

// 错误
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "错误描述"
  }
}
```

---

### 静态方法

除了通过 `HttpContext` 调用，也可以使用静态方法：

| 方法 | 说明 |
|------|------|
| `web.Ok(data)` | 200 成功 |
| `web.Created(data)` | 201 创建成功 |
| `web.NoContent()` | 204 无内容 |
| `web.BadRequest(msg)` | 400 请求错误 |
| `web.NotFound(msg)` | 404 未找到 |
| `web.Redirect(url)` | 302 重定向 |
| `web.RedirectPermanent(url)` | 301 永久重定向 |
| `web.Json(code, data)` | 自定义 JSON |
| `web.Content(code, text)` | 纯文本 |
| `web.File(path)` | 文件 |
| `web.FileDownload(path, name)` | 文件下载 |
| `web.Status(code)` | 仅状态码 |

---

## 路由

### IEndpointConventionBuilder

路由配置接口，用于添加路由元数据。

| 方法 | 说明 |
|------|------|
| `WithSummary(summary string)` | 设置 Swagger 摘要 |
| `WithDescription(desc string)` | 设置 Swagger 描述 |
| `WithTags(tags ...string)` | 设置 Swagger 标签 |
| `WithName(name string)` | 设置路由名称 |
| `RequireAuthorization(policies ...string)` | 要求授权 |

**示例：**

```go
app.MapGet("/users", GetUsers).
    WithSummary("获取所有用户").
    WithDescription("返回用户列表").
    WithTags("Users")
```

---

### RouteGroupBuilder

路由组构建器。

| 方法 | 说明 |
|------|------|
| `MapGet/Post/Put/Delete/Patch` | 在组内注册路由 |
| `MapGroup(prefix)` | 创建嵌套组 |
| `WithTags(tags ...string)` | 设置组标签 |
| `WithOpenApi()` | 启用 OpenAPI |
| `RequireAuthorization(policies ...string)` | 组授权 |

---

## 控制器

### IController

控制器接口。

```go
type IController interface {
    MapRoutes(app *WebApplication)
}
```

---

### AddController

注册控制器。

```go
func AddController[T any](services di.IServiceCollection, factory func(di.IServiceProvider) T)
```

**示例：**

```go
web.AddController(services, func(sp di.IServiceProvider) *UserController {
    return NewUserController(di.GetRequiredService[*UserService](sp))
})
```

---

### AddControllerInstance

注册控制器实例。

```go
func AddControllerInstance(services di.IServiceCollection, factory func(di.IServiceProvider) IController)
```

---

## 完整示例

```go
package main

import (
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type CreateUserRequest struct {
    Name string `json:"name" binding:"required"`
}

func main() {
    builder := web.CreateBuilder()
    app := builder.Build()
    
    users := app.MapGroup("/api/users").
        WithOpenApi(
            openapi.Tags("Users"),
        )
    
    users.MapGet("", func(c *web.HttpContext) web.IActionResult {
        return c.Ok([]User{{ID: 1, Name: "Alice"}})
    }).WithOpenApi(
        openapi.Summary("获取所有用户"),
        openapi.Produces[[]User](200),
    )
    
    users.MapGet("/:id", func(c *web.HttpContext) web.IActionResult {
        id, err := c.MustPathInt("id")
        if err != nil {
            return err
        }
        return c.Ok(User{ID: id, Name: "User"})
    }).WithOpenApi(
        openapi.Summary("获取用户"),
        openapi.Produces[User](200),
        openapi.ProducesProblem(404),
    )
    
    users.MapPost("", func(c *web.HttpContext) web.IActionResult {
        var req CreateUserRequest
        if err := c.MustBindJSON(&req); err != nil {
            return err
        }
        return c.Created(User{ID: 1, Name: req.Name})
    }).WithOpenApi(
        openapi.Summary("创建用户"),
        openapi.Accepts[CreateUserRequest]("application/json"),
        openapi.Produces[User](201),
    )
    
    users.MapDelete("/:id", func(c *web.HttpContext) web.IActionResult {
        return c.NoContent()
    }).WithOpenApi(
        openapi.Summary("删除用户"),
        openapi.Produces[any](204),
    )
    
    app.Run()
}
```

---

## 相关文档

- [Web 应用指南](../guides/web-applications.md)
- [控制器指南](../guides/controllers.md)
- [依赖注入 API](di.md)

