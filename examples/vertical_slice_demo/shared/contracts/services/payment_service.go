package services

// IPaymentService 支付服务接口
type IPaymentService interface {
	// CreatePayment 创建支付
	CreatePayment(orderID int64, amount float64, method PaymentMethod) (*PaymentResult, error)
	
	// QueryPayment 查询支付状态
	QueryPayment(paymentID string) (*PaymentStatus, error)
	
	// RefundPayment 退款
	RefundPayment(paymentID string, amount float64, reason string) error
}

// PaymentMethod 支付方式
type PaymentMethod string

const (
	// PaymentMethodAlipay 支付宝
	PaymentMethodAlipay PaymentMethod = "alipay"
	// PaymentMethodWechat 微信支付
	PaymentMethodWechat PaymentMethod = "wechat"
	// PaymentMethodCard 银行卡
	PaymentMethodCard PaymentMethod = "card"
)

// PaymentResult 支付结果
type PaymentResult struct {
	PaymentID   string        `json:"payment_id"`   // 支付ID
	OrderID     int64         `json:"order_id"`     // 订单ID
	Amount      float64       `json:"amount"`       // 支付金额
	Method      PaymentMethod `json:"method"`       // 支付方式
	Status      string        `json:"status"`       // 支付状态
	RedirectURL string        `json:"redirect_url"` // 支付跳转URL（如果需要）
}

// PaymentStatus 支付状态
type PaymentStatus struct {
	PaymentID string  `json:"payment_id"` // 支付ID
	Status    string  `json:"status"`     // 状态：pending, success, failed
	Amount    float64 `json:"amount"`     // 金额
	PaidAt    string  `json:"paid_at"`    // 支付时间
}

