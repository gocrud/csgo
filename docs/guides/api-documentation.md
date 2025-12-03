# Swagger 方法详解

本文档详细解释 `AddSwaggerGen` 和 `UseSwagger` 的区别和作用。

---

## 🎯 核心区别

| 方法 | 阶段 | 作用 | 类比 |
|------|------|------|------|
| **AddSwaggerGen** | 配置阶段（Builder） | 注册 Swagger 服务和配置 | 准备工具 |
| **UseSwagger** | 运行阶段（Application） | 启用 Swagger JSON 端点 | 使用工具 |
| **UseSwaggerUI** | 运行阶段（Application） | 启用 Swagger UI 界面 | 使用界面 |

---

## 📝 详细说明

### 1. AddSwaggerGen - 注册服务和配置

**位置：** `swagger/service_collection_extensions.go`

**作用：** 在**构建阶段**注册 Swagger 配置到依赖注入容器。

```go
// AddSwaggerGen adds Swagger generation services to the service collection.
// Corresponds to .NET services.AddSwaggerGen().
func AddSwaggerGen(services di.IServiceCollection, configure func(*SwaggerGenOptions)) {
    services.AddSingleton(func() *SwaggerGenOptions {
        opts := NewSwaggerGenOptions()
        if configure != nil {
            configure(opts)
        }
        return opts
    })
}
```

**功能：**
1. ✅ 创建 `SwaggerGenOptions` 配置对象
2. ✅ 应用用户的自定义配置
3. ✅ 将配置注册为 Singleton 服务
4. ✅ 存储在 DI 容器中，供后续使用

**使用时机：** 在 `builder.Build()` **之前**

```go
builder := web.CreateBuilder()

// ✅ 配置阶段：注册 Swagger 服务
swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
    opts.Title = "My API"
    opts.Version = "v1"
    opts.Description = "API documentation"
    opts.AddSecurityDefinition("Bearer", openapi.SecurityScheme{
        Type:   "http",
        Scheme: "bearer",
    })
})

app := builder.Build()  // 构建应用
```

---

### 2. UseSwagger - 启用 Swagger JSON 端点

**位置：** `swagger/application_builder_extensions.go`

**作用：** 在**运行阶段**注册 Swagger JSON 端点，生成 OpenAPI 规范。

```go
// UseSwagger enables the Swagger JSON endpoint.
// Corresponds to .NET app.UseSwagger().
func UseSwagger(app *web.WebApplication) {
    // 1. 从 DI 容器获取配置
    var opts *SwaggerGenOptions
    err := app.Services.GetService(&opts)
    if err != nil {
        opts = NewSwaggerGenOptions()
    }
    
    // 2. 创建 OpenAPI 生成器
    generator := openapi.NewGenerator(opts.Title, opts.Version).
        WithDescription(opts.Description)
    
    // 3. 添加安全方案
    for name, scheme := range opts.SecurityDefinitions {
        generator.WithSecurityScheme(name, scheme)
    }
    
    // 4. 注册 Swagger JSON 端点
    app.MapGet("/swagger/v1/swagger.json", func(c *gin.Context) {
        routes := app.GetRoutes()
        routeInfos := make([]openapi.RouteInfo, len(routes))
        for i, r := range routes {
            routeInfos[i] = r
        }
        
        spec := generator.Generate(routeInfos)
        c.JSON(200, spec)
    })
}
```

**功能：**
1. ✅ 从 DI 容器获取 Swagger 配置
2. ✅ 创建 OpenAPI 规范生成器
3. ✅ 注册 `/swagger/v1/swagger.json` 端点
4. ✅ 动态生成 OpenAPI JSON 规范

**使用时机：** 在 `builder.Build()` **之后**

```go
app := builder.Build()

// ✅ 运行阶段：启用 Swagger JSON 端点
swagger.UseSwagger(app)

// 现在可以访问：http://localhost:8080/swagger/v1/swagger.json
```

---

### 3. UseSwaggerUI - 启用 Swagger UI 界面

**位置：** `swagger/application_builder_extensions.go`

**作用：** 在**运行阶段**注册 Swagger UI 界面端点。

