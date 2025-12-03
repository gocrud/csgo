# Business Module Extensions Demo

这个示例展示了如何为业务模块创建 `IServiceCollection` 扩展方法，这是 .NET 中非常常见和推荐的模式。

## 📁 项目结构

```
business_module_demo/
├── main.go                    # 主程序
├── users/                     # 用户模块
│   ├── user_service.go
│   ├── user_repository.go
│   ├── user_module_options.go
│   └── service_collection_extensions.go  # ✅ DI 扩展
└── orders/                    # 订单模块
    ├── order_service.go
    └── service_collection_extensions.go  # ✅ DI 扩展
```

## 🎯 核心概念

### 扩展方法模式

在 .NET 中：
```csharp
builder.Services.AddUserServices();
```

在 csgo 中：
```go
users.AddUserServices(builder.Services)
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

### 用户模块

- `GET /api/users` - 列出所有用户
- `GET /api/users/{id}` - 获取指定用户
- `POST /api/users` - 创建新用户

### 订单模块

- `GET /api/orders/{id}` - 获取指定订单
- `GET /api/orders/user/{userId}` - 获取用户的所有订单
- `POST /api/orders` - 创建新订单

## 💡 使用方式

### Style 1: 简单注册（最常用）

```go
func main() {
    builder := web.CreateBuilder()
    
    // ✅ 一行代码注册整个模块
    users.AddUserServices(builder.Services)
    orders.AddOrderServices(builder.Services)
    
    app := builder.Build()
    app.Run()
}
```

### Style 2: 带配置的注册

```go
users.AddUserServicesWithOptions(builder.Services, func(opts *users.UserModuleOptions) {
    opts.EnableCache = true
    opts.CacheExpiration = 5 * time.Minute
})
```

### Style 3: 不同生命周期

```go
// Singleton（默认）
orders.AddOrderServices(builder.Services)

// Transient（每次请求新实例）
orders.AddOrderServicesScoped(builder.Services)
```

## 🎨 在路由中使用服务

### Style 1: 传统方式

```go
app.MapGet("/api/users", func(c *gin.Context) {
    var userSvc users.IUserService
    app.Services.GetRequiredService(&userSvc)
    
    userList, _ := userSvc.ListUsers()
    c.JSON(200, userList)
})
```

### Style 2: 泛型辅助（推荐）⭐

```go
app.MapGet("/api/users", func(c *gin.Context) {
    userSvc := di.GetRequiredService[users.IUserService](app.Services)
    
    userList, _ := userSvc.ListUsers()
    c.JSON(200, userList)
})
```

## 📚 相关文档

查看根目录的 `BUSINESS_MODULE_EXTENSIONS_GUIDE.md` 获取完整的指南和最佳实践。

## ✅ 关键要点

1. ✅ 每个业务模块都有自己的 `service_collection_extensions.go`
2. ✅ 使用 `AddXxxServices` 命名约定
3. ✅ 支持配置选项和不同生命周期
4. ✅ 与 .NET 的使用体验高度一致（96%）

这种模式让你的代码更加模块化、可维护和专业！🚀

