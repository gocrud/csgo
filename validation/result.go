package validation

import (
	"fmt"
	"strings"
)

// ValidationError 表示单个字段的验证错误
type ValidationError struct {
	Field   string `json:"field"`   // 字段路径，优先使用 JSON tag，如 "email" 或嵌套 "address.city"
	Message string `json:"message"` // 错误消息
	Code    string `json:"code"`    // 错误码，如 "VALIDATION.REQUIRED"
}

// Error 实现 error 接口
func (e ValidationError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Field, e.Message)
	}
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

// Add 添加验证错误（需要提供错误码）
func (e *ValidationErrors) Add(field, message, code string) {
	*e = append(*e, ValidationError{
		Field:   field,
		Message: message,
		Code:    code,
	})
}

// GetByField 获取指定字段的所有错误
func (e ValidationErrors) GetByField(field string) []ValidationError {
	var result []ValidationError
	for _, err := range e {
		if err.Field == field {
			result = append(result, err)
		}
	}
	return result
}

// Fields 返回所有错误字段名列表（去重）
func (e ValidationErrors) Fields() []string {
	fieldMap := make(map[string]bool)
	var fields []string
	for _, err := range e {
		if !fieldMap[err.Field] {
			fieldMap[err.Field] = true
			fields = append(fields, err.Field)
		}
	}
	return fields
}

// FirstError 返回第一个错误（如果存在）
func (e ValidationErrors) FirstError() *ValidationError {
	if len(e) > 0 {
		return &e[0]
	}
	return nil
}

// ToFieldMap 转换为 map[字段名][]错误消息
func (e ValidationErrors) ToFieldMap() map[string][]string {
	result := make(map[string][]string)
	for _, err := range e {
		result[err.Field] = append(result[err.Field], err.Message)
	}
	return result
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
