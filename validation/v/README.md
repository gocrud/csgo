# Validation V - 新一代验证器

基于包装类型和链式调用的全新验证器架构，提供更直观、更高效的 API 验证方案。

## ✨ 特性

- ✅ **包装类型设计**：使用 `v.String`、`v.Int`、`v.Slice[T]` 等包装类型
- ✅ **链式调用**：直接在字段上调用验证方法，如 `req.Name.MinLen(2).MaxLen(50).Msg("错误信息")`
- ✅ **自动字段追踪**：支持嵌套字段路径，如 `req.Contact.Phone` 自动识别为 "contact.phone"
- ✅ **元数据预注册**：在 `init()` 中注册验证函数，应用启动时收集元数据，提升运行时效率
- ✅ **智能字段名**：优先使用 json tag，如果没有 tag 则直接使用字段名（小驼峰）
- ✅ **JSON 兼容**：完整支持 JSON 序列化和反序列化
- ✅ **类型安全**：编译时类型检查，IDE 友好

## 🚀 快速开始

### 1. 定义 DTO

```go
package main

import "github.com/gocrud/csgo/validation/v"

type CreateUserRequest struct {
    Name     v.String      `json:"name"`
    Email    v.String      `json:"email"`
    Password v.String      `json:"password"`
    Age      v.Int         `json:"age"`
    Tags     v.Slice[string] `json:"tags"`
    Contact  struct {
        Phone v.String `json:"phone"`
        Email v.String `json:"email"`
    } `json:"contact"`
}
```

### 2. 定义验证函数

```go
func validateCreateUserRequest(req CreateUserRequest) {
    // 名称验证
    req.Name.NotEmpty().Msg("名称不能为空")
    req.Name.MinLen(2).Msg("名称至少2个字符")
    req.Name.MaxLen(50).Msg("名称最多50个字符")
    
    // 邮箱验证
    req.Email.NotEmpty().Msg("邮箱不能为空")
    req.Email.Email().Msg("邮箱格式不正确")
    
    // 密码验证
    req.Password.MinLen(8).Msg("密码长度至少8位")
    req.Password.Pattern(`[A-Z]`).Msg("密码必须包含大写字母")
    req.Password.Pattern(`[a-z]`).Msg("密码必须包含小写字母")
    req.Password.Pattern(`[0-9]`).Msg("密码必须包含数字")
    
    // 年龄验证
    req.Age.Min(0).Msg("年龄不能为负数")
    req.Age.Max(150).Msg("年龄不能超过150")
    
    // 标签验证
    req.Tags.MinLen(1).Msg("至少需要一个标签")
    req.Tags.MaxLen(10).Msg("最多10个标签")
    
    // 嵌套字段验证
    req.Contact.Phone.MinLen(11).Msg("手机号至少11位")
    req.Contact.Phone.MaxLen(11).Msg("手机号最多11位")
    req.Contact.Email.Email().Msg("联系邮箱格式不正确")
}
```

### 3. 在 init 中注册

```go
func init() {
    v.Register[CreateUserRequest](validateCreateUserRequest)
}
```

### 4. 在 Handler 中使用

```go
import (
    "github.com/gocrud/csgo/validation/v"
    "github.com/gocrud/csgo/web"
)

func CreateUserHandler(c *web.HttpContext) web.IActionResult {
    // 1. 解析 JSON
    var req CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    // 2. 执行验证
    result := v.Validate(&req)
    if !result.IsValid {
        return c.BadRequest(result.Errors)
    }
    
    // 3. 使用 .Value() 获取实际值
    name := req.Name.Value()
    email := req.Email.Value()
    age := req.Age.Value()
    tags := req.Tags.Value()
    
    // 4. 业务逻辑...
    user := createUser(name, email, age, tags)
    
    return c.Created(user)
}
```

## 📖 完整示例

