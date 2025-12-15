package v

// Rule 验证规则接口
type Rule interface {
	Validate(value interface{}) error
	GetMessage() string
	SetMessage(msg string)
}

// BaseRule 基础规则（包含错误消息）
type BaseRule struct {
	Message string
}

// GetMessage 获取错误消息
func (r *BaseRule) GetMessage() string {
	return r.Message
}

// SetMessage 设置错误消息
func (r *BaseRule) SetMessage(msg string) {
	r.Message = msg
}

// StringRule 字符串验证规则接口
type StringRule interface {
	Rule
	ValidateString(value string) error
}

// IntRule 整数验证规则接口
type IntRule interface {
	Rule
	ValidateInt(value int) error
}

// Int64Rule int64 验证规则接口
type Int64Rule interface {
	Rule
	ValidateInt64(value int64) error
}

// Float64Rule float64 验证规则接口
type Float64Rule interface {
	Rule
	ValidateFloat64(value float64) error
}

// SliceRule 切片验证规则接口
type SliceRule interface {
	Rule
	ValidateSlice(value interface{}) error
}
