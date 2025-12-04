package notification

import "github.com/gocrud/csgo/di"

// AddNotificationService 注册通知服务
func AddNotificationService(services di.IServiceCollection) {
	services.AddSingleton(NewNotificationService)
}

