package users

import (
	"github.com/gocrud/csgo/web"
)

// UserController 用户控制器
type UserController struct {
	createHandler *CreateUserHandler
	listHandler   *ListUsersHandler
	updateHandler *UpdateUserHandler
}

// NewUserController 创建用户控制器
func NewUserController(
	createHandler *CreateUserHandler,
	listHandler *ListUsersHandler,
	updateHandler *UpdateUserHandler,
) *UserController {
	return &UserController{
		createHandler: createHandler,
		listHandler:   listHandler,
		updateHandler: updateHandler,
	}
}

// MapRoutes 映射路由
func (ctrl *UserController) MapRoutes(app *web.WebApplication) {
	users := app.MapGroup("/api/admin/users")
	users.MapPost("", ctrl.createHandler.Handle)
	users.MapGet("", ctrl.listHandler.Handle)
	users.MapPut("/:id", ctrl.updateHandler.Handle)
}

