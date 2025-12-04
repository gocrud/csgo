package products

import (
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/web"
)

// AddProductFeature 注册商品浏览功能（C端）
func AddProductFeature(services di.IServiceCollection) {
	// 注册处理器
	services.AddSingleton(NewBrowseProductsHandler)
	services.AddSingleton(NewGetProductDetailHandler)

	// 注册控制器
	web.AddController(services, NewProductController)
}