```go
// UseSwaggerUI enables the Swagger UI.
// Corresponds to .NET app.UseSwaggerUI().
func UseSwaggerUI(app *web.WebApplication, configure ...func(*SwaggerUIOptions)) {
    opts := NewSwaggerUIOptions()
    if len(configure) > 0 && configure[0] != nil {
        configure[0](opts)
    }
    
    // Register Swagger UI endpoints
    handler := func(c *gin.Context) {
        c.Header("Content-Type", "text/html; charset=utf-8")
        c.String(200, getSwaggerUIHTML(opts))
    }
    
    app.MapGet(opts.RoutePrefix+"/index.html", handler)
    app.MapGet(opts.RoutePrefix+"/", handler)
    app.MapGet(opts.RoutePrefix, func(c *gin.Context) {
        c.Redirect(301, opts.RoutePrefix+"/index.html")
    })
}
```

**功能：**
1. ✅ 注册 `/swagger` 端点
2. ✅ 返回 Swagger UI HTML 页面
3. ✅ 自动加载 `/swagger/v1/swagger.json`
4. ✅ 提供交互式 API 文档界面

**使用时机：** 在 `builder.Build()` **之后**

```go
app := builder.Build()

swagger.UseSwagger(app)
swagger.UseSwaggerUI(app)

// 现在可以访问：http://localhost:8080/swagger
```

---

## 🔄 完整流程

### 步骤 1: 配置阶段（Builder）

```go
builder := web.CreateBuilder()

// ✅ 注册 Swagger 服务和配置
swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
    opts.Title = "My API"
    opts.Version = "v1"
    opts.Description = "API documentation"
})
```

**发生了什么：**
1. 创建 `SwaggerGenOptions` 对象
2. 应用用户配置（Title、Version、Description）
3. 将配置注册到 DI 容器
4. **此时还没有任何端点被注册**

---

### 步骤 2: 定义路由

```go
app := builder.Build()

// 定义 API 路由
app.MapGet("/api/users", func(c *gin.Context) {
    c.JSON(200, []string{"Alice", "Bob"})
}).WithSummary("List users")
```

---

### 步骤 3: 启用 Swagger（Application）

```go
// ✅ 启用 Swagger JSON 端点
swagger.UseSwagger(app)

// ✅ 启用 Swagger UI 界面
swagger.UseSwaggerUI(app)
```

**发生了什么：**
1. `UseSwagger` 注册 `/swagger/v1/swagger.json` 端点
2. 该端点会收集所有已注册的路由
3. 生成 OpenAPI JSON 规范
4. `UseSwaggerUI` 注册 `/swagger` 端点
5. 返回 Swagger UI HTML 页面

---

### 步骤 4: 运行应用

```go
app.Run()
```

**可以访问：**
- `http://localhost:8080/swagger/v1/swagger.json` - OpenAPI JSON 规范
- `http://localhost:8080/swagger` - Swagger UI 界面

---

## 📊 与 .NET 的对比

### .NET 代码

```csharp
var builder = WebApplication.CreateBuilder(args);

// 1. 配置阶段：注册 Swagger 服务
builder.Services.AddSwaggerGen(c =>
{
    c.SwaggerDoc("v1", new OpenApiInfo
    {
        Title = "My API",
        Version = "v1",
        Description = "API documentation"
    });
    c.AddSecurityDefinition("Bearer", new OpenApiSecurityScheme
    {
        Type = SecuritySchemeType.Http,
        Scheme = "bearer"
    });
});

var app = builder.Build();

// 2. 运行阶段：启用 Swagger 中间件
app.UseSwagger();      // 启用 /swagger/v1/swagger.json
app.UseSwaggerUI();    // 启用 /swagger UI

app.MapGet("/api/users", () => new[] { "Alice", "Bob" })
   .WithSummary("List users");

app.Run();
```

---

### Ego 代码

```go
builder := web.CreateBuilder()

// 1. 配置阶段：注册 Swagger 服务
swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
    opts.Title = "My API"
    opts.Version = "v1"
    opts.Description = "API documentation"
    opts.AddSecurityDefinition("Bearer", openapi.SecurityScheme{
        Type:   "http",
        Scheme: "bearer",
    })
})

app := builder.Build()

// 2. 运行阶段：启用 Swagger 中间件
swagger.UseSwagger(app)      // 启用 /swagger/v1/swagger.json
swagger.UseSwaggerUI(app)    // 启用 /swagger UI

app.MapGet("/api/users", func(c *gin.Context) {
    c.JSON(200, []string{"Alice", "Bob"})
}).WithSummary("List users")

app.Run()
```

**一致性：99%** ✅

---

## 🎯 常见问题

### Q1: 可以只用 UseSwagger 不用 AddSwaggerGen 吗？

**A:** 可以，但会使用默认配置。

