# 快速开始

本指南将帮助你在 5 分钟内创建并运行第一个 CSGO 应用。

## 前置要求

- Go 1.18 或更高版本
- 基本的 Go 语言知识

## 安装

使用 `go get` 安装 CSGO：

```bash
go get github.com/gocrud/csgo
```

## 创建第一个应用

### 1. 初始化项目

```bash
mkdir my-ego-app
cd my-ego-app
go mod init my-ego-app
```

### 2. 创建 main.go

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/web"
)

// 服务接口
type IGreeterService interface {
	Greet(name string) string
}

// 服务实现
type GreeterService struct{}

func NewGreeterService() IGreeterService {
	return &GreeterService{}
}

func (s *GreeterService) Greet(name string) string {
	return "Hello, " + name + "!"
}

func main() {
	// 1. 创建 Web 应用构建器
	builder := web.CreateBuilder()

	// 2. 注册服务（依赖注入）
	builder.Services.AddSingleton(NewGreeterService)

	// 3. 构建应用
	app := builder.Build()

	// 4. 定义路由
	app.MapGet("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to CSGO Framework!",
			"docs":    "https://github.com/gocrud/csgo",
		})
	})

	app.MapGet("/greet/:name", func(c *gin.Context) {
		// 解析服务（指针填充方式）
		var greeter IGreeterService
		app.Services.GetRequiredService(&greeter)

		name := c.Param("name")
		c.JSON(200, gin.H{
			"greeting": greeter.Greet(name),
		})
	})

	// 5. 运行应用
	println("🚀 Server started on http://localhost:8080")
	println("📖 Try: http://localhost:8080/greet/World")
	app.Run(":8080")
}
```

### 3. 运行应用

```bash
go mod tidy
go run main.go
```

你会看到：

```
🚀 Server started on http://localhost:8080
📖 Try: http://localhost:8080/greet/World
```

### 4. 测试 API

在浏览器或使用 curl 访问：

```bash
# 首页
curl http://localhost:8080/
# {"message":"Welcome to CSGO Framework!","docs":"https://github.com/gocrud/csgo"}

# 问候接口
curl http://localhost:8080/greet/World
# {"greeting":"Hello, World!"}

curl http://localhost:8080/greet/CSGO
# {"greeting":"Hello, CSGO!"}
```

## 理解代码

让我们逐步理解这个应用：

### 1. 创建构建器

```go
builder := web.CreateBuilder()
```

`CreateBuilder()` 创建一个 Web 应用构建器，类似 .NET 的 `WebApplication.CreateBuilder()`。它会：
- 初始化依赖注入容器
- 配置默认设置
- 准备 Web 服务器

### 2. 注册服务

```go
builder.Services.AddSingleton(NewGreeterService)
```

将 `GreeterService` 注册为单例服务。工厂函数 `NewGreeterService` 会在第一次请求时调用，之后重用同一个实例。

**三种生命周期**：
- `AddSingleton` - 应用程序生命周期内只创建一次
- `AddScoped` - 每个 HTTP 请求创建一次
- `AddTransient` - 每次需要时都创建新实例

### 3. 构建应用

```go
app := builder.Build()
```

编译依赖图，创建服务提供者，返回可运行的应用实例。

### 4. 定义路由

```go
app.MapGet("/greet/:name", func(c *gin.Context) {
    // 路由处理逻辑
})
```

定义 HTTP GET 路由。CSGO 提供了类似 .NET 的路由 API：
- `MapGet` - GET 请求
- `MapPost` - POST 请求
- `MapPut` - PUT 请求
- `MapDelete` - DELETE 请求
- `MapGroup` - 路由组

### 5. 解析服务

```go
var greeter IGreeterService
app.Services.GetRequiredService(&greeter)
```

使用**指针填充方式**解析服务。这是 CSGO 的特色功能：
- Go 惯用法（类似 `json.Unmarshal`）
- 编译时类型检查
- 无需类型断言
- IDE 友好

**或使用泛型辅助方法**：

```go
greeter := di.GetRequiredService[IGreeterService](app.Services)
```

## 添加更多功能

### 添加数据库服务

```go
type IDatabase interface {
	Query(sql string) []map[string]interface{}
}

type PostgresDB struct{}

func NewPostgresDB() IDatabase {
	return &PostgresDB{}
}

func (db *PostgresDB) Query(sql string) []map[string]interface{} {
	// 数据库查询逻辑
	return []map[string]interface{}{}
}

// 注册服务
builder.Services.AddSingleton(NewPostgresDB)
```

### 使用 Scoped 服务

```go
type RequestContext struct {
	RequestID string
	UserID    string
}

func NewRequestContext(c *gin.Context) *RequestContext {
	return &RequestContext{
		RequestID: c.GetHeader("X-Request-ID"),
		UserID:    c.GetHeader("X-User-ID"),
	}
}

// 注册为 Scoped（每个请求一个实例）
builder.Services.AddScoped(NewRequestContext)
```

### 添加 CORS

```go
import "github.com/gocrud/csgo/swagger"

// 添加 CORS 支持
builder.AddCors(func(opts *web.CorsOptions) {
	opts.AllowOrigins = []string{"http://localhost:3000"}
	opts.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
})

app := builder.Build()

// 使用 CORS 中间件
app.UseCors()
```

### 添加 Swagger 文档

```go
import "github.com/gocrud/csgo/swagger"

// 添加 Swagger
swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
	opts.Title = "My API"
	opts.Version = "v1"
	opts.Description = "My first CSGO API"
})

app := builder.Build()

// 启用 Swagger UI
swagger.UseSwagger(app)
swagger.UseSwaggerUI(app)

// 访问 http://localhost:8080/swagger
```

## 下一步

恭喜！你已经创建了第一个 CSGO 应用。接下来可以：

### 深入学习核心概念

- [依赖注入详解](guides/dependency-injection.md) - 学习 DI 的高级用法
- [Web 应用开发](guides/web-applications.md) - 完整的 Web 开发指南
- [控制器模式](guides/controllers.md) - 使用控制器组织代码
- [业务模块](guides/business-modules.md) - 创建可复用的业务模块

### 查看示例代码

- [完整 DI 演示](../examples/complete_di_demo/) - 所有 DI 功能示例
- [业务模块演示](../examples/business_module_demo/) - 模块化设计示例
- [控制器 API 演示](../examples/controller_api_demo/) - 控制器模式示例

### 构建实际应用

- [教程：构建 REST API](tutorials/rest-api.md) - 完整的 CRUD API
- [教程：CRUD 应用](tutorials/crud-app.md) - 带数据库的应用
- [最佳实践](best-practices.md) - 推荐的代码组织

## 遇到问题？

- 查看 [FAQ](faq.md)
- 查看 [API 参考](api/)
- 提交 [Issue](https://github.com/gocrud/csgo/issues)

---

[← 返回首页](../README.md) | [依赖注入指南 →](guides/dependency-injection.md)

