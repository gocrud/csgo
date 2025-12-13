# API 文档

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

自动生成Swagger文档让API更易用。

## 完整文档

关于API文档的完整详细文档，请查看：

👉 **[Swagger 集成完整文档](../../swagger/README.md)**

## 快速示例

```go
import "github.com/gocrud/csgo/swagger"

func main() {
    builder := web.CreateBuilder()
    
    // 添加Swagger
    swagger.AddSwaggerGen(builder.Services, func(opts *swagger.SwaggerGenOptions) {
        opts.Title = "我的API"
        opts.Version = "v1"
        opts.Description = "API文档"
    })
    
    app := builder.Build()
    
    // 定义路由
    app.MapGet("/api/users", getUsers)
    app.MapPost("/api/users", createUser)
    
    // 启用Swagger
    swagger.UseSwagger(app)
    swagger.UseSwaggerUI(app)
    
    app.Run()
}
```

访问：http://localhost:8080/swagger

## 下一步

继续学习：[最佳实践](best-practices.md) →

---

[← 返回目录](README.md) | [← 返回主目录](../../README.md)

