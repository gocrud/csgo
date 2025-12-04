package utils

import (
	"fmt"
	"time"
)

// GenerateToken 生成 JWT token（简化版）
func GenerateToken(userID int64, role string) string {
	// 实际项目中应使用真正的 JWT 库
	return fmt.Sprintf("token_%d_%s_%d", userID, role, time.Now().Unix())
}

// VerifyToken 验证 token（简化版）
func VerifyToken(token string) (userID int64, role string, err error) {
	// 实际项目中应使用真正的 JWT 库
	if token == "" {
		return 0, "", fmt.Errorf("token is empty")
	}
	// 简化实现：直接解析
	fmt.Sscanf(token, "token_%d_%s", &userID, &role)
	return userID, role, nil
}

