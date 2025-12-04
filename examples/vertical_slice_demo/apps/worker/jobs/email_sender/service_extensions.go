package email_sender

import (
	"github.com/gocrud/csgo/di"
)

// AddEmailSenderJob 注册邮件发送任务
func AddEmailSenderJob(services di.IServiceCollection) {
	services.AddHostedService(NewEmailSenderJob)
}

