package validation_test

import (
	"testing"

	"github.com/gocrud/csgo/errors"
	"github.com/gocrud/csgo/validation"
)

// TestProduct 测试产品结构体
type TestProduct struct {
	Name     string
	Price    float64
	Stock    int
	Discount int64
	Rating   float64
}

// TestGreaterThan_Int 测试大于规则（int类型）
func TestGreaterThan_Int(t *testing.T) {
	v := validation.NewValidator[TestProduct]()

	validation.GreaterThan(v.FieldInt(func(p *TestProduct) int { return p.Stock }), 0)

	tests := []struct {
		name      string
		product   *TestProduct
		wantValid bool
	}{
		{
			name:      "大于0",
			product:   &TestProduct{Stock: 10},
			wantValid: true,
		},
		{
			name:      "等于0",
			product:   &TestProduct{Stock: 0},
			wantValid: false,
		},
		{
			name:      "小于0",
			product:   &TestProduct{Stock: -5},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.product)
			if result.IsValid != tt.wantValid {
				t.Errorf("validation.GreaterThan() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			// 检查错误码
			if !tt.wantValid && len(result.Errors) > 0 {
				if result.Errors[0].Code != errors.ValidationMin {
					t.Errorf("Error code = %v, want %v", result.Errors[0].Code, errors.ValidationMin)
				}
			}
		})
	}
}

// TestGreaterThan_Float64 测试大于规则（float64类型）
func TestGreaterThan_Float64(t *testing.T) {
	v := validation.NewValidator[TestProduct]()

	validation.GreaterThan(v.FieldFloat64(func(p *TestProduct) float64 { return p.Price }), 0.0)

	tests := []struct {
		name      string
		product   *TestProduct
		wantValid bool
	}{
		{
			name:      "大于0",
			product:   &TestProduct{Price: 99.99},
			wantValid: true,
		},
		{
			name:      "等于0",
			product:   &TestProduct{Price: 0.0},
			wantValid: false,
		},
		{
			name:      "小于0",
			product:   &TestProduct{Price: -10.5},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.product)
			if result.IsValid != tt.wantValid {
				t.Errorf("validation.GreaterThan() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}
		})
	}
}

// TestGreaterThanOrEqual_Int 测试大于等于规则（int类型）
func TestGreaterThanOrEqual_Int(t *testing.T) {
	v := validation.NewValidator[TestProduct]()

	validation.GreaterThanOrEqual(v.FieldInt(func(p *TestProduct) int { return p.Stock }), 0)

	tests := []struct {
		name      string
		product   *TestProduct
		wantValid bool
	}{
		{
			name:      "大于0",
			product:   &TestProduct{Stock: 10},
			wantValid: true,
		},
		{
			name:      "等于0",
			product:   &TestProduct{Stock: 0},
			wantValid: true,
		},
		{
			name:      "小于0",
			product:   &TestProduct{Stock: -5},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.product)
			if result.IsValid != tt.wantValid {
				t.Errorf("validation.GreaterThanOrEqual() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			// 检查错误码
			if !tt.wantValid && len(result.Errors) > 0 {
				if result.Errors[0].Code != errors.ValidationMin {
					t.Errorf("Error code = %v, want %v", result.Errors[0].Code, errors.ValidationMin)
				}
			}
		})
	}
}

// TestGreaterThanOrEqual_Int64 测试大于等于规则（int64类型）
func TestGreaterThanOrEqual_Int64(t *testing.T) {
	v := validation.NewValidator[TestProduct]()

	validation.GreaterThanOrEqual(v.FieldInt64(func(p *TestProduct) int64 { return p.Discount }), int64(0))

	tests := []struct {
		name      string
		product   *TestProduct
		wantValid bool
	}{
		{
			name:      "大于0",
			product:   &TestProduct{Discount: 100},
			wantValid: true,
		},
		{
			name:      "等于0",
			product:   &TestProduct{Discount: 0},
			wantValid: true,
		},
		{
			name:      "小于0",
			product:   &TestProduct{Discount: -10},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.product)
			if result.IsValid != tt.wantValid {
				t.Errorf("validation.GreaterThanOrEqual() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}
		})
	}
}

// TestLessThan_Int 测试小于规则（int类型）
func TestLessThan_Int(t *testing.T) {
	v := validation.NewValidator[TestProduct]()

	validation.LessThan(v.FieldInt(func(p *TestProduct) int { return p.Stock }), 1000)

	tests := []struct {
		name      string
		product   *TestProduct
		wantValid bool
	}{
		{
			name:      "小于1000",
			product:   &TestProduct{Stock: 500},
			wantValid: true,
		},
		{
			name:      "等于1000",
			product:   &TestProduct{Stock: 1000},
			wantValid: false,
		},
		{
			name:      "大于1000",
			product:   &TestProduct{Stock: 1500},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.product)
			if result.IsValid != tt.wantValid {
				t.Errorf("validation.LessThan() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			// 检查错误码
			if !tt.wantValid && len(result.Errors) > 0 {
				if result.Errors[0].Code != errors.ValidationMax {
					t.Errorf("Error code = %v, want %v", result.Errors[0].Code, errors.ValidationMax)
				}
			}
		})
	}
}

