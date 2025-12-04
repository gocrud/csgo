package validation

import (
	"fmt"
	"strings"
)

// ValidationError 表示单个字段的验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// Error 实现 error 接口
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors 验证错误集合
type ValidationErrors []ValidationError

// Error 实现 error 接口
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// Add 添加验证错误
func (e *ValidationErrors) Add(field, message string) {
	*e = append(*e, ValidationError{
		Field:   field,
		Message: message,
	})
}

// AddWithCode 添加带错误码的验证错误
func (e *ValidationErrors) AddWithCode(field, message, code string) {
	*e = append(*e, ValidationError{
		Field:   field,
		Message: message,
		Code:    code,
	})
}

// HasErrors 检查是否有错误
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// ValidationResult 验证结果
type ValidationResult struct {
	IsValid bool
	Errors  ValidationErrors
}

// ToError 转换为 error（如果有错误）
func (r ValidationResult) ToError() error {
	if r.IsValid {
		return nil
	}
	return r.Errors
}

// NewValidationResult 创建验证结果
func NewValidationResult(errors ValidationErrors) ValidationResult {
	return ValidationResult{
		IsValid: !errors.HasErrors(),
		Errors:  errors,
	}
}

// SuccessResult 创建成功的验证结果
func SuccessResult() ValidationResult {
	return ValidationResult{
		IsValid: true,
		Errors:  ValidationErrors{},
	}
}
