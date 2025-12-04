# 共享服务使用指南

## 📦 概述

本指南展示了如何在垂直切片架构中使用共享服务（Shared Services）。共享服务是可以在多个端（管理端、C端、Worker）之间复用的业务服务。

## 🎯 什么是共享服务？

共享服务是封装了通用业务逻辑或第三方系统集成的服务，它们：

- ✅ **可跨端复用**：在管理端和C端都可以使用
- ✅ **封装第三方 API**：如邮件、短信、支付等
- ✅ **独立测试**：可以单独测试和 Mock
- ✅ **统一管理**：集中管理第三方服务的配置和调用

## 📂 已实现的共享服务

### 1. NotificationService - 通知服务

**位置：** `shared/services/notification/`

**功能：**
- 发送邮件
- 发送短信
- 发送推送通知

**接口定义：**
```go
type INotificationService interface {
    SendEmail(to, subject, body string) error
    SendSMS(phone, message string) error
    SendPush(userID int64, title, message string) error
}
```

**使用场景：**
- 用户注册后发送欢迎邮件
- 订单创建后发送确认通知
- 支付成功后发送推送

### 2. PaymentService - 支付服务

**位置：** `shared/services/payment/`

**功能：**
- 创建支付
- 查询支付状态
- 退款处理

**接口定义：**
```go
type IPaymentService interface {
    CreatePayment(orderID int64, amount float64, method PaymentMethod) (*PaymentResult, error)
    QueryPayment(paymentID string) (*PaymentStatus, error)
    RefundPayment(paymentID string, amount float64, reason string) error
}
```

**使用场景：**
- 订单支付
- 支付状态查询
- 订单退款

## 🏗️ 架构设计

```
┌─────────────────────────────────────────────────────────┐
│                   Apps (各端业务)                        │
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │    Admin     │  │     API      │  │    Worker    │ │
│  │              │  │              │  │              │ │
│  │  Features:   │  │  Features:   │  │    Jobs:     │ │
│  │  - Users     │  │  - Auth      │  │  - OrderSync │ │
│  │  - Products  │  │  - Orders    │  │  - EmailSend │ │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘ │
│         │                  │                  │         │
└─────────┼──────────────────┼──────────────────┼─────────┘
          │                  │                  │
          └──────────────────┼──────────────────┘
                             ▼
          ┌────────────────────────────────────┐
          │       Shared Services              │
          │                                    │
          │  ┌──────────────────────────────┐ │
          │  │  NotificationService         │ │
          │  │  - SendEmail()               │ │
          │  │  - SendSMS()                 │ │
          │  │  - SendPush()                │ │
          │  └──────────────────────────────┘ │
          │                                    │
          │  ┌──────────────────────────────┐ │
          │  │  PaymentService              │ │
          │  │  - CreatePayment()           │ │
          │  │  - QueryPayment()            │ │
          │  │  - RefundPayment()           │ │
          │  └──────────────────────────────┘ │
          └────────────────────────────────────┘
```

## 💡 使用示例

### 示例 1：在订单支付功能中使用

**文件：** `apps/api/features/orders/pay_order.go`

```go
type PayOrderHandler struct {
    orderRepo           repositories.IOrderRepository
    paymentService      services.IPaymentService          // 共享服务
    notificationService services.INotificationService     // 共享服务
    userRepo            repositories.IUserRepository
}

func (h *PayOrderHandler) Handle(c *web.HttpContext) web.IActionResult {
    // 1. 获取订单信息
    order, _ := h.orderRepo.GetByID(orderID)
    
    // 2. 使用支付服务创建支付
    paymentResult, err := h.paymentService.CreatePayment(
        order.ID,
        order.TotalPrice,
        services.PaymentMethodAlipay,
    )
    
    // 3. 使用通知服务发送通知
    user, _ := h.userRepo.GetByID(order.UserID)
    
    // 发送邮件
    h.notificationService.SendEmail(
        user.Email,
        "订单支付确认",
        "您的订单正在支付中",
    )
    
    // 发送推送
    h.notificationService.SendPush(
        user.ID,
        "订单支付",
        "请完成支付",
    )
    
    return c.Ok(paymentResult)
}
```

### 示例 2：服务注册

**文件：** `apps/api/bootstrap.go`

```go
func Bootstrap() *web.WebApplication {
    builder := web.CreateBuilder()
    
    // 注册共享基础设施
    database.AddDatabase(builder.Services)
    cache.AddCache(builder.Services)
    
    // 注册共享仓储
    repositories.AddRepositories(builder.Services)
    
    // ✅ 注册共享服务
    notification.AddNotificationService(builder.Services)
    payment.AddPaymentService(builder.Services)
    
    // 注册功能切片
    orders.AddOrderFeature(builder.Services)
    
    return builder.Build()
}
```

### 示例 3：功能切片注册

**文件：** `apps/api/features/orders/service_extensions.go`

```go
func AddOrderFeature(services di.IServiceCollection) {
    // 注册处理器（会自动注入共享服务）
    services.AddSingleton(NewCreateOrderHandler)
    services.AddSingleton(NewPayOrderHandler)  // 使用共享服务
    
    // 注册控制器
    web.AddController(services, NewOrderController)
}
```

## 🧪 测试示例

### API 调用测试

