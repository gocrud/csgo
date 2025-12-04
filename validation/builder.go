package validation

// RuleBuilder 规则构建器
type RuleBuilder[T any, TProperty any] struct {
	validator *AbstractValidator[T]
	selector  func(*T) TProperty
	fieldName string
	rules     []func(*T) *ValidationError
	condition func(*T) bool
}

// addRule 添加规则到链
func (b *RuleBuilder[T, TProperty]) addRule(rule func(*T) *ValidationError) *RuleBuilder[T, TProperty] {
	b.rules = append(b.rules, rule)
	return b
}

// When 条件验证
func (b *RuleBuilder[T, TProperty]) When(condition func(*T) bool) *RuleBuilder[T, TProperty] {
	b.condition = condition
	return b
}

// Unless 反向条件验证
func (b *RuleBuilder[T, TProperty]) Unless(condition func(*T) bool) *RuleBuilder[T, TProperty] {
	b.condition = func(t *T) bool {
		return !condition(t)
	}
	return b
}

// WithMessage 自定义错误消息（需要在规则之后调用）
func (b *RuleBuilder[T, TProperty]) WithMessage(message string) *RuleBuilder[T, TProperty] {
	if len(b.rules) > 0 {
		// 包装最后一个规则，修改其错误消息
		lastRule := b.rules[len(b.rules)-1]
		b.rules[len(b.rules)-1] = func(instance *T) *ValidationError {
			if err := lastRule(instance); err != nil {
				err.Message = message
				return err
			}
			return nil
		}
	}
	return b
}

// WithCode 设置错误码
func (b *RuleBuilder[T, TProperty]) WithCode(code string) *RuleBuilder[T, TProperty] {
	if len(b.rules) > 0 {
		lastRule := b.rules[len(b.rules)-1]
		b.rules[len(b.rules)-1] = func(instance *T) *ValidationError {
			if err := lastRule(instance); err != nil {
				err.Code = code
				return err
			}
			return nil
		}
	}
	return b
}

// getString 从泛型值中获取 string（用于 string 类型的 RuleBuilder）
func (b *RuleBuilder[T, string]) getString(instance *T) string {
	val := b.selector(instance)
	// 使用 any 转换来确保返回具体的 string 类型
	return any(val).(string)
}

// getInt 从泛型值中获取 int
func (b *RuleBuilder[T, int]) getInt(instance *T) int {
	return b.selector(instance)
}

// getInt64 从泛型值中获取 int64
func (b *RuleBuilder[T, int64]) getInt64(instance *T) int64 {
	return b.selector(instance)
}

// getFloat64 从泛型值中获取 float64
func (b *RuleBuilder[T, float64]) getFloat64(instance *T) float64 {
	return b.selector(instance)
}
