package admin

import (
	"vertical_slice_demo/apps/admin/features/products"
	"vertical_slice_demo/apps/admin/features/users"
	"vertical_slice_demo/apps/admin/middlewares"
	"vertical_slice_demo/configs"
	"vertical_slice_demo/shared/infrastructure/cache"
	"vertical_slice_demo/shared/infrastructure/database"
	"vertical_slice_demo/shared/repositories"
	"vertical_slice_demo/shared/services/notification"
	"vertical_slice_demo/shared/services/payment"

	"github.com/gocrud/csgo/configuration"
	"github.com/gocrud/csgo/web"
)

// Bootstrap 启动管理端应用
func Bootstrap() *web.WebApplication {
	builder := web.CreateBuilder()

	// 构建配置（使用框架的配置系统）
	config := configuration.NewConfigurationBuilder().
		AddJsonFile("configs/config.dev.json", true, false).
		AddEnvironmentVariables("APP_").
		Build()

	// 绑定配置到结构体
	var appConfig configs.Config
	config.Bind("", &appConfig)

	// 注册配置到 DI 容器
	builder.Services.AddSingleton(func() configuration.IConfiguration {
		return config
	})
	builder.Services.AddSingleton(func() *configs.Config {
		return &appConfig
	})

	// 注册共享基础设施
	database.AddDatabase(builder.Services)
	cache.AddCache(builder.Services)

	// 注册共享仓储
	repositories.AddRepositories(builder.Services)

	// 注册共享服务
	notification.AddNotificationService(builder.Services)
	payment.AddPaymentService(builder.Services)

	// 注册管理端功能切片
	users.AddUserFeature(builder.Services)
	products.AddProductFeature(builder.Services)

	// 构建应用
	app := builder.Build()

	// 管理端专属中间件
	app.Use(middlewares.AdminAuthMiddleware())

	// 根路由
	app.MapGet("/", func(c *web.HttpContext) web.IActionResult {
		return c.Ok(map[string]string{
			"message": "管理端 API",
			"version": "1.0.0",
		})
	})

	// 映射所有控制器
	app.MapControllers()

	return app
}
