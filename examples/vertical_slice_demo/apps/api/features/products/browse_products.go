package products

import (
	"vertical_slice_demo/shared/contracts/repositories"
	"vertical_slice_demo/shared/domain"

	"github.com/gocrud/csgo/web"
)

// BrowseProductsHandler 浏览商品处理器（C端视角）
type BrowseProductsHandler struct {
	productRepo repositories.IProductRepository
}

// NewBrowseProductsHandler 浏览商品处理器构造函数
func NewBrowseProductsHandler(productRepo repositories.IProductRepository) *BrowseProductsHandler {
	return &BrowseProductsHandler{productRepo: productRepo}
}

// ProductListItem C端商品列表项（不显示某些管理信息）
type ProductListItem struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	InStock     bool    `json:"in_stock"`
}

// BrowseProductsResponse 浏览商品响应
type BrowseProductsResponse struct {
	Products []ProductListItem `json:"products"`
	Total    int               `json:"total"`
}

// Handle 处理浏览商品请求
func (h *BrowseProductsHandler) Handle(c *web.HttpContext) web.IActionResult {
	// 获取分页参数
	p := c.Params()
	offset := p.QueryInt("offset").NonNegative().ValueOr(0)
	limit := p.QueryInt("limit").Range(1, 100).ValueOr(20)

	// 查询商品列表（C端只显示 active 状态的商品）
	products, err := h.productRepo.List(offset, limit, "active")
	if err != nil {
		return c.InternalError("查询商品列表失败")
	}

	// 转换为C端视角的响应格式
	items := make([]ProductListItem, len(products))
	for i, product := range products {
		items[i] = ProductListItem{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			InStock:     product.Stock > 0,
		}
	}

	response := &BrowseProductsResponse{
		Products: items,
		Total:    len(items),
	}

	return c.Ok(response)
}

// GetProductDetailHandler 获取商品详情处理器
type GetProductDetailHandler struct {
	productRepo repositories.IProductRepository
}

// NewGetProductDetailHandler 获取商品详情处理器构造函数
func NewGetProductDetailHandler(productRepo repositories.IProductRepository) *GetProductDetailHandler {
	return &GetProductDetailHandler{productRepo: productRepo}
}

// ProductDetail 商品详情
type ProductDetail struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	InStock     bool    `json:"in_stock"`
}

// Handle 处理获取商品详情请求
func (h *GetProductDetailHandler) Handle(c *web.HttpContext) web.IActionResult {
	// 获取商品 ID
	id := c.Params().PathInt64("id").Positive().Value()
	if err := c.Params().Check(); err != nil {
		return err
	}

	// 查询商品
	product, err := h.productRepo.GetByID(id)
	if err != nil {
		return c.NotFound("商品不存在")
	}

	// C端只能查看 active 状态的商品
	if product.Status != "active" {
		return c.NotFound("商品不存在")
	}

	detail := &ProductDetail{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		InStock:     product.Stock > 0,
	}

	return c.Ok(detail)
}

// toProductListItem 转换为列表项
func toProductListItem(product *domain.Product) ProductListItem {
	return ProductListItem{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		InStock:     product.Stock > 0,
	}
}
