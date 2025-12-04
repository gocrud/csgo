package validation

import "fmt"

// Rule 验证规则接口
type Rule interface {
	Validate(value interface{}) error
}

// PropertyRule 属性验证规则
type PropertyRule[T any, TProperty any] interface {
	Validate(instance *T, value TProperty) error
}

// ValidatorRegistry 验证器注册表
var validatorRegistry = make(map[string]interface{})

// RegisterValidator 注册验证器
func RegisterValidator[T any](validator IValidator[T]) {
	var t T
	typeName := getTypeName(t)
	validatorRegistry[typeName] = validator
}

// GetValidator 获取验证器
func GetValidator[T any]() (IValidator[T], bool) {
	var t T
	typeName := getTypeName(t)
	if v, ok := validatorRegistry[typeName]; ok {
		if validator, ok := v.(IValidator[T]); ok {
			return validator, true
		}
	}
	return nil, false
}

// getTypeName 获取类型名称
func getTypeName(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

// ValidateStruct 验证结构体（使用注册的验证器）
func ValidateStruct[T any](instance *T) ValidationResult {
	if validator, ok := GetValidator[T](); ok {
		return validator.Validate(instance)
	}
	// 如果没有注册验证器，返回成功
	return SuccessResult()
}
