package web

import (
	"context"
	"io"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/errors"
	"github.com/gocrud/csgo/validation"
)

type HttpContext struct {
	gin *gin.Context

	// Services 提供对应用程序 DI 容器的访问。
	// 使用 di.Get[T](ctx.Services) 来解析服务。
	Services di.IServiceProvider

	// paramErrors 存储参数验证错误，供新的泛型参数 API 使用。
	// 使用 web.Path[T], web.Query[T], web.Header[T] 等方法时，
	// 验证错误会自动收集到这里，并在处理器结束时统一返回。
	paramErrors validation.ValidationErrors
}

// NewHttpContext 从 gin.Context 创建新的 HttpContext。
func NewHttpContext(c *gin.Context) *HttpContext {
	return &HttpContext{gin: c}
}

// RawCtx 返回底层的 gin.Context。
func (c *HttpContext) RawCtx() *gin.Context {
	return c.gin
}

// Context 返回请求的上下文。
func (c *HttpContext) Context() context.Context {
	return c.gin.Request.Context()
}

// ==================== 成功响应 ====================

// Ok 返回 200 OK 及数据。
func (c *HttpContext) Ok(data interface{}) IActionResult {
	return Ok(data)
}

// Created 返回 201 Created 及数据。
func (c *HttpContext) Created(data interface{}) IActionResult {
	return Created(data)
}

// NoContent 返回 204 No Content。
func (c *HttpContext) NoContent() IActionResult {
	return NoContent()
}

// ==================== 错误响应 ====================

// BadRequest 返回 400 Bad Request。
func (c *HttpContext) BadRequest(message string) IActionResult {
	return BadRequest(message)
}

// BadRequestWithCode 返回 400 Bad Request，带有自定义错误码。
func (c *HttpContext) BadRequestWithCode(code, message string) IActionResult {
	return BadRequestWithCode(code, message)
}

// Unauthorized 返回 401 Unauthorized。
func (c *HttpContext) Unauthorized(message string) IActionResult {
	return Unauthorized(message)
}

// Forbidden 返回 403 Forbidden。
func (c *HttpContext) Forbidden(message string) IActionResult {
	return Forbidden(message)
}

// NotFound 返回 404 Not Found。
func (c *HttpContext) NotFound(message string) IActionResult {
	return NotFound(message)
}

// Conflict 返回 409 Conflict。
func (c *HttpContext) Conflict(message string) IActionResult {
	return Conflict(message)
}

// InternalError 返回 500 Internal Server Error。
func (c *HttpContext) InternalError(message string) IActionResult {
	return InternalError(message)
}

// Error 返回自定义错误响应。
func (c *HttpContext) Error(statusCode int, code, message string) IActionResult {
	return Error(statusCode, code, message)
}

// ValidationBadRequest 返回 400 Bad Request，带有验证错误。
func (c *HttpContext) ValidationBadRequest(errs validation.ValidationErrors) IActionResult {
	return ValidationBadRequest(errs)
}

// ValidationBadRequestWithCode 返回 400 Bad Request，带有验证错误和自定义错误码。
func (c *HttpContext) ValidationBadRequestWithCode(code string, errs validation.ValidationErrors) IActionResult {
	return ValidationBadRequestWithCode(code, errs)
}

// FrameworkError 返回框架错误，自动映射 HTTP 状态码。
func (c *HttpContext) FrameworkError(err *errors.Error) IActionResult {
	return FrameworkError(err)
}

// FrameworkErrorWithStatus 返回框架错误，带有指定的 HTTP 状态码。
func (c *HttpContext) FrameworkErrorWithStatus(statusCode int, err *errors.Error) IActionResult {
	return FrameworkErrorWithStatus(statusCode, err)
}

// ==================== 绑定辅助方法 ====================

// cleanJSONErrorMessage 清理 JSON 解析错误消息中的字段路径
// 移除字段路径前面的点和结构体名称，提供更清晰的错误提示
func cleanJSONErrorMessage(errMsg string) string {
	// 匹配 "Go struct field .field" 或 "Go struct field StructName.field"
	// 并替换为 "Go struct field field"
	// 处理前导点：".info.nickname" -> "info.nickname"
	// 处理结构体名称："User.name" -> "name"
	re := regexp.MustCompile(`Go struct field\s+(\w+\.|\.)`)
	cleaned := re.ReplaceAllString(errMsg, "Go struct field ")

	return cleaned
}

