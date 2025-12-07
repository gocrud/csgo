package web

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gocrud/csgo/errors"
	"github.com/gocrud/csgo/validation"
)

// ParamValidator 参数验证器，支持链式调用
type ParamValidator struct {
	ctx    *HttpContext
	errors validation.ValidationErrors
}

// Params 创建参数验证器
func (c *HttpContext) Params() *ParamValidator {
	return &ParamValidator{
		ctx:    c,
		errors: validation.ValidationErrors{},
	}
}

// ==================== Path 参数验证 ====================

// PathString 获取并验证 path 字符串参数
func (v *ParamValidator) PathString(key string) *StringParamRule {
	value := v.ctx.Param(key)
	return &StringParamRule{
		validator: v,
		key:       key,
		value:     value,
		source:    "path",
		optional:  false, // path 参数默认必填
	}
}

// PathInt 获取并验证 path 整数参数
func (v *ParamValidator) PathInt(key string) *IntParamRule {
	valueStr := v.ctx.Param(key)
	value, err := strconv.Atoi(valueStr)

	rule := &IntParamRule{
		validator: v,
		key:       key,
		value:     value,
		source:    "path",
		parseErr:  err,
		optional:  false,
	}

	// 如果解析失败，添加验证错误
	if err != nil {
		v.addError(key, "必须是有效的整数", errors.ValidationInvalidInteger)
	}

	return rule
}

// PathInt64 获取并验证 path int64 参数
func (v *ParamValidator) PathInt64(key string) *Int64ParamRule {
	valueStr := v.ctx.Param(key)
	value, err := strconv.ParseInt(valueStr, 10, 64)

	rule := &Int64ParamRule{
		validator: v,
		key:       key,
		value:     value,
		source:    "path",
		parseErr:  err,
		optional:  false,
	}

	if err != nil {
		v.addError(key, "必须是有效的整数", errors.ValidationInvalidInteger)
	}

	return rule
}

// ==================== Query 参数验证 ====================

// QueryString 获取并验证 query 字符串参数
func (v *ParamValidator) QueryString(key string) *StringParamRule {
	value := v.ctx.Query(key)
	return &StringParamRule{
		validator: v,
		key:       key,
		value:     value,
		source:    "query",
		optional:  true, // query 参数默认可选
	}
}

// QueryInt 获取并验证 query 整数参数
func (v *ParamValidator) QueryInt(key string) *IntParamRule {
	valueStr := v.ctx.Query(key)

	rule := &IntParamRule{
		validator: v,
		key:       key,
		source:    "query",
		optional:  true,
	}

	// 只有参数存在时才尝试解析
	if valueStr != "" {
		value, err := strconv.Atoi(valueStr)
		rule.value = value
		rule.parseErr = err

		if err != nil {
			v.addError(key, "必须是有效的整数", errors.ValidationInvalidInteger)
		}
	}

	return rule
}

// QueryInt64 获取并验证 query int64 参数
func (v *ParamValidator) QueryInt64(key string) *Int64ParamRule {
	valueStr := v.ctx.Query(key)

	rule := &Int64ParamRule{
		validator: v,
		key:       key,
		source:    "query",
		optional:  true,
	}

	if valueStr != "" {
		value, err := strconv.ParseInt(valueStr, 10, 64)
		rule.value = value
		rule.parseErr = err

		if err != nil {
			v.addError(key, "必须是有效的整数", errors.ValidationInvalidInteger)
		}
	}

	return rule
}

// QueryBool 获取并验证 query 布尔参数
func (v *ParamValidator) QueryBool(key string) *BoolParamRule {
	valueStr := v.ctx.Query(key)

	rule := &BoolParamRule{
		validator: v,
		key:       key,
		source:    "query",
		optional:  true,
	}

	if valueStr != "" {
		value, err := strconv.ParseBool(valueStr)
		rule.value = value
		rule.parseErr = err

		if err != nil {
			v.addError(key, "必须是有效的布尔值 (true/false)", errors.ValidationInvalidBoolean)
		}
	}

	return rule
}

// QueryFloat 获取并验证 query 浮点数参数
func (v *ParamValidator) QueryFloat(key string) *FloatParamRule {
	valueStr := v.ctx.Query(key)

	rule := &FloatParamRule{
		validator: v,
		key:       key,
		source:    "query",
		optional:  true,
	}

	if valueStr != "" {
		value, err := strconv.ParseFloat(valueStr, 64)
		rule.value = value
		rule.parseErr = err

		if err != nil {
			v.addError(key, "必须是有效的数字", errors.ValidationInvalidNumber)
		}
	}

	return rule
}

// ==================== Header 参数验证 ====================

