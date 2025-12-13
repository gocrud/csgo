# Swagger 集成

[← 返回主目录](../README.md)

Swagger 模块提供了自动生成 OpenAPI 文档和 Swagger UI 界面的功能，让您的 API 拥有可交互的在线文档。

## 特性

- ✅ 自动生成 OpenAPI 3.0 规范
- ✅ Swagger UI 可视化界面
- ✅ 路由自动发现
- ✅ 请求/响应示例
- ✅ 安全定义支持（Bearer Token、API Key 等）
- ✅ 标签分组
- ✅ 自定义配置

## 快速开始

### 1. 启用 Swagger

```go
package main

import (
    "github.com/gocrud/csgo/swagger"
    "github.com/gocrud/csgo/web"
)

func main() {
    builder := web.CreateBuilder()
    
    // 添加 Swagger 生成服务
    swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
        opts.Title = "我的 API"
        opts.Version = "v1"
        opts.Description = "这是我的 API 文档"
    })
    
    app := builder.Build()
    
    // 定义路由
    app.MapGet("/api/users", getUsers)
    app.MapPost("/api/users", createUser)
    
    // 启用 Swagger JSON 端点
    swagger.UseSwagger(app)
    
    // 启用 Swagger UI
    swagger.UseSwaggerUI(app)
    
    app.Run()
}
```

访问 Swagger UI：http://localhost:8080/swagger

### 2. 添加路由注释

使用注释标记路由以生成更详细的文档：

```go
// @Summary 获取用户列表
// @Description 获取所有用户的列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Router /api/users [get]
app.MapGet("/api/users", getUsers)

// @Summary 创建用户
// @Description 创建一个新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "创建用户请求"
// @Success 201 {object} User
// @Failure 400 {object} ApiError
// @Router /api/users [post]
app.MapPost("/api/users", createUser)
```

## 配置

### SwaggerGenOptions

配置 Swagger 文档生成选项：

```go
swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
    // API 信息
    opts.Title = "电商 API"
    opts.Version = "v1.0.0"
    opts.Description = "电商平台 REST API 文档"
})
```

### SwaggerUIOptions

配置 Swagger UI 界面：

```go
swagger.UseSwaggerUI(app, func(opts *swagger.SwaggerUIOptions) {
    // UI 路径
    opts.RoutePrefix = "/docs"  // 默认 /swagger
    
    // Swagger JSON 路径
    opts.SpecURL = "/swagger/v1/swagger.json"
    
    // 页面标题
    opts.Title = "我的 API 文档"
})

// 访问：http://localhost:8080/docs
```

## 安全定义

### Bearer Token 认证

```go
import "github.com/gocrud/csgo/openapi"

swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
    opts.Title = "安全 API"
    opts.Version = "v1"
    
    // 添加 Bearer Token 安全方案
    opts.AddSecurityDefinition("Bearer", openapi.SecurityScheme{
        Type:         "http",
        Scheme:       "bearer",
        BearerFormat: "JWT",
        Description:  "输入 JWT Token",
    })
})
```

### API Key 认证

```go
swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
    opts.Title = "API Key API"
    opts.Version = "v1"
    
    // 添加 API Key 安全方案
    opts.AddSecurityDefinition("ApiKey", openapi.SecurityScheme{
        Type:        "apiKey",
        In:          "header",
        Name:        "X-API-Key",
        Description: "API Key 认证",
    })
})
```

### OAuth2 认证

```go
swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
    opts.Title = "OAuth2 API"
    opts.Version = "v1"
    
    // 添加 OAuth2 安全方案
    opts.AddSecurityDefinition("OAuth2", openapi.SecurityScheme{
        Type: "oauth2",
        Flows: &openapi.OAuthFlows{
            AuthorizationCode: &openapi.OAuthFlow{
                AuthorizationURL: "https://auth.example.com/oauth/authorize",
                TokenURL:         "https://auth.example.com/oauth/token",
                Scopes: map[string]string{
                    "read":  "读取权限",
                    "write": "写入权限",
                },
            },
        },
    })
})
```

## 路由标签

### 自动标签分组

框架会自动根据路由路径生成标签：

```go
// 自动标签：users
app.MapGet("/api/users", getUsers)

// 自动标签：orders
app.MapGet("/api/orders", getOrders)

// 自动标签：products
app.MapGet("/api/products", getProducts)
```

### 自定义标签

使用路由元数据自定义标签：

```go
// 设置自定义标签
app.MapGet("/api/users", getUsers).
    WithTags("用户管理", "核心接口")

app.MapGet("/api/admin/users", getAdminUsers).
    WithTags("管理员", "用户管理")
```

## 响应示例

### 定义响应模型

