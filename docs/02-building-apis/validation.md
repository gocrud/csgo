# 请求验证

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

请求验证确保API收到正确的数据。

## 完整文档

关于请求验证的完整详细文档，请查看：

👉 **[验证系统完整文档](../../validation/README.md)**

## 快速示例

```go
// 定义请求
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

// 创建验证器
func NewCreateUserValidator() *validation.AbstractValidator[CreateUserRequest] {
    v := validation.NewValidator[CreateUserRequest]()
    
    v.Field(func(r *CreateUserRequest) string { return r.Name }).
        NotEmpty().
        MinLength(2)
    
    v.Field(func(r *CreateUserRequest) string { return r.Email }).
        NotEmpty().
        EmailAddress()
    
    return v
}

// 注册验证器
func init() {
    validation.RegisterValidator[CreateUserRequest](NewCreateUserValidator())
}

// 使用验证
func createUser(c *web.HttpContext) web.IActionResult {
    req, err := web.BindAndValidate[CreateUserRequest](c)
    if err != nil {
        return err
    }
    // 验证通过，处理业务逻辑
    return c.Created(user)
}
```

## 下一步

继续学习：[错误处理](error-handling.md) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

