package payment

import "github.com/gocrud/csgo/di"

// AddPaymentService 注册支付服务
func AddPaymentService(services di.IServiceCollection) {
	services.AddSingleton(NewPaymentService)
}

