package v

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// ========== MinLen 规则 ==========

type minLenRule struct {
	BaseRule
	Min int
}

func (r *minLenRule) Validate(value interface{}) error {
	if str, ok := value.(string); ok {
		return r.ValidateString(str)
	}
	return fmt.Errorf("期望字符串类型")
}

func (r *minLenRule) ValidateString(value string) error {
	length := utf8.RuneCountInString(value)
	if length < r.Min {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("长度不能少于 %d", r.Min)
	}
	return nil
}

// MinLen 最小长度验证
func (s String) MinLen(min int) String {
	rule := &minLenRule{Min: min}
	s.tracker.addStringRule(s.fieldPath, rule)
	return s
}

// ========== MaxLen 规则 ==========

type maxLenRule struct {
	BaseRule
	Max int
}

func (r *maxLenRule) Validate(value interface{}) error {
	if str, ok := value.(string); ok {
		return r.ValidateString(str)
	}
	return fmt.Errorf("期望字符串类型")
}

func (r *maxLenRule) ValidateString(value string) error {
	length := utf8.RuneCountInString(value)
	if length > r.Max {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("长度不能超过 %d", r.Max)
	}
	return nil
}

// MaxLen 最大长度验证
func (s String) MaxLen(max int) String {
	rule := &maxLenRule{Max: max}
	s.tracker.addStringRule(s.fieldPath, rule)
	return s
}

// ========== NotEmpty 规则 ==========

type notEmptyRule struct {
	BaseRule
}

func (r *notEmptyRule) Validate(value interface{}) error {
	if str, ok := value.(string); ok {
		return r.ValidateString(str)
	}
	return fmt.Errorf("期望字符串类型")
}

func (r *notEmptyRule) ValidateString(value string) error {
	if strings.TrimSpace(value) == "" {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("不能为空")
	}
	return nil
}

// NotEmpty 非空验证
func (s String) NotEmpty() String {
	rule := &notEmptyRule{}
	s.tracker.addStringRule(s.fieldPath, rule)
	return s
}

// ========== Email 规则 ==========

type emailRule struct {
	BaseRule
	regex *regexp.Regexp
}

func (r *emailRule) Validate(value interface{}) error {
	if str, ok := value.(string); ok {
		return r.ValidateString(str)
	}
	return fmt.Errorf("期望字符串类型")
}

func (r *emailRule) ValidateString(value string) error {
	if value != "" && !r.regex.MatchString(value) {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("邮箱格式不正确")
	}
	return nil
}

// Email 邮箱格式验证
func (s String) Email() String {
	rule := &emailRule{
		regex: regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`),
	}
	s.tracker.addStringRule(s.fieldPath, rule)
	return s
}

// ========== Pattern 规则 ==========

type patternRule struct {
	BaseRule
	Pattern string
	regex   *regexp.Regexp
}

func (r *patternRule) Validate(value interface{}) error {
	if str, ok := value.(string); ok {
		return r.ValidateString(str)
	}
	return fmt.Errorf("期望字符串类型")
}

func (r *patternRule) ValidateString(value string) error {
	if value != "" && !r.regex.MatchString(value) {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("格式不正确")
	}
	return nil
}

// Pattern 正则表达式验证
func (s String) Pattern(pattern string) String {
	rule := &patternRule{
		Pattern: pattern,
		regex:   regexp.MustCompile(pattern),
	}
	s.tracker.addStringRule(s.fieldPath, rule)
	return s
}

// ========== Msg 方法 ==========

// Msg 设置最后一个规则的错误消息
func (s String) Msg(msg string) String {
	s.tracker.setLastMessage(s.fieldPath, msg)
	return s
}