```go
package main

import (
    "github.com/gocrud/csgo/validation/v"
    "github.com/gocrud/csgo/web"
)

// ========== 1. 定义 DTO ==========

type User struct {
    Name    v.String `json:"name"`
    Age     v.Int    `json:"age"`
    Email   v.String `json:"email"`
    Tags    v.Slice[string] `json:"tags"`
    Contact struct {
        Phone   v.String `json:"phone"`
        Address v.String `json:"address"`
    } `json:"contact"`
}

// ========== 2. 定义验证规则 ==========

func validateUser(req User) {
    req.Name.MinLen(2).Msg("名称至少2个字符")
    req.Age.Range(0, 150).Msg("年龄必须在0-150之间")
    req.Email.Email().Msg("邮箱格式不正确")
    req.Tags.NotEmpty().Msg("至少需要一个标签")
    req.Contact.Phone.MinLen(11).MaxLen(11).Msg("手机号必须是11位")
    req.Contact.Address.MinLen(5).Msg("地址至少5个字符")
}

// ========== 3. 注册验证器 ==========

func init() {
    v.Register[User](validateUser)
}

// ========== 4. 使用验证器 ==========

func main() {
    builder := web.CreateBuilder()
    app := builder.Build()
    
    app.MapPost("/users", func(c *web.HttpContext) web.IActionResult {
        var req User
        if err := c.MustBindJSON(&req); err != nil {
            return err
        }
        
        // 执行验证
        result := v.Validate(&req)
        if !result.IsValid {
            return c.BadRequest(result.Errors)
        }
        
        // 使用实际值
        name := req.Name.Value()
        age := req.Age.Value()
        
        // 业务逻辑...
        return c.Created(web.M{
            "name": name,
            "age":  age,
        })
    })
    
    app.Run()
}
```

## 🎯 支持的类型

### 基础类型

| 包装类型 | 底层类型 | 说明 |
|---------|---------|------|
| `v.String` | `string` | 字符串 |
| `v.Int` | `int` | 整数 |
| `v.Int64` | `int64` | 64位整数 |
| `v.Float64` | `float64` | 浮点数 |
| `v.Bool` | `bool` | 布尔值 |
| `v.Slice[T]` | `[]T` | 切片（泛型） |

### 获取实际值

所有包装类型都提供 `Value()` 方法来获取底层值：

```go
name := req.Name.Value()        // string
age := req.Age.Value()          // int
tags := req.Tags.Value()        // []string
active := req.Active.Value()    // bool
```

## 📋 验证规则

### 字符串规则 (v.String)

```go
// 非空验证
req.Name.NotEmpty().Msg("名称不能为空")

// 长度验证
req.Name.MinLen(2).Msg("至少2个字符")
req.Name.MaxLen(50).Msg("最多50个字符")

// 邮箱格式
req.Email.Email().Msg("邮箱格式不正确")

// 正则匹配
req.Password.Pattern(`[A-Z]`).Msg("必须包含大写字母")
req.Password.Pattern(`[0-9]`).Msg("必须包含数字")
```

### 数字规则 (v.Int, v.Int64, v.Float64)

```go
// 最小值
req.Age.Min(0).Msg("不能为负数")

// 最大值
req.Age.Max(150).Msg("不能超过150")

// 范围
req.Age.Range(0, 150).Msg("必须在0-150之间")
```

### 切片规则 (v.Slice[T])

```go
// 非空验证
req.Tags.NotEmpty().Msg("不能为空")

// 长度验证
req.Tags.MinLen(1).Msg("至少1个元素")
req.Tags.MaxLen(10).Msg("最多10个元素")
```

## 🔑 核心特性详解

### 1. 自动字段路径追踪

验证器会自动追踪嵌套字段的路径：

```go
type Request struct {
    User struct {
        Contact struct {
            Phone v.String `json:"phone"`
        } `json:"contact"`
    } `json:"user"`
}

func validate(req Request) {
    // 字段路径自动识别为 "user.contact.phone"
    req.User.Contact.Phone.MinLen(11).Msg("手机号至少11位")
}
```