// HeaderString 获取并验证 header 字符串参数
func (v *ParamValidator) HeaderString(key string) *StringParamRule {
	value := v.ctx.GetHeader(key)
	return &StringParamRule{
		validator: v,
		key:       key,
		value:     value,
		source:    "header",
		optional:  true,
	}
}

// HeaderInt 获取并验证 header 整数参数
func (v *ParamValidator) HeaderInt(key string) *IntParamRule {
	valueStr := v.ctx.GetHeader(key)

	rule := &IntParamRule{
		validator: v,
		key:       key,
		source:    "header",
		optional:  true,
	}

	if valueStr != "" {
		value, err := strconv.Atoi(valueStr)
		rule.value = value
		rule.parseErr = err

		if err != nil {
			v.addError(key, "必须是有效的整数", errors.ValidationInvalidInteger)
		}
	}

	return rule
}

// ==================== 验证结果 ====================

// Check 检查所有验证规则，如果有错误返回 ValidationBadRequest 结果
func (v *ParamValidator) Check() IActionResult {
	if len(v.errors) > 0 {
		return ValidationBadRequest(v.errors)
	}
	return nil
}

// IsValid 检查是否所有验证都通过
func (v *ParamValidator) IsValid() bool {
	return len(v.errors) == 0
}

// Errors 获取所有验证错误
func (v *ParamValidator) Errors() validation.ValidationErrors {
	return v.errors
}

// addError 添加验证错误（使用统一的错误码）
func (v *ParamValidator) addError(field, message, code string) {
	v.errors = append(v.errors, validation.ValidationError{
		Field:   field,
		Message: message,
		Code:    code,
	})
}

// ==================== 字符串参数规则 ====================

// StringParamRule 字符串参数规则
type StringParamRule struct {
	validator *ParamValidator
	key       string
	value     string
	source    string
	optional  bool
}

// Required 必填验证
func (r *StringParamRule) Required() *StringParamRule {
	r.optional = false
	if r.value == "" {
		r.validator.addError(r.key, "不能为空", errors.ValidationRequired)
	}
	return r
}

// Optional 标记为可选参数
func (r *StringParamRule) Optional() *StringParamRule {
	r.optional = true
	return r
}

// NotEmpty 非空验证（与 Required 类似，但语义更清晰）
func (r *StringParamRule) NotEmpty() *StringParamRule {
	return r.Required()
}

// MinLength 最小长度验证
func (r *StringParamRule) MinLength(min int) *StringParamRule {
	if !r.shouldSkip() && len(r.value) < min {
		r.validator.addError(r.key, "长度不能少于 "+strconv.Itoa(min)+" 个字符", errors.ValidationMinLength)
	}
	return r
}

// MaxLength 最大长度验证
func (r *StringParamRule) MaxLength(max int) *StringParamRule {
	if !r.shouldSkip() && len(r.value) > max {
		r.validator.addError(r.key, "长度不能超过 "+strconv.Itoa(max)+" 个字符", errors.ValidationMaxLength)
	}
	return r
}

// Length 长度范围验证
func (r *StringParamRule) Length(min, max int) *StringParamRule {
	if !r.shouldSkip() {
		length := len(r.value)
		if length < min || length > max {
			r.validator.addError(r.key, "长度必须在 "+strconv.Itoa(min)+" 到 "+strconv.Itoa(max)+" 个字符之间", errors.ValidationLength)
		}
	}
	return r
}

// Pattern 正则表达式验证
func (r *StringParamRule) Pattern(pattern, message string) *StringParamRule {
	if !r.shouldSkip() {
		matched, err := regexp.MatchString(pattern, r.value)
		if err != nil || !matched {
			r.validator.addError(r.key, message, errors.ValidationPattern)
		}
	}
	return r
}

// Email 邮箱格式验证
func (r *StringParamRule) Email() *StringParamRule {
	if !r.shouldSkip() {
		emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(emailPattern, r.value)
		if !matched {
			r.validator.addError(r.key, "必须是有效的邮箱地址", errors.ValidationEmail)
		}
	}
	return r
}

// URL 网址格式验证
func (r *StringParamRule) URL() *StringParamRule {
	if !r.shouldSkip() {
		if !strings.HasPrefix(r.value, "http://") && !strings.HasPrefix(r.value, "https://") {
			r.validator.addError(r.key, "必须是有效的 URL 地址", errors.ValidationUrl)
		}
	}
	return r
}

// In 枚举值验证
func (r *StringParamRule) In(values ...string) *StringParamRule {
	if !r.shouldSkip() {
		found := false
		for _, v := range values {
			if r.value == v {
				found = true
				break
			}
		}
		if !found {
			r.validator.addError(r.key, "必须是以下值之一: "+strings.Join(values, ", "), errors.ValidationEnum)
		}
	}
	return r
}

