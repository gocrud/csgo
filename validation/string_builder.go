package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gocrud/csgo/errors"
)

// StringRuleBuilder 字符串规则构建器（特化版本）
type StringRuleBuilder[T any] struct {
	*RuleBuilder[T, string]
}

// NotEmpty 非空规则
func (b *StringRuleBuilder[T]) NotEmpty() *StringRuleBuilder[T] {
	rule := func(instance *T) *ValidationError {
		value := b.selector(instance)
		// 在闭包内，value 的类型已经确定
		if strings.TrimSpace(value) == "" {
			return &ValidationError{
				Field:   b.fieldName,
				Message: "不能为空",
				Code:    errors.ValidationRequired,
			}
		}
		return nil
	}
	b.addRule(rule)
	return b
}

// Length 长度限制
func (b *StringRuleBuilder[T]) Length(min, max int) *StringRuleBuilder[T] {
	rule := func(instance *T) *ValidationError {
		value := b.selector(instance)
		length := utf8.RuneCountInString(value)
		if length < min || length > max {
			return &ValidationError{
				Field:   b.fieldName,
				Message: fmt.Sprintf("长度必须在 %d 到 %d 之间", min, max),
				Code:    errors.ValidationLength,
			}
		}
		return nil
	}
	b.addRule(rule)
	return b
}

// MinLength 最小长度
func (b *StringRuleBuilder[T]) MinLength(min int) *StringRuleBuilder[T] {
	rule := func(instance *T) *ValidationError {
		value := b.selector(instance)
		length := utf8.RuneCountInString(value)
		if length < min {
			return &ValidationError{
				Field:   b.fieldName,
				Message: fmt.Sprintf("长度不能少于 %d", min),
				Code:    errors.ValidationMinLength,
			}
		}
		return nil
	}
	b.addRule(rule)
	return b
}

// MaxLength 最大长度
func (b *StringRuleBuilder[T]) MaxLength(max int) *StringRuleBuilder[T] {
	rule := func(instance *T) *ValidationError {
		value := b.selector(instance)
		length := utf8.RuneCountInString(value)
		if length > max {
			return &ValidationError{
				Field:   b.fieldName,
				Message: fmt.Sprintf("长度不能超过 %d", max),
				Code:    errors.ValidationMaxLength,
			}
		}
		return nil
	}
	b.addRule(rule)
	return b
}

// EmailAddress 邮箱验证
func (b *StringRuleBuilder[T]) EmailAddress() *StringRuleBuilder[T] {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	rule := func(instance *T) *ValidationError {
		value := b.selector(instance)
		if value != "" && !emailRegex.MatchString(value) {
			return &ValidationError{
				Field:   b.fieldName,
				Message: "邮箱格式不正确",
				Code:    errors.ValidationEmail,
			}
		}
		return nil
	}
	b.addRule(rule)
	return b
}

// Matches 正则匹配
func (b *StringRuleBuilder[T]) Matches(pattern string) *StringRuleBuilder[T] {
	regex := regexp.MustCompile(pattern)

	rule := func(instance *T) *ValidationError {
		value := b.selector(instance)
		if value != "" && !regex.MatchString(value) {
			return &ValidationError{
				Field:   b.fieldName,
				Message: "格式不正确",
				Code:    errors.ValidationPattern,
			}
		}
		return nil
	}
	b.addRule(rule)
	return b
}

// MustString 自定义验证函数
func (b *StringRuleBuilder[T]) MustString(predicate func(*T, string) bool) *StringRuleBuilder[T] {
	rule := func(instance *T) *ValidationError {
		value := b.selector(instance)
		if !predicate(instance, value) {
			return &ValidationError{
				Field:   b.fieldName,
				Message: "验证失败",
				Code:    errors.ValidationFailed,
			}
		}
		return nil
	}
	b.addRule(rule)
	return b
}

// WithMessage 自定义错误消息
func (b *StringRuleBuilder[T]) WithMessage(message string) *StringRuleBuilder[T] {
	b.RuleBuilder.WithMessage(message)
	return b
}

// WithCode 设置错误码
func (b *StringRuleBuilder[T]) WithCode(code string) *StringRuleBuilder[T] {
	b.RuleBuilder.WithCode(code)
	return b
}

// When 条件验证
func (b *StringRuleBuilder[T]) When(condition func(*T) bool) *StringRuleBuilder[T] {
	b.RuleBuilder.When(condition)
	return b
}

// Unless 反向条件验证
func (b *StringRuleBuilder[T]) Unless(condition func(*T) bool) *StringRuleBuilder[T] {
	b.RuleBuilder.Unless(condition)
	return b
}
