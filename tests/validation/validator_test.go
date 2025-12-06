package validation_test

import (
	"testing"

	"github.com/gocrud/csgo/validation"
)

// TestUser 测试用户结构体
type TestUser struct {
	Name  string
	Email string
	Age   int
	Score float64
	Tags  []string
}

// TestUserValidator 用户验证器
type TestUserValidator struct {
	*validation.AbstractValidator[TestUser]
}

func NewTestUserValidator() *TestUserValidator {
	v := validation.NewValidator[TestUser]()

	validator := &TestUserValidator{
		AbstractValidator: v,
	}

	// 定义验证规则
	v.Field(func(u *TestUser) string { return u.Name }).
		NotEmpty().
		MinLength(2).
		MaxLength(50)

	v.Field(func(u *TestUser) string { return u.Email }).
		NotEmpty().
		EmailAddress()

	validation.GreaterThanOrEqual(v.FieldInt(func(u *TestUser) int { return u.Age }), 0)

	return validator
}

// TestNewValidator 测试创建验证器
func TestNewValidator(t *testing.T) {
	v := validation.NewValidator[TestUser]()

	if v == nil {
		t.Error("validation.NewValidator() returned nil")
		return
	}

	// mode 字段是私有的，无法从外部访问，只测试创建成功
}

// TestNewValidatorAll 测试创建全量验证器
func TestNewValidatorAll(t *testing.T) {
	v := validation.NewValidatorAll[TestUser]()

	if v == nil {
		t.Error("validation.NewValidatorAll() returned nil")
		return
	}

	// mode 字段是私有的，无法从外部访问，只测试创建成功
}

// TestAbstractValidator_Validate_Success 测试验证成功
func TestAbstractValidator_Validate_Success(t *testing.T) {
	validator := NewTestUserValidator()

	user := &TestUser{
		Name:  "张三",
		Email: "zhangsan@example.com",
		Age:   25,
	}

	result := validator.Validate(user)

	if !result.IsValid {
		t.Errorf("Validate() IsValid = false, want true. Errors: %v", result.Errors)
	}

	if len(result.Errors) != 0 {
		t.Errorf("Validate() Errors length = %v, want 0", len(result.Errors))
	}
}

// TestAbstractValidator_Validate_FailFast 测试快速失败模式
func TestAbstractValidator_Validate_FailFast(t *testing.T) {
	validator := NewTestUserValidator()

	user := &TestUser{
		Name:  "",        // 错误1: 不能为空
		Email: "invalid", // 错误2: 邮箱格式不正确
		Age:   -1,        // 错误3: 年龄必须大于等于0
	}

	result := validator.Validate(user)

	if result.IsValid {
		t.Error("Validate() IsValid = true, want false")
	}

	// 快速失败模式应该只返回第一个错误
	if len(result.Errors) != 1 {
		t.Errorf("Validate() in FailFast mode Errors length = %v, want 1", len(result.Errors))
	}
}

// TestAbstractValidator_Validate_AllErrors 测试全量验证模式
func TestAbstractValidator_Validate_AllErrors(t *testing.T) {
	v := validation.NewValidatorAll[TestUser]()

	v.Field(func(u *TestUser) string { return u.Name }).
		NotEmpty().
		MinLength(2)

	v.Field(func(u *TestUser) string { return u.Email }).
		NotEmpty().
		EmailAddress()

	validation.GreaterThanOrEqual(v.FieldInt(func(u *TestUser) int { return u.Age }), 0)

	user := &TestUser{
		Name:  "",        // 错误1: 不能为空
		Email: "invalid", // 错误2: 邮箱格式不正确
		Age:   -1,        // 错误3: 年龄必须大于等于0
	}

	result := v.Validate(user)

	if result.IsValid {
		t.Error("Validate() IsValid = true, want false")
	}

	// 全量验证模式应该返回所有错误
	if len(result.Errors) < 2 {
		t.Errorf("Validate() in ValidateAll mode Errors length = %v, want at least 2", len(result.Errors))
	}
}

// TestAbstractValidator_Field 测试字段验证
func TestAbstractValidator_Field(t *testing.T) {
	v := validation.NewValidator[TestUser]()

	builder := v.Field(func(u *TestUser) string { return u.Name })

	if builder == nil {
		t.Error("Field() returned nil")
		return
	}

	// fieldName 是私有字段，无法从外部访问
}

