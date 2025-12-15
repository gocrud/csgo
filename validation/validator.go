package validation

import (
	"fmt"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
	"unsafe"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	uuidRegex  = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

// Validate 验证对象
func Validate[T any](obj *T) ValidationErrors {
	if obj == nil {
		return nil
	}
	t := reflect.TypeOf(*obj)
	schema := GetSchema(t)
	if schema == nil {
		return nil
	}

	basePtr := uintptr(unsafe.Pointer(obj))
	var errors ValidationErrors

	// 为了保证验证顺序（尤其是 FailFast 模式下的一致性），我们需要对 offset 进行排序
	// map 的遍历顺序是随机的
	var offsets []uintptr
	for offset := range schema.FieldRules {
		offsets = append(offsets, offset)
	}
	sort.Slice(offsets, func(i, j int) bool {
		return offsets[i] < offsets[j]
	})

	for _, offset := range offsets {
		rules := schema.FieldRules[offset]
		fieldName := schema.FieldNames[offset]
		if fieldName == "" {
			fieldName = fmt.Sprintf("Field_%d", offset)
		}

		// 修正：FieldKinds 存储的是反射的 Kind
		kind := schema.FieldKinds[offset]

		// 获取字段类型，用于特殊判断（如 Slice 类型或 Time 类型）
		fieldType := schema.Type.Field(getFieldIdxByOffset(schema.Type, offset)).Type

		for _, rule := range rules {
			// nolint:govet
			fieldPtr := unsafe.Pointer(basePtr + offset)

			// 保持 obj 活跃，防止在此之前被 GC 回收
			// runtime.KeepAlive(obj)

			if err := checkRule(fieldPtr, kind, rule, fieldType); err != nil {
				msg := rule.CustomMsg
				if msg == "" {
					msg = err.Error()
				}
				// 错误码格式: VALIDATION.REQUIRED (全大写)
				code := fmt.Sprintf("VALIDATION.%s", strings.ToUpper(string(rule.Type)))
				errors.Add(fieldName, msg, code)

				// 如果是快速失败模式，立即返回当前收集到的错误
				if schema.FailFast {
					return errors
				}
			}
		}
	}

	// 保持 obj 活跃，防止在此之前被 GC 回收
	// runtime.KeepAlive(obj)

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// 辅助函数：通过 offset 查找字段索引（用于获取字段详细 Type 信息）
// 这只是一个简单的线性查找，对于大型结构体可能需要优化
func getFieldIdxByOffset(t reflect.Type, offset uintptr) int {
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Offset == offset {
			return i
		}
	}
	// 递归查找嵌套结构体
	// ... (省略复杂的递归逻辑，当前假设扁平或者不需要深度 Type 信息)
	return -1
}

func checkRule(ptr unsafe.Pointer, kind reflect.Kind, rule Rule, fieldType reflect.Type) error {
	// 特殊处理 Time
	if kind == reflect.Struct && isTime(fieldType) {
		val := *(*time.Time)(ptr)
		return checkTimeRule(val, rule)
	}

	// 根据 Kind 派发检查逻辑
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// 必须根据具体类型转换
		return checkInt64Rule(getInt64Value(ptr, kind), rule)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return checkUint64Rule(getUint64Value(ptr, kind), rule)
	case reflect.Float32, reflect.Float64:
		return checkFloat64Rule(getFloat64Value(ptr, kind), rule)
	case reflect.String:
		val := *(*string)(ptr)
		return checkStringRule(val, rule)
	case reflect.Bool:
		val := *(*bool)(ptr)
		return checkBoolRule(val, rule)
	case reflect.Slice:
		// 切片类型处理
		// 对于普通长度检查，使用 SliceHeader
		if rule.Type == "min_len" || rule.Type == "max_len" || rule.Type == "len" || rule.Type == "required" {
			header := (*reflect.SliceHeader)(ptr)
			return checkSliceRule(header.Len, rule)
		}
		// 对于 Unique 等需要访问元素值的规则，需要 reflect
		// 注意：fieldType.Elem() 是切片元素的类型
		// 我们可以构建一个 reflect.Value 来操作
		val := reflect.NewAt(fieldType, ptr).Elem()
		return checkComplexSliceRule(val, rule)
	default:
		// 尝试根据 RuleType 猜测 (fallback)
		switch rule.Type {
		case "min", "max", "range":
			// 默认按 int64 尝试
			return checkInt64Rule(getInt64Value(ptr, kind), rule)
		case "min_len", "max_len", "len", "email", "pattern", "alpha", "alphanum", "numeric", "uppercase", "lowercase", "contains", "startswith", "endswith", "url", "ip", "uuid":
			val := *(*string)(ptr)
			return checkStringRule(val, rule)
		case "required":
			// 简单的通用 required 检查
			// 无法准确判断零值，只能依赖 Kind
			return fmt.Errorf("unsupported type for required check: %v", kind)
		}
	}
	return nil
}

