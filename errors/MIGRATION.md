# 错误处理系统重构完成

## ✨ 重构概览

errors 包已完成重构，提供更简洁、类型安全、功能强大的错误处理体系。

## 🎯 主要变化

### 1. 统一错误类型

**之前：**
```go
type BizError struct {
    Code    string
    Message string
    Cause   error
    Details map[string]interface{}
}
```

**现在：**
```go
type Error struct {
    category ErrorCategory
    code     string
    message  string
    cause    error
    details  map[string]any
    httpCode int
}
```

### 2. 模块化错误创建

**之前：**
```go
var ErrUser = errors.Business("USER")
err := ErrUser.NotFound("用户不存在")
```

**现在：**
```go
var UserErrors = errors.NewModule("USER")
err := UserErrors.NotFound("用户不存在")
```

### 3. Web 层集成方法名变更

**之前：**
```go
// action_result.go
BizError(err *errors.BizError)
BizErrorWithStatus(statusCode int, err *errors.BizError)

// http_context.go
ctx.BizError(err)
ctx.BizErrorWithStatus(statusCode, err)
```

**现在：**
```go
// action_result.go
FrameworkError(err *errors.Error)
FrameworkErrorWithStatus(statusCode int, err *errors.Error)

// http_context.go
ctx.FrameworkError(err)
ctx.FrameworkErrorWithStatus(statusCode, err)
```

**推荐使用（无需关心类型）：**
```go
ctx.FromError(err)  // 自动识别错误类型
```

## 📝 迁移步骤

### 1. 更新错误模块定义

```go
// 旧代码
var ErrUser = errors.Business("USER")
var ErrOrder = errors.Business("ORDER")

// 新代码
var UserErrors = errors.NewModule("USER")
var OrderErrors = errors.NewModule("ORDER")
```

### 2. 更新错误创建方式

```go
// 旧代码
err := ErrUser.Code("CUSTOM").Message("自定义错误")

// 新代码
err := UserErrors.Code("CUSTOM").Msg("自定义错误")
```

### 3. 更新函数签名

```go
// 旧代码
func GetUser(id int) (*User, *errors.BizError) {
    return nil, ErrUser.NotFound("用户不存在")
}

// 新代码
func GetUser(id int) (*User, error) {
    return nil, UserErrors.NotFound("用户不存在")
}
```

### 4. 更新字段访问（如果有直接访问）

```go
// 旧代码
code := err.Code
message := err.Message
details := err.Details

// 新代码
code := err.Code()
message := err.Message()
details := err.Details()
```

### 5. 更新 Controller 层

**方式A：使用 FromError（推荐）**
```go
// 无需修改，FromError 自动处理
user, err := service.GetUser(id)
if err != nil {
    return ctx.FromError(err)  // 自动识别 *errors.Error
}
```

**方式B：直接使用（如果需要）**
```go
// 旧代码
return ctx.BizError(err)

// 新代码
return ctx.FrameworkError(err)
```

## 🚀 新特性

### 1. 更简洁的 API

```go
// 常用错误一行搞定
UserErrors.NotFound("用户不存在")
OrderErrors.InvalidParam("金额必须大于0")
DramaErrors.PermissionDenied("无权访问")

// 自定义错误码
OrderErrors.Code("PAYMENT_FAILED").Msg("支付失败")
OrderErrors.Code("PAYMENT_FAILED").Msgf("余额不足: %.2f", balance)
```

### 2. 不可变设计

```go
baseErr := OrderErrors.Code("PAYMENT_FAILED").Msg("支付失败")

// 创建变体，原错误不变
err1 := baseErr.WithMsg("余额不足")
err2 := baseErr.WithMsg("网络异常")

// baseErr 保持不变
```

### 3. 完整的链式调用

```go
OrderErrors.NotFound("订单不存在").
    WithDetail("orderId", id).
    WithDetail("userId", userId).
    WithHTTPCode(410).
    Wrap(dbErr)
```

### 4. 灵活的消息处理

```go
// 创建时指定
OrderErrors.Code("PAYMENT_FAILED").Msg("支付失败")

// 创建后修改
baseErr.WithMsg("余额不足")
baseErr.WithMsgf("余额不足: %.2f", balance)

// 追加/前置
err.AppendMsg("，请联系客服")
err.PrependMsg("订单")
```

## 📊 对比总结

| 特性 | 旧版本 | 新版本 |
|------|--------|--------|
| 错误类型 | BizError | Error |
| 创建方式 | Business() | NewModule() |
| 方法名 | Message() | Msg() / Msgf() |
| 字段访问 | err.Code | err.Code() |
| Web 集成 | BizError() | FrameworkError() 或 FromError() |
| 不可变性 | ❌ 修改原对象 | ✅ 返回新对象 |
| HTTP 状态码 | ❌ 需要手动映射 | ✅ 自动映射 |
| 链式调用 | ⚠️ 部分支持 | ✅ 完整支持 |

## ⚠️ 注意事项

1. **类型变更**：`*errors.BizError` → `*errors.Error`
2. **方法调用**：字段改为方法 `err.Code` → `err.Code()`
3. **命名规范**：建议使用 `UserErrors` 而不是 `ErrUser`
4. **推荐用法**：Controller 层使用 `ctx.FromError(err)` 无需关心具体类型

## 📚 完整文档

- [README.md](./README.md) - 完整的 API 文档和使用指南
- [EXAMPLES.md](./EXAMPLES.md) - 实际业务场景示例
- [module_test.go](./module_test.go) - 单元测试示例

## 🎉 完成

重构已完成，所有测试通过！享受更简洁、强大的错误处理体验！
