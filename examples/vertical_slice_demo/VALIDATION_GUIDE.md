# FluentValidation 验证器使用指南

本示例展示如何在 csgo 框架中使用 FluentValidation 风格的验证器。

## 快速开始

### 1. 查看示例代码

参考 `apps/admin/features/users/validators.go` 查看完整的验证器实现示例。

### 2. 创建验证器

```go
// validators.go
package users

import (
	"github.com/gocrud/csgo/validation"
)

type CreateUserRequestValidator struct {
	*validation.AbstractValidator[CreateUserRequest]
}

func NewCreateUserRequestValidator() *CreateUserRequestValidator {
	v := &CreateUserRequestValidator{
		AbstractValidator: validation.NewValidator[CreateUserRequest](),
	}

	// 添加验证规则
	validation.NotEmpty(
		validation.RuleFor(v.AbstractValidator,
			func(r *CreateUserRequest) string { return r.Name },
			"Name",
		),
	).WithMessage("用户名不能为空")

	validation.Length(
		validation.RuleFor(v.AbstractValidator,
			func(r *CreateUserRequest) string { return r.Name },
			"Name",
		), 2, 50,
	).WithMessage("用户名长度必须在2-50之间")

	return v
}
```

### 3. 注册验证器

```go
func init() {
	validation.RegisterValidator[CreateUserRequest](
		NewCreateUserRequestValidator(),
	)
}
```

### 4. 在 Handler 中使用

```go
func (h *CreateUserHandler) Handle(c *web.HttpContext) web.IActionResult {
	// ✅ 使用 BindAndValidate 自动绑定和验证
	req, err := web.BindAndValidate[CreateUserRequest](c)
	if err != nil {
		return err  // 自动返回 400 + 验证错误
	}

	// 业务逻辑...
	return c.Created(result)
}
```

## 对比：旧方式 vs 新方式

### 旧方式（gin binding tag）

```go
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *Handler) Handle(c *web.HttpContext) web.IActionResult {
	var req CreateUserRequest
	if err := c.MustBindJSON(&req); err != nil {
		return err
	}
	// ...
}
```

**问题：**
- ❌ 字符串形式，无编译检查
- ❌ IDE 无自动补全
- ❌ 拼写错误只在运行时发现
- ❌ 自定义验证复杂
- ❌ 跨字段验证困难

### 新方式（FluentValidation）

```go
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// 验证器在独立文件中
func NewValidator() *Validator {
	v := validation.NewValidator[CreateUserRequest]()
	
	validation.NotEmpty(
		validation.RuleFor(v, func(r *CreateUserRequest) string { return r.Name }, "Name"),
	).WithMessage("用户名不能为空")
	
	validation.EmailAddress(
		validation.RuleFor(v, func(r *CreateUserRequest) string { return r.Email }, "Email"),
	).WithMessage("邮箱格式不正确")
	
	return &Validator{AbstractValidator: v}
}

func (h *Handler) Handle(c *web.HttpContext) web.IActionResult {
	req, err := web.BindAndValidate[CreateUserRequest](c)
	if err != nil {
		return err
	}
	// ...
}
```

**优势：**
- ✅ 完整的类型安全和编译检查
- ✅ IDE 自动补全和类型提示
- ✅ 错误消息自定义灵活
- ✅ 支持复杂验证逻辑
- ✅ 支持依赖注入
- ✅ 易于测试

## 常用验证规则

### 字符串验证

```go
validation.NotEmpty(builder).WithMessage("不能为空")
validation.Length(builder, min, max).WithMessage("长度必须在x-y之间")
validation.MinLength(builder, min).WithMessage("长度不能少于x")
validation.MaxLength(builder, max).WithMessage("长度不能超过x")
validation.EmailAddress(builder).WithMessage("邮箱格式不正确")
validation.Matches(builder, pattern).WithMessage("格式不正确")
validation.MustString(builder, predicate).WithMessage("验证失败")
```

### 数字验证

```go
validation.GreaterThan(builder, value).WithMessage("必须大于x")
validation.GreaterThanOrEqual(builder, value).WithMessage("必须大于等于x")
validation.LessThan(builder, value).WithMessage("必须小于x")
validation.LessThanOrEqual(builder, value).WithMessage("必须小于等于x")
validation.InclusiveBetween(builder, min, max).WithMessage("必须在x到y之间")
validation.MustNumber(builder, predicate).WithMessage("验证失败")
```

### 集合验证

```go
validation.NotEmptySlice(builder).WithMessage("集合不能为空")
validation.MinLengthSlice(builder, min).WithMessage("集合长度不能少于x")
validation.MaxLengthSlice(builder, max).WithMessage("集合长度不能超过x")
validation.MustSlice(builder, predicate).WithMessage("验证失败")
```

### 条件验证

```go
validation.RuleFor(v, selector, "Age").
	When(func(r *Request) bool {
		return r.Type == "adult"
	}).
	InclusiveBetween(18, 150)
```

### 自定义验证（跨字段）

```go
validation.MustString(
	validation.RuleFor(v, func(r *Request) string { return r.Password }, "Password"),
	func(req *Request, password string) bool {
		return !strings.Contains(password, req.Name)
	},
).WithMessage("密码不能包含用户名")
```

## 完整文档

查看 [`csgo/validation/README.md`](../../validation/README.md) 获取完整的 API 文档和高级用法。

## 示例代码位置

- 验证器实现：`apps/admin/features/users/validators.go`
- Handler 使用：`apps/admin/features/users/create_user.go`
- 核心包：`csgo/validation/`