// NotIn 不在枚举值中
func (r *StringParamRule) NotIn(values ...string) *StringParamRule {
	if !r.shouldSkip() {
		for _, v := range values {
			if r.value == v {
				r.validator.addError(r.key, "不能是以下值: "+strings.Join(values, ", "), errors.ValidationNotIn)
				break
			}
		}
	}
	return r
}

// Alpha 只允许字母
func (r *StringParamRule) Alpha() *StringParamRule {
	return r.Pattern(`^[a-zA-Z]+$`, "只能包含字母")
}

// AlphaNumeric 只允许字母和数字
func (r *StringParamRule) AlphaNumeric() *StringParamRule {
	return r.Pattern(`^[a-zA-Z0-9]+$`, "只能包含字母和数字")
}

// Value 获取验证后的值
func (r *StringParamRule) Value() string {
	return r.value
}

// ValueOr 获取值，如果为空则返回默认值
func (r *StringParamRule) ValueOr(defaultValue string) string {
	if r.value == "" {
		return defaultValue
	}
	return r.value
}

// shouldSkip 判断是否应该跳过验证（可选参数且值为空）
func (r *StringParamRule) shouldSkip() bool {
	return r.optional && r.value == ""
}

// ==================== 整数参数规则 ====================

// IntParamRule 整数参数规则
type IntParamRule struct {
	validator *ParamValidator
	key       string
	value     int
	source    string
	optional  bool
	parseErr  error
}

// Required 必填验证
func (r *IntParamRule) Required() *IntParamRule {
	r.optional = false
	return r
}

// Optional 可选参数
func (r *IntParamRule) Optional() *IntParamRule {
	r.optional = true
	return r
}

// Min 最小值验证
func (r *IntParamRule) Min(min int) *IntParamRule {
	if r.parseErr == nil && r.value < min {
		r.validator.addError(r.key, "不能小于 "+strconv.Itoa(min), errors.ValidationMin)
	}
	return r
}

// Max 最大值验证
func (r *IntParamRule) Max(max int) *IntParamRule {
	if r.parseErr == nil && r.value > max {
		r.validator.addError(r.key, "不能大于 "+strconv.Itoa(max), errors.ValidationMax)
	}
	return r
}

// Range 范围验证
func (r *IntParamRule) Range(min, max int) *IntParamRule {
	if r.parseErr == nil && (r.value < min || r.value > max) {
		r.validator.addError(r.key, "必须在 "+strconv.Itoa(min)+" 到 "+strconv.Itoa(max)+" 之间", errors.ValidationRange)
	}
	return r
}

// Positive 正数验证（大于 0）
func (r *IntParamRule) Positive() *IntParamRule {
	if r.parseErr == nil && r.value <= 0 {
		r.validator.addError(r.key, "必须是正数", errors.ValidationPositive)
	}
	return r
}

// NonNegative 非负数验证（大于等于 0）
func (r *IntParamRule) NonNegative() *IntParamRule {
	if r.parseErr == nil && r.value < 0 {
		r.validator.addError(r.key, "不能是负数", errors.ValidationNonNegative)
	}
	return r
}

// In 枚举值验证
func (r *IntParamRule) In(values ...int) *IntParamRule {
	if r.parseErr == nil {
		found := false
		for _, v := range values {
			if r.value == v {
				found = true
				break
			}
		}
		if !found {
			strValues := make([]string, len(values))
			for i, v := range values {
				strValues[i] = strconv.Itoa(v)
			}
			r.validator.addError(r.key, "必须是以下值之一: "+strings.Join(strValues, ", "), errors.ValidationEnum)
		}
	}
	return r
}

// Value 获取验证后的值
func (r *IntParamRule) Value() int {
	return r.value
}

// ValueOr 获取值，如果解析失败或为空则返回默认值
func (r *IntParamRule) ValueOr(defaultValue int) int {
	if r.parseErr != nil {
		return defaultValue
	}
	return r.value
}

// ==================== Int64 参数规则 ====================

// Int64ParamRule int64 参数规则
type Int64ParamRule struct {
	validator *ParamValidator
	key       string
	value     int64
	source    string
	optional  bool
	parseErr  error
}

// Required 必填验证
func (r *Int64ParamRule) Required() *Int64ParamRule {
	r.optional = false
	return r
}

// Optional 可选参数
func (r *Int64ParamRule) Optional() *Int64ParamRule {
	r.optional = true
	return r
}

