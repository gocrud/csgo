package errors

import (
	"fmt"
)

// ErrorCategory 错误分类
type ErrorCategory string

const (
	CategoryBusiness   ErrorCategory = "BUSINESS"   // 业务错误
	CategorySystem     ErrorCategory = "SYSTEM"     // 系统错误
	CategoryValidation ErrorCategory = "VALIDATION" // 验证错误
	CategoryAuth       ErrorCategory = "AUTH"       // 认证授权错误
)

// Error 表示一个结构化的错误
type Error struct {
	category ErrorCategory  // 错误分类
	code     string         // 完整错误码，如 "USER.NOT_FOUND"
	message  string         // 错误消息
	cause    error          // 原始错误（支持错误链）
	details  map[string]any // 额外详细信息
	httpCode int            // HTTP 状态码
}

// Error 实现 error 接口
func (e *Error) Error() string {
	if e.code != "" {
		msg := fmt.Sprintf("[%s] %s", e.code, e.message)
		if e.cause != nil {
			msg += fmt.Sprintf(": %v", e.cause)
		}
		return msg
	}
	return e.message
}

// Unwrap 实现 Go 1.13+ 错误链支持
func (e *Error) Unwrap() error {
	return e.cause
}

// Code 返回错误码
func (e *Error) Code() string {
	return e.code
}

// Message 返回错误消息
func (e *Error) Message() string {
	return e.message
}

// HTTPCode 返回 HTTP 状态码
func (e *Error) HTTPCode() int {
	return e.httpCode
}

// Details 返回详细信息
func (e *Error) Details() map[string]any {
	if e.details == nil {
		return make(map[string]any)
	}
	// 返回副本，避免外部修改
	result := make(map[string]any, len(e.details))
	for k, v := range e.details {
		result[k] = v
	}
	return result
}

// Category 返回错误分类
func (e *Error) Category() ErrorCategory {
	return e.category
}

// Wrap 包装原始错误（返回新实例）
func (e *Error) Wrap(err error) *Error {
	newErr := e.clone()
	newErr.cause = err
	return newErr
}

// WithMsg 覆盖错误消息（返回新实例）
func (e *Error) WithMsg(message string) *Error {
	newErr := e.clone()
	newErr.message = message
	return newErr
}

// WithMsgf 使用格式化字符串覆盖错误消息（返回新实例）
func (e *Error) WithMsgf(format string, args ...any) *Error {
	return e.WithMsg(fmt.Sprintf(format, args...))
}

// AppendMsg 追加消息（返回新实例）
func (e *Error) AppendMsg(suffix string) *Error {
	newErr := e.clone()
	newErr.message = e.message + suffix
	return newErr
}

// PrependMsg 前置消息（返回新实例）
func (e *Error) PrependMsg(prefix string) *Error {
	newErr := e.clone()
	newErr.message = prefix + e.message
	return newErr
}

// WithDetail 添加单个详细信息（返回新实例）
func (e *Error) WithDetail(key string, value any) *Error {
	newErr := e.clone()
	if newErr.details == nil {
		newErr.details = make(map[string]any)
	}
	newErr.details[key] = value
	return newErr
}

// WithDetails 批量添加详细信息（返回新实例）
func (e *Error) WithDetails(details map[string]any) *Error {
	newErr := e.clone()
	if newErr.details == nil {
		newErr.details = make(map[string]any)
	}
	for k, v := range details {
		newErr.details[k] = v
	}
	return newErr
}

// WithHTTPCode 设置 HTTP 状态码（返回新实例）
func (e *Error) WithHTTPCode(code int) *Error {
	newErr := e.clone()
	newErr.httpCode = code
	return newErr
}

// clone 克隆错误对象
func (e *Error) clone() *Error {
	newErr := &Error{
		category: e.category,
		code:     e.code,
		message:  e.message,
		cause:    e.cause,
		httpCode: e.httpCode,
	}
	if e.details != nil {
		newErr.details = make(map[string]any, len(e.details))
		for k, v := range e.details {
			newErr.details[k] = v
		}
	}
	return newErr
}

// New 创建自定义错误
func New(code string, message string, httpCode int) *Error {
	return &Error{
		category: CategoryBusiness,
		code:     code,
		message:  message,
		httpCode: httpCode,
		details:  make(map[string]any),
	}
}

// Newf 使用格式化消息创建自定义错误
func Newf(code string, httpCode int, format string, args ...any) *Error {
	return &Error{
		category: CategoryBusiness,
		code:     code,
		message:  fmt.Sprintf(format, args...),
		httpCode: httpCode,
		details:  make(map[string]any),
	}
}
