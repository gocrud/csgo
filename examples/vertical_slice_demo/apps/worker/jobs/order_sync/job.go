package order_sync

import (
	"context"
	"fmt"
	"time"

	"github.com/gocrud/csgo/hosting"
	"vertical_slice_demo/shared/contracts/repositories"
)

// OrderSyncJob 订单同步任务
type OrderSyncJob struct {
	orderRepo repositories.IOrderRepository
	stopChan  chan struct{}
}

// NewOrderSyncJob 创建订单同步任务
func NewOrderSyncJob(orderRepo repositories.IOrderRepository) hosting.IHostedService {
	return &OrderSyncJob{
		orderRepo: orderRepo,
		stopChan:  make(chan struct{}),
	}
}

// StartAsync 启动任务
func (j *OrderSyncJob) StartAsync(ctx context.Context) error {
	fmt.Println("订单同步任务启动...")

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				j.syncOrders()
			case <-j.stopChan:
				fmt.Println("订单同步任务停止")
				return
			case <-ctx.Done():
				fmt.Println("订单同步任务被取消")
				return
			}
		}
	}()

	return nil
}

// StopAsync 停止任务
func (j *OrderSyncJob) StopAsync(ctx context.Context) error {
	fmt.Println("正在停止订单同步任务...")
	close(j.stopChan)
	return nil
}

// syncOrders 同步订单
func (j *OrderSyncJob) syncOrders() {
	fmt.Printf("[%s] 执行订单同步...\n", time.Now().Format("15:04:05"))
	
	// 实际项目中这里会执行真实的同步逻辑
	// 比如：同步到第三方系统、更新订单状态等
	
	// 模拟处理时间
	time.Sleep(2 * time.Second)
	
	fmt.Printf("[%s] 订单同步完成\n", time.Now().Format("15:04:05"))
}

