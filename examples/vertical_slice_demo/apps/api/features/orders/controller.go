package orders

import (
	"github.com/gocrud/csgo/web"
)

// OrderController 订单控制器
type OrderController struct {
	createHandler   *CreateOrderHandler
	myOrdersHandler *MyOrdersHandler
	payHandler      *PayOrderHandler
}

// NewOrderController 创建订单控制器
func NewOrderController(
	createHandler *CreateOrderHandler,
	myOrdersHandler *MyOrdersHandler,
	payHandler *PayOrderHandler,
) *OrderController {
	return &OrderController{
		createHandler:   createHandler,
		myOrdersHandler: myOrdersHandler,
		payHandler:      payHandler,
	}
}

// MapRoutes 映射路由
func (ctrl *OrderController) MapRoutes(app *web.WebApplication) {
	orders := app.MapGroup("/api/orders")
	orders.MapPost("", ctrl.createHandler.Handle)
	orders.MapGet("/my", ctrl.myOrdersHandler.Handle)
	orders.MapPost("/:id/pay", ctrl.payHandler.Handle)
}

