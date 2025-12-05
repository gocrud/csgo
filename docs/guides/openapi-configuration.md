# OpenAPI 配置指南

本指南介绍如何在 CSGO 框架中配置 OpenAPI 文档。

## 目录

- [快速开始](#快速开始)
- [Doc Tag 语法](#doc-tag-语法)
- [函数式配置](#函数式配置)
- [完整示例](#完整示例)
- [最佳实践](#最佳实践)

## 快速开始

### 1. 启用 Swagger

```go
package main

import (
    "github.com/gocrud/csgo/web"
    "github.com/gocrud/csgo/swagger"
)

func main() {
    builder := web.CreateBuilder()
    
    // 配置 Swagger
    builder.Services.AddSwaggerGen(func(opts *swagger.SwaggerGenOptions) {
        opts.Title = "我的 API"
        opts.Version = "v1.0"
        opts.Description = "API 文档"
    })
    
    app := builder.Build()
    
    // 启用 Swagger
    swagger.UseSwagger(app)
    swagger.UseSwaggerUI(app)
    
    app.Run()
}
```

### 2. 配置端点

```go
import "github.com/gocrud/csgo/openapi"

app.MapGet("/users/:id", getUserHandler).
    WithOpenApi(
        openapi.Name("GetUser"),
        openapi.Summary("获取用户详情"),
        openapi.Tags("Users"),
        openapi.Produces[User](200),
        openapi.ProducesProblem(404),
    )
```

### 3. 访问文档

- Swagger UI: http://localhost:8080/swagger
- OpenAPI JSON: http://localhost:8080/swagger/v1/swagger.json

## Doc Tag 语法

CSGO 支持自定义 `doc` tag 来定义 OpenAPI Schema。

### 基本格式

```go
doc:"描述,配置项,配置项:值"
```

- 第一部分是字段描述
- 后续部分是配置项，用逗号分隔
- 配置项可以是标志（如 `required`）或键值对（如 `example:value`）

### 支持的配置项

| 配置项 | 类型 | 说明 | 示例 |
|--------|------|------|------|
| `required` | 标志 | 标记为必需字段 | `doc:"用户名,required"` |
| `example:value` | 键值 | 示例值 | `doc:"用户名,example:张三"` |
| `format:type` | 键值 | 格式类型 | `doc:"邮箱,format:email"` |
| `min:n` | 键值 | 最小值（数字） | `doc:"年龄,min:0"` |
| `max:n` | 键值 | 最大值（数字） | `doc:"年龄,max:120"` |
| `minLength:n` | 键值 | 最小长度（字符串） | `doc:"用户名,minLength:2"` |
| `maxLength:n` | 键值 | 最大长度（字符串） | `doc:"用户名,maxLength:50"` |
| `pattern:regex` | 键值 | 正则模式 | `doc:"手机,pattern:^1[3-9]\\d{9}$"` |
| `enum:a\|b\|c` | 键值 | 枚举值（用 \| 分隔） | `doc:"状态,enum:active\|inactive\|banned"` |

### 常用格式

| 格式 | 说明 | 示例 |
|------|------|------|
| `email` | 电子邮件 | `doc:"邮箱,format:email"` |
| `date` | 日期 (YYYY-MM-DD) | `doc:"生日,format:date"` |
| `date-time` | 日期时间 (ISO 8601) | `doc:"创建时间,format:date-time"` |
| `uri` | URI | `doc:"主页,format:uri"` |
| `uuid` | UUID | `doc:"ID,format:uuid"` |
| `binary` | 二进制数据 | `doc:"文件,format:binary"` |
| `password` | 密码（不显示） | `doc:"密码,format:password"` |

## 函数式配置

### 配置选项函数

所有配置通过 `WithOpenApi()` 方法集中管理：

```go
app.MapGet("/users", handler).
    WithOpenApi(
        openapi.Name("ListUsers"),              // 端点名称
        openapi.Summary("获取用户列表"),         // 摘要
        openapi.Description("返回所有用户"),     // 详细描述
        openapi.Tags("Users"),                  // 标签
        openapi.Produces[[]User](200),          // 响应类型
        openapi.ProducesProblem(400),           // 错误响应
        openapi.Authorization("Admin"),         // 授权策略
    )
```

### 可用的配置选项

| 选项 | 说明 | 示例 |
|------|------|------|
| `Name(string)` | 设置端点名称 | `openapi.Name("GetUser")` |
| `Summary(string)` | 设置摘要 | `openapi.Summary("获取用户")` |
| `Description(string)` | 设置详细描述 | `openapi.Description("根据 ID 获取用户")` |
| `Tags(...string)` | 添加标签 | `openapi.Tags("Users", "Admin")` |
| `Produces[T](code)` | 添加响应类型 | `openapi.Produces[User](200)` |
| `ProducesProblem(code)` | 添加问题响应 | `openapi.ProducesProblem(404)` |
| `ProducesValidationProblem()` | 添加验证错误响应 (422) | `openapi.ProducesValidationProblem()` |
| `Accepts[T](contentType)` | 添加请求体类型 | `openapi.Accepts[CreateUserRequest]("application/json")` |
| `Authorization(...policies)` | 要求授权 | `openapi.Authorization("Admin")` |
| `Anonymous()` | 允许匿名访问 | `openapi.Anonymous()` |

## 完整示例

### 定义模型

```go
// User 用户模型
type User struct {
    ID        int    `json:"id" doc:"用户ID,example:1"`
    Name      string `json:"name" doc:"用户名,required,example:张三,minLength:2,maxLength:50"`
    Email     string `json:"email" doc:"电子邮件,required,format:email,example:user@example.com"`
    Age       int    `json:"age,omitempty" doc:"年龄,min:0,max:120,example:25"`
    Status    string `json:"status" doc:"用户状态,required,enum:active|inactive|banned,example:active"`
    Phone     string `json:"phone,omitempty" doc:"手机号,pattern:^1[3-9]\\d{9}$,example:13800138000"`
    CreatedAt string `json:"createdAt" doc:"创建时间,format:date-time"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
    Name  string `json:"name" doc:"用户名,required,example:张三,minLength:2,maxLength:50"`
    Email string `json:"email" doc:"电子邮件,required,format:email,example:user@example.com"`
    Age   int    `json:"age,omitempty" doc:"年龄,min:0,max:120"`
}
```

### 配置路由

#### GET 请求

```go
app.MapGet("/users/:id", getUserHandler).
    WithOpenApi(
        openapi.Name("GetUser"),
        openapi.Summary("获取用户详情"),
        openapi.Description("根据用户ID获取用户的详细信息"),
        openapi.Tags("Users"),
        openapi.Produces[User](200),
        openapi.ProducesProblem(404),
    )
```

#### POST 请求

```go
app.MapPost("/users", createUserHandler).
    WithOpenApi(
        openapi.Name("CreateUser"),
        openapi.Summary("创建新用户"),
        openapi.Tags("Users"),
        openapi.Accepts[CreateUserRequest]("application/json"),
        openapi.Produces[User](201),
        openapi.ProducesValidationProblem(),
        openapi.ProducesProblem(400),
    )
```

#### PUT 请求

```go
app.MapPut("/users/:id", updateUserHandler).
    WithOpenApi(
        openapi.Name("UpdateUser"),
        openapi.Summary("更新用户信息"),
        openapi.Tags("Users"),
        openapi.Accepts[UpdateUserRequest]("application/json"),
        openapi.Produces[User](200),
        openapi.ProducesProblem(404),
        openapi.ProducesValidationProblem(),
    )
```

#### DELETE 请求

```go
app.MapDelete("/users/:id", deleteUserHandler).
    WithOpenApi(
        openapi.Name("DeleteUser"),
        openapi.Summary("删除用户"),
        openapi.Tags("Users"),
        openapi.Produces[any](204),
        openapi.ProducesProblem(404),
    )
```

### 路由组配置

```go
// 创建 API 组
apiV1 := app.MapGroup("/api/v1").
    WithOpenApi(
        openapi.Tags("API v1"),
    )

// 组内路由会自动继承标签
apiV1.MapGet("/users", listUsersHandler).
    WithOpenApi(
        openapi.Name("V1.ListUsers"),
        openapi.Summary("获取用户列表"),
        openapi.Produces[[]User](200),
    )

apiV1.MapPost("/users", createUserHandler).
    WithOpenApi(
        openapi.Name("V1.CreateUser"),
        openapi.Summary("创建用户"),
        openapi.Accepts[CreateUserRequest]("application/json"),
        openapi.Produces[User](201),
        openapi.ProducesValidationProblem(),
    )
```

### 多种响应状态码

```go
app.MapPost("/users/:id/avatar", uploadAvatarHandler).
    WithOpenApi(
        openapi.Name("UploadAvatar"),
        openapi.Summary("上传用户头像"),
        openapi.Tags("Users"),
        openapi.Accepts[any]("multipart/form-data"),
        openapi.Produces[User](200),          // 上传成功
        openapi.Produces[any](202),           // 接受但处理中
        openapi.ProducesProblem(400),         // 文件格式错误
        openapi.ProducesProblem(413),         // 文件太大
    )
```

## 最佳实践

### 1. 显式启用 OpenAPI

只为需要公开的端点启用 OpenAPI：

```go
// ✅ 公共 API - 启用 OpenAPI
app.MapGet("/api/users", publicUsersHandler).
    WithOpenApi(
        openapi.Name("ListUsers"),
        openapi.Summary("获取用户列表"),
        openapi.Produces[[]User](200),
    )

// ✅ 内部 API - 不启用 OpenAPI（不会出现在文档中）
app.MapGet("/internal/debug", debugHandler)
```

### 2. 合理使用 Doc Tag

将常用的验证规则和描述放在 `doc` tag 中：

```go
type CreateUserRequest struct {
    // ✅ 好的做法：完整的描述和验证规则
    Name  string `json:"name" doc:"用户名，2-50个字符,required,minLength:2,maxLength:50,example:张三"`
    Email string `json:"email" doc:"电子邮件地址,required,format:email,example:user@example.com"`
    
    // ❌ 避免：信息不完整
    Age   int    `json:"age" doc:"年龄"`
    
    // ✅ 更好：添加约束和示例
    Age   int    `json:"age,omitempty" doc:"年龄（0-120岁）,min:0,max:120,example:25"`
}
```

### 3. 使用有意义的名称

```go
// ✅ 好的做法：清晰的名称
app.MapGet("/users", handler).
    WithOpenApi(
        openapi.Name("ListUsers"),
        openapi.Summary("获取用户列表"),
        ...
    )

// ❌ 避免：模糊的名称
app.MapGet("/users", handler).
    WithOpenApi(
        openapi.Name("GetUsers"),  // 不清楚是获取列表还是单个
        ...
    )
```

### 4. 组织标签

使用标签组织相关的端点：

```go
// 用户相关
app.MapGet("/users", ...).WithOpenApi(...openapi.Tags("Users"))
app.MapPost("/users", ...).WithOpenApi(...openapi.Tags("Users"))

// 订单相关
app.MapGet("/orders", ...).WithOpenApi(...openapi.Tags("Orders"))
app.MapPost("/orders", ...).WithOpenApi(...openapi.Tags("Orders"))

// 管理员功能
app.MapGet("/admin/users", ...).WithOpenApi(...openapi.Tags("Admin", "Users"))
```

### 5. 提供完整的响应类型

为所有可能的响应状态码提供类型信息：

```go
app.MapPost("/users", handler).
    WithOpenApi(
        openapi.Produces[User](201),              // 成功创建
        openapi.ProducesProblem(400),             // 请求错误
        openapi.ProducesValidationProblem(),      // 验证错误 (422)
        openapi.ProducesProblem(409),             // 冲突（如邮箱已存在）
        openapi.ProducesProblem(500),             // 服务器错误
    )
```

### 6. 使用枚举约束

对于有限的可选值，使用枚举：

```go
type User struct {
    Status string `json:"status" doc:"用户状态,required,enum:active|inactive|banned,example:active"`
    Role   string `json:"role" doc:"用户角色,required,enum:admin|user|guest,example:user"`
}
```

### 7. 国际化考虑

如果需要支持多语言，考虑在代码中管理描述：

```go
// 方式1：在 tag 中使用英文
type User struct {
    Name string `json:"name" doc:"User name,required,minLength:2,maxLength:50"`
}

// 方式2：使用配置选项（更灵活）
app.MapGet("/users", handler).
    WithOpenApi(
        openapi.Name("GetUser"),
        openapi.Summary(i18n.GetText("api.users.get.summary")),
        openapi.Description(i18n.GetText("api.users.get.description")),
        ...
    )
```

## 与 .NET 的对比

CSGO 的 OpenAPI 配置与 .NET Minimal API 高度一致：

### Go (CSGO)

```go
app.MapPost("/users", createUserHandler).
    WithOpenApi(
        openapi.Name("CreateUser"),
        openapi.Summary("创建新用户"),
        openapi.Tags("Users"),
        openapi.Accepts[CreateUserRequest]("application/json"),
        openapi.Produces[User](201),
        openapi.ProducesValidationProblem(),
    )
```

### C# (.NET)

```csharp
app.MapPost("/users", createUserHandler)
    .WithName("CreateUser")
    .WithSummary("创建新用户")
    .WithTags("Users")
    .Accepts<CreateUserRequest>("application/json")
    .Produces<User>(201)
    .ProducesValidationProblem()
    .WithOpenApi();
```

### 主要区别

| 特性 | Go (CSGO) | C# (.NET) |
|------|-----------|-----------|
| 配置位置 | 所有配置作为参数传给 `WithOpenApi()` | 链式调用，`WithOpenApi()` 在最后 |
| 泛型语法 | `[T]` | `<T>` |
| 配置方式 | 函数式选项 `openapi.Name("...")` | 直接方法 `.WithName("...")` |
| 类型推断 | 需要显式指定泛型参数 | 某些情况可以推断 |

### 为什么 Go 采用不同的设计？

Go 的方法不支持泛型，因此 `Produces[T]` 必须是函数而非方法。为了保持 API 一致性，我们将所有配置统一为函数式选项模式，集中在 `WithOpenApi()` 中配置。

## 参考资料

- [OpenAPI 规范](https://spec.openapis.org/oas/v3.0.3)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [.NET Minimal APIs](https://learn.microsoft.com/en-us/aspnet/core/fundamentals/minimal-apis)

