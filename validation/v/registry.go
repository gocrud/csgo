package v

import (
	"fmt"
	"reflect"
	
	"github.com/gocrud/csgo/validation"
)

// ValidatorMetadata 验证器元数据
type ValidatorMetadata struct {
	TypeName string
	Rules    map[string][]Rule // 字段路径 -> 规则列表
}

// 全局元数据注册表
var metadataRegistry = make(map[string]*ValidatorMetadata)

// Register 注册验证函数并收集元数据
func Register[T any](validateFunc func(T)) {
	typeName := getTypeName[T]()
	
	// 创建字段追踪器
	tracker := newFieldTracker()
	
	// 创建带追踪的实例
	trackedInstance := buildTrackedInstance[T](tracker)
	
	// 调用验证函数，收集规则
	validateFunc(trackedInstance)
	
	// 提取并存储元数据
	metadata := &ValidatorMetadata{
		TypeName: typeName,
		Rules:    tracker.rules,
	}
	
	metadataRegistry[typeName] = metadata
}

// Validate 执行验证（使用预注册的元数据）
func Validate[T any](instance *T) validation.ValidationResult {
	typeName := getTypeName[T]()
	metadata, ok := metadataRegistry[typeName]
	if !ok {
		// 如果没有注册验证器，返回成功
		return validation.SuccessResult()
	}
	
	errors := validation.ValidationErrors{}
	
	// 遍历所有字段规则
	for fieldPath, rules := range metadata.Rules {
		// 获取字段值
		fieldValue := getFieldValueByPath(instance, fieldPath)
		
		// 执行所有规则
		for _, rule := range rules {
			if err := rule.Validate(fieldValue); err != nil {
				errors = append(errors, validation.ValidationError{
					Field:   fieldPath,
					Message: err.Error(),
					Code:    "",
				})
			}
		}
	}
	
	return validation.NewValidationResult(errors)
}

// getTypeName 获取类型名称
func getTypeName[T any]() string {
	var zero T
	typ := reflect.TypeOf(zero)
	
	// 如果是指针，获取元素类型
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	
	return typ.String()
}

// getFieldValueByPath 根据字段路径获取字段值
func getFieldValueByPath(instance interface{}, path string) interface{} {
	val := reflect.ValueOf(instance)
	
	// 如果是指针，获取其指向的值
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	// 分割路径
	parts := splitPath(path)
	
	// 逐级访问字段
	for _, part := range parts {
		if val.Kind() != reflect.Struct {
			return nil
		}
		
		// 查找字段
		field := findFieldByName(val, part)
		if !field.IsValid() {
			return nil
		}
		
		val = field
	}
	
	// 如果是包装类型，提取实际值
	return extractValue(val)
}

// splitPath 分割字段路径
func splitPath(path string) []string {
	if path == "" {
		return []string{}
	}
	
	parts := []string{}
	current := ""
	
	for _, ch := range path {
		if ch == '.' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	
	if current != "" {
		parts = append(parts, current)
	}
	
	return parts
}

// findFieldByName 根据名称查找字段（支持 json tag 和字段名）
func findFieldByName(val reflect.Value, name string) reflect.Value {
	typ := val.Type()
	
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		
		// 检查 json tag
		if getFieldName(field) == name {
			return val.Field(i)
		}
	}
	
	return reflect.Value{}
}

// extractValue 从包装类型中提取实际值
func extractValue(val reflect.Value) interface{} {
	if !val.IsValid() {
		return nil
	}
	
	typeName := val.Type().String()
	
	// 对于包装类型，使用反射访问私有字段
	switch typeName {
	case "v.String":
		// 获取 value 字段（第一个字段）
		if val.NumField() > 0 {
			valueField := val.Field(0)
			// 使用 Elem() 来处理可能的指针
			if valueField.CanInterface() {
				return valueField.Interface()
			}
			// 如果不能直接访问，使用 unsafe 方式
			return reflect.NewAt(valueField.Type(), valueField.Addr().UnsafePointer()).Elem().Interface()
		}
		
	case "v.Int":
		if val.NumField() > 0 {
			valueField := val.Field(0)
			if valueField.CanInterface() {
				return valueField.Interface()
			}
			return reflect.NewAt(valueField.Type(), valueField.Addr().UnsafePointer()).Elem().Interface()
		}
		
	case "v.Int64":
		if val.NumField() > 0 {
			valueField := val.Field(0)
			if valueField.CanInterface() {
				return valueField.Interface()
			}
			return reflect.NewAt(valueField.Type(), valueField.Addr().UnsafePointer()).Elem().Interface()
		}
		
	case "v.Float64":
		if val.NumField() > 0 {
			valueField := val.Field(0)
			if valueField.CanInterface() {
				return valueField.Interface()
			}
			return reflect.NewAt(valueField.Type(), valueField.Addr().UnsafePointer()).Elem().Interface()
		}
		
	case "v.Bool":
		if val.NumField() > 0 {
			valueField := val.Field(0)
			if valueField.CanInterface() {
				return valueField.Interface()
			}
			return reflect.NewAt(valueField.Type(), valueField.Addr().UnsafePointer()).Elem().Interface()
		}
	}
	
	// 检查是否是 Slice 类型
	if len(typeName) >= 8 && typeName[:8] == "v.Slice[" {
		// 对于 Slice 类型，返回第一个字段（value）
		if val.NumField() > 0 {
			valueField := val.Field(0)
			if valueField.CanInterface() {
				return valueField.Interface()
			}
			return reflect.NewAt(valueField.Type(), valueField.Addr().UnsafePointer()).Elem().Interface()
		}
	}
	
	// 其他类型直接返回
	if val.CanInterface() {
		return val.Interface()
	}
	
	return nil
}

// GetMetadata 获取类型的验证元数据（用于测试和调试）
func GetMetadata[T any]() (*ValidatorMetadata, bool) {
	typeName := getTypeName[T]()
	metadata, ok := metadataRegistry[typeName]
	return metadata, ok
}

// ClearRegistry 清空注册表（用于测试）
func ClearRegistry() {
	metadataRegistry = make(map[string]*ValidatorMetadata)
}

// PrintMetadata 打印元数据（用于调试）
func PrintMetadata[T any]() {
	metadata, ok := GetMetadata[T]()
	if !ok {
		fmt.Printf("类型 %s 未注册验证器\n", getTypeName[T]())
		return
	}
	
	fmt.Printf("验证器元数据: %s\n", metadata.TypeName)
	for fieldPath, rules := range metadata.Rules {
		fmt.Printf("  字段 %s: %d 个规则\n", fieldPath, len(rules))
	}
}
