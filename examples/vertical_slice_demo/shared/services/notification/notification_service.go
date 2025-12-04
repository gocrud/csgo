package notification

import (
	"fmt"
	"time"

	"vertical_slice_demo/shared/contracts/services"
)

// NotificationService 通知服务实现
type NotificationService struct {
	// 实际项目中这里会有邮件服务器配置、短信服务配置等
}

// NewNotificationService 创建通知服务
func NewNotificationService() services.INotificationService {
	return &NotificationService{}
}

// SendEmail 发送邮件
func (s *NotificationService) SendEmail(to, subject, body string) error {
	// 实际项目中这里会调用真实的邮件服务（如 SMTP、SendGrid、阿里云邮件推送等）
	fmt.Printf("[EMAIL] 发送邮件\n")
	fmt.Printf("  收件人: %s\n", to)
	fmt.Printf("  主题: %s\n", subject)
	fmt.Printf("  内容: %s\n", body)
	fmt.Printf("  时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	
	// 模拟发送延迟
	time.Sleep(100 * time.Millisecond)
	
	return nil
}

// SendSMS 发送短信
func (s *NotificationService) SendSMS(phone, message string) error {
	// 实际项目中这里会调用真实的短信服务（如阿里云短信、腾讯云短信等）
	fmt.Printf("[SMS] 发送短信\n")
	fmt.Printf("  手机号: %s\n", phone)
	fmt.Printf("  内容: %s\n", message)
	fmt.Printf("  时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	
	// 模拟发送延迟
	time.Sleep(100 * time.Millisecond)
	
	return nil
}

// SendPush 发送推送通知
func (s *NotificationService) SendPush(userID int64, title, message string) error {
	// 实际项目中这里会调用推送服务（如极光推送、个推等）
	fmt.Printf("[PUSH] 发送推送通知\n")
	fmt.Printf("  用户ID: %d\n", userID)
	fmt.Printf("  标题: %s\n", title)
	fmt.Printf("  内容: %s\n", message)
	fmt.Printf("  时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	
	// 模拟发送延迟
	time.Sleep(100 * time.Millisecond)
	
	return nil
}

// SendWelcomeEmail 发送欢迎邮件（业务方法示例）
func (s *NotificationService) SendWelcomeEmail(email, name string) error {
	subject := "欢迎加入我们！"
	body := fmt.Sprintf("亲爱的 %s，\n\n欢迎注册我们的平台！我们很高兴您的加入。\n\n如有任何问题，请随时联系我们。\n\n祝好！", name)
	return s.SendEmail(email, subject, body)
}

// SendOrderConfirmEmail 发送订单确认邮件（业务方法示例）
func (s *NotificationService) SendOrderConfirmEmail(email string, orderID int64, amount float64) error {
	subject := "订单确认"
	body := fmt.Sprintf("您的订单 #%d 已确认！\n\n订单金额: ¥%.2f\n\n我们会尽快为您安排发货。", orderID, amount)
	return s.SendEmail(email, subject, body)
}

