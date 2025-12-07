package main

import "github.com/gocrud/csgo/web"

// MigrationExample 展示如何从旧 API 迁移到新 API
//
// 这个文件仅用于演示目的，展示 API 变更前后的对比
// 不要编译运行此文件

// ==================== 旧 API (已移除) ====================

// func oldStyleHandler(c *web.HttpContext) web.IActionResult {
// 	// ❌ 以下方法已被移除
// 	// id, err := c.PathInt("id")
// 	// if err != nil {
// 	//     return c.BadRequest("Invalid ID")
// 	// }
//
// 	// page := c.QueryInt("page", 1)
// 	// size := c.QueryInt("size", 10)
// 	// active := c.QueryBool("active", true)
//
// 	// name := c.Query("name")
// 	// auth := c.GetHeader("Authorization")
//
// 	// return c.Ok(data)
// }

// ==================== 新 API (推荐) ====================

// Example 1: 路径参数验证
func getUserHandler(c *web.HttpContext) web.IActionResult {
	// 使用 Params() 验证器获取并验证路径参数
	p := c.Params()
	id := p.PathInt("id").Positive().Value()

	// 检查验证错误
	if err := p.Check(); err != nil {
		return err // 自动返回 400 + 详细验证错误
	}

	// 使用验证后的参数
	user := getUserByID(id)
	return c.Ok(user)
}

// Example 2: 查询参数验证
func listUsersHandler(c *web.HttpContext) web.IActionResult {
	// 创建参数验证器
	p := c.Params()

	// 链式调用进行验证，使用 ValueOr 提供默认值
	page := p.QueryInt("page").Range(1, 100).ValueOr(1)
	size := p.QueryInt("size").Range(1, 50).ValueOr(10)
	status := p.QueryString("status").In("active", "inactive").ValueOr("active")

	// 批量检查所有验证错误（可选）
	if err := p.Check(); err != nil {
		return err
	}

	users := listUsers(page, size, status)
	return c.Ok(users)
}

// Example 3: 可选查询参数
func searchProductsHandler(c *web.HttpContext) web.IActionResult {
	p := c.Params()

	// 可选参数使用 Optional() 标记
	offset := p.QueryInt("offset").Optional().NonNegative().ValueOr(0)
	limit := p.QueryInt("limit").Optional().Range(1, 100).ValueOr(20)

	// 可选字符串参数
	keyword := p.QueryString("keyword").Optional().MinLength(2).ValueOr("")

	products := searchProducts(keyword, offset, limit)
	return c.Ok(products)
}

// Example 4: Header 参数验证
func authenticatedHandler(c *web.HttpContext) web.IActionResult {
	p := c.Params()

	// 验证必需的 header
	token := p.HeaderString("Authorization").Required().MinLength(10).Value()

	if err := p.Check(); err != nil {
		return c.Unauthorized("Missing or invalid authorization")
	}

	// 使用 token
	user := validateToken(token)
	return c.Ok(user)
}

// Example 5: 访问底层 gin.Context (高级用法)
func advancedHandler(c *web.HttpContext) web.IActionResult {
	// 使用 RawCtx() 访问底层 gin.Context
	clientIP := c.RawCtx().ClientIP()
	userAgent := c.RawCtx().GetHeader("User-Agent")

	// 设置上下文值
	c.RawCtx().Set("user_id", 123)

	// 获取上下文值
	value, exists := c.RawCtx().Get("user_id")

	return c.Ok(map[string]interface{}{
		"ip":         clientIP,
		"user_agent": userAgent,
		"value":      value,
		"exists":     exists,
	})
}

// Example 6: 复杂验证场景
func updateUserHandler(c *web.HttpContext) web.IActionResult {
	// 1. 验证路径参数
	p := c.Params()
	userID := p.PathInt64("id").Positive().Value()

	if err := p.Check(); err != nil {
		return err
	}

	// 2. 绑定并验证 JSON body
	type UpdateRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	req, err := web.BindAndValidate[UpdateRequest](c)
	if err != nil {
		return err
	}

	// 3. 执行更新
	user := updateUser(userID, req.Name, req.Email)
	return c.Ok(user)
}

// Example 7: 多个参数验证
func filterOrdersHandler(c *web.HttpContext) web.IActionResult {
	p := c.Params()

	// 同时验证多个参数
	userID := p.PathInt64("userId").Positive().Value()
	status := p.QueryString("status").In("pending", "paid", "shipped", "completed").ValueOr("pending")
	fromDate := p.QueryString("from").Optional().Pattern(`^\d{4}-\d{2}-\d{2}$`, "日期格式必须为 YYYY-MM-DD").Value()
	toDate := p.QueryString("to").Optional().Pattern(`^\d{4}-\d{2}-\d{2}$`, "日期格式必须为 YYYY-MM-DD").Value()
	offset := p.QueryInt("offset").NonNegative().ValueOr(0)
	limit := p.QueryInt("limit").Range(1, 100).ValueOr(20)

	// 一次性检查所有错误
	if err := p.Check(); err != nil {
		return err // 返回所有验证错误的详细信息
	}

	orders := filterOrders(userID, status, fromDate, toDate, offset, limit)
	return c.Ok(orders)
}

// ==================== 辅助函数 (仅用于示例) ====================

func getUserByID(id int) interface{}                               { return nil }
func listUsers(page, size int, status string) interface{}          { return nil }
func searchProducts(keyword string, offset, limit int) interface{} { return nil }
func validateToken(token string) interface{}                       { return nil }
func updateUser(id int64, name, email string) interface{}          { return nil }
func filterOrders(userID int64, status, from, to string, offset, limit int) interface{} {
	return nil
}
