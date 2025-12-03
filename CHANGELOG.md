# 更新日志

## [未发布] - 2024-12

### 🔄 重大变更

#### 移除 Scoped 生命周期
- **移除** `AddScoped()` 方法
- **移除** `TryAddScoped()` 方法
- **移除** `AddKeyedScoped()` 方法
- **移除** `CreateScope()` 方法
- **移除** `IServiceScope` 接口
- **移除** `IServiceScopeFactory` 接口
- **移除** `GetServiceScopeFactory()` 函数
- **移除** `WithValidateScopes()` 配置选项

**原因：** 简化框架设计，提升性能，符合 Go 生态习惯。

**迁移指南：**
```go
// 旧代码
services.AddScoped(NewUserService)
scope := provider.CreateScope()

// 新代码
services.AddTransient(NewUserService)  // 或 AddSingleton
// Controllers 是单例，不需要 Scope
```

### ✨ 改进

#### Controllers 生命周期明确化
- **明确** Controllers 是单例，在 `MapControllers()` 时创建一次
- **优化** `AddController()` 不再注册到 DI 容器，只存储工厂函数
- **添加** 详细的文档说明和警告注释

#### 底层代码优化
- **移除** `engine.createInstance()` 的 `scopeInstances` 参数
- **移除** `serviceProvider.validateScopes` 字段
- **移除** `ServiceProviderOptions.ValidateScopes` 配置
- **简化** 服务生命周期管理逻辑

#### 配置系统增强
- **完善** IOptions 和 IOptionsMonitor 使用文档
- **添加** Controller 中使用配置的完整示例
- **明确** 配置注入的最佳实践

### 📚 文档更新

#### 新增文档
- **添加** [FRAMEWORK_CHANGES.md](docs/FRAMEWORK_CHANGES.md) - 框架设计决策和变更说明
- **添加** [QUICK_REFERENCE.md](docs/QUICK_REFERENCE.md) - 快速参考手册

#### 更新文档
- **更新** [README.md](README.md) - 添加新文档链接
- **更新** [docs/api/di.md](docs/api/di.md) - 移除 Scoped 相关 API，添加配置注入示例
- **更新** [docs/guides/controllers.md](docs/guides/controllers.md) - 添加生命周期警告和配置使用
- **更新** [docs/guides/dependency-injection.md](docs/guides/dependency-injection.md) - 更新生命周期说明
- **更新** [examples/complete_di_demo/](examples/complete_di_demo/) - 移除 Scoped 示例
- **更新** [examples/di_pointer_filling_demo/](examples/di_pointer_filling_demo/) - 移除 Scoped 示例

### 🐛 修复

- **修复** 控制器注册语义不一致的问题（注释说 Transient，实际是 Singleton）
- **修复** `ControllerOptions.EnableEndpointMetadata` 未使用的问题（已移除）

### 🎯 性能优化

- **优化** 移除运行时作用域管理开销
- **优化** 简化服务解析路径
- **优化** 减少不必要的参数传递

---

## 支持的生命周期

当前框架支持以下两种生命周期：

| 生命周期 | 创建时机 | 适用场景 |
|---------|---------|---------|
| **Singleton** | 应用启动时创建一次 | 数据库连接池、配置、缓存、无状态服务 |
| **Transient** | 每次请求时创建 | 有状态服务、轻量级操作、请求日志 |

---

## 迁移检查清单

如果你的代码使用了以下 API，需要进行迁移：

- [ ] `AddScoped()` → 改为 `AddTransient()` 或 `AddSingleton()`
- [ ] `TryAddScoped()` → 改为 `TryAddTransient()` 或 `TryAddSingleton()`
- [ ] `AddKeyedScoped()` → 改为 `AddKeyedTransient()` 或 `AddKeyedSingleton()`
- [ ] `CreateScope()` → 移除，Controllers 是单例
- [ ] `IServiceScope` → 移除相关代码
- [ ] `GetServiceScopeFactory()` → 移除相关代码
- [ ] Controller 中的请求状态 → 从 `HttpContext` 获取

---

## 相关资源

- [框架变更详细说明](docs/FRAMEWORK_CHANGES.md)
- [快速参考手册](docs/QUICK_REFERENCE.md)
- [依赖注入指南](docs/guides/dependency-injection.md)
- [控制器指南](docs/guides/controllers.md)

---

**注意：** 这些变更旨在简化框架，提升性能，并更好地符合 Go 语言的设计哲学。