```go
app := builder.Build()

// ✅ 可以工作，使用默认配置
swagger.UseSwagger(app)
swagger.UseSwaggerUI(app)
```

但推荐先用 `AddSwaggerGen` 配置：

```go
// ✅ 推荐：先配置
swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
    opts.Title = "My API"
    opts.Version = "v1"
})

app := builder.Build()
swagger.UseSwagger(app)
swagger.UseSwaggerUI(app)
```

---

### Q2: 可以只用 AddSwaggerGen 不用 UseSwagger 吗？

**A:** 不行！`AddSwaggerGen` 只是注册配置，不会创建任何端点。

```go
// ❌ 错误：没有 Swagger 端点
swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
    opts.Title = "My API"
})

app := builder.Build()
// 缺少 UseSwagger 和 UseSwaggerUI
app.Run()

// 访问 /swagger 会 404
```

---

### Q3: UseSwagger 和 UseSwaggerUI 的顺序重要吗？

**A:** 不重要，但推荐先 `UseSwagger` 再 `UseSwaggerUI`。

```go
// ✅ 推荐顺序
swagger.UseSwagger(app)      // 先启用 JSON 端点
swagger.UseSwaggerUI(app)    // 再启用 UI

// ✅ 也可以反过来
swagger.UseSwaggerUI(app)
swagger.UseSwagger(app)
```

---

### Q4: 可以只用 UseSwagger 不用 UseSwaggerUI 吗？

**A:** 可以！如果你只需要 OpenAPI JSON 规范。

```go
app := builder.Build()

// ✅ 只启用 JSON 端点
swagger.UseSwagger(app)

// 可以访问：http://localhost:8080/swagger/v1/swagger.json
// 但不能访问：http://localhost:8080/swagger（UI）
```

---

## 📋 完整示例

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gocrud/csgo/openapi"
    "github.com/gocrud/csgo/swagger"
    "github.com/gocrud/csgo/web"
)

func main() {
    // ========================================
    // 1. 配置阶段（Builder）
    // ========================================
    builder := web.CreateBuilder()
    
    // ✅ 注册 Swagger 服务和配置
    swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
        opts.Title = "My API"
        opts.Version = "v1"
        opts.Description = "API documentation"
        opts.AddSecurityDefinition("Bearer", openapi.SecurityScheme{
            Type:         "http",
            Scheme:       "bearer",
            BearerFormat: "JWT",
            Description:  "Enter JWT token",
        })
    })
    
    // ========================================
    // 2. 构建应用
    // ========================================
    app := builder.Build()
    
    // ========================================
    // 3. 运行阶段（Application）
    // ========================================
    
    // ✅ 启用 Swagger JSON 端点
    swagger.UseSwagger(app)
    
    // ✅ 启用 Swagger UI 界面
    swagger.UseSwaggerUI(app)
    
    // ========================================
    // 4. 定义路由
    // ========================================
    app.MapGet("/api/users", func(c *gin.Context) {
        c.JSON(200, []string{"Alice", "Bob"})
    }).WithSummary("List users")
    
    // ========================================
    // 5. 运行应用
    // ========================================
    println("Server: http://localhost:8080")
    println("Swagger JSON: http://localhost:8080/swagger/v1/swagger.json")
    println("Swagger UI: http://localhost:8080/swagger")
    
    app.Run()
}
```

---

## ✅ 总结

### AddSwaggerGen（配置阶段）

- **时机：** `builder.Build()` 之前
- **作用：** 注册 Swagger 配置到 DI 容器
- **类比：** 准备工具和材料

```go
swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
    opts.Title = "My API"
    opts.Version = "v1"
})
```

---

### UseSwagger（运行阶段）

- **时机：** `builder.Build()` 之后
- **作用：** 注册 `/swagger/v1/swagger.json` 端点
- **类比：** 使用工具生成 JSON 规范

```go
swagger.UseSwagger(app)
```

---

### UseSwaggerUI（运行阶段）

- **时机：** `builder.Build()` 之后
- **作用：** 注册 `/swagger` UI 界面端点
- **类比：** 使用界面展示文档

```go
swagger.UseSwaggerUI(app)
```

---

### 完整流程

```
1. AddSwaggerGen (配置)
   ↓
2. builder.Build() (构建)
   ↓
3. UseSwagger (启用 JSON)
   ↓
4. UseSwaggerUI (启用 UI)
   ↓
5. app.Run() (运行)
```

**记住：先配置（Add），再使用（Use）！** 🚀

