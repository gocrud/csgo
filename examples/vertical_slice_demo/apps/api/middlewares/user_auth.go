package middlewares

import (
	"github.com/gin-gonic/gin"
)

// UserAuthMiddleware C端认证中间件
func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 对于公开接口（登录、注册、商品浏览）不需要认证
		if isPublicPath(c.FullPath()) {
			c.Next()
			return
		}

		// 检查 Authorization header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "未授权：请先登录"})
			c.Abort()
			return
		}

		// 实际项目中应该验证 JWT token
		// userID, role, err := utils.VerifyToken(token)
		// if err != nil {
		//     c.JSON(401, gin.H{"error": "认证失败"})
		//     c.Abort()
		//     return
		// }

		// 将用户信息存储到上下文
		// c.Set("user_id", userID)
		// c.Set("role", role)

		c.Next()
	}
}

// isPublicPath 判断是否是公开路径
func isPublicPath(path string) bool {
	publicPaths := []string{
		"/api/auth/login",
		"/api/auth/register",
		"/api/products",
		"/",
	}

	for _, p := range publicPaths {
		if path == p || path == p+"/:id" {
			return true
		}
	}

	return false
}

