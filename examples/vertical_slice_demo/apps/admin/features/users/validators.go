package users

import (
	"strings"

	"github.com/gocrud/csgo/validation"
)

// CreateUserRequestValidator 创建用户请求验证器
type CreateUserRequestValidator struct {
	*validation.AbstractValidator[CreateUserRequest]
}

// NewCreateUserRequestValidator 创建验证器实例
func NewCreateUserRequestValidator() *CreateUserRequestValidator {
	v := &CreateUserRequestValidator{
		AbstractValidator: validation.NewValidator[CreateUserRequest](),
	}

	// Name 验证 - 使用新的 Field API，自动从 json tag 提取字段名
	v.Field(func(r *CreateUserRequest) string { return r.Name }).
		NotEmpty().WithMessage("用户名不能为空").
		Length(2, 50).WithMessage("用户名长度必须在2-50之间")

	// Email 验证
	v.Field(func(r *CreateUserRequest) string { return r.Email }).
		NotEmpty().WithMessage("邮箱不能为空").
		EmailAddress().WithMessage("邮箱格式不正确")

	// Password 验证
	v.Field(func(r *CreateUserRequest) string { return r.Password }).
		NotEmpty().WithMessage("密码不能为空").
		MinLength(8).WithMessage("密码长度至少8位").
		Matches(`[A-Z]`).WithMessage("密码必须包含大写字母").
		Matches(`[0-9]`).WithMessage("密码必须包含数字").
		MustString(func(req *CreateUserRequest, password string) bool {
			return !strings.Contains(password, req.Name)
		}).WithMessage("密码不能包含用户名")

	// Role 验证
	v.Field(func(r *CreateUserRequest) string { return r.Role }).
		NotEmpty().WithMessage("角色不能为空").
		MustString(func(req *CreateUserRequest, role string) bool {
			return role == "admin" || role == "user"
		}).WithMessage("角色必须是 admin 或 user")

	return v
}

// UpdateUserRequestValidator 更新用户请求验证器
type UpdateUserRequestValidator struct {
	*validation.AbstractValidator[UpdateUserRequest]
}

// NewUpdateUserRequestValidator 创建更新用户验证器
func NewUpdateUserRequestValidator() *UpdateUserRequestValidator {
	v := &UpdateUserRequestValidator{
		AbstractValidator: validation.NewValidator[UpdateUserRequest](),
	}

	// Name 验证
	v.Field(func(r *UpdateUserRequest) string { return r.Name }).
		NotEmpty().WithMessage("用户名不能为空").
		Length(2, 50).WithMessage("用户名长度必须在2-50之间")

	// Role 验证
	v.Field(func(r *UpdateUserRequest) string { return r.Role }).
		NotEmpty().WithMessage("角色不能为空").
		MustString(func(req *UpdateUserRequest, role string) bool {
			return role == "admin" || role == "user"
		}).WithMessage("角色必须是 admin 或 user")

	return v
}

// 注册验证器
func init() {
	validation.RegisterValidator[CreateUserRequest](NewCreateUserRequestValidator())
	validation.RegisterValidator[UpdateUserRequest](NewUpdateUserRequestValidator())
}
