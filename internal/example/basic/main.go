package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocrud/csgo/web"
)

func main() {
	builder := web.CreateBuilder()

	app := builder.Build()

	// 带认证和日志的 API 组
	api := app.Group("/api",
		// 日志中间件
		func(hc *web.HttpContext) web.IActionResult {
			start := time.Now()
			path := hc.RawCtx().Request.URL.Path

			// 继续执行后续处理器
			hc.RawCtx().Next()

			// 后续处理器执行完毕后记录日志
			duration := time.Since(start)
			status := hc.RawCtx().Writer.Status()
			log.Printf("[%s] %s - %d (%v)",
				hc.RawCtx().Request.Method, path, status, duration)

			return nil
		},
		// 认证中间件
		func(hc *web.HttpContext) web.IActionResult {
			token := hc.RawCtx().GetHeader("Authorization")
			if token == "" {
				return hc.Unauthorized("需要认证令牌")
			}

			// 设置用户信息
			hc.RawCtx().Set("user", "张三")
			return nil // 继续
		},
	)

	// ==================== 新的泛型参数 API 示例 ====================

	// 示例 1: 基本用法 - Path 参数
	api.GET("/users/:id", func(c *web.HttpContext) web.IActionResult {
		// 使用新的泛型 API 获取路径参数
		// Value() 返回 (值, IActionResult)，验证失败时立即返回错误
		id, err := web.Path[int](c, "id").Value()
		if err != nil {
			return err
		}

		return c.Ok(web.M{
			"message": "获取用户详情",
			"id":      id,
		})
	})

	// 示例 2: Query 参数带默认值
	api.GET("/products", func(c *web.HttpContext) web.IActionResult {
		// Query 参数，使用 Default 设置默认值
		page := web.Query[int](c, "page").Default(1)
		size := web.Query[int](c, "size").Default(10)

		// 可选的排序参数
		sort := web.Query[string](c, "sort").Default("date")

		return c.Ok(web.M{
			"page":     page,
			"size":     size,
			"sort":     sort,
			"products": []string{"产品1", "产品2", "产品3"},
		})
	})

	// 示例 3: 使用 Required 和 Custom 验证
	api.POST("/register", func(c *web.HttpContext) web.IActionResult {
		// 必填参数，带自定义验证
		// 每个 Value() 都会返回错误，验证失败时立即返回
		username, err := web.Query[string](c, "username").
			Required().
			Custom(func(v string) error {
				if len(v) < 3 || len(v) > 20 {
					return fmt.Errorf("用户名长度必须在 3-20 个字符之间")
				}
				return nil
			}).
			Value()
		if err != nil {
			return err
		}

		email, err := web.Query[string](c, "email").
			Required().
			Custom(func(v string) error {
				// 简单的邮箱验证
				if !contains(v, "@") {
					return fmt.Errorf("邮箱格式不正确")
				}
				return nil
			}).
			Value()
		if err != nil {
			return err
		}

		age, err := web.Query[int](c, "age").
			Required().
			Custom(func(v int) error {
				if v < 18 || v > 120 {
					return fmt.Errorf("年龄必须在 18-120 之间")
				}
				return nil
			}).
			Value()
		if err != nil {
			return err
		}

		return c.Ok(web.M{
			"message":  "注册成功",
			"username": username,
			"email":    email,
			"age":      age,
		})
	})

	// 示例 4: 多个类型的参数
	api.PUT("/settings/:key", func(c *web.HttpContext) web.IActionResult {
		// Path 参数
		key, err := web.Path[string](c, "key").Value()
		if err != nil {
			return err
		}

		// Query 参数 - 必填
		value, err := web.Query[string](c, "value").Required().Value()
		if err != nil {
			return err
		}

		// Header 参数 - 使用 Default 不返回错误
		version := web.Header[int](c, "X-API-Version").Default(1)

		return c.Ok(web.M{
			"key":     key,
			"value":   value,
			"version": version,
		})
	})

	// 示例 5: 手动错误处理 (使用 Get 方法)
	api.GET("/search", func(c *web.HttpContext) web.IActionResult {
		// 使用 Get() 方法手动处理错误
		keyword, err := web.Query[string](c, "keyword").
			Required().
			Custom(func(v string) error {
				if len(v) < 2 {
					return fmt.Errorf("关键词至少需要 2 个字符")
				}
				return nil
			}).
			Get()

		if err != nil {
			// 自定义错误响应
			return c.BadRequest(fmt.Sprintf("搜索失败: %v", err))
		}

		return c.Ok(web.M{
			"keyword": keyword,
			"results": []string{"结果1", "结果2"},
		})
	})

	// ==================== 用户列表 (原有示例) ====================
	api.GET("/users", func(hc *web.HttpContext) web.IActionResult {
		// 获取中间件设置的用户信息
		user, _ := hc.RawCtx().Get("user")
		return hc.Ok(web.M{
			"message": "用户列表",
			"user":    user,
		})
	})

	log.Println("服务器启动在 http://localhost:8080")
	log.Println("试试这些 API:")
	log.Println("  GET  /api/users/123")
	log.Println("  GET  /api/products?page=2&size=20")
	log.Println("  POST /api/register?username=john&email=john@example.com&age=25")
	log.Println("  GET  /api/search?keyword=test")

	app.Run()
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && indexOf(s, substr) >= 0
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
