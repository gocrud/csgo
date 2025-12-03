# API 参考

本目录包含 CSGO 框架的完整 API 参考文档。

## 📚 API 文档索引

| 模块 | 文档 | 说明 |
|------|------|------|
| 依赖注入 | [di.md](di.md) | DI 容器、服务注册、服务解析 |
| Web 框架 | [web.md](web.md) | WebApplication、HttpContext、ActionResult、路由 |
| 配置管理 | [configuration.md](configuration.md) | 配置源、Options 模式、绑定 |
| 应用托管 | [hosting.md](hosting.md) | Host、HostedService、BackgroundService |

## 🔍 快速查找

### 依赖注入

```go
// 服务注册
services.AddSingleton(factory)
services.AddScoped(factory)
services.AddTransient(factory)

// 服务解析
provider.GetRequiredService(&target)
di.GetRequiredService[T](provider)
```

[查看完整 DI API →](di.md)

### Web 框架

```go
// 应用构建
builder := web.CreateBuilder()
app := builder.Build()

// 路由注册
app.MapGet("/path", handler)
app.MapGroup("/api")

// HttpContext
c.MustPathInt("id")
c.QueryInt("page", 1)
c.MustBindJSON(&req)

// ActionResult
c.Ok(data)
c.NotFound("message")
c.Created(data)
```

[查看完整 Web API →](web.md)

### 配置管理

```go
// 读取配置
config.Get("key")
config.GetSection("section")
config.Bind("section", &target)

// Options 模式
configuration.Configure[T](services, config, "section")
configuration.BindOptions[T](config, "section")
```

[查看完整配置 API →](configuration.md)

### 应用托管

```go
// 托管服务
services.AddHostedService(factory)

// 后台服务
type MyService struct {
    *hosting.BackgroundService
}

// 生命周期
host.Start(ctx)
host.Stop(ctx)
host.Run()
```

[查看完整托管 API →](hosting.md)

## 📖 相关指南

- [依赖注入指南](../guides/dependency-injection.md)
- [Web 应用指南](../guides/web-applications.md)
- [控制器指南](../guides/controllers.md)
- [配置管理指南](../guides/configuration.md)
- [应用托管指南](../guides/hosting.md)

