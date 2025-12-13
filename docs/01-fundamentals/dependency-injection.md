# 依赖注入详解

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

本章深入讲解 CSGO 的依赖注入系统。

## 完整文档

关于依赖注入的完整详细文档，请查看：

👉 **[依赖注入 (DI) 完整文档](../../di/README.md)**

## 核心内容概览

### 1. DI 容器工作原理
- 服务注册
- 服务解析
- 依赖图构建
- 生命周期管理

### 2. Singleton 生命周期
- 什么是 Singleton
- 适用场景
- 线程安全考虑

### 3. 服务注册方式
- `Add` - 注册单例服务
- `AddInstance` - 注册实例
- `AddNamed` - 注册命名服务
- `TryAdd` - 条件注册

### 4. 服务解析方式
- `Get[T]` - 泛型获取（推荐）
- `GetOr[T]` - 带默认值获取
- `TryGet[T]` - 安全获取
- `GetNamed[T]` - 获取命名服务
- `GetAll[T]` - 获取所有实例

### 5. 自动依赖解析
- 构造函数注入
- 依赖图自动构建
- 循环依赖检测

## 快速示例

```go
// 定义服务
type UserService struct {
    repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
    return &UserService{repo: repo}
}

// 注册服务
builder := web.CreateBuilder()
builder.Services.Add(NewUserRepository)
builder.Services.Add(NewUserService)  // 自动注入 UserRepository

// 使用服务
app := builder.Build()
app.MapGet("/users", func(c *web.HttpContext) web.IActionResult {
    userService := di.Get[*UserService](c.Services)
    users := userService.GetAll()
    return c.Ok(users)
})
```

## 下一步

继续学习：[Web 应用基础](web-basics.md) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

