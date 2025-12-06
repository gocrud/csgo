package validation

// IValidator 验证器接口
type IValidator[T any] interface {
	Validate(instance *T) ValidationResult
}

// ValidationContext 验证上下文
type ValidationContext struct {
	FieldName  string
	FieldValue interface{}
	RootObject interface{}
}

// ValidateMode 验证模式
type ValidateMode int

const (
	// ValidateFailFast 快速失败模式（默认），遇到第一个错误立即返回
	ValidateFailFast ValidateMode = iota
	// ValidateAll 验证所有字段，收集所有错误
	ValidateAll
)

// AbstractValidator 抽象验证器基类
type AbstractValidator[T any] struct {
	rules []func(*T, ValidateMode) ValidationErrors
	mode  ValidateMode // 验证模式
}

// NewValidator 创建快速失败模式验证器（默认，推荐）
func NewValidator[T any]() *AbstractValidator[T] {
	return &AbstractValidator[T]{
		rules: make([]func(*T, ValidateMode) ValidationErrors, 0),
		mode:  ValidateFailFast,
	}
}

// NewValidatorAll 创建全量验证模式验证器（收集所有错误）
func NewValidatorAll[T any]() *AbstractValidator[T] {
	return &AbstractValidator[T]{
		rules: make([]func(*T, ValidateMode) ValidationErrors, 0),
		mode:  ValidateAll,
	}
}

// Validate 执行验证（使用创建时指定的模式）
func (v *AbstractValidator[T]) Validate(instance *T) ValidationResult {
	errors := ValidationErrors{}

	for _, rule := range v.rules {
		ruleErrors := rule(instance, v.mode)
		errors = append(errors, ruleErrors...)

		// 快速失败模式：遇到第一个错误立即返回
		if v.mode == ValidateFailFast && len(errors) > 0 {
			break
		}
	}

	return NewValidationResult(errors)
}

// Field 为字符串字段定义规则（自动提取字段名）
func (v *AbstractValidator[T]) Field(selector func(*T) string) *StringRuleBuilder[T] {
	fieldName := ExtractFieldNameSimple(selector)

	baseBuilder := &RuleBuilder[T, string]{
		validator: v,
		selector:  selector,
		fieldName: fieldName,
		rules:     make([]func(*T) *ValidationError, 0),
	}

	builder := &StringRuleBuilder[T]{
		RuleBuilder: baseBuilder,
	}

	// 将规则收集函数添加到验证器
	v.rules = append(v.rules, func(instance *T, mode ValidateMode) ValidationErrors {
		errors := ValidationErrors{}

		// 检查条件
		if baseBuilder.condition != nil && !baseBuilder.condition(instance) {
			return errors
		}

		// 执行所有规则
		for _, rule := range baseBuilder.rules {
			if err := rule(instance); err != nil {
				errors = append(errors, *err)
				// 快速失败模式：遇到错误立即返回
				if mode == ValidateFailFast {
					return errors
				}
			}
		}
		return errors
	})

	return builder
}

// FieldInt 为 int 字段定义规则（自动提取字段名）
func (v *AbstractValidator[T]) FieldInt(selector func(*T) int) *RuleBuilder[T, int] {
	fieldName := ExtractFieldNameSimple(selector)

	builder := &RuleBuilder[T, int]{
		validator: v,
		selector:  selector,
		fieldName: fieldName,
		rules:     make([]func(*T) *ValidationError, 0),
	}

	v.rules = append(v.rules, func(instance *T, mode ValidateMode) ValidationErrors {
		errors := ValidationErrors{}

		if builder.condition != nil && !builder.condition(instance) {
			return errors
		}

		for _, rule := range builder.rules {
			if err := rule(instance); err != nil {
				errors = append(errors, *err)
				// 快速失败模式：遇到错误立即返回
				if mode == ValidateFailFast {
					return errors
				}
			}
		}
		return errors
	})

	return builder
}

// FieldInt64 为 int64 字段定义规则（自动提取字段名）
func (v *AbstractValidator[T]) FieldInt64(selector func(*T) int64) *RuleBuilder[T, int64] {
	fieldName := ExtractFieldNameSimple(selector)

	builder := &RuleBuilder[T, int64]{
		validator: v,
		selector:  selector,
		fieldName: fieldName,
		rules:     make([]func(*T) *ValidationError, 0),
	}

	v.rules = append(v.rules, func(instance *T, mode ValidateMode) ValidationErrors {
		errors := ValidationErrors{}

		if builder.condition != nil && !builder.condition(instance) {
			return errors
		}

		for _, rule := range builder.rules {
			if err := rule(instance); err != nil {
				errors = append(errors, *err)
				// 快速失败模式：遇到错误立即返回
				if mode == ValidateFailFast {
					return errors
				}
			}
		}
		return errors
	})

	return builder
}

// FieldFloat64 为 float64 字段定义规则（自动提取字段名）
func (v *AbstractValidator[T]) FieldFloat64(selector func(*T) float64) *RuleBuilder[T, float64] {
	fieldName := ExtractFieldNameSimple(selector)

	builder := &RuleBuilder[T, float64]{
		validator: v,
		selector:  selector,
		fieldName: fieldName,
		rules:     make([]func(*T) *ValidationError, 0),
	}

	v.rules = append(v.rules, func(instance *T, mode ValidateMode) ValidationErrors {
		errors := ValidationErrors{}

		if builder.condition != nil && !builder.condition(instance) {
			return errors
		}

		for _, rule := range builder.rules {
			if err := rule(instance); err != nil {
				errors = append(errors, *err)
				// 快速失败模式：遇到错误立即返回
				if mode == ValidateFailFast {
					return errors
				}
			}
		}
		return errors
	})

	return builder
}

// FieldSlice 为切片字段定义规则（自动提取字段名）
func FieldSlice[T any, TItem any](v *AbstractValidator[T], selector func(*T) []TItem) *RuleBuilder[T, []TItem] {
	fieldName := ExtractFieldNameSimple(selector)

	builder := &RuleBuilder[T, []TItem]{
		validator: v,
		selector:  selector,
		fieldName: fieldName,
		rules:     make([]func(*T) *ValidationError, 0),
	}

	v.rules = append(v.rules, func(instance *T, mode ValidateMode) ValidationErrors {
		errors := ValidationErrors{}

		if builder.condition != nil && !builder.condition(instance) {
			return errors
		}

		for _, rule := range builder.rules {
			if err := rule(instance); err != nil {
				errors = append(errors, *err)
				// 快速失败模式：遇到错误立即返回
				if mode == ValidateFailFast {
					return errors
				}
			}
		}
		return errors
	})

	return builder
}
