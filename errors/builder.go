package errors

import (
	"fmt"
	"strings"
)

// BizError 业务错误
type BizError struct {
	Code    string                 // 错误码，如 "USER.NOT_FOUND"
	Message string                 // 错误消息
	Cause   error                  // 原始错误（支持错误链）
	Details map[string]interface{} // 额外详细信息
}

// Error 实现 error 接口
func (e *BizError) Error() string {
	if e.Code != "" {
		msg := fmt.Sprintf("[%s] %s", e.Code, e.Message)
		if e.Cause != nil {
			msg += fmt.Sprintf(": %v", e.Cause)
		}
		return msg
	}
	return e.Message
}

// Unwrap 实现 Go 1.13+ 错误链支持
func (e *BizError) Unwrap() error {
	return e.Cause
}

// Wrap 包装原始错误
func (e *BizError) Wrap(err error) *BizError {
	e.Cause = err
	return e
}

// WithMsg 设置或修改错误消息
func (e *BizError) WithMsg(message string) *BizError {
	e.Message = message
	return e
}

// WithDetail 添加单个详细信息
func (e *BizError) WithDetail(key string, value interface{}) *BizError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithDetails 批量添加详细信息
func (e *BizError) WithDetails(details map[string]interface{}) *BizError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	for k, v := range details {
		e.Details[k] = v
	}
	return e
}

// ErrorBuilder 错误码构建器
type ErrorBuilder struct {
	module string
}

// Business 创建业务错误构建器
// 用法: errors.Business("USER").NotFound("用户不存在")
func Business(module string) *ErrorBuilder {
	return &ErrorBuilder{module: strings.ToUpper(module)}
}

// NotFound 资源不存在
func (b *ErrorBuilder) NotFound(message string) *BizError {
	return &BizError{
		Code:    fmt.Sprintf("%s.NOT_FOUND", b.module),
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// AlreadyExists 资源已存在
func (b *ErrorBuilder) AlreadyExists(message string) *BizError {
	return &BizError{
		Code:    fmt.Sprintf("%s.ALREADY_EXISTS", b.module),
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// InvalidStatus 状态无效
func (b *ErrorBuilder) InvalidStatus(message string) *BizError {
	return &BizError{
		Code:    fmt.Sprintf("%s.INVALID_STATUS", b.module),
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// InvalidParam 参数无效
func (b *ErrorBuilder) InvalidParam(message string) *BizError {
	return &BizError{
		Code:    fmt.Sprintf("%s.INVALID_PARAM", b.module),
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// PermissionDenied 权限不足
func (b *ErrorBuilder) PermissionDenied(message string) *BizError {
	return &BizError{
		Code:    fmt.Sprintf("%s.PERMISSION_DENIED", b.module),
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// OperationFailed 操作失败
func (b *ErrorBuilder) OperationFailed(message string) *BizError {
	return &BizError{
		Code:    fmt.Sprintf("%s.OPERATION_FAILED", b.module),
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// Expired 资源已过期
func (b *ErrorBuilder) Expired(message string) *BizError {
	return &BizError{
		Code:    fmt.Sprintf("%s.EXPIRED", b.module),
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// Locked 资源已锁定
func (b *ErrorBuilder) Locked(message string) *BizError {
	return &BizError{
		Code:    fmt.Sprintf("%s.LOCKED", b.module),
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// LimitExceeded 超出限制
func (b *ErrorBuilder) LimitExceeded(message string) *BizError {
	return &BizError{
		Code:    fmt.Sprintf("%s.LIMIT_EXCEEDED", b.module),
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// Custom 自定义语义错误码
// semantic 参数应该使用大写下划线命名，如 "AMOUNT_EXCEEDED"
func (b *ErrorBuilder) Custom(semantic, message string) *BizError {
	semantic = strings.ToUpper(semantic)
	return &BizError{
		Code:    fmt.Sprintf("%s.%s", b.module, semantic),
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// New 直接创建业务错误（不使用构建器）
// 返回 *BizError，支持链式调用 Wrap、WithDetail 等方法
func New(code, message string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// Newf 使用格式化消息创建业务错误
// 返回 *BizError，支持链式调用 Wrap、WithDetail 等方法
func Newf(code, format string, args ...interface{}) *BizError {
	return &BizError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
		Details: make(map[string]interface{}),
	}
}