func isTime(t reflect.Type) bool {
	return t.ConvertibleTo(reflect.TypeOf(time.Time{}))
}

// 通用切片长度验证
func checkSliceRule(length int, rule Rule) error {
	switch rule.Type {
	case "required":
		if length == 0 {
			return fmt.Errorf("is required")
		}
	case "min_len":
		min := rule.Params[0].(int)
		if length < min {
			return fmt.Errorf("length must be at least %d", min)
		}
	case "max_len":
		max := rule.Params[0].(int)
		if length > max {
			return fmt.Errorf("length must be at most %d", max)
		}
	case "len":
		l := rule.Params[0].(int)
		if length != l {
			return fmt.Errorf("length must be %d", l)
		}
	}
	return nil
}

// 复杂切片验证（Unique等）
func checkComplexSliceRule(val reflect.Value, rule Rule) error {
	switch rule.Type {
	case "unique":
		if val.Len() == 0 {
			return nil
		}
		seen := make(map[interface{}]bool)
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i).Interface()
			if seen[v] {
				return fmt.Errorf("elements must be unique")
			}
			seen[v] = true
		}
	}
	return nil
}

func checkTimeRule(val time.Time, rule Rule) error {
	switch rule.Type {
	case "required":
		if val.IsZero() {
			return fmt.Errorf("is required")
		}
	case "after":
		target := rule.Params[0].(time.Time)
		if !val.After(target) {
			return fmt.Errorf("must be after %s", target.Format(time.RFC3339))
		}
	case "before":
		target := rule.Params[0].(time.Time)
		if !val.Before(target) {
			return fmt.Errorf("must be before %s", target.Format(time.RFC3339))
		}
	}
	return nil
}

// 值获取辅助函数
func getInt64Value(ptr unsafe.Pointer, kind reflect.Kind) int64 {
	switch kind {
	case reflect.Int:
		return int64(*(*int)(ptr))
	case reflect.Int8:
		return int64(*(*int8)(ptr))
	case reflect.Int16:
		return int64(*(*int16)(ptr))
	case reflect.Int32:
		return int64(*(*int32)(ptr))
	case reflect.Int64:
		return *(*int64)(ptr)
	}
	return 0
}

func getUint64Value(ptr unsafe.Pointer, kind reflect.Kind) uint64 {
	switch kind {
	case reflect.Uint:
		return uint64(*(*uint)(ptr))
	case reflect.Uint8:
		return uint64(*(*uint8)(ptr))
	case reflect.Uint16:
		return uint64(*(*uint16)(ptr))
	case reflect.Uint32:
		return uint64(*(*uint32)(ptr))
	case reflect.Uint64:
		return *(*uint64)(ptr)
	}
	return 0
}

func getFloat64Value(ptr unsafe.Pointer, kind reflect.Kind) float64 {
	switch kind {
	case reflect.Float32:
		return float64(*(*float32)(ptr))
	case reflect.Float64:
		return *(*float64)(ptr)
	}
	return 0
}

func checkInt64Rule(val int64, rule Rule) error {
	switch rule.Type {
	case "min":
		min := toInt64(rule.Params[0])
		if val < min {
			return fmt.Errorf("must be at least %d", min)
		}
	case "max":
		max := toInt64(rule.Params[0])
		if val > max {
			return fmt.Errorf("must be at most %d", max)
		}
	case "range":
		min := toInt64(rule.Params[0])
		max := toInt64(rule.Params[1])
		if val < min || val > max {
			return fmt.Errorf("must be between %d and %d", min, max)
		}
	case "required":
		if val == 0 {
			return fmt.Errorf("is required")
		}
	case "eq":
		target := toInt64(rule.Params[0])
		if val != target {
			return fmt.Errorf("must be equal to %d", target)
		}
	case "in":
		for _, p := range rule.Params {
			if val == toInt64(p) {
				return nil
			}
		}
		return fmt.Errorf("must be one of valid values")
	}
	return nil
}

func checkUint64Rule(val uint64, rule Rule) error {
	switch rule.Type {
	case "min":
		min := toUint64(rule.Params[0])
		if val < min {
			return fmt.Errorf("must be at least %d", min)
		}
	case "max":
		max := toUint64(rule.Params[0])
		if val > max {
			return fmt.Errorf("must be at most %d", max)
		}
	case "required":
		if val == 0 {
			return fmt.Errorf("is required")
		}
	case "eq":
		target := toUint64(rule.Params[0])
		if val != target {
			return fmt.Errorf("must be equal to %d", target)
		}
	case "in":
		for _, p := range rule.Params {
			if val == toUint64(p) {
				return nil
			}
		}
		return fmt.Errorf("must be one of valid values")
	}
	return nil
}

