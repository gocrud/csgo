# 与 ASP.NET Core 对比

CSGO 框架深受 ASP.NET Core 启发，但针对 Go 语言特性进行了优化。本文档详细对比两者的异同。

## 快速对照表

### 依赖注入

| 功能 | .NET | CSGO | 说明 |
|------|------|-----|------|
| 服务容器 | `IServiceCollection` | `di.IServiceCollection` | 接口名称一致 |
| 服务提供者 | `IServiceProvider` | `di.IServiceProvider` | 接口名称一致 |
| 注册 Singleton | `AddSingleton<T>()` | `AddSingleton(factory)` | CSGO 使用工厂函数 |
| 注册 Scoped | `AddScoped<T>()` | `AddScoped(factory)` | CSGO 使用工厂函数 |
| 注册 Transient | `AddTransient<T>()` | `AddTransient(factory)` | CSGO 使用工厂函数 |
| 解析服务 | `GetService<T>()` | `GetService(&target)` | CSGO 使用指针填充 |
| 必需服务 | `GetRequiredService<T>()` | `GetRequiredService(&target)` | CSGO 使用指针填充 |
| 命名服务 | `GetKeyedService<T>()` | `GetKeyedService(&target, key)` | .NET 8+ 功能 |
| 作用域 | `IServiceScope` | `IServiceScope` | 接口名称一致 |
| 作用域工厂 | `IServiceScopeFactory` | `IServiceScopeFactory` | 接口名称一致 |

### Web 应用

| 功能 | .NET | CSGO | 说明 |
|------|------|-----|------|
| 应用构建器 | `WebApplication.CreateBuilder()` | `web.CreateBuilder()` | API 相似 |
| 路由 GET | `app.MapGet()` | `app.MapGet()` | 完全一致 |
| 路由 POST | `app.MapPost()` | `app.MapPost()` | 完全一致 |
| 路由组 | `app.MapGroup()` | `app.MapGroup()` | 完全一致 |
| CORS | `app.UseCors()` | `app.UseCors()` | 完全一致 |
| 运行应用 | `app.Run()` | `app.Run()` | 完全一致 |

### 托管服务

| 功能 | .NET | CSGO | 说明 |
|------|------|-----|------|
| Host Builder | `Host.CreateDefaultBuilder()` | `hosting.CreateDefaultBuilder()` | API 相似 |
| 后台服务 | `IHostedService` | `IHostedService` | 接口名称一致 |
| 应用生命周期 | `IHostApplicationLifetime` | `IHostApplicationLifetime` | 接口名称一致 |

## 详细对比

### 1. 依赖注入

#### 服务注册

**.NET**:
```csharp
// 使用泛型
services.AddSingleton<IUserService, UserService>();
services.AddScoped<IOrderService, OrderService>();
services.AddTransient<IEmailService, EmailService>();

// 使用工厂
services.AddSingleton<IDatabase>(sp => 
{
    var config = sp.GetRequiredService<IConfiguration>();
    return new PostgresDatabase(config.ConnectionString);
});
```

**CSGO**:
```go
// 使用工厂函数（自动依赖解析）
services.AddSingleton(NewUserService)
services.AddScoped(NewOrderService)
services.AddTransient(NewEmailService)

// 带依赖的工厂
services.AddSingleton(func(config *AppConfig) IDatabase {
    return NewPostgresDatabase(config.ConnectionString)
})
```

#### 服务解析

**.NET**:
```csharp
// 泛型解析
var userService = serviceProvider.GetRequiredService<IUserService>();

// 可选服务
var cache = serviceProvider.GetService<ICache>();
if (cache != null) {
    // 使用缓存
}
```

**CSGO**:
```go
// 指针填充方式
var userService IUserService
provider.GetRequiredService(&userService)

// 可选服务
var cache ICache
if err := provider.GetService(&cache); err == nil {
    // 使用缓存
}

// 或使用泛型辅助方法
userService := di.GetRequiredService[IUserService](provider)
cache, err := di.GetService[ICache](provider)
```

