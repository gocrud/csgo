# 垂直切片架构 + 多端隔离示例

这个示例展示了一个现代化的项目结构设计，结合了 **垂直切片架构（Vertical Slice Architecture）** 和 **多端隔离（Multi-App）** 的思想。

## 🎯 核心理念

### 1. 垂直切片架构

每个功能都是一个完整的垂直切片，包含：
- HTTP 处理
- 业务逻辑
- 数据访问
- 验证规则

**优点：**
- ✅ 功能内聚，改一个功能只需要动一个目录
- ✅ 新增功能不影响其他功能
- ✅ 团队可以并行开发不同的功能切片
- ✅ 测试更简单，每个切片可以独立测试

### 2. 多端隔离

每个端都有自己的：
- 独立入口（`cmd/`）
- 独立业务逻辑（`apps/`）
- 独立中间件
- 独立路由配置

**优点：**
- ✅ 不同端的需求变化互不影响
- ✅ 可以独立部署、独立扩容
- ✅ 安全边界清晰（管理端和C端完全隔离）

### 3. 共享层

只有真正需要共享的东西才放在 `shared/`：
- 数据模型
- 数据访问层
- 基础设施（数据库、缓存、消息队列）
- 通用工具

## 📁 项目结构

```
vertical_slice_demo/
├── cmd/                          # 多端入口
│   ├── admin/                    # 管理端入口
│   ├── api/                      # C端入口
│   └── worker/                   # Worker 入口
│
├── apps/                         # 各端的独立业务逻辑
│   ├── admin/                    # 管理端
│   │   ├── bootstrap.go         # 启动配置
│   │   ├── features/            # 功能切片
│   │   │   ├── users/           # 用户管理
│   │   │   │   ├── create_user.go
│   │   │   │   ├── list_users.go
│   │   │   │   ├── update_user.go
│   │   │   │   ├── controller.go
│   │   │   │   └── service_extensions.go
│   │   │   └── products/        # 商品管理
│   │   └── middlewares/         # 管理端中间件
│   │
│   ├── api/                      # C端
│   │   ├── bootstrap.go
│   │   ├── features/
│   │   │   ├── auth/            # 认证
│   │   │   ├── products/        # 商品浏览（C端视角）
│   │   │   └── orders/          # 订单
│   │   └── middlewares/
│   │
│   └── worker/                   # Worker
│       ├── bootstrap.go
│       └── jobs/                # 后台任务
│           ├── order_sync/      # 订单同步
│           └── email_sender/    # 邮件发送
│
├── shared/                       # 共享层
│   ├── domain/                  # 领域模型
│   ├── contracts/               # 接口契约
│   │   └── repositories/
│   ├── repositories/            # 仓储实现
│   ├── infrastructure/          # 基础设施
│   │   ├── database/
│   │   └── cache/
│   └── utils/                   # 工具函数
│
└── configs/                     # 配置
```

## 🚀 快速开始

### 1. 安装依赖

```bash
cd examples/vertical_slice_demo
go mod tidy
```

### 2. 配置应用

项目使用 csgo 框架的配置管理系统，支持多种配置源：

**配置文件：** `configs/config.dev.json`（默认）

**环境变量覆盖：**
```bash
export APP_Server__AdminPort=:9091
export APP_Database__Host=192.168.1.100
```

**命令行覆盖：**
```bash
./bin/admin --server:admin_port=:9091
```

**配置优先级：** 命令行 > 环境变量 > 配置文件

详细说明：[configs/CONFIGURATION_GUIDE.md](configs/CONFIGURATION_GUIDE.md)

### 3. 启动管理端

```bash
go run cmd/admin/main.go
```

访问: http://localhost:8081

### 4. 启动 C 端

```bash
go run cmd/api/main.go
```

访问: http://localhost:8080

### 5. 启动 Worker

```bash
go run cmd/worker/main.go
```

## 📖 API 端点

### 管理端 (端口 8081)

需要在 Header 中添加 `Authorization` 令牌

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/admin/users | 创建用户 |
| GET | /api/admin/users | 用户列表 |
| PUT | /api/admin/users/:id | 更新用户 |
| POST | /api/admin/products | 创建商品 |
| GET | /api/admin/products | 商品列表 |

### C 端 (端口 8080)

| 方法 | 路径 | 说明 | 需要认证 |
|------|------|------|---------|
| POST | /api/auth/register | 用户注册 | ❌ |
| POST | /api/auth/login | 用户登录 | ❌ |
| GET | /api/products | 浏览商品 | ❌ |
| GET | /api/products/:id | 商品详情 | ❌ |
| POST | /api/orders | 创建订单 | ✅ |
| GET | /api/orders/my | 我的订单 | ✅ |

