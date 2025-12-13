# 自定义中间件

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

中间件是处理HTTP请求的强大工具。

## 完整文档

关于中间件的完整详细文档，请查看：

👉 **[Web框架完整文档 - 中间件部分](../../web/README.md#中间件)**

## 快速示例

```go
func loggingMiddleware(c *gin.Context) {
    start := time.Now()
    c.Next()
    latency := time.Since(start)
    fmt.Printf("[%s] %s %v\n", c.Request.Method, c.Request.URL.Path, latency)
}

app.Use(loggingMiddleware)
```

## 下一步

继续学习：[后台服务](background-services.md) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

