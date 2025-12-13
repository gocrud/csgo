# 结构化日志

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

结构化日志帮助你更好地监控应用。

## 完整文档

关于日志系统的完整详细文档，请查看：

👉 **[日志系统完整文档](../../logging/README.md)**

## 快速示例

```go
type UserService struct {
    logger logging.ILogger
}

func NewUserService(factory logging.ILoggerFactory) *UserService {
    return &UserService{
        logger: logging.GetLogger[UserService](factory),
    }
}

func (s *UserService) GetUser(id int) {
    s.logger.LogInformation("Getting user with id: %d", id)
    // ...
    s.logger.LogError(err, "Failed to get user")
}
```

## 下一步

继续学习：[测试](testing.md) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