func checkFloat64Rule(val float64, rule Rule) error {
	switch rule.Type {
	case "min":
		min := toFloat64(rule.Params[0])
		if val < min {
			return fmt.Errorf("must be at least %f", min)
		}
	case "max":
		max := toFloat64(rule.Params[0])
		if val > max {
			return fmt.Errorf("must be at most %f", max)
		}
	case "required":
		if val == 0 {
			return fmt.Errorf("is required")
		}
	case "eq":
		target := toFloat64(rule.Params[0])
		// 浮点数相等判断需要 epsilon，这里简化为直接相等
		if val != target {
			return fmt.Errorf("must be equal to %f", target)
		}
	}
	return nil
}

func checkBoolRule(val bool, rule Rule) error {
	switch rule.Type {
	case "true":
		if !val {
			return fmt.Errorf("must be true")
		}
	case "false":
		if val {
			return fmt.Errorf("must be false")
		}
	}
	return nil
}

func checkStringRule(val string, rule Rule) error {
	switch rule.Type {
	case "required":
		if val == "" {
			return fmt.Errorf("is required")
		}
	case "min_len":
		min := rule.Params[0].(int)
		if len(val) < min {
			return fmt.Errorf("length must be at least %d", min)
		}
	case "max_len":
		max := rule.Params[0].(int)
		if len(val) > max {
			return fmt.Errorf("length must be at most %d", max)
		}
	case "len":
		l := rule.Params[0].(int)
		if len(val) != l {
			return fmt.Errorf("length must be %d", l)
		}
	case "email":
		if val != "" && !emailRegex.MatchString(val) {
			return fmt.Errorf("invalid email format")
		}
	case "url":
		if val != "" {
			_, err := url.ParseRequestURI(val)
			if err != nil {
				return fmt.Errorf("invalid url format")
			}
		}
	case "ip":
		if val != "" && net.ParseIP(val) == nil {
			return fmt.Errorf("invalid ip address")
		}
	case "uuid":
		if val != "" && !uuidRegex.MatchString(val) {
			return fmt.Errorf("invalid uuid format")
		}
	case "pattern":
		pat := rule.Params[0].(string)
		matched, _ := regexp.MatchString(pat, val)
		if !matched {
			return fmt.Errorf("must match pattern %s", pat)
		}
	case "alpha":
		for _, r := range val {
			if !unicode.IsLetter(r) {
				return fmt.Errorf("must contain only letters")
			}
		}
	case "alphanum":
		for _, r := range val {
			if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
				return fmt.Errorf("must contain only letters and numbers")
			}
		}
	case "numeric":
		for _, r := range val {
			if !unicode.IsNumber(r) {
				return fmt.Errorf("must contain only numbers")
			}
		}
	case "lowercase":
		if val != strings.ToLower(val) {
			return fmt.Errorf("must be lowercase")
		}
	case "uppercase":
		if val != strings.ToUpper(val) {
			return fmt.Errorf("must be uppercase")
		}
	case "contains":
		sub := rule.Params[0].(string)
		if !strings.Contains(val, sub) {
			return fmt.Errorf("must contain %s", sub)
		}
	case "startswith":
		prefix := rule.Params[0].(string)
		if !strings.HasPrefix(val, prefix) {
			return fmt.Errorf("must start with %s", prefix)
		}
	case "endswith":
		suffix := rule.Params[0].(string)
		if !strings.HasSuffix(val, suffix) {
			return fmt.Errorf("must end with %s", suffix)
		}
	case "in":
		for _, p := range rule.Params {
			if val == p.(string) {
				return nil
			}
		}
		return fmt.Errorf("must be one of valid values")
	}
	return nil
}

// 辅助转换函数
func toInt64(i interface{}) int64 {
	switch v := i.(type) {
	case int:
		return int64(v)
	case int64:
		return v
	case int32:
		return int64(v)
	case int16:
		return int64(v)
	case int8:
		return int64(v)
	}
	return 0
}

func toUint64(i interface{}) uint64 {
	switch v := i.(type) {
	case uint:
		return uint64(v)
	case uint64:
		return v
	case uint32:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint8:
		return uint64(v)
	case int: // 允许传入 int 常量
		if v >= 0 {
			return uint64(v)
		}
	}
	return 0
}

func toFloat64(i interface{}) float64 {
	switch v := i.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	}
	return 0
}
