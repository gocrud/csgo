# Web 应用基础

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

本章讲解 CSGO Web 应用的基础知识。

## 完整文档

关于 Web 框架的完整详细文档，请查看：

👉 **[Web 框架完整文档](../../web/README.md)**

## 核心内容概览

### 1. WebApplicationBuilder
- 创建构建器
- 配置服务
- 配置主机
- 访问配置和环境

### 2. WebApplication
- 运行应用
- 访问服务
- 优雅关闭

### 3. 应用生命周期
- 启动流程
- 运行阶段
- 关闭流程

### 4. 环境管理
- Development
- Production
- Staging
- 自定义环境

## 快速示例

```go
func main() {
    // 创建构建器
    builder := web.CreateBuilder()
    
    // 配置服务
    builder.Services.Add(NewUserService)
    
    // 访问环境
    if builder.Environment.IsDevelopment() {
        // 开发环境配置
    }
    
    // 构建应用
    app := builder.Build()
    
    // 定义路由
    app.MapGet("/", func(c *web.HttpContext) web.IActionResult {
        return c.Ok(web.M{"message": "Hello"})
    })
    
    // 运行应用
    app.Run()
}
```

## 下一步

继续学习：[路由系统](routing.md) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

