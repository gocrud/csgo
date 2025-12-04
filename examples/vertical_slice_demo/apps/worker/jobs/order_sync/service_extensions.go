package order_sync

import (
	"github.com/gocrud/csgo/di"
)

// AddOrderSyncJob 注册订单同步任务
func AddOrderSyncJob(services di.IServiceCollection) {
	services.AddHostedService(NewOrderSyncJob)
}

