package orders

import (
	"vertical_slice_demo/shared/contracts/repositories"
	"vertical_slice_demo/shared/contracts/services"

	"github.com/gocrud/csgo/web"
)

// PayOrderHandler 支付订单处理器（展示共享服务的使用）
type PayOrderHandler struct {
	orderRepo           repositories.IOrderRepository
	paymentService      services.IPaymentService
	notificationService services.INotificationService
	userRepo            repositories.IUserRepository
}

// NewPayOrderHandler 创建支付订单处理器
func NewPayOrderHandler(
	orderRepo repositories.IOrderRepository,
	paymentService services.IPaymentService,
	notificationService services.INotificationService,
	userRepo repositories.IUserRepository,
) *PayOrderHandler {
	return &PayOrderHandler{
		orderRepo:           orderRepo,
		paymentService:      paymentService,
		notificationService: notificationService,
		userRepo:            userRepo,
	}
}

// PayOrderRequest 支付订单请求
type PayOrderRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required,oneof=alipay wechat card"`
}

// Handle 处理支付订单请求
func (h *PayOrderHandler) Handle(c *web.HttpContext) web.IActionResult {
	// 获取订单 ID
	orderID := c.Params().PathInt64("id").Positive().Value()
	if err := c.Params().Check(); err != nil {
		return err
	}

	// 绑定请求
	var req PayOrderRequest
	if err := c.MustBindJSON(&req); err != nil {
		return err
	}

	// 查询订单
	order, err := h.orderRepo.GetByID(orderID)
	if err != nil {
		return c.NotFound("订单不存在")
	}

	// 检查订单状态
	if order.Status != "pending" {
		return c.BadRequest("订单状态不允许支付")
	}

	// 创建支付（使用共享的支付服务）
	paymentResult, err := h.paymentService.CreatePayment(
		order.ID,
		order.TotalPrice,
		services.PaymentMethod(req.PaymentMethod),
	)
	if err != nil {
		return c.InternalError("创建支付失败")
	}

	// 更新订单状态
	order.Status = "paying"
	_ = h.orderRepo.Update(order)

	// 发送通知（使用共享的通知服务）
	// 实际项目中应该从 JWT token 中获取用户信息
	user, _ := h.userRepo.GetByID(order.UserID)
	if user != nil {
		// 发送邮件通知
		_ = h.notificationService.SendEmail(
			user.Email,
			"订单支付确认",
			"您的订单正在支付中，请完成支付。",
		)

		// 发送推送通知
		_ = h.notificationService.SendPush(
			user.ID,
			"订单支付",
			"请完成订单支付",
		)
	}

	return c.Ok(paymentResult)
}
