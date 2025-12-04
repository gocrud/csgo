# 共享服务（Shared Services）

## 📦 概述

共享服务层包含可以在多个端（管理端、C端、Worker）之间复用的业务服务。这些服务通常封装了与第三方系统的交互或通用的业务逻辑。

## 🎯 何时创建共享服务？

创建共享服务的时机：

- ✅ **需要跨端复用**：管理端和C端都需要使用
- ✅ **第三方集成**：邮件、短信、支付、存储等
- ✅ **通用业务逻辑**：不属于特定端的业务功能
- ✅ **独立测试**：可以独立测试的服务

不要创建共享服务的情况：

- ❌ **端特定逻辑**：只在一个端使用的逻辑应该在对应端实现
- ❌ **简单工具函数**：应该放在 `shared/utils` 中
- ❌ **数据访问**：应该使用 Repository 模式

## 📂 当前的共享服务

### 1. NotificationService - 通知服务

负责发送各种通知（邮件、短信、推送）。

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
- 订单创建后发送确认短信
- 重要操作后发送推送通知

**使用示例：**

```go
// 在功能切片的 Handler 中使用
type RegisterHandler struct {
    userRepo            repositories.IUserRepository
    notificationService services.INotificationService
}

func (h *RegisterHandler) Handle(c *web.HttpContext) web.IActionResult {
    // ... 创建用户 ...
    
    // 发送欢迎邮件
    err := h.notificationService.SendEmail(
        user.Email,
        "欢迎注册",
        "欢迎加入我们的平台！",
    )
    
    return c.Created(user)
}
```

### 2. PaymentService - 支付服务

负责处理支付相关的操作（创建支付、查询状态、退款）。

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
- 查询支付状态
- 订单退款

**使用示例：**

```go
// 在订单支付功能中使用
type PayOrderHandler struct {
    orderRepo      repositories.IOrderRepository
    paymentService services.IPaymentService
}

func (h *PayOrderHandler) Handle(c *web.HttpContext) web.IActionResult {
    // ... 查询订单 ...
    
    // 创建支付
    result, err := h.paymentService.CreatePayment(
        order.ID,
        order.TotalPrice,
        services.PaymentMethodAlipay,
    )
    
    return c.Ok(result)
}
```

## 🏗️ 创建新的共享服务

### 步骤 1：定义接口

在 `shared/contracts/services/` 中定义服务接口：

```go
// shared/contracts/services/sms_service.go
package services

type ISMSService interface {
    SendVerificationCode(phone string) (code string, err error)
    VerifyCode(phone string, code string) bool
}
```

### 步骤 2：实现服务

在 `shared/services/` 中创建服务实现：

```go
// shared/services/sms/sms_service.go
package sms

import "vertical_slice_demo/shared/contracts/services"

type SMSService struct {
    // 配置...
}

func NewSMSService() services.ISMSService {
    return &SMSService{}
}

func (s *SMSService) SendVerificationCode(phone string) (string, error) {
    // 实现...
    return "123456", nil
}

func (s *SMSService) VerifyCode(phone string, code string) bool {
    // 实现...
    return true
}
```

### 步骤 3：注册服务

创建 DI 注册文件：

```go
// shared/services/sms/service_extensions.go
package sms

import "github.com/gocrud/csgo/di"

func AddSMSService(services di.IServiceCollection) {
    services.AddSingleton(NewSMSService)
}
```

### 步骤 4：在 Bootstrap 中注册

```go
// apps/api/bootstrap.go
func Bootstrap() *web.WebApplication {
    builder := web.CreateBuilder()
    
    // 注册共享服务
    sms.AddSMSService(builder.Services)
    
    // ...
}
```

### 步骤 5：在功能切片中使用

```go
type SendCodeHandler struct {
    smsService services.ISMSService
}

func (h *SendCodeHandler) Handle(c *web.HttpContext) web.IActionResult {
    code, err := h.smsService.SendVerificationCode(phone)
    // ...
}
```

## 📝 最佳实践

### 1. 接口优先

始终定义接口，便于测试和替换实现：

```go
// ✅ 好的做法
type IEmailService interface {
    Send(to, subject, body string) error
}

// ❌ 不好的做法 - 直接依赖具体实现
type EmailService struct { }
```

### 2. 依赖注入

通过构造函数注入依赖：

```go
// ✅ 好的做法
type Handler struct {
    emailService services.IEmailService
}

func NewHandler(emailService services.IEmailService) *Handler {
    return &Handler{emailService: emailService}
}
```

### 3. 单一职责

每个服务只负责一类功能：

```go
// ✅ 好的做法
type IEmailService interface {
    Send(to, subject, body string) error
}

type ISMSService interface {
    Send(phone, message string) error
}

// ❌ 不好的做法 - 职责混乱
type INotificationService interface {
    SendEmail(to, subject, body string) error
    SendSMS(phone, message string) error
    ProcessPayment(amount float64) error // 不应该在这里
}
```

### 4. 错误处理

明确的错误返回：

```go
// ✅ 好的做法
func (s *PaymentService) CreatePayment(...) (*PaymentResult, error) {
    if amount <= 0 {
        return nil, errors.New("invalid amount")
    }
    // ...
}
```

### 5. 配置管理

通过配置对象管理服务配置：

```go
type EmailConfig struct {
    SMTPHost string
    SMTPPort int
    Username string
    Password string
}

type EmailService struct {
    config *EmailConfig
}

func NewEmailService(config *EmailConfig) *EmailService {
    return &EmailService{config: config}
}
```

## 🧪 测试共享服务

### 单元测试

```go
func TestNotificationService_SendEmail(t *testing.T) {
    service := NewNotificationService()
    
    err := service.SendEmail("test@example.com", "Test", "Body")
    
    assert.NoError(t, err)
}
```

### Mock 实现

```go
type MockNotificationService struct {
    SendEmailCalled bool
}

func (m *MockNotificationService) SendEmail(to, subject, body string) error {
    m.SendEmailCalled = true
    return nil
}
```

## 📊 架构图

```
┌─────────────────────────────────────┐
│         功能切片 (Feature)          │
│                                     │
│  ┌─────────────────────────────┐  │
│  │  Handler (业务逻辑)         │  │
│  │                             │  │
│  │  - 使用 Repository         │  │
│  │  - 使用 Shared Services    │  │
│  └─────────────────────────────┘  │
└─────────────────┬───────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│       Shared Services               │
│                                     │
│  ┌────────────┐  ┌──────────────┐ │
│  │Notification│  │   Payment    │ │
│  │  Service   │  │   Service    │ │
│  └────────────┘  └──────────────┘ │
│                                     │
│  - 封装第三方 API                  │
│  - 通用业务逻辑                    │
│  - 可跨端复用                      │
└─────────────────────────────────────┘
```

## 🚀 实际项目集成

在实际项目中，你需要：

1. **替换模拟实现**：用真实的 API 调用替换打印语句
2. **添加配置**：从配置文件读取 API Key、密钥等
3. **错误处理**：完善的错误处理和重试机制
4. **日志记录**：记录关键操作日志
5. **监控告警**：监控服务调用情况

## 📚 相关文档

- [依赖注入指南](../../../docs/guides/dependency-injection.md)
- [Repository 模式](../repositories/README.md)
- [功能切片示例](../../apps/api/features/orders/pay_order.go)

---

**共享服务让你的代码更易维护和测试！** 🎉

