package validation_test

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gocrud/csgo/validation"
	"github.com/gocrud/csgo/validation/v"
)

// ==================== Test Models ====================

type User struct {
	Name  v.String `json:"name"`
	Age   v.Int    `json:"age"`
	Email v.String `json:"email"`
}

type ComplexTypes struct {
	Score    v.Float64
	IsActive v.Bool
	Tags     v.Slice[string]
	Counts   v.Slice[int]
	Created  v.Time
}

type NestedStruct struct {
	Title v.String
	Meta  struct {
		Author v.String `json:"author"`
		Info   struct {
			Version v.Int `json:"ver"`
		} `json:"info"`
	} `json:"meta"`
}

type PrivateStruct struct {
	public  v.String
	private v.String // 私有字段
}

// ==================== Tests ====================

func TestPrimitives(t *testing.T) {
	validation.Register(func(u *User) {
		u.Name.Required().MinLen(3).MaxLen(10).Regex("^[a-zA-Z]+$")
		u.Age.Range(18, 100)
		u.Email.Email()
	})

	tests := []struct {
		name    string
		input   User
		wantErr bool
		errMsg  string
	}{
		{"Valid", User{Name: "Alice", Age: 25, Email: "alice@example.com"}, false, ""},
		{"NameRequired", User{Age: 25}, true, "is required"},
		{"NameTooShort", User{Name: "Al", Age: 25}, true, "length must be at least 3"},
		{"NameTooLong", User{Name: "Christopher", Age: 25}, true, "length must be at most 10"},
		{"NameRegexFail", User{Name: "Alice123", Age: 25}, true, "must match regex"},
		{"AgeTooYoung", User{Name: "Alice", Age: 10}, true, "must be between 18 and 100"},
		{"AgeTooOld", User{Name: "Alice", Age: 101}, true, "must be between 18 and 100"},
		{"InvalidEmail", User{Name: "Alice", Age: 25, Email: "invalid"}, true, "invalid email format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validation.Validate(&tt.input)
			if (errs != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", errs, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(errs.Error(), tt.errMsg) {
				t.Errorf("Validate() error = %v, wantErrMsg %v", errs, tt.errMsg)
			}
		})
	}
}

func TestExtendedTypes(t *testing.T) {
	validation.Register(func(c *ComplexTypes) {
		c.Score.Min(0.0).Max(100.0)
		c.IsActive.True()
		c.Tags.MinLen(1).MaxLen(5)
		c.Counts.Required()
		c.Created.Required().Before(time.Now())
	})

	future := time.Now().Add(time.Hour)
	past := time.Now().Add(-time.Hour)

	tests := []struct {
		name    string
		input   ComplexTypes
		wantErr bool
		errMsg  string
	}{
		{"Valid", ComplexTypes{Score: 88.5, IsActive: true, Tags: []string{"a"}, Counts: []int{1}, Created: v.Time(past)}, false, ""},
		{"ScoreLow", ComplexTypes{Score: -1.0, IsActive: true, Tags: []string{"a"}, Counts: []int{1}, Created: v.Time(past)}, true, "must be at least 0.000000"},
		{"ScoreHigh", ComplexTypes{Score: 101.0, IsActive: true, Tags: []string{"a"}, Counts: []int{1}, Created: v.Time(past)}, true, "must be at most 100.000000"},
		{"Inactive", ComplexTypes{Score: 50, IsActive: false, Tags: []string{"a"}, Counts: []int{1}, Created: v.Time(past)}, true, "must be true"},
		{"TagsEmpty", ComplexTypes{Score: 50, IsActive: true, Tags: []string{}, Counts: []int{1}, Created: v.Time(past)}, true, "length must be at least 1"},
		{"CountsNil", ComplexTypes{Score: 50, IsActive: true, Tags: []string{"a"}, Counts: nil, Created: v.Time(past)}, true, "is required"},
		{"FutureTime", ComplexTypes{Score: 50, IsActive: true, Tags: []string{"a"}, Counts: []int{1}, Created: v.Time(future)}, true, "must be before"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validation.Validate(&tt.input)
			if (errs != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", errs, tt.wantErr)
			}
			if tt.wantErr && !strings.Contains(errs.Error(), tt.errMsg) {
				t.Errorf("Validate() error = %v, wantErrMsg %v", errs, tt.errMsg)
			}
		})
	}
}

func TestNestedAndPath(t *testing.T) {
	validation.Register(func(n *NestedStruct) {
		n.Title.Required()
		n.Meta.Author.Required()
		n.Meta.Info.Version.Min(1)
	})

	input := NestedStruct{}
	// Intentionally empty to trigger all errors (but fail-fast will stop at first)
	// We use this to test field names

	errs := validation.Validate(&input)
	if errs == nil {
		t.Fatal("Expected errors")
	}

	// FailFast is default, so we expect only the first error (Title) or consistent order
	// Let's check if we can trigger specific errors

	input.Title = "Valid"
	errs = validation.Validate(&input) // Should fail at Meta.Author

	firstErr := errs.FirstError()
	if firstErr.Field != "meta.author" {
		t.Errorf("Expected field 'meta.author', got '%s'", firstErr.Field)
	}

	input.Meta.Author = "Me"
	errs = validation.Validate(&input) // Should fail at Meta.Info.Version

	firstErr = errs.FirstError()
	if firstErr.Field != "meta.info.ver" { // JSON tag 'ver'
		t.Errorf("Expected field 'meta.info.ver', got '%s'", firstErr.Field)
	}
}

func TestFailFastVsAll(t *testing.T) {
	type MultiError struct {
		A v.Int
		B v.Int
	}

	// Case 1: FailFast (Default)
	validation.Register(func(m *MultiError) {
		m.A.Min(10)
		m.B.Min(10)
	})

	input := MultiError{A: 1, B: 1}
	errs := validation.Validate(&input)
	if len(errs) != 1 {
		t.Errorf("FailFast: expected 1 error, got %d", len(errs))
	}

	// Case 2: RegisterAll
	validation.RegisterAll(func(m *MultiError) {
		m.A.Min(10)
		m.B.Min(10)
	})

	// Re-validate same input (schema updated)
	errs = validation.Validate(&input)
	if len(errs) != 2 {
		t.Errorf("RegisterAll: expected 2 errors, got %d", len(errs))
	}
}

func TestCustomMessage(t *testing.T) {
	type MsgStruct struct {
		Name v.String
	}
	validation.Register(func(m *MsgStruct) {
		m.Name.Required().Msg("姓名不能为空")
	})

	input := MsgStruct{}
	errs := validation.Validate(&input)
	if errs == nil || errs.FirstError().Message != "姓名不能为空" {
		t.Errorf("Expected custom message '姓名不能为空', got '%v'", errs)
	}
}

func TestPrivateFields(t *testing.T) {
	validation.Register(func(p *PrivateStruct) {
		p.public.Required()
		p.private.Required() // 验证私有字段
	})

	input := PrivateStruct{public: "ok", private: ""}
	errs := validation.Validate(&input)

	if errs == nil {
		t.Fatal("Expected error on private field")
	}

	// Private fields use struct field name if no JSON tag
	if errs.FirstError().Field != "private" {
		t.Errorf("Expected field 'private', got '%s'", errs.FirstError().Field)
	}
}

func TestConcurrency(t *testing.T) {
	type ConcurrentUser struct {
		Id v.Int
	}
	validation.Register(func(u *ConcurrentUser) {
		u.Id.Min(1)
	})

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			u := ConcurrentUser{Id: 0}
			validation.Validate(&u)
		}()
	}
	wg.Wait()
}
