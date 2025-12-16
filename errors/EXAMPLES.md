# 错误处理完整示例

本文档提供了完整的错误处理使用示例，涵盖从定义模型、验证器到 Controller 的完整流程。

## 目录

- [场景1：用户管理 CRUD](#场景1用户管理-crud)
- [场景2：订单状态管理](#场景2订单状态管理)
- [场景3：批量数据导入](#场景3批量数据导入)
- [场景4：嵌套对象验证](#场景4嵌套对象验证)

## 场景1：用户管理 CRUD

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
        return errors.Business("ORDER").Custom("CANNOT_CANCEL", 
            "已发货或已完成的订单不能取消")
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
