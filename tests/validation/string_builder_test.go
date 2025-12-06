package validation_test

import (
	"testing"

	"github.com/gocrud/csgo/errors"
	"github.com/gocrud/csgo/validation"
)

// TestStringUser 测试用户结构体
type TestStringUser struct {
	Name     string
	Email    string
	Password string
	Phone    string
	Website  string
}

// TestStringRuleBuilder_NotEmpty 测试非空规则
func TestStringRuleBuilder_NotEmpty(t *testing.T) {
	v := validation.NewValidator[TestStringUser]()

	v.Field(func(u *TestStringUser) string { return u.Name }).NotEmpty()

	tests := []struct {
		name      string
		user      *TestStringUser
		wantValid bool
	}{
		{
			name:      "正常字符串",
			user:      &TestStringUser{Name: "张三"},
			wantValid: true,
		},
		{
			name:      "空字符串",
			user:      &TestStringUser{Name: ""},
			wantValid: false,
		},
		{
			name:      "只有空格",
			user:      &TestStringUser{Name: "   "},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.user)
			if result.IsValid != tt.wantValid {
				t.Errorf("NotEmpty() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			// 检查错误码
			if !tt.wantValid && len(result.Errors) > 0 {
				if result.Errors[0].Code != errors.ValidationRequired {
					t.Errorf("Error code = %v, want %v", result.Errors[0].Code, errors.ValidationRequired)
				}
			}
		})
	}
}

// TestStringRuleBuilder_MinLength 测试最小长度
func TestStringRuleBuilder_MinLength(t *testing.T) {
	v := validation.NewValidator[TestStringUser]()

	v.Field(func(u *TestStringUser) string { return u.Name }).MinLength(2)

	tests := []struct {
		name      string
		user      *TestStringUser
		wantValid bool
	}{
		{
			name:      "长度足够",
			user:      &TestStringUser{Name: "张三"},
			wantValid: true,
		},
		{
			name:      "长度不足",
			user:      &TestStringUser{Name: "a"},
			wantValid: false,
		},
		{
			name:      "空字符串",
			user:      &TestStringUser{Name: ""},
			wantValid: false,
		},
		{
			name:      "中文字符",
			user:      &TestStringUser{Name: "李四"},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.user)
			if result.IsValid != tt.wantValid {
				t.Errorf("MinLength() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			// 检查错误码
			if !tt.wantValid && len(result.Errors) > 0 {
				if result.Errors[0].Code != errors.ValidationMinLength {
					t.Errorf("Error code = %v, want %v", result.Errors[0].Code, errors.ValidationMinLength)
				}
			}
		})
	}
}

// TestStringRuleBuilder_MaxLength 测试最大长度
func TestStringRuleBuilder_MaxLength(t *testing.T) {
	v := validation.NewValidator[TestStringUser]()

	v.Field(func(u *TestStringUser) string { return u.Name }).MaxLength(10)

	tests := []struct {
		name      string
		user      *TestStringUser
		wantValid bool
	}{
		{
			name:      "长度合适",
			user:      &TestStringUser{Name: "张三"},
			wantValid: true,
		},
		{
			name:      "长度超出",
			user:      &TestStringUser{Name: "这是一个非常非常长的名字超过了限制"},
			wantValid: false,
		},
		{
			name:      "空字符串",
			user:      &TestStringUser{Name: ""},
			wantValid: true,
		},
		{
			name:      "刚好最大长度",
			user:      &TestStringUser{Name: "1234567890"},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.user)
			if result.IsValid != tt.wantValid {
				t.Errorf("MaxLength() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			// 检查错误码
			if !tt.wantValid && len(result.Errors) > 0 {
				if result.Errors[0].Code != errors.ValidationMaxLength {
					t.Errorf("Error code = %v, want %v", result.Errors[0].Code, errors.ValidationMaxLength)
				}
			}
		})
	}
}

// TestStringRuleBuilder_Length 测试长度范围
func TestStringRuleBuilder_Length(t *testing.T) {
	v := validation.NewValidator[TestStringUser]()

	v.Field(func(u *TestStringUser) string { return u.Name }).Length(2, 10)

	tests := []struct {
		name      string
		user      *TestStringUser
		wantValid bool
	}{
		{
			name:      "长度在范围内",
			user:      &TestStringUser{Name: "张三"},
			wantValid: true,
		},
		{
			name:      "长度太短",
			user:      &TestStringUser{Name: "a"},
			wantValid: false,
		},
		{
			name:      "长度太长",
			user:      &TestStringUser{Name: "这是一个非常非常长的名字"},
			wantValid: false,
		},
		{
			name:      "最小长度",
			user:      &TestStringUser{Name: "ab"},
			wantValid: true,
		},
		{
			name:      "最大长度",
			user:      &TestStringUser{Name: "1234567890"},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.user)
			if result.IsValid != tt.wantValid {
				t.Errorf("Length() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			// 检查错误码
			if !tt.wantValid && len(result.Errors) > 0 {
				if result.Errors[0].Code != errors.ValidationLength {
					t.Errorf("Error code = %v, want %v", result.Errors[0].Code, errors.ValidationLength)
				}
			}
		})
	}
}

