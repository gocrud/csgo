# 完整的 DI 功能演示

本示例展示了 ego 框架中依赖注入的所有核心功能，完全符合 .NET 10 风格。

## 🎯 演示的功能

### 1. 基础服务注册和解析
- Singleton/Scoped/Transient 三种生命周期
- 指针填充方案（类似 json.Unmarshal）

### 2. Keyed Services（命名服务）
- 注册多个同类型服务
- 通过 serviceKey 获取特定实现

### 3. Scoped 生命周期
- 不同作用域创建不同实例
- 自动资源释放（Dispose）

### 4. TryGetService（可选服务）
- 优雅处理服务不存在的情况
- 无异常的服务查询

### 5. GetServices（多服务解析）
- 插件模式支持
- 获取所有实现

### 6. 泛型辅助方法
- 最简洁的语法糖
- 类型安全

### 7. IsService（服务查询）
- 检查服务是否已注册
- 运行时服务发现

### 8. IServiceScopeFactory
- 作用域工厂模式
- 符合 .NET 标准

## 🚀 运行示例

```bash
cd examples/complete_di_demo
go run main.go
```

## 📝 核心 API

### 服务注册（.NET 风格）
```go
services := di.NewServiceCollection()
services.
    AddSingleton(NewLogger).
    AddScoped(NewDatabase).
    AddTransient(NewService).
    AddKeyedSingleton("primary", NewPrimaryDb).
    AddKeyedSingleton("secondary", NewSecondaryDb)
```

### 服务获取（Go 风格 - 指针填充）
```go
// 方式 1：指针填充（推荐）
var logger ILogger
provider.GetRequiredService(&logger)

// 方式 2：泛型（可选，更简洁）
logger := di.GetRequiredService[ILogger](provider)

// 方式 3：可选服务
var cache ICache
if provider.TryGetService(&cache) {
    // 使用缓存
}

// 方式 4：命名服务
var primaryDb IDatabase
provider.GetKeyedService(&primaryDb, "primary")

// 方式 5：所有服务
var databases []IDatabase
provider.GetServices(&databases)
```

### 作用域使用
```go
scope := provider.CreateScope()
defer scope.Dispose()

scopedProvider := scope.ServiceProvider()
var service IUserService
scopedProvider.GetRequiredService(&service)
```

## ✨ 设计理念

**"注册像 .NET，获取像 Go"**

- **服务注册**：采用 .NET 风格的链式调用和命名
- **服务获取**：采用 Go 习惯的指针填充（类似 json.Unmarshal）
- **完整功能**：支持 .NET 10 的所有核心 DI 特性