## 💡 使用示例

### 1. 创建用户（管理端）

```bash
curl -X POST http://localhost:8081/api/admin/users \
  -H "Content-Type: application/json" \
  -H "Authorization: admin-token" \
  -d '{
    "name": "张三",
    "email": "zhangsan@example.com",
    "password": "password123",
    "role": "user"
  }'
```

### 2. 用户注册（C端）

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "李四",
    "email": "lisi@example.com",
    "password": "password123"
  }'
```

### 3. 创建商品（管理端）

```bash
curl -X POST http://localhost:8081/api/admin/products \
  -H "Content-Type: application/json" \
  -H "Authorization: admin-token" \
  -d '{
    "name": "iPhone 15",
    "description": "最新款苹果手机",
    "price": 5999.00,
    "stock": 100,
    "status": "active"
  }'
```

### 4. 浏览商品（C端）

```bash
curl http://localhost:8080/api/products
```

### 5. 创建订单（C端）

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: user-token" \
  -d '{
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      }
    ]
  }'
```

## 🎨 架构特点

### 与传统架构对比

| 维度 | 传统三层架构 | DDD | 垂直切片架构 |
|------|-------------|-----|-------------|
| 组织方式 | 横向分层 | 领域为中心 | 功能为中心 |
| 文件改动 | 改一个功能需跨多层 | 在一个聚合根内 | 在一个功能目录内 |
| 学习成本 | 低 | 高 | 中 |
| 多端支持 | 需要额外设计 | 需要额外设计 | 原生支持 |

### 核心优势

1. **功能内聚**：一个功能的所有代码都在一个目录下
2. **多端隔离**：管理端、C端、Worker 完全独立
3. **共享复用**：数据模型、基础设施在 shared/ 中复用
4. **测试友好**：每个功能切片可以独立测试
5. **团队协作**：不同团队可以并行开发不同的端和功能
6. **扩展性好**：新增功能只需在对应端的 features/ 下加目录
7. **部署灵活**：每个端可以独立部署、独立扩容

## 🔍 代码组织原则

### 1. 功能切片（Feature Slice）

每个功能切片包含：

```
features/users/
├── create_user.go        # 创建用户的完整逻辑
├── list_users.go         # 列表的完整逻辑
├── update_user.go        # 更新的完整逻辑
├── controller.go         # HTTP 控制器
└── service_extensions.go # DI 注册
```

### 2. 依赖注入

使用 csgo 框架的 DI 容器：

```go
// 注册功能切片
func AddUserFeature(services di.IServiceCollection) {
    services.AddSingleton(NewCreateUserHandler)
    services.AddSingleton(NewListUsersHandler)
    services.AddSingleton(NewUpdateUserHandler)
    web.AddController(services, NewUserController)
}
```

### 3. 启动配置

每个端都有自己的 bootstrap.go：

```go
func Bootstrap() *web.WebApplication {
    builder := web.CreateBuilder()
    
    // 注册基础设施
    database.AddDatabase(builder.Services)
    repositories.AddRepositories(builder.Services)
    
    // 注册功能切片
    users.AddUserFeature(builder.Services)
    products.AddProductFeature(builder.Services)
    
    app := builder.Build()
    app.Use(middlewares.AdminAuthMiddleware())
    app.MapControllers()
    
    return app
}
```

## 🎓 最佳实践

### 1. 何时创建新的功能切片？

- ✅ 当有一个相对独立的功能时（如用户管理、订单管理）
- ✅ 当团队需要并行开发时
- ❌ 不要过度拆分，避免过多的小切片

### 2. 何时放入 shared/?

- ✅ 真正需要跨端共享的数据模型
- ✅ 数据访问层（Repository）
- ✅ 基础设施（数据库、缓存）
- ❌ 不要把业务逻辑放入 shared/

### 3. 不同端的功能如何共享？

不同端可能有同名的功能（如 products），但实现不同：
- 管理端：完整的 CRUD + 状态管理
- C 端：只读 + 过滤（只显示 active 状态）

它们通过 `shared/repositories` 访问同一份数据。

## 📚 相关文档

- [csgo 框架文档](../../README.md)
- [业务模块示例](../business_module_demo/)
- [控制器示例](../controller_api_demo/)

## 📄 许可证

MIT License

---

**这个架构适合中等复杂度、快速迭代、多端需求的项目！** 🚀
