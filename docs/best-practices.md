# 最佳实践

本文档提供了使用 CSGO 框架开发应用的推荐模式和最佳实践。

## 项目结构

### 推荐的目录结构

```
my-app/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/
│   ├── domain/                  # 领域模型
│   │   ├── user/
│   │   │   ├── model.go        # 实体定义
│   │   │   ├── service.go      # 业务逻辑
│   │   │   └── repository.go   # 数据访问接口
│   │   └── order/
│   ├── infrastructure/          # 基础设施
│   │   ├── database/
│   │   ├── cache/
│   │   └── messaging/
│   ├── application/             # 应用层
│   │   ├── dto/                # 数据传输对象
│   │   └── services/           # 应用服务
│   └── api/                     # API 层
│       ├── controllers/        # 控制器
│       ├── middleware/         # 中间件
│       └── routes.go           # 路由注册
├── pkg/                         # 可导出的公共包
├── configs/                     # 配置文件
│   ├── appsettings.json
│   └── appsettings.development.json
├── docs/                        # 文档
├── scripts/                     # 脚本
├── go.mod
└── README.md
```

## 依赖注入最佳实践

### 1. 使用接口定义契约

**✅ 推荐**：

```go
// 定义接口
type IUserService interface {
    GetUser(id int) (*User, error)
    CreateUser(user *User) error
}

// 实现接口
type UserService struct {
    repo IUserRepository
}

func NewUserService(repo IUserRepository) IUserService {
    return &UserService{repo: repo}
}
```

**❌ 不推荐**：

```go
// 直接使用具体类型
type UserService struct {
    repo *UserRepository  // 紧耦合
}
```

### 2. 选择合适的生命周期

**Singleton** - 用于无状态服务：

```go
// ✅ 适合：数据库连接池、配置、缓存
builder.Services.AddSingleton(NewDatabasePool)
builder.Services.AddSingleton(NewRedisClient)
builder.Services.AddSingleton(NewAppConfig)
```

**Scoped** - 用于请求相关的服务：

```go
// ✅ 适合：数据库事务、请求上下文、工作单元
builder.Services.AddScoped(NewUnitOfWork)
builder.Services.AddScoped(NewRequestContext)
builder.Services.AddScoped(NewUserRepository)
```

**Transient** - 用于轻量级无状态服务：

```go
// ✅ 适合：工具类、验证器、轻量级服务
builder.Services.AddTransient(NewEmailValidator)
builder.Services.AddTransient(NewPasswordHasher)
```

### 3. 避免服务定位器模式

**❌ 不推荐**（服务定位器）：

```go
func (ctrl *UserController) GetUser(c *gin.Context) {
    // 在方法内部解析服务
    var userService IUserService
    c.Get("services").GetRequiredService(&userService)  // 不好
}
```

**✅ 推荐**（构造函数注入）：

```go
type UserController struct {
    userService IUserService
}

func NewUserController(userService IUserService) *UserController {
    return &UserController{userService: userService}
}

func (ctrl *UserController) GetUser(c *gin.Context) {
    // 直接使用注入的服务
    user := ctrl.userService.GetUser(...)
}
```

### 4. 使用工厂函数注册服务

**✅ 推荐**：

```go
builder.Services.AddSingleton(func(config *AppConfig, logger ILogger) IUserService {
    return NewUserService(config, logger)
})
```

这样做的好处：
- 自动解析依赖
- 类型安全
- 易于测试

## 业务模块组织

### 1. 为每个模块创建扩展方法

```go
// internal/domain/user/module.go
package user

func AddUserModule(services di.IServiceCollection) {
    // 注册领域服务
    services.AddScoped(NewUserService)
    services.AddScoped(NewUserRepository)
    
    // 注册验证器
    services.AddTransient(NewUserValidator)
}
```

在 main.go 中使用：

```go
builder := web.CreateBuilder()

// 注册模块
user.AddUserModule(builder.Services)
order.AddOrderModule(builder.Services)
payment.AddPaymentModule(builder.Services)
```

### 2. 模块间通过接口通信

```go
// order 模块依赖 user 模块
type OrderService struct {
    userService user.IUserService  // 依赖接口，不是实现
    orderRepo   IOrderRepository
}
```

### 3. 使用选项模式配置模块

```go
type UserModuleOptions struct {
    EnableCache      bool
    CacheExpiration  time.Duration
    EnableAudit      bool
}

func AddUserModule(services di.IServiceCollection, configure ...func(*UserModuleOptions)) {
    opts := &UserModuleOptions{
        EnableCache:     true,
        CacheExpiration: 5 * time.Minute,
    }
    
    if len(configure) > 0 && configure[0] != nil {
        configure[0](opts)
    }
    
    if opts.EnableCache {
        services.AddSingleton(NewUserCache)
    }
    
    services.AddScoped(NewUserService)
}
```

## 控制器设计

### 1. 保持控制器精简

**✅ 推荐**：

```go
func (ctrl *UserController) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 委托给服务层
    user, err := ctrl.userService.CreateUser(&req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, user)
}
```

控制器只负责：
- 请求绑定和验证
- 调用服务
- 返回响应

### 2. 使用 DTO 而不是直接暴露实体

```go
// DTO
type UserResponse struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    // 不包含敏感字段如 Password
}

// 转换函数
func ToUserResponse(user *User) *UserResponse {
    return &UserResponse{
        ID:       user.ID,
        Username: user.Username,
        Email:    user.Email,
    }
}
```

