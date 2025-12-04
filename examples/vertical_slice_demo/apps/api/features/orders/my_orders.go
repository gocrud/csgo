package orders

import (
	"strconv"

	"github.com/gocrud/csgo/web"
	"vertical_slice_demo/shared/contracts/repositories"
	"vertical_slice_demo/shared/domain"
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
	offset := 0
	limit := 20

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil {
			offset = v
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

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

