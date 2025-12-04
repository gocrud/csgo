package main

import (
	"fmt"

	"vertical_slice_demo/apps/worker"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("      垂直切片架构示例 - Worker")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("后台任务:")
	fmt.Println("  - 订单同步任务（每 30 秒执行一次）")
	fmt.Println("  - 邮件发送任务（每 60 秒执行一次）")
	fmt.Println()
	fmt.Println("按 Ctrl+C 停止...")
	fmt.Println()

	// 启动 Worker 服务
	host := worker.Bootstrap()
	host.Run()
}

