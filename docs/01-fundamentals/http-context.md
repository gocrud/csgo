# HTTP 上下文

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

本章讲解 HttpContext 和 ActionResult。

## 完整文档

关于 HttpContext 和 ActionResult 的完整详细文档，请查看：

👉 **[Web 框架完整文档 - HttpContext 和 ActionResult 部分](../../web/README.md#httpcontext)**

## 核心内容概览

### 1. HttpContext
- 获取请求信息
- 请求体绑定
- 请求验证
- 访问服务

### 2. ActionResult
- 成功响应（Ok、Created、NoContent）
- 错误响应（BadRequest、NotFound、InternalError）
- 验证错误响应
- 业务错误响应

### 3. 请求绑定
- BindJSON
- MustBindJSON
- BindQuery
- BindAndValidate

## 快速示例

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func createUser(c *web.HttpContext) web.IActionResult {
    // 绑定并验证请求
    req, err := web.BindAndValidate[CreateUserRequest](c)
    if err != nil {
        return err
    }
    
    // 访问服务
    userService := di.Get[*UserService](c.Services)
    
    // 创建用户
    user, err := userService.Create(req)
    if err != nil {
        return c.InternalError("创建失败")
    }
    
    // 返回 201 Created
    return c.Created(user)
}
```

## 响应格式

**成功响应：**

```json
{
  "success": true,
  "data": { /* 数据 */ }
}
```

**错误响应：**

```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "资源不存在"
  }
}
```

## 下一步

继续学习：[实践项目：简单 API](project-simple-api.md) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

