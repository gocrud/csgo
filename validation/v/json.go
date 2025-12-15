package v

import (
	"encoding/json"
	"strconv"
)

// ========== String JSON 序列化 ==========

// MarshalJSON 实现 json.Marshaler 接口
func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.value)
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (s *String) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.value)
}

// ========== Int JSON 序列化 ==========

// MarshalJSON 实现 json.Marshaler 接口
func (i Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.value)
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (i *Int) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &i.value)
}

// ========== Int64 JSON 序列化 ==========

// MarshalJSON 实现 json.Marshaler 接口
func (i Int64) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.value)
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (i *Int64) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &i.value)
}

// ========== Float64 JSON 序列化 ==========

// MarshalJSON 实现 json.Marshaler 接口
func (f Float64) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.value)
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (f *Float64) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &f.value)
}

// ========== Bool JSON 序列化 ==========

// MarshalJSON 实现 json.Marshaler 接口
func (b Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.value)
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (b *Bool) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &b.value)
}

// ========== Slice JSON 序列化 ==========

// MarshalJSON 实现 json.Marshaler 接口
func (s Slice[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.value)
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (s *Slice[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.value)
}

// ========== 辅助函数 ==========

// String 返回字符串表示
func (s String) String() string {
	return s.value
}

// String 返回字符串表示
func (i Int) String() string {
	return strconv.Itoa(i.value)
}

// String 返回字符串表示
func (i Int64) String() string {
	return strconv.FormatInt(i.value, 10)
}

// String 返回字符串表示
func (f Float64) String() string {
	return strconv.FormatFloat(f.value, 'f', -1, 64)
}

// String 返回字符串表示
func (b Bool) String() string {
	return strconv.FormatBool(b.value)
}
