package v

import (
	"time"
	"unsafe"

	"github.com/gocrud/csgo/validation"
)

// ==================== Integers ====================

type Int8 int8

func (i *Int8) Min(min int) *Int8   { validation.RecordRule(unsafe.Pointer(i), "min", min); return i }
func (i *Int8) Max(max int) *Int8   { validation.RecordRule(unsafe.Pointer(i), "max", max); return i }
func (i *Int8) Required() *Int8     { validation.RecordRule(unsafe.Pointer(i), "required"); return i }
func (i *Int8) Equal(val int) *Int8 { validation.RecordRule(unsafe.Pointer(i), "eq", val); return i }
func (i *Int8) In(vals ...int) *Int8 {
	args := make([]interface{}, len(vals))
	for idx, v := range vals {
		args[idx] = v
	}
	validation.RecordRule(unsafe.Pointer(i), "in", args...)
	return i
}
func (i *Int8) Msg(msg string) *Int8 { validation.SetLastRuleMsg(unsafe.Pointer(i), msg); return i }

type Int16 int16

func (i *Int16) Min(min int) *Int16   { validation.RecordRule(unsafe.Pointer(i), "min", min); return i }
func (i *Int16) Max(max int) *Int16   { validation.RecordRule(unsafe.Pointer(i), "max", max); return i }
func (i *Int16) Required() *Int16     { validation.RecordRule(unsafe.Pointer(i), "required"); return i }
func (i *Int16) Equal(val int) *Int16 { validation.RecordRule(unsafe.Pointer(i), "eq", val); return i }
func (i *Int16) In(vals ...int) *Int16 {
	args := make([]interface{}, len(vals))
	for idx, v := range vals {
		args[idx] = v
	}
	validation.RecordRule(unsafe.Pointer(i), "in", args...)
	return i
}
func (i *Int16) Msg(msg string) *Int16 { validation.SetLastRuleMsg(unsafe.Pointer(i), msg); return i }

type Int32 int32

func (i *Int32) Min(min int) *Int32   { validation.RecordRule(unsafe.Pointer(i), "min", min); return i }
func (i *Int32) Max(max int) *Int32   { validation.RecordRule(unsafe.Pointer(i), "max", max); return i }
func (i *Int32) Required() *Int32     { validation.RecordRule(unsafe.Pointer(i), "required"); return i }
func (i *Int32) Equal(val int) *Int32 { validation.RecordRule(unsafe.Pointer(i), "eq", val); return i }
func (i *Int32) In(vals ...int) *Int32 {
	args := make([]interface{}, len(vals))
	for idx, v := range vals {
		args[idx] = v
	}
	validation.RecordRule(unsafe.Pointer(i), "in", args...)
	return i
}
func (i *Int32) Msg(msg string) *Int32 { validation.SetLastRuleMsg(unsafe.Pointer(i), msg); return i }

type Int64 int64

func (i *Int64) Min(min int64) *Int64 { validation.RecordRule(unsafe.Pointer(i), "min", min); return i }
func (i *Int64) Max(max int64) *Int64 { validation.RecordRule(unsafe.Pointer(i), "max", max); return i }
func (i *Int64) Required() *Int64     { validation.RecordRule(unsafe.Pointer(i), "required"); return i }
func (i *Int64) Equal(val int64) *Int64 {
	validation.RecordRule(unsafe.Pointer(i), "eq", val)
	return i
}
func (i *Int64) In(vals ...int64) *Int64 {
	args := make([]interface{}, len(vals))
	for idx, v := range vals {
		args[idx] = v
	}
	validation.RecordRule(unsafe.Pointer(i), "in", args...)
	return i
}
func (i *Int64) Msg(msg string) *Int64 { validation.SetLastRuleMsg(unsafe.Pointer(i), msg); return i }

type Uint uint

