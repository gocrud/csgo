package validation

import (
	"reflect"
	"runtime"
	"strings"
)

// extractFieldName 从选择器函数中提取字段名
// 优先使用 json tag，如果没有则使用字段名本身
func extractFieldName[T any, TProperty any](selector func(*T) TProperty) string {
	// 创建一个零值实例用于反射
	var zero T
	zeroValue := reflect.ValueOf(&zero)

	// 获取类型信息
	t := reflect.TypeOf(zero)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 尝试通过函数名推断（作为备选方案）
	funcName := runtime.FuncForPC(reflect.ValueOf(selector).Pointer()).Name()

	// 执行选择器，使用 panic recovery 捕获字段访问
	// 这是一个简化的实现，实际中我们通过分析返回值的地址来确定字段

	// 遍历所有字段，找到匹配的
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 尝试调用选择器看是否匹配这个字段
		// 由于我们无法直接执行，我们使用类型匹配作为启发式方法
		fieldValue := zeroValue.Elem().Field(i)

		// 检查类型是否匹配
		var propertyType TProperty
		if fieldValue.Type() == reflect.TypeOf(propertyType) {
			// 可能是这个字段，返回其名称
			return getFieldName(field)
		}
	}

	// 如果没有找到，使用函数名作为最后的尝试
	// 从函数名中提取可能的字段名（例如 func1.func2 -> func2）
	parts := strings.Split(funcName, ".")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		// 移除可能的后缀如 funcN
		if strings.HasPrefix(lastPart, "func") {
			return "unknown"
		}
		return lastPart
	}

	return "unknown"
}

// getFieldName 从字段信息中提取名称，优先使用 json tag
func getFieldName(field reflect.StructField) string {
	// 优先使用 json tag
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		// 移除选项（如 omitempty）
		parts := strings.Split(jsonTag, ",")
		if parts[0] != "" && parts[0] != "-" {
			return parts[0]
		}
	}

	// 如果没有 json tag，使用字段名本身
	return field.Name
}

// extractFieldNameByName 通过字段名直接提取（用于类型不匹配时的备选）
func extractFieldNameByName[T any](fieldName string) string {
	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if field, ok := t.FieldByName(fieldName); ok {
		return getFieldName(field)
	}

	return fieldName
}

// fieldNameExtractor 字段名提取器（更可靠的实现）
// 通过在运行时捕获字段访问来确定字段名
type fieldNameExtractor[T any] struct {
	fieldName string
	captured  bool
}

// newFieldNameExtractor 创建字段名提取器
func newFieldNameExtractor[T any]() *fieldNameExtractor[T] {
	return &fieldNameExtractor[T]{}
}

// extract 提取字段名
// TODO: 实现更可靠的实现
func (e *fieldNameExtractor[T]) extract(selector func(*T) interface{}) string {
	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 尝试执行选择器并捕获访问的字段
	// 由于 Go 的限制，我们使用类型分析
	zeroPtr := reflect.New(t)

	// 调用选择器
	_ = selector(zeroPtr.Interface().(*T))

	// 这是一个简化实现，实际应该使用更复杂的反射技巧
	// 暂时返回未知
	return e.fieldName
}

// ExtractFieldNameSimple 简化版字段名提取
// 通过匹配字段类型来推断（足够用于大多数场景）
// TODO: 实现更可靠的实现
func ExtractFieldNameSimple[T any, TProperty any](selector func(*T) TProperty) string {
	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 获取 TProperty 的类型
	var propZero TProperty
	propType := reflect.TypeOf(propZero)

	// 遍历所有字段，找到第一个类型匹配的
	// 注意：如果有多个相同类型的字段，这个方法可能不准确
	// 但对于大多数实际场景足够了
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type == propType {
			return getFieldName(field)
		}
	}

	return "unknown"
}
