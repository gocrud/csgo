package products

import (
	"github.com/gocrud/csgo/web"
)

// ProductController 商品控制器（C端）
type ProductController struct {
	browseHandler *BrowseProductsHandler
	detailHandler *GetProductDetailHandler
}

// NewProductController 创建商品控制器
func NewProductController(
	browseHandler *BrowseProductsHandler,
	detailHandler *GetProductDetailHandler,
) *ProductController {
	return &ProductController{
		browseHandler: browseHandler,
		detailHandler: detailHandler,
	}
}

// MapRoutes 映射路由
func (ctrl *ProductController) MapRoutes(app *web.WebApplication) {
	products := app.MapGroup("/api/products")
	products.MapGet("", ctrl.browseHandler.Handle)
	products.MapGet("/:id", ctrl.detailHandler.Handle)
}

