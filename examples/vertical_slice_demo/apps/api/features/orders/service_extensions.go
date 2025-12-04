package orders

import (
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/web"
)

// AddOrderFeature 注册订单功能
func AddOrderFeature(services di.IServiceCollection) {
	// 注册处理器
	services.AddSingleton(NewCreateOrderHandler)
	services.AddSingleton(NewMyOrdersHandler)
	services.AddSingleton(NewPayOrderHandler)

	// 注册控制器
	web.AddController(services, NewOrderController)
}

