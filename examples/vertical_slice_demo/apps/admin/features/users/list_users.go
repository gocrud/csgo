package users

import (
	"vertical_slice_demo/shared/contracts/repositories"

	"github.com/gocrud/csgo/web"
)

// ListUsersHandler 用户列表处理器
type ListUsersHandler struct {
	userRepo repositories.IUserRepository
}

// NewListUsersHandler 用户列表处理器构造函数
func NewListUsersHandler(userRepo repositories.IUserRepository) *ListUsersHandler {
	return &ListUsersHandler{userRepo: userRepo}
}

// UserListItem 用户列表项
type UserListItem struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// ListUsersResponse 用户列表响应
type ListUsersResponse struct {
	Users []UserListItem `json:"users"`
	Total int            `json:"total"`
}

// Handle 处理用户列表请求
func (h *ListUsersHandler) Handle(c *web.HttpContext) web.IActionResult {
	// 获取分页参数
	p := c.Params()
	offset := p.QueryInt("offset").NonNegative().ValueOr(0)
	limit := p.QueryInt("limit").Range(1, 100).ValueOr(20)

	// 查询用户列表
	users, err := h.userRepo.List(offset, limit)
	if err != nil {
		return c.InternalError("查询用户列表失败")
	}

	// 转换为响应格式
	items := make([]UserListItem, len(users))
	for i, user := range users {
		items[i] = UserListItem{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		}
	}

	response := &ListUsersResponse{
		Users: items,
		Total: len(items),
	}

	return c.Ok(response)
}
