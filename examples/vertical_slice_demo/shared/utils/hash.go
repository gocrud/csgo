package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword 哈希密码（简化版，生产环境应使用 bcrypt）
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// VerifyPassword 验证密码
func VerifyPassword(hashedPassword, password string) bool {
	return hashedPassword == HashPassword(password)
}