// Min 最小值验证
func (r *Int64ParamRule) Min(min int64) *Int64ParamRule {
	if r.parseErr == nil && r.value < min {
		r.validator.addError(r.key, "不能小于 "+strconv.FormatInt(min, 10), errors.ValidationMin)
	}
	return r
}

// Max 最大值验证
func (r *Int64ParamRule) Max(max int64) *Int64ParamRule {
	if r.parseErr == nil && r.value > max {
		r.validator.addError(r.key, "不能大于 "+strconv.FormatInt(max, 10), errors.ValidationMax)
	}
	return r
}

// Range 范围验证
func (r *Int64ParamRule) Range(min, max int64) *Int64ParamRule {
	if r.parseErr == nil && (r.value < min || r.value > max) {
		r.validator.addError(r.key, "必须在 "+strconv.FormatInt(min, 10)+" 到 "+strconv.FormatInt(max, 10)+" 之间", errors.ValidationRange)
	}
	return r
}

// Positive 正数验证
func (r *Int64ParamRule) Positive() *Int64ParamRule {
	if r.parseErr == nil && r.value <= 0 {
		r.validator.addError(r.key, "必须是正数", errors.ValidationPositive)
	}
	return r
}

// NonNegative 非负数验证
func (r *Int64ParamRule) NonNegative() *Int64ParamRule {
	if r.parseErr == nil && r.value < 0 {
		r.validator.addError(r.key, "不能是负数", errors.ValidationNonNegative)
	}
	return r
}

// Value 获取验证后的值
func (r *Int64ParamRule) Value() int64 {
	return r.value
}

// ValueOr 获取值，如果解析失败则返回默认值
func (r *Int64ParamRule) ValueOr(defaultValue int64) int64 {
	if r.parseErr != nil {
		return defaultValue
	}
	return r.value
}

// ==================== 布尔参数规则 ====================

// BoolParamRule 布尔参数规则
type BoolParamRule struct {
	validator *ParamValidator
	key       string
	value     bool
	source    string
	optional  bool
	parseErr  error
}

// Required 必填验证
func (r *BoolParamRule) Required() *BoolParamRule {
	r.optional = false
	return r
}

// Optional 可选参数
func (r *BoolParamRule) Optional() *BoolParamRule {
	r.optional = true
	return r
}

// Value 获取验证后的值
func (r *BoolParamRule) Value() bool {
	return r.value
}

// ValueOr 获取值，如果解析失败则返回默认值
func (r *BoolParamRule) ValueOr(defaultValue bool) bool {
	if r.parseErr != nil {
		return defaultValue
	}
	return r.value
}

// ==================== 浮点数参数规则 ====================

// FloatParamRule 浮点数参数规则
type FloatParamRule struct {
	validator *ParamValidator
	key       string
	value     float64
	source    string
	optional  bool
	parseErr  error
}

// Required 必填验证
func (r *FloatParamRule) Required() *FloatParamRule {
	r.optional = false
	return r
}

// Optional 可选参数
func (r *FloatParamRule) Optional() *FloatParamRule {
	r.optional = true
	return r
}

// Min 最小值验证
func (r *FloatParamRule) Min(min float64) *FloatParamRule {
	if r.parseErr == nil && r.value < min {
		r.validator.addError(r.key, "不能小于 "+strconv.FormatFloat(min, 'f', -1, 64), errors.ValidationMin)
	}
	return r
}

// Max 最大值验证
func (r *FloatParamRule) Max(max float64) *FloatParamRule {
	if r.parseErr == nil && r.value > max {
		r.validator.addError(r.key, "不能大于 "+strconv.FormatFloat(max, 'f', -1, 64), errors.ValidationMax)
	}
	return r
}

// Range 范围验证
func (r *FloatParamRule) Range(min, max float64) *FloatParamRule {
	if r.parseErr == nil && (r.value < min || r.value > max) {
		r.validator.addError(r.key, "必须在 "+strconv.FormatFloat(min, 'f', -1, 64)+" 到 "+strconv.FormatFloat(max, 'f', -1, 64)+" 之间", errors.ValidationRange)
	}
	return r
}

// Positive 正数验证
func (r *FloatParamRule) Positive() *FloatParamRule {
	if r.parseErr == nil && r.value <= 0 {
		r.validator.addError(r.key, "必须是正数", errors.ValidationPositive)
	}
	return r
}

// Value 获取验证后的值
func (r *FloatParamRule) Value() float64 {
	return r.value
}

// ValueOr 获取值，如果解析失败则返回默认值
func (r *FloatParamRule) ValueOr(defaultValue float64) float64 {
	if r.parseErr != nil {
		return defaultValue
	}
	return r.value
}
