package validation

import (
	"reflect"
)

// RuleType 定义规则类型
type RuleType string

// Rule 表示单个验证规则
type Rule struct {
	Type      RuleType
	Params    []interface{}
	CustomMsg string // 自定义错误消息
}

// Schema 存储类型的验证模式
type Schema struct {
	// Offset -> Rules
	FieldRules map[uintptr][]Rule
	// Offset -> Field Name (JSON tag or Struct Field Name)
	FieldNames map[uintptr]string
	// Offset -> Field Kind (用于 Required 等通用规则判断类型)
	FieldKinds map[uintptr]reflect.Kind
	// Offset -> Field Type (用于获取字段详细类型信息，如 Time, Slice 元素类型)
	FieldTypes map[uintptr]reflect.Type
	// OrderedOffsets 预先排序的偏移量列表，用于保证验证顺序一致性
	OrderedOffsets []uintptr
	// BaseType
	Type reflect.Type
	// FailFast 是否快速失败（遇到第一个错误就返回）
	FailFast bool
}

// 录制器上下文
type recorderContext struct {
	BasePtr uintptr
	Rules   map[uintptr][]Rule
}
