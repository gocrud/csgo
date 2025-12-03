# 托管 API 参考

本文档提供 CSGO 应用托管系统的完整 API 参考。

## 目录

- [IHost](#ihost)
- [HostBuilder](#hostbuilder)
- [IHostedService](#ihostedservice)
- [BackgroundService](#backgroundservice)
- [IHostEnvironment](#ihostenvironment)
- [生命周期](#生命周期)

---

## IHost

应用主机接口。

```go
import "github.com/gocrud/csgo/hosting"
```

### Run

阻塞运行主机直到收到停止信号。

```go
func (h *Host) Run() error
```

---

### RunAsync

异步运行主机。

```go
func (h *Host) RunAsync(ctx context.Context) error
```

---

### Start

启动主机。

```go
func (h *Host) Start(ctx context.Context) error
```

---

### Stop

停止主机。

```go
func (h *Host) Stop(ctx context.Context) error
```

---

### Services

获取服务提供者。

```go
func (h *Host) Services() di.IServiceProvider
```

---

## HostBuilder

主机构建器。

### CreateDefaultBuilder

创建默认主机构建器。

```go
func CreateDefaultBuilder(args ...string) *HostBuilder
```

**默认行为：**
- 加载 `appsettings.json`（如果存在）
- 加载环境特定配置（如 `appsettings.Development.json`）
- 加载环境变量
- 加载命令行参数

---

### 属性

| 属性 | 类型 | 说明 |
|------|------|------|
| `Services` | `di.IServiceCollection` | 服务集合 |
| `Configuration` | `configuration.IConfiguration` | 配置对象 |
| `Environment` | `IHostEnvironment` | 环境信息 |

---

### ConfigureServices

配置服务。

```go
func (b *HostBuilder) ConfigureServices(configure func(di.IServiceCollection)) *HostBuilder
```

**示例：**

```go
builder.ConfigureServices(func(services di.IServiceCollection) {
    services.AddSingleton(NewMyService)
})
```

---

### Build

构建主机。

```go
func (b *HostBuilder) Build() IHost
```

---

## IHostedService

托管服务接口。

```go
type IHostedService interface {
    // Start 启动服务
    Start(ctx context.Context) error
    
    // Stop 停止服务
    Stop(ctx context.Context) error
}
```

### 实现托管服务

```go
type MyHostedService struct{}

func NewMyHostedService() *MyHostedService {
    return &MyHostedService{}
}

func (s *MyHostedService) Start(ctx context.Context) error {
    log.Println("服务启动")
    // 执行启动逻辑
    return nil
}

func (s *MyHostedService) Stop(ctx context.Context) error {
    log.Println("服务停止")
    // 执行清理逻辑
    return nil
}
```

### 注册托管服务

```go
builder.Services.AddHostedService(func() hosting.IHostedService {
    return NewMyHostedService()
})
```

---

## BackgroundService

后台服务基类，简化长时间运行的后台任务实现。

### NewBackgroundService

创建后台服务基类。

```go
func NewBackgroundService() *BackgroundService
```

---

### SetExecuteFunc

设置执行函数。

```go
func (s *BackgroundService) SetExecuteFunc(execute func(ctx context.Context) error)
```

---

### StoppingToken

获取停止通知通道。

```go
func (s *BackgroundService) StoppingToken() <-chan struct{}
```

---

### 实现后台服务

```go
type DataSyncService struct {
    *hosting.BackgroundService
    interval time.Duration
}

func NewDataSyncService() *DataSyncService {
    service := &DataSyncService{
        BackgroundService: hosting.NewBackgroundService(),
        interval:          time.Minute * 5,
    }
    service.SetExecuteFunc(service.executeAsync)
    return service
}

func (s *DataSyncService) executeAsync(ctx context.Context) error {
    ticker := time.NewTicker(s.interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := s.syncData(); err != nil {
                log.Printf("同步失败: %v", err)
            }
        case <-s.StoppingToken():
            log.Println("收到停止信号")
            return nil
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

func (s *DataSyncService) syncData() error {
    log.Println("执行数据同步...")
    // 同步逻辑
    return nil
}
```

---

## IHostEnvironment

环境信息接口。

### EnvironmentName

获取环境名称。

```go
func (e *HostEnvironment) EnvironmentName() string
```

---

### ApplicationName

获取应用名称。

```go
func (e *HostEnvironment) ApplicationName() string
```

---

### ContentRootPath

获取内容根路径。

```go
func (e *HostEnvironment) ContentRootPath() string
```

---

### IsDevelopment

检查是否为开发环境。

```go
func (e *HostEnvironment) IsDevelopment() bool
```

---

### IsStaging

检查是否为预发布环境。

```go
func (e *HostEnvironment) IsStaging() bool
```

---

### IsProduction

检查是否为生产环境。

```go
func (e *HostEnvironment) IsProduction() bool
```

---

### 环境变量

通过 `GO_ENVIRONMENT` 环境变量设置环境：

```bash
# 开发环境
GO_ENVIRONMENT=Development go run main.go

# 生产环境
GO_ENVIRONMENT=Production go run main.go

# 预发布环境
GO_ENVIRONMENT=Staging go run main.go
```

---

## 生命周期

### 启动流程

1. 构建主机 (`Build()`)
2. 创建服务提供者
3. 启动所有 `IHostedService`（按注册顺序）
4. 应用开始运行

### 停止流程

1. 收到停止信号（SIGINT/SIGTERM）
2. 停止所有 `IHostedService`（按注册逆序）
3. 释放服务提供者
4. 主机退出

### 优雅关闭

```go
func main() {
    builder := hosting.CreateDefaultBuilder()
    
    builder.Services.AddHostedService(func() hosting.IHostedService {
        return NewMyService()
    })
    
    host := builder.Build()
    
    // Run() 会处理 SIGINT/SIGTERM 信号
    // 自动触发优雅关闭
    if err := host.Run(); err != nil {
        log.Fatal(err)
    }
}
```

---

## 完整示例

### 带后台服务的应用

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/gocrud/csgo/hosting"
    "github.com/gocrud/csgo/web"
)

// 健康检查服务
type HealthCheckService struct {
    *hosting.BackgroundService
    healthy bool
}

func NewHealthCheckService() *HealthCheckService {
    service := &HealthCheckService{
        BackgroundService: hosting.NewBackgroundService(),
        healthy:           true,
    }
    service.SetExecuteFunc(service.executeAsync)
    return service
}

func (s *HealthCheckService) executeAsync(ctx context.Context) error {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            s.checkHealth()
        case <-s.StoppingToken():
            return nil
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

func (s *HealthCheckService) checkHealth() {
    // 执行健康检查
    log.Println("执行健康检查...")
    s.healthy = true
}

func (s *HealthCheckService) IsHealthy() bool {
    return s.healthy
}

// 数据清理服务
type CleanupService struct {
    *hosting.BackgroundService
}

func NewCleanupService() *CleanupService {
    service := &CleanupService{
        BackgroundService: hosting.NewBackgroundService(),
    }
    service.SetExecuteFunc(service.executeAsync)
    return service
}

func (s *CleanupService) executeAsync(ctx context.Context) error {
    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            log.Println("清理过期数据...")
            // 清理逻辑
        case <-s.StoppingToken():
            return nil
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

func main() {
    builder := web.CreateBuilder()
    
    // 注册后台服务
    builder.Services.AddHostedService(func() hosting.IHostedService {
        return NewHealthCheckService()
    })
    
    builder.Services.AddHostedService(func() hosting.IHostedService {
        return NewCleanupService()
    })
    
    app := builder.Build()
    
    // 健康检查端点
    app.MapGet("/health", func(c *web.HttpContext) web.IActionResult {
        return c.Ok(gin.H{"status": "healthy"})
    })
    
    log.Println("应用启动...")
    app.Run()
}
```

---

## 相关文档

- [应用托管指南](../guides/hosting.md)
- [依赖注入 API](di.md)
- [Web 框架 API](web.md)

