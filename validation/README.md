# Validation Package

FluentValidation 风格的验证器，为 csgo 框架提供类型安全、灵活、可扩展的 API 验证方案。

## 特性

- ✅ **类型安全**：使用 Go 泛型，完整的编译时类型检查
- ✅ **Fluent API**：链式调用，语义清晰
- ✅ **可扩展**：支持自定义验证规则
- ✅ **DI 友好**：验证器可以注入依赖
- ✅ **条件验证**：支持 When/Unless 条件
- ✅ **自定义消息**：WithMessage 方法自定义错误消息
- ✅ **IDE 支持**：完整的自动补全和类型提示

## 快速开始

### 1. 定义 DTO

```go
type CreateUserRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Age      int    `json:"age"`
}
```

### 2. 创建验证器

```go
type CreateUserRequestValidator struct {
    *validation.AbstractValidator[CreateUserRequest]
}

func NewCreateUserRequestValidator() *CreateUserRequestValidator {
    v := &CreateUserRequestValidator{
        AbstractValidator: validation.NewValidator[CreateUserRequest](),
    }
    
    // ✅ 使用 Field 方法 - 自动从 json tag 提取字段名
    
    // Name 验证
    validation.NotEmpty(
        v.Field(func(r *CreateUserRequest) string { return r.Name }),
    ).WithMessage("用户名不能为空")
    
    validation.Length(
        v.Field(func(r *CreateUserRequest) string { return r.Name }),
        2, 50,
    ).WithMessage("用户名长度必须在2-50之间")
    
    // Email 验证
    validation.NotEmpty(
        v.Field(func(r *CreateUserRequest) string { return r.Email }),
    ).WithMessage("邮箱不能为空")
    
    validation.EmailAddress(
        v.Field(func(r *CreateUserRequest) string { return r.Email }),
    ).WithMessage("邮箱格式不正确")
    
    // Password 验证
    validation.NotEmpty(
        v.Field(func(r *CreateUserRequest) string { return r.Password }),
    ).WithMessage("密码不能为空")
    
    validation.MinLength(
        v.Field(func(r *CreateUserRequest) string { return r.Password }),
        8,
    ).WithMessage("密码长度至少8位")
    
    validation.Matches(
        v.Field(func(r *CreateUserRequest) string { return r.Password }),
        `[A-Z]`,
    ).WithMessage("密码必须包含大写字母")
    
    validation.MustString(
        v.Field(func(r *CreateUserRequest) string { return r.Password }),
        func(req *CreateUserRequest, password string) bool {
            return !strings.Contains(password, req.Name)
        },
    ).WithMessage("密码不能包含用户名")
    
    // Age 验证（int 类型）
    validation.InclusiveBetween(
        v.FieldInt(func(r *CreateUserRequest) int { return r.Age }).
            When(func(r *CreateUserRequest) bool { return r.Age > 0 }),
        0, 150,
    ).WithMessage("年龄必须在0-150之间")
    
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
    // 方式 1：使用 BindAndValidate（推荐）
    req, err := web.BindAndValidate[CreateUserRequest](c)
    if err != nil {
        return err
    }
    
    // 业务逻辑...
    return c.Created(result)
}

// 方式 2：手动验证
func (h *CreateUserHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    validator := NewCreateUserRequestValidator()
    result := validator.Validate(&req)
    if !result.IsValid {
        return c.BadRequest(result.Errors.Error())
    }
    
    // 业务逻辑...
}
```

## 验证规则

### 字符串规则

```go
// ✅ Field 方法自动提取字段名（从 json tag 或字段名）
validation.NotEmpty(
    v.Field(func(r *Request) string { return r.Name }),
).WithMessage("不能为空")

validation.Length(
    v.Field(func(r *Request) string { return r.Name }),
    min, max,
).WithMessage("长度必须在指定范围内")

validation.MinLength(v.Field(selector), min)        // 最小长度
validation.MaxLength(v.Field(selector), max)        // 最大长度
validation.EmailAddress(v.Field(selector))          // 邮箱格式
validation.Matches(v.Field(selector), pattern)      // 正则匹配

validation.MustString(
    v.Field(selector),
    func(instance *Request, value string) bool {
        return /* 验证逻辑 */
    },
)
```

### 数字规则

```go
// 使用 FieldInt, FieldInt64, FieldFloat64 等
validation.GreaterThan(
    v.FieldInt(func(r *Request) int { return r.Age }),
    18,
)

validation.GreaterThanOrEqual(v.FieldInt(selector), value)
validation.LessThan(v.FieldInt(selector), value)
validation.LessThanOrEqual(v.FieldInt(selector), value)
validation.InclusiveBetween(v.FieldInt(selector), min, max)
validation.ExclusiveBetween(v.FieldInt(selector), min, max)

validation.MustNumber(
    v.FieldInt(selector),
    func(instance *Request, value int) bool {
        return /* 验证逻辑 */
    },
)
```

### 集合规则

```go
// 使用 FieldSlice 全局函数
validation.NotEmptySlice(
    validation.FieldSlice(v, func(r *Request) []Item { return r.Items }),
)

validation.MinLengthSlice(validation.FieldSlice(v, selector), min)
validation.MaxLengthSlice(validation.FieldSlice(v, selector), max)

validation.MustSlice(
    validation.FieldSlice(v, selector),
    func(instance *Request, value []Item) bool {
        return /* 验证逻辑 */
    },
)
```

## 高级特性

### 条件验证

