package email_sender

import (
	"context"
	"fmt"
	"time"

	"github.com/gocrud/csgo/hosting"
)

// EmailSenderJob 邮件发送任务
type EmailSenderJob struct {
	stopChan chan struct{}
}

// NewEmailSenderJob 创建邮件发送任务
func NewEmailSenderJob() hosting.IHostedService {
	return &EmailSenderJob{
		stopChan: make(chan struct{}),
	}
}

// StartAsync 启动任务
func (j *EmailSenderJob) StartAsync(ctx context.Context) error {
	fmt.Println("邮件发送任务启动...")

	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				j.sendEmails()
			case <-j.stopChan:
				fmt.Println("邮件发送任务停止")
				return
			case <-ctx.Done():
				fmt.Println("邮件发送任务被取消")
				return
			}
		}
	}()

	return nil
}

// StopAsync 停止任务
func (j *EmailSenderJob) StopAsync(ctx context.Context) error {
	fmt.Println("正在停止邮件发送任务...")
	close(j.stopChan)
	return nil
}

// sendEmails 发送邮件
func (j *EmailSenderJob) sendEmails() {
	fmt.Printf("[%s] 检查待发送邮件...\n", time.Now().Format("15:04:05"))
	
	// 实际项目中这里会执行真实的邮件发送逻辑
	// 比如：查询待发送邮件队列、发送邮件、更新状态等
	
	// 模拟处理时间
	time.Sleep(1 * time.Second)
	
	fmt.Printf("[%s] 邮件发送完成\n", time.Now().Format("15:04:05"))
}

