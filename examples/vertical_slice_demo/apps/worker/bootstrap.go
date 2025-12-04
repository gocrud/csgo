package worker

import (
	"github.com/gocrud/csgo/configuration"
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/hosting"
	"vertical_slice_demo/apps/worker/jobs/email_sender"
	"vertical_slice_demo/apps/worker/jobs/order_sync"
	"vertical_slice_demo/configs"
	"vertical_slice_demo/shared/infrastructure/database"
	"vertical_slice_demo/shared/repositories"
)

// ConfigureServices 配置服务
func ConfigureServices(services di.IServiceCollection) {
	// 构建配置（使用框架的配置系统）
	config := configuration.NewConfigurationBuilder().
		AddJsonFile("configs/config.dev.json", true, false).
		AddEnvironmentVariables("APP_").
		Build()

	// 绑定配置到结构体
	var appConfig configs.Config
	config.Bind("", &appConfig)

	// 注册配置到 DI 容器
	services.AddSingleton(func() configuration.IConfiguration {
		return config
	})
	services.AddSingleton(func() *configs.Config {
		return &appConfig
	})

	// 注册共享基础设施
	database.AddDatabase(services)

	// 注册共享仓储
	repositories.AddRepositories(services)

	// 注册后台任务
	order_sync.AddOrderSyncJob(services)
	email_sender.AddEmailSenderJob(services)
}

// Bootstrap 启动 Worker 服务
func Bootstrap() hosting.IHost {
	return hosting.CreateDefaultBuilder().
		ConfigureServices(ConfigureServices).
		Build()
}

