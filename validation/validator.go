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

// AbstractValidator 抽象验证器基类
type AbstractValidator[T any] struct {
	rules []func(*T) ValidationErrors
}

// NewValidator 创建新的验证器
func NewValidator[T any]() *AbstractValidator[T] {
	return &AbstractValidator[T]{
		rules: make([]func(*T) ValidationErrors, 0),
	}
}

// Validate 执行验证
func (v *AbstractValidator[T]) Validate(instance *T) ValidationResult {
	errors := ValidationErrors{}

	for _, rule := range v.rules {
		ruleErrors := rule(instance)
		errors = append(errors, ruleErrors...)
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
	v.rules = append(v.rules, func(instance *T) ValidationErrors {
		errors := ValidationErrors{}

		// 检查条件
		if baseBuilder.condition != nil && !baseBuilder.condition(instance) {
			return errors
		}

		// 执行所有规则
		for _, rule := range baseBuilder.rules {
			if err := rule(instance); err != nil {
				errors = append(errors, *err)
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

	v.rules = append(v.rules, func(instance *T) ValidationErrors {
		errors := ValidationErrors{}

		if builder.condition != nil && !builder.condition(instance) {
			return errors
		}

		for _, rule := range builder.rules {
			if err := rule(instance); err != nil {
				errors = append(errors, *err)
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

	v.rules = append(v.rules, func(instance *T) ValidationErrors {
		errors := ValidationErrors{}

		if builder.condition != nil && !builder.condition(instance) {
			return errors
		}

		for _, rule := range builder.rules {
			if err := rule(instance); err != nil {
				errors = append(errors, *err)
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

	v.rules = append(v.rules, func(instance *T) ValidationErrors {
		errors := ValidationErrors{}

		if builder.condition != nil && !builder.condition(instance) {
			return errors
		}

		for _, rule := range builder.rules {
			if err := rule(instance); err != nil {
				errors = append(errors, *err)
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

	v.rules = append(v.rules, func(instance *T) ValidationErrors {
		errors := ValidationErrors{}

		if builder.condition != nil && !builder.condition(instance) {
			return errors
		}

		for _, rule := range builder.rules {
			if err := rule(instance); err != nil {
				errors = append(errors, *err)
			}
		}
		return errors
	})

	return builder
}
