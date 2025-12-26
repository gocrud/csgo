package web

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"github.com/gocrud/csgo/errors"
	"github.com/gocrud/csgo/validation"
)

// ParamChain 参数验证链，支持泛型和链式调用。
//
// 使用示例：
//
//	id := web.Path[int](c, "id").Min(1).Required().Value()
//	page := web.Query[int](c, "page").Default(1)
//	email := web.Query[string](c, "email").Required().Email().Value()
type ParamChain[T any] struct {
	key      string       // 参数名
	value    T            // 解析后的值
	rawValue string       // 原始字符串值
	source   string       // 参数来源: "path", "query", "header", "form"
	ctx      *HttpContext // HTTP 上下文
	hasValue bool         // 是否成功解析了值
	required bool         // 是否必填
	hasError bool         // 是否已经有错误
}

// ==================== 入口函数 ====================

// Path 获取路径参数并转换为指定类型 T。
// 路径参数默认为必填。
//
// 使用示例：
//
//	id := web.Path[int](c, "id").Min(1).Value()
func Path[T any](c *HttpContext, key string) *ParamChain[T] {
	rawValue := c.gin.Param(key)
	return newParamChain[T](c, key, rawValue, "path", true)
}

// Query 获取查询参数并转换为指定类型 T。
// 查询参数默认为可选。
//
// 使用示例：
//
//	page := web.Query[int](c, "page").Default(1)
//	email := web.Query[string](c, "email").Required().Email().Value()
func Query[T any](c *HttpContext, key string) *ParamChain[T] {
	rawValue := c.gin.Query(key)
	return newParamChain[T](c, key, rawValue, "query", false)
}

// Header 获取请求头并转换为指定类型 T。
// 请求头默认为可选。
//
// 使用示例：
//
//	token := web.Header[string](c, "Authorization").Required().MinLength(10).Value()
func Header[T any](c *HttpContext, key string) *ParamChain[T] {
	rawValue := c.gin.GetHeader(key)
	return newParamChain[T](c, key, rawValue, "header", false)
}

// Form 获取表单参数并转换为指定类型 T。
// 表单参数默认为可选。
//
// 使用示例：
//
//	username := web.Form[string](c, "username").Required().MinLength(3).Value()
func Form[T any](c *HttpContext, key string) *ParamChain[T] {
	rawValue := c.gin.PostForm(key)
	return newParamChain[T](c, key, rawValue, "form", false)
}

// ==================== 核心内部函数 ====================

// newParamChain 创建新的参数链。
func newParamChain[T any](c *HttpContext, key, rawValue, source string, required bool) *ParamChain[T] {
	chain := &ParamChain[T]{
		key:      key,
		rawValue: rawValue,
		source:   source,
		ctx:      c,
		required: required,
	}

	// 尝试解析值
	if rawValue != "" {
		if parsed, err := parseValue[T](rawValue); err == nil {
			chain.value = parsed
			chain.hasValue = true
		} else {
			// 解析失败，记录错误
			chain.addError(getParseErrorMessage[T](rawValue), errors.ValidationInvalidInteger)
		}
	}

	return chain
}