func (i *Uint) Min(min uint) *Uint   { validation.RecordRule(unsafe.Pointer(i), "min", min); return i }
func (i *Uint) Max(max uint) *Uint   { validation.RecordRule(unsafe.Pointer(i), "max", max); return i }
func (i *Uint) Required() *Uint      { validation.RecordRule(unsafe.Pointer(i), "required"); return i }
func (i *Uint) Equal(val uint) *Uint { validation.RecordRule(unsafe.Pointer(i), "eq", val); return i }
func (i *Uint) In(vals ...uint) *Uint {
	args := make([]interface{}, len(vals))
	for idx, v := range vals {
		args[idx] = v
	}
	validation.RecordRule(unsafe.Pointer(i), "in", args...)
	return i
}
func (i *Uint) Msg(msg string) *Uint { validation.SetLastRuleMsg(unsafe.Pointer(i), msg); return i }

type Uint8 uint8

func (i *Uint8) Min(min uint) *Uint8   { validation.RecordRule(unsafe.Pointer(i), "min", min); return i }
func (i *Uint8) Max(max uint) *Uint8   { validation.RecordRule(unsafe.Pointer(i), "max", max); return i }
func (i *Uint8) Required() *Uint8      { validation.RecordRule(unsafe.Pointer(i), "required"); return i }
func (i *Uint8) Equal(val uint) *Uint8 { validation.RecordRule(unsafe.Pointer(i), "eq", val); return i }
func (i *Uint8) In(vals ...uint) *Uint8 {
	args := make([]interface{}, len(vals))
	for idx, v := range vals {
		args[idx] = v
	}
	validation.RecordRule(unsafe.Pointer(i), "in", args...)
	return i
}
func (i *Uint8) Msg(msg string) *Uint8 { validation.SetLastRuleMsg(unsafe.Pointer(i), msg); return i }

type Uint16 uint16

func (i *Uint16) Min(min uint) *Uint16 {
	validation.RecordRule(unsafe.Pointer(i), "min", min)
	return i
}
func (i *Uint16) Max(max uint) *Uint16 {
	validation.RecordRule(unsafe.Pointer(i), "max", max)
	return i
}
func (i *Uint16) Required() *Uint16 { validation.RecordRule(unsafe.Pointer(i), "required"); return i }
func (i *Uint16) Equal(val uint) *Uint16 {
	validation.RecordRule(unsafe.Pointer(i), "eq", val)
	return i
}
func (i *Uint16) In(vals ...uint) *Uint16 {
	args := make([]interface{}, len(vals))
	for idx, v := range vals {
		args[idx] = v
	}
	validation.RecordRule(unsafe.Pointer(i), "in", args...)
	return i
}
func (i *Uint16) Msg(msg string) *Uint16 { validation.SetLastRuleMsg(unsafe.Pointer(i), msg); return i }

type Uint32 uint32

func (i *Uint32) Min(min uint) *Uint32 {
	validation.RecordRule(unsafe.Pointer(i), "min", min)
	return i
}
func (i *Uint32) Max(max uint) *Uint32 {
	validation.RecordRule(unsafe.Pointer(i), "max", max)
	return i
}
func (i *Uint32) Required() *Uint32 { validation.RecordRule(unsafe.Pointer(i), "required"); return i }
func (i *Uint32) Equal(val uint) *Uint32 {
	validation.RecordRule(unsafe.Pointer(i), "eq", val)
	return i
}
func (i *Uint32) In(vals ...uint) *Uint32 {
	args := make([]interface{}, len(vals))
	for idx, v := range vals {
		args[idx] = v
	}
	validation.RecordRule(unsafe.Pointer(i), "in", args...)
	return i
}
func (i *Uint32) Msg(msg string) *Uint32 { validation.SetLastRuleMsg(unsafe.Pointer(i), msg); return i }

type Uint64 uint64

func (i *Uint64) Min(min uint64) *Uint64 {
	validation.RecordRule(unsafe.Pointer(i), "min", min)
	return i
}
func (i *Uint64) Max(max uint64) *Uint64 {
	validation.RecordRule(unsafe.Pointer(i), "max", max)
	return i
}
func (i *Uint64) Required() *Uint64 { validation.RecordRule(unsafe.Pointer(i), "required"); return i }
func (i *Uint64) Equal(val uint64) *Uint64 {
	validation.RecordRule(unsafe.Pointer(i), "eq", val)
	return i
}
func (i *Uint64) In(vals ...uint64) *Uint64 {
	args := make([]interface{}, len(vals))
	for idx, v := range vals {
		args[idx] = v
	}
	validation.RecordRule(unsafe.Pointer(i), "in", args...)
	return i
}
func (i *Uint64) Msg(msg string) *Uint64 { validation.SetLastRuleMsg(unsafe.Pointer(i), msg); return i }

