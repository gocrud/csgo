# 错误处理系统

CSGO 框架提供了一套完整的错误处理体系，包括框架级错误码和业务错误码构建器。

## 目录

- [框架级错误码](#框架级错误码)
- [业务错误码构建器](#业务错误码构建器)
- [Web 层集成](#web-层集成)
- [验证错误](#验证错误)
- [最佳实践](#最佳实践)

## 框架级错误码

框架预定义了常用的错误码，所有验证规则会自动使用这些错误码。

### 错误码格式

错误码采用 **模块.语义描述** 的格式，使用点分隔，全大写下划线命名。

例如：`VALIDATION.REQUIRED`、`SYSTEM.INTERNAL_ERROR`

### 错误码列表

#### 系统错误 (SYSTEM.*)

| 错误码 | 描述 | HTTP 状态码 |
|--------|------|------------|
| `SYSTEM.INTERNAL_ERROR` | 系统内部错误 | 500 |
| `SYSTEM.SERVICE_UNAVAILABLE` | 服务不可用 | 503 |
| `SYSTEM.TIMEOUT` | 系统超时 | 504 |

#### 验证错误 (VALIDATION.*)

| 错误码 | 描述 | 使用场景 |
|--------|------|---------|
| `VALIDATION.FAILED` | 通用验证失败 | 自定义验证失败 |
| `VALIDATION.REQUIRED` | 必填项为空 | NotEmpty() |
| `VALIDATION.MIN_LENGTH` | 字符串长度小于最小值 | MinLength() |
| `VALIDATION.MAX_LENGTH` | 字符串长度大于最大值 | MaxLength() |
| `VALIDATION.LENGTH` | 字符串长度不在范围内 | Length() |
| `VALIDATION.MIN` | 数值小于最小值 | GreaterThan(), GreaterThanOrEqual() |
| `VALIDATION.MAX` | 数值大于最大值 | LessThan(), LessThanOrEqual() |
| `VALIDATION.RANGE` | 数值不在范围内 | InclusiveBetween(), ExclusiveBetween() |
| `VALIDATION.EMAIL` | 邮箱格式不正确 | EmailAddress() |
| `VALIDATION.URL` | URL 格式不正确 | - |
| `VALIDATION.PATTERN` | 不匹配正则表达式 | Matches() |
| `VALIDATION.IN` | 值不在枚举列表中 | - |
| `VALIDATION.NOT_IN` | 值在排除列表中 | - |
| `VALIDATION.NOT_EMPTY` | 集合不能为空 | NotEmptySlice() |
| `VALIDATION.MIN_COUNT` | 集合元素数量不足 | MinLengthSlice() |
| `VALIDATION.MAX_COUNT` | 集合元素数量过多 | MaxLengthSlice() |

#### HTTP 错误 (HTTP.*)

| 错误码 | 描述 | HTTP 状态码 |
|--------|------|------------|
| `HTTP.BAD_REQUEST` | 错误的请求 | 400 |
| `HTTP.UNAUTHORIZED` | 未授权 | 401 |
| `HTTP.FORBIDDEN` | 禁止访问 | 403 |
| `HTTP.NOT_FOUND` | 资源不存在 | 404 |
| `HTTP.CONFLICT` | 资源冲突 | 409 |
| `HTTP.METHOD_NOT_ALLOWED` | 方法不允许 | 405 |

#### 认证授权 (AUTH.*)

| 错误码 | 描述 | HTTP 状态码 |
|--------|------|------------|
| `AUTH.TOKEN_EXPIRED` | 令牌已过期 | 401 |
| `AUTH.TOKEN_INVALID` | 令牌无效 | 401 |
| `AUTH.PERMISSION_DENIED` | 权限不足 | 403 |
| `AUTH.CREDENTIALS_INVALID` | 凭证无效 | 401 |

## 业务错误码构建器

框架提供了便捷的错误码构建器，让您轻松创建业务错误而无需手动定义大量常量。

### 基本用法

```go
import "github.com/gocrud/csgo/errors"

// 创建业务错误
err := errors.Business("USER").NotFound("用户不存在")
// 自动生成: Code="USER.NOT_FOUND", Message="用户不存在"

err := errors.Business("ORDER").InvalidStatus("订单已关闭")
// 自动生成: Code="ORDER.INVALID_STATUS", Message="订单已关闭"
```

### 内置语义方法

| 方法 | 生成的错误码 | 典型 HTTP 状态码 | 使用场景 |
|------|------------|-----------------|---------|
| `NotFound(message)` | `模块.NOT_FOUND` | 404 | 资源不存在 |
| `AlreadyExists(message)` | `模块.ALREADY_EXISTS` | 409 | 资源已存在 |
| `InvalidStatus(message)` | `模块.INVALID_STATUS` | 400 | 状态无效 |
| `InvalidParam(message)` | `模块.INVALID_PARAM` | 400 | 参数无效 |
| `PermissionDenied(message)` | `模块.PERMISSION_DENIED` | 403 | 权限不足 |
| `OperationFailed(message)` | `模块.OPERATION_FAILED` | 400 | 操作失败 |
| `Expired(message)` | `模块.EXPIRED` | 410 | 资源已过期 |
| `Locked(message)` | `模块.LOCKED` | 423 | 资源已锁定 |
| `LimitExceeded(message)` | `模块.LIMIT_EXCEEDED` | 429 | 超出限制 |
| `Custom(semantic, message)` | `模块.自定义语义` | 400 | 自定义错误 |

### 完整示例

```go
package services

import (
    "github.com/gocrud/csgo/errors"
    "github.com/gocrud/csgo/web"
)

type UserService struct {
    repo UserRepository
}

func (s *UserService) GetUser(id int) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    if user == nil {
        // 使用构建器创建业务错误
        return nil, errors.Business("USER").NotFound("用户不存在")
    }
    
    return user, nil
}

func (s *UserService) CreateUser(req *CreateUserRequest) (*User, error) {
    // 检查用户是否已存在
    existing, _ := s.repo.FindByEmail(req.Email)
    if existing != nil {
        return nil, errors.Business("USER").AlreadyExists("邮箱已被注册")
    }
    
    // 创建用户
    user := &User{
        Email:    req.Email,
        Password: hashPassword(req.Password),
    }
    
    if err := s.repo.Create(user); err != nil {
        return nil, errors.Business("USER").OperationFailed("创建用户失败")
    }
    
    return user, nil
}

func (s *UserService) UpdateUserStatus(id int, status string) error {
    user, err := s.GetUser(id)
    if err != nil {
        return err
    }
    
    // 检查状态转换是否合法
    if !user.CanTransitionTo(status) {
        return errors.Business("USER").InvalidStatus(
            fmt.Sprintf("不能从 %s 转换到 %s", user.Status, status))
    }
    
    user.Status = status
    return s.repo.Update(user)
}

// 自定义语义错误
func (s *UserService) SendVerificationEmail(userID int) error {
    user, err := s.GetUser(userID)
    if err != nil {
        return err
    }
    
    // 检查发送限制
    if user.EmailSentToday >= 5 {
        return errors.Business("USER").Custom("EMAIL_LIMIT_EXCEEDED", 
            "今日验证邮件发送次数已达上限")
        // 生成: Code="USER.EMAIL_LIMIT_EXCEEDED"
    }
    
    // 发送邮件...
    return nil
}
```

## Web 层集成

### 推荐方式：使用 FromError（简洁优雅）

框架提供了 `FromError` 方法，可以智能识别错误类型并自动返回对应的 HTTP 响应：

```go
package controllers

import (
    "github.com/gocrud/csgo/errors"
    "github.com/gocrud/csgo/web"
)

type UserController struct {
    userService *UserService
}

// 推荐：使用 FromError 一行搞定所有错误处理
func (ctrl *UserController) GetUser(c *web.HttpContext) web.IActionResult {
    id := c.Params().PathInt("id").Value()
    
    user, err := ctrl.userService.GetUser(id)
    if err != nil {
        // FromError 会自动识别错误类型：
        // - BizError: 自动映射 HTTP 状态码（NOT_FOUND -> 404）
        // - ValidationErrors: 返回 400 验证错误
        // - 普通 error: 返回 500 + 自定义消息
        return c.FromError(err, "获取用户失败")
    }
    
    return c.Ok(user)
}

func (ctrl *UserController) CreateUser(c *web.HttpContext) web.IActionResult {
    req, result := web.BindAndValidate[CreateUserRequest](c)
    if result != nil {
        return result
    }
    
    user, err := ctrl.userService.CreateUser(req)
    if err != nil {
        return c.FromError(err, "创建用户失败")  // 一行搞定！
    }
    
    return c.Created(user)
}

// 指定状态码（用于特殊场景）
func (ctrl *UserController) ConnectDatabase(c *web.HttpContext) web.IActionResult {
    err := ctrl.service.Connect()
    if err != nil {
        // 数据库连接错误返回 503
        return c.FromErrorWithStatus(err, 503, "数据库服务暂时不可用")
    }
    return c.Ok(nil)
}
```

### 传统方式：手动类型判断（仍然支持）

如果需要更细粒度的控制，仍然可以使用传统方式：

```go
func (ctrl *UserController) GetUser(c *web.HttpContext) web.IActionResult {
    user, err := ctrl.userService.GetUser(id)
    if err != nil {
        // 手动判断错误类型
        if bizErr, ok := err.(*errors.BizError); ok {
            return c.BizError(bizErr)  // 自动映射状态码
        }
        return c.InternalError("服务器错误")
    }
    return c.Ok(user)
}
```

### 自定义错误处理器（高级功能）

为特定错误类型注册自定义处理逻辑：

```go
package main

import (
    "database/sql"
    "errors"
    "github.com/gocrud/csgo/web"
)

func init() {
    // 注册数据库记录不存在错误处理器
    web.RegisterErrorHandler(
        func(err error) bool {
            return errors.Is(err, sql.ErrNoRows)
        },
        func(err error, msg ...string) web.IActionResult {
            message := "记录不存在"
            if len(msg) > 0 && msg[0] != "" {
                message = msg[0]
            }
            return web.Error(404, "NOT_FOUND", message)
        },
    )
    
    // 注册超时错误处理器
    web.RegisterErrorHandler(
        func(err error) bool {
            return errors.Is(err, context.DeadlineExceeded)
        },
        func(err error, msg ...string) web.IActionResult {
            return web.Error(408, "TIMEOUT", "请求超时")
        },
    )
}

// 控制器中使用，自动应用自定义处理器
func (ctrl *UserController) GetUser(c *web.HttpContext) web.IActionResult {
    user, err := ctrl.repo.FindByID(id) // 可能返回 sql.ErrNoRows
    if err != nil {
        return c.FromError(err, "用户不存在") // 自动使用注册的处理器
    }
    return c.Ok(user)
}
```

### 错误响应格式

#### 验证错误响应

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION.FAILED",
    "message": "验证失败",
    "fields": [
      {
        "field": "email",
        "message": "邮箱格式不正确",
        "code": "VALIDATION.EMAIL"
      },
      {
        "field": "password",
        "message": "长度不能少于 6",
        "code": "VALIDATION.MIN_LENGTH"
      }
    ]
  }
}
```

#### 业务错误响应

```json
{
  "success": false,
  "error": {
    "code": "USER.NOT_FOUND",
    "message": "用户不存在"
  }
}
```

## 验证错误

验证错误会自动使用框架定义的错误码。

### 验证器示例

```go
package validators

import (
    "github.com/gocrud/csgo/validation"
)

type CreateUserRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Age      int    `json:"age"`
}

func NewCreateUserValidator() *validation.AbstractValidator[CreateUserRequest] {
    v := validation.NewValidator[CreateUserRequest]()  // 快速失败模式
    
    v.Field(func(r *CreateUserRequest) string { return r.Email }).
        NotEmpty().           // 自动使用 VALIDATION.REQUIRED
        EmailAddress()        // 自动使用 VALIDATION.EMAIL
    
    v.Field(func(r *CreateUserRequest) string { return r.Password }).
        NotEmpty().           // 自动使用 VALIDATION.REQUIRED
        MinLength(6)          // 自动使用 VALIDATION.MIN_LENGTH
    
    v.FieldInt(func(r *CreateUserRequest) int { return r.Age }).
        GreaterThanOrEqual(18) // 自动使用 VALIDATION.MIN
    
    return v
}

// 注册验证器
func init() {
    validation.RegisterValidator[CreateUserRequest](NewCreateUserValidator())
}
```

### 验证模式

框架支持两种验证模式：

#### 快速失败模式（默认，推荐）

```go
// 创建快速失败验证器
v := validation.NewValidator[User]()
// 遇到第一个错误立即返回，性能最优
// 适合 99% 的表单验证场景
```

#### 全量验证模式

```go
// 创建全量验证器
v := validation.NewValidatorAll[User]()
// 收集所有字段的所有错误
// 适合批量数据导入、复杂表单审核等场景
```

**注意：** Web 层无需关心验证模式，注册什么模式就使用什么模式。

## 最佳实践

### 1. 错误码命名规范

- **模块名**：使用业务领域名称，如 `USER`、`ORDER`、`PAYMENT`
- **语义描述**：使用清晰的动词或状态描述，如 `NOT_FOUND`、`INVALID_STATUS`
- **全大写下划线**：统一使用大写字母和下划线
- **避免重复**：不要在错误码中重复模块名，如 ~~`USER.USER_NOT_FOUND`~~ 应为 `USER.NOT_FOUND`

### 2. 优先使用构建器

```go
// ✅ 推荐：使用构建器
err := errors.Business("USER").NotFound("用户不存在")

// ❌ 不推荐：手动定义常量（除非有特殊需求）
const UserNotFound = "USER.NOT_FOUND"
err := &errors.BizError{Code: UserNotFound, Message: "用户不存在"}
```

### 3. 合理划分错误粒度

```go
// ✅ 推荐：合理的粒度
errors.Business("USER").NotFound("用户不存在")
errors.Business("ORDER").NotFound("订单不存在")

// ❌ 不推荐：过细的粒度（维护成本高）
errors.Business("USER").Custom("NOT_FOUND_BY_ID", "通过ID未找到用户")
errors.Business("USER").Custom("NOT_FOUND_BY_EMAIL", "通过邮箱未找到用户")
```

### 4. 错误消息本地化

错误消息应该面向用户，考虑国际化需求：

```go
// ✅ 推荐：用户友好的消息
errors.Business("USER").NotFound("用户不存在")

// ❌ 不推荐：技术性消息
errors.Business("USER").NotFound("User record not found in database")
```

### 5. 验证器注册

```go
// 在 validators 包中定义和注册
package validators

func init() {
    // 快速失败模式（默认）- 适合大多数场景
    validation.RegisterValidator[CreateUserRequest](NewCreateUserValidator())
    
    // 全量验证模式 - 适合批量导入
    validation.RegisterValidator[BatchImportRequest](NewBatchImportValidator())
}
```

### 6. 统一错误处理

```go
// 在 Controller 中统一处理错误
func (ctrl *BaseController) HandleServiceError(ctx *web.HttpContext, err error) web.IActionResult {
    if bizErr, ok := err.(*errors.BizError); ok {
        return ctx.BizError(bizErr)
    }
    
    // 记录未预期的错误
    log.Error().Err(err).Msg("Unexpected error")
    return ctx.InternalError("服务器错误")
}
```

## 错误码与 HTTP 状态码映射

框架会自动将业务错误码映射到合适的 HTTP 状态码：

| 错误码模式 | HTTP 状态码 | 说明 |
|-----------|------------|------|
| `*.NOT_FOUND` | 404 | 资源不存在 |
| `*.ALREADY_EXISTS` | 409 | 资源冲突 |
| `*.PERMISSION_DENIED` | 403 | 权限不足 |
| `*.UNAUTHORIZED` | 401 | 未授权 |
| `*.INVALID*` | 400 | 参数或状态无效 |
| `*.EXPIRED` | 410 | 资源已过期 |
| `*.LOCKED` | 423 | 资源已锁定 |
| `*.LIMIT_EXCEEDED` | 429 | 超出限制 |
| 其他 | 400 | 默认错误请求 |

如果自动映射不满足需求，可以使用 `BizErrorWithStatus` 手动指定状态码。

## 总结

CSGO 的错误处理系统设计理念：

1. **框架负责框架的事**：框架预定义验证、系统、HTTP 等框架级错误码
2. **业务负责业务的事**：业务通过构建器灵活创建业务错误码
3. **提升开发体验**：减少样板代码，专注业务逻辑
4. **类型安全**：充分利用 Go 的类型系统，编译时发现问题
5. **统一响应格式**：前后端对接更加规范

通过这套体系，您可以快速构建健壮、规范的 API 错误处理。
