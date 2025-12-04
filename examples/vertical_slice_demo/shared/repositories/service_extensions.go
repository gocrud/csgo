package repositories

import "github.com/gocrud/csgo/di"

// AddRepositories 注册所有仓储
func AddRepositories(services di.IServiceCollection) {
	services.AddSingleton(NewUserRepository)
	services.AddSingleton(NewProductRepository)
	services.AddSingleton(NewOrderRepository)
}