// TestAbstractValidator_FieldInt 测试整数字段验证
func TestAbstractValidator_FieldInt(t *testing.T) {
	v := validation.NewValidator[TestUser]()

	builder := v.FieldInt(func(u *TestUser) int { return u.Age })

	if builder == nil {
		t.Error("FieldInt() returned nil")
		return
	}

	// fieldName 是私有字段，无法从外部访问
}

// TestAbstractValidator_FieldFloat64 测试浮点数字段验证
func TestAbstractValidator_FieldFloat64(t *testing.T) {
	v := validation.NewValidator[TestUser]()

	builder := v.FieldFloat64(func(u *TestUser) float64 { return u.Score })

	if builder == nil {
		t.Error("FieldFloat64() returned nil")
		return
	}

	// fieldName 是私有字段，无法从外部访问
}

// TestAbstractValidator_When 测试条件验证
func TestAbstractValidator_When(t *testing.T) {
	v := validation.NewValidator[TestUser]()

	// 只有当 Age > 18 时才验证 Email
	v.Field(func(u *TestUser) string { return u.Email }).
		NotEmpty().
		When(func(u *TestUser) bool { return u.Age > 18 })

	// 测试条件满足时
	user1 := &TestUser{
		Name:  "张三",
		Email: "", // 应该验证失败
		Age:   20,
	}

	result1 := v.Validate(user1)
	if result1.IsValid {
		t.Error("When condition met: Validate() IsValid = true, want false")
	}

	// 测试条件不满足时
	user2 := &TestUser{
		Name:  "李四",
		Email: "", // 应该跳过验证
		Age:   15,
	}

	result2 := v.Validate(user2)
	if !result2.IsValid {
		t.Errorf("When condition not met: Validate() IsValid = false, want true. Errors: %v", result2.Errors)
	}
}

// TestAbstractValidator_Unless 测试反向条件验证
func TestAbstractValidator_Unless(t *testing.T) {
	v := validation.NewValidator[TestUser]()

	// 除非 Age < 18，否则验证 Email
	v.Field(func(u *TestUser) string { return u.Email }).
		NotEmpty().
		Unless(func(u *TestUser) bool { return u.Age < 18 })

	// 测试条件不满足时（应该验证）
	user1 := &TestUser{
		Name:  "张三",
		Email: "", // 应该验证失败
		Age:   20,
	}

	result1 := v.Validate(user1)
	if result1.IsValid {
		t.Error("Unless condition not met: Validate() IsValid = true, want false")
	}

	// 测试条件满足时（应该跳过验证）
	user2 := &TestUser{
		Name:  "李四",
		Email: "", // 应该跳过验证
		Age:   15,
	}

	result2 := v.Validate(user2)
	if !result2.IsValid {
		t.Errorf("Unless condition met: Validate() IsValid = false, want true. Errors: %v", result2.Errors)
	}
}

// TestAbstractValidator_MultipleRulesOnSameField 测试同一字段的多个规则
func TestAbstractValidator_MultipleRulesOnSameField(t *testing.T) {
	v := validation.NewValidator[TestUser]()

	v.Field(func(u *TestUser) string { return u.Name }).
		NotEmpty().
		MinLength(2).
		MaxLength(50)

	// 测试违反第一个规则
	user1 := &TestUser{Name: ""}
	result1 := v.Validate(user1)
	if result1.IsValid {
		t.Error("Empty name should fail validation")
	}

	// 测试违反第二个规则（快速失败模式下应该先检测到为空）
	user2 := &TestUser{Name: "a"}
	result2 := v.Validate(user2)
	if result2.IsValid {
		t.Error("Name too short should fail validation")
	}

	// 测试通过所有规则
	user3 := &TestUser{Name: "张三"}
	result3 := v.Validate(user3)
	if !result3.IsValid {
		t.Errorf("Valid name should pass validation. Errors: %v", result3.Errors)
	}
}

// TestValidatorRegistry 测试验证器注册表
func TestValidatorRegistry(t *testing.T) {
	// 创建并注册验证器
	validator := NewTestUserValidator()
	validation.RegisterValidator[TestUser](validator)

	// 获取验证器
	retrieved, ok := validation.GetValidator[TestUser]()
	if !ok {
		t.Error("validation.GetValidator() ok = false, want true")
		return
	}

	if retrieved == nil {
		t.Error("validation.GetValidator() returned nil")
		return
	}

	// 测试使用注册的验证器
	user := &TestUser{
		Name:  "张三",
		Email: "zhangsan@example.com",
		Age:   25,
	}

	result := retrieved.Validate(user)
	if !result.IsValid {
		t.Errorf("Registered validator failed. Errors: %v", result.Errors)
	}
}

