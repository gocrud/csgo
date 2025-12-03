# DI 指针填充方案示例

本示例展示了新的依赖注入API，采用指针填充方案（类似 `json.Unmarshal`）。

## 核心特性

### 1. .NET 风格的服务注册

```go
services := di.NewServiceCollection()
services.
    AddSingleton(NewConsoleLogger).
    AddScoped(NewUserService).
    AddTransient(NewEmailNotificationService)
```

### 2. Go 风格的指针填充

```go
var userService IUserService
provider.GetRequiredService(&userService)
```

### 3. 支持多种使用模式

- **基础用法**：`GetService(&target)` 返回error
- **必需服务**：`GetRequiredService(&target)` panic on error
- **尝试获取**：`TryGetService(&target)` 返回bool
- **获取所有**：`GetServices(&[]IService)`
- **命名服务**：`GetKeyedService(&target, "key")`

### 4. 作用域支持

```go
scope := provider.CreateScope()
defer scope.Dispose()

scopedProvider := scope.ServiceProvider()
var service IUserService
scopedProvider.GetRequiredService(&service)
```

### 5. 可选的泛型辅助方法

```go
// 更简洁的语法
logger := di.GetRequiredService[ILogger](provider)
services, _ := di.GetServices[INotificationService](provider)
```

## 运行示例

```bash
cd examples/di_pointer_filling_demo
go run main.go
```

## 优势

1. ✅ **类型安全**：编译时检查，无需类型断言
2. ✅ **Go 习惯**：类似 `json.Unmarshal(&v)` 的API风格
3. ✅ **IDE 友好**：自动补全和类型推导
4. ✅ **简洁明了**：一行代码即可获取服务
5. ✅ **支持 Scoped**：完整的三种生命周期支持