// BindJSON 将 JSON 请求体绑定到目标对象，失败时返回 BadRequest。
// 绑定成功时返回 true，否则返回 false。
func (c *HttpContext) BindJSON(target interface{}) (ok bool, result IActionResult) {
	if err := c.gin.ShouldBindJSON(target); err != nil {
		// 检查是否为 EOF 错误(空请求体)
		if err == io.EOF || err.Error() == "EOF" {
			return false, c.BadRequest("请求体不能为空,请提供有效的 JSON 数据")
		}

		// 检查是否为不完整的 JSON
		errMsg := err.Error()
		if strings.Contains(errMsg, "unexpected end of JSON input") ||
			strings.Contains(errMsg, "unexpected EOF") {
			return false, c.BadRequest("请求体格式不完整,请提供完整的 JSON 数据")
		}

		// 其他 JSON 解析错误,提供更友好的错误提示
		if strings.Contains(errMsg, "invalid character") ||
			strings.Contains(errMsg, "cannot unmarshal") {
			cleanedMsg := cleanJSONErrorMessage(errMsg)
			return false, c.BadRequest("JSON 格式错误: " + cleanedMsg)
		}

		// 未知错误,直接返回原始错误信息
		return false, c.BadRequest(err.Error())
	}
	return true, nil
}

// MustBindJSON 将 JSON 请求体绑定到目标对象，失败时返回 BadRequest。
// 这是一个便捷方法，仅返回错误结果。
func (c *HttpContext) MustBindJSON(target interface{}) IActionResult {
	if err := c.gin.ShouldBindJSON(target); err != nil {
		// 检查是否为 EOF 错误(空请求体)
		if err == io.EOF || err.Error() == "EOF" {
			return c.BadRequest("请求体不能为空,请提供有效的 JSON 数据")
		}

		// 检查是否为不完整的 JSON
		errMsg := err.Error()
		if strings.Contains(errMsg, "unexpected end of JSON input") ||
			strings.Contains(errMsg, "unexpected EOF") {
			return c.BadRequest("请求体格式不完整,请提供完整的 JSON 数据")
		}

		// 其他 JSON 解析错误,提供更友好的错误提示
		if strings.Contains(errMsg, "invalid character") ||
			strings.Contains(errMsg, "cannot unmarshal") {
			cleanedMsg := cleanJSONErrorMessage(errMsg)
			return c.BadRequest("JSON 格式错误: " + cleanedMsg)
		}

		// 未知错误,直接返回原始错误信息
		return c.BadRequest(err.Error())
	}
	return nil
}

// BindQuery 将查询参数绑定到目标对象，失败时返回 BadRequest。
func (c *HttpContext) BindQuery(target interface{}) (ok bool, result IActionResult) {
	if err := c.gin.ShouldBindQuery(target); err != nil {
		return false, c.BadRequest(err.Error())
	}
	return true, nil
}

// ==================== 验证 ====================

// handleBindError 处理绑定错误,返回友好的错误信息
func handleBindError(c *HttpContext, err error) IActionResult {
	// 检查是否为 EOF 错误(空请求体)
	if err == io.EOF || err.Error() == "EOF" {
		return c.BadRequest("请求体不能为空,请提供有效的 JSON 数据")
	}

	// 检查是否为不完整的 JSON
	errMsg := err.Error()
	if strings.Contains(errMsg, "unexpected end of JSON input") ||
		strings.Contains(errMsg, "unexpected EOF") {
		return c.BadRequest("请求体格式不完整,请提供完整的 JSON 数据")
	}

	// 其他 JSON 解析错误,提供更友好的错误提示
	if strings.Contains(errMsg, "invalid character") ||
		strings.Contains(errMsg, "cannot unmarshal") {
		cleanedMsg := cleanJSONErrorMessage(errMsg)
		return c.BadRequest("JSON 格式错误: " + cleanedMsg)
	}

	// 未知错误,直接返回原始错误信息
	return c.BadRequest(err.Error())
}

