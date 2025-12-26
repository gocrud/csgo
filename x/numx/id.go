package numx

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ID is an int64 type that serializes to/from JSON as a string to avoid precision loss in JavaScript.
type ID int64

func (id ID) IsEmpty() bool {
	return id == 0
}

// MarshalJSON implements json.Marshaler interface.
// Converts int64 to JSON string to prevent precision loss in JavaScript.
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatInt(int64(id), 10))
}

// UnmarshalJSON implements json.Unmarshaler interface.
// Accepts both JSON string and number formats.
func (id *ID) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid ID string: %w", err)
		}
		*id = ID(i)
		return nil
	}

	// Try to unmarshal as number
	var num int64
	if err := json.Unmarshal(data, &num); err != nil {
		return fmt.Errorf("ID must be a string or number: %w", err)
	}
	*id = ID(num)
	return nil
}

// Int64 converts ID to int64.
func (id ID) Int64() int64 {
	return int64(id)
}

// String returns the string representation of ID.
func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

// IsZero returns true if ID is zero.
func (id ID) IsZero() bool {
	return id == 0
}