// parseValue 将字符串解析为目标类型 T。
func parseValue[T any](raw string) (T, error) {
	var zero T
	targetType := reflect.TypeOf(zero)
	// 处理指针类型
	if targetType != nil && targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}

	// 如果是 string 类型，直接返回
	if _, ok := any(zero).(string); ok {
		return any(raw).(T), nil
	}

	// 使用反射处理其他类型
	if targetType == nil {
		return zero, fmt.Errorf("无法确定目标类型")
	}

	switch targetType.Kind() {
	case reflect.Int:
		v, err := strconv.Atoi(raw)
		if err != nil {
			return zero, err
		}
		return *(*T)(unsafe.Pointer(&v)), nil

	case reflect.Int8:
		v, err := strconv.ParseInt(raw, 10, 8)
		if err != nil {
			return zero, err
		}
		val := int8(v)
		return *(*T)(unsafe.Pointer(&val)), nil

	case reflect.Int16:
		v, err := strconv.ParseInt(raw, 10, 16)
		if err != nil {
			return zero, err
		}
		val := int16(v)
		return *(*T)(unsafe.Pointer(&val)), nil

	case reflect.Int32:
		v, err := strconv.ParseInt(raw, 10, 32)
		if err != nil {
			return zero, err
		}
		val := int32(v)
		return *(*T)(unsafe.Pointer(&val)), nil

	case reflect.Int64:
		// 特殊处理 time.Duration
		if targetType.String() == "time.Duration" {
			v, err := time.ParseDuration(raw)
			if err != nil {
				return zero, err
			}
			return *(*T)(unsafe.Pointer(&v)), nil
		}
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return zero, err
		}
		return *(*T)(unsafe.Pointer(&v)), nil

	case reflect.Uint:
		v, err := strconv.ParseUint(raw, 10, 0)
		if err != nil {
			return zero, err
		}
		val := uint(v)
		return *(*T)(unsafe.Pointer(&val)), nil

	case reflect.Uint8:
		v, err := strconv.ParseUint(raw, 10, 8)
		if err != nil {
			return zero, err
		}
		val := uint8(v)
		return *(*T)(unsafe.Pointer(&val)), nil

	case reflect.Uint16:
		v, err := strconv.ParseUint(raw, 10, 16)
		if err != nil {
			return zero, err
		}
		val := uint16(v)
		return *(*T)(unsafe.Pointer(&val)), nil

	case reflect.Uint32:
		v, err := strconv.ParseUint(raw, 10, 32)
		if err != nil {
			return zero, err
		}
		val := uint32(v)
		return *(*T)(unsafe.Pointer(&val)), nil

	case reflect.Uint64:
		v, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return zero, err
		}
		return *(*T)(unsafe.Pointer(&v)), nil

	case reflect.Float32:
		v, err := strconv.ParseFloat(raw, 32)
		if err != nil {
			return zero, err
		}
		val := float32(v)
		return *(*T)(unsafe.Pointer(&val)), nil

	case reflect.Float64:
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return zero, err
		}
		return *(*T)(unsafe.Pointer(&v)), nil

	case reflect.Bool:
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return zero, err
		}
		return *(*T)(unsafe.Pointer(&v)), nil

	case reflect.Struct:
		// 特殊处理 time.Time
		if targetType.String() == "time.Time" {
			// 尝试多种时间格式
			formats := []string{
				time.RFC3339,
				"2006-01-02 15:04:05",
				"2006-01-02",
				time.RFC3339Nano,
			}
			for _, format := range formats {
				if t, err := time.Parse(format, raw); err == nil {
					return any(t).(T), nil
				}
			}
			return zero, fmt.Errorf("无法解析时间格式: %s", raw)
		}
	}

	return zero, fmt.Errorf("不支持的类型: %v", targetType)
}

// getParseErrorMessage 获取解析错误的友好提示信息。
func getParseErrorMessage[T any](raw string) string {
	var zero T
	targetType := reflect.TypeOf(zero)

	if targetType == nil {
		return fmt.Sprintf("无法解析值: %s", raw)
	}

	switch targetType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "必须是有效的整数"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "必须是有效的正整数"
	case reflect.Float32, reflect.Float64:
		return "必须是有效的数字"
	case reflect.Bool:
		return "必须是有效的布尔值 (true/false)"
	default:
		if targetType.String() == "time.Time" {
			return "必须是有效的时间格式"
		}
		return fmt.Sprintf("无法解析为 %s 类型", targetType.String())
	}
}

// ==================== 通用验证方法 ====================

// Required 标记参数为必填。
// 如果参数为空或解析失败，将添加验证错误。
func (p *ParamChain[T]) Required() *ParamChain[T] {
	p.required = true
	if !p.hasValue && !p.hasError {
		p.addError("不能为空", errors.ValidationRequired)
	}
	return p
}

// Optional 显式标记参数为可选（通常不需要调用，query/header 默认就是可选的）。
func (p *ParamChain[T]) Optional() *ParamChain[T] {
	p.required = false
	return p
}

