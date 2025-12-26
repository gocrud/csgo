package errors

import (
	"fmt"
	"strings"
)

// Module 表示一个错误模块，提供模块化的错误创建
type Module struct {
	prefix string
}

// NewModule 创建一个新的错误模块
// prefix 将自动转换为大写，如 "user" -> "USER"
func NewModule(prefix string) *Module {
	return &Module{
		prefix: strings.ToUpper(prefix),
	}
}

// ============ 常用错误快捷方法 ============

// NotFound 资源不存在（404）
func (m *Module) NotFound(msg ...string) *Error {
	return m.newError("NOT_FOUND", 404, CategoryBusiness, "资源不存在", msg...)
}

// AlreadyExists 资源已存在（409）
func (m *Module) AlreadyExists(msg ...string) *Error {
	return m.newError("ALREADY_EXISTS", 409, CategoryBusiness, "资源已存在", msg...)
}

// InvalidParam 参数无效（400）
func (m *Module) InvalidParam(msg ...string) *Error {
	return m.newError("INVALID_PARAM", 400, CategoryBusiness, "参数无效", msg...)
}

// InvalidStatus 状态无效（400）
func (m *Module) InvalidStatus(msg ...string) *Error {
	return m.newError("INVALID_STATUS", 400, CategoryBusiness, "状态无效", msg...)
}

// PermissionDenied 权限不足（403）
func (m *Module) PermissionDenied(msg ...string) *Error {
	return m.newError("PERMISSION_DENIED", 403, CategoryAuth, "权限不足", msg...)
}

// Unauthorized 未授权（401）
func (m *Module) Unauthorized(msg ...string) *Error {
	return m.newError("UNAUTHORIZED", 401, CategoryAuth, "未授权", msg...)
}

// OperationFailed 操作失败（400）
func (m *Module) OperationFailed(msg ...string) *Error {
	return m.newError("OPERATION_FAILED", 400, CategoryBusiness, "操作失败", msg...)
}

// Expired 资源已过期（410）
func (m *Module) Expired(msg ...string) *Error {
	return m.newError("EXPIRED", 410, CategoryBusiness, "资源已过期", msg...)
}

// Locked 资源已锁定（423）
func (m *Module) Locked(msg ...string) *Error {
	return m.newError("LOCKED", 423, CategoryBusiness, "资源已锁定", msg...)
}

// LimitExceeded 超出限制（429）
func (m *Module) LimitExceeded(msg ...string) *Error {
	return m.newError("LIMIT_EXCEEDED", 429, CategoryBusiness, "超出限制", msg...)
}

// Conflict 资源冲突（409）
func (m *Module) Conflict(msg ...string) *Error {
	return m.newError("CONFLICT", 409, CategoryBusiness, "资源冲突", msg...)
}

// Internal 内部错误（500）
func (m *Module) Internal(msg ...string) *Error {
	return m.newError("INTERNAL_ERROR", 500, CategorySystem, "内部错误", msg...)
}

// ServiceUnavailable 服务不可用（503）
func (m *Module) ServiceUnavailable(msg ...string) *Error {
	return m.newError("SERVICE_UNAVAILABLE", 503, CategorySystem, "服务不可用", msg...)
}

// ============ 自定义错误码构建器 ============

// Code 创建自定义错误码构建器
// 用法: module.Code("PAYMENT_FAILED").Msg("支付失败")
func (m *Module) Code(code string) *ErrorBuilder {
	return &ErrorBuilder{
		module: m,
		code:   strings.ToUpper(code),
	}
}

// ErrorBuilder 错误码构建器
type ErrorBuilder struct {
	module *Module
	code   string
}

// Msg 设置错误消息并返回 Error（默认 400 状态码）
func (b *ErrorBuilder) Msg(msg string) *Error {
	return b.module.newError(b.code, 400, CategoryBusiness, "", msg)
}

// Msgf 使用格式化字符串设置错误消息
func (b *ErrorBuilder) Msgf(format string, args ...any) *Error {
	return b.Msg(fmt.Sprintf(format, args...))
}

// MsgWithCode 设置错误消息和 HTTP 状态码
func (b *ErrorBuilder) MsgWithCode(msg string, httpCode int) *Error {
	return b.module.newError(b.code, httpCode, CategoryBusiness, "", msg)
}

// ============ 完全自定义 ============

// Custom 创建完全自定义的错误
func (m *Module) Custom(code string, msg string, httpCode int) *Error {
	return m.newError(code, httpCode, CategoryBusiness, "", msg)
}

// Customf 使用格式化消息创建完全自定义的错误
func (m *Module) Customf(code string, httpCode int, format string, args ...any) *Error {
	return m.Custom(code, fmt.Sprintf(format, args...), httpCode)
}

// ============ 内部辅助方法 ============

// newError 创建新的错误实例
func (m *Module) newError(code string, httpCode int, category ErrorCategory, defaultMsg string, msg ...string) *Error {
	message := defaultMsg
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}
	if message == "" {
		message = code // 兜底使用错误码
	}

	return &Error{
		category: category,
		code:     m.prefix + "." + code,
		message:  message,
		httpCode: httpCode,
		details:  make(map[string]any),
	}
}

// Prefix 返回模块前缀
func (m *Module) Prefix() string {
	return m.prefix
}
