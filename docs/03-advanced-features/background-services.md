# 后台服务

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

后台服务用于执行长时间运行的任务。

## 完整文档

关于后台服务的完整详细文档，请查看：

👉 **[主机管理完整文档 - 后台服务部分](../../hosting/README.md#ihostedservice-接口)**

## 快速示例

```go
type MyWorker struct {
    *hosting.BackgroundService
}

func NewMyWorker() *MyWorker {
    worker := &MyWorker{
        BackgroundService: hosting.NewBackgroundService(),
    }
    worker.SetExecuteFunc(worker.execute)
    return worker
}

func (w *MyWorker) execute(ctx context.Context) error {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // 执行任务
        case <-w.StoppingToken():
            return nil
        }
    }
}

// 注册
builder.Services.AddHostedService(NewMyWorker)
```

## 下一步

继续学习：[结构化日志](logging.md) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

