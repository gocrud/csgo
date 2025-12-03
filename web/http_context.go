package web

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// HttpContext wraps gin.Context and provides unified API for HTTP handling.
// It embeds gin.Context so all original methods are still available.
type HttpContext struct {
	*gin.Context
}

// NewHttpContext creates a new HttpContext from gin.Context.
func NewHttpContext(c *gin.Context) *HttpContext {
	return &HttpContext{Context: c}
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

// ==================== Parameter Helpers ====================

// PathInt gets path parameter and converts to int.
// Returns error if conversion fails.
func (c *HttpContext) PathInt(key string) (int, error) {
	return strconv.Atoi(c.Param(key))
}

// PathInt64 gets path parameter and converts to int64.
func (c *HttpContext) PathInt64(key string) (int64, error) {
	return strconv.ParseInt(c.Param(key), 10, 64)
}

// MustPathInt gets path parameter and converts to int.
// Returns BadRequest result if conversion fails.
func (c *HttpContext) MustPathInt(key string) (int, IActionResult) {
	val, err := strconv.Atoi(c.Param(key))
	if err != nil {
		return 0, c.BadRequest("Invalid " + key + ": must be an integer")
	}
	return val, nil
}

// QueryInt gets query parameter and converts to int with default value.
func (c *HttpContext) QueryInt(key string, defaultValue int) int {
	if val := c.Query(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultValue
}

// QueryInt64 gets query parameter and converts to int64 with default value.
func (c *HttpContext) QueryInt64(key string, defaultValue int64) int64 {
	if val := c.Query(key); val != "" {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i
		}
	}
	return defaultValue
}

// QueryBool gets query parameter and converts to bool with default value.
func (c *HttpContext) QueryBool(key string, defaultValue bool) bool {
	if val := c.Query(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultValue
}

// ==================== Binding Helpers ====================

// BindJSON binds JSON body to target and returns BadRequest if failed.
// Returns true if binding succeeded, false otherwise.
func (c *HttpContext) BindJSON(target interface{}) (ok bool, result IActionResult) {
	if err := c.ShouldBindJSON(target); err != nil {
		return false, c.BadRequest(err.Error())
	}
	return true, nil
}

// MustBindJSON binds JSON body to target and returns BadRequest if failed.
// This is a convenience method that returns only the error result.
func (c *HttpContext) MustBindJSON(target interface{}) IActionResult {
	if err := c.ShouldBindJSON(target); err != nil {
		return c.BadRequest(err.Error())
	}
	return nil
}

// BindQuery binds query parameters to target and returns BadRequest if failed.
func (c *HttpContext) BindQuery(target interface{}) (ok bool, result IActionResult) {
	if err := c.ShouldBindQuery(target); err != nil {
		return false, c.BadRequest(err.Error())
	}
	return true, nil
}

