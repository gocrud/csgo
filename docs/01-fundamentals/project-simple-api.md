# 实践项目：简单 API

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

通过实践项目巩固所学知识。

## 项目目标

创建一个简单的用户管理 API，包含：
- 获取用户列表
- 获取单个用户
- 创建用户
- 更新用户
- 删除用户

## 项目结构

```
simple-api/
├── main.go
├── models/
│   └── user.go
├── services/
│   └── user_service.go
└── appsettings.json
```

## 实现步骤

### 1. 定义模型

**models/user.go：**

```go
package models

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

### 2. 实现服务

**services/user_service.go：**

```go
package services

import "simple-api/models"

type UserService struct {
    users  []models.User
    nextID int
}

func NewUserService() *UserService {
    return &UserService{
        users:  make([]models.User, 0),
        nextID: 1,
    }
}

func (s *UserService) GetAll() []models.User {
    return s.users
}

func (s *UserService) GetByID(id int) *models.User {
    for _, user := range s.users {
        if user.ID == id {
            return &user
        }
    }
    return nil
}

func (s *UserService) Create(req *models.CreateUserRequest) *models.User {
    user := models.User{
        ID:    s.nextID,
        Name:  req.Name,
        Email: req.Email,
    }
    s.nextID++
    s.users = append(s.users, user)
    return &user
}

func (s *UserService) Update(id int, req *models.CreateUserRequest) *models.User {
    for i, user := range s.users {
        if user.ID == id {
            s.users[i].Name = req.Name
            s.users[i].Email = req.Email
            return &s.users[i]
        }
    }
    return nil
}

func (s *UserService) Delete(id int) bool {
    for i, user := range s.users {
        if user.ID == id {
            s.users = append(s.users[:i], s.users[i+1:]...)
            return true
        }
    }
    return false
}
```

### 3. 实现 API

**main.go：**

```go
package main

import (
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
    "simple-api/models"
    "simple-api/services"
    "strconv"
)

func main() {
    builder := web.CreateBuilder()
    
    // 注册服务
    builder.Services.Add(services.NewUserService)
    
    app := builder.Build()
    
    // API 路由
    api := app.MapGroup("/api")
    users := api.MapGroup("/users")
    {
        users.MapGet("", listUsers)
        users.MapGet("/:id", getUser)
        users.MapPost("", createUser)
        users.MapPut("/:id", updateUser)
        users.MapDelete("/:id", deleteUser)
    }
    
    app.Run()
}

func listUsers(c *web.HttpContext) web.IActionResult {
    service := di.Get[*services.UserService](c.Services)
    users := service.GetAll()
    return c.Ok(users)
}

func getUser(c *web.HttpContext) web.IActionResult {
    idStr := c.RawCtx().Param("id")
    id, _ := strconv.Atoi(idStr)
    
    service := di.Get[*services.UserService](c.Services)
    user := service.GetByID(id)
    
    if user == nil {
        return c.NotFound("用户不存在")
    }
    
    return c.Ok(user)
}

func createUser(c *web.HttpContext) web.IActionResult {
    var req models.CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    service := di.Get[*services.UserService](c.Services)
    user := service.Create(&req)
    
    return c.Created(user)
}

func updateUser(c *web.HttpContext) web.IActionResult {
    idStr := c.RawCtx().Param("id")
    id, _ := strconv.Atoi(idStr)
    
    var req models.CreateUserRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    service := di.Get[*services.UserService](c.Services)
    user := service.Update(id, &req)
    
    if user == nil {
        return c.NotFound("用户不存在")
    }
    
    return c.Ok(user)
}

func deleteUser(c *web.HttpContext) web.IActionResult {
    idStr := c.RawCtx().Param("id")
    id, _ := strconv.Atoi(idStr)
    
    service := di.Get[*services.UserService](c.Services)
    if !service.Delete(id) {
        return c.NotFound("用户不存在")
    }
    
    return c.NoContent()
}
```

## 测试 API

```bash
# 创建用户
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"张三","email":"zhangsan@example.com"}'

# 获取用户列表
curl http://localhost:8080/api/users

# 获取单个用户
curl http://localhost:8080/api/users/1

# 更新用户
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"李四","email":"lisi@example.com"}'

# 删除用户
curl -X DELETE http://localhost:8080/api/users/1
```

## 小结

恭喜完成阶段 1！🎉

你已经学会了：
- ✅ 依赖注入
- ✅ Web 应用基础
- ✅ 路由系统
- ✅ 配置管理
- ✅ HTTP 上下文
- ✅ 构建简单 API

## 下一步

继续学习：[阶段 2：构建 API](../02-building-apis/) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

