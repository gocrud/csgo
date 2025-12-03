# Controller API Demo

这个示例展示了如何使用 ASP.NET Core 风格的 Controller 模式来组织 API。

## 📁 项目结构

```
controller_api_demo/
├── main.go                           # 主程序
├── controllers/                      # 控制器层
│   ├── user_controller.go            # 用户控制器
│   ├── order_controller.go           # 订单控制器
│   └── controller_extensions.go      # 控制器注册扩展
└── services/                         # 服务层
    ├── user_service.go               # 用户服务
    ├── order_service.go              # 订单服务
    └── service_collection_extensions.go  # 服务注册扩展
```

## 🎯 核心概念

### Controller 模式

在 ASP.NET Core 中：
```csharp
[ApiController]
[Route("api/[controller]")]
public class UserController : ControllerBase
{
    private readonly IUserService _userService;
    
    public UserController(IUserService userService)
    {
        _userService = userService;
    }
    
    [HttpGet]
    public IActionResult GetUsers()
    {
        return Ok(_userService.ListUsers());
    }
}
```

在 Ego 中：
```go
type UserController struct {
    userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
    return &UserController{userService: userService}
}

func (ctrl *UserController) RegisterRoutes(app *web.WebApplication) {
    users := app.MapGroup("/api/users")
    users.MapGet("", ctrl.ListUsers)
}

func (ctrl *UserController) ListUsers(c *gin.Context) {
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

### UserController

- `GET /api/users` - 列出所有用户
- `GET /api/users/{id}` - 获取指定用户
- `POST /api/users` - 创建新用户
- `PUT /api/users/{id}` - 更新用户
- `DELETE /api/users/{id}` - 删除用户

### OrderController

- `GET /api/orders/{id}` - 获取指定订单
- `GET /api/orders/user/{userId}` - 获取用户的所有订单
- `POST /api/orders` - 创建新订单

## 💡 使用步骤

### 1. 创建服务层

```go
// services/user_service.go
type UserService interface {
    GetUser(id int) (*User, error)
    ListUsers() ([]*User, error)
}

func NewUserService() UserService {
    return &userService{}
}
```

### 2. 创建控制器

```go
// controllers/user_controller.go
type UserController struct {
    userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
    return &UserController{userService: userService}
}

func (ctrl *UserController) RegisterRoutes(app *web.WebApplication) {
    users := app.MapGroup("/api/users")
    users.MapGet("", ctrl.ListUsers)
}
```

### 3. 注册服务和控制器

```go
// main.go
func main() {
    builder := web.CreateBuilder()
    
    // 注册服务
    services.AddServices(builder.Services)
    
    // 注册控制器
    controllers.AddControllers(builder.Services)
    
    app := builder.Build()
    
    // 映射控制器路由
    controllers.MapControllers(app)
    
    app.Run()
}
```

## 🎨 与 .NET 的对比

### .NET 代码

```csharp
// Startup.cs / Program.cs
builder.Services.AddControllers();

var app = builder.Build();
app.MapControllers();
app.Run();
```

### Ego 代码

```go
// main.go
services.AddServices(builder.Services)
controllers.AddControllers(builder.Services)

app := builder.Build()
controllers.MapControllers(app)
app.Run()
```

**一致性：95%** ✅

## ✅ 优势

1. ✅ **清晰的关注点分离** - Controller、Service、Model 分层
2. ✅ **依赖注入** - 构造函数注入模式
3. ✅ **易于测试** - 可以轻松 mock 服务
4. ✅ **代码组织** - 按功能模块组织
5. ✅ **与 .NET 一致** - 熟悉的开发模式

## 📚 相关文档

查看根目录的 `CONTROLLER_API_GUIDE.md` 获取完整的指南和最佳实践。

