# CSGO 框架示例

## 📦 垂直切片架构示例

本目录包含一个完整的企业级项目示例，展示如何使用 csgo 框架构建现代化的 Go 应用。

### [vertical_slice_demo](vertical_slice_demo/) - 垂直切片架构 + 多端隔离

这是一个完整的电商系统示例，展示了：

**核心特性：**
- ✅ **垂直切片架构** - 功能优先，高内聚低耦合
- ✅ **多端隔离** - 管理端、C端、Worker 独立部署
- ✅ **完整的 DI 容器** - 服务注册和依赖注入
- ✅ **HttpContext + ActionResult** - 统一的 HTTP 处理模式
- ✅ **后台服务** - IHostedService 定时任务
- ✅ **共享层设计** - 合理的代码复用

**包含的功能：**

1. **管理端** (端口 8081)
   - 用户管理（创建、列表、更新）
   - 商品管理（创建、列表）
   - 管理员认证

2. **C端** (端口 8080)
   - 用户注册和登录
   - 商品浏览和详情
   - 订单创建和查询

3. **Worker**
   - 订单同步任务（每 30 秒）
   - 邮件发送任务（每 60 秒）

**项目结构：**

```
vertical_slice_demo/
├── cmd/                    # 三个独立入口
│   ├── admin/             # 管理端
│   ├── api/               # C端
│   └── worker/            # 后台服务
│
├── apps/                   # 各端独立业务逻辑
│   ├── admin/
│   │   └── features/      # 功能切片
│   ├── api/
│   │   └── features/
│   └── worker/
│       └── jobs/
│
└── shared/                 # 共享层
    ├── domain/            # 实体模型
    ├── repositories/      # 数据访问
    └── infrastructure/    # 基础设施
```

## 🚀 快速开始

### 1. 进入示例目录

```bash
cd examples/vertical_slice_demo
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 启动服务

```bash
# 启动 C 端 API
make run-api

# 启动管理端 API（新终端）
make run-admin

# 启动 Worker（新终端）
make run-worker
```

### 4. 测试 API

```bash
# 浏览商品
curl http://localhost:8080/api/products

# 创建商品（管理端）
curl -X POST http://localhost:8081/api/admin/products \
  -H "Content-Type: application/json" \
  -H "Authorization: admin-token" \
  -d '{"name":"iPhone 15","price":5999,"stock":100,"status":"active","description":"最新款"}'

# 用户注册
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"测试用户","email":"test@example.com","password":"password123"}'
```

## 📚 详细文档

进入 `vertical_slice_demo/` 目录查看：

- [README.md](vertical_slice_demo/README.md) - 完整的项目文档
- [ARCHITECTURE.md](vertical_slice_demo/ARCHITECTURE.md) - 架构设计详解
- [QUICKSTART.md](vertical_slice_demo/QUICKSTART.md) - 5 分钟快速上手
- [EXAMPLES.md](vertical_slice_demo/EXAMPLES.md) - API 调用示例集

## 🎯 学习路径

### 1. 初学者

从 [QUICKSTART.md](vertical_slice_demo/QUICKSTART.md) 开始，5 分钟快速体验。

### 2. 开发者

阅读 [README.md](vertical_slice_demo/README.md)，了解完整的功能和使用方式。

### 3. 架构师

阅读 [ARCHITECTURE.md](vertical_slice_demo/ARCHITECTURE.md)，深入理解架构设计决策。

### 4. 实战者

参考 [EXAMPLES.md](vertical_slice_demo/EXAMPLES.md)，学习各种 API 调用场景。

## 💡 为什么选择垂直切片架构？

| 维度 | 传统三层架构 | DDD | 垂直切片架构 |
|------|-------------|-----|-------------|
| **组织方式** | 横向分层 | 领域为中心 | **功能为中心** ✅ |
| **改动范围** | 跨多层 | 聚合根内 | **单一目录** ✅ |
| **学习成本** | 低 | 高 | **中等** ✅ |
| **多端支持** | 需设计 | 需设计 | **原生支持** ✅ |
| **团队协作** | 易冲突 | 按领域 | **按功能** ✅ |
| **适用场景** | 简单 CRUD | 复杂业务 | **中等复杂度** ✅ |

## 🎨 核心概念

### 1. 功能切片（Feature Slice）

每个功能包含完整的垂直逻辑：

```
features/users/
├── create_user.go        # 完整的创建逻辑
├── list_users.go         # 完整的列表逻辑
├── controller.go         # HTTP 路由
└── service_extensions.go # DI 注册
```

### 2. 多端隔离（Multi-App）

不同端完全独立：

```
apps/
├── admin/      # 管理端业务
├── api/        # C端业务
└── worker/     # 后台任务
```

### 3. 共享层（Shared）

只共享真正需要共享的：

```
shared/
├── domain/          # 实体模型
├── repositories/    # 数据访问
└── infrastructure/  # 基础设施
```

## 🔧 技术栈

- **框架**: csgo (受 ASP.NET Core 启发)
- **HTTP**: Gin
- **依赖注入**: csgo DI 容器
- **后台任务**: IHostedService
- **配置管理**: JSON + 环境变量

## 🚀 适用场景

这个架构特别适合：

- ✅ 中等复杂度的项目
- ✅ 需要快速迭代的项目
- ✅ 多端（管理端、C端、移动端）的项目
- ✅ 需要独立部署和扩容的项目
- ✅ 团队协作开发的项目

## 📖 扩展阅读

- [csgo 框架主文档](../README.md)
- [依赖注入指南](../docs/guides/dependency-injection.md)
- [控制器指南](../docs/guides/controllers.md)
- [Web 应用指南](../docs/guides/web-applications.md)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

---

**开始你的 csgo 之旅吧！** 🚀

