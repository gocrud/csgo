package users

import (
	"github.com/gocrud/csgo/web"
	"vertical_slice_demo/shared/contracts/repositories"
	"vertical_slice_demo/shared/domain"
	"vertical_slice_demo/shared/utils"
)

// CreateUserHandler 创建用户处理器
type CreateUserHandler struct {
	userRepo repositories.IUserRepository
}

// NewCreateUserHandler 创建用户处理器构造函数
func NewCreateUserHandler(userRepo repositories.IUserRepository) *CreateUserHandler {
	return &CreateUserHandler{userRepo: userRepo}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// CreateUserResponse 创建用户响应
type CreateUserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// Handle 处理创建用户请求（使用 FluentValidation）
func (h *CreateUserHandler) Handle(c *web.HttpContext) web.IActionResult {
	// ✅ 使用 BindAndValidate 自动绑定和验证
	req, err := web.BindAndValidate[CreateUserRequest](c)
	if err != nil {
		return err
	}

	// 业务验证：检查邮箱是否已存在
	if h.userRepo.ExistsByEmail(req.Email) {
		return c.BadRequest("邮箱已存在")
	}

	// 创建用户实体
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: utils.HashPassword(req.Password),
		Role:     req.Role,
	}

	// 持久化
	if err := h.userRepo.Create(user); err != nil {
		return c.InternalError("创建用户失败")
	}

	// 返回响应
	response := &CreateUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	return c.Created(response)
}

