package validation_test

import (
	"testing"

	"github.com/gocrud/csgo/validation"
)

// TestValidationError_Error 测试验证错误格式化
func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  validation.ValidationError
		want string
	}{
		{
			name: "带错误码",
			err: validation.ValidationError{
				Field:   "email",
				Message: "邮箱格式不正确",
				Code:    "VALIDATION.EMAIL",
			},
			want: "[VALIDATION.EMAIL] email: 邮箱格式不正确",
		},
		{
			name: "无错误码",
			err: validation.ValidationError{
				Field:   "name",
				Message: "不能为空",
				Code:    "",
			},
			want: "name: 不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("validation.ValidationError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestValidationErrors_Error 测试验证错误集合格式化
func TestValidationErrors_Error(t *testing.T) {
	tests := []struct {
		name   string
		errors validation.ValidationErrors
		want   string
	}{
		{
			name:   "空错误",
			errors: validation.ValidationErrors{},
			want:   "",
		},
		{
			name: "单个错误",
			errors: validation.ValidationErrors{
				{Field: "email", Message: "邮箱格式不正确", Code: "VALIDATION.EMAIL"},
			},
			want: "[VALIDATION.EMAIL] email: 邮箱格式不正确",
		},
		{
			name: "多个错误",
			errors: validation.ValidationErrors{
				{Field: "email", Message: "邮箱格式不正确", Code: "VALIDATION.EMAIL"},
				{Field: "name", Message: "不能为空", Code: "VALIDATION.REQUIRED"},
			},
			want: "[VALIDATION.EMAIL] email: 邮箱格式不正确; [VALIDATION.REQUIRED] name: 不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.errors.Error(); got != tt.want {
				t.Errorf("validation.ValidationErrors.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestValidationErrors_Add 测试添加错误
func TestValidationErrors_Add(t *testing.T) {
	errors := validation.ValidationErrors{}

	errors.Add("email", "邮箱格式不正确", "VALIDATION.EMAIL")
	errors.Add("name", "不能为空", "VALIDATION.REQUIRED")

	if len(errors) != 2 {
		t.Errorf("validation.ValidationErrors.Add() length = %v, want 2", len(errors))
	}

	if errors[0].Field != "email" {
		t.Errorf("First error field = %v, want email", errors[0].Field)
	}

	if errors[1].Field != "name" {
		t.Errorf("Second error field = %v, want name", errors[1].Field)
	}
}

// TestValidationErrors_GetByField 测试按字段获取错误
func TestValidationErrors_GetByField(t *testing.T) {
	errors := validation.ValidationErrors{
		{Field: "email", Message: "邮箱格式不正确", Code: "VALIDATION.EMAIL"},
		{Field: "name", Message: "不能为空", Code: "VALIDATION.REQUIRED"},
		{Field: "email", Message: "邮箱已存在", Code: "VALIDATION.UNIQUE"},
	}

	emailErrors := errors.GetByField("email")
	if len(emailErrors) != 2 {
		t.Errorf("GetByField('email') length = %v, want 2", len(emailErrors))
	}

	nameErrors := errors.GetByField("name")
	if len(nameErrors) != 1 {
		t.Errorf("GetByField('name') length = %v, want 1", len(nameErrors))
	}

	notExistErrors := errors.GetByField("notexist")
	if len(notExistErrors) != 0 {
		t.Errorf("GetByField('notexist') length = %v, want 0", len(notExistErrors))
	}
}

// TestValidationErrors_Fields 测试获取所有错误字段
func TestValidationErrors_Fields(t *testing.T) {
	errors := validation.ValidationErrors{
		{Field: "email", Message: "邮箱格式不正确", Code: "VALIDATION.EMAIL"},
		{Field: "name", Message: "不能为空", Code: "VALIDATION.REQUIRED"},
		{Field: "email", Message: "邮箱已存在", Code: "VALIDATION.UNIQUE"},
		{Field: "age", Message: "年龄必须大于0", Code: "VALIDATION.MIN"},
	}

	fields := errors.Fields()
	if len(fields) != 3 {
		t.Errorf("Fields() length = %v, want 3", len(fields))
	}

	// 检查字段是否包含预期值
	expectedFields := map[string]bool{"email": false, "name": false, "age": false}
	for _, field := range fields {
		if _, ok := expectedFields[field]; ok {
			expectedFields[field] = true
		}
	}

	for field, found := range expectedFields {
		if !found {
			t.Errorf("Fields() missing field %v", field)
		}
	}
}

// TestValidationErrors_FirstError 测试获取第一个错误
func TestValidationErrors_FirstError(t *testing.T) {
	errors := validation.ValidationErrors{
		{Field: "email", Message: "邮箱格式不正确", Code: "VALIDATION.EMAIL"},
		{Field: "name", Message: "不能为空", Code: "VALIDATION.REQUIRED"},
	}

	first := errors.FirstError()
	if first == nil {
		t.Error("FirstError() returned nil")
		return
	}

	if first.Field != "email" {
		t.Errorf("FirstError() Field = %v, want email", first.Field)
	}

	// 测试空错误集合
	emptyErrors := validation.ValidationErrors{}
	if emptyErrors.FirstError() != nil {
		t.Error("FirstError() should return nil for empty errors")
	}
}

// TestValidationErrors_ToFieldMap 测试转换为字段映射
func TestValidationErrors_ToFieldMap(t *testing.T) {
	errors := validation.ValidationErrors{
		{Field: "email", Message: "邮箱格式不正确", Code: "VALIDATION.EMAIL"},
		{Field: "name", Message: "不能为空", Code: "VALIDATION.REQUIRED"},
		{Field: "email", Message: "邮箱已存在", Code: "VALIDATION.UNIQUE"},
	}

	fieldMap := errors.ToFieldMap()

	if len(fieldMap) != 2 {
		t.Errorf("ToFieldMap() length = %v, want 2", len(fieldMap))
	}

	emailMessages := fieldMap["email"]
	if len(emailMessages) != 2 {
		t.Errorf("email messages length = %v, want 2", len(emailMessages))
	}

	nameMessages := fieldMap["name"]
	if len(nameMessages) != 1 {
		t.Errorf("name messages length = %v, want 1", len(nameMessages))
	}
}

// TestValidationErrors_HasErrors 测试是否有错误
func TestValidationErrors_HasErrors(t *testing.T) {
	errors := validation.ValidationErrors{
		{Field: "email", Message: "邮箱格式不正确", Code: "VALIDATION.EMAIL"},
	}

	if !errors.HasErrors() {
		t.Error("HasErrors() = false, want true")
	}

	emptyErrors := validation.ValidationErrors{}
	if emptyErrors.HasErrors() {
		t.Error("HasErrors() = true, want false for empty errors")
	}
}

// TestValidationResult_ToError 测试转换为error
func TestValidationResult_ToError(t *testing.T) {
	// 测试成功结果
	successResult := validation.ValidationResult{
		IsValid: true,
		Errors:  validation.ValidationErrors{},
	}

	if err := successResult.ToError(); err != nil {
		t.Errorf("ToError() for success result = %v, want nil", err)
	}

	// 测试失败结果
	failedResult := validation.ValidationResult{
		IsValid: false,
		Errors: validation.ValidationErrors{
			{Field: "email", Message: "邮箱格式不正确", Code: "VALIDATION.EMAIL"},
		},
	}

	if err := failedResult.ToError(); err == nil {
		t.Error("ToError() for failed result = nil, want error")
	}
}

// TestNewValidationResult 测试创建验证结果
func TestNewValidationResult(t *testing.T) {
	// 测试有错误的情况
	errors := validation.ValidationErrors{
		{Field: "email", Message: "邮箱格式不正确", Code: "VALIDATION.EMAIL"},
	}

	result := validation.NewValidationResult(errors)
	if result.IsValid {
		t.Error("validation.NewValidationResult() IsValid = true, want false")
	}

	if len(result.Errors) != 1 {
		t.Errorf("validation.NewValidationResult() Errors length = %v, want 1", len(result.Errors))
	}

	// 测试无错误的情况
	emptyErrors := validation.ValidationErrors{}
	successResult := validation.NewValidationResult(emptyErrors)
	if !successResult.IsValid {
		t.Error("validation.NewValidationResult() IsValid = false, want true")
	}
}

// TestSuccessResult 测试创建成功结果
func TestSuccessResult(t *testing.T) {
	result := validation.SuccessResult()

	if !result.IsValid {
		t.Error("validation.SuccessResult() IsValid = false, want true")
	}

	if result.Errors == nil {
		t.Error("validation.SuccessResult() Errors = nil, want empty slice")
	}

	if len(result.Errors) != 0 {
		t.Errorf("validation.SuccessResult() Errors length = %v, want 0", len(result.Errors))
	}
}