错误信息中的字段路径：
```json
{
  "field": "user.contact.phone",
  "message": "手机号至少11位"
}
```

### 2. 字段名提取规则

1. **优先使用 json tag**：
```go
type User struct {
    UserName v.String `json:"name"`  // 字段路径: "name"
}
```

2. **没有 json tag 时使用字段名（小驼峰）**：
```go
type User struct {
    UserName v.String  // 字段路径: "userName"
    Age      v.Int     // 字段路径: "age"
}
```

### 3. 错误消息自定义

使用 `Msg()` 方法设置**最后一个规则**的错误消息：

```go
// ✅ 正确：每个规则单独设置消息
req.Name.MinLen(2).Msg("至少2个字符")
req.Name.MaxLen(50).Msg("最多50个字符")

// ❌ 错误：Msg 只应用到 MaxLen
req.Name.MinLen(2).MaxLen(50).Msg("长度在2-50之间")
// 这种情况下，MinLen 会使用默认消息，只有 MaxLen 使用自定义消息
```

### 4. JSON 序列化

包装类型完全支持 JSON 序列化：

```go
type User struct {
    Name v.String `json:"name"`
    Age  v.Int    `json:"age"`
}

// 序列化
user := User{
    Name: v.String{/* ... */},
    Age:  v.Int{/* ... */},
}
json.Marshal(user)  // {"name":"张三","age":25}

// 反序列化
var user User
json.Unmarshal(jsonData, &user)
// user.Name.Value() 可以获取实际值
```

### 5. 元数据预注册

验证规则在应用启动时注册，运行时直接使用元数据，避免反射开销：

```go
func init() {
    // 注册时：使用反射收集元数据（只执行一次）
    v.Register[CreateUserRequest](validateCreateUserRequest)
}

func handler(c *web.HttpContext) {
    // 运行时：直接使用元数据执行验证（快速）
    result := v.Validate(&req)
}
```

## 🎨 最佳实践

### 1. 结构清晰

每个验证规则单独一行，便于阅读和维护：

```go
func validateUser(req User) {
    // 名称
    req.Name.NotEmpty().Msg("名称不能为空")
    req.Name.MinLen(2).Msg("名称至少2个字符")
    req.Name.MaxLen(50).Msg("名称最多50个字符")
    
    // 邮箱
    req.Email.NotEmpty().Msg("邮箱不能为空")
    req.Email.Email().Msg("邮箱格式不正确")
}
```

### 2. 错误消息友好

使用用户友好的错误消息：

```go
// ✅ 好的错误消息
req.Password.MinLen(8).Msg("密码长度至少8位")
req.Email.Email().Msg("邮箱格式不正确")

// ❌ 不好的错误消息
req.Password.MinLen(8).Msg("长度不足")
req.Email.Email().Msg("格式错误")
```

### 3. 分组验证

对相关字段进行分组，提高可读性：

```go
func validateUser(req User) {
    // 基本信息
    req.Name.NotEmpty().Msg("名称不能为空")
    req.Age.Min(0).Msg("年龄不能为负数")
    
    // 登录信息
    req.Email.Email().Msg("邮箱格式不正确")
    req.Password.MinLen(8).Msg("密码至少8位")
    
    // 联系方式
    req.Contact.Phone.MinLen(11).Msg("手机号至少11位")
    req.Contact.Address.MinLen(5).Msg("地址至少5个字符")
}
```

### 4. 复用验证逻辑

对于复杂的验证，可以提取为辅助函数：

```go
// 强密码验证
func validateStrongPassword(password v.String) {
    password.MinLen(8).Msg("密码至少8位")
    password.Pattern(`[A-Z]`).Msg("必须包含大写字母")
    password.Pattern(`[a-z]`).Msg("必须包含小写字母")
    password.Pattern(`[0-9]`).Msg("必须包含数字")
    password.Pattern(`[!@#$%^&*]`).Msg("必须包含特殊字符")
}

