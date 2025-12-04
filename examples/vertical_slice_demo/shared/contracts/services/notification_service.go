package services

// INotificationService 通知服务接口
type INotificationService interface {
	// SendEmail 发送邮件
	SendEmail(to, subject, body string) error
	
	// SendSMS 发送短信
	SendSMS(phone, message string) error
	
	// SendPush 发送推送通知
	SendPush(userID int64, title, message string) error
}

// EmailTemplate 邮件模板
type EmailTemplate string

const (
	// TemplateWelcome 欢迎邮件
	TemplateWelcome EmailTemplate = "welcome"
	// TemplateOrderConfirm 订单确认邮件
	TemplateOrderConfirm EmailTemplate = "order_confirm"
	// TemplatePasswordReset 密码重置邮件
	TemplatePasswordReset EmailTemplate = "password_reset"
)