// TestStringRuleBuilder_EmailAddress 测试邮箱验证
func TestStringRuleBuilder_EmailAddress(t *testing.T) {
	v := validation.NewValidator[TestStringUser]()

	v.Field(func(u *TestStringUser) string { return u.Email }).EmailAddress()

	tests := []struct {
		name      string
		user      *TestStringUser
		wantValid bool
	}{
		{
			name:      "有效邮箱",
			user:      &TestStringUser{Email: "test@example.com"},
			wantValid: true,
		},
		{
			name:      "有效邮箱带数字",
			user:      &TestStringUser{Email: "user123@test.com"},
			wantValid: true,
		},
		{
			name:      "有效邮箱带点",
			user:      &TestStringUser{Email: "user.name@example.com"},
			wantValid: true,
		},
		{
			name:      "无效邮箱缺少@",
			user:      &TestStringUser{Email: "invalid.email.com"},
			wantValid: false,
		},
		{
			name:      "无效邮箱缺少域名",
			user:      &TestStringUser{Email: "test@"},
			wantValid: false,
		},
		{
			name:      "无效邮箱格式",
			user:      &TestStringUser{Email: "@example.com"},
			wantValid: false,
		},
		{
			name:      "空字符串（允许）",
			user:      &TestStringUser{Email: ""},
			wantValid: true, // EmailAddress 允许空字符串
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.user)
			if result.IsValid != tt.wantValid {
				t.Errorf("EmailAddress() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			// 检查错误码
			if !tt.wantValid && len(result.Errors) > 0 {
				if result.Errors[0].Code != errors.ValidationEmail {
					t.Errorf("Error code = %v, want %v", result.Errors[0].Code, errors.ValidationEmail)
				}
			}
		})
	}
}

// TestStringRuleBuilder_Matches 测试正则匹配
func TestStringRuleBuilder_Matches(t *testing.T) {
	v := validation.NewValidator[TestStringUser]()

	// 验证手机号（简化版：11位数字）
	v.Field(func(u *TestStringUser) string { return u.Phone }).Matches(`^1[3-9]\d{9}$`)

	tests := []struct {
		name      string
		user      *TestStringUser
		wantValid bool
	}{
		{
			name:      "有效手机号",
			user:      &TestStringUser{Phone: "13800138000"},
			wantValid: true,
		},
		{
			name:      "无效手机号（长度不足）",
			user:      &TestStringUser{Phone: "138001380"},
			wantValid: false,
		},
		{
			name:      "无效手机号（首位错误）",
			user:      &TestStringUser{Phone: "23800138000"},
			wantValid: false,
		},
		{
			name:      "空字符串（允许）",
			user:      &TestStringUser{Phone: ""},
			wantValid: true, // Matches 允许空字符串
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.user)
			if result.IsValid != tt.wantValid {
				t.Errorf("Matches() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			// 检查错误码
			if !tt.wantValid && len(result.Errors) > 0 {
				if result.Errors[0].Code != errors.ValidationPattern {
					t.Errorf("Error code = %v, want %v", result.Errors[0].Code, errors.ValidationPattern)
				}
			}
		})
	}
}

// TestStringRuleBuilder_MustString 测试自定义验证
func TestStringRuleBuilder_MustString(t *testing.T) {
	v := validation.NewValidator[TestStringUser]()

	// 自定义规则：密码必须包含数字
	v.Field(func(u *TestStringUser) string { return u.Password }).
		MustString(func(u *TestStringUser, s string) bool {
			for _, ch := range s {
				if ch >= '0' && ch <= '9' {
					return true
				}
			}
			return false
		})

	tests := []struct {
		name      string
		user      *TestStringUser
		wantValid bool
	}{
		{
			name:      "包含数字",
			user:      &TestStringUser{Password: "abc123"},
			wantValid: true,
		},
		{
			name:      "不包含数字",
			user:      &TestStringUser{Password: "abcdef"},
			wantValid: false,
		},
		{
			name:      "只有数字",
			user:      &TestStringUser{Password: "123456"},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.user)
			if result.IsValid != tt.wantValid {
				t.Errorf("MustString() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			// 检查错误码
			if !tt.wantValid && len(result.Errors) > 0 {
				if result.Errors[0].Code != errors.ValidationFailed {
					t.Errorf("Error code = %v, want %v", result.Errors[0].Code, errors.ValidationFailed)
				}
			}
		})
	}
}

