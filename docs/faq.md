# 常见问题 (FAQ)

## 通用问题

### CSGO 是什么？

CSGO 是一个受 ASP.NET Core 启发的 Go Web 框架，提供完整的依赖注入、控制器模式和现代化开发体验。

### 为什么创建 CSGO？

- 将 .NET 的优秀设计理念带到 Go 生态
- 提供更好的代码组织和依赖管理
- 为 .NET 开发者提供熟悉的 API
- 保持 Go 语言的性能优势

### CSGO 是否兼容现有的 Gin 项目？

是的！CSGO 基于 Gin 构建，你可以：
- 在 CSGO 中使用所有 Gin 中间件
- 逐步迁移现有 Gin 项目
- 混合使用 CSGO 和原生 Gin API

### CSGO 需要什么 Go 版本？

Go 1.18 或更高版本（需要泛型支持）。

## 依赖注入

### 为什么使用指针填充而不是泛型返回？

```go
// CSGO 方式（指针填充）
var userService IUserService
provider.GetRequiredService(&userService)

// 其他方式（泛型返回）
userService := provider.GetRequiredService[IUserService]()
```

指针填充的优势：
- ✅ 符合 Go 惯用法（如 `json.Unmarshal`）
- ✅ 编译时类型检查
- ✅ 无需类型断言
- ✅ IDE 友好

CSGO 同时提供泛型辅助方法供选择使用。

### Singleton、Scoped 和 Transient 有什么区别？

- **Singleton**: 整个应用生命周期只创建一次
  - 适合：数据库连接池、配置、缓存
  
- **Scoped**: 每个 HTTP 请求创建一次
  - 适合：数据库事务、请求上下文、工作单元
  
- **Transient**: 每次请求服务时都创建新实例
  - 适合：轻量级无状态服务、工具类

示例：
```go
services.AddSingleton(NewDatabasePool)    // 全局一个
services.AddScoped(NewUnitOfWork)         // 每个请求一个
services.AddTransient(NewEmailValidator)  // 每次都是新的
```

### 如何解决循环依赖？

**最佳实践**：重新设计避免循环依赖

如果确实需要：
1. 使用接口解耦
2. 引入中介者模式
3. 延迟初始化

```go
// ❌ 循环依赖
type A struct {
    b *B
}

type B struct {
    a *A
}

// ✅ 使用接口解耦
type IA interface {
    MethodA()
}

type IB interface {
    MethodB()
}

type A struct {
    b IB
}

type B struct {
    a IA
}
```

### 如何注册多个实现同一接口的服务？

使用 Keyed Services：

```go
// 注册
services.AddKeyedSingleton("postgres", NewPostgresDB)
services.AddKeyedSingleton("mysql", NewMySQLDB)

// 解析
var postgresDB IDatabase
provider.GetKeyedService(&postgresDB, "postgres")

var mysqlDB IDatabase
provider.GetKeyedService(&mysqlDB, "mysql")
```

或使用 `GetServices` 获取所有实现：

```go
var databases []IDatabase
provider.GetServices(&databases)
```

### 服务什么时候被创建？

- **Singleton**: 第一次请求时创建，之后重用
- **Scoped**: 进入作用域时创建
- **Transient**: 每次 `GetService` 时创建

### 如何在服务中使用其他服务？

通过构造函数注入：

```go
type UserService struct {
    repo   IUserRepository
    cache  ICache
    logger ILogger
}

func NewUserService(
    repo IUserRepository,
    cache ICache,
    logger ILogger,
) *UserService {
    return &UserService{
        repo:   repo,
        cache:  cache,
        logger: logger,
    }
}

// 注册（依赖会自动解析）
services.AddScoped(NewUserService)
```

## Web 应用

### 如何处理 CORS？

```go
builder := web.CreateBuilder()

// 添加 CORS
builder.AddCors(func(opts *web.CorsOptions) {
    opts.AllowOrigins = []string{
        "http://localhost:3000",
        "https://example.com",
    }
    opts.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
    opts.AllowCredentials = true
})

app := builder.Build()

// 使用 CORS 中间件
app.UseCors()
```

### 如何添加自定义中间件？

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
            return
        }
        
        // 验证 token
        // ...
        
        c.Next()
    }
}

// 使用
app.Use(AuthMiddleware())
```

### 如何获取请求体？

```go
type CreateUserRequest struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
}

app.MapPost("/users", func(c *gin.Context) {
    var req CreateUserRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 使用 req
    c.JSON(200, gin.H{"message": "success"})
})
```

### 如何返回不同的 HTTP 状态码？

```go
// 200 OK
c.JSON(200, data)

// 201 Created
c.JSON(201, user)

// 400 Bad Request
c.JSON(400, gin.H{"error": "Invalid input"})

// 404 Not Found
c.JSON(404, gin.H{"error": "Not found"})

