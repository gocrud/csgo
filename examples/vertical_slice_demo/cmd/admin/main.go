package main

import (
	"fmt"

	"vertical_slice_demo/apps/admin"
	"vertical_slice_demo/configs"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("        垂直切片架构示例 - 管理端")
	fmt.Println("========================================")
	fmt.Println()

	// 启动管理端应用
	app := admin.Bootstrap()

	// 获取配置
	config := configs.DefaultConfig()

	fmt.Printf("管理端 API 启动在: http://localhost%s\n", config.Server.AdminPort)
	fmt.Println()
	fmt.Println("可用的 API 端点:")
	fmt.Println("  POST   /api/admin/users       - 创建用户")
	fmt.Println("  GET    /api/admin/users       - 用户列表")
	fmt.Println("  PUT    /api/admin/users/:id   - 更新用户")
	fmt.Println("  POST   /api/admin/products    - 创建商品")
	fmt.Println("  GET    /api/admin/products    - 商品列表")
	fmt.Println()
	fmt.Println("注意：所有请求需要在 Header 中添加 Authorization")
	fmt.Println()

	// 运行应用
	app.Run(config.Server.AdminPort)
}

