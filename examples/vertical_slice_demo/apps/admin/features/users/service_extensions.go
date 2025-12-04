package users

import (
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/web"
)

// AddUserFeature 注册用户管理功能
func AddUserFeature(services di.IServiceCollection) {
	// 注册处理器
	services.AddSingleton(NewCreateUserHandler)
	services.AddSingleton(NewListUsersHandler)
	services.AddSingleton(NewUpdateUserHandler)

	// 注册控制器
	web.AddController(services, NewUserController)
}

