package v

import (
	"unsafe"

	"github.com/gocrud/csgo/validation"
)

// ==================== Int ====================

type Int int

func (i *Int) Min(min int) *Int {
	validation.RecordRule(unsafe.Pointer(i), "min", min)
	return i
}

func (i *Int) Max(max int) *Int {
	validation.RecordRule(unsafe.Pointer(i), "max", max)
	return i
}

func (i *Int) Range(min, max int) *Int {
	validation.RecordRule(unsafe.Pointer(i), "range", min, max)
	return i
}

func (i *Int) Required() *Int {
	validation.RecordRule(unsafe.Pointer(i), "required")
	return i
}

func (i *Int) Equal(val int) *Int {
	validation.RecordRule(unsafe.Pointer(i), "eq", val)
	return i
}

// In 验证值是否在列表中 (替代 OneOf)
func (i *Int) In(vals ...int) *Int {
	// 将 []int 转换为 []interface{}
	args := make([]interface{}, len(vals))
	for idx, v := range vals {
		args[idx] = v
	}
	validation.RecordRule(unsafe.Pointer(i), "in", args...)
	return i
}

func (i *Int) Msg(msg string) *Int {
	validation.SetLastRuleMsg(unsafe.Pointer(i), msg)
	return i
}

func (i *Int) MsgGroup(msg string) *Int {
	validation.SetGroupMsg(unsafe.Pointer(i), msg)
	return i
}

// ==================== String ====================

type String string

func (s *String) MinLen(min int) *String {
	validation.RecordRule(unsafe.Pointer(s), "min_len", min)
	return s
}

func (s *String) MaxLen(max int) *String {
	validation.RecordRule(unsafe.Pointer(s), "max_len", max)
	return s
}

func (s *String) Len(len int) *String {
	validation.RecordRule(unsafe.Pointer(s), "len", len)
	return s
}

func (s *String) RangeLen(min, max int) *String {
	validation.RecordRule(unsafe.Pointer(s), "range_len", min, max)
	return s
}

func (s *String) Required() *String {
	validation.RecordRule(unsafe.Pointer(s), "required")
	return s
}

func (s *String) Email() *String {
	validation.RecordRule(unsafe.Pointer(s), "email")
	return s
}

func (s *String) URL() *String {
	validation.RecordRule(unsafe.Pointer(s), "url")
	return s
}

func (s *String) IP() *String {
	validation.RecordRule(unsafe.Pointer(s), "ip")
	return s
}

func (s *String) UUID() *String {
	validation.RecordRule(unsafe.Pointer(s), "uuid")
	return s
}

func (s *String) Regex(pattern string) *String {
	validation.RecordRule(unsafe.Pointer(s), "regex", pattern)
	return s
}

func (s *String) Alpha() *String {
	validation.RecordRule(unsafe.Pointer(s), "alpha")
	return s
}

func (s *String) AlphaNum() *String {
	validation.RecordRule(unsafe.Pointer(s), "alphanum")
	return s
}

func (s *String) Numeric() *String {
	validation.RecordRule(unsafe.Pointer(s), "numeric")
	return s
}

func (s *String) Lowercase() *String {
	validation.RecordRule(unsafe.Pointer(s), "lowercase")
	return s
}

func (s *String) Uppercase() *String {
	validation.RecordRule(unsafe.Pointer(s), "uppercase")
	return s
}

func (s *String) Contains(sub string) *String {
	validation.RecordRule(unsafe.Pointer(s), "contains", sub)
	return s
}

func (s *String) StartsWith(prefix string) *String {
	validation.RecordRule(unsafe.Pointer(s), "startswith", prefix)
	return s
}

func (s *String) EndsWith(suffix string) *String {
	validation.RecordRule(unsafe.Pointer(s), "endswith", suffix)
	return s
}

func (s *String) In(vals ...string) *String {
	args := make([]interface{}, len(vals))
	for idx, v := range vals {
		args[idx] = v
	}
	validation.RecordRule(unsafe.Pointer(s), "in", args...)
	return s
}

func (s *String) Msg(msg string) *String {
	validation.SetLastRuleMsg(unsafe.Pointer(s), msg)
	return s
}

func (s *String) MsgGroup(msg string) *String {
	validation.SetGroupMsg(unsafe.Pointer(s), msg)
	return s
}
