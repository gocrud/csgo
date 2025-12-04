package products

import (
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/web"
)

// AddProductFeature 注册商品管理功能
func AddProductFeature(services di.IServiceCollection) {
	// 注册处理器
	services.AddSingleton(NewCreateProductHandler)
	services.AddSingleton(NewListProductsHandler)

	// 注册控制器
	web.AddController(services, NewProductController)
}

