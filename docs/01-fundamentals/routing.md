# 路由系统

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

本章讲解 CSGO 的路由系统。

## 完整文档

关于路由系统的完整详细文档，请查看：

👉 **[Web 框架完整文档 - 路由系统部分](../../web/README.md#路由系统)**

## 核心内容概览

### 1. 基本路由
- GET/POST/PUT/DELETE/PATCH
- 路径参数
- 查询参数

### 2. 路由组
- 创建路由组
- 嵌套路由组
- 路由组中间件

### 3. 路由模式
- 静态路径
- 动态参数
- 通配符

## 快速示例

```go
app := builder.Build()

// 基本路由
app.MapGet("/users", listUsers)
app.MapPost("/users", createUser)

// 路径参数
app.MapGet("/users/:id", getUser)

// 路由组
api := app.MapGroup("/api")
{
    v1 := api.MapGroup("/v1")
    {
        users := v1.MapGroup("/users")
        {
            users.MapGet("", listUsers)
            users.MapGet("/:id", getUser)
            users.MapPost("", createUser)
        }
    }
}

app.Run()
```

## 下一步

继续学习：[配置管理](configuration.md) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