func validateUser(req User) {
    req.Name.NotEmpty().Msg("名称不能为空")
    validateStrongPassword(req.Password)
}
```

## 🆚 与旧验证器对比

| 特性 | 旧验证器 (validation) | 新验证器 (validation/v) |
|------|---------------------|----------------------|
| **字段定义** | `string` | `v.String` |
| **验证方式** | `validation.NotEmpty(v.Field(...))` | `req.Name.NotEmpty()` |
| **字段路径** | 通过反射 + json tag 提取 | 注册时自动注入 |
| **嵌套支持** | 需要手动处理 | 自动追踪 `contact.phone` |
| **错误消息** | `.WithMessage(msg)` | `.Msg(msg)` |
| **性能** | 每次验证都反射 | 预注册元数据，运行时快速 |
| **IDE 支持** | ✅ | ✅✅（更直观） |
| **学习曲线** | 中等 | 低（更直观） |

## 🔧 高级用法

### 调试验证规则

查看已注册的验证元数据：

```go
import "github.com/gocrud/csgo/validation/v"

func main() {
    // 打印元数据
    v.PrintMetadata[CreateUserRequest]()
    
    // 获取元数据
    metadata, ok := v.GetMetadata[CreateUserRequest]()
    if ok {
        fmt.Printf("类型: %s\n", metadata.TypeName)
        fmt.Printf("字段数: %d\n", len(metadata.Rules))
    }
}
```

### 清空注册表（测试用）

```go
func TestSomething(t *testing.T) {
    v.ClearRegistry()
    v.Register[MyType](validateMyType)
    // ...
}
```

## ⚠️ 注意事项

### 1. 必须调用 Value() 获取实际值

```go
// ✅ 正确
name := req.Name.Value()  // string 类型

// ❌ 错误
name := req.Name  // v.String 类型，不是 string
```

### 2. Msg() 只影响最后一个规则

```go
// 这个 Msg 只应用到 MaxLen，MinLen 使用默认消息
req.Name.MinLen(2).MaxLen(50).Msg("长度在2-50之间")

// 如果要为每个规则设置消息，应该分开写
req.Name.MinLen(2).Msg("至少2个字符")
req.Name.MaxLen(50).Msg("最多50个字符")
```

### 3. 必须在 init() 中注册

验证器必须在应用启动前注册：

```go
func init() {
    v.Register[User](validateUser)
}
```

### 4. JSON 序列化正常工作

包装类型会自动序列化为底层值：

```json
{
  "name": "张三",
  "age": 25
}
```

而不是：
```json
{
  "name": {
    "value": "张三",
    "fieldPath": "name"
  }
}
```

## 📚 更多示例

查看 [`v_test.go`](v_test.go) 和 [`integration_test.go`](integration_test.go) 获取更多使用示例。

## 🤝 与现有代码集成

新验证器在 `validation/v` 包下，与旧验证器 (`validation`) **完全独立**，可以在同一项目中共存：

```go
import (
    oldvalidation "github.com/gocrud/csgo/validation"
    "github.com/gocrud/csgo/validation/v"
)

// 旧验证器
type OldRequest struct {
    Name string `json:"name"`
}

// 新验证器
type NewRequest struct {
    Name v.String `json:"name"`
}
```

## 📝 总结

新验证器提供了更直观、更高效的验证方式：

1. **更简洁**：`req.Name.MinLen(2)` vs `validation.MinLength(v.Field(...), 2)`
2. **更快速**：元数据预注册，避免运行时反射
3. **更直观**：链式调用，IDE 友好
4. **更强大**：自动追踪嵌套字段路径

立即开始使用新验证器，让 API 验证更加优雅！
