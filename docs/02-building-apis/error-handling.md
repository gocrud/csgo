# 错误处理

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

统一的错误处理让API更规范和易用。

## 完整文档

关于错误处理的完整详细文档，请查看：

👉 **[错误处理完整文档](../../errors/README.md)**

## 快速示例

```go
// 在服务层抛出业务错误
func (s *UserService) GetUser(id int) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    if user == nil {
        // 使用业务错误构建器
        return nil, errors.Business("USER").NotFound("用户不存在")
    }
    
    return user, nil
}

// 在控制器层转换为HTTP响应
func (ctrl *UserController) GetUser(c *web.HttpContext) web.IActionResult {
    id := c.Params().PathInt("id").Value()
    user, err := ctrl.service.GetUser(id)
    
    if err != nil {
        if bizErr, ok := err.(*errors.BizError); ok {
            return c.BizError(bizErr)  // 自动映射HTTP状态码
        }
        return c.InternalError("服务器错误")
    }
    
    return c.Ok(user)
}
```

## 下一步

继续学习：[API 文档](api-docs.md) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

