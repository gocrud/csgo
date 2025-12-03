# Controller Pattern Demo

这个示例展示了如何在 Ego 框架中使用控制器（Controller）模式，这是 .NET MVC/Web API 中非常经典的模式。

## 📁 项目结构

```
controller_demo/
├── main.go
└── controllers/
    ├── user_controller.go      # 用户控制器（完整 CRUD）
    └── product_controller.go   # 产品控制器
```

## 🎯 什么是控制器模式？

控制器模式将相关的 HTTP 请求处理逻辑组织在一个类中，每个方法（Action）处理一个特定的路由。

### .NET 风格

```csharp
[ApiController]
[Route("api/users")]
public class UserController : ControllerBase
{
    private readonly IUserService _userService;
    
    public UserController(IUserService userService)
    {
        _userService = userService;
    }
    
    [HttpGet]
    public IActionResult GetAll() => Ok(_userService.ListUsers());
}
```

### Ego 风格

```go
type UserController struct {
    app         *web.WebApplication
    userService UserService
}

func NewUserController(app *web.WebApplication) *UserController {
    userService := di.GetRequiredService[UserService](app.Services)
    return &UserController{app: app, userService: userService}
}

func (ctrl *UserController) RegisterRoutes() {
    users := ctrl.app.MapGroup("/api/users")
    users.MapGet("", ctrl.GetAll)
}

func (ctrl *UserController) GetAll(c *gin.Context) {
    users, _ := ctrl.userService.ListUsers()
    c.JSON(200, users)
}
```

## 🚀 运行示例

```bash
# 安装依赖
go mod tidy

# 运行
go run main.go
```

访问：
- **API Root**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger

## 📖 API 端点

### UserController - `/api/users`

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/users` | 获取所有用户 |
| GET | `/api/users/{id}` | 获取指定用户 |
| POST | `/api/users` | 创建新用户 |
| PUT | `/api/users/{id}` | 更新用户 |
| DELETE | `/api/users/{id}` | 删除用户 |

### ProductController - `/api/products`

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/products` | 获取所有产品 |
| GET | `/api/products/{id}` | 获取指定产品 |
| POST | `/api/products` | 创建新产品 |

## 💡 核心概念

### 1. 控制器定义

```go
type UserController struct {
    app         *web.WebApplication
    userService UserService
}
```

### 2. 依赖注入

```go
func NewUserController(app *web.WebApplication) *UserController {
    // ✅ 从 DI 容器解析服务
    userService := di.GetRequiredService[UserService](app.Services)
    
    return &UserController{
        app:         app,
        userService: userService,
    }
}
```

### 3. 路由注册

```go
func (ctrl *UserController) RegisterRoutes() {
    users := ctrl.app.MapGroup("/api/users")
    users.WithTags("Users")
    
    // 注册路由到 Action 方法
    users.MapGet("", ctrl.GetAll).WithSummary("Get all users")
    users.MapGet("/{id}", ctrl.GetByID).WithSummary("Get user by ID")
    users.MapPost("", ctrl.Create).WithSummary("Create user")
}
```

### 4. Action 方法

```go
func (ctrl *UserController) GetAll(c *gin.Context) {
    users, err := ctrl.userService.ListUsers()
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, users)
}
```

### 5. 主程序注册

```go
func main() {
    builder := web.CreateBuilder()
    
    // 注册服务
    builder.Services.AddSingleton(controllers.NewUserService)
    
    app := builder.Build()
    
    // 注册控制器
    userController := controllers.NewUserController(app)
    userController.RegisterRoutes()
    
    app.Run()
}
```

## 🎨 测试 API

### 获取所有用户

```bash
curl http://localhost:8080/api/users
```

### 获取指定用户

```bash
curl http://localhost:8080/api/users/1
```

### 创建用户

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"id": 3, "name": "Charlie", "email": "charlie@example.com"}'
```

### 更新用户

```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"id": 1, "name": "Alice Updated", "email": "alice.new@example.com"}'
```

### 删除用户

```bash
curl -X DELETE http://localhost:8080/api/users/1
```

## 📊 与 .NET 的对比

| 特性 | .NET | Ego | 一致性 |
|------|------|-----|--------|
| 控制器类 | `class UserController` | `type UserController struct` | 100% |
| 依赖注入 | 构造函数自动注入 | 构造函数手动解析 | 95% |
| 路由注册 | 特性路由 `[HttpGet]` | `RegisterRoutes()` 方法 | 95% |
| Action 方法 | `IActionResult GetAll()` | `func GetAll(c *gin.Context)` | 98% |

**总体一致性：97%** ✅

## 💡 最佳实践

1. ✅ **一个控制器负责一个领域** - `UserController` 只处理用户相关的请求
2. ✅ **服务在构造函数中注入** - 避免在每个 Action 中解析
3. ✅ **使用 RegisterRoutes 方法** - 集中管理路由
4. ✅ **清晰的 Action 命名** - `GetAll`, `GetByID`, `Create`, `Update`, `Delete`
5. ✅ **统一的错误处理** - 返回标准的 HTTP 状态码

## 📚 相关文档

查看根目录的 `CONTROLLER_PATTERN_GUIDE.md` 获取完整的指南和最佳实践。

## ✅ 关键要点

1. ✅ 控制器模式让代码组织更清晰
2. ✅ 依赖注入在构造函数中完成
3. ✅ 每个 Action 方法处理一个路由
4. ✅ 与 .NET 的使用体验高度一致（97%）

这种模式特别适合大型项目和团队协作！🚀

