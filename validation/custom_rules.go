package validation

// CustomRule 自定义验证规则
func CustomRule[T any](v *AbstractValidator[T], predicate func(*T) error) {
	v.rules = append(v.rules, func(instance *T) ValidationErrors {
		if err := predicate(instance); err != nil {
			return ValidationErrors{
				{
					Field:   "_custom",
					Message: err.Error(),
				},
			}
		}
		return ValidationErrors{}
	})
}

// MustBeValid 嵌套验证器
func MustBeValid[T any, TProperty any](b *RuleBuilder[T, TProperty], validator IValidator[TProperty]) *RuleBuilder[T, TProperty] {
	rule := func(instance *T) *ValidationError {
		value := b.selector(instance)
		result := validator.Validate(&value)
		if !result.IsValid {
			// 返回第一个错误
			if len(result.Errors) > 0 {
				err := result.Errors[0]
				err.Field = b.fieldName + "." + err.Field
				return &err
			}
		}
		return nil
	}
	return b.addRule(rule)
}
