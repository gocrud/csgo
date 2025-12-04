# API 调用示例

这个文件包含了完整的 API 调用示例，帮助你快速上手。

## 🚀 启动应用

### 方式 1：使用 Make（推荐）

```bash
# 启动管理端
make run-admin

# 启动 C 端（新终端）
make run-api

# 启动 Worker（新终端）
make run-worker
```

### 方式 2：直接运行

```bash
# 启动管理端
go run cmd/admin/main.go

# 启动 C 端
go run cmd/api/main.go

# 启动 Worker
go run cmd/worker/main.go
```

## 📝 完整使用流程

### 场景 1：管理员创建商品

#### 1. 创建管理员用户

```bash
curl -X POST http://localhost:8081/api/admin/users \
  -H "Content-Type: application/json" \
  -H "Authorization: admin-token-123" \
  -d '{
    "name": "管理员张三",
    "email": "admin@example.com",
    "password": "password123",
    "role": "admin"
  }'
```

**响应示例：**
```json
{
  "id": 1,
  "name": "管理员张三",
  "email": "admin@example.com",
  "role": "admin"
}
```

#### 2. 创建商品

```bash
curl -X POST http://localhost:8081/api/admin/products \
  -H "Content-Type: application/json" \
  -H "Authorization: admin-token-123" \
  -d '{
    "name": "iPhone 15 Pro",
    "description": "最新款苹果手机，A17 Pro 芯片",
    "price": 7999.00,
    "stock": 100,
    "status": "active"
  }'
```

**响应示例：**
```json
{
  "id": 1,
  "name": "iPhone 15 Pro",
  "description": "最新款苹果手机，A17 Pro 芯片",
  "price": 7999,
  "stock": 100,
  "status": "active",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

#### 3. 查看所有商品

```bash
curl http://localhost:8081/api/admin/products \
  -H "Authorization: admin-token-123"
```

#### 4. 查看所有用户

```bash
curl http://localhost:8081/api/admin/users \
  -H "Authorization: admin-token-123"
```

### 场景 2：用户注册、登录、下单

#### 1. 用户注册

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "李四",
    "email": "lisi@example.com",
    "password": "password123"
  }'
```

**响应示例：**
```json
{
  "token": "token_1_user_1704096000",
  "user_id": 2,
  "name": "李四",
  "email": "lisi@example.com"
}
```

#### 2. 用户登录

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "lisi@example.com",
    "password": "password123"
  }'
```

#### 3. 浏览商品（无需登录）

```bash
curl http://localhost:8080/api/products
```

**响应示例：**
```json
{
  "products": [
    {
      "id": 1,
      "name": "iPhone 15 Pro",
      "description": "最新款苹果手机，A17 Pro 芯片",
      "price": 7999,
      "in_stock": true
    }
  ],
  "total": 1
}
```

#### 4. 查看商品详情

```bash
curl http://localhost:8080/api/products/1
```

#### 5. 创建订单（需要登录）

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: token_2_user_1704096000" \
  -d '{
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      }
    ]
  }'
```

**响应示例：**
```json
{
  "id": 1,
  "user_id": 1,
  "total_price": 15998,
  "status": "pending",
  "items": [
    {
      "id": 0,
      "order_id": 0,
      "product_id": 1,
      "quantity": 2,
      "price": 7999
    }
  ],
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

#### 6. 查看我的订单

```bash
curl http://localhost:8080/api/orders/my \
  -H "Authorization: token_2_user_1704096000"
```

### 7. 支付订单（使用共享服务）

```bash
curl -X POST http://localhost:8080/api/orders/1/pay \
  -H "Content-Type: application/json" \
  -H "Authorization: token_2_user_1704096000" \
  -d '{
    "payment_method": "alipay"
  }'
```

**响应示例：**
```json
{
  "payment_id": "PAY_1_1704096000",
  "order_id": 1,
  "amount": 15998,
  "method": "alipay",
  "status": "pending",
  "redirect_url": "https://pay.example.com?payment_id=PAY_1_1704096000"
}
```

**后台日志（展示共享服务的调用）：**
```
[PAYMENT] 创建支付
  支付ID: PAY_1_1704096000
  订单ID: 1
  金额: ¥15998.00
  支付方式: alipay
  时间: 2024-01-01 12:00:00

[EMAIL] 发送邮件
  收件人: lisi@example.com
  主题: 订单支付确认
  内容: 您的订单正在支付中，请完成支付。
  时间: 2024-01-01 12:00:00

[PUSH] 发送推送通知
  用户ID: 2
  标题: 订单支付
  内容: 请完成订单支付
  时间: 2024-01-01 12:00:00
```

## 🎨 共享服务示例

### 示例 1：支付订单（展示 PaymentService）

这个示例展示了如何在功能切片中使用共享的支付服务。

```bash
# 1. 创建商品
curl -X POST http://localhost:8081/api/admin/products \
  -H "Content-Type: application/json" \
  -H "Authorization: admin-token" \
  -d '{"name":"MacBook Pro","description":"M3 Max","price":16999,"stock":10,"status":"active"}'

# 2. 用户注册
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"测试用户","email":"test@example.com","password":"password123"}'

# 3. 创建订单
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: <注册返回的token>" \
  -d '{"items":[{"product_id":1,"quantity":1}]}'

# 4. 支付订单（使用共享服务）
curl -X POST http://localhost:8080/api/orders/1/pay \
  -H "Content-Type: application/json" \
  -H "Authorization: <token>" \
  -d '{"payment_method":"alipay"}'
