package v

import (
	"fmt"
	"reflect"
)

// ========== Slice MinLen 规则 ==========

type sliceMinLenRule struct {
	BaseRule
	Min int
}

func (r *sliceMinLenRule) Validate(value interface{}) error {
	return r.ValidateSlice(value)
}

func (r *sliceMinLenRule) ValidateSlice(value interface{}) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("期望切片类型")
	}

	length := v.Len()
	if length < r.Min {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("长度不能少于 %d", r.Min)
	}
	return nil
}

// MinLen 最小长度验证
func (s Slice[T]) MinLen(min int) Slice[T] {
	rule := &sliceMinLenRule{Min: min}
	s.tracker.addSliceRule(s.fieldPath, rule)
	return s
}

// ========== Slice MaxLen 规则 ==========

type sliceMaxLenRule struct {
	BaseRule
	Max int
}

func (r *sliceMaxLenRule) Validate(value interface{}) error {
	return r.ValidateSlice(value)
}

func (r *sliceMaxLenRule) ValidateSlice(value interface{}) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("期望切片类型")
	}

	length := v.Len()
	if length > r.Max {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("长度不能超过 %d", r.Max)
	}
	return nil
}

// MaxLen 最大长度验证
func (s Slice[T]) MaxLen(max int) Slice[T] {
	rule := &sliceMaxLenRule{Max: max}
	s.tracker.addSliceRule(s.fieldPath, rule)
	return s
}

// ========== Slice NotEmpty 规则 ==========

type sliceNotEmptyRule struct {
	BaseRule
}

func (r *sliceNotEmptyRule) Validate(value interface{}) error {
	return r.ValidateSlice(value)
}

func (r *sliceNotEmptyRule) ValidateSlice(value interface{}) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return fmt.Errorf("期望切片类型")
	}

	if v.Len() == 0 {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("不能为空")
	}
	return nil
}

// NotEmpty 非空验证
func (s Slice[T]) NotEmpty() Slice[T] {
	rule := &sliceNotEmptyRule{}
	s.tracker.addSliceRule(s.fieldPath, rule)
	return s
}

// Msg 设置最后一个规则的错误消息
func (s Slice[T]) Msg(msg string) Slice[T] {
	s.tracker.setLastMessage(s.fieldPath, msg)
	return s
}
