# 错误处理完整示例

本文档提供了完整的错误处理使用示例，涵盖从定义错误、Service 层到 Controller 的完整流程。

## 目录

- [快速入门](#快速入门)
- [Service 层错误处理](#service-层错误处理)
- [Controller 层错误处理](#controller-层错误处理)
- [复杂业务场景](#复杂业务场景)
- [错误链和追踪](#错误链和追踪)

## 快速入门

### 1. 定义模块错误

```go
// services/errors.go
package services

import "github.com/gocrud/csgo/errors"

// 定义所有业务模块的错误
var (
    UserErrors  = errors.NewModule("USER")
    OrderErrors = errors.NewModule("ORDER")
    DramaErrors = errors.NewModule("DRAMA")
)
```

### 2. 最简单的使用

```go
// 资源不存在
return UserErrors.NotFound("用户不存在")

// 资源已存在
return UserErrors.AlreadyExists("用户名已被占用")

// 参数无效
return OrderErrors.InvalidParam("金额必须大于0")

// 权限不足
return DramaErrors.PermissionDenied("无权访问此剧集")
```

### 3. 自定义错误码

```go
// 方式1：Code().Msg()
return OrderErrors.Code("PAYMENT_FAILED").Msg("支付失败")

// 方式2：Code().Msgf() - 格式化消息
return OrderErrors.Code("PAYMENT_FAILED").Msgf("余额不足: %.2f", balance)

// 方式3：带详细信息
return OrderErrors.Code("PAYMENT_FAILED").
    Msg("支付失败").
    WithDetail("orderId", orderID).
    WithDetail("amount", amount)
```

## Service 层错误处理

### 场景1：用户管理 CRUD

```go
package services

import "github.com/gocrud/csgo/errors"

type UserService struct {
    repo UserRepository
}

// GetUser 查询用户
func (s *UserService) GetUser(id int) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        // 包装数据库错误
        return nil, UserErrors.Internal("查询用户失败").Wrap(err)
    }
    
    if user == nil {
        // 资源不存在
        return nil, UserErrors.NotFound("用户不存在").
            WithDetail("userId", id)
    }
    
    return user, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(req *CreateUserRequest) (*User, error) {
    // 检查邮箱是否已存在
    if s.repo.ExistsByEmail(req.Email) {
        return nil, UserErrors.AlreadyExists("邮箱已被注册").
            WithDetail("email", req.Email)
    }
    
    // 检查用户名是否已存在
    if s.repo.ExistsByUsername(req.Username) {
        return nil, UserErrors.AlreadyExists("用户名已被使用").
            WithDetail("username", req.Username)
    }
    
    user := &User{
        Email:    req.Email,
        Username: req.Username,
        Password: hashPassword(req.Password),
    }
    
    if err := s.repo.Create(user); err != nil {
        return nil, UserErrors.Internal("创建用户失败").Wrap(err)
    }
    
    return user, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(id int, req *UpdateUserRequest) error {
    user, err := s.GetUser(id)
    if err != nil {
        return err  // 直接返回 GetUser 的错误
    }
    
    // 更新字段...
    if err := s.repo.Update(user); err != nil {
        return UserErrors.OperationFailed("更新用户失败").Wrap(err)
    }
    
    return nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id int) error {
    user, err := s.GetUser(id)
    if err != nil {
        return err
    }
    
    // 检查是否有关联订单
    if s.orderRepo.CountByUserID(id) > 0 {
        return UserErrors.OperationFailed("该用户还有未完成的订单，无法删除").
            WithDetail("userId", id)
    }
    
    if err := s.repo.Delete(id); err != nil {
        return UserErrors.Internal("删除用户失败").Wrap(err)
    }
    
    return nil
}
```

### 场景2：订单支付流程

```go
package services

type OrderService struct {
    repo           OrderRepository
    accountService *AccountService
    paymentGateway PaymentGateway
}

// ProcessPayment 处理支付
func (s *OrderService) ProcessPayment(orderID string, amount float64) error {
    // 1. 查询订单
    order, err := s.repo.FindByID(orderID)
    if err != nil {
        return OrderErrors.NotFound("订单不存在").
            WithDetail("orderId", orderID).
            Wrap(err)
    }
    
    // 2. 检查订单状态
    if order.Status == "cancelled" {
        return OrderErrors.InvalidStatus("订单已取消").
            WithDetail("orderId", orderID).
            WithDetail("status", order.Status)
    }
    
    if order.Status == "paid" {
        return OrderErrors.AlreadyExists("订单已支付")
    }
    
    if order.IsExpired() {
        return OrderErrors.Expired("订单已过期").
            WithDetail("expireTime", order.ExpireTime)
    }
    
    // 3. 检查账户余额
    balance, err := s.accountService.GetBalance(order.UserID)
    if err != nil {
        return OrderErrors.Internal("获取账户余额失败").Wrap(err)
    }
    
    if balance < amount {
        return OrderErrors.Code("PAYMENT_FAILED").
            Msgf("余额不足，当前: ¥%.2f，需要: ¥%.2f", balance, amount).
            WithDetail("balance", balance).
            WithDetail("required", amount).
            WithDetail("shortfall", amount-balance)
    }
    
    // 4. 调用支付网关
    if err := s.paymentGateway.Charge(amount); err != nil {
        // 根据不同的支付错误返回不同消息
        switch err.(type) {
        case *NetworkError:
            return OrderErrors.Code("PAYMENT_FAILED").
                Msg("网络异常，请稍后重试").
                WithDetail("retryable", true).
                Wrap(err)
                
        case *ChannelError:
            return OrderErrors.Code("PAYMENT_FAILED").
                Msg("支付渠道暂时不可用").
                WithDetail("channel", "alipay").
                Wrap(err)
                
        default:
            return OrderErrors.Code("PAYMENT_FAILED").
                Msg("支付失败").
                Wrap(err)
        }
    }
    
    return nil
}

// RefundOrder 退款
func (s *OrderService) RefundOrder(orderID string, reason string) error {
    order, err := s.repo.FindByID(orderID)
    if err != nil {
        return OrderErrors.NotFound("订单不存在")
    }
    
    // 检查退款条件
    if order.Status != "paid" {
        return OrderErrors.InvalidStatus("只有已支付的订单才能退款").
            WithDetail("currentStatus", order.Status)
    }
    
    // 检查退款时限
    if order.CreatedAt.AddDate(0, 0, 7).Before(time.Now()) {
        return OrderErrors.Code("REFUND_EXPIRED").
            Msg("超过7天退款期限").
            WithDetail("orderTime", order.CreatedAt)
    }
    
    // 执行退款...
    return nil
}
```

### 场景3：动态修改错误消息

```go
// ValidateOrder 验证订单
func (s *OrderService) ValidateOrder(order *Order) error {
    // 基础验证错误
    baseErr := OrderErrors.InvalidParam()
    
    // 根据不同情况返回不同消息
    if order.Amount <= 0 {
        return baseErr.
            WithMsg("订单金额必须大于0").
            WithDetail("amount", order.Amount)
    }
    
    if len(order.Items) == 0 {
        return baseErr.WithMsg("订单商品不能为空")
    }
    
    if order.UserID == "" {
        return baseErr.
            WithMsg("用户ID不能为空").
            WithDetail("field", "userId")
    }
    
    // 复杂验证
    totalAmount := order.CalculateTotalAmount()
    if totalAmount != order.Amount {
        return baseErr.
            WithMsgf("订单金额不匹配，计算总额: %.2f，实际金额: %.2f", 
                totalAmount, order.Amount).
            WithDetail("calculated", totalAmount).
            WithDetail("actual", order.Amount)
    }
    
    return nil
}
```

## Controller 层错误处理

### 推荐方式：使用 FromError

```go
package controllers

import (
    "github.com/gocrud/csgo/web"
    "your-project/services"
)

type UserController struct {
    service *services.UserService
}

// GetUser 获取用户（最简洁）
func (c *UserController) GetUser(ctx *web.HttpContext) web.IActionResult {
    id := ctx.Params().PathInt("id").Value()
    
    user, err := c.service.GetUser(id)
    if err != nil {
        // FromError 自动处理所有错误类型
        return ctx.FromError(err)
    }
    
    return ctx.Ok(user)
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *web.HttpContext) web.IActionResult {
    req, result := web.BindAndValidate[CreateUserRequest](ctx)
    if result != nil {
        return result  // 验证失败自动返回
    }
    
    user, err := c.service.CreateUser(req)
    if err != nil {
        return ctx.FromError(err)  // 一行搞定！
    }
    
    return ctx.Created(user)
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *web.HttpContext) web.IActionResult {
    id := ctx.Params().PathInt("id").Value()
    
    if err := c.service.DeleteUser(id); err != nil {
        return ctx.FromError(err)
    }
    
    return ctx.NoContent()
}
```

### 响应示例

```json
// 成功响应 (200)
{
  "success": true,
  "data": {
    "id": 1,
    "username": "john",
    "email": "john@example.com"
  }
}

// 错误响应 - 资源不存在 (404)
{
  "success": false,
  "error": {
    "code": "USER.NOT_FOUND",
    "message": "用户不存在",
    "details": {
      "userId": 123
    }
  }
}

// 错误响应 - 资源已存在 (409)
{
  "success": false,
  "error": {
    "code": "USER.ALREADY_EXISTS",
    "message": "邮箱已被注册",
    "details": {
      "email": "test@example.com"
    }
  }
}

// 错误响应 - 支付失败 (400)
{
  "success": false,
  "error": {
    "code": "ORDER.PAYMENT_FAILED",
    "message": "余额不足，当前: ¥50.00，需要: ¥100.00",
    "details": {
      "balance": 50.00,
      "required": 100.00,
      "shortfall": 50.00
}
```

## 复杂业务场景

### 场景1：订单结算流程

展示多层错误包装和错误链追踪。

```go
package services

type CheckoutService struct {
    orderService   *OrderService
    inventoryService *InventoryService
    paymentService *PaymentService
}

func (s *CheckoutService) ProcessCheckout(req *CheckoutRequest) (*Order, error) {
    // 1. 验证库存
    if err := s.inventoryService.CheckStock(req.Items); err != nil {
        return nil, OrderErrors.OperationFailed("结算失败").
            WithMsg("库存检查失败").
            Wrap(err)  // 包装库存错误
    }
    
    // 2. 创建订单
    order, err := s.orderService.CreateOrder(req)
    if err != nil {
        return nil, OrderErrors.OperationFailed("结算失败").
            WithMsg("创建订单失败").
            Wrap(err)
    }
    
    // 3. 处理支付
    if err := s.paymentService.ProcessPayment(order.ID, req.Amount); err != nil {
        // 支付失败，回滚订单
        s.orderService.CancelOrder(order.ID, "支付失败")
        
        return nil, OrderErrors.Code("CHECKOUT_FAILED").
            Msg("结算失败").
            Wrap(err)  // 包装支付错误
    }
    
    // 4. 扣减库存
    if err := s.inventoryService.DeductStock(req.Items); err != nil {
        // 库存扣减失败，需要退款
        s.paymentService.Refund(order.ID)
        s.orderService.CancelOrder(order.ID, "库存扣减失败")
        
        return nil, OrderErrors.OperationFailed("结算失败").
            WithMsg("库存扣减失败").
            Wrap(err)
    }
    
    return order, nil
}

// 库存服务错误
type InventoryService struct {
    repo InventoryRepository
}

var InventoryErrors = errors.NewModule("INVENTORY")

func (s *InventoryService) CheckStock(items []OrderItem) error {
    for _, item := range items {
        stock, err := s.repo.GetStock(item.ProductID)
        if err != nil {
            return InventoryErrors.Internal("查询库存失败").
                WithDetail("productId", item.ProductID).
                Wrap(err)
        }
        
        if stock < item.Quantity {
            return InventoryErrors.Code("INSUFFICIENT_STOCK").
                Msgf("商品库存不足: %s", item.ProductName).
                WithDetail("productId", item.ProductID).
                WithDetail("required", item.Quantity).
                WithDetail("available", stock)
        }
    }
    return nil
}
```

### 场景2：权限检查和多级错误

```go
package services

var AuthErrors = errors.NewModule("AUTH")

type PermissionService struct {
    userService *UserService
    roleService *RoleService
}

func (s *PermissionService) CheckPermission(userID int, resource string, action string) error {
    // 1. 获取用户
    user, err := s.userService.GetUser(userID)
    if err != nil {
        return AuthErrors.Internal("权限检查失败").
            WithMsg("获取用户信息失败").
            Wrap(err)
    }
    
    // 2. 检查用户状态
    if user.Status == "suspended" {
        return AuthErrors.Code("USER_SUSPENDED").
            Msg("用户已被停用").
            WithDetail("userId", userID).
            WithDetail("suspendReason", user.SuspendReason)
    }
    
    if user.Status == "deleted" {
        return AuthErrors.Unauthorized("用户不存在")
    }
    
    // 3. 获取用户角色
    roles, err := s.roleService.GetUserRoles(userID)
    if err != nil {
        return AuthErrors.Internal("权限检查失败").
            WithMsg("获取用户角色失败").
            Wrap(err)
    }
    
    // 4. 检查权限
    hasPermission := false
    for _, role := range roles {
        if role.HasPermission(resource, action) {
            hasPermission = true
            break
        }
    }
    
    if !hasPermission {
        return AuthErrors.PermissionDenied("无权执行此操作").
            WithDetail("resource", resource).
            WithDetail("action", action).
            WithDetail("userId", userID)
    }
    
    return nil
}
```

### 场景3：外部 API 调用错误处理

```go
package services

type PaymentGateway struct {
    client *http.Client
}

var PaymentErrors = errors.NewModule("PAYMENT")

func (g *PaymentGateway) Charge(orderID string, amount float64) error {
    resp, err := g.client.Post("/charge", map[string]any{
        "order_id": orderID,
        "amount":   amount,
    })
    
    if err != nil {
        // 网络错误
        return PaymentErrors.Code("GATEWAY_ERROR").
            Msg("支付网关连接失败").
            WithDetail("retryable", true).
            WithHTTPCode(503).  // Service Unavailable
            Wrap(err)
    }
    
    if resp.StatusCode == 400 {
        // 参数错误
        return PaymentErrors.InvalidParam("支付参数错误").
            WithDetail("orderId", orderID)
    }
    
    if resp.StatusCode == 402 {
        // 余额不足
        return PaymentErrors.Code("INSUFFICIENT_BALANCE").
            Msg("账户余额不足").
            WithHTTPCode(402)
    }
    
    if resp.StatusCode >= 500 {
        // 网关错误
        return PaymentErrors.Code("GATEWAY_ERROR").
            Msg("支付网关异常").
            WithDetail("retryable", true).
            WithHTTPCode(503)
    }
    
    return nil
}
```

## 错误链和追踪

### 使用 errors.Is 和 errors.As

```go
package services

import (
    "errors"
    "github.com/gocrud/csgo/errors"
)

func (s *OrderService) HandleCheckoutError(err error) {
    // 1. 使用 errors.Is 检查错误类型
    var csgoErr *errors.Error
    if errors.As(err, &csgoErr) {
        log.Printf("业务错误: code=%s, message=%s", 
            csgoErr.Code(), csgoErr.Message())
        
        // 获取详细信息
        if details := csgoErr.Details(); len(details) > 0 {
            log.Printf("错误详情: %+v", details)
        }
        
        // 检查是否可重试
        if retryable, ok := details["retryable"].(bool); ok && retryable {
            // 可以重试
            log.Println("错误可重试")
        }
    }
    
    // 2. 获取原始错误
    if cause := errors.Unwrap(err); cause != nil {
        log.Printf("原始错误: %v", cause)
    }
    
    // 3. 根据错误码进行处理
    if csgoErr != nil {
        switch {
        case strings.Contains(csgoErr.Code(), "INSUFFICIENT_STOCK"):
            // 处理库存不足
            s.handleInsufficientStock(csgoErr)
            
        case strings.Contains(csgoErr.Code(), "PAYMENT_FAILED"):
            // 处理支付失败
            s.handlePaymentFailed(csgoErr)
            
        default:
            // 其他错误
            log.Printf("未处理的错误: %s", csgoErr.Code())
        }
    }
}
```

### 错误传播示例

```go
// 底层 Repository
func (r *OrderRepository) FindByID(id string) (*Order, error) {
    var order Order
    err := r.db.First(&order, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil  // 返回 nil，在上层处理
        }
        return nil, err  // 返回原始数据库错误
    }
    return &order, nil
}

// Service 层包装错误
func (s *OrderService) GetOrder(id string) (*Order, error) {
    order, err := s.repo.FindByID(id)
    if err != nil {
        // 包装数据库错误
        return nil, OrderErrors.Internal("查询订单失败").
            WithDetail("orderId", id).
            Wrap(err)
    }
    
    if order == nil {
        // 转换为业务错误
        return nil, OrderErrors.NotFound("订单不存在").
            WithDetail("orderId", id)
    }
    
    return order, nil
}

// Controller 层统一处理
func (c *OrderController) GetOrder(ctx *web.HttpContext) web.IActionResult {
    id := ctx.Params().String("id").Value()
    
    order, err := c.service.GetOrder(id)
    if err != nil {
        // 自动处理所有错误类型
        return ctx.FromError(err)
    }
    
    return ctx.Ok(order)
}
```

## 最佳实践总结

### 1. 错误定义

```go
// ✅ 推荐：统一定义模块错误
var (
    UserErrors  = errors.NewModule("USER")
    OrderErrors = errors.NewModule("ORDER")
)

// ❌ 不推荐：每次都创建
func GetUser(id int) error {
    return errors.NewModule("USER").NotFound()  // 性能差
}
```

### 2. 错误消息

```go
// ✅ 推荐：提供清晰的错误消息
return UserErrors.NotFound("用户不存在").
    WithDetail("userId", id)

// ❌ 不推荐：模糊的消息
return UserErrors.NotFound("Not found")
```

### 3. 错误包装

```go
// ✅ 推荐：包装底层错误
if err != nil {
    return UserErrors.Internal("创建用户失败").Wrap(err)
}

// ❌ 不推荐：丢失原始错误
if err != nil {
    return UserErrors.Internal("创建用户失败")  // 丢失了原始错误链
}
```

### 4. 详细信息

```go
// ✅ 推荐：添加有用的详细信息
return OrderErrors.Code("PAYMENT_FAILED").
    Msg("余额不足").
    WithDetail("balance", balance).
    WithDetail("required", amount).
    WithDetail("retryable", false)

// ❌ 不推荐：没有详细信息
return OrderErrors.Code("PAYMENT_FAILED").Msg("余额不足")
```

### 5. HTTP 状态码

```go
// ✅ 推荐：使用合适的快捷方法（自动映射状态码）
return UserErrors.NotFound("用户不存在")  // 自动 404

// ⚠️ 特殊场景：手动指定状态码
return UserErrors.NotFound("用户不存在").
    WithHTTPCode(410)  // Gone，表示永久删除

// ❌ 不推荐：所有业务错误都用 400
return errors.New("USER.NOT_FOUND", "用户不存在", 400)  // 应该用 404
```

---

更多信息请查看 [README.md](./README.md)

完整的用户管理示例，展示基本的 CRUD 操作和错误处理。

### 1.1 定义模型和请求

```go
package models

type User struct {
    ID        int    `json:"id"`
    Email     string `json:"email"`
    Username  string `json:"username"`
    Password  string `json:"-"`
    Age       int    `json:"age"`
    Status    string `json:"status"` // active, suspended, deleted
}

type CreateUserRequest struct {
    Email    string `json:"email"`
    Username string `json:"username"`
    Password string `json:"password"`
    Age      int    `json:"age"`
}

type UpdateUserRequest struct {
    Username string `json:"username"`
    Age      int    `json:"age"`
}
```

### 1.2 定义验证器

```go
package validators

import "github.com/gocrud/csgo/validation"

// 快速失败模式验证器（默认，推荐）
func NewCreateUserValidator() *validation.AbstractValidator[CreateUserRequest] {
    v := validation.NewValidator[CreateUserRequest]()
    
    v.Field(func(r *CreateUserRequest) string { return r.Email }).
        NotEmpty().
        EmailAddress()
    
    v.Field(func(r *CreateUserRequest) string { return r.Username }).
        NotEmpty().
        MinLength(3).
        MaxLength(20)
    
    v.Field(func(r *CreateUserRequest) string { return r.Password }).
        NotEmpty().
        MinLength(6).
        Matches(`^[a-zA-Z0-9@#$%^&+=]*$`).
        WithMessage("密码只能包含字母、数字和特殊字符")
    
    v.FieldInt(func(r *CreateUserRequest) int { return r.Age }).
        GreaterThanOrEqual(18).
        LessThanOrEqual(120)
    
    return v
}

func NewUpdateUserValidator() *validation.AbstractValidator[UpdateUserRequest] {
    v := validation.NewValidator[UpdateUserRequest]()
    
    v.Field(func(r *UpdateUserRequest) string { return r.Username }).
        NotEmpty().
        MinLength(3).
        MaxLength(20)
    
    v.FieldInt(func(r *UpdateUserRequest) int { return r.Age }).
        GreaterThanOrEqual(18).
        LessThanOrEqual(120)
    
    return v
}

// 注册验证器
func init() {
    validation.RegisterValidator[CreateUserRequest](NewCreateUserValidator())
    validation.RegisterValidator[UpdateUserRequest](NewUpdateUserValidator())
}
```

### 1.3 实现 Service 层

```go
package services

import (
    "github.com/gocrud/csgo/errors"
    "myapp/models"
    "myapp/repositories"
)

type UserService struct {
    userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
    return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUser(id int) (*models.User, error) {
    user, err := s.userRepo.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    if user == nil {
        return nil, errors.Business("USER").NotFound("用户不存在")
    }
    
    return user, nil
}

func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
    // 检查邮箱是否已存在
    existing, _ := s.userRepo.FindByEmail(req.Email)
    if existing != nil {
        return nil, errors.Business("USER").AlreadyExists("邮箱已被注册")
    }
    
    // 检查用户名是否已存在
    existing, _ = s.userRepo.FindByUsername(req.Username)
    if existing != nil {
        return nil, errors.Business("USER").AlreadyExists("用户名已被使用")
    }
    
    // 创建用户
    user := &models.User{
        Email:    req.Email,
        Username: req.Username,
        Password: hashPassword(req.Password),
        Age:      req.Age,
        Status:   "active",
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, errors.Business("USER").OperationFailed("创建用户失败")
    }
    
    return user, nil
}

func (s *UserService) UpdateUser(id int, req *models.UpdateUserRequest) (*models.User, error) {
    user, err := s.GetUser(id)
    if err != nil {
        return nil, err
    }
    
    // 检查用户名是否被其他用户使用
    if req.Username != user.Username {
        existing, _ := s.userRepo.FindByUsername(req.Username)
        if existing != nil && existing.ID != id {
            return nil, errors.Business("USER").AlreadyExists("用户名已被使用")
        }
    }
    
    user.Username = req.Username
    user.Age = req.Age
    
    if err := s.userRepo.Update(user); err != nil {
        return nil, errors.Business("USER").OperationFailed("更新用户失败")
    }
    
    return user, nil
}

func (s *UserService) DeleteUser(id int) error {
    user, err := s.GetUser(id)
    if err != nil {
        return err
    }
    
    // 软删除：更新状态
    user.Status = "deleted"
    
    if err := s.userRepo.Update(user); err != nil {
        return errors.Business("USER").OperationFailed("删除用户失败")
    }
    
    return nil
}

func hashPassword(password string) string {
    // 实现密码哈希逻辑
    return password // 示例
}
```

### 1.4 实现 Controller 层

**推荐方式：使用 FromError（简洁优雅）**

```go
package controllers

import (
    "github.com/gin-gonic/gin"
    "github.com/gocrud/csgo/web"
    "myapp/models"
    "myapp/services"
)

type UserController struct {
    userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
    return &UserController{userService: userService}
}

// GET /users/:id
func (ctrl *UserController) GetUser(c *gin.Context) {
    ctx := web.NewHttpContext(c)
    
    // 获取并验证路径参数
    id, result := ctx.MustPathInt("id")
    if result != nil {
        result.ExecuteResult(c)
        return
    }
    
    // 调用 Service
    user, err := ctrl.userService.GetUser(id)
    if err != nil {
        // 使用 FromError 自动处理所有错误类型
        ctx.FromError(err, "获取用户失败").ExecuteResult(c)
        return
    }
    
    // 返回成功响应
    ctx.Ok(user).ExecuteResult(c)
}

// POST /users
func (ctrl *UserController) CreateUser(c *gin.Context) {
    ctx := web.NewHttpContext(c)
    
    // 绑定并验证请求体（自动使用注册的验证器）
    req, result := web.BindAndValidate[models.CreateUserRequest](ctx)
    if result != nil {
        result.ExecuteResult(c)
        return
    }
    
    // 调用 Service
    user, err := ctrl.userService.CreateUser(req)
    if err != nil {
        // 一行搞定所有错误处理
        ctx.FromError(err, "创建用户失败").ExecuteResult(c)
        return
    }
    
    // 返回 201 Created
    ctx.Created(user).ExecuteResult(c)
}

// PUT /users/:id
func (ctrl *UserController) UpdateUser(c *gin.Context) {
    ctx := web.NewHttpContext(c)
    
    id, result := ctx.MustPathInt("id")
    if result != nil {
        result.ExecuteResult(c)
        return
    }
    
    req, result := web.BindAndValidate[models.UpdateUserRequest](ctx)
    if result != nil {
        result.ExecuteResult(c)
        return
    }
    
    user, err := ctrl.userService.UpdateUser(id, req)
    if err != nil {
        ctx.FromError(err, "更新用户失败").ExecuteResult(c)
        return
    }
    
    ctx.Ok(user).ExecuteResult(c)
}

// DELETE /users/:id
func (ctrl *UserController) DeleteUser(c *gin.Context) {
    ctx := web.NewHttpContext(c)
    
    id, result := ctx.MustPathInt("id")
    if result != nil {
        result.ExecuteResult(c)
        return
    }
    
    err := ctrl.userService.DeleteUser(id)
    if err != nil {
        ctx.FromError(err, "删除用户失败").ExecuteResult(c)
        return
    }
    
    ctx.NoContent().ExecuteResult(c)
}
```

**传统方式：手动类型判断（仍然支持）**

```go
// 如果需要更细粒度的控制，可以使用传统方式
func (ctrl *UserController) GetUser(c *gin.Context) {
    ctx := web.NewHttpContext(c)
    
    user, err := ctrl.userService.GetUser(id)
    if err != nil {
        ctrl.handleError(ctx, err).ExecuteResult(c)
        return
    }
    
    ctx.Ok(user).ExecuteResult(c)
}

// 统一错误处理（传统方式）
func (ctrl *UserController) handleError(ctx *web.HttpContext, err error) web.IActionResult {
    if bizErr, ok := err.(*errors.BizError); ok {
        return ctx.BizError(bizErr)
    }
    return ctx.InternalError("服务器错误")
}
```

### 1.5 响应示例

#### 成功响应

```json
// GET /users/1
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "username": "johndoe",
    "age": 25,
    "status": "active"
  }
}
```

#### 验证错误响应（快速失败模式）

```json
// POST /users (邮箱格式错误)
{
  "success": false,
  "error": {
    "code": "VALIDATION.FAILED",
    "message": "验证失败",
    "fields": [
      {
        "field": "email",
        "message": "邮箱格式不正确",
        "code": "VALIDATION.EMAIL"
      }
    ]
  }
}
```

#### 业务错误响应

```json
// GET /users/999 (用户不存在)
{
  "success": false,
  "error": {
    "code": "USER.NOT_FOUND",
    "message": "用户不存在"
  }
}

// POST /users (邮箱已存在)
{
  "success": false,
  "error": {
    "code": "USER.ALREADY_EXISTS",
    "message": "邮箱已被注册"
  }
}
```

## 场景2：订单状态管理

展示业务状态验证和自定义错误码。

### 2.1 定义模型

```go
package models

type Order struct {
    ID     int    `json:"id"`
    UserID int    `json:"user_id"`
    Amount float64 `json:"amount"`
    Status string `json:"status"` // pending, paid, shipped, completed, cancelled
}

type UpdateOrderStatusRequest struct {
    Status string `json:"status"`
}
```

### 2.2 Service 层

```go
package services

import (
    "github.com/gocrud/csgo/errors"
    "myapp/models"
)

type OrderService struct {
    orderRepo repositories.OrderRepository
}

// 订单状态转换规则
var statusTransitions = map[string][]string{
    "pending":   {"paid", "cancelled"},
    "paid":      {"shipped", "cancelled"},
    "shipped":   {"completed"},
    "completed": {},
    "cancelled": {},
}

func (s *OrderService) UpdateOrderStatus(orderID int, newStatus string) (*models.Order, error) {
    order, err := s.orderRepo.FindByID(orderID)
    if err != nil {
        return nil, err
    }
    
    if order == nil {
        return nil, errors.Business("ORDER").NotFound("订单不存在")
    }
    
    // 验证状态转换是否合法
    if !s.canTransitionTo(order.Status, newStatus) {
        return nil, errors.Business("ORDER").InvalidStatus(
            fmt.Sprintf("订单状态不能从 %s 转换到 %s", order.Status, newStatus))
    }
    
    order.Status = newStatus
    
    if err := s.orderRepo.Update(order); err != nil {
        return nil, errors.Business("ORDER").OperationFailed("更新订单状态失败")
    }
    
    return order, nil
}

func (s *OrderService) canTransitionTo(currentStatus, newStatus string) bool {
    allowedStatuses, ok := statusTransitions[currentStatus]
    if !ok {
        return false
    }
    
    for _, status := range allowedStatuses {
        if status == newStatus {
            return true
        }
    }
    return false
}

func (s *OrderService) CancelOrder(orderID int, reason string) error {
    order, err := s.orderRepo.FindByID(orderID)
    if err != nil {
        return err
    }
    
    if order == nil {
        return errors.Business("ORDER").NotFound("订单不存在")
    }
    
    // 已发货的订单不能取消
    if order.Status == "shipped" || order.Status == "completed" {
        return ErrOrder.Code("CANNOT_CANCEL").
            Message("已发货或已完成的订单不能取消")
    }
    
    order.Status = "cancelled"
    // 保存取消原因...
    
    if err := s.orderRepo.Update(order); err != nil {
        return errors.Business("ORDER").OperationFailed("取消订单失败")
    }
    
    return nil
}
```

## 场景3：批量数据导入

展示全量验证模式的使用。

### 3.1 定义模型

```go
package models

type BatchImportUserRequest struct {
    Users []ImportUserItem `json:"users"`
}

type ImportUserItem struct {
    Email    string `json:"email"`
    Username string `json:"username"`
    Age      int    `json:"age"`
}
```

### 3.2 全量验证器

```go
package validators

import "github.com/gocrud/csgo/validation"

// 使用全量验证模式（收集所有错误）
func NewBatchImportValidator() *validation.AbstractValidator[BatchImportUserRequest] {
    v := validation.NewValidatorAll[BatchImportUserRequest]()  // 注意这里用 NewValidatorAll
    
    // 验证用户列表不为空
    validation.FieldSlice(v, func(r *BatchImportUserRequest) []ImportUserItem { 
        return r.Users 
    }).NotEmptySlice()
    
    return v
}

// 单个导入项验证器（也使用全量验证）
func NewImportUserItemValidator() *validation.AbstractValidator[ImportUserItem] {
    v := validation.NewValidatorAll[ImportUserItem]()
    
    v.Field(func(r *ImportUserItem) string { return r.Email }).
        NotEmpty().
        EmailAddress()
    
    v.Field(func(r *ImportUserItem) string { return r.Username }).
        NotEmpty().
        MinLength(3)
    
    v.FieldInt(func(r *ImportUserItem) int { return r.Age }).
        GreaterThanOrEqual(18)
    
    return v
}

func init() {
    validation.RegisterValidator[BatchImportUserRequest](NewBatchImportValidator())
    validation.RegisterValidator[ImportUserItem](NewImportUserItemValidator())
}
```

### 3.3 Service 层

```go
package services

import "github.com/gocrud/csgo/validation"

type ImportResult struct {
    SuccessCount int                                `json:"success_count"`
    FailedCount  int                                `json:"failed_count"`
    Errors       map[int]validation.ValidationErrors `json:"errors,omitempty"`
}

func (s *UserService) BatchImportUsers(req *BatchImportUserRequest) (*ImportResult, error) {
    result := &ImportResult{
        Errors: make(map[int]validation.ValidationErrors),
    }
    
    itemValidator := validators.NewImportUserItemValidator()
    
    for i, item := range req.Users {
        // 验证每一项（全量验证）
        validationResult := itemValidator.Validate(&item)
        if !validationResult.IsValid {
            result.FailedCount++
            result.Errors[i] = validationResult.Errors
            continue
        }
        
        // 尝试创建用户
        user := &models.User{
            Email:    item.Email,
            Username: item.Username,
            Age:      item.Age,
            Status:   "active",
        }
        
        if err := s.userRepo.Create(user); err != nil {
            result.FailedCount++
            result.Errors[i] = validation.ValidationErrors{
                {
                    Field:   "general",
                    Message: err.Error(),
                    Code:    errors.ValidationFailed,
                },
            }
            continue
        }
        
        result.SuccessCount++
    }
    
    return result, nil
}
```

### 3.4 全量验证响应示例

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION.FAILED",
    "message": "验证失败",
    "fields": [
      {
        "field": "email",
        "message": "邮箱格式不正确",
        "code": "VALIDATION.EMAIL"
      },
      {
        "field": "username",
        "message": "长度不能少于 3",
        "code": "VALIDATION.MIN_LENGTH"
      },
      {
        "field": "age",
        "message": "必须大于或等于 18",
        "code": "VALIDATION.MIN"
      }
    ]
  }
}
```

## 场景4：嵌套对象验证

展示嵌套对象的验证和字段路径。

### 4.1 定义模型

```go
package models

type Address struct {
    Province string `json:"province"`
    City     string `json:"city"`
    Street   string `json:"street"`
    ZipCode  string `json:"zip_code"`
}

type CreateProfileRequest struct {
    Username string  `json:"username"`
    Address  Address `json:"address"`
}
```

### 4.2 验证器

```go
package validators

func NewAddressValidator() *validation.AbstractValidator[Address] {
    v := validation.NewValidator[Address]()
    
    v.Field(func(a *Address) string { return a.Province }).
        NotEmpty()
    
    v.Field(func(a *Address) string { return a.City }).
        NotEmpty()
    
    v.Field(func(a *Address) string { return a.Street }).
        NotEmpty().
        MinLength(5)
    
    v.Field(func(a *Address) string { return a.ZipCode }).
        NotEmpty().
        Matches(`^\d{6}$`).
        WithMessage("邮政编码必须是6位数字")
    
    return v
}

func NewCreateProfileValidator() *validation.AbstractValidator[CreateProfileRequest] {
    v := validation.NewValidator[CreateProfileRequest]()
    
    v.Field(func(r *CreateProfileRequest) string { return r.Username }).
        NotEmpty().
        MinLength(3)
    
    // 嵌套对象验证
    addressValidator := NewAddressValidator()
    // 注意：这需要在 builder.go 中实现 MustBeValid 方法
    // 这里仅作示例
    
    return v
}
```

### 4.3 嵌套验证响应

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION.FAILED",
    "message": "验证失败",
    "fields": [
      {
        "field": "address.city",
        "message": "不能为空",
        "code": "VALIDATION.REQUIRED"
      },
      {
        "field": "address.zipCode",
        "message": "邮政编码必须是6位数字",
        "code": "VALIDATION.PATTERN"
      }
    ]
  }
}
```

## 总结

通过以上示例，您可以看到：

1. **验证模式选择**：
   - 表单验证使用 `NewValidator()`（快速失败）
   - 批量导入使用 `NewValidatorAll()`（全量验证）

2. **错误码使用**：
   - 验证规则自动使用框架错误码
   - 业务逻辑使用错误码构建器

3. **统一错误处理**：
   - Service 层返回业务错误
   - Controller 层统一转换为 HTTP 响应

4. **清晰的职责分离**：
   - 验证器负责数据验证
   - Service 负责业务逻辑
   - Controller 负责 HTTP 处理

这样的架构使代码更加清晰、可维护，并且易于测试。
