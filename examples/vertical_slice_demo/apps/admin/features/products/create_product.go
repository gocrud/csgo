package products

import (
	"github.com/gocrud/csgo/web"
	"vertical_slice_demo/shared/contracts/repositories"
	"vertical_slice_demo/shared/domain"
)

// CreateProductHandler 创建商品处理器
type CreateProductHandler struct {
	productRepo repositories.IProductRepository
}

// NewCreateProductHandler 创建商品处理器构造函数
func NewCreateProductHandler(productRepo repositories.IProductRepository) *CreateProductHandler {
	return &CreateProductHandler{productRepo: productRepo}
}

// CreateProductRequest 创建商品请求
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
	Status      string  `json:"status" binding:"required,oneof=active inactive"`
}

// Handle 处理创建商品请求
func (h *CreateProductHandler) Handle(c *web.HttpContext) web.IActionResult {
	var req CreateProductRequest
	if err := c.MustBindJSON(&req); err != nil {
		return err
	}

	// 创建商品实体
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Status:      req.Status,
	}

	// 持久化
	if err := h.productRepo.Create(product); err != nil {
		return c.InternalError("创建商品失败")
	}

	return c.Created(product)
}

