package web

import (
	"github.com/gin-gonic/gin"
)

// IActionResult represents the result of an action method.
// Similar to .NET's IActionResult interface.
type IActionResult interface {
	// ExecuteResult writes the result to the response.
	ExecuteResult(c *gin.Context)
}

// ApiResponse is the standard API response format.
type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ApiError   `json:"error,omitempty"`
}

// ApiError represents an error in the API response.
type ApiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ==================== Success Results ====================

// OkResult represents a 200 OK response.
type OkResult struct {
	Data interface{}
}

// ExecuteResult implements IActionResult.
func (r OkResult) ExecuteResult(c *gin.Context) {
	c.JSON(200, ApiResponse{Success: true, Data: r.Data})
}

// Ok creates a 200 OK result with data.
func Ok(data interface{}) IActionResult {
	return OkResult{Data: data}
}

// CreatedResult represents a 201 Created response.
type CreatedResult struct {
	Data interface{}
}

// ExecuteResult implements IActionResult.
func (r CreatedResult) ExecuteResult(c *gin.Context) {
	c.JSON(201, ApiResponse{Success: true, Data: r.Data})
}

// Created creates a 201 Created result with data.
func Created(data interface{}) IActionResult {
	return CreatedResult{Data: data}
}

// NoContentResult represents a 204 No Content response.
type NoContentResult struct{}

// ExecuteResult implements IActionResult.
func (r NoContentResult) ExecuteResult(c *gin.Context) {
	c.Status(204)
}

// NoContent creates a 204 No Content result.
func NoContent() IActionResult {
	return NoContentResult{}
}

// ==================== Redirect Results ====================

// RedirectResult represents a redirect response.
type RedirectResult struct {
	StatusCode int
	Location   string
}

// ExecuteResult implements IActionResult.
func (r RedirectResult) ExecuteResult(c *gin.Context) {
	c.Redirect(r.StatusCode, r.Location)
}

// Redirect creates a 302 Found redirect result.
func Redirect(location string) IActionResult {
	return RedirectResult{StatusCode: 302, Location: location}
}

// RedirectPermanent creates a 301 Moved Permanently redirect result.
func RedirectPermanent(location string) IActionResult {
	return RedirectResult{StatusCode: 301, Location: location}
}

// ==================== Error Results ====================

// ErrorResult represents an error response.
type ErrorResult struct {
	StatusCode int
	Code       string
	Message    string
}

// ExecuteResult implements IActionResult.
func (r ErrorResult) ExecuteResult(c *gin.Context) {
	c.JSON(r.StatusCode, ApiResponse{
		Success: false,
		Error:   &ApiError{Code: r.Code, Message: r.Message},
	})
}

// Error creates a custom error result.
func Error(statusCode int, code, message string) IActionResult {
	return ErrorResult{StatusCode: statusCode, Code: code, Message: message}
}

// BadRequest creates a 400 Bad Request result.
func BadRequest(message string) IActionResult {
	return ErrorResult{StatusCode: 400, Code: "BAD_REQUEST", Message: message}
}

// BadRequestWithCode creates a 400 Bad Request result with custom code.
func BadRequestWithCode(code, message string) IActionResult {
	return ErrorResult{StatusCode: 400, Code: code, Message: message}
}

// Unauthorized creates a 401 Unauthorized result.
func Unauthorized(message string) IActionResult {
	return ErrorResult{StatusCode: 401, Code: "UNAUTHORIZED", Message: message}
}

// Forbidden creates a 403 Forbidden result.
func Forbidden(message string) IActionResult {
	return ErrorResult{StatusCode: 403, Code: "FORBIDDEN", Message: message}
}

// NotFound creates a 404 Not Found result.
func NotFound(message string) IActionResult {
	return ErrorResult{StatusCode: 404, Code: "NOT_FOUND", Message: message}
}

// Conflict creates a 409 Conflict result.
func Conflict(message string) IActionResult {
	return ErrorResult{StatusCode: 409, Code: "CONFLICT", Message: message}
}

// InternalError creates a 500 Internal Server Error result.
func InternalError(message string) IActionResult {
	return ErrorResult{StatusCode: 500, Code: "INTERNAL_ERROR", Message: message}
}

// ==================== JSON Result ====================

// JsonResult represents a custom JSON response.
type JsonResult struct {
	StatusCode int
	Data       interface{}
}

// ExecuteResult implements IActionResult.
func (r JsonResult) ExecuteResult(c *gin.Context) {
	c.JSON(r.StatusCode, r.Data)
}

// Json creates a custom JSON result.
func Json(statusCode int, data interface{}) IActionResult {
	return JsonResult{StatusCode: statusCode, Data: data}
}

// ==================== Content Result ====================

// ContentResult represents a plain text response.
type ContentResult struct {
	StatusCode  int
	Content     string
	ContentType string
}

// ExecuteResult implements IActionResult.
func (r ContentResult) ExecuteResult(c *gin.Context) {
	if r.ContentType != "" {
		c.Data(r.StatusCode, r.ContentType, []byte(r.Content))
	} else {
		c.String(r.StatusCode, r.Content)
	}
}

// Content creates a plain text result.
func Content(statusCode int, content string) IActionResult {
	return ContentResult{StatusCode: statusCode, Content: content}
}

// ContentWithType creates a content result with custom content type.
func ContentWithType(statusCode int, content, contentType string) IActionResult {
	return ContentResult{StatusCode: statusCode, Content: content, ContentType: contentType}
}

// ==================== File Result ====================

// FileResult represents a file download response.
type FileResult struct {
	FilePath string
	FileName string
}

// ExecuteResult implements IActionResult.
func (r FileResult) ExecuteResult(c *gin.Context) {
	if r.FileName != "" {
		c.FileAttachment(r.FilePath, r.FileName)
	} else {
		c.File(r.FilePath)
	}
}

// File creates a file result.
func File(filePath string) IActionResult {
	return FileResult{FilePath: filePath}
}

// FileDownload creates a file download result with custom filename.
func FileDownload(filePath, fileName string) IActionResult {
	return FileResult{FilePath: filePath, FileName: fileName}
}

// ==================== Status Result ====================

// StatusResult represents a response with only status code.
type StatusResult struct {
	StatusCode int
}

// ExecuteResult implements IActionResult.
func (r StatusResult) ExecuteResult(c *gin.Context) {
	c.Status(r.StatusCode)
}

// Status creates a status-only result.
func Status(statusCode int) IActionResult {
	return StatusResult{StatusCode: statusCode}
}
