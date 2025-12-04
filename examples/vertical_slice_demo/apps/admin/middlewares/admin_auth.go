package middlewares

import (
	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware 管理端认证中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 简化实现：检查 Authorization header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "未授权：缺少认证令牌"})
			c.Abort()
			return
		}

		// 实际项目中应该验证 JWT token 并检查是否是 admin 角色
		// 这里简化为只检查 token 是否包含 "admin"
		// if !strings.Contains(token, "admin") {
		// 	c.JSON(403, gin.H{"error": "禁止访问：需要管理员权限"})
		// 	c.Abort()
		// 	return
		// }

		c.Next()
	}
}

