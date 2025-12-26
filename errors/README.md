# 错误处理系统

CSGO 框架提供了简洁、类型安全、功能强大的错误处理体系。

## 🎯 核心特性

- ✅ **简洁易用** - 预定义常用错误方法，一行代码搞定
- ✅ **类型安全** - 模块化设计，避免拼写错误
- ✅ **链式调用** - 流畅的 API，支持灵活组合
- ✅ **不可变** - 所有操作返回新实例，线程安全
- ✅ **错误链** - 完整支持 Go 1.13+ 错误包装
- ✅ **HTTP 集成** - 自动映射合适的 HTTP 状态码

## 📖 目录

- [快速开始](#快速开始)
- [Module 快捷方法](#module-快捷方法)
- [自定义错误码](#自定义错误码)
- [链式调用](#链式调用)
- [Web 层集成](#web-层集成)
- [最佳实践](#最佳实践)
- [迁移指南](#迁移指南)

## 🚀 快速开始

### 1. 定义模块

```go
package services

import "github.com/gocrud/csgo/errors"

// 定义各模块的错误
var (
    UserErrors  = errors.NewModule("USER")
    OrderErrors = errors.NewModule("ORDER")
    DramaErrors = errors.NewModule("DRAMA")
)
```

### 2. 使用快捷方法（最常用）

```go
// 资源不存在
err := UserErrors.NotFound("用户不存在")
// 生成: code="USER.NOT_FOUND", httpCode=404

// 参数无效
err := OrderErrors.InvalidParam("订单金额必须大于0")
// 生成: code="ORDER.INVALID_PARAM", httpCode=400

// 权限不足
err := DramaErrors.PermissionDenied("无权访问此剧集")
// 生成: code="DRAMA.PERMISSION_DENIED", httpCode=403
```

### 3. 自定义错误码

```go
// 支付失败
err := OrderErrors.Code("PAYMENT_FAILED").Msg("支付失败")
// 生成: code="ORDER.PAYMENT_FAILED", httpCode=400

// 格式化消息
err := OrderErrors.Code("PAYMENT_FAILED").Msgf("余额不足: %.2f", balance)

// 自定义 HTTP 状态码
err := OrderErrors.Code("PAYMENT_REQUIRED").MsgWithCode("需要支付", 402)
```

## 🎨 Module 快捷方法

| 方法 | 生成错误码 | HTTP 状态码 | 说明 |
|------|-----------|------------|------|
| `NotFound(msg)` | `模块.NOT_FOUND` | 404 | 资源不存在 |
| `AlreadyExists(msg)` | `模块.ALREADY_EXISTS` | 409 | 资源已存在 |
| `InvalidParam(msg)` | `模块.INVALID_PARAM` | 400 | 参数无效 |
| `InvalidStatus(msg)` | `模块.INVALID_STATUS` | 400 | 状态无效 |
| `PermissionDenied(msg)` | `模块.PERMISSION_DENIED` | 403 | 权限不足 |
| `Unauthorized(msg)` | `模块.UNAUTHORIZED` | 401 | 未授权 |
| `OperationFailed(msg)` | `模块.OPERATION_FAILED` | 400 | 操作失败 |
| `Expired(msg)` | `模块.EXPIRED` | 410 | 资源已过期 |
| `Locked(msg)` | `模块.LOCKED` | 423 | 资源已锁定 |
| `LimitExceeded(msg)` | `模块.LIMIT_EXCEEDED` | 429 | 超出限制 |
| `Conflict(msg)` | `模块.CONFLICT` | 409 | 资源冲突 |
| `Internal(msg)` | `模块.INTERNAL_ERROR` | 500 | 内部错误 |
| `ServiceUnavailable(msg)` | `模块.SERVICE_UNAVAILABLE` | 503 | 服务不可用 |

**使用示例：**

```go
// 可以不传消息，使用默认消息
err := UserErrors.NotFound()  // message="资源不存在"

// 传入自定义消息
err := UserErrors.NotFound("用户不存在")  // message="用户不存在"
```

## 🔧 自定义错误码

### 基本用法

```go
// Code().Msg() - 自定义错误码 + 消息
err := OrderErrors.Code("PAYMENT_FAILED").Msg("支付失败")

// Code().Msgf() - 格式化消息
err := OrderErrors.Code("PAYMENT_FAILED").Msgf("订单 %s 支付失败", orderID)

// Code().MsgWithCode() - 自定义 HTTP 状态码
err := OrderErrors.Code("PAYMENT_REQUIRED").MsgWithCode("需要支付", 402)
```

### 完全自定义

```go
// Custom() - 完全控制
err := OrderErrors.Custom("RARE_ERROR", "罕见错误", 418)

// Customf() - 带格式化
err := OrderErrors.Customf("RARE_ERROR", 418, "错误: %s", reason)
```

## 🔗 链式调用

所有方法返回新实例，支持流畅的链式调用：

### 添加详细信息

```go
err := OrderErrors.NotFound("订单不存在").
    WithDetail("orderId", "20231222001").
    WithDetail("userId", 123)
```

### 修改消息

```go
// 创建基础错误
baseErr := OrderErrors.Code("PAYMENT_FAILED").Msg("支付失败")

// 根据条件修改消息
if networkError {
    err = baseErr.WithMsg("网络异常，请稍后重试")
} else if balanceError {
    err = baseErr.WithMsgf("余额不足: %.2f", balance)
}

// 追加消息
err := OrderErrors.NotFound("订单不存在").AppendMsg("，请联系客服")

// 前置消息
err := OrderErrors.OperationFailed("创建失败").PrependMsg("订单")
// 结果: "订单创建失败"
```

### 包装底层错误

```go
user, err := repo.FindByID(id)
if err != nil {
    return UserErrors.NotFound("用户不存在").
        WithDetail("userId", id).
        Wrap(err)  // 包装原始错误，支持 errors.Unwrap()
}
```

### 覆盖 HTTP 状态码

```go
err := OrderErrors.NotFound("订单不存在").
    WithHTTPCode(410)  // 覆盖默认的 404，使用 410 Gone
```

### 组合使用

```go
err := OrderErrors.Code("PAYMENT_FAILED").
    Msgf("订单 %s 支付失败", orderID).
    WithDetail("orderId", orderID).
    WithDetail("amount", amount).
    WithDetail("reason", "余额不足").
    WithDetail("retryable", true).
    WithHTTPCode(402).
    Wrap(originalErr)
```

## 🌐 Web 层集成

错误会自动转换为标准的 API 响应格式。

### Controller 中使用

```go
func (c *UserController) GetUser(ctx *web.HttpContext) web.IActionResult {
    id := ctx.Params().Int("id").Value()
    
    user, err := c.service.GetUser(id)
    if err != nil {
        // 自动处理错误，映射到合适的 HTTP 状态码
        return ctx.Error(err)
    }
    
    return ctx.Ok(user)
}
```

### 错误响应格式

```json
{
  "success": false,
  "error": {
    "code": "USER.NOT_FOUND",
    "message": "用户不存在",
    "details": {
      "userId": 123
    }
  }
}
```

### HTTP 状态码映射

框架会自动根据错误码映射到合适的 HTTP 状态码：

| 错误码模式 | HTTP 状态码 |
|----------|------------|
| `*.NOT_FOUND` | 404 |
| `*.ALREADY_EXISTS` | 409 |
| `*.INVALID_*` | 400 |
| `*.PERMISSION_DENIED` | 403 |
| `*.UNAUTHORIZED` | 401 |
| `*.EXPIRED` | 410 |
| `*.LOCKED` | 423 |
| `*.LIMIT_EXCEEDED` | 429 |
| 其他 | 400 |

## 💡 最佳实践

### 1. 统一定义模块错误

```go
// services/errors.go
package services

import "github.com/gocrud/csgo/errors"

var (
    UserErrors  = errors.NewModule("USER")
    OrderErrors = errors.NewModule("ORDER")
    DramaErrors = errors.NewModule("DRAMA")
    // ... 其他模块
)
```

### 2. Service 层返回错误

```go
func (s *OrderService) GetOrder(id string) (*Order, error) {
    order, err := s.repo.FindByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, OrderErrors.NotFound("订单不存在").
                WithDetail("orderId", id)
        }
        return nil, OrderErrors.Internal("查询订单失败").Wrap(err)
    }
    
    return order, nil
}
```

### 3. 根据业务逻辑返回不同错误

```go
func (s *OrderService) ProcessPayment(orderID string, amount float64) error {
    order, err := s.GetOrder(orderID)
    if err != nil {
        return err
    }
    
    // 检查订单状态
    if order.Status == "cancelled" {
        return OrderErrors.InvalidStatus("订单已取消，无法支付")
    }
    
    if order.Status == "paid" {
        return OrderErrors.AlreadyExists("订单已支付")
    }
    
    // 检查余额
    balance, _ := s.accountService.GetBalance(order.UserID)
    if balance < amount {
        return OrderErrors.Code("PAYMENT_FAILED").
            Msgf("余额不足，当前: %.2f，需要: %.2f", balance, amount).
            WithDetail("balance", balance).
            WithDetail("required", amount)
    }
    
    return nil
}
```

### 4. 渐进式错误构建

```go
func (s *OrderService) ValidateOrder(order *Order) error {
    // 基础错误
    baseErr := OrderErrors.InvalidParam()
    
    if order.Amount <= 0 {
        return baseErr.WithMsg("订单金额必须大于0").
            WithDetail("amount", order.Amount)
    }
    
    if len(order.Items) == 0 {
        return baseErr.WithMsg("订单商品不能为空")
    }
    
    return nil
}
```

### 5. 错误链追踪

```go
func (s *OrderService) CreateOrder(req *CreateOrderRequest) error {
    // 调用多层服务
    if err := s.validateStock(req.Items); err != nil {
        return OrderErrors.OperationFailed("创建订单失败").
            WithMsg("库存不足").
            Wrap(err)  // 保留原始错误链
    }
    
    // 可以在外层判断根因
    // if errors.Is(err, StockErrors.NotEnough) { ... }
    
    return nil
}
```

## 📚 完整 API 参考

### Module 方法

| 方法 | 说明 |
|------|------|
| `NotFound(msg...)` | 资源不存在 (404) |
| `AlreadyExists(msg...)` | 资源已存在 (409) |
| `InvalidParam(msg...)` | 参数无效 (400) |
| `InvalidStatus(msg...)` | 状态无效 (400) |
| `PermissionDenied(msg...)` | 权限不足 (403) |
| `Unauthorized(msg...)` | 未授权 (401) |
| `OperationFailed(msg...)` | 操作失败 (400) |
| `Expired(msg...)` | 资源已过期 (410) |
| `Locked(msg...)` | 资源已锁定 (423) |
| `LimitExceeded(msg...)` | 超出限制 (429) |
| `Conflict(msg...)` | 资源冲突 (409) |
| `Internal(msg...)` | 内部错误 (500) |
| `ServiceUnavailable(msg...)` | 服务不可用 (503) |
| `Code(code)` | 自定义错误码构建器 |
| `Custom(code, msg, httpCode)` | 完全自定义错误 |

### Error 方法

| 方法 | 说明 | 不可变 |
|------|------|--------|
| `WithMsg(msg)` | 覆盖错误消息 | ✅ |
| `WithMsgf(format, args...)` | 格式化覆盖消息 | ✅ |
| `AppendMsg(suffix)` | 追加消息 | ✅ |
| `PrependMsg(prefix)` | 前置消息 | ✅ |
| `WithDetail(key, value)` | 添加详细信息 | ✅ |
| `WithDetails(map)` | 批量添加详细信息 | ✅ |
| `WithHTTPCode(code)` | 设置 HTTP 状态码 | ✅ |
| `Wrap(err)` | 包装原始错误 | ✅ |
| `Error()` | 实现 error 接口 | - |
| `Unwrap()` | 返回原始错误 | - |
| `Code()` | 获取错误码 | - |
| `Message()` | 获取错误消息 | - |
| `HTTPCode()` | 获取 HTTP 状态码 | - |
| `Details()` | 获取详细信息 | - |
| `Category()` | 获取错误分类 | - |

### ErrorBuilder 方法

| 方法 | 说明 |
|------|------|
| `Msg(msg)` | 设置消息，返回 Error |
| `Msgf(format, args...)` | 格式化消息，返回 Error |
| `MsgWithCode(msg, httpCode)` | 设置消息和 HTTP 状态码 |

## 🔄 迁移指南

### 从旧版本迁移

#### 1. 替换 `Business()` 为 `NewModule()`

```go
// ❌ 旧方式
var ErrUser = errors.Business("USER")
err := ErrUser.NotFound("用户不存在")

// ✅ 新方式
var UserErrors = errors.NewModule("USER")
err := UserErrors.NotFound("用户不存在")
```

#### 2. 替换 `Code().Message()` 为 `Code().Msg()`

```go
// ❌ 旧方式
err := ErrUser.Code("CUSTOM").Message("自定义错误")

// ✅ 新方式
err := UserErrors.Code("CUSTOM").Msg("自定义错误")
```

#### 3. `BizError` 改为 `Error`

```go
// ❌ 旧方式
func GetUser(id int) (*User, *errors.BizError) {
    return nil, errors.Business("USER").NotFound("用户不存在")
}

// ✅ 新方式
func GetUser(id int) (*User, error) {
    return nil, UserErrors.NotFound("用户不存在")
}
```

#### 4. 字段访问改为方法调用

```go
// ❌ 旧方式（会导致编译错误）
code := err.Code
message := err.Message

// ✅ 新方式
code := err.Code()
message := err.Message()
```

**注意**: 新版本的所有方法都返回新实例，是线程安全的！

## 🎉 优势总结

与旧版本相比的改进：

1. **更简洁** - `UserErrors.NotFound()` vs `errors.Business("USER").NotFound()`
2. **更安全** - 不可变设计，避免并发问题
3. **更灵活** - 支持 `WithMsg()` 动态修改消息
4. **更强大** - 完整的链式调用和错误链支持
5. **更清晰** - Module 组织错误，避免全局污染

---

更多完整示例请查看 [EXAMPLES.md](./EXAMPLES.md)

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
