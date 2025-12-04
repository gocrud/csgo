package auth

import (
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/web"
)

// AddAuthFeature 注册认证功能
func AddAuthFeature(services di.IServiceCollection) {
	// 注册处理器
	services.AddSingleton(NewLoginHandler)
	services.AddSingleton(NewRegisterHandler)

	// 注册控制器
	web.AddController(services, NewAuthController)
}

