package numx_test

import (
	"encoding/json"
	"fmt"

	"github.com/gocrud/csgo/numx"
)

// Example demonstrates basic usage of numx types
func Example() {
	type User struct {
		ID        numx.ID        `json:"id"`
		Score     numx.BigInt    `json:"score"`
		Views     numx.BigUint   `json:"views"`
		CreatedAt numx.Timestamp `json:"created_at"`
	}

	user := User{
		ID:        9007199254740992, // Greater than JavaScript's MAX_SAFE_INTEGER
		Score:     -123456789012345,
		Views:     18446744073709551,
		CreatedAt: numx.Now(),
	}

	// Marshal to JSON - numbers become strings
	data, _ := json.Marshal(user)
	fmt.Printf("JSON: %s\n", data)

	// Unmarshal from JSON - accepts both string and number
	jsonStr := `{"id":"123","score":-456,"views":"789","created_at":"1703234567890"}`
	var u User
	json.Unmarshal([]byte(jsonStr), &u)
	fmt.Printf("User ID: %d, Score: %d\n", u.ID.Int64(), u.Score.Int64())
}

// ExampleID demonstrates ID type usage
func ExampleID() {
	id := numx.ID(9007199254740992)

	// JSON serialization
	data, _ := json.Marshal(id)
	fmt.Printf("JSON: %s\n", data)

	// Type conversion
	fmt.Printf("Int64: %d\n", id.Int64())
	fmt.Printf("String: %s\n", id.String())
	fmt.Printf("IsZero: %v\n", id.IsZero())

	// Output:
	// JSON: "9007199254740992"
	// Int64: 9007199254740992
	// String: 9007199254740992
	// IsZero: false
}

// ExampleTimestamp demonstrates Timestamp type usage
func ExampleTimestamp() {
	// Create timestamp
	now := numx.Now()

	// Convert to time.Time
	t := now.Time()
	fmt.Printf("Time: %v\n", t.Format("2006-01-02"))

	// JSON serialization
	data, _ := json.Marshal(now)
	fmt.Printf("JSON: %s\n", data)
}
