package orders

import (
	"github.com/gocrud/csgo/web"
	"vertical_slice_demo/shared/contracts/repositories"
	"vertical_slice_demo/shared/domain"
)

// CreateOrderHandler 创建订单处理器
type CreateOrderHandler struct {
	orderRepo   repositories.IOrderRepository
	productRepo repositories.IProductRepository
}

// NewCreateOrderHandler 创建订单处理器构造函数
func NewCreateOrderHandler(
	orderRepo repositories.IOrderRepository,
	productRepo repositories.IProductRepository,
) *CreateOrderHandler {
	return &CreateOrderHandler{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

// OrderItemRequest 订单项请求
type OrderItemRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,gt=0"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

// Handle 处理创建订单请求
func (h *CreateOrderHandler) Handle(c *web.HttpContext) web.IActionResult {
	var req CreateOrderRequest
	if err := c.MustBindJSON(&req); err != nil {
		return err
	}

	// 简化实现：从上下文获取用户 ID（实际项目中应从 JWT token 中解析）
	userID := int64(1) // 假设当前用户 ID 为 1

	// 计算总价并创建订单项
	var totalPrice float64
	orderItems := make([]domain.OrderItem, 0, len(req.Items))

	for _, item := range req.Items {
		// 获取商品信息
		product, err := h.productRepo.GetByID(item.ProductID)
		if err != nil {
			return c.BadRequest("商品不存在")
		}

		// 检查库存
		if product.Stock < item.Quantity {
			return c.BadRequest("商品库存不足")
		}

		// 创建订单项
		orderItem := domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}
		orderItems = append(orderItems, orderItem)

		// 累计总价
		totalPrice += product.Price * float64(item.Quantity)

		// 更新库存
		_ = h.productRepo.UpdateStock(item.ProductID, -item.Quantity)
	}

	// 创建订单
	order := &domain.Order{
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     "pending",
		Items:      orderItems,
	}

	if err := h.orderRepo.Create(order); err != nil {
		return c.InternalError("创建订单失败")
	}

	return c.Created(order)
}

