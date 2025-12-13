# 控制器模式

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

控制器模式帮助你更好地组织 API 代码。

## 完整文档

关于控制器的完整详细文档，请查看：

👉 **[Web 框架完整文档 - 控制器模式部分](../../web/README.md#控制器模式)**

## 快速示例

```go
type UserController struct {
    userService *UserService
}

func NewUserController(userService *UserService) *UserController {
    return &UserController{userService: userService}
}

func (ctrl *UserController) MapRoutes(app *web.WebApplication) {
    users := app.MapGroup("/api/users")
    users.MapGet("", ctrl.List)
    users.MapGet("/:id", ctrl.Get)
    users.MapPost("", ctrl.Create)
}

// 注册控制器
web.AddController(builder.Services, NewUserController)
app.MapControllers()
```

## 下一步

继续学习：[请求验证](validation.md) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

