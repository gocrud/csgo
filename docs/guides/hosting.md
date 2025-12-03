# 应用托管指南

本指南介绍 CSGO 框架的应用托管系统，包括主机构建、生命周期管理和后台服务。

## 目录

- [快速开始](#快速开始)
- [主机构建器](#主机构建器)
- [应用生命周期](#应用生命周期)
- [托管服务](#托管服务)
- [后台服务](#后台服务)
- [环境管理](#环境管理)
- [优雅关闭](#优雅关闭)
- [最佳实践](#最佳实践)

## 快速开始

### 基本 Worker Service

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/gocrud/csgo/hosting"
)

// 定义 Worker
type MyWorker struct {
    *hosting.BackgroundService
}

func NewMyWorker() *MyWorker {
    w := &MyWorker{
        BackgroundService: hosting.NewBackgroundService(),
    }
    w.SetExecuteFunc(w.ExecuteAsync)
    return w
}

func (w *MyWorker) ExecuteAsync(ctx context.Context) error {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            fmt.Println("Worker running at:", time.Now())
            
        case <-w.StoppingToken():
            fmt.Println("Worker stopping")
            return nil
            
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

func main() {
    // 创建主机构建器
    builder := hosting.CreateDefaultBuilder()
    
    // 注册托管服务
    builder.Services.AddHostedService(NewMyWorker)
    
    // 构建并运行
    host := builder.Build()
    host.Run()
}
```

### Web 应用托管

Web 应用也使用相同的托管系统：

```go
package main

import "github.com/gocrud/csgo/web"

func main() {
    // Web 应用构建器内置托管支持
    builder := web.CreateBuilder()
    
    // HTTP 服务器作为托管服务运行
    app := builder.Build()
    
    // Run() 启动所有托管服务并等待关闭
    app.Run()
}
```

## 主机构建器

主机构建器提供了配置和创建应用主机的能力。

### CreateDefaultBuilder

创建带有默认配置的主机构建器：

```go
builder := hosting.CreateDefaultBuilder(os.Args...)

// 默认配置包括：
// 1. 配置系统（JSON 文件、环境变量、命令行）
// 2. 环境检测（Development/Staging/Production）
// 3. 核心服务注册（Configuration、Environment、Lifetime）
```

**默认配置源顺序：**

1. `appsettings.json`
2. `appsettings.{Environment}.json`
3. 环境变量
4. 命令行参数

### CreateEmptyBuilder

创建空白主机构建器，用于完全自定义配置：

```go
builder := hosting.CreateEmptyBuilder()

// 手动配置所有内容
builder.Services.AddSingleton(NewMyService)
```

### 配置主机

#### 方式 1: 直接访问 Services（推荐）

```go
builder := hosting.CreateDefaultBuilder()

// 直接操作 Services
builder.Services.AddSingleton(NewDatabaseConnection)
builder.Services.AddScoped(NewUserService)
builder.Services.AddHostedService(NewBackgroundWorker)
```

#### 方式 2: 使用 ConfigureServices

```go
builder := hosting.CreateDefaultBuilder()

// 使用配置方法（更符合 .NET 风格）
builder.ConfigureServices(func(services di.IServiceCollection) {
    services.AddSingleton(NewDatabaseConnection)
    services.AddScoped(NewUserService)
    services.AddHostedService(NewBackgroundWorker)
})
```

### 主机配置

```go
builder := hosting.CreateDefaultBuilder()

// 配置应用配置
builder.ConfigureAppConfiguration(func(config configuration.IConfigurationBuilder) {
    config.AddJsonFile("custom-settings.json", true, false)
})

// 配置主机配置
builder.ConfigureHostConfiguration(func(config configuration.IConfigurationBuilder) {
    config.AddInMemoryCollection(map[string]string{
        "Host:ShutdownTimeout": "60",
    })
})
```

### 访问配置和环境

```go
builder := hosting.CreateDefaultBuilder()

// 访问配置
dbHost := builder.Configuration.Get("Database:Host")

// 检查环境
if builder.Environment.IsDevelopment() {
    // 开发环境特定配置
    builder.Services.AddSingleton(NewMockEmailService)
} else {
    builder.Services.AddSingleton(NewRealEmailService)
}
```

## 应用生命周期

### 生命周期阶段

应用生命周期包含以下阶段：

```
Starting → Started → Running → Stopping → Stopped
```

### 生命周期事件

```go
import "github.com/gocrud/csgo/hosting"

// 在服务中注入生命周期管理器
type MyService struct {
    lifetime hosting.IHostApplicationLifetime
}

func NewMyService(lifetime hosting.IHostApplicationLifetime) *MyService {
    return &MyService{lifetime: lifetime}
}

func (s *MyService) Initialize() {
    // 监听应用启动事件
    go func() {
        <-s.lifetime.ApplicationStarted()
        fmt.Println("Application started")
    }()
    
    // 监听应用停止事件
    go func() {
        <-s.lifetime.ApplicationStopping()
        fmt.Println("Application is stopping")
        // 执行清理工作
    }()
    
    go func() {
        <-s.lifetime.ApplicationStopped()
        fmt.Println("Application stopped")
    }()
}

func (s *MyService) RequestShutdown() {
    // 请求应用关闭
    s.lifetime.StopApplication()
}
```

### 主机运行模式

#### Run() - 阻塞运行

```go
host := builder.Build()

// 启动并阻塞，直到收到关闭信号
if err := host.Run(); err != nil {
    fmt.Printf("Host error: %v\n", err)
}
```

#### RunAsync() - 异步运行

```go
host := builder.Build()

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

// 在指定上下文中运行
if err := host.RunAsync(ctx); err != nil {
    fmt.Printf("Host error: %v\n", err)
}
```

#### Start/Stop - 手动控制

```go
host := builder.Build()

// 启动主机
ctx := context.Background()
if err := host.Start(ctx); err != nil {
    panic(err)
}

// 执行其他操作...
fmt.Println("Host is running")
time.Sleep(10 * time.Second)

// 停止主机
if err := host.Stop(ctx); err != nil {
    panic(err)
}
```

## 托管服务

托管服务是由主机管理生命周期的长期运行服务。

### IHostedService 接口

```go
type IHostedService interface {
    StartAsync(ctx context.Context) error
    StopAsync(ctx context.Context) error
}
```

### 实现托管服务

```go
// 数据同步服务
type DataSyncService struct {
    db     *Database
    logger *Logger
    stopCh chan struct{}
}

func NewDataSyncService(db *Database, logger *Logger) *DataSyncService {
    return &DataSyncService{
        db:     db,
        logger: logger,
        stopCh: make(chan struct{}),
    }
}

// StartAsync 在主机启动时调用
func (s *DataSyncService) StartAsync(ctx context.Context) error {
    s.logger.Info("Data sync service starting")
    
    go s.syncLoop()
    
    return nil
}

// StopAsync 在主机关闭时调用
func (s *DataSyncService) StopAsync(ctx context.Context) error {
    s.logger.Info("Data sync service stopping")
    
    close(s.stopCh)
    
    // 等待当前同步完成
    select {
    case <-time.After(5 * time.Second):
        s.logger.Warn("Sync service stop timeout")
    case <-ctx.Done():
        return ctx.Err()
    }
    
    return nil
}

func (s *DataSyncService) syncLoop() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            s.performSync()
        case <-s.stopCh:
            return
        }
    }
}

func (s *DataSyncService) performSync() {
    s.logger.Info("Performing data sync")
    // 同步逻辑...
}

// 注册服务
func main() {
    builder := hosting.CreateDefaultBuilder()
    
    builder.Services.AddSingleton(NewDatabase)
    builder.Services.AddSingleton(NewLogger)
    builder.Services.AddHostedService(NewDataSyncService)
    
    host := builder.Build()
    host.Run()
}
```

### 托管服务启动顺序

托管服务按注册顺序启动，按相反顺序停止：

```go
// 注册顺序
builder.Services.AddHostedService(NewDatabaseService)     // 1. 先启动
builder.Services.AddHostedService(NewCacheService)        // 2. 然后启动
builder.Services.AddHostedService(NewWebApiService)       // 3. 最后启动

// 关闭顺序
// 3. WebApiService 先停止
// 2. CacheService 然后停止
// 1. DatabaseService 最后停止
```

## 后台服务

`BackgroundService` 是实现长期运行后台任务的基类。

### 基本后台服务

```go
type MyBackgroundService struct {
    *hosting.BackgroundService
    logger *Logger
}

func NewMyBackgroundService(logger *Logger) hosting.IHostedService {
    s := &MyBackgroundService{
        BackgroundService: hosting.NewBackgroundService(),
        logger:           logger,
    }
    s.SetExecuteFunc(s.ExecuteAsync)
    return s
}

func (s *MyBackgroundService) ExecuteAsync(ctx context.Context) error {
    s.logger.Info("Background service started")
    
    for {
        select {
        case <-s.StoppingToken():
            s.logger.Info("Stopping background service")
            return nil
            
        case <-ctx.Done():
            return ctx.Err()
            
        case <-time.After(10 * time.Second):
            // 执行周期性任务
            if err := s.doWork(); err != nil {
                s.logger.Error("Work failed: %v", err)
                // 决定是否返回错误或继续
            }
        }
    }
}

func (s *MyBackgroundService) doWork() error {
    s.logger.Info("Performing background work")
    // 实际工作逻辑
    return nil
}
```

### 定时任务服务

```go
type ScheduledTaskService struct {
    *hosting.BackgroundService
    config *ScheduleConfig
    logger *Logger
}

type ScheduleConfig struct {
    IntervalSeconds int
    Enabled         bool
}

func NewScheduledTaskService(
    config *ScheduleConfig,
    logger *Logger,
) hosting.IHostedService {
    s := &ScheduledTaskService{
        BackgroundService: hosting.NewBackgroundService(),
        config:           config,
        logger:           logger,
    }
    s.SetExecuteFunc(s.ExecuteAsync)
    return s
}

func (s *ScheduledTaskService) ExecuteAsync(ctx context.Context) error {
    if !s.config.Enabled {
        s.logger.Info("Scheduled task is disabled")
        <-s.StoppingToken()
        return nil
    }
    
    ticker := time.NewTicker(time.Duration(s.config.IntervalSeconds) * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            s.executeTask()
            
        case <-s.StoppingToken():
            s.logger.Info("Scheduled task stopping")
            return nil
            
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

func (s *ScheduledTaskService) executeTask() {
    s.logger.Info("Executing scheduled task at %v", time.Now())
    // 任务逻辑
}
```

### 消息队列处理服务

```go
type MessageProcessorService struct {
    *hosting.BackgroundService
    queue  *MessageQueue
    logger *Logger
}

func NewMessageProcessorService(
    queue *MessageQueue,
    logger *Logger,
) hosting.IHostedService {
    s := &MessageProcessorService{
        BackgroundService: hosting.NewBackgroundService(),
        queue:            queue,
        logger:           logger,
    }
    s.SetExecuteFunc(s.ExecuteAsync)
    return s
}

func (s *MessageProcessorService) ExecuteAsync(ctx context.Context) error {
    s.logger.Info("Message processor started")
    
    messageCh := s.queue.Subscribe()
    
    for {
        select {
        case msg := <-messageCh:
            if err := s.processMessage(msg); err != nil {
                s.logger.Error("Failed to process message: %v", err)
            }
            
        case <-s.StoppingToken():
            s.logger.Info("Message processor stopping")
            s.queue.Unsubscribe(messageCh)
            return nil
            
        case <-ctx.Done():
            s.queue.Unsubscribe(messageCh)
            return ctx.Err()
        }
    }
}

func (s *MessageProcessorService) processMessage(msg *Message) error {
    s.logger.Info("Processing message: %s", msg.ID)
    // 处理消息
    return nil
}
```

## 环境管理

### 环境类型

CSGO 支持多个标准环境：

```go
type Environment struct {
    environmentName string
}

// 环境名称
func (e *Environment) Name() string

// 环境检查
func (e *Environment) IsDevelopment() bool
func (e *Environment) IsStaging() bool
func (e *Environment) IsProduction() bool
func (e *Environment) IsEnvironment(name string) bool
```

### 设置环境

**通过环境变量：**

```bash
export CSGO_ENVIRONMENT=Production
```

**通过命令行：**

```bash
./myapp --environment=Production
```

**在代码中：**

```go
env := hosting.NewEnvironment()
// 默认读取 CSGO_ENVIRONMENT 环境变量
```

### 按环境配置

```go
builder := hosting.CreateDefaultBuilder()

// 根据环境注册不同服务
if builder.Environment.IsDevelopment() {
    builder.Services.AddSingleton(NewMockPaymentService)
    builder.Services.AddSingleton(NewInMemoryCache)
} else {
    builder.Services.AddSingleton(NewRealPaymentService)
    builder.Services.AddSingleton(NewRedisCache)
}

// 根据环境加载配置
if builder.Environment.IsProduction() {
    // 生产环境特殊配置
    builder.ConfigureAppConfiguration(func(config configuration.IConfigurationBuilder) {
        config.AddJsonFile("production-secrets.json", false, false)
    })
}
```

## 优雅关闭

### 关闭信号处理

主机自动处理以下关闭信号：

- `SIGINT` (Ctrl+C)
- `SIGTERM` (kill 命令)

```go
host := builder.Build()

// Run() 会自动监听关闭信号
host.Run()

// 收到信号后：
// 1. 触发 ApplicationStopping 事件
// 2. 停止所有托管服务（逆序）
// 3. 触发 ApplicationStopped 事件
```

### 自定义关闭逻辑

```go
type CleanupService struct {
    db     *Database
    cache  *Cache
    logger *Logger
}

func NewCleanupService(
    db *Database,
    cache *Cache,
    logger *Logger,
    lifetime hosting.IHostApplicationLifetime,
) *CleanupService {
    s := &CleanupService{
        db:     db,
        cache:  cache,
        logger: logger,
    }
    
    // 注册关闭回调
    go func() {
        <-lifetime.ApplicationStopping()
        s.cleanup()
    }()
    
    return s
}

func (s *CleanupService) cleanup() {
    s.logger.Info("Performing cleanup")
    
    // 刷新缓存
    if err := s.cache.Flush(); err != nil {
        s.logger.Error("Failed to flush cache: %v", err)
    }
    
    // 关闭数据库连接
    if err := s.db.Close(); err != nil {
        s.logger.Error("Failed to close database: %v", err)
    }
    
    s.logger.Info("Cleanup completed")
}
```

### 关闭超时

```go
// 主机默认有 30 秒的关闭超时
// 可以在主机配置中修改：

builder.ConfigureHostConfiguration(func(config configuration.IConfigurationBuilder) {
    config.AddInMemoryCollection(map[string]string{
        "Host:ShutdownTimeout": "60", // 60 秒
    })
})
```

### 手动触发关闭

```go
type AdminService struct {
    lifetime hosting.IHostApplicationLifetime
}

func (s *AdminService) Shutdown() {
    // 触发应用关闭
    s.lifetime.StopApplication()
}
```

## 最佳实践

### 1. 服务依赖顺序

按依赖顺序注册托管服务：

```go
// ✅ 推荐：基础服务先注册
builder.Services.AddHostedService(NewDatabaseInitService)   // 数据库初始化
builder.Services.AddHostedService(NewCacheWarmupService)    // 缓存预热
builder.Services.AddHostedService(NewApiService)            // API 服务

// 关闭时会逆序停止，确保依赖正确
```

### 2. 错误处理

在托管服务中正确处理错误：

```go
func (s *MyService) StartAsync(ctx context.Context) error {
    // ❌ 不要在 StartAsync 中阻塞
    // time.Sleep(10 * time.Second)
    
    // ✅ 启动后台 goroutine
    go s.backgroundWork()
    
    return nil
}

func (s *MyService) backgroundWork() {
    defer func() {
        if r := recover(); r != nil {
            s.logger.Error("Service panic: %v", r)
        }
    }()
    
    for {
        select {
        case <-s.stopCh:
            return
        default:
            if err := s.doWork(); err != nil {
                s.logger.Error("Work failed: %v", err)
                // 决定是否继续或停止服务
            }
        }
    }
}
```

### 3. 资源清理

确保资源被正确清理：

```go
type ResourceService struct {
    *hosting.BackgroundService
    conn     *Connection
    cleanedUp bool
    mu        sync.Mutex
}

func (s *ResourceService) StopAsync(ctx context.Context) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    if s.cleanedUp {
        return nil
    }
    
    // 清理资源
    if s.conn != nil {
        s.conn.Close()
    }
    
    s.cleanedUp = true
    return nil
}
```

### 4. 超时处理

为长时间操作设置超时：

```go
func (s *DataService) StopAsync(ctx context.Context) error {
    // 创建带超时的上下文
    stopCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    
    // 执行清理操作
    done := make(chan error, 1)
    go func() {
        done <- s.cleanup()
    }()
    
    select {
    case err := <-done:
        return err
    case <-stopCtx.Done():
        return fmt.Errorf("cleanup timeout")
    }
}
```

### 5. 日志记录

在生命周期事件中记录日志：

```go
func (s *MyService) StartAsync(ctx context.Context) error {
    s.logger.Info("Service starting: %s", s.Name)
    
    if err := s.initialize(); err != nil {
        s.logger.Error("Service start failed: %v", err)
        return err
    }
    
    s.logger.Info("Service started successfully: %s", s.Name)
    return nil
}

func (s *MyService) StopAsync(ctx context.Context) error {
    s.logger.Info("Service stopping: %s", s.Name)
    
    if err := s.shutdown(); err != nil {
        s.logger.Error("Service stop failed: %v", err)
        return err
    }
    
    s.logger.Info("Service stopped successfully: %s", s.Name)
    return nil
}
```

### 6. 健康检查

实现健康检查以监控服务状态：

```go
type HealthCheckService struct {
    *hosting.BackgroundService
    db     *Database
    cache  *Cache
    healthy bool
    mu      sync.RWMutex
}

func (s *HealthCheckService) ExecuteAsync(ctx context.Context) error {
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
    healthy := true
    
    // 检查数据库
    if err := s.db.Ping(); err != nil {
        healthy = false
    }
    
    // 检查缓存
    if err := s.cache.Ping(); err != nil {
        healthy = false
    }
    
    s.mu.Lock()
    s.healthy = healthy
    s.mu.Unlock()
}

func (s *HealthCheckService) IsHealthy() bool {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.healthy
}
```

## 完整示例

### Worker Service 应用

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/gocrud/csgo/configuration"
    "github.com/gocrud/csgo/hosting"
)

// 配置
type WorkerOptions struct {
    IntervalSeconds int  `json:"intervalSeconds"`
    Enabled         bool `json:"enabled"`
}

// 日志服务
type Logger struct{}

func NewLogger() *Logger {
    return &Logger{}
}

func (l *Logger) Info(format string, args ...interface{}) {
    fmt.Printf("[INFO] "+format+"\n", args...)
}

// Worker 服务
type DataProcessorWorker struct {
    *hosting.BackgroundService
    options *WorkerOptions
    logger  *Logger
}

func NewDataProcessorWorker(
    options configuration.IOptions[WorkerOptions],
    logger *Logger,
) hosting.IHostedService {
    w := &DataProcessorWorker{
        BackgroundService: hosting.NewBackgroundService(),
        options:          options.Value(),
        logger:           logger,
    }
    w.SetExecuteFunc(w.ExecuteAsync)
    return w
}

func (w *DataProcessorWorker) ExecuteAsync(ctx context.Context) error {
    if !w.options.Enabled {
        w.logger.Info("Worker is disabled")
        <-w.StoppingToken()
        return nil
    }
    
    w.logger.Info("Worker started with interval: %d seconds", w.options.IntervalSeconds)
    
    ticker := time.NewTicker(time.Duration(w.options.IntervalSeconds) * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            w.processData()
            
        case <-w.StoppingToken():
            w.logger.Info("Worker stopping")
            return nil
            
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

func (w *DataProcessorWorker) processData() {
    w.logger.Info("Processing data at %v", time.Now().Format(time.RFC3339))
    // 数据处理逻辑
}

func main() {
    // 创建主机
    builder := hosting.CreateDefaultBuilder()
    
    // 注册选项
    builder.Services.AddSingleton(func() configuration.IOptions[WorkerOptions] {
        opts := &WorkerOptions{
            IntervalSeconds: 10,
            Enabled:         true,
        }
        builder.Configuration.Bind("Worker", opts)
        return configuration.NewOptions(opts)
    })
    
    // 注册服务
    builder.Services.AddSingleton(NewLogger)
    builder.Services.AddHostedService(NewDataProcessorWorker)
    
    // 构建并运行
    host := builder.Build()
    
    fmt.Println("=== Worker Service Starting ===")
    fmt.Println("Press Ctrl+C to stop")
    
    if err := host.Run(); err != nil {
        fmt.Printf("Host error: %v\n", err)
    }
    
    fmt.Println("=== Worker Service Stopped ===")
}
```

**appsettings.json：**

```json
{
  "Worker": {
    "IntervalSeconds": 10,
    "Enabled": true
  }
}
```

## 与 .NET 对比

| .NET | CSGO | 说明 |
|------|-----|------|
| `IHostBuilder` | `hosting.IHostBuilder` | 主机构建器接口 |
| `Host.CreateDefaultBuilder()` | `hosting.CreateDefaultBuilder()` | 创建默认主机 |
| `IHost` | `hosting.IHost` | 主机接口 |
| `IHostedService` | `hosting.IHostedService` | 托管服务接口 |
| `BackgroundService` | `hosting.BackgroundService` | 后台服务基类 |
| `IHostApplicationLifetime` | `hosting.IHostApplicationLifetime` | 生命周期接口 |
| `IHostEnvironment` | `hosting.IHostEnvironment` | 环境接口 |

## 相关资源

- [快速开始](../getting-started.md) - 开始使用 CSGO
- [配置管理](configuration.md) - 配置系统详解
- [依赖注入](dependency-injection.md) - DI 系统
- [Worker Service 示例](../../examples/worker_service/) - 完整示例

---

**下一步**: 查看 [配置管理指南](configuration.md) 了解如何配置应用。