// TestLessThan_Float64 测试小于规则（float64类型）
func TestLessThan_Float64(t *testing.T) {
	v := validation.NewValidator[TestProduct]()

	validation.LessThan(v.FieldFloat64(func(p *TestProduct) float64 { return p.Price }), 10000.0)

	tests := []struct {
		name      string
		product   *TestProduct
		wantValid bool
	}{
		{
			name:      "小于10000",
			product:   &TestProduct{Price: 5000.5},
			wantValid: true,
		},
		{
			name:      "等于10000",
			product:   &TestProduct{Price: 10000.0},
			wantValid: false,
		},
		{
			name:      "大于10000",
			product:   &TestProduct{Price: 15000.0},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.product)
			if result.IsValid != tt.wantValid {
				t.Errorf("validation.LessThan() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}
		})
	}
}

// TestLessThanOrEqual_Int 测试小于等于规则（int类型）
func TestLessThanOrEqual_Int(t *testing.T) {
	v := validation.NewValidator[TestProduct]()

	validation.LessThanOrEqual(v.FieldInt(func(p *TestProduct) int { return p.Stock }), 1000)

	tests := []struct {
		name      string
		product   *TestProduct
		wantValid bool
	}{
		{
			name:      "小于1000",
			product:   &TestProduct{Stock: 500},
			wantValid: true,
		},
		{
			name:      "等于1000",
			product:   &TestProduct{Stock: 1000},
			wantValid: true,
		},
		{
			name:      "大于1000",
			product:   &TestProduct{Stock: 1500},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.product)
			if result.IsValid != tt.wantValid {
				t.Errorf("validation.LessThanOrEqual() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			// 检查错误码
			if !tt.wantValid && len(result.Errors) > 0 {
				if result.Errors[0].Code != errors.ValidationMax {
					t.Errorf("Error code = %v, want %v", result.Errors[0].Code, errors.ValidationMax)
				}
			}
		})
	}
}

// TestInclusiveBetween_Int 测试包含边界的范围（int类型）
func TestInclusiveBetween_Int(t *testing.T) {
	v := validation.NewValidator[TestProduct]()

	validation.InclusiveBetween(v.FieldInt(func(p *TestProduct) int { return p.Stock }), 1, 100)

	tests := []struct {
		name      string
		product   *TestProduct
		wantValid bool
	}{
		{
			name:      "在范围内",
			product:   &TestProduct{Stock: 50},
			wantValid: true,
		},
		{
			name:      "等于最小值",
			product:   &TestProduct{Stock: 1},
			wantValid: true,
		},
		{
			name:      "等于最大值",
			product:   &TestProduct{Stock: 100},
			wantValid: true,
		},
		{
			name:      "小于最小值",
			product:   &TestProduct{Stock: 0},
			wantValid: false,
		},
		{
			name:      "大于最大值",
			product:   &TestProduct{Stock: 101},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.product)
			if result.IsValid != tt.wantValid {
				t.Errorf("validation.InclusiveBetween() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			// 检查错误码
			if !tt.wantValid && len(result.Errors) > 0 {
				if result.Errors[0].Code != errors.ValidationRange {
					t.Errorf("Error code = %v, want %v", result.Errors[0].Code, errors.ValidationRange)
				}
			}
		})
	}
}

// TestInclusiveBetween_Float64 测试包含边界的范围（float64类型）
func TestInclusiveBetween_Float64(t *testing.T) {
	v := validation.NewValidator[TestProduct]()

	validation.InclusiveBetween(v.FieldFloat64(func(p *TestProduct) float64 { return p.Rating }), 0.0, 5.0)

	tests := []struct {
		name      string
		product   *TestProduct
		wantValid bool
	}{
		{
			name:      "在范围内",
			product:   &TestProduct{Rating: 3.5},
			wantValid: true,
		},
		{
			name:      "等于最小值",
			product:   &TestProduct{Rating: 0.0},
			wantValid: true,
		},
		{
			name:      "等于最大值",
			product:   &TestProduct{Rating: 5.0},
			wantValid: true,
		},
		{
			name:      "小于最小值",
			product:   &TestProduct{Rating: -0.1},
			wantValid: false,
		},
		{
			name:      "大于最大值",
			product:   &TestProduct{Rating: 5.1},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.product)
			if result.IsValid != tt.wantValid {
				t.Errorf("validation.InclusiveBetween() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}
		})
	}
}

