package numx

import (
	"encoding/json"
	"testing"
)

func TestID_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		id       ID
		expected string
	}{
		{"zero", 0, `"0"`},
		{"positive", 123, `"123"`},
		{"negative", -456, `"-456"`},
		{"large number", 9007199254740992, `"9007199254740992"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.id)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(data))
			}
		})
	}
}

func TestID_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected ID
		wantErr  bool
	}{
		{"string format", `"123"`, 123, false},
		{"number format", `456`, 456, false},
		{"negative string", `"-789"`, -789, false},
		{"negative number", `-999`, -999, false},
		{"zero string", `"0"`, 0, false},
		{"zero number", `0`, 0, false},
		{"large number string", `"9007199254740992"`, 9007199254740992, false},
		{"invalid string", `"abc"`, 0, true},
		{"invalid type", `true`, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id ID
			err := json.Unmarshal([]byte(tt.json), &id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && id != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, id)
			}
		})
	}
}

func TestID_Methods(t *testing.T) {
	id := ID(123)

	if id.Int64() != 123 {
		t.Errorf("Int64() expected 123, got %d", id.Int64())
	}

	if id.String() != "123" {
		t.Errorf("String() expected '123', got '%s'", id.String())
	}

	if id.IsZero() {
		t.Error("IsZero() should return false for non-zero ID")
	}

	if !ID(0).IsZero() {
		t.Error("IsZero() should return true for zero ID")
	}
}
