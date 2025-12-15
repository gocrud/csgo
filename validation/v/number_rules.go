package v

import "fmt"

// ========== Int Min 规则 ==========

type intMinRule struct {
	BaseRule
	Min int
}

func (r *intMinRule) Validate(value interface{}) error {
	if v, ok := value.(int); ok {
		return r.ValidateInt(v)
	}
	return fmt.Errorf("期望整数类型")
}

func (r *intMinRule) ValidateInt(value int) error {
	if value < r.Min {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("不能小于 %d", r.Min)
	}
	return nil
}

// Min 最小值验证
func (i Int) Min(min int) Int {
	rule := &intMinRule{Min: min}
	i.tracker.addIntRule(i.fieldPath, rule)
	return i
}

// ========== Int Max 规则 ==========

type intMaxRule struct {
	BaseRule
	Max int
}

func (r *intMaxRule) Validate(value interface{}) error {
	if v, ok := value.(int); ok {
		return r.ValidateInt(v)
	}
	return fmt.Errorf("期望整数类型")
}

func (r *intMaxRule) ValidateInt(value int) error {
	if value > r.Max {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("不能大于 %d", r.Max)
	}
	return nil
}

// Max 最大值验证
func (i Int) Max(max int) Int {
	rule := &intMaxRule{Max: max}
	i.tracker.addIntRule(i.fieldPath, rule)
	return i
}

// ========== Int Range 规则 ==========

type intRangeRule struct {
	BaseRule
	Min int
	Max int
}

func (r *intRangeRule) Validate(value interface{}) error {
	if v, ok := value.(int); ok {
		return r.ValidateInt(v)
	}
	return fmt.Errorf("期望整数类型")
}

func (r *intRangeRule) ValidateInt(value int) error {
	if value < r.Min || value > r.Max {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("必须在 %d 到 %d 之间", r.Min, r.Max)
	}
	return nil
}

// Range 范围验证
func (i Int) Range(min, max int) Int {
	rule := &intRangeRule{Min: min, Max: max}
	i.tracker.addIntRule(i.fieldPath, rule)
	return i
}

// Msg 设置最后一个规则的错误消息
func (i Int) Msg(msg string) Int {
	i.tracker.setLastMessage(i.fieldPath, msg)
	return i
}

// ========== Int64 Min 规则 ==========

type int64MinRule struct {
	BaseRule
	Min int64
}

func (r *int64MinRule) Validate(value interface{}) error {
	if v, ok := value.(int64); ok {
		return r.ValidateInt64(v)
	}
	return fmt.Errorf("期望 int64 类型")
}

func (r *int64MinRule) ValidateInt64(value int64) error {
	if value < r.Min {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("不能小于 %d", r.Min)
	}
	return nil
}

// Min 最小值验证
func (i Int64) Min(min int64) Int64 {
	rule := &int64MinRule{Min: min}
	i.tracker.addInt64Rule(i.fieldPath, rule)
	return i
}

// ========== Int64 Max 规则 ==========

type int64MaxRule struct {
	BaseRule
	Max int64
}

func (r *int64MaxRule) Validate(value interface{}) error {
	if v, ok := value.(int64); ok {
		return r.ValidateInt64(v)
	}
	return fmt.Errorf("期望 int64 类型")
}

func (r *int64MaxRule) ValidateInt64(value int64) error {
	if value > r.Max {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("不能大于 %d", r.Max)
	}
	return nil
}

// Max 最大值验证
func (i Int64) Max(max int64) Int64 {
	rule := &int64MaxRule{Max: max}
	i.tracker.addInt64Rule(i.fieldPath, rule)
	return i
}

// ========== Int64 Range 规则 ==========

type int64RangeRule struct {
	BaseRule
	Min int64
	Max int64
}

func (r *int64RangeRule) Validate(value interface{}) error {
	if v, ok := value.(int64); ok {
		return r.ValidateInt64(v)
	}
	return fmt.Errorf("期望 int64 类型")
}

func (r *int64RangeRule) ValidateInt64(value int64) error {
	if value < r.Min || value > r.Max {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("必须在 %d 到 %d 之间", r.Min, r.Max)
	}
	return nil
}

// Range 范围验证
func (i Int64) Range(min, max int64) Int64 {
	rule := &int64RangeRule{Min: min, Max: max}
	i.tracker.addInt64Rule(i.fieldPath, rule)
	return i
}

// Msg 设置最后一个规则的错误消息
func (i Int64) Msg(msg string) Int64 {
	i.tracker.setLastMessage(i.fieldPath, msg)
	return i
}

// ========== Float64 Min 规则 ==========

type float64MinRule struct {
	BaseRule
	Min float64
}

func (r *float64MinRule) Validate(value interface{}) error {
	if v, ok := value.(float64); ok {
		return r.ValidateFloat64(v)
	}
	return fmt.Errorf("期望 float64 类型")
}

func (r *float64MinRule) ValidateFloat64(value float64) error {
	if value < r.Min {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("不能小于 %f", r.Min)
	}
	return nil
}

// Min 最小值验证
func (f Float64) Min(min float64) Float64 {
	rule := &float64MinRule{Min: min}
	f.tracker.addFloat64Rule(f.fieldPath, rule)
	return f
}

// ========== Float64 Max 规则 ==========

type float64MaxRule struct {
	BaseRule
	Max float64
}

func (r *float64MaxRule) Validate(value interface{}) error {
	if v, ok := value.(float64); ok {
		return r.ValidateFloat64(v)
	}
	return fmt.Errorf("期望 float64 类型")
}

func (r *float64MaxRule) ValidateFloat64(value float64) error {
	if value > r.Max {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("不能大于 %f", r.Max)
	}
	return nil
}

// Max 最大值验证
func (f Float64) Max(max float64) Float64 {
	rule := &float64MaxRule{Max: max}
	f.tracker.addFloat64Rule(f.fieldPath, rule)
	return f
}

// ========== Float64 Range 规则 ==========

type float64RangeRule struct {
	BaseRule
	Min float64
	Max float64
}

func (r *float64RangeRule) Validate(value interface{}) error {
	if v, ok := value.(float64); ok {
		return r.ValidateFloat64(v)
	}
	return fmt.Errorf("期望 float64 类型")
}

func (r *float64RangeRule) ValidateFloat64(value float64) error {
	if value < r.Min || value > r.Max {
		if r.Message != "" {
			return fmt.Errorf("%s", r.Message)
		}
		return fmt.Errorf("必须在 %f 到 %f 之间", r.Min, r.Max)
	}
	return nil
}

// Range 范围验证
func (f Float64) Range(min, max float64) Float64 {
	rule := &float64RangeRule{Min: min, Max: max}
	f.tracker.addFloat64Rule(f.fieldPath, rule)
	return f
}

// Msg 设置最后一个规则的错误消息
func (f Float64) Msg(msg string) Float64 {
	f.tracker.setLastMessage(f.fieldPath, msg)
	return f
}
