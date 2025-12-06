package validation_test

import (
	"testing"

	"github.com/gocrud/csgo/errors"
	"github.com/gocrud/csgo/validation"
)

// TestOrder 测试订单结构体
type TestOrder struct {
	OrderNo string
	Items   []string
	Amounts []float64
	Tags    []int
}

// TestNotEmptySlice 测试集合非空规则
func TestNotEmptySlice(t *testing.T) {
	v := validation.NewValidator[TestOrder]()

	validation.NotEmptySlice(validation.FieldSlice(v, func(o *TestOrder) []string { return o.Items }))

	tests := []struct {
		name      string
		order     *TestOrder
		wantValid bool
	}{
		{
			name:      "有元素",
			order:     &TestOrder{Items: []string{"item1", "item2"}},
			wantValid: true,
		},
		{
			name:      "空切片",
			order:     &TestOrder{Items: []string{}},
			wantValid: false,
		},
		{
			name:      "nil切片",
			order:     &TestOrder{Items: nil},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.order)
			if result.IsValid != tt.wantValid {
				t.Errorf("validation.NotEmptySlice() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			// 检查错误码
			if !tt.wantValid && len(result.Errors) > 0 {
				if result.Errors[0].Code != errors.ValidationNotEmpty {
					t.Errorf("Error code = %v, want %v", result.Errors[0].Code, errors.ValidationNotEmpty)
				}
			}
		})
	}
}

// TestMinLengthSlice 测试集合最小长度
func TestMinLengthSlice(t *testing.T) {
	v := validation.NewValidator[TestOrder]()

	validation.MinLengthSlice(validation.FieldSlice(v, func(o *TestOrder) []string { return o.Items }), 2)

	tests := []struct {
		name      string
		order     *TestOrder
		wantValid bool
	}{
		{
			name:      "满足最小长度",
			order:     &TestOrder{Items: []string{"item1", "item2"}},
			wantValid: true,
		},
		{
			name:      "超过最小长度",
			order:     &TestOrder{Items: []string{"item1", "item2", "item3"}},
			wantValid: true,
		},
		{
			name:      "不足最小长度",
			order:     &TestOrder{Items: []string{"item1"}},
			wantValid: false,
		},
		{
			name:      "空切片",
			order:     &TestOrder{Items: []string{}},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.order)
			if result.IsValid != tt.wantValid {
				t.Errorf("validation.MinLengthSlice() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			// 检查错误码
			if !tt.wantValid && len(result.Errors) > 0 {
				if result.Errors[0].Code != errors.ValidationMinCount {
					t.Errorf("Error code = %v, want %v", result.Errors[0].Code, errors.ValidationMinCount)
				}
			}
		})
	}
}

// TestMaxLengthSlice 测试集合最大长度
func TestMaxLengthSlice(t *testing.T) {
	v := validation.NewValidator[TestOrder]()

	validation.MaxLengthSlice(validation.FieldSlice(v, func(o *TestOrder) []string { return o.Items }), 3)

	tests := []struct {
		name      string
		order     *TestOrder
		wantValid bool
	}{
		{
			name:      "不超过最大长度",
			order:     &TestOrder{Items: []string{"item1", "item2"}},
			wantValid: true,
		},
		{
			name:      "刚好最大长度",
			order:     &TestOrder{Items: []string{"item1", "item2", "item3"}},
			wantValid: true,
		},
		{
			name:      "超过最大长度",
			order:     &TestOrder{Items: []string{"item1", "item2", "item3", "item4"}},
			wantValid: false,
		},
		{
			name:      "空切片",
			order:     &TestOrder{Items: []string{}},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.order)
			if result.IsValid != tt.wantValid {
				t.Errorf("validation.MaxLengthSlice() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			// 检查错误码
			if !tt.wantValid && len(result.Errors) > 0 {
				if result.Errors[0].Code != errors.ValidationMaxCount {
					t.Errorf("Error code = %v, want %v", result.Errors[0].Code, errors.ValidationMaxCount)
				}
			}
		})
	}
}

