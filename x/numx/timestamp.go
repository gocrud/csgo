package numx

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Timestamp is an int64 Unix timestamp (milliseconds) that serializes to/from JSON as a string.
type Timestamp int64

// MarshalJSON implements json.Marshaler interface.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(int64(t), 10))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid Timestamp string: %w", err)
		}
		*t = Timestamp(i)
		return nil
	}

	var num int64
	if err := json.Unmarshal(data, &num); err != nil {
		return fmt.Errorf("Timestamp must be a string or number: %w", err)
	}
	*t = Timestamp(num)
	return nil
}

// Value implements driver.Valuer interface.
func (t Timestamp) Value() (driver.Value, error) {
	return int64(t), nil
}

// Scan implements sql.Scanner interface.
func (t *Timestamp) Scan(value interface{}) error {
	if value == nil {
		*t = 0
		return nil
	}

	switch v := value.(type) {
	case int64:
		*t = Timestamp(v)
	case int:
		*t = Timestamp(v)
	case int32:
		*t = Timestamp(v)
	case uint64:
		*t = Timestamp(v)
	case uint32:
		*t = Timestamp(v)
	case []byte:
		i, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to scan Timestamp from bytes: %w", err)
		}
		*t = Timestamp(i)
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to scan Timestamp from string: %w", err)
		}
		*t = Timestamp(i)
	case time.Time:
		*t = Timestamp(v.UnixMilli())
	default:
		return fmt.Errorf("unsupported type for Timestamp: %T", value)
	}

	return nil
}

// Int64 converts Timestamp to int64.
func (t Timestamp) Int64() int64 {
	return int64(t)
}

// String returns the string representation of Timestamp.
func (t Timestamp) String() string {
	return strconv.FormatInt(int64(t), 10)
}

// Time converts Timestamp to time.Time (assumes milliseconds).
func (t Timestamp) Time() time.Time {
	return time.UnixMilli(int64(t))
}

// Now returns the current time as a Timestamp (milliseconds).
func Now() Timestamp {
	return Timestamp(time.Now().UnixMilli())
}

// FromTime creates a Timestamp from time.Time (milliseconds).
func FromTime(t time.Time) Timestamp {
	return Timestamp(t.UnixMilli())
}
