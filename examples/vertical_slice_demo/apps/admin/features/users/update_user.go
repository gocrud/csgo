package users

import (
	"github.com/gocrud/csgo/web"
	"vertical_slice_demo/shared/contracts/repositories"
)

// UpdateUserHandler 更新用户处理器
type UpdateUserHandler struct {
	userRepo repositories.IUserRepository
}

// NewUpdateUserHandler 更新用户处理器构造函数
func NewUpdateUserHandler(userRepo repositories.IUserRepository) *UpdateUserHandler {
	return &UpdateUserHandler{userRepo: userRepo}
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Name string `json:"name" binding:"required"`
	Role string `json:"role" binding:"required,oneof=admin user"`
}

// Handle 处理更新用户请求
func (h *UpdateUserHandler) Handle(c *web.HttpContext) web.IActionResult {
	// 获取用户 ID
	id, err := c.PathInt64("id")
	if err != nil {
		return c.BadRequest("无效的用户 ID")
	}

	// 绑定请求
	var req UpdateUserRequest
	if err := c.MustBindJSON(&req); err != nil {
		return err
	}

	// 查询用户
	user, err := h.userRepo.GetByID(id)
	if err != nil {
		return c.NotFound("用户不存在")
	}

	// 更新用户信息
	user.Name = req.Name
	user.Role = req.Role

	// 持久化
	if err := h.userRepo.Update(user); err != nil {
		return c.InternalError("更新用户失败")
	}

	return c.Ok(user)
}