// TestMustSlice 测试集合自定义验证
func TestMustSlice(t *testing.T) {
	v := validation.NewValidator[TestOrder]()

	// 自定义规则：所有金额必须大于0
	validation.MustSlice(validation.FieldSlice(v, func(o *TestOrder) []float64 { return o.Amounts }), func(o *TestOrder, amounts []float64) bool {
		for _, amount := range amounts {
			if amount <= 0 {
				return false
			}
		}
		return true
	})

	tests := []struct {
		name      string
		order     *TestOrder
		wantValid bool
	}{
		{
			name:      "所有金额有效",
			order:     &TestOrder{Amounts: []float64{10.5, 20.0, 30.5}},
			wantValid: true,
		},
		{
			name:      "包含零",
			order:     &TestOrder{Amounts: []float64{10.5, 0, 30.5}},
			wantValid: false,
		},
		{
			name:      "包含负数",
			order:     &TestOrder{Amounts: []float64{10.5, -5.0, 30.5}},
			wantValid: false,
		},
		{
			name:      "空切片",
			order:     &TestOrder{Amounts: []float64{}},
			wantValid: true, // 空切片满足条件
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.order)
			if result.IsValid != tt.wantValid {
				t.Errorf("validation.MustSlice() IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
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

// TestSliceRules_ChainedValidation 测试链式集合验证
func TestSliceRules_ChainedValidation(t *testing.T) {
	v := validation.NewValidator[TestOrder]()

	// Items 必须非空且长度在2-5之间
	validation.MinLengthSlice(
		validation.MaxLengthSlice(
			validation.NotEmptySlice(validation.FieldSlice(v, func(o *TestOrder) []string { return o.Items })),
			5,
		),
		2,
	)

	tests := []struct {
		name      string
		order     *TestOrder
		wantValid bool
	}{
		{
			name:      "有效长度",
			order:     &TestOrder{Items: []string{"item1", "item2", "item3"}},
			wantValid: true,
		},
		{
			name:      "空切片",
			order:     &TestOrder{Items: []string{}},
			wantValid: false,
		},
		{
			name:      "长度不足",
			order:     &TestOrder{Items: []string{"item1"}},
			wantValid: false,
		},
		{
			name:      "长度过长",
			order:     &TestOrder{Items: []string{"item1", "item2", "item3", "item4", "item5", "item6"}},
			wantValid: false,
		},
		{
			name:      "边界值2",
			order:     &TestOrder{Items: []string{"item1", "item2"}},
			wantValid: true,
		},
		{
			name:      "边界值5",
			order:     &TestOrder{Items: []string{"item1", "item2", "item3", "item4", "item5"}},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.order)
			if result.IsValid != tt.wantValid {
				t.Errorf("Chained validation IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}
		})
	}
}

// TestSliceRules_IntSlice 测试int切片验证
func TestSliceRules_IntSlice(t *testing.T) {
	v := validation.NewValidator[TestOrder]()

	validation.NotEmptySlice(validation.FieldSlice(v, func(o *TestOrder) []int { return o.Tags }))

	tests := []struct {
		name      string
		order     *TestOrder
		wantValid bool
	}{
		{
			name:      "有元素",
			order:     &TestOrder{Tags: []int{1, 2, 3}},
			wantValid: true,
		},
		{
			name:      "空切片",
			order:     &TestOrder{Tags: []int{}},
			wantValid: false,
		},
		{
			name:      "nil切片",
			order:     &TestOrder{Tags: nil},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.order)
			if result.IsValid != tt.wantValid {
				t.Errorf("Int slice validation IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}
		})
	}
}

// TestSliceRules_Float64Slice 测试float64切片验证
func TestSliceRules_Float64Slice(t *testing.T) {
	v := validation.NewValidator[TestOrder]()

	validation.MinLengthSlice(validation.FieldSlice(v, func(o *TestOrder) []float64 { return o.Amounts }), 1)

	tests := []struct {
		name      string
		order     *TestOrder
		wantValid bool
	}{
		{
			name:      "满足最小长度",
			order:     &TestOrder{Amounts: []float64{10.5}},
			wantValid: true,
		},
		{
			name:      "多个元素",
			order:     &TestOrder{Amounts: []float64{10.5, 20.0, 30.5}},
			wantValid: true,
		},
		{
			name:      "空切片",
			order:     &TestOrder{Amounts: []float64{}},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.order)
			if result.IsValid != tt.wantValid {
				t.Errorf("Float64 slice validation IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}
		})
	}
}

// TestSliceRules_CustomPredicate 测试复杂的自定义断言
func TestSliceRules_CustomPredicate(t *testing.T) {
	v := validation.NewValidator[TestOrder]()

	// 自定义规则：Items中不能有重复元素
	validation.MustSlice(validation.FieldSlice(v, func(o *TestOrder) []string { return o.Items }), func(o *TestOrder, items []string) bool {
		seen := make(map[string]bool)
		for _, item := range items {
			if seen[item] {
				return false
			}
			seen[item] = true
		}
		return true
	})

	tests := []struct {
		name      string
		order     *TestOrder
		wantValid bool
	}{
		{
			name:      "无重复元素",
			order:     &TestOrder{Items: []string{"item1", "item2", "item3"}},
			wantValid: true,
		},
		{
			name:      "有重复元素",
			order:     &TestOrder{Items: []string{"item1", "item2", "item1"}},
			wantValid: false,
		},
		{
			name:      "空切片",
			order:     &TestOrder{Items: []string{}},
			wantValid: true,
		},
		{
			name:      "单个元素",
			order:     &TestOrder{Items: []string{"item1"}},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.order)
			if result.IsValid != tt.wantValid {
				t.Errorf("Custom predicate IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}
		})
	}
}

// TestSliceRules_WithCustomMessage 测试自定义错误消息
func TestSliceRules_WithCustomMessage(t *testing.T) {
	v := validation.NewValidator[TestOrder]()

	customMessage := "订单项不能为空"
	validation.NotEmptySlice(validation.FieldSlice(v, func(o *TestOrder) []string { return o.Items })).
		WithMessage(customMessage)

	order := &TestOrder{Items: []string{}}
	result := v.Validate(order)

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

// TestSliceRules_WithCustomCode 测试自定义错误码
func TestSliceRules_WithCustomCode(t *testing.T) {
	v := validation.NewValidator[TestOrder]()

	customCode := "ORDER.ITEMS_REQUIRED"
	validation.NotEmptySlice(validation.FieldSlice(v, func(o *TestOrder) []string { return o.Items })).
		WithCode(customCode)

	order := &TestOrder{Items: []string{}}
	result := v.Validate(order)

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

// TestSliceRules_ConditionalValidation 测试条件验证
func TestSliceRules_ConditionalValidation(t *testing.T) {
	v := validation.NewValidator[TestOrder]()

	// 只有当OrderNo不为空时才验证Items
	validation.NotEmptySlice(validation.FieldSlice(v, func(o *TestOrder) []string { return o.Items })).
		When(func(o *TestOrder) bool { return o.OrderNo != "" })

	tests := []struct {
		name      string
		order     *TestOrder
		wantValid bool
	}{
		{
			name:      "OrderNo不为空且Items有效",
			order:     &TestOrder{OrderNo: "ORD001", Items: []string{"item1"}},
			wantValid: true,
		},
		{
			name:      "OrderNo不为空但Items无效",
			order:     &TestOrder{OrderNo: "ORD001", Items: []string{}},
			wantValid: false,
		},
		{
			name:      "OrderNo为空（跳过验证）",
			order:     &TestOrder{OrderNo: "", Items: []string{}},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.order)
			if result.IsValid != tt.wantValid {
				t.Errorf("Conditional validation IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}
		})
	}
}

// TestSliceRules_MultipleFieldsValidation 测试多个字段的集合验证
func TestSliceRules_MultipleFieldsValidation(t *testing.T) {
	v := validation.NewValidatorAll[TestOrder]() // 使用全量验证模式

	validation.NotEmptySlice(validation.FieldSlice(v, func(o *TestOrder) []string { return o.Items }))
	validation.MinLengthSlice(validation.FieldSlice(v, func(o *TestOrder) []float64 { return o.Amounts }), 1)

	tests := []struct {
		name           string
		order          *TestOrder
		wantValid      bool
		wantErrorCount int
	}{
		{
			name: "所有字段有效",
			order: &TestOrder{
				Items:   []string{"item1"},
				Amounts: []float64{10.5},
			},
			wantValid:      true,
			wantErrorCount: 0,
		},
		{
			name: "Items无效",
			order: &TestOrder{
				Items:   []string{},
				Amounts: []float64{10.5},
			},
			wantValid:      false,
			wantErrorCount: 1,
		},
		{
			name: "Amounts无效",
			order: &TestOrder{
				Items:   []string{"item1"},
				Amounts: []float64{},
			},
			wantValid:      false,
			wantErrorCount: 1,
		},
		{
			name: "所有字段无效",
			order: &TestOrder{
				Items:   []string{},
				Amounts: []float64{},
			},
			wantValid:      false,
			wantErrorCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := v.Validate(tt.order)
			if result.IsValid != tt.wantValid {
				t.Errorf("Multiple fields validation IsValid = %v, want %v. Errors: %v", result.IsValid, tt.wantValid, result.Errors)
			}

			if len(result.Errors) != tt.wantErrorCount {
				t.Errorf("Error count = %v, want %v. Errors: %v", len(result.Errors), tt.wantErrorCount, result.Errors)
			}
		})
	}
}
