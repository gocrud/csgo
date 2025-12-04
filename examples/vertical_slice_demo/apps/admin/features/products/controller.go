package products

import (
	"github.com/gocrud/csgo/web"
)

// ProductController 商品控制器
type ProductController struct {
	createHandler *CreateProductHandler
	listHandler   *ListProductsHandler
}

// NewProductController 创建商品控制器
func NewProductController(
	createHandler *CreateProductHandler,
	listHandler *ListProductsHandler,
) *ProductController {
	return &ProductController{
		createHandler: createHandler,
		listHandler:   listHandler,
	}
}

// MapRoutes 映射路由
func (ctrl *ProductController) MapRoutes(app *web.WebApplication) {
	products := app.MapGroup("/api/admin/products")
	products.MapPost("", ctrl.createHandler.Handle)
	products.MapGet("", ctrl.listHandler.Handle)
}