// TestStringRuleBuilder_WithMessage 测试自定义错误消息
func TestStringRuleBuilder_WithMessage(t *testing.T) {
	v := validation.NewValidator[TestStringUser]()

	customMessage := "用户名不能为空哦"
	v.Field(func(u *TestStringUser) string { return u.Name }).
		NotEmpty().
		WithMessage(customMessage)

	user := &TestStringUser{Name: ""}
	result := v.Validate(user)

	if result.IsValid {
		t.Error("Validation should fail")
		return
	}

	if len(result.Errors) == 0 {
		t.Error("Should have errors")
		return
	}

	if result.Errors[0].Message != customMessage {
		t.Errorf("Error message = %v, want %v", result.Errors[0].Message, customMessage)
	}
}

// TestStringRuleBuilder_WithCode 测试自定义错误码
func TestStringRuleBuilder_WithCode(t *testing.T) {
	v := validation.NewValidator[TestStringUser]()

	customCode := "CUSTOM.NAME_REQUIRED"
	v.Field(func(u *TestStringUser) string { return u.Name }).
		NotEmpty().
		WithCode(customCode)

	user := &TestStringUser{Name: ""}
	result := v.Validate(user)

	if result.IsValid {
		t.Error("Validation should fail")
		return
	}

	if len(result.Errors) == 0 {
		t.Error("Should have errors")
		return
	}

	if result.Errors[0].Code != customCode {
		t.Errorf("Error code = %v, want %v", result.Errors[0].Code, customCode)
	}
}

// TestStringRuleBuilder_ChainedRules 测试链式规则
func TestStringRuleBuilder_ChainedRules(t *testing.T) {
	v := validation.NewValidator[TestStringUser]()

	v.Field(func(u *TestStringUser) string { return u.Name }).
		NotEmpty().
		MinLength(2).
		MaxLength(20)

	tests := []struct {
		name      string
		user      *TestStringUser
		wantValid bool
	}{
		{
			name:      "所有规则通过",
			user:      &TestStringUser{Name: "张三"},
			wantValid: true,
		},
		{
			name:      "违反NotEmpty",
			user:      &TestStringUser{Name: ""},
			wantValid: false,
		},
		{
			name:      "违反MinLength",
			user:      &TestStringUser{Name: "a"},
			wantValid: false,
		},
		{
			name:      "违反MaxLength",
			user:      &TestStringUser{Name: "这是一个非常非常非常非常非常长的名字超过限制了吧"},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.user)
			if result.IsValid != tt.wantValid {
				t.Errorf("Chained rules IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}
		})
	}
}

// TestStringRuleBuilder_When 测试条件验证
func TestStringRuleBuilder_When(t *testing.T) {
	v := validation.NewValidator[TestStringUser]()

	// 只有当 Email 不为空时才验证格式
	v.Field(func(u *TestStringUser) string { return u.Email }).
		EmailAddress().
		When(func(u *TestStringUser) bool { return u.Email != "" })

	tests := []struct {
		name      string
		user      *TestStringUser
		wantValid bool
	}{
		{
			name:      "Email为空（跳过验证）",
			user:      &TestStringUser{Email: ""},
			wantValid: true,
		},
		{
			name:      "Email有效",
			user:      &TestStringUser{Email: "test@example.com"},
			wantValid: true,
		},
		{
			name:      "Email无效",
			user:      &TestStringUser{Email: "invalid"},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.user)
			if result.IsValid != tt.wantValid {
				t.Errorf("When() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}
		})
	}
}

// TestStringRuleBuilder_Unless 测试反向条件验证
func TestStringRuleBuilder_Unless(t *testing.T) {
	v := validation.NewValidator[TestStringUser]()

	// 除非 Email 为空，否则验证格式
	v.Field(func(u *TestStringUser) string { return u.Email }).
		EmailAddress().
		Unless(func(u *TestStringUser) bool { return u.Email == "" })

	tests := []struct {
		name      string
		user      *TestStringUser
		wantValid bool
	}{
		{
			name:      "Email为空（跳过验证）",
			user:      &TestStringUser{Email: ""},
			wantValid: true,
		},
		{
			name:      "Email有效",
			user:      &TestStringUser{Email: "test@example.com"},
			wantValid: true,
		},
		{
			name:      "Email无效",
			user:      &TestStringUser{Email: "invalid"},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.user)
			if result.IsValid != tt.wantValid {
				t.Errorf("Unless() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}
		})
	}
}
