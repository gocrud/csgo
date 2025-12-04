package products

import (
	"strconv"

	"github.com/gocrud/csgo/web"
	"vertical_slice_demo/shared/contracts/repositories"
	"vertical_slice_demo/shared/domain"
)

// ListProductsHandler 商品列表处理器
type ListProductsHandler struct {
	productRepo repositories.IProductRepository
}

// NewListProductsHandler 商品列表处理器构造函数
func NewListProductsHandler(productRepo repositories.IProductRepository) *ListProductsHandler {
	return &ListProductsHandler{productRepo: productRepo}
}

// ListProductsResponse 商品列表响应
type ListProductsResponse struct {
	Products []*domain.Product `json:"products"`
	Total    int               `json:"total"`
}

// Handle 处理商品列表请求
func (h *ListProductsHandler) Handle(c *web.HttpContext) web.IActionResult {
	// 获取分页参数
	offset := 0
	limit := 20
	status := c.Query("status") // 可选的状态过滤

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

	// 查询商品列表
	products, err := h.productRepo.List(offset, limit, status)
	if err != nil {
		return c.InternalError("查询商品列表失败")
	}

	response := &ListProductsResponse{
		Products: products,
		Total:    len(products),
	}

	return c.Ok(response)
}

