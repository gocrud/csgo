package auth

import (
	"github.com/gocrud/csgo/web"
)

// AuthController 认证控制器
type AuthController struct {
	loginHandler    *LoginHandler
	registerHandler *RegisterHandler
}

// NewAuthController 创建认证控制器
func NewAuthController(
	loginHandler *LoginHandler,
	registerHandler *RegisterHandler,
) *AuthController {
	return &AuthController{
		loginHandler:    loginHandler,
		registerHandler: registerHandler,
	}
}

// MapRoutes 映射路由
func (ctrl *AuthController) MapRoutes(app *web.WebApplication) {
	auth := app.MapGroup("/api/auth")
	auth.MapPost("/login", ctrl.loginHandler.Handle)
	auth.MapPost("/register", ctrl.registerHandler.Handle)
}