### 3. 统一错误处理

创建错误处理中间件：

```go
func ErrorHandlingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            switch e := err.Err.(type) {
            case *ValidationError:
                c.JSON(400, gin.H{"error": e.Error()})
            case *NotFoundError:
                c.JSON(404, gin.H{"error": e.Error()})
            case *UnauthorizedError:
                c.JSON(401, gin.H{"error": e.Error()})
            default:
                c.JSON(500, gin.H{"error": "Internal server error"})
            }
        }
    }
}
```

## 配置管理

### 1. 使用强类型配置

**✅ 推荐**：

```go
type AppConfig struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
}

type ServerConfig struct {
    Port int
    Host string
}

func NewAppConfig() *AppConfig {
    config := &AppConfig{}
    // 加载配置
    return config
}

// 注册为单例
builder.Services.AddSingleton(NewAppConfig)
```

### 2. 使用环境特定的配置

```
configs/
├── appsettings.json              # 基础配置
├── appsettings.development.json  # 开发环境
├── appsettings.production.json   # 生产环境
└── appsettings.test.json         # 测试环境
```

### 3. 配置验证

```go
func (c *AppConfig) Validate() error {
    if c.Database.Host == "" {
        return errors.New("database host is required")
    }
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return errors.New("invalid server port")
    }
    return nil
}
```

## 错误处理

### 1. 定义自定义错误类型

```go
type AppError struct {
    Code    string
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

// 预定义错误
var (
    ErrUserNotFound = &AppError{Code: "USER_NOT_FOUND", Message: "User not found"}
    ErrInvalidInput = &AppError{Code: "INVALID_INPUT", Message: "Invalid input"}
)
```

### 2. 错误包装

```go
func (s *UserService) GetUser(id int) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)
    }
    return user, nil
}
```

## 测试

### 1. 使用接口便于测试

```go
// 生产代码
type UserService struct {
    repo IUserRepository
}

// 测试代码
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) FindByID(id int) (*User, error) {
    args := m.Called(id)
    return args.Get(0).(*User), args.Error(1)
}

func TestUserService_GetUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    mockRepo.On("FindByID", 1).Return(&User{ID: 1}, nil)
    
    service := NewUserService(mockRepo)
    user, err := service.GetUser(1)
    
    assert.NoError(t, err)
    assert.Equal(t, 1, user.ID)
    mockRepo.AssertExpectations(t)
}
```

### 2. 测试依赖注入

```go
func TestDI(t *testing.T) {
    services := di.NewServiceCollection()
    services.AddSingleton(NewUserService)
    
    provider := services.BuildServiceProvider()
    
    var userService IUserService
    err := provider.GetRequiredService(&userService)
    
    assert.NoError(t, err)
    assert.NotNil(t, userService)
}
```

## 性能优化

### 1. 避免在热路径中创建 Transient 服务

```go
// ❌ 不推荐：每次请求都创建
builder.Services.AddTransient(NewHeavyService)

// ✅ 推荐：使用 Scoped 或 Singleton
builder.Services.AddScoped(NewHeavyService)
```

### 2. 使用连接池

```go
type DatabasePool struct {
    pool *sql.DB
}

func NewDatabasePool(config *DatabaseConfig) *DatabasePool {
    db, _ := sql.Open("postgres", config.DSN)
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    return &DatabasePool{pool: db}
}
```

### 3. 合理使用缓存

```go
type CachedUserService struct {
    inner IUserService
    cache ICache
}

func (s *CachedUserService) GetUser(id int) (*User, error) {
    // 先查缓存
    if cached, found := s.cache.Get(fmt.Sprintf("user:%d", id)); found {
        return cached.(*User), nil
    }
    
    // 缓存未命中，查数据库
    user, err := s.inner.GetUser(id)
    if err != nil {
        return nil, err
    }
    
    // 写入缓存
    s.cache.Set(fmt.Sprintf("user:%d", id), user, 5*time.Minute)
    return user, nil
}
```

## 安全性

### 1. 输入验证

```go
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}
```

### 2. 密码处理

```go
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### 3. 防止 SQL 注入

```go
// ✅ 使用参数化查询
func (r *UserRepository) FindByUsername(username string) (*User, error) {
    var user User
    err := r.db.QueryRow(
        "SELECT id, username FROM users WHERE username = $1",
        username,  // 参数化
    ).Scan(&user.ID, &user.Username)
    return &user, err
}
```

## 日志记录

### 1. 结构化日志

```go
type ILogger interface {
    Info(msg string, fields ...Field)
    Error(msg string, err error, fields ...Field)
}

type Field struct {
    Key   string
    Value interface{}
}

// 使用
logger.Info("User created", 
    Field{Key: "user_id", Value: user.ID},
    Field{Key: "username", Value: user.Username},
)
```

### 2. 请求日志中间件

```go
func RequestLoggingMiddleware(logger ILogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        c.Next()
        
        logger.Info("Request completed",
            Field{Key: "method", Value: c.Request.Method},
            Field{Key: "path", Value: path},
            Field{Key: "status", Value: c.Writer.Status()},
            Field{Key: "duration", Value: time.Since(start)},
        )
    }
}
```

## 总结

遵循这些最佳实践可以帮助你：
- ✅ 构建可维护的代码
- ✅ 提高代码质量和可测试性
- ✅ 提升应用性能
- ✅ 增强安全性
- ✅ 简化团队协作

---

[← 返回文档首页](../README.md)