// 500 Internal Server Error
c.JSON(500, gin.H{"error": "Internal error"})
```

### 如何处理文件上传？

```go
app.MapPost("/upload", func(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 保存文件
    dst := "./uploads/" + file.Filename
    if err := c.SaveUploadedFile(file, dst); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"message": "File uploaded"})
})
```

## 配置

### 如何加载配置文件？

CSGO 自动加载配置文件：

```
appsettings.json                    # 基础配置
appsettings.{Environment}.json      # 环境特定配置
```

```go
builder := hosting.CreateDefaultBuilder()

// 访问配置
value := builder.Configuration.Get("Database:Host")
```

### 如何使用环境变量？

```go
// 环境变量会自动加载
// export APP_DATABASE_HOST=localhost

value := builder.Configuration.Get("APP_DATABASE_HOST")
```

### 如何创建强类型配置？

```go
type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
}

type AppConfig struct {
    Database DatabaseConfig
}

// 加载配置
func NewAppConfig(config configuration.IConfiguration) *AppConfig {
    cfg := &AppConfig{}
    // 绑定配置
    return cfg
}

// 注册
services.AddSingleton(NewAppConfig)
```

## 测试

### 如何测试使用了 DI 的代码？

使用 Mock：

```go
// 生产代码
type UserService struct {
    repo IUserRepository
}

// 测试代码
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) FindByID(id int) (*User, error) {
    args := m.Called(id)
    return args.Get(0).(*User), args.Error(1)
}

func TestUserService(t *testing.T) {
    mockRepo := new(MockRepository)
    mockRepo.On("FindByID", 1).Return(&User{ID: 1}, nil)
    
    service := NewUserService(mockRepo)
    user, err := service.GetUser(1)
    
    assert.NoError(t, err)
    assert.Equal(t, 1, user.ID)
}
```

### 如何测试 HTTP 端点？

```go
func TestAPI(t *testing.T) {
    // 创建测试应用
    builder := web.CreateBuilder()
    builder.Services.AddSingleton(NewMockUserService)
    
    app := builder.Build()
    app.MapGet("/users/:id", GetUser)
    
    // 创建测试请求
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/users/1", nil)
    app.Engine.ServeHTTP(w, req)
    
    // 验证响应
    assert.Equal(t, 200, w.Code)
}
```

## 性能

### CSGO 的性能如何？

CSGO 基于 Gin 构建，保持了 Gin 的高性能特点。DI 系统进行了多项优化：
- typeID 索引避免频繁类型比较
- Lock-free 单例缓存
- 对象池减少 GC 压力
- Unsafe 指针优化

### 如何优化性能？

1. **合理选择服务生命周期**
   ```go
   // ❌ 不要过度使用 Transient
   services.AddTransient(NewHeavyService)
   
   // ✅ 使用 Singleton 或 Scoped
   services.AddScoped(NewHeavyService)
   ```

2. **使用连接池**
   ```go
   db.SetMaxOpenConns(25)
   db.SetMaxIdleConns(5)
   ```

3. **启用缓存**
   ```go
   services.AddSingleton(NewRedisCache)
   ```

### DI 会影响性能吗？

影响极小。依赖解析在应用启动时完成（编译阶段），运行时主要是简单的查找操作。

## 部署

### 如何部署 CSGO 应用？

1. **构建二进制文件**
   ```bash
   go build -o app
   ```

2. **运行**
   ```bash
   ./app
   ```

3. **Docker 部署**
   ```dockerfile
   FROM golang:1.20 AS builder
   WORKDIR /app
   COPY . .
   RUN go build -o server

   FROM alpine:latest
   WORKDIR /app
   COPY --from=builder /app/server .
   CMD ["./server"]
   ```

### 如何配置生产环境？

```bash
# 设置环境
export ENVIRONMENT=production

# 使用 appsettings.production.json
```

## 故障排查

### 服务解析失败怎么办？

检查：
1. 服务是否已注册？
2. 依赖是否都已注册？
3. 生命周期是否正确？

```go
// 检查服务是否注册
if provider.IsService(reflect.TypeOf((*IUserService)(nil)).Elem()) {
    // 已注册
}
```

### 循环依赖错误

错误信息：`circular dependency detected`

解决方案：
1. 重新设计避免循环
2. 使用接口解耦
3. 延迟初始化

### 作用域服务解析失败

错误：从 Singleton 解析 Scoped 服务

```go
// ❌ 错误
services.AddSingleton(func(scoped IScopedService) ISingletonService {
    return NewSingletonService(scoped)  // 不允许
})

// ✅ 正确
services.AddScoped(func(scoped IScopedService) IService {
    return NewService(scoped)  // Scoped 可以注入 Scoped
})
```

## 更多问题？

- 查看[完整文档](../README.md)
- 查看[示例代码](../examples/)
- 提交 [Issue](https://github.com/gocrud/csgo/issues)

---

[← 返回文档首页](../README.md)

