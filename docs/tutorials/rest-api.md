# REST API 教程

本教程将带你从零开始，使用 CSGO 框架构建一个完整的 REST API。

## 目录

- [准备工作](#准备工作)
- [第 1 步：创建项目](#第-1-步创建项目)
- [第 2 步：定义模型](#第-2-步定义模型)
- [第 3 步：创建服务层](#第-3-步创建服务层)
- [第 4 步：创建控制器](#第-4-步创建控制器)
- [第 5 步：配置 Swagger](#第-5-步配置-swagger)
- [第 6 步：添加验证](#第-6-步添加验证)
- [第 7 步：错误处理](#第-7-步错误处理)
- [第 8 步：运行和测试](#第-8-步运行和测试)
- [完整代码](#完整代码)
- [下一步](#下一步)

---

## 准备工作

确保你已经：
- 安装 Go 1.18+
- 了解基本的 Go 语法
- 阅读过 [快速开始](../getting-started.md)

## 第 1 步：创建项目

```bash
# 创建项目目录
mkdir todo-api
cd todo-api

# 初始化 Go 模块
go mod init todo-api

# 安装 CSGO
go get github.com/gocrud/csgo
```

创建项目结构：

```
todo-api/
├── main.go
├── models/
│   └── todo.go
├── services/
│   ├── todo_service.go
│   └── extensions.go
├── controllers/
│   ├── todo_controller.go
│   └── extensions.go
└── go.mod
```

## 第 2 步：定义模型

创建 `models/todo.go`：

```go
package models

import "time"

// Todo 待办事项
type Todo struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description,omitempty"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTodoRequest 创建待办请求
type CreateTodoRequest struct {
    Title       string `json:"title" binding:"required,min=1,max=200"`
    Description string `json:"description" binding:"max=1000"`
}

// UpdateTodoRequest 更新待办请求
type UpdateTodoRequest struct {
    Title       *string `json:"title" binding:"omitempty,min=1,max=200"`
    Description *string `json:"description" binding:"omitempty,max=1000"`
    Completed   *bool   `json:"completed"`
}

// TodoListResponse 待办列表响应
type TodoListResponse struct {
    Items      []Todo `json:"items"`
    Total      int    `json:"total"`
    Page       int    `json:"page"`
    PageSize   int    `json:"page_size"`
    TotalPages int    `json:"total_pages"`
}
```

## 第 3 步：创建服务层

创建 `services/todo_service.go`：

```go
package services

import (
    "errors"
    "sync"
    "time"
    
    "todo-api/models"
)

var (
    ErrTodoNotFound = errors.New("todo not found")
)

// TodoService 待办服务
type TodoService struct {
    mu     sync.RWMutex
    todos  map[int]*models.Todo
    nextID int
}

// NewTodoService 创建待办服务
func NewTodoService() *TodoService {
    return &TodoService{
        todos:  make(map[int]*models.Todo),
        nextID: 1,
    }
}

// List 获取待办列表（分页）
func (s *TodoService) List(page, pageSize int, completed *bool) *models.TodoListResponse {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    // 收集所有待办
    var allTodos []models.Todo
    for _, todo := range s.todos {
        // 过滤完成状态
        if completed != nil && todo.Completed != *completed {
            continue
        }
        allTodos = append(allTodos, *todo)
    }
    
    // 计算分页
    total := len(allTodos)
    totalPages := (total + pageSize - 1) / pageSize
    
    start := (page - 1) * pageSize
    end := start + pageSize
    
    if start > total {
        start = total
    }
    if end > total {
        end = total
    }
    
    return &models.TodoListResponse{
        Items:      allTodos[start:end],
        Total:      total,
        Page:       page,
        PageSize:   pageSize,
        TotalPages: totalPages,
    }
}

// GetByID 根据 ID 获取待办
func (s *TodoService) GetByID(id int) (*models.Todo, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    todo, ok := s.todos[id]
    if !ok {
        return nil, ErrTodoNotFound
    }
    return todo, nil
}

// Create 创建待办
func (s *TodoService) Create(req *models.CreateTodoRequest) *models.Todo {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    now := time.Now()
    todo := &models.Todo{
        ID:          s.nextID,
        Title:       req.Title,
        Description: req.Description,
        Completed:   false,
        CreatedAt:   now,
        UpdatedAt:   now,
    }
    
    s.todos[s.nextID] = todo
    s.nextID++
    
    return todo
}

// Update 更新待办
func (s *TodoService) Update(id int, req *models.UpdateTodoRequest) (*models.Todo, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    todo, ok := s.todos[id]
    if !ok {
        return nil, ErrTodoNotFound
    }
    
    if req.Title != nil {
        todo.Title = *req.Title
    }
    if req.Description != nil {
        todo.Description = *req.Description
    }
    if req.Completed != nil {
        todo.Completed = *req.Completed
    }
    todo.UpdatedAt = time.Now()
    
    return todo, nil
}

// Delete 删除待办
func (s *TodoService) Delete(id int) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    if _, ok := s.todos[id]; !ok {
        return ErrTodoNotFound
    }
    
    delete(s.todos, id)
    return nil
}

// ToggleComplete 切换完成状态
func (s *TodoService) ToggleComplete(id int) (*models.Todo, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    todo, ok := s.todos[id]
    if !ok {
        return nil, ErrTodoNotFound
    }
    
    todo.Completed = !todo.Completed
    todo.UpdatedAt = time.Now()
    
    return todo, nil
}
```

创建 `services/extensions.go`：

```go
package services

import "github.com/gocrud/csgo/di"

// AddServices 注册所有服务
func AddServices(services di.IServiceCollection) {
    services.AddSingleton(NewTodoService)
}
```

## 第 4 步：创建控制器

创建 `controllers/todo_controller.go`：

```go
package controllers

import (
    "todo-api/models"
    "todo-api/services"
    
    "github.com/gocrud/csgo/web"
)

// TodoController 待办控制器
type TodoController struct {
    todoService *services.TodoService
}

// NewTodoController 创建待办控制器
func NewTodoController(todoService *services.TodoService) *TodoController {
    return &TodoController{todoService: todoService}
}

// MapRoutes 实现 IController 接口
func (ctrl *TodoController) MapRoutes(app *web.WebApplication) {
    todos := app.MapGroup("/api/todos").
        WithOpenApi(
            openapi.Tags("Todos"),
        )
    
    todos.MapGet("", ctrl.List).
        WithSummary("获取待办列表").
        WithDescription("获取所有待办事项，支持分页和过滤")
    
    todos.MapGet("/:id", ctrl.GetByID).
        WithSummary("获取待办详情").
        WithDescription("根据 ID 获取待办事项详情")
    
    todos.MapPost("", ctrl.Create).
        WithSummary("创建待办").
        WithDescription("创建新的待办事项")
    
    todos.MapPut("/:id", ctrl.Update).
        WithSummary("更新待办").
        WithDescription("更新待办事项")
    
    todos.MapDelete("/:id", ctrl.Delete).
        WithSummary("删除待办").
        WithDescription("删除待办事项")
    
    todos.MapPost("/:id/toggle", ctrl.ToggleComplete).
        WithSummary("切换完成状态").
        WithDescription("切换待办事项的完成状态")
}

// List 获取待办列表
func (ctrl *TodoController) List(c *web.HttpContext) web.IActionResult {
    page := c.QueryInt("page", 1)
    pageSize := c.QueryInt("page_size", 10)
    
    // 限制分页大小
    if pageSize > 100 {
        pageSize = 100
    }
    if page < 1 {
        page = 1
    }
    
    // 获取完成状态过滤
    var completed *bool
    if completedStr := c.Query("completed"); completedStr != "" {
        val := completedStr == "true"
        completed = &val
    }
    
    result := ctrl.todoService.List(page, pageSize, completed)
    return c.Ok(result)
}

// GetByID 获取待办详情
func (ctrl *TodoController) GetByID(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    todo, err := ctrl.todoService.GetByID(id)
    if err != nil {
        if err == services.ErrTodoNotFound {
            return c.NotFound("待办事项不存在")
        }
        return c.InternalError(err.Error())
    }
    
    return c.Ok(todo)
}

// Create 创建待办
func (ctrl *TodoController) Create(c *web.HttpContext) web.IActionResult {
    var req models.CreateTodoRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    todo := ctrl.todoService.Create(&req)
    return c.Created(todo)
}

// Update 更新待办
func (ctrl *TodoController) Update(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    var req models.UpdateTodoRequest
    if err := c.MustBindJSON(&req); err != nil {
        return err
    }
    
    todo, err := ctrl.todoService.Update(id, &req)
    if err != nil {
        if err == services.ErrTodoNotFound {
            return c.NotFound("待办事项不存在")
        }
        return c.InternalError(err.Error())
    }
    
    return c.Ok(todo)
}

// Delete 删除待办
func (ctrl *TodoController) Delete(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    if err := ctrl.todoService.Delete(id); err != nil {
        if err == services.ErrTodoNotFound {
            return c.NotFound("待办事项不存在")
        }
        return c.InternalError(err.Error())
    }
    
    return c.NoContent()
}

// ToggleComplete 切换完成状态
func (ctrl *TodoController) ToggleComplete(c *web.HttpContext) web.IActionResult {
    id, err := c.MustPathInt("id")
    if err != nil {
        return err
    }
    
    todo, err := ctrl.todoService.ToggleComplete(id)
    if err != nil {
        if err == services.ErrTodoNotFound {
            return c.NotFound("待办事项不存在")
        }
        return c.InternalError(err.Error())
    }
    
    return c.Ok(todo)
}
```

创建 `controllers/extensions.go`：

```go
package controllers

import (
    "todo-api/services"
    
    "github.com/gocrud/csgo/di"
    "github.com/gocrud/csgo/web"
)

// AddControllers 注册所有控制器
func AddControllers(svc di.IServiceCollection) {
    web.AddController(svc, func(sp di.IServiceProvider) *TodoController {
        return NewTodoController(di.GetRequiredService[*services.TodoService](sp))
    })
}
```

## 第 5 步：配置 Swagger

创建 `main.go`：

```go
package main

import (
    "todo-api/controllers"
    "todo-api/services"
    
    "github.com/gin-gonic/gin"
    "github.com/gocrud/csgo/swagger"
    "github.com/gocrud/csgo/web"
)

func main() {
    // 1. 创建应用构建器
    builder := web.CreateBuilder()
    
    // 2. 注册服务
    services.AddServices(builder.Services)
    
    // 3. 注册控制器
    controllers.AddControllers(builder.Services)
    
    // 4. 配置 Swagger
    swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
        opts.Title = "Todo API"
        opts.Version = "v1"
        opts.Description = "待办事项管理 API"
    })
    
    // 5. 构建应用
    app := builder.Build()
    
    // 6. 配置中间件
    app.Use(gin.Logger())
    app.Use(gin.Recovery())
    
    // 7. 启用 Swagger
    swagger.UseSwagger(app)
    swagger.UseSwaggerUI(app)
    
    // 8. 映射控制器
    app.MapControllers()
    
    // 9. 根路由
    app.MapGet("/", func(c *web.HttpContext) web.IActionResult {
        return c.Ok(gin.H{
            "name":    "Todo API",
            "version": "v1",
            "docs":    "/swagger",
        })
    })
    
    // 10. 运行
    println("🚀 Todo API 启动成功!")
    println("   API: http://localhost:8080")
    println("   Swagger: http://localhost:8080/swagger")
    app.Run()
}
```

## 第 6 步：添加验证

模型中已使用 `binding` 标签进行验证：

```go
type CreateTodoRequest struct {
    Title       string `json:"title" binding:"required,min=1,max=200"`
    Description string `json:"description" binding:"max=1000"`
}
```

常用验证规则：

| 规则 | 说明 | 示例 |
|------|------|------|
| `required` | 必填 | `binding:"required"` |
| `min=N` | 最小长度/值 | `binding:"min=1"` |
| `max=N` | 最大长度/值 | `binding:"max=200"` |
| `email` | 邮箱格式 | `binding:"email"` |
| `url` | URL 格式 | `binding:"url"` |
| `oneof=a b` | 枚举值 | `binding:"oneof=active inactive"` |

## 第 7 步：错误处理

使用 `IActionResult` 统一错误处理：

```go
// 参数错误
id, err := c.MustPathInt("id")
if err != nil {
    return err  // 自动返回 400 Bad Request
}

// 业务错误
if err == services.ErrTodoNotFound {
    return c.NotFound("待办事项不存在")
}

// 服务器错误
return c.InternalError(err.Error())
```

响应格式：

```json
// 成功
{
  "success": true,
  "data": { ... }
}

// 错误
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "待办事项不存在"
  }
}
```

## 第 8 步：运行和测试

```bash
# 运行应用
go run main.go
```

使用 curl 测试：

```bash
# 创建待办
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "学习 CSGO 框架", "description": "完成 REST API 教程"}'

# 获取列表
curl http://localhost:8080/api/todos

# 获取详情
curl http://localhost:8080/api/todos/1

# 更新待办
curl -X PUT http://localhost:8080/api/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'

# 切换状态
curl -X POST http://localhost:8080/api/todos/1/toggle

# 删除待办
curl -X DELETE http://localhost:8080/api/todos/1
```

访问 Swagger UI：http://localhost:8080/swagger

---

## 完整代码

完整项目代码可在 `examples/todo_api/` 目录找到。

---

## 下一步

恭喜！你已经完成了一个完整的 REST API。接下来可以：

- [CRUD 应用教程](crud-app.md) - 学习数据库集成
- [控制器指南](../guides/controllers.md) - 深入控制器模式
- [API 文档指南](../guides/api-documentation.md) - 完善 Swagger 文档

---

## 相关资源

- [Web 应用指南](../guides/web-applications.md)
- [依赖注入指南](../guides/dependency-injection.md)
- [最佳实践](../best-practices.md)

