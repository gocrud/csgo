package numx

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// BigInt 是一个 int64 类型，在 JSON 序列化/反序列化时作为字符串处理，以避免在 JavaScript 中丢失精度。
type BigInt int64

// MarshalJSON 实现 json.Marshaler 接口。
func (b BigInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(int64(b), 10))
}

// UnmarshalJSON 实现 json.Unmarshaler 接口。
func (b *BigInt) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid BigInt string: %w", err)
		}
		*b = BigInt(i)
		return nil
	}

	var num int64
	if err := json.Unmarshal(data, &num); err != nil {
		return fmt.Errorf("BigInt must be a string or number: %w", err)
	}
	*b = BigInt(num)
	return nil
}

// Int64 将 BigInt 转换为 int64。
func (b BigInt) Int64() int64 {
	return int64(b)
}

// String 返回 BigInt 的字符串表示。
func (b BigInt) String() string {
	return strconv.FormatInt(int64(b), 10)
}

// BigUint 是一个 uint64 类型，在 JSON 序列化/反序列化时作为字符串处理，以避免在 JavaScript 中丢失精度。
type BigUint uint64

// MarshalJSON 实现 json.Marshaler 接口。
func (b BigUint) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatUint(uint64(b), 10))
}

// UnmarshalJSON 实现 json.Unmarshaler 接口。
func (b *BigUint) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		i, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid BigUint string: %w", err)
		}
		*b = BigUint(i)
		return nil
	}

	var num uint64
	if err := json.Unmarshal(data, &num); err != nil {
		return fmt.Errorf("BigUint must be a string or number: %w", err)
	}
	*b = BigUint(num)
	return nil
}

// Uint64 将 BigUint 转换为 uint64。
func (b BigUint) Uint64() uint64 {
	return uint64(b)
}

// String 返回 BigUint 的字符串表示。
func (b BigUint) String() string {
	return strconv.FormatUint(uint64(b), 10)
}
