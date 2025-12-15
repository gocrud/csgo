package v

import (
	"reflect"
	"strings"
)

// fieldTracker 字段追踪器，用于收集验证规则
type fieldTracker struct {
	rules map[string][]Rule // 字段路径 -> 规则列表
}

// newFieldTracker 创建新的字段追踪器
func newFieldTracker() *fieldTracker {
	return &fieldTracker{
		rules: make(map[string][]Rule),
	}
}

// addStringRule 添加字符串规则
func (t *fieldTracker) addStringRule(fieldPath string, rule StringRule) {
	t.rules[fieldPath] = append(t.rules[fieldPath], rule)
}

// addIntRule 添加整数规则
func (t *fieldTracker) addIntRule(fieldPath string, rule IntRule) {
	t.rules[fieldPath] = append(t.rules[fieldPath], rule)
}

// addInt64Rule 添加 int64 规则
func (t *fieldTracker) addInt64Rule(fieldPath string, rule Int64Rule) {
	t.rules[fieldPath] = append(t.rules[fieldPath], rule)
}

// addFloat64Rule 添加 float64 规则
func (t *fieldTracker) addFloat64Rule(fieldPath string, rule Float64Rule) {
	t.rules[fieldPath] = append(t.rules[fieldPath], rule)
}

// addSliceRule 添加切片规则
func (t *fieldTracker) addSliceRule(fieldPath string, rule SliceRule) {
	t.rules[fieldPath] = append(t.rules[fieldPath], rule)
}

// setLastMessage 设置最后一个规则的错误消息
func (t *fieldTracker) setLastMessage(fieldPath string, msg string) {
	if rules, ok := t.rules[fieldPath]; ok && len(rules) > 0 {
		lastRule := rules[len(rules)-1]
		lastRule.SetMessage(msg)
	}
}

// buildTrackedInstance 创建带追踪的结构体实例
func buildTrackedInstance[T any](tracker *fieldTracker) T {
	var zero T
	typ := reflect.TypeOf(zero)
	
	// 如果是指针类型，获取其元素类型
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	
	// 创建值
	val := buildTrackedValue(typ, "", tracker)
	
	return val.Interface().(T)
}

// buildTrackedValue 递归构建带追踪的值
func buildTrackedValue(typ reflect.Type, pathPrefix string, tracker *fieldTracker) reflect.Value {
	switch typ.Kind() {
	case reflect.Struct:
		return buildTrackedStruct(typ, pathPrefix, tracker)
	default:
		// 对于非结构体类型，返回零值
		return reflect.Zero(typ)
	}
}

// buildTrackedStruct 构建带追踪的结构体
func buildTrackedStruct(typ reflect.Type, pathPrefix string, tracker *fieldTracker) reflect.Value {
	val := reflect.New(typ).Elem()
	
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldType := field.Type
		
		// 跳过未导出的字段
		if !field.IsExported() {
			continue
		}
		
		// 获取字段名（优先使用 json tag，否则使用字段名）
		fieldName := getFieldName(field)
		
		// 构建字段路径
		fieldPath := fieldName
		if pathPrefix != "" {
			fieldPath = pathPrefix + "." + fieldName
		}
		
		// 根据字段类型创建对应的包装类型
		fieldValue := createWrappedField(fieldType, fieldPath, tracker)
		
		if fieldValue.IsValid() {
			val.Field(i).Set(fieldValue)
		}
	}
	
	return val
}

// createWrappedField 创建包装字段
func createWrappedField(fieldType reflect.Type, fieldPath string, tracker *fieldTracker) reflect.Value {
	// 检查是否是我们的包装类型
	typeName := fieldType.String()
	
	switch typeName {
	case "v.String":
		str := newString("", fieldPath, tracker)
		return reflect.ValueOf(str)
		
	case "v.Int":
		i := newInt(0, fieldPath, tracker)
		return reflect.ValueOf(i)
		
	case "v.Int64":
		i := newInt64(0, fieldPath, tracker)
		return reflect.ValueOf(i)
		
	case "v.Float64":
		f := newFloat64(0, fieldPath, tracker)
		return reflect.ValueOf(f)
		
	case "v.Bool":
		b := newBool(false, fieldPath, tracker)
		return reflect.ValueOf(b)
	}
	
	// 检查是否是 Slice 类型
	if strings.HasPrefix(typeName, "v.Slice[") {
		// 创建一个空的 Slice 实例
		// 由于泛型的限制，这里需要通过反射创建
		val := reflect.New(fieldType).Elem()
		
		// 设置 fieldPath 和 tracker
		// 通过反射设置私有字段
		if val.NumField() >= 3 {
			// fieldPath 字段（索引 1）
			fieldPathField := val.Field(1)
			if fieldPathField.CanSet() {
				fieldPathField.SetString(fieldPath)
			} else {
				// 使用 unsafe 方式设置私有字段
				reflect.NewAt(fieldPathField.Type(), fieldPathField.Addr().UnsafePointer()).
					Elem().SetString(fieldPath)
			}
			
			// tracker 字段（索引 2）
			trackerField := val.Field(2)
			if trackerField.CanSet() {
				trackerField.Set(reflect.ValueOf(tracker))
			} else {
				reflect.NewAt(trackerField.Type(), trackerField.Addr().UnsafePointer()).
					Elem().Set(reflect.ValueOf(tracker))
			}
		}
		
		return val
	}
	
	// 如果是嵌套结构体，递归处理
	if fieldType.Kind() == reflect.Struct {
		return buildTrackedStruct(fieldType, fieldPath, tracker)
	}
	
	// 其他类型返回零值
	return reflect.Zero(fieldType)
}

// getFieldName 获取字段名（优先 json tag，否则使用字段名）
func getFieldName(field reflect.StructField) string {
	// 首先尝试获取 json tag
	jsonTag := field.Tag.Get("json")
	if jsonTag != "" {
		// json tag 可能包含选项，如 "name,omitempty"
		// 只取第一部分
		parts := strings.Split(jsonTag, ",")
		if parts[0] != "" && parts[0] != "-" {
			return parts[0]
		}
	}
	
	// 如果没有 json tag，使用字段名并转换为小驼峰
	return toLowerCamelCase(field.Name)
}

// toLowerCamelCase 将字符串转换为小驼峰命名
func toLowerCamelCase(s string) string {
	if s == "" {
		return ""
	}
	// 简单实现：将首字母小写
	runes := []rune(s)
	runes[0] = toLower(runes[0])
	return string(runes)
}

// toLower 将字符转换为小写
func toLower(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + ('a' - 'A')
	}
	return r
}
