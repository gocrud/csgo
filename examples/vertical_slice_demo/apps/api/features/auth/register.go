package auth

import (
	"github.com/gocrud/csgo/web"
	"vertical_slice_demo/shared/contracts/repositories"
	"vertical_slice_demo/shared/domain"
	"vertical_slice_demo/shared/utils"
)

// RegisterHandler 注册处理器
type RegisterHandler struct {
	userRepo repositories.IUserRepository
}

// NewRegisterHandler 注册处理器构造函数
func NewRegisterHandler(userRepo repositories.IUserRepository) *RegisterHandler {
	return &RegisterHandler{userRepo: userRepo}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	Token  string `json:"token"`
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

// Handle 处理注册请求
func (h *RegisterHandler) Handle(c *web.HttpContext) web.IActionResult {
	var req RegisterRequest
	if err := c.MustBindJSON(&req); err != nil {
		return err
	}

	// 业务验证：检查邮箱是否已存在
	if h.userRepo.ExistsByEmail(req.Email) {
		return c.BadRequest("邮箱已被注册")
	}

	// 创建用户实体（C端用户默认角色为 user）
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: utils.HashPassword(req.Password),
		Role:     "user",
	}

	// 持久化
	if err := h.userRepo.Create(user); err != nil {
		return c.InternalError("注册失败")
	}

	// 自动登录：生成 token
	token := utils.GenerateToken(user.ID, user.Role)

	response := &RegisterResponse{
		Token:  token,
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}

	return c.Created(response)
}