```

**这个过程会调用：**
- ✅ `PaymentService.CreatePayment()` - 创建支付
- ✅ `NotificationService.SendEmail()` - 发送邮件通知
- ✅ `NotificationService.SendPush()` - 发送推送通知

### 示例 2：通知服务的使用

通知服务在多个场景下被使用：

**场景 1：用户注册**
```go
// apps/api/features/auth/register.go
notificationService.SendEmail(user.Email, "欢迎注册", "欢迎...")
```

**场景 2：订单支付**
```go
// apps/api/features/orders/pay_order.go
notificationService.SendEmail(user.Email, "订单支付确认", "...")
notificationService.SendPush(user.ID, "订单支付", "...")
```

**场景 3：管理员操作**
```go
// apps/admin/features/users/create_user.go
notificationService.SendSMS(phone, "账号已创建")
```

## 🧪 测试场景

### 测试 1：完整的电商流程（包含共享服务）

```bash
# 1. 管理员创建商品
curl -X POST http://localhost:8081/api/admin/products \
  -H "Content-Type: application/json" \
  -H "Authorization: admin-token" \
  -d '{"name":"MacBook Pro","description":"M3 Max","price":16999,"stock":50,"status":"active"}'

# 2. 用户注册
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"王五","email":"wangwu@example.com","password":"password123"}'

# 保存返回的 token

# 3. 用户浏览商品
curl http://localhost:8080/api/products

# 4. 用户下单
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: <你的token>" \
  -d '{"items":[{"product_id":1,"quantity":1}]}'

# 5. 查看订单
curl http://localhost:8080/api/orders/my \
  -H "Authorization: <你的token>"

# 6. 支付订单（使用共享服务）
curl -X POST http://localhost:8080/api/orders/1/pay \
  -H "Content-Type: application/json" \
  -H "Authorization: <你的token>" \
  -d '{"payment_method":"wechat"}'
```

### 测试 2：分页查询

```bash
# 查询前 10 个用户
curl "http://localhost:8081/api/admin/users?offset=0&limit=10" \
  -H "Authorization: admin-token"

# 查询第 11-20 个用户
curl "http://localhost:8081/api/admin/users?offset=10&limit=10" \
  -H "Authorization: admin-token"
```

### 测试 3：错误处理

```bash
# 邮箱已存在
curl -X POST http://localhost:8081/api/admin/users \
  -H "Content-Type: application/json" \
  -H "Authorization: admin-token" \
  -d '{"name":"重复用户","email":"admin@example.com","password":"password123","role":"user"}'
# 预期：400 Bad Request - "邮箱已存在"

# 商品不存在
curl http://localhost:8080/api/products/999
# 预期：404 Not Found

# 未授权访问
curl http://localhost:8080/api/orders/my
# 预期：401 Unauthorized
```

## 📊 数据格式说明

### 用户角色

- `admin`: 管理员，可以访问管理端
- `user`: 普通用户，可以访问 C 端

### 商品状态

- `active`: 上架，C 端可见
- `inactive`: 下架，C 端不可见

### 订单状态

- `pending`: 待支付
- `paying`: 支付中
- `paid`: 已支付
- `shipped`: 已发货
- `completed`: 已完成
- `cancelled`: 已取消

### 支付方式

- `alipay`: 支付宝
- `wechat`: 微信支付
- `card`: 银行卡

## 🐛 常见问题

### Q: 端口被占用怎么办？

修改 `configs/config.go` 中的端口配置。

### Q: 为什么我的请求返回 401？

检查是否添加了 `Authorization` header。

### Q: 如何重置数据？

重启应用即可，数据存储在内存中。

### Q: Worker 后台任务在哪里查看？

查看运行 Worker 的终端输出。

## 📚 共享服务架构

本项目展示了如何使用共享服务：

```
功能切片 (Feature)
    │
    ├─→ Repository (数据访问)
    │   └─→ Database
    │
    └─→ Shared Services (业务服务)
        ├─→ NotificationService
        │   ├─→ 邮件服务
        │   ├─→ 短信服务
        │   └─→ 推送服务
        │
        └─→ PaymentService
            ├─→ 支付宝
            ├─→ 微信支付
            └─→ 银行卡支付
```

**共享服务的优势：**
- ✅ **跨端复用**：管理端和C端共享同一个服务
- ✅ **易于测试**：可以创建 Mock 实现
- ✅ **易于替换**：更换支付网关不影响业务代码
- ✅ **统一管理**：所有第三方集成在一处

查看详细文档：
- [共享服务文档](shared/services/README.md)
- [仓储层文档](shared/repositories/README.md)

## 🎯 下一步

### 实际项目改造

1. **数据库集成**
   - 替换内存存储为 PostgreSQL/MySQL
   - 使用 GORM 或 sqlx
   - 添加数据库迁移

2. **真实的第三方服务**
   - 集成阿里云邮件推送
   - 集成腾讯云短信
   - 集成支付宝/微信支付

3. **完善的认证**
   - 实现真实的 JWT 认证
   - 添加 Refresh Token
   - 实现权限控制

4. **测试**
   - 添加单元测试
   - 添加集成测试
   - Mock 共享服务

5. **监控和日志**
   - 添加日志框架（zap、logrus）
   - 添加监控（Prometheus）
   - 添加链路追踪（Jaeger）

6. **API 文档**
   - 集成 Swagger
   - 生成 API 文档
   - 添加请求示例

---

**祝你使用愉快！** 🎉