**关键差异**：
- .NET 使用泛型返回值
- CSGO 使用指针填充（类似 `json.Unmarshal`）
- CSGO 的方式更符合 Go 惯用法，提供编译时类型检查

#### Keyed Services (.NET 8+)

**.NET**:
```csharp
services.AddKeyedSingleton<IDatabase>("postgres", new PostgresDB());
services.AddKeyedSingleton<IDatabase>("mysql", new MySQL());

var postgresDB = serviceProvider.GetRequiredKeyedService<IDatabase>("postgres");
```

**CSGO**:
```go
services.AddKeyedSingleton("postgres", NewPostgresDB)
services.AddKeyedSingleton("mysql", NewMySQLDB)

var postgresDB IDatabase
provider.GetRequiredKeyedService(&postgresDB, "postgres")
```

### 2. Web 应用

#### 应用启动

**.NET**:
```csharp
var builder = WebApplication.CreateBuilder(args);

// 注册服务
builder.Services.AddControllers();
builder.Services.AddCors();

var app = builder.Build();

// 配置中间件
app.UseCors();
app.UseAuthorization();
app.MapControllers();

app.Run();
```

**CSGO**:
```go
builder := web.CreateBuilder()

// 注册服务
builder.AddControllers()
builder.AddCors()

app := builder.Build()

// 配置中间件
app.UseCors()
app.UseAuthorization()
app.MapControllers()

app.Run()
```

**相似度**: ⭐⭐⭐⭐⭐ (几乎完全一致)

#### 路由定义

**.NET**:
```csharp
app.MapGet("/api/users/{id}", async (int id, IUserService userService) =>
{
    var user = await userService.GetUserAsync(id);
    return Results.Ok(user);
});

app.MapPost("/api/users", async (CreateUserRequest request, IUserService userService) =>
{
    var user = await userService.CreateUserAsync(request);
    return Results.Created($"/api/users/{user.Id}", user);
});
```

**CSGO**:
```go
app.MapGet("/api/users/:id", func(c *gin.Context) {
    var userService IUserService
    app.Services.GetRequiredService(&userService)
    
    id := c.Param("id")
    user := userService.GetUser(id)
    c.JSON(200, user)
})

app.MapPost("/api/users", func(c *gin.Context) {
    var userService IUserService
    app.Services.GetRequiredService(&userService)
    
    var request CreateUserRequest
    c.ShouldBindJSON(&request)
    
    user := userService.CreateUser(&request)
    c.JSON(201, user)
})
```

**差异**：
- .NET 支持参数自动注入
- CSGO 使用 Gin 的上下文（`gin.Context`）
- CSGO 需要手动绑定请求体

### 3. 控制器

#### 控制器定义

**.NET**:
```csharp
[ApiController]
[Route("api/[controller]")]
public class UsersController : ControllerBase
{
    private readonly IUserService _userService;
    
    public UsersController(IUserService userService)
    {
        _userService = userService;
    }
    
    [HttpGet("{id}")]
    public async Task<ActionResult<User>> GetUser(int id)
    {
        var user = await _userService.GetUserAsync(id);
        if (user == null)
            return NotFound();
        return Ok(user);
    }
}
```

**CSGO**:
```go
type UserController struct {
    userService IUserService
}

func NewUserController(userService IUserService) *UserController {
    return &UserController{userService: userService}
}

func (ctrl *UserController) GetUser(c *gin.Context) {
    id := c.Param("id")
    user := ctrl.userService.GetUser(id)
    
    if user == nil {
        c.JSON(404, gin.H{"error": "User not found"})
        return
    }
    
    c.JSON(200, user)
}

// 路由注册
app.MapGet("/api/users/:id", userController.GetUser)
```

**差异**：
- .NET 使用特性（Attribute）进行路由配置
- CSGO 需要显式注册路由
- CSGO 使用构造函数注入（与 .NET 相同）

### 4. 配置系统