// shouldBindAndValidate 通用的绑定和验证逻辑
// bindFunc: 执行实际绑定操作的函数
func shouldBindAndValidate[T any](c *HttpContext, bindFunc func(any) error) (*T, IActionResult) {
	var target T

	// 1. 执行绑定操作
	if err := bindFunc(&target); err != nil {
		return nil, handleBindError(c, err)
	}

	// 2. 使用新验证器执行验证
	errs := validation.Validate(&target)
	if errs != nil && errs.HasErrors() {
		// 返回结构化的验证错误
		return nil, c.ValidationBadRequest(errs)
	}

	// 3. 如果实现了 Validator 接口，也调用
	if v, ok := any(&target).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, c.BadRequest(err.Error())
		}
	}

	return &target, nil
}

// ShouldBindJSON 绑定 JSON 并验证
// 自动使用注册验证器的模式(快速失败或全量验证)
func ShouldBindJSON[T any](c *HttpContext) (*T, IActionResult) {
	return shouldBindAndValidate[T](c, c.gin.ShouldBindJSON)
}

// ShouldBindHeader 绑定 Header 并验证
func ShouldBindHeader[T any](c *HttpContext) (*T, IActionResult) {
	return shouldBindAndValidate[T](c, c.gin.ShouldBindHeader)
}

// ShouldBindQuery 绑定 Query 参数并验证
func ShouldBindQuery[T any](c *HttpContext) (*T, IActionResult) {
	return shouldBindAndValidate[T](c, c.gin.ShouldBindQuery)
}

// ShouldBindPlain 绑定 Plain 文本并验证
func ShouldBindPlain[T any](c *HttpContext) (*T, IActionResult) {
	return shouldBindAndValidate[T](c, c.gin.ShouldBindPlain)
}

// ShouldBindUri 绑定 URI 参数并验证
func ShouldBindUri[T any](c *HttpContext) (*T, IActionResult) {
	return shouldBindAndValidate[T](c, c.gin.ShouldBindUri)
}

// ShouldBindXML 绑定 XML 并验证
func ShouldBindXML[T any](c *HttpContext) (*T, IActionResult) {
	return shouldBindAndValidate[T](c, c.gin.ShouldBindXML)
}

// ShouldBindYAML 绑定 YAML 并验证
func ShouldBindYAML[T any](c *HttpContext) (*T, IActionResult) {
	return shouldBindAndValidate[T](c, c.gin.ShouldBindYAML)
}

// ShouldBindTOML 绑定 TOML 并验证
func ShouldBindTOML[T any](c *HttpContext) (*T, IActionResult) {
	return shouldBindAndValidate[T](c, c.gin.ShouldBindTOML)
}

// ==================== 智能错误处理 ====================

// FromError 智能处理错误，是 web.FromError 的便捷方法。
// 自动识别错误类型并返回对应的 ActionResult。
//
// 使用示例：
//
//	user, err := service.GetUser(id)
//	if err != nil {
//	    return c.FromError(err, "获取用户失败")
//	}
func (c *HttpContext) FromError(err error, defaultMessage ...string) IActionResult {
	return FromError(err, defaultMessage...)
}

// FromErrorWithStatus 智能处理错误并指定状态码。
// 对于普通 error 使用指定的状态码，对于 BizError 和 ValidationErrors 忽略状态码。
//
// 使用示例：
//
//	err := db.Connect()
//	if err != nil {
//	    return c.FromErrorWithStatus(err, 503, "数据库服务暂时不可用")
//	}
func (c *HttpContext) FromErrorWithStatus(err error, statusCode int, defaultMessage ...string) IActionResult {
	return FromErrorWithStatus(err, statusCode, defaultMessage...)
}

// ==================== 参数验证错误收集（内部使用）====================

// addParamError 添加参数验证错误（供泛型参数 API 内部使用）。
func (c *HttpContext) addParamError(err validation.ValidationError) {
	c.paramErrors = append(c.paramErrors, err)
}

// HasParamErrors 检查是否有参数验证错误。
func (c *HttpContext) HasParamErrors() bool {
	return len(c.paramErrors) > 0
}

// GetParamErrors 获取所有参数验证错误。
func (c *HttpContext) GetParamErrors() validation.ValidationErrors {
	return c.paramErrors
}
