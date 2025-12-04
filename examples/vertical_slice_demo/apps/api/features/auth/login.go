package auth

import (
	"github.com/gocrud/csgo/web"
	"vertical_slice_demo/shared/contracts/repositories"
	"vertical_slice_demo/shared/utils"
)

// LoginHandler 登录处理器
type LoginHandler struct {
	userRepo repositories.IUserRepository
}

// NewLoginHandler 登录处理器构造函数
func NewLoginHandler(userRepo repositories.IUserRepository) *LoginHandler {
	return &LoginHandler{userRepo: userRepo}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token  string `json:"token"`
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

// Handle 处理登录请求
func (h *LoginHandler) Handle(c *web.HttpContext) web.IActionResult {
	var req LoginRequest
	if err := c.MustBindJSON(&req); err != nil {
		return err
	}

	// 查询用户
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		return c.Unauthorized("邮箱或密码错误")
	}

	// 验证密码
	if !utils.VerifyPassword(user.Password, req.Password) {
		return c.Unauthorized("邮箱或密码错误")
	}

	// 生成 token
	token := utils.GenerateToken(user.ID, user.Role)

	response := &LoginResponse{
		Token:  token,
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}

	return c.Ok(response)
}