```bash
# 1. 创建订单
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: token_xxx" \
  -d '{
    "items": [{"product_id": 1, "quantity": 2}]
  }'

# 2. 支付订单（触发共享服务）
curl -X POST http://localhost:8080/api/orders/1/pay \
  -H "Content-Type: application/json" \
  -H "Authorization: token_xxx" \
  -d '{
    "payment_method": "alipay"
  }'
```

### 预期的后台日志

```
[PAYMENT] 创建支付
  支付ID: PAY_1_1704096000
  订单ID: 1
  金额: ¥15998.00
  支付方式: alipay
  时间: 2024-01-01 12:00:00

[EMAIL] 发送邮件
  收件人: user@example.com
  主题: 订单支付确认
  内容: 您的订单正在支付中
  时间: 2024-01-01 12:00:00

[PUSH] 发送推送通知
  用户ID: 2
  标题: 订单支付
  内容: 请完成支付
  时间: 2024-01-01 12:00:00
```

## 📝 创建新的共享服务

### 步骤 1：定义接口

在 `shared/contracts/services/` 中定义接口：

```go
// shared/contracts/services/sms_service.go
package services

type ISMSService interface {
    SendVerificationCode(phone string) (code string, err error)
    VerifyCode(phone string, code string) bool
}
```

### 步骤 2：实现服务

在 `shared/services/` 中实现：

```go
// shared/services/sms/sms_service.go
package sms

type SMSService struct {
    // 配置
}

func NewSMSService() services.ISMSService {
    return &SMSService{}
}

func (s *SMSService) SendVerificationCode(phone string) (string, error) {
    // 实现逻辑
    return "123456", nil
}
```

### 步骤 3：创建 DI 注册

```go
// shared/services/sms/service_extensions.go
package sms

func AddSMSService(services di.IServiceCollection) {
    services.AddSingleton(NewSMSService)
}
```

### 步骤 4：在 Bootstrap 中注册

```go
// apps/api/bootstrap.go
func Bootstrap() *web.WebApplication {
    builder := web.CreateBuilder()
    
    // 注册新的共享服务
    sms.AddSMSService(builder.Services)
    
    return builder.Build()
}
```

### 步骤 5：在功能中使用

```go
type SendCodeHandler struct {
    smsService services.ISMSService
}

func (h *SendCodeHandler) Handle(c *web.HttpContext) web.IActionResult {
    code, err := h.smsService.SendVerificationCode(phone)
    // ...
}
```

## 🔄 实际项目改造

### 替换为真实的第三方服务

#### 1. 邮件服务（使用阿里云邮件推送）

```go
import "github.com/aliyun/alibaba-cloud-sdk-go/services/dm"

type AliyunEmailService struct {
    client *dm.Client
}

func (s *AliyunEmailService) SendEmail(to, subject, body string) error {
    request := dm.CreateSingleSendMailRequest()
    request.AccountName = s.config.AccountName
    request.AddressType = "1"
    request.ReplyToAddress = "false"
    request.ToAddress = to
    request.Subject = subject
    request.HtmlBody = body
    
    _, err := s.client.SingleSendMail(request)
    return err
}
```

#### 2. 支付服务（使用支付宝）

```go
import "github.com/smartwalle/alipay/v3"

type AlipayService struct {
    client *alipay.Client
}

func (s *AlipayService) CreatePayment(...) (*PaymentResult, error) {
    // 调用支付宝 API
    p := alipay.TradePagePay{}
    p.Subject = "订单支付"
    p.OutTradeNo = fmt.Sprintf("%d", orderID)
    p.TotalAmount = fmt.Sprintf("%.2f", amount)
    p.ProductCode = "FAST_INSTANT_TRADE_PAY"
    
    url, err := s.client.TradePagePay(p)
    // ...
}
```

## 📊 目录结构

```
shared/
├── contracts/              # 接口定义
│   ├── repositories/      # 仓储接口
│   └── services/          # ✅ 服务接口
│       ├── notification_service.go
│       └── payment_service.go
│
├── services/              # ✅ 服务实现
│   ├── notification/
│   │   ├── notification_service.go
│   │   └── service_extensions.go
│   ├── payment/
│   │   ├── payment_service.go
│   │   └── service_extensions.go
│   └── README.md
│
├── repositories/          # 仓储实现
├── domain/               # 领域模型
└── infrastructure/       # 基础设施
```

## 🎯 最佳实践

### 1. 接口优先

始终先定义接口，再实现服务。

### 2. 单一职责

每个服务只负责一类功能。

### 3. 依赖注入

通过构造函数注入依赖。

### 4. 错误处理

明确的错误返回和处理。

### 5. 可测试性

便于创建 Mock 实现进行测试。

## 📚 相关文档

- [共享服务详细文档](shared/services/README.md)
- [仓储层文档](shared/repositories/README.md)
- [API 调用示例](EXAMPLES.md)
- [架构设计文档](ARCHITECTURE.md)

## ✅ 总结

共享服务的使用让你的代码：

- ✅ **更易维护**：集中管理第三方服务
- ✅ **更易测试**：可以 Mock 共享服务
- ✅ **更易扩展**：新增服务不影响现有代码
- ✅ **更易复用**：在多个端之间共享逻辑

---

**使用共享服务，让你的架构更清晰！** 🚀

