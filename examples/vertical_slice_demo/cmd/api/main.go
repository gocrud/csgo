package main

import (
	"fmt"

	"vertical_slice_demo/apps/api"
	"vertical_slice_demo/configs"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("        垂直切片架构示例 - C端")
	fmt.Println("========================================")
	fmt.Println()

	// 启动C端应用
	app := api.Bootstrap()

	// 获取配置
	config := configs.DefaultConfig()

	fmt.Printf("C端 API 启动在: http://localhost%s\n", config.Server.ApiPort)
	fmt.Println()
	fmt.Println("可用的 API 端点:")
	fmt.Println("  POST   /api/auth/register    - 用户注册")
	fmt.Println("  POST   /api/auth/login       - 用户登录")
	fmt.Println("  GET    /api/products         - 浏览商品")
	fmt.Println("  GET    /api/products/:id     - 商品详情")
	fmt.Println("  POST   /api/orders           - 创建订单")
	fmt.Println("  GET    /api/orders/my        - 我的订单")
	fmt.Println()
	fmt.Println("注意：除了注册、登录、商品浏览外，其他接口需要认证")
	fmt.Println()

	// 运行应用
	app.Run(config.Server.ApiPort)
}

