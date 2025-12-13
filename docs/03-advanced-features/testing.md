# 测试

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

测试确保代码质量和可靠性。

## 单元测试

```go
func TestUserService_GetUser(t *testing.T) {
    // 创建mock依赖
    mockRepo := NewMockUserRepository()
    service := NewUserService(mockRepo)
    
    // 测试
    user, err := service.GetUser(1)
    
    // 断言
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if user.ID != 1 {
        t.Errorf("expected ID 1, got %d", user.ID)
    }
}
```

## 集成测试

```go
func TestAPI_CreateUser(t *testing.T) {
    builder := web.CreateBuilder()
    builder.Services.Add(NewUserService)
    app := builder.Build()
    
    // 测试API端点
    // ...
}
```

## 下一步

继续学习：[性能优化](performance.md) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