// Custom 使用自定义验证函数。
//
// 使用示例：
//
//	age := web.Query[int](c, "age").Custom(func(v int) error {
//	    if v < 18 || v > 120 {
//	        return errors.New("年龄必须在 18-120 之间")
//	    }
//	    return nil
//	})
//
//	// Min/Max 验证示例
//	page := web.Query[int](c, "page").Custom(func(v int) error {
//	    if v < 1 {
//	        return errors.New("不能小于 1")
//	    }
//	    if v > 100 {
//	        return errors.New("不能大于 100")
//	    }
//	    return nil
//	}).Value()
//
//	// 字符串长度验证示例
//	username := web.Query[string](c, "username").Custom(func(v string) error {
//	    if len(v) < 3 {
//	        return errors.New("长度不能少于 3 个字符")
//	    }
//	    if len(v) > 20 {
//	        return errors.New("长度不能超过 20 个字符")
//	    }
//	    return nil
//	}).Value()
func (p *ParamChain[T]) Custom(fn func(T) error) *ParamChain[T] {
	if p.hasValue && !p.hasError {
		if err := fn(p.value); err != nil {
			p.addError(err.Error(), errors.ValidationRequired) // 使用已有的错误码
		}
	}
	return p
}

// ==================== 获取值的方法 ====================

// Value 获取参数值，如果有验证错误则返回错误响应。
// 返回 (值, IActionResult)，如果验证成功则 IActionResult 为 nil。
//
// 使用示例：
//
//	id, err := web.Path[int](c, "id").Min(1).Value()
//	if err != nil {
//	    return err  // 立即返回验证错误
//	}
//	// 继续执行业务逻辑
func (p *ParamChain[T]) Value() (T, IActionResult) {
	// 检查必填验证
	if p.required && !p.hasValue && !p.hasError {
		p.hasError = true
		var zero T
		return zero, p.buildErrorResult("不能为空", errors.ValidationRequired)
	}

	// 如果有错误，返回错误响应
	if p.hasError {
		var zero T
		return zero, p.buildErrorResult("", "")
	}

	return p.value, nil
}

// Default 返回参数值，如果参数不存在或解析失败则返回默认值。
// 此方法会忽略所有验证错误，直接返回默认值。
//
// 使用示例：
//
//	page := web.Query[int](c, "page").Default(1)
func (p *ParamChain[T]) Default(defaultValue T) T {
	if !p.hasValue || p.hasError {
		return defaultValue
	}
	return p.value
}

// ValueOr 是 Default 的别名，提供更直观的语义。
func (p *ParamChain[T]) ValueOr(defaultValue T) T {
	return p.Default(defaultValue)
}

// Get 返回参数值和标准 error。
// 使用此方法可以手动处理验证错误（返回自定义错误消息）。
//
// 使用示例：
//
//	id, err := web.Path[int](c, "id").Min(1).Get()
//	if err != nil {
//	    return c.BadRequest("ID 无效: " + err.Error())
//	}
func (p *ParamChain[T]) Get() (T, error) {
	// 检查必填验证
	if p.required && !p.hasValue && !p.hasError {
		var zero T
		return zero, fmt.Errorf("%s 不能为空", p.key)
	}

	if p.hasError {
		var zero T
		// 返回第一个错误信息
		if len(p.ctx.paramErrors) > 0 {
			for _, e := range p.ctx.paramErrors {
				if e.Field == p.key {
					return zero, fmt.Errorf("%s", e.Message)
				}
			}
		}
		return zero, fmt.Errorf("%s 验证失败", p.key)
	}

	return p.value, nil
}

// HasValue 检查参数是否有值（是否成功解析）。
func (p *ParamChain[T]) HasValue() bool {
	return p.hasValue
}

// ==================== 内部辅助方法 ====================

// addError 添加验证错误到内部错误列表。
func (p *ParamChain[T]) addError(message, code string) {
	p.hasError = true
	if p.ctx != nil {
		p.ctx.addParamError(validation.ValidationError{
			Field:   p.key,
			Message: message,
			Code:    code,
		})
	}
}

// buildErrorResult 构建验证错误的 IActionResult。
// 如果 message 为空，则从 ctx.paramErrors 中查找该字段的错误。
func (p *ParamChain[T]) buildErrorResult(message, code string) IActionResult {
	if p.ctx == nil {
		return BadRequest("参数验证失败")
	}

	// 如果有收集的错误，返回结构化的验证错误响应
	if len(p.ctx.paramErrors) > 0 {
		return ValidationBadRequest(p.ctx.paramErrors)
	}

	// 否则返回单个错误消息
	if message != "" {
		return BadRequestWithCode(code, message)
	}

	return BadRequest("参数验证失败")
}

// isComparable 检查类型是否可比较（用于 In 方法）。
func isComparable[T any]() bool {
	var zero T
	typ := reflect.TypeOf(zero)
	return typ != nil && typ.Comparable()
}