// TestExclusiveBetween_Int 测试不包含边界的范围（int类型）
func TestExclusiveBetween_Int(t *testing.T) {
	v := validation.NewValidator[TestProduct]()

	validation.ExclusiveBetween(v.FieldInt(func(p *TestProduct) int { return p.Stock }), 0, 100)

	tests := []struct {
		name      string
		product   *TestProduct
		wantValid bool
	}{
		{
			name:      "在范围内",
			product:   &TestProduct{Stock: 50},
			wantValid: true,
		},
		{
			name:      "等于最小值",
			product:   &TestProduct{Stock: 0},
			wantValid: false,
		},
		{
			name:      "等于最大值",
			product:   &TestProduct{Stock: 100},
			wantValid: false,
		},
		{
			name:      "小于最小值",
			product:   &TestProduct{Stock: -1},
			wantValid: false,
		},
		{
			name:      "大于最大值",
			product:   &TestProduct{Stock: 101},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.product)
			if result.IsValid != tt.wantValid {
				t.Errorf("validation.ExclusiveBetween() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			// 检查错误码
			if !tt.wantValid && len(result.Errors) > 0 {
				if result.Errors[0].Code != errors.ValidationRange {
					t.Errorf("Error code = %v, want %v", result.Errors[0].Code, errors.ValidationRange)
				}
			}
		})
	}
}

// TestExclusiveBetween_Float64 测试不包含边界的范围（float64类型）
func TestExclusiveBetween_Float64(t *testing.T) {
	v := validation.NewValidator[TestProduct]()

	validation.ExclusiveBetween(v.FieldFloat64(func(p *TestProduct) float64 { return p.Price }), 0.0, 10000.0)

	tests := []struct {
		name      string
		product   *TestProduct
		wantValid bool
	}{
		{
			name:      "在范围内",
			product:   &TestProduct{Price: 5000.0},
			wantValid: true,
		},
		{
			name:      "等于最小值",
			product:   &TestProduct{Price: 0.0},
			wantValid: false,
		},
		{
			name:      "等于最大值",
			product:   &TestProduct{Price: 10000.0},
			wantValid: false,
		},
		{
			name:      "小于最小值",
			product:   &TestProduct{Price: -0.1},
			wantValid: false,
		},
		{
			name:      "大于最大值",
			product:   &TestProduct{Price: 10000.1},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.product)
			if result.IsValid != tt.wantValid {
				t.Errorf("validation.ExclusiveBetween() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}
		})
	}
}

// TestMustNumber 测试自定义数字验证
func TestMustNumber(t *testing.T) {
	v := validation.NewValidator[TestProduct]()

	// 自定义规则：价格必须是整数
	validation.MustNumber(v.FieldFloat64(func(p *TestProduct) float64 { return p.Price }), func(p *TestProduct, price float64) bool {
		return price == float64(int(price))
	})

	tests := []struct {
		name      string
		product   *TestProduct
		wantValid bool
	}{
		{
			name:      "整数价格",
			product:   &TestProduct{Price: 100.0},
			wantValid: true,
		},
		{
			name:      "小数价格",
			product:   &TestProduct{Price: 99.99},
			wantValid: false,
		},
		{
			name:      "零",
			product:   &TestProduct{Price: 0.0},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.product)
			if result.IsValid != tt.wantValid {
				t.Errorf("validation.MustNumber() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
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

// TestNumberRules_ChainedValidation 测试链式数字验证
func TestNumberRules_ChainedValidation(t *testing.T) {
	v := validation.NewValidator[TestProduct]()

	// 库存必须大于等于0且小于1000
	validation.GreaterThanOrEqual(
		validation.LessThan(v.FieldInt(func(p *TestProduct) int { return p.Stock }), 1000),
		0,
	)

	tests := []struct {
		name      string
		product   *TestProduct
		wantValid bool
	}{
		{
			name:      "有效库存",
			product:   &TestProduct{Stock: 100},
			wantValid: true,
		},
		{
			name:      "负数库存",
			product:   &TestProduct{Stock: -1},
			wantValid: false,
		},
		{
			name:      "库存过大",
			product:   &TestProduct{Stock: 1000},
			wantValid: false,
		},
		{
			name:      "边界值0",
			product:   &TestProduct{Stock: 0},
			wantValid: true,
		},
		{
			name:      "边界值999",
			product:   &TestProduct{Stock: 999},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.product)
			if result.IsValid != tt.wantValid {
				t.Errorf("Chained validation IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}
		})
	}
}
