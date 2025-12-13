# 第一个应用

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

让我们创建一个更完整的 Hello World 应用，学习 CSGO 框架的基本用法。

## Hello World

### 基础版本

最简单的 CSGO 应用：

```go
package main

import "github.com/gocrud/csgo/web"

func main() {
    builder := web.CreateBuilder()
    app := builder.Build()
    
    app.MapGet("/", func(c *web.HttpContext) web.IActionResult {
        return c.Ok(web.M{"message": "Hello, World!"})
    })
    
    app.Run()
}
```

运行并访问 http://localhost:8080/

## 添加更多路由

### 多个端点

```go
package main

import (
    "github.com/gocrud/csgo/web"
    "github.com/gin-gonic/gin"
)

func main() {
    builder := web.CreateBuilder()
    app := builder.Build()
    
    // 根路径
    app.MapGet("/", func(c *web.HttpContext) web.IActionResult {
        return c.Ok(web.M{"message": "欢迎使用 CSGO 框架"})
    })
    
    // 问候端点
    app.MapGet("/hello/:name", func(c *web.HttpContext) web.IActionResult {
        name := c.RawCtx().Param("name")
        return c.Ok(web.M{
            "message": "Hello, " + name + "!",
        })
    })
    
    // 健康检查
    app.MapGet("/health", func(c *web.HttpContext) web.IActionResult {
        return c.Ok(web.M{
            "status": "healthy",
            "version": "1.0.0",
        })
    })
    
    app.Run()
}
```

### 测试端点

```bash
# 根路径
curl http://localhost:8080/

# 问候
curl http://localhost:8080/hello/张三

# 健康检查
curl http://localhost:8080/health
```

## 使用依赖注入

### 创建服务

```go
package main

import (
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
    "github.com/gin-gonic/gin"
    "time"
)

// 定义服务
type GreetingService struct{}

func NewGreetingService() *GreetingService {
    return &GreetingService{}
}

func (s *GreetingService) Greet(name string) string {
    hour := time.Now().Hour()
    var greeting string
    
    if hour < 12 {
        greeting = "早上好"
    } else if hour < 18 {
        greeting = "下午好"
    } else {
        greeting = "晚上好"
    }
    
    return greeting + "，" + name + "！"
}

func main() {
    builder := web.CreateBuilder()
    
    // 注册服务
    builder.Services.Add(NewGreetingService)
    
    app := builder.Build()
    
    // 使用服务
    app.MapGet("/greet/:name", func(c *web.HttpContext) web.IActionResult {
        name := c.RawCtx().Param("name")
        
        // 从 DI 容器获取服务
        service := di.Get[*GreetingService](c.Services)
        message := service.Greet(name)
        
        return c.Ok(web.M{"message": message})
    })
    
    app.Run()
}
```

测试：

```bash
curl http://localhost:8080/greet/张三
# 根据时间返回：{"message":"早上好，张三！"}
```

## POST 请求

### 创建数据

```go
package main

import (
    "github.com/gocrud/csgo/web"
    "github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    builder := web.CreateBuilder()
    app := builder.Build()
    
    var users []User
    nextID := 1
    
    // 获取用户列表
    app.MapGet("/users", func(c *web.HttpContext) web.IActionResult {
        return c.Ok(users)
    })
    
    // 创建用户
    app.MapPost("/users", func(c *web.HttpContext) web.IActionResult {
        var req CreateUserRequest
        
        // 绑定 JSON 请求体
        if err := c.MustBindJSON(&req); err != nil {
            return err
        }
        
        // 创建用户
        user := User{
            ID:    nextID,
            Name:  req.Name,
            Email: req.Email,
        }
        nextID++
        users = append(users, user)
        
        // 返回 201 Created
        return c.Created(user)
    })
    
    app.Run()
}
```

### 测试 POST

```bash
# 创建用户
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"张三","email":"zhangsan@example.com"}'

# 查看用户列表
curl http://localhost:8080/users
```

## 路由组

### 组织 API

```go
package main

import (
    "github.com/gocrud/csgo/web"
    "github.com/gin-gonic/gin"
)

func main() {
    builder := web.CreateBuilder()
    app := builder.Build()
    
    // API v1 路由组
    api := app.MapGroup("/api/v1")
    {
        // 用户相关
        users := api.MapGroup("/users")
        {
            users.MapGet("", listUsers)
            users.MapPost("", createUser)
            users.MapGet("/:id", getUser)
        }
        
        // 产品相关
        products := api.MapGroup("/products")
        {
            products.MapGet("", listProducts)
            products.MapGet("/:id", getProduct)
        }
    }
    
    app.Run()
}

func listUsers(c *web.HttpContext) web.IActionResult {
    return c.Ok(web.M{"users": []string{}})
}

func createUser(c *web.HttpContext) web.IActionResult {
    return c.Created(web.M{"message": "User created"})
}

func getUser(c *web.HttpContext) web.IActionResult {
    id := c.RawCtx().Param("id")
    return c.Ok(web.M{"id": id})
}

func listProducts(c *web.HttpContext) web.IActionResult {
    return c.Ok(web.M{"products": []string{}})
}

func getProduct(c *web.HttpContext) web.IActionResult {
    id := c.RawCtx().Param("id")
    return c.Ok(web.M{"id": id})
}
```

访问：
- http://localhost:8080/api/v1/users
- http://localhost:8080/api/v1/products

## 练习

### 练习 1：计算器 API

创建一个简单的计算器 API：

- `GET /calc/add?a=10&b=20` - 加法
- `GET /calc/subtract?a=10&b=5` - 减法
- `GET /calc/multiply?a=10&b=3` - 乘法
- `GET /calc/divide?a=10&b=2` - 除法

### 练习 2：待办事项 API

创建一个待办事项 API：

- `GET /todos` - 获取所有待办
- `POST /todos` - 创建待办
- `GET /todos/:id` - 获取单个待办
- `PUT /todos/:id` - 更新待办
- `DELETE /todos/:id` - 删除待办

待办事项数据结构：

```go
type Todo struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}
```

## 小结

本章你学会了：

- ✅ 创建基本的 Web 应用
- ✅ 定义多个路由
- ✅ 使用依赖注入
- ✅ 处理 GET 和 POST 请求
- ✅ 使用路由组组织 API

## 下一步

现在你已经能够创建简单的 Web 应用了，接下来让我们深入了解 CSGO 的核心概念。

继续学习：[核心概念](concepts.md) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

