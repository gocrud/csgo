package orders

import (
	"vertical_slice_demo/shared/contracts/repositories"
	"vertical_slice_demo/shared/domain"

	"github.com/gocrud/csgo/web"
)

// MyOrdersHandler 我的订单处理器
type MyOrdersHandler struct {
	orderRepo repositories.IOrderRepository
}

// NewMyOrdersHandler 我的订单处理器构造函数
func NewMyOrdersHandler(orderRepo repositories.IOrderRepository) *MyOrdersHandler {
	return &MyOrdersHandler{orderRepo: orderRepo}
}

// MyOrdersResponse 我的订单响应
type MyOrdersResponse struct {
	Orders []*domain.Order `json:"orders"`
	Total  int             `json:"total"`
}

// Handle 处理我的订单请求
func (h *MyOrdersHandler) Handle(c *web.HttpContext) web.IActionResult {
	// 简化实现：从上下文获取用户 ID（实际项目中应从 JWT token 中解析）
	userID := int64(1) // 假设当前用户 ID 为 1

	// 获取分页参数
	p := c.Params()
	offset := p.QueryInt("offset").NonNegative().ValueOr(0)
	limit := p.QueryInt("limit").Range(1, 100).ValueOr(20)

	// 查询订单列表
	orders, err := h.orderRepo.GetByUserID(userID, offset, limit)
	if err != nil {
		return c.InternalError("查询订单列表失败")
	}

	response := &MyOrdersResponse{
		Orders: orders,
		Total:  len(orders),
	}

	return c.Ok(response)
}