// ==================== Floats ====================

type Float32 float32

func (f *Float32) Min(min float64) *Float32 {
	validation.RecordRule(unsafe.Pointer(f), "min", min)
	return f
}
func (f *Float32) Max(max float64) *Float32 {
	validation.RecordRule(unsafe.Pointer(f), "max", max)
	return f
}
func (f *Float32) Required() *Float32 { validation.RecordRule(unsafe.Pointer(f), "required"); return f }
func (f *Float32) Equal(val float64) *Float32 {
	validation.RecordRule(unsafe.Pointer(f), "eq", val)
	return f
}
func (f *Float32) Msg(msg string) *Float32 {
	validation.SetLastRuleMsg(unsafe.Pointer(f), msg)
	return f
}

type Float64 float64

func (f *Float64) Min(min float64) *Float64 {
	validation.RecordRule(unsafe.Pointer(f), "min", min)
	return f
}
func (f *Float64) Max(max float64) *Float64 {
	validation.RecordRule(unsafe.Pointer(f), "max", max)
	return f
}
func (f *Float64) Required() *Float64 { validation.RecordRule(unsafe.Pointer(f), "required"); return f }
func (f *Float64) Equal(val float64) *Float64 {
	validation.RecordRule(unsafe.Pointer(f), "eq", val)
	return f
}
func (f *Float64) Msg(msg string) *Float64 {
	validation.SetLastRuleMsg(unsafe.Pointer(f), msg)
	return f
}

// ==================== Bool ====================

type Bool bool

func (b *Bool) True() *Bool          { validation.RecordRule(unsafe.Pointer(b), "true"); return b }
func (b *Bool) False() *Bool         { validation.RecordRule(unsafe.Pointer(b), "false"); return b }
func (b *Bool) Msg(msg string) *Bool { validation.SetLastRuleMsg(unsafe.Pointer(b), msg); return b }

// ==================== Generic Slice ====================

type Slice[T any] []T

func (s *Slice[T]) MinLen(min int) *Slice[T] {
	validation.RecordRule(unsafe.Pointer(s), "min_len", min)
	return s
}

func (s *Slice[T]) MaxLen(max int) *Slice[T] {
	validation.RecordRule(unsafe.Pointer(s), "max_len", max)
	return s
}

func (s *Slice[T]) Len(len int) *Slice[T] {
	validation.RecordRule(unsafe.Pointer(s), "len", len)
	return s
}

func (s *Slice[T]) Required() *Slice[T] {
	validation.RecordRule(unsafe.Pointer(s), "required")
	return s
}

func (s *Slice[T]) Unique() *Slice[T] {
	// Unique 需要访问元素，这在 generic slice + unsafe pointer 场景下比较困难
	// 我们只能通过 reflect.Value 来做，这在 validator.go 中会处理
	validation.RecordRule(unsafe.Pointer(s), "unique")
	return s
}

func (s *Slice[T]) Msg(msg string) *Slice[T] {
	validation.SetLastRuleMsg(unsafe.Pointer(s), msg)
	return s
}

// ==================== Time ====================

type Time time.Time

func (t *Time) After(val time.Time) *Time {
	validation.RecordRule(unsafe.Pointer(t), "after", val)
	return t
}
func (t *Time) Before(val time.Time) *Time {
	validation.RecordRule(unsafe.Pointer(t), "before", val)
	return t
}
func (t *Time) Required() *Time {
	validation.RecordRule(unsafe.Pointer(t), "required")
	return t
}
func (t *Time) Msg(msg string) *Time {
	validation.SetLastRuleMsg(unsafe.Pointer(t), msg)
	return t
}