```go
type User struct {
    ID    int    `json:"id" example:"1"`
    Name  string `json:"name" example:"张三"`
    Email string `json:"email" example:"zhangsan@example.com"`
    Age   int    `json:"age" example:"25"`
}

type ApiResponse struct {
    Success bool        `json:"success" example:"true"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ApiError   `json:"error,omitempty"`
}

type ApiError struct {
    Code    string `json:"code" example:"NOT_FOUND"`
    Message string `json:"message" example:"资源不存在"`
}
```

### 使用注释定义响应

```go
// @Summary 获取用户
// @Description 根据 ID 获取用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户 ID"
// @Success 200 {object} ApiResponse{data=User}
// @Failure 404 {object} ApiResponse{error=ApiError}
// @Failure 500 {object} ApiResponse{error=ApiError}
// @Router /api/users/{id} [get]
func getUser(c *web.HttpContext) web.IActionResult {
    // 实现逻辑
}
```

## 完整示例

### 用户管理 API

```go
package main

import (
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/openapi"
    "github.com/gocrud/csgo/swagger"
    "github.com/gocrud/csgo/web"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    builder := web.CreateBuilder()
    
    // 配置 Swagger
    swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
        opts.Title = "用户管理 API"
        opts.Version = "v1.0.0"
        opts.Description = "用户管理系统 REST API"
        
        // 添加 Bearer Token 认证
        opts.AddSecurityDefinition("Bearer", openapi.SecurityScheme{
            Type:         "http",
            Scheme:       "bearer",
            BearerFormat: "JWT",
        })
    })
    
    app := builder.Build()
    
    // 定义 API 路由
    api := app.MapGroup("/api")
    users := api.MapGroup("/users")
    {
        // @Summary 获取用户列表
        // @Tags 用户管理
        // @Produce json
        // @Success 200 {array} User
        // @Router /api/users [get]
        users.MapGet("", listUsers)
        
        // @Summary 获取用户详情
        // @Tags 用户管理
        // @Param id path int true "用户 ID"
        // @Produce json
        // @Success 200 {object} User
        // @Failure 404 {object} ApiError
        // @Router /api/users/{id} [get]
        users.MapGet("/:id", getUser)
        
        // @Summary 创建用户
        // @Tags 用户管理
        // @Accept json
        // @Produce json
        // @Param request body CreateUserRequest true "用户信息"
        // @Success 201 {object} User
        // @Failure 400 {object} ApiError
        // @Router /api/users [post]
        users.MapPost("", createUser)
        
        // @Summary 更新用户
        // @Tags 用户管理
        // @Accept json
        // @Produce json
        // @Param id path int true "用户 ID"
        // @Param request body CreateUserRequest true "用户信息"
        // @Success 200 {object} User
        // @Failure 404 {object} ApiError
        // @Router /api/users/{id} [put]
        users.MapPut("/:id", updateUser)
        
        // @Summary 删除用户
        // @Tags 用户管理
        // @Param id path int true "用户 ID"
        // @Success 204
        // @Failure 404 {object} ApiError
        // @Router /api/users/{id} [delete]
        users.MapDelete("/:id", deleteUser)
    }
    
    // 启用 Swagger
    swagger.UseSwagger(app)
    swagger.UseSwaggerUI(app)
    
    app.Run()
}

func listUsers(c *web.HttpContext) web.IActionResult {
    users := []User{
        {ID: 1, Name: "张三", Email: "zhangsan@example.com"},
        {ID: 2, Name: "李四", Email: "lisi@example.com"},
    }
    return c.Ok(users)
}

func getUser(c *web.HttpContext) web.IActionResult {
    user := User{ID: 1, Name: "张三", Email: "zhangsan@example.com"}
    return c.Ok(user)
}

func createUser(c *web.HttpContext) web.IActionResult {
    req, err := web.BindAndValidate[CreateUserRequest](c)
    if err != nil {
        return err
    }
    
    user := User{ID: 1, Name: req.Name, Email: req.Email}
    return c.Created(user)
}

func updateUser(c *web.HttpContext) web.IActionResult {
    req, err := web.BindAndValidate[CreateUserRequest](c)
    if err != nil {
        return err
    }
    
    user := User{ID: 1, Name: req.Name, Email: req.Email}
    return c.Ok(user)
}

func deleteUser(c *web.HttpContext) web.IActionResult {
    return c.NoContent()
}
```

访问 http://localhost:8080/swagger 查看完整的 API 文档。

## 环境配置

### 开发环境启用

```go
builder := web.CreateBuilder()

// 只在开发环境启用 Swagger
if builder.Environment.IsDevelopment() {
    swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
        opts.Title = "开发环境 API"
        opts.Version = "dev"
    })
}

