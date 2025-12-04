package errors

import "fmt"

// AppError 应用错误
type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// 常见错误
var (
	ErrNotFound       = &AppError{Code: "NOT_FOUND", Message: "资源不存在"}
	ErrAlreadyExists  = &AppError{Code: "ALREADY_EXISTS", Message: "资源已存在"}
	ErrInvalidInput   = &AppError{Code: "INVALID_INPUT", Message: "输入无效"}
	ErrUnauthorized   = &AppError{Code: "UNAUTHORIZED", Message: "未授权"}
	ErrForbidden      = &AppError{Code: "FORBIDDEN", Message: "禁止访问"}
	ErrInternalServer = &AppError{Code: "INTERNAL_SERVER", Message: "内部服务器错误"}
)

// NewAppError 创建应用错误
func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

