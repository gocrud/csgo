package validation

import "fmt"

// NotEmptySlice 集合非空
func NotEmptySlice[T any, TItem any](b *RuleBuilder[T, []TItem]) *RuleBuilder[T, []TItem] {
	rule := func(instance *T) *ValidationError {
		value := b.selector(instance)
		if len(value) == 0 {
			return &ValidationError{
				Field:   b.fieldName,
				Message: "集合不能为空",
			}
		}
		return nil
	}
	return b.addRule(rule)
}

// MinLengthSlice 集合最小长度
func MinLengthSlice[T any, TItem any](b *RuleBuilder[T, []TItem], min int) *RuleBuilder[T, []TItem] {
	rule := func(instance *T) *ValidationError {
		value := b.selector(instance)
		if len(value) < min {
			return &ValidationError{
				Field:   b.fieldName,
				Message: fmt.Sprintf("集合长度不能少于 %d", min),
			}
		}
		return nil
	}
	return b.addRule(rule)
}

// MaxLengthSlice 集合最大长度
func MaxLengthSlice[T any, TItem any](b *RuleBuilder[T, []TItem], max int) *RuleBuilder[T, []TItem] {
	rule := func(instance *T) *ValidationError {
		value := b.selector(instance)
		if len(value) > max {
			return &ValidationError{
				Field:   b.fieldName,
				Message: fmt.Sprintf("集合长度不能超过 %d", max),
			}
		}
		return nil
	}
	return b.addRule(rule)
}

// MustSlice 集合自定义验证
func MustSlice[T any, TItem any](b *RuleBuilder[T, []TItem], predicate func(*T, []TItem) bool) *RuleBuilder[T, []TItem] {
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