**.NET**:
```csharp
var builder = WebApplication.CreateBuilder(args);

// 配置会自动加载 appsettings.json
var connectionString = builder.Configuration["ConnectionStrings:Default"];
var config = builder.Configuration.Get<AppConfig>();
```

**CSGO**:
```go
builder := hosting.CreateDefaultBuilder()

// 配置自动加载 appsettings.json
config := builder.Configuration
connectionString := config.Get("ConnectionStrings:Default")
```

**相似度**: ⭐⭐⭐⭐ (高度相似)

### 5. 中间件

**.NET**:
```csharp
app.Use(async (context, next) =>
{
    // 前置逻辑
    await next();
    // 后置逻辑
});

app.UseMiddleware<CustomMiddleware>();
```

**CSGO**:
```go
app.Use(func(c *gin.Context) {
    // 前置逻辑
    c.Next()
    // 后置逻辑
})

app.Use(CustomMiddleware())
```

**相似度**: ⭐⭐⭐⭐⭐ (完全一致)

## 迁移指南

### 从 .NET 迁移到 CSGO

#### 1. 依赖注入

**.NET 代码**:
```csharp
services.AddSingleton<IUserService, UserService>();
var userService = provider.GetRequiredService<IUserService>();
```

**CSGO 代码**:
```go
services.AddSingleton(NewUserService)
var userService IUserService
provider.GetRequiredService(&userService)
```

#### 2. Web 应用

**.NET 代码**:
```csharp
var builder = WebApplication.CreateBuilder();
builder.Services.AddCors();
var app = builder.Build();
app.UseCors();
app.MapGet("/", () => "Hello");
app.Run();
```

**CSGO 代码**:
```go
builder := web.CreateBuilder()
builder.AddCors()
app := builder.Build()
app.UseCors()
app.MapGet("/", func(c *gin.Context) {
    c.JSON(200, "Hello")
})
app.Run()
```

#### 3. 异步处理

**.NET**: 原生支持 `async/await`
```csharp
public async Task<User> GetUserAsync(int id)
{
    return await _repository.FindByIdAsync(id);
}
```

**CSGO**: Go 使用 goroutine
```go
func (s *UserService) GetUser(id int) *User {
    // Go 的并发通过 goroutine 处理
    return s.repo.FindByID(id)
}
```

## 优势对比

### CSGO 的优势

1. **性能**: Go 的性能通常优于 C#
2. **部署**: 单一二进制文件，无需运行时
3. **并发**: Goroutine 比 Task 更轻量
4. **内存占用**: 更小的内存占用

### .NET 的优势

1. **生态系统**: 更成熟的生态和库
2. **语言特性**: 更丰富的语言特性（如 LINQ）
3. **IDE 支持**: Visual Studio 的强大支持
4. **企业支持**: Microsoft 官方支持

## 概念映射

| 概念 | .NET | CSGO |
|------|------|-----|
| 依赖注入容器 | `IServiceCollection` | `di.IServiceCollection` |
| 服务提供者 | `IServiceProvider` | `di.IServiceProvider` |
| Web 应用 | `WebApplication` | `web.WebApplication` |
| 主机构建器 | `HostBuilder` | `hosting.HostBuilder` |
| 后台服务 | `IHostedService` | `hosting.IHostedService` |
| 配置 | `IConfiguration` | `configuration.IConfiguration` |
| 日志 | `ILogger` | 待实现 |

## 总结

CSGO 框架成功地将 ASP.NET Core 的优秀设计理念带到了 Go 生态系统：

- ✅ **API 相似度高**：80%+ 的 API 保持一致
- ✅ **Go 惯用法**：指针填充、接口设计符合 Go 最佳实践
- ✅ **学习曲线平缓**：.NET 开发者可以快速上手
- ✅ **性能优化**：针对 Go 运行时的优化

如果你熟悉 ASP.NET Core，你会发现 CSGO 非常易于学习和使用！

---

[← 返回文档首页](../README.md) | [快速开始 →](getting-started.md)