// TestValidateStruct 测试验证结构体
func TestValidateStruct(t *testing.T) {
	// 注册验证器
	validator := NewTestUserValidator()
	validation.RegisterValidator[TestUser](validator)

	// 测试验证成功
	user1 := &TestUser{
		Name:  "张三",
		Email: "zhangsan@example.com",
		Age:   25,
	}

	result1 := validation.ValidateStruct(user1)
	if !result1.IsValid {
		t.Errorf("validation.ValidateStruct() failed for valid user. Errors: %v", result1.Errors)
	}

	// 测试验证失败
	user2 := &TestUser{
		Name:  "",
		Email: "invalid",
		Age:   -1,
	}

	result2 := validation.ValidateStruct(user2)
	if result2.IsValid {
		t.Error("validation.ValidateStruct() succeeded for invalid user, want failure")
	}
}

// TestWithMessage 测试自定义错误消息
func TestWithMessage(t *testing.T) {
	v := validation.NewValidator[TestUser]()

	customMessage := "用户名必须填写"
	v.Field(func(u *TestUser) string { return u.Name }).
		NotEmpty().
		WithMessage(customMessage)

	user := &TestUser{Name: ""}
	result := v.Validate(user)

	if result.IsValid {
		t.Error("Validation should fail")
		return
	}

	if len(result.Errors) == 0 {
		t.Error("Should have validation errors")
		return
	}

	if result.Errors[0].Message != customMessage {
		t.Errorf("Error message = %v, want %v", result.Errors[0].Message, customMessage)
	}
}

// TestWithCode 测试自定义错误码
func TestWithCode(t *testing.T) {
	v := validation.NewValidator[TestUser]()

	customCode := "CUSTOM.REQUIRED"
	v.Field(func(u *TestUser) string { return u.Name }).
		NotEmpty().
		WithCode(customCode)

	user := &TestUser{Name: ""}
	result := v.Validate(user)

	if result.IsValid {
		t.Error("Validation should fail")
		return
	}

	if len(result.Errors) == 0 {
		t.Error("Should have validation errors")
		return
	}

	if result.Errors[0].Code != customCode {
		t.Errorf("Error code = %v, want %v", result.Errors[0].Code, customCode)
	}
}

// TestWithMessageAndCode 测试同时自定义消息和错误码
func TestWithMessageAndCode(t *testing.T) {
	v := validation.NewValidator[TestUser]()

	customMessage := "用户名必须填写"
	customCode := "CUSTOM.REQUIRED"

	v.Field(func(u *TestUser) string { return u.Name }).
		NotEmpty().
		WithMessage(customMessage).
		WithCode(customCode)

	user := &TestUser{Name: ""}
	result := v.Validate(user)

	if result.IsValid {
		t.Error("Validation should fail")
		return
	}

	if len(result.Errors) == 0 {
		t.Error("Should have validation errors")
		return
	}

	err := result.Errors[0]
	if err.Message != customMessage {
		t.Errorf("Error message = %v, want %v", err.Message, customMessage)
	}

	if err.Code != customCode {
		t.Errorf("Error code = %v, want %v", err.Code, customCode)
	}
}

// TestValidationErrorCodes 测试验证错误码
func TestValidationErrorCodes(t *testing.T) {
	v := validation.NewValidator[TestUser]()

	v.Field(func(u *TestUser) string { return u.Name }).NotEmpty()
	v.Field(func(u *TestUser) string { return u.Email }).EmailAddress()

	user := &TestUser{
		Name:  "",
		Email: "invalid",
	}

	result := v.Validate(user)

	if result.IsValid {
		t.Error("Validation should fail")
		return
	}

	// 检查错误码是否正确设置
	if len(result.Errors) > 0 {
		err := result.Errors[0]
		if err.Code == "" {
			t.Error("Error code should not be empty")
		}

		// 验证错误码格式（应该以 VALIDATION. 开头）
		if len(err.Code) > 0 && err.Code[:11] != "VALIDATION." {
			t.Errorf("Error code format incorrect: %v", err.Code)
		}
	}
}
