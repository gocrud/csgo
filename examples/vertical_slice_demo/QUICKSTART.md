# 快速开始指南

## 🎯 5 分钟快速体验

### 步骤 1：编译项目

```bash
cd examples/vertical_slice_demo
make build
```

或者：

```bash
go mod tidy
go build -o bin/admin cmd/admin/main.go
go build -o bin/api cmd/api/main.go
go build -o bin/worker cmd/worker/main.go
```

### 步骤 2：配置应用（可选）

项目默认使用 `configs/config.dev.json` 配置，支持通过环境变量或命令行覆盖：

```bash
# 环境变量方式
export APP_Server__AdminPort=:9091
export APP_Server__ApiPort=:9090

# 命令行方式
go run cmd/admin/main.go --server:admin_port=:9091
```

查看完整配置说明：[configs/CONFIGURATION_GUIDE.md](configs/CONFIGURATION_GUIDE.md)

### 步骤 3：启动服务（三个终端）

**终端 1 - C 端 API：**
```bash
make run-api
# 或 go run cmd/api/main.go
```

**终端 2 - 管理端 API：**
```bash
make run-admin
# 或 go run cmd/admin/main.go
```

**终端 3 - Worker 后台服务：**
```bash
make run-worker
# 或 go run cmd/worker/main.go
```

### 步骤 4：测试 API（第四个终端）

**创建商品（管理端）：**
```bash
curl -X POST http://localhost:8081/api/admin/products \
  -H "Content-Type: application/json" \
  -H "Authorization: admin-token" \
  -d '{
    "name": "iPhone 15",
    "description": "最新款苹果手机",
    "price": 5999,
    "stock": 100,
    "status": "active"
  }'
```

**浏览商品（C端，无需认证）：**
```bash
curl http://localhost:8080/api/products
```

**用户注册：**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试用户",
    "email": "test@example.com",
    "password": "password123"
  }'
```

**创建订单：**
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: token_xxx" \
  -d '{
    "items": [{"product_id": 1, "quantity": 2}]
  }'
```

## 📂 项目结构一览

```
vertical_slice_demo/
├── cmd/                    # 三个独立的入口
│   ├── admin/             # 管理端 :8081
│   ├── api/               # C端 :8080
│   └── worker/            # 后台任务
│
├── apps/                   # 各端的业务逻辑
│   ├── admin/
│   │   └── features/      # 功能切片
│   │       ├── users/     # 用户管理
│   │       └── products/  # 商品管理
│   ├── api/
│   │   └── features/
│   │       ├── auth/      # 认证
│   │       ├── products/  # 商品浏览
│   │       └── orders/    # 订单
│   └── worker/
│       └── jobs/          # 后台任务
│
└── shared/                 # 共享层
    ├── domain/            # 实体
    ├── repositories/      # 数据访问
    └── infrastructure/    # 基础设施
```

## 🎨 核心概念

### 1. 垂直切片架构

每个功能都是一个完整的切片：

```
features/users/
├── create_user.go        # 创建用户的全部逻辑
├── list_users.go         # 列表的全部逻辑
├── update_user.go        # 更新的全部逻辑
├── controller.go         # HTTP 路由
└── service_extensions.go # DI 注册
```

### 2. 多端隔离

- **管理端**：用户管理、商品管理、数据统计
- **C端**：注册登录、商品浏览、下单
- **Worker**：后台任务、定时同步

每个端有自己的：
- 独立入口
- 独立业务逻辑
- 独立中间件
- 独立部署

### 3. 共享层

只有真正需要共享的放在 `shared/`：
- 数据模型（User、Product、Order）
- 数据访问（Repository）
- 基础设施（Database、Cache）
- 工具函数

## 💡 关键代码

### 功能切片示例

```go
// apps/admin/features/users/create_user.go
type CreateUserHandler struct {
    userRepo repositories.IUserRepository
}

func (h *CreateUserHandler) Handle(c *web.HttpContext) web.IActionResult {
    var req CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    // 业务验证
    if h.userRepo.ExistsByEmail(req.Email) {
        return c.BadRequest("邮箱已存在")
    }
    
    // 创建实体
    user := &domain.User{...}
    
    // 持久化
    h.userRepo.Create(user)
    
    return c.Created(user)
}
```

### DI 注册

```go
// apps/admin/features/users/service_extensions.go
func AddUserFeature(services di.IServiceCollection) {
    services.AddSingleton(NewCreateUserHandler)
    services.AddSingleton(NewListUsersHandler)
    web.AddController(services, NewUserController)
}
```

### Bootstrap

```go
// apps/admin/bootstrap.go
func Bootstrap() *web.WebApplication {
    builder := web.CreateBuilder()
    
    // 注册基础设施
    database.AddDatabase(builder.Services)
    repositories.AddRepositories(builder.Services)
    
    // 注册功能
    users.AddUserFeature(builder.Services)
    products.AddProductFeature(builder.Services)
    
    app := builder.Build()
    app.Use(middlewares.AdminAuthMiddleware())
    app.MapControllers()
    
    return app
}
```

## 🚀 扩展指南

### 添加新功能

在对应端的 `features/` 下创建新目录：

```bash
mkdir -p apps/admin/features/reports
touch apps/admin/features/reports/{generate_report.go,controller.go,service_extensions.go}
```

然后在 `bootstrap.go` 中注册：

```go
reports.AddReportFeature(builder.Services)
```

### 添加新的端

```bash
mkdir -p cmd/mobile apps/mobile/features
```

创建 `bootstrap.go` 和 `main.go`，参考现有的端。

### 添加共享服务

```bash
mkdir -p shared/services/sms
touch shared/services/sms/{sms_service.go,service_extensions.go}
```

## 📚 更多文档

- [README.md](README.md) - 完整文档
- [ARCHITECTURE.md](ARCHITECTURE.md) - 架构设计
- [EXAMPLES.md](EXAMPLES.md) - API 调用示例

## ❓ 常见问题

**Q: 和传统三层架构有什么区别？**

A: 传统三层是横向分层（Controller -> Service -> Repository），改一个功能需要跨三层。垂直切片是纵向分割，一个功能的所有代码在一个目录。

**Q: 和 DDD 有什么区别？**

A: DDD 以领域为中心，学习曲线陡峭。垂直切片以功能为中心，更简单实用。

**Q: 适合什么场景？**

A: 适合中等复杂度、快速迭代、多端需求的项目。

**Q: 数据存在哪里？**

A: 当前示例数据存在内存中。生产环境应替换为真实数据库（PostgreSQL、MySQL 等）。

---

**开始你的垂直切片之旅吧！** 🚀

