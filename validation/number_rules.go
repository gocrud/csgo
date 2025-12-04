package validation

import "fmt"

// Number 表示可比较的数字类型约束
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// GreaterThan 大于
func GreaterThan[T any, TProperty Number](b *RuleBuilder[T, TProperty], value TProperty) *RuleBuilder[T, TProperty] {
	rule := func(instance *T) *ValidationError {
		fieldValue := b.selector(instance)
		if fieldValue <= value {
			return &ValidationError{
				Field:   b.fieldName,
				Message: fmt.Sprintf("必须大于 %v", value),
			}
		}
		return nil
	}
	return b.addRule(rule)
}

// GreaterThanOrEqual 大于等于
func GreaterThanOrEqual[T any, TProperty Number](b *RuleBuilder[T, TProperty], value TProperty) *RuleBuilder[T, TProperty] {
	rule := func(instance *T) *ValidationError {
		fieldValue := b.selector(instance)
		if fieldValue < value {
			return &ValidationError{
				Field:   b.fieldName,
				Message: fmt.Sprintf("必须大于或等于 %v", value),
			}
		}
		return nil
	}
	return b.addRule(rule)
}

// LessThan 小于
func LessThan[T any, TProperty Number](b *RuleBuilder[T, TProperty], value TProperty) *RuleBuilder[T, TProperty] {
	rule := func(instance *T) *ValidationError {
		fieldValue := b.selector(instance)
		if fieldValue >= value {
			return &ValidationError{
				Field:   b.fieldName,
				Message: fmt.Sprintf("必须小于 %v", value),
			}
		}
		return nil
	}
	return b.addRule(rule)
}

// LessThanOrEqual 小于等于
func LessThanOrEqual[T any, TProperty Number](b *RuleBuilder[T, TProperty], value TProperty) *RuleBuilder[T, TProperty] {
	rule := func(instance *T) *ValidationError {
		fieldValue := b.selector(instance)
		if fieldValue > value {
			return &ValidationError{
				Field:   b.fieldName,
				Message: fmt.Sprintf("必须小于或等于 %v", value),
			}
		}
		return nil
	}
	return b.addRule(rule)
}

// InclusiveBetween 包含边界的范围
func InclusiveBetween[T any, TProperty Number](b *RuleBuilder[T, TProperty], min, max TProperty) *RuleBuilder[T, TProperty] {
	rule := func(instance *T) *ValidationError {
		value := b.selector(instance)
		if value < min || value > max {
			return &ValidationError{
				Field:   b.fieldName,
				Message: fmt.Sprintf("必须在 %v 到 %v 之间", min, max),
			}
		}
		return nil
	}
	return b.addRule(rule)
}

// ExclusiveBetween 不包含边界的范围
func ExclusiveBetween[T any, TProperty Number](b *RuleBuilder[T, TProperty], min, max TProperty) *RuleBuilder[T, TProperty] {
	rule := func(instance *T) *ValidationError {
		value := b.selector(instance)
		if value <= min || value >= max {
			return &ValidationError{
				Field:   b.fieldName,
				Message: fmt.Sprintf("必须在 %v 和 %v 之间（不包含边界）", min, max),
			}
		}
		return nil
	}
	return b.addRule(rule)
}

// MustNumber 自定义验证函数（数字类型）
func MustNumber[T any, TProperty Number](b *RuleBuilder[T, TProperty], predicate func(*T, TProperty) bool) *RuleBuilder[T, TProperty] {
	rule := func(instance *T) *ValidationError {
		value := b.selector(instance)
		if !predicate(instance, value) {
			return &ValidationError{
				Field:   b.fieldName,
				Message: "验证失败",
			}
		}
		return nil
	}
	return b.addRule(rule)
}
