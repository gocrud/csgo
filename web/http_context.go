package web

import (
	"context"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/errors"
	"github.com/gocrud/csgo/validation"
)

// HttpContext wraps gin.Context and provides unified API for HTTP handling.
// Access the underlying gin.Context via RawCtx() method.
type HttpContext struct {
	gin *gin.Context
}

// NewHttpContext creates a new HttpContext from gin.Context.
func NewHttpContext(c *gin.Context) *HttpContext {
	return &HttpContext{gin: c}
}

// RawCtx returns the underlying gin.Context.
func (c *HttpContext) RawCtx() *gin.Context {
	return c.gin
}

// Context returns the request's context.
func (c *HttpContext) Context() context.Context {
	return c.gin.Request.Context()
}

// ==================== Success Responses ====================

// Ok returns 200 OK with data.
func (c *HttpContext) Ok(data interface{}) IActionResult {
	return Ok(data)
}

// Created returns 201 Created with data.
func (c *HttpContext) Created(data interface{}) IActionResult {
	return Created(data)
}

// NoContent returns 204 No Content.
func (c *HttpContext) NoContent() IActionResult {
	return NoContent()
}

// ==================== Error Responses ====================

// BadRequest returns 400 Bad Request.
func (c *HttpContext) BadRequest(message string) IActionResult {
	return BadRequest(message)
}

// BadRequestWithCode returns 400 Bad Request with custom error code.
func (c *HttpContext) BadRequestWithCode(code, message string) IActionResult {
	return BadRequestWithCode(code, message)
}

// Unauthorized returns 401 Unauthorized.
func (c *HttpContext) Unauthorized(message string) IActionResult {
	return Unauthorized(message)
}

// Forbidden returns 403 Forbidden.
func (c *HttpContext) Forbidden(message string) IActionResult {
	return Forbidden(message)
}

// NotFound returns 404 Not Found.
func (c *HttpContext) NotFound(message string) IActionResult {
	return NotFound(message)
}

// Conflict returns 409 Conflict.
func (c *HttpContext) Conflict(message string) IActionResult {
	return Conflict(message)
}

// InternalError returns 500 Internal Server Error.
func (c *HttpContext) InternalError(message string) IActionResult {
	return InternalError(message)
}

// Error returns a custom error response.
func (c *HttpContext) Error(statusCode int, code, message string) IActionResult {
	return Error(statusCode, code, message)
}

// ValidationBadRequest returns 400 Bad Request with validation errors.
func (c *HttpContext) ValidationBadRequest(errs validation.ValidationErrors) IActionResult {
	return ValidationBadRequest(errs)
}

// ValidationBadRequestWithCode returns 400 Bad Request with validation errors and custom code.
func (c *HttpContext) ValidationBadRequestWithCode(code string, errs validation.ValidationErrors) IActionResult {
	return ValidationBadRequestWithCode(code, errs)
}

// BizError returns a business error with auto-mapped HTTP status code.
func (c *HttpContext) BizError(err *errors.BizError) IActionResult {
	return BizError(err)
}

// BizErrorWithStatus returns a business error with specified HTTP status code.
func (c *HttpContext) BizErrorWithStatus(statusCode int, err *errors.BizError) IActionResult {
	return BizErrorWithStatus(statusCode, err)
}

// ==================== Binding Helpers ====================

// BindJSON binds JSON body to target and returns BadRequest if failed.
// Returns true if binding succeeded, false otherwise.
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
			return false, c.BadRequest("JSON 格式错误: " + errMsg)
		}

		// 未知错误,直接返回原始错误信息
		return false, c.BadRequest(err.Error())
	}
	return true, nil
}

// MustBindJSON binds JSON body to target and returns BadRequest if failed.
// This is a convenience method that returns only the error result.
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
			return c.BadRequest("JSON 格式错误: " + errMsg)
		}

		// 未知错误,直接返回原始错误信息
		return c.BadRequest(err.Error())
	}
	return nil
}

// BindQuery binds query parameters to target and returns BadRequest if failed.
func (c *HttpContext) BindQuery(target interface{}) (ok bool, result IActionResult) {
	if err := c.gin.ShouldBindQuery(target); err != nil {
		return false, c.BadRequest(err.Error())
	}
	return true, nil
}

// ==================== Validation ====================

// BindAndValidate binds JSON body and validates using FluentValidation validator.
// Returns the bound object and nil if successful, or nil and an error result if failed.
// 自动使用注册验证器的模式(快速失败或全量验证)
func BindAndValidate[T any](c *HttpContext) (*T, IActionResult) {
	var target T

	// 1. 绑定 JSON,增强错误处理
	if err := c.gin.ShouldBindJSON(&target); err != nil {
		// 检查是否为 EOF 错误(空请求体)
		if err == io.EOF || err.Error() == "EOF" {
			return nil, c.BadRequest("请求体不能为空,请提供有效的 JSON 数据")
		}

		// 检查是否为不完整的 JSON
		errMsg := err.Error()
		if strings.Contains(errMsg, "unexpected end of JSON input") ||
			strings.Contains(errMsg, "unexpected EOF") {
			return nil, c.BadRequest("请求体格式不完整,请提供完整的 JSON 数据")
		}

		// 其他 JSON 解析错误,提供更友好的错误提示
		if strings.Contains(errMsg, "invalid character") ||
			strings.Contains(errMsg, "cannot unmarshal") {
			return nil, c.BadRequest("JSON 格式错误: " + errMsg)
		}

		// 未知错误,直接返回原始错误信息
		return nil, c.BadRequest(err.Error())
	}

	// 2. 查找注册的验证器并执行验证
	if validator, ok := validation.GetValidator[T](); ok {
		result := validator.Validate(&target)
		if !result.IsValid {
			// 返回结构化的验证错误
			return nil, c.ValidationBadRequest(result.Errors)
		}
	}

	// 3. 如果实现了 Validator 接口，也调用
	if v, ok := any(&target).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return nil, c.BadRequest(err.Error())
		}
	}

	return &target, nil
}

// ValidateStruct validates a struct using registered FluentValidation validator.
// Note: This method requires the target to be passed as a pointer and
// a validator must be registered for the type.
func (c *HttpContext) ValidateStruct(target interface{}) IActionResult {
	// 由于类型推断问题，这个方法暂时不实现通用版本
	// 建议直接使用 BindAndValidate[T] 泛型方法
	return nil
}