app := builder.Build()

if builder.Environment.IsDevelopment() {
    swagger.UseSwagger(app)
    swagger.UseSwaggerUI(app)
}

app.Run()
```

### 生产环境保护

```go
// 生产环境可以启用 Swagger JSON，但关闭 UI
if builder.Environment.IsProduction() {
    swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
        opts.Title = "生产环境 API"
        opts.Version = "v1"
    })
}

app := builder.Build()

// 只启用 JSON 端点，不启用 UI
if builder.Environment.IsProduction() {
    swagger.UseSwagger(app)
    // 不调用 UseSwaggerUI
}

app.Run()
```

## 最佳实践

### 1. 使用注释标记路由

```go
// ✅ 推荐：添加详细注释
// @Summary 创建用户
// @Description 创建一个新用户账号
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "用户信息"
// @Success 201 {object} User
// @Failure 400 {object} ApiError
// @Router /api/users [post]
app.MapPost("/api/users", createUser)

// ❌ 不推荐：没有注释
app.MapPost("/api/users", createUser)
```

### 2. 组织 API 标签

```go
// ✅ 推荐：使用路由组自动生成标签
users := api.MapGroup("/users")     // 标签：users
orders := api.MapGroup("/orders")   // 标签：orders

// ✅ 或自定义标签
app.MapGet("/api/users", getUsers).WithTags("用户管理")
```

### 3. 定义清晰的模型

```go
// ✅ 推荐：使用 example 标签
type User struct {
    ID    int    `json:"id" example:"1"`
    Name  string `json:"name" example:"张三"`
    Email string `json:"email" example:"zhangsan@example.com"`
    Age   int    `json:"age" example:"25"`
}

// ❌ 不推荐：没有示例
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

### 4. 安全方案

```go
// ✅ 推荐：添加安全定义
swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
    opts.AddSecurityDefinition("Bearer", openapi.SecurityScheme{
        Type:         "http",
        Scheme:       "bearer",
        BearerFormat: "JWT",
    })
})
```

### 5. 环境隔离

```go
// ✅ 推荐：根据环境配置
if builder.Environment.IsDevelopment() {
    swagger.UseSwagger(app)
    swagger.UseSwaggerUI(app)
}
```

## API 参考

### AddSwaggerGen

```go
AddSwaggerGen(
    services di.IServiceCollection, 
    configure func(*SwaggerGenOptions),
) di.IServiceCollection
```

### UseSwagger

```go
UseSwagger(app *web.WebApplication)
```

启用 Swagger JSON 端点：`/swagger/v1/swagger.json`

### UseSwaggerUI

```go
UseSwaggerUI(
    app *web.WebApplication,
    configure func(*SwaggerUIOptions),
)
```

启用 Swagger UI 界面，默认路径：`/swagger`

### SwaggerGenOptions

```go
type SwaggerGenOptions struct {
    Title               string
    Version             string
    Description         string
    SecurityDefinitions map[string]openapi.SecurityScheme
}
```

### SwaggerUIOptions

```go
type SwaggerUIOptions struct {
    RoutePrefix string  // UI 路由前缀
    SpecURL     string  // Swagger JSON 路径
    Title       string  // 页面标题
}
```

## 常见问题

### 如何查看生成的 OpenAPI 规范？

访问 `/swagger/v1/swagger.json` 查看 JSON 格式的 OpenAPI 规范。

### Swagger UI 不显示路由？

确保在调用 `UseSwagger` 和 `UseSwaggerUI` 之前定义了所有路由：

```go
// ✅ 正确顺序
app := builder.Build()
app.MapGet("/api/users", getUsers)  // 先定义路由
swagger.UseSwagger(app)              // 后启用 Swagger
swagger.UseSwaggerUI(app)

// ❌ 错误顺序
app := builder.Build()
swagger.UseSwagger(app)
swagger.UseSwaggerUI(app)
app.MapGet("/api/users", getUsers)  // 路由不会被发现
```

### 如何自定义 Swagger UI 路径？

```go
swagger.UseSwaggerUI(app, func(opts *swagger.SwaggerUIOptions) {
    opts.RoutePrefix = "/docs"  // 自定义路径
})
// 访问：http://localhost:8080/docs
```

### 如何添加多个安全方案？

```go
swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
    // Bearer Token
    opts.AddSecurityDefinition("Bearer", openapi.SecurityScheme{
        Type:   "http",
        Scheme: "bearer",
    })
    
    // API Key
    opts.AddSecurityDefinition("ApiKey", openapi.SecurityScheme{
        Type: "apiKey",
        In:   "header",
        Name: "X-API-Key",
    })
})
```

---

[← 返回主目录](../README.md)

