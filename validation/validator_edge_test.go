package validation_test

import (
	"testing"
	"time"

	"github.com/gocrud/csgo/validation"
	"github.com/gocrud/csgo/validation/v"
)

// ==================== Test Models ====================

type EmbeddedInner struct {
	InnerField v.String `json:"inner"`
}

type EmbeddedOuter struct {
	EmbeddedInner          // Anonymous embedding
	OuterField    v.String `json:"outer"`
}

type SliceStruct struct {
	Ints    v.Slice[int]    `json:"ints"`
	Strings v.Slice[string] `json:"strings"`
}

type TimeStruct struct {
	At v.Time `json:"at"`
}

// ==================== Tests ====================

func TestNilInput(t *testing.T) {
	// Test passing nil to Validate
	var u *SliceStruct = nil
	errs := validation.Validate(u)
	if errs != nil {
		t.Error("Validate(nil) should return nil errors")
	}
}

func TestGenericSlices_EdgeCases(t *testing.T) {
	validation.Register(func(s *SliceStruct) {
		s.Ints.MinLen(2).MaxLen(4)
		s.Strings.Required()
	})

	tests := []struct {
		name    string
		input   SliceStruct
		wantErr bool
		errMsg  string
	}{
		{"NilSlice", SliceStruct{Ints: nil, Strings: []string{"a"}}, true, "length must be at least 2"},
		{"EmptySlice", SliceStruct{Ints: []int{}, Strings: []string{"a"}}, true, "length must be at least 2"},
		{"ExactMin", SliceStruct{Ints: []int{1, 2}, Strings: []string{"a"}}, false, ""},
		{"ExactMax", SliceStruct{Ints: []int{1, 2, 3, 4}, Strings: []string{"a"}}, false, ""},
		{"OverMax", SliceStruct{Ints: []int{1, 2, 3, 4, 5}, Strings: []string{"a"}}, true, "length must be at most 4"},
		{"StringsNil", SliceStruct{Ints: []int{1, 2}, Strings: nil}, true, "is required"},
		{"StringsEmpty", SliceStruct{Ints: []int{1, 2}, Strings: []string{}}, true, "is required"}, // Required checks length > 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validation.Validate(&tt.input)
			if (errs != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", errs, tt.wantErr)
			}
			if tt.wantErr && errs != nil && errs.FirstError().Message != tt.errMsg {
				t.Errorf("Expected error '%s', got '%s'", tt.errMsg, errs.FirstError().Message)
			}
		})
	}
}

func TestEmbeddedStructs(t *testing.T) {
	validation.Register(func(e *EmbeddedOuter) {
		e.InnerField.Required()
		e.OuterField.Required()
	})

	input := EmbeddedOuter{}
	errs := validation.Validate(&input)

	if errs == nil {
		t.Fatal("Expected errors")
	}

	// Check field names for embedded structs
	// Current behavior: Embedded fields are flattened or use their type name?
	// Let's verify what happens.
	// If json tag is present on the field inside EmbeddedInner, it should be used.
	// But since it's anonymous, the path might be just "inner" or "EmbeddedInner.inner"

	foundInner := false
	for _, e := range errs {
		// Expect "inner" because it's at the top level JSON-wise (promoted)
		// Or if our offset map logic treats it as nested, it might be different.
		// Go's JSON marshaler promotes fields.
		// Our validation logic iterates fields.
		// If mapOffsetsToNames recurses, it adds prefix.
		// For anonymous fields, field.Name is the type name.
		t.Logf("Field Error: %s -> %s", e.Field, e.Message)
		if e.Field == "inner" || e.Field == "EmbeddedInner.inner" {
			foundInner = true
		}
	}

	if !foundInner {
		t.Error("Expected error for embedded field 'inner'")
	}
}

func TestTimeValidation_EdgeCases(t *testing.T) {
	now := time.Now()
	validation.Register(func(tm *TimeStruct) {
		tm.At.After(now)
	})

	// Exact same time
	input := TimeStruct{At: v.Time(now)}
	errs := validation.Validate(&input)
	if errs == nil {
		// After means strict >
		t.Error("Expected error for exact time match (After is strict)")
	}

	// One nanosecond after
	input.At = v.Time(now.Add(1 * time.Nanosecond))
	errs = validation.Validate(&input)
	if errs != nil {
		t.Errorf("Expected valid for slightly after time, got %v", errs)
	}
}

func TestChainCustomMessages(t *testing.T) {
	type ChainMsg struct {
		Val v.Int
	}
	validation.Register(func(c *ChainMsg) {
		c.Val.Min(10).Msg("Too small").Max(20).Msg("Too big")
	})

	// Test Min fail
	input := ChainMsg{Val: 5}
	errs := validation.Validate(&input)
	if errs == nil || errs.FirstError().Message != "Too small" {
		t.Errorf("Expected 'Too small', got '%v'", errs)
	}

	// Test Max fail
	input.Val = 25
	errs = validation.Validate(&input)
	if errs == nil || errs.FirstError().Message != "Too big" {
		t.Errorf("Expected 'Too big', got '%v'", errs)
	}
}

func TestOverwritingRegistration(t *testing.T) {
	type Overwrite struct {
		Val v.Int
	}

	// First registration
	validation.Register(func(o *Overwrite) {
		o.Val.Min(10)
	})

	input := Overwrite{Val: 5}
	if errs := validation.Validate(&input); errs == nil {
		t.Error("Expected error from first registration")
	}

	// Second registration (should overwrite)
	validation.Register(func(o *Overwrite) {
		o.Val.Max(5) // Change rule completely
	})

	// Now 5 should be valid (Max 5)
	if errs := validation.Validate(&input); errs != nil {
		t.Errorf("Expected valid after overwrite, got %v", errs)
	}

	// Test new rule
	input.Val = 6
	if errs := validation.Validate(&input); errs == nil {
		t.Error("Expected error from second registration")
	}
}