```go
validation.InclusiveBetween(
    v.FieldInt(func(r *Request) int { return r.Age }).
        When(func(r *Request) bool {
            return r.Type == "adult"  // 只有成人才验证年龄
        }),
    18, 150,
)
```

### 自定义错误消息

```go
validation.NotEmpty(
    v.Field(func(r *Request) string { return r.Name }),
).WithMessage("请输入用户名")

validation.Length(
    v.Field(func(r *Request) string { return r.Name }),
    2, 50,
).WithMessage("用户名长度必须在2-50个字符之间")
```

### 错误码

```go
validation.NotEmpty(
    v.Field(func(r *Request) string { return r.Email }),
).WithCode("EMAIL_REQUIRED")

validation.EmailAddress(
    v.Field(func(r *Request) string { return r.Email }),
).WithCode("EMAIL_INVALID")
```

### 自定义验证规则

```go
// 方式 1：使用 MustString
validation.MustString(
    v.Field(func(r *Request) string { return r.Username }),
    func(req *Request, username string) bool {
        // 检查用户名是否已存在
        exists, _ := userRepo.ExistsByUsername(username)
        return !exists
    },
).WithMessage("用户名已存在")

// 方式 2：使用 CustomRule（跨字段验证）
validation.CustomRule(v, func(req *Request) error {
    // 跨字段验证
    if req.Password != req.ConfirmPassword {
        return errors.New("两次密码不一致")
    }
    return nil
})
```

### 依赖注入

验证器可以注入依赖：

```go
type CreateUserRequestValidator struct {
    *validation.AbstractValidator[CreateUserRequest]
    userRepo repositories.IUserRepository  // 注入仓储
}

func NewCreateUserRequestValidator(
    userRepo repositories.IUserRepository,
) *CreateUserRequestValidator {
    v := &CreateUserRequestValidator{
        AbstractValidator: validation.NewValidator[CreateUserRequest](),
        userRepo:          userRepo,
    }
    
    validation.MustString(
        v.Field(func(r *CreateUserRequest) string { return r.Email }),
        func(req *CreateUserRequest, email string) bool {
            // ✅ 使用注入的仓储检查邮箱唯一性
            exists, _ := v.userRepo.ExistsByEmail(email)
            return !exists
        },
    ).WithMessage("邮箱已被使用")
    
    return v
}

// 在 DI 中注册
func AddValidators(services di.IServiceCollection) {
    services.AddTransient(NewCreateUserRequestValidator)
}
```

## 错误处理

验证失败时返回结构化错误：

```go
result := validator.Validate(&req)
if !result.IsValid {
    // result.Errors 是 []ValidationError
    for _, err := range result.Errors {
        fmt.Printf("Field: %s, Message: %s, Code: %s\n", 
            err.Field, err.Message, err.Code)
    }
    
    // 或者直接返回所有错误
    return c.BadRequest(result.Errors.Error())
}
```

错误格式：

```json
{
  "errors": [
    {
      "field": "Name",
      "message": "用户名不能为空",
      "code": "NAME_REQUIRED"
    },
    {
      "field": "Email",
      "message": "邮箱格式不正确",
      "code": ""
    }
  ]
}
```

## API 参考

### 核心类型

- `IValidator[T]` - 验证器接口
- `AbstractValidator[T]` - 抽象验证器基类
- `ValidationResult` - 验证结果
- `ValidationError` - 单个验证错误
- `ValidationErrors` - 错误集合
- `RuleBuilder[T, TProperty]` - 规则构建器

### 核心函数

- `NewValidator[T]()` - 创建验证器
- `RuleFor(validator, selector, fieldName)` - 定义字段规则
- `RegisterValidator[T](validator)` - 注册验证器
- `GetValidator[T]()` - 获取验证器

### Web 集成

- `web.BindAndValidate[T](c)` - 绑定并验证（推荐）
- `validation.ValidateStruct(instance)` - 验证结构体

## 最佳实践

1. **一个 DTO 一个验证器**：每个 Request DTO 都应该有对应的验证器
2. **验证器放在功能目录**：`features/users/validators.go`
3. **使用 DI 注册**：通过 DI 容器注册验证器，支持依赖注入
4. **自定义消息**：为用户友好的错误消息使用 WithMessage
5. **错误码**：为前端提供明确的错误码使用 WithCode
6. **条件验证**：使用 When/Unless 实现复杂验证逻辑
7. **复用验证器**：对于相似的验证，可以创建辅助函数

## 与 gin binding tag 对比

| 特性 | gin binding tag | FluentValidation |
|------|----------------|------------------|
| 类型安全 | ❌ 字符串，无编译检查 | ✅ 完整类型检查 |
| IDE 支持 | ❌ 无自动补全 | ✅ 完整支持 |
| 自定义规则 | ⚠️ 复杂 | ✅ 简单 |
| 跨字段验证 | ❌ 不支持 | ✅ 支持 |
| 条件验证 | ❌ 不支持 | ✅ 支持 |
| 错误消息 | ⚠️ 有限 | ✅ 完全自定义 |
| 依赖注入 | ❌ 不支持 | ✅ 支持 |
| 可测试性 | ⚠️ 一般 | ✅ 优秀 |

## 示例

查看 [`examples/vertical_slice_demo`](../examples/vertical_slice_demo) 获取完整示例。

## 参考

本验证器参考了 [FluentValidation](https://github.com/FluentValidation/FluentValidation) (C#) 的 API 设计，并适配了 Go 的语言特性。

