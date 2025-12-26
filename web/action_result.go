package web

import (
	"encoding/base64"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/errors"
	"github.com/gocrud/csgo/validation"
)

// IActionResult 表示操作方法的结果。
// 类似于 .NET 的 IActionResult 接口。
type IActionResult interface {
	// ExecuteResult 将结果写入响应。
	ExecuteResult(c *gin.Context)
}

// ApiResponse 是标准的 API 响应格式。
type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ApiError   `json:"error,omitempty"`
}

// ApiError 表示 API 响应中的错误。
type ApiError struct {
	Code    string                       `json:"code"`              // 错误码
	Message string                       `json:"message"`           // 错误消息
	Fields  []validation.ValidationError `json:"fields,omitempty"`  // 验证错误字段列表
	Details M                            `json:"details,omitempty"` // 额外详情（可选）
}

// ==================== 成功结果 ====================

// OkResult 表示 200 OK 响应。
type OkResult struct {
	Data interface{}
}

// ExecuteResult 实现 IActionResult 接口。
func (r OkResult) ExecuteResult(c *gin.Context) {
	c.JSON(200, ApiResponse{Success: true, Data: r.Data})
}

// Ok 创建 200 OK 结果及数据。
func Ok(data interface{}) IActionResult {
	return OkResult{Data: data}
}

// CreatedResult 表示 201 Created 响应。
type CreatedResult struct {
	Data interface{}
}

// ExecuteResult 实现 IActionResult 接口。
func (r CreatedResult) ExecuteResult(c *gin.Context) {
	c.JSON(201, ApiResponse{Success: true, Data: r.Data})
}

// Created 创建 201 Created 结果及数据。
func Created(data interface{}) IActionResult {
	return CreatedResult{Data: data}
}

// NoContentResult 表示 204 No Content 响应。
type NoContentResult struct{}

// ExecuteResult 实现 IActionResult 接口。
func (r NoContentResult) ExecuteResult(c *gin.Context) {
	c.Status(204)
}

// NoContent 创建 204 No Content 结果。
func NoContent() IActionResult {
	return NoContentResult{}
}

// ==================== 重定向结果 ====================

// RedirectResult 表示重定向响应。
type RedirectResult struct {
	StatusCode int
	Location   string
}

// ExecuteResult 实现 IActionResult 接口。
func (r RedirectResult) ExecuteResult(c *gin.Context) {
	c.Redirect(r.StatusCode, r.Location)
}

// Redirect 创建 302 Found 重定向结果。
func Redirect(location string) IActionResult {
	return RedirectResult{StatusCode: 302, Location: location}
}

// RedirectPermanent 创建 301 Moved Permanently 重定向结果。
func RedirectPermanent(location string) IActionResult {
	return RedirectResult{StatusCode: 301, Location: location}
}

// ==================== 错误结果 ====================

// ErrorResult 表示错误响应。
type ErrorResult struct {
	StatusCode int
	Code       string
	Message    string
}

// ExecuteResult 实现 IActionResult 接口。
func (r ErrorResult) ExecuteResult(c *gin.Context) {
	c.JSON(r.StatusCode, ApiResponse{
		Success: false,
		Error:   &ApiError{Code: r.Code, Message: r.Message},
	})
}

// Error 创建自定义错误结果。
func Error(statusCode int, code, message string) IActionResult {
	return ErrorResult{StatusCode: statusCode, Code: code, Message: message}
}

// BadRequest 创建 400 Bad Request 结果。
func BadRequest(message string) IActionResult {
	return ErrorResult{StatusCode: 400, Code: "BAD_REQUEST", Message: message}
}

// BadRequestWithCode 创建 400 Bad Request 结果，带有自定义错误码。
func BadRequestWithCode(code, message string) IActionResult {
	return ErrorResult{StatusCode: 400, Code: code, Message: message}
}

// Unauthorized 创建 401 Unauthorized 结果。
func Unauthorized(message string) IActionResult {
	return ErrorResult{StatusCode: 401, Code: "UNAUTHORIZED", Message: message}
}

// Forbidden 创建 403 Forbidden 结果。
func Forbidden(message string) IActionResult {
	return ErrorResult{StatusCode: 403, Code: "FORBIDDEN", Message: message}
}

// NotFound 创建 404 Not Found 结果。
func NotFound(message string) IActionResult {
	return ErrorResult{StatusCode: 404, Code: "NOT_FOUND", Message: message}
}

// Conflict 创建 409 Conflict 结果。
func Conflict(message string) IActionResult {
	return ErrorResult{StatusCode: 409, Code: "CONFLICT", Message: message}
}

// InternalError 创建 500 Internal Server Error 结果。
func InternalError(message string) IActionResult {
	return ErrorResult{StatusCode: 500, Code: "INTERNAL_ERROR", Message: message}
}

// ==================== 验证错误结果 ====================

// ValidationErrorResult 表示验证错误响应。
type ValidationErrorResult struct {
	StatusCode int
	Errors     validation.ValidationErrors
}

// ExecuteResult 实现 IActionResult 接口。
func (r ValidationErrorResult) ExecuteResult(c *gin.Context) {
	c.JSON(r.StatusCode, ApiResponse{
		Success: false,
		Error: &ApiError{
			Code:    errors.ValidationFailed,
			Message: "验证失败",
			Fields:  r.Errors,
		},
	})
}

// ValidationBadRequest 创建 400 Bad Request 结果，带有验证错误。
func ValidationBadRequest(errs validation.ValidationErrors) IActionResult {
	return ValidationErrorResult{StatusCode: 400, Errors: errs}
}

// ValidationBadRequestWithCode 创建 400 Bad Request 结果，带有验证错误和自定义错误码。
func ValidationBadRequestWithCode(code string, errs validation.ValidationErrors) IActionResult {
	return &customValidationErrorResult{
		StatusCode: 400,
		Code:       code,
		Errors:     errs,
	}
}

// customValidationErrorResult 用于自定义错误码。
type customValidationErrorResult struct {
	StatusCode int
	Code       string
	Errors     validation.ValidationErrors
}

// ExecuteResult 实现 IActionResult 接口。
func (r customValidationErrorResult) ExecuteResult(c *gin.Context) {
	c.JSON(r.StatusCode, ApiResponse{
		Success: false,
		Error: &ApiError{
			Code:    r.Code,
			Message: "验证失败",
			Fields:  r.Errors,
		},
	})
}

// ==================== 业务错误结果 ====================

// FrameworkErrorResult 表示框架错误响应。
type FrameworkErrorResult struct {
	StatusCode int
	Error      *errors.Error
}

// ExecuteResult 实现 IActionResult 接口。
func (r FrameworkErrorResult) ExecuteResult(c *gin.Context) {
	apiError := &ApiError{
		Code:    r.Error.Code(),
		Message: r.Error.Message(),
	}
	// 包含 Details 字段（如果有）
	details := r.Error.Details()
	if len(details) > 0 {
		apiError.Details = details
	}
	c.JSON(r.StatusCode, ApiResponse{
		Success: false,
		Error:   apiError,
	})
}

// FrameworkError 创建框架错误结果，自动映射 HTTP 状态码。
// 将常见错误模式映射到适当的 HTTP 状态码：
// - NOT_FOUND -> 404
// - ALREADY_EXISTS -> 409
// - PERMISSION_DENIED -> 403
// - INVALID_* -> 400
// - 默认 -> 400
func FrameworkError(err *errors.Error) IActionResult {
	statusCode := mapErrorToStatusCode(err.Code())
	return FrameworkErrorResult{StatusCode: statusCode, Error: err}
}

// FrameworkErrorWithStatus 创建框架错误结果，带有指定的 HTTP 状态码。
func FrameworkErrorWithStatus(statusCode int, err *errors.Error) IActionResult {
	return FrameworkErrorResult{StatusCode: statusCode, Error: err}
}

// mapErrorToStatusCode 将错误码映射到 HTTP 状态码。
func mapErrorToStatusCode(code string) int {
	switch {
	case containsPattern(code, "NOT_FOUND"):
		return 404
	case containsPattern(code, "ALREADY_EXISTS"):
		return 409
	case containsPattern(code, "PERMISSION_DENIED"):
		return 403
	case containsPattern(code, "UNAUTHORIZED"):
		return 401
	case containsPattern(code, "FORBIDDEN"):
		return 403
	case containsPattern(code, "INVALID"):
		return 400
	case containsPattern(code, "EXPIRED"):
		return 410
	case containsPattern(code, "LOCKED"):
		return 423
	case containsPattern(code, "LIMIT_EXCEEDED"):
		return 429
	default:
		return 400
	}
}

// containsPattern 检查错误码是否包含特定模式。
func containsPattern(code, pattern string) bool {
	return len(code) >= len(pattern) &&
		(code[len(code)-len(pattern):] == pattern ||
			containsSubstring(code, "."+pattern))
}

// containsSubstring 简单的子字符串检查。
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ==================== JSON 结果 ====================

// JsonResult 表示自定义 JSON 响应。
type JsonResult struct {
	StatusCode int
	Data       interface{}
}

// ExecuteResult 实现 IActionResult 接口。
func (r JsonResult) ExecuteResult(c *gin.Context) {
	c.JSON(r.StatusCode, r.Data)
}

// JSON 创建自定义 JSON 结果。
func JSON(statusCode int, data interface{}) IActionResult {
	return JsonResult{StatusCode: statusCode, Data: data}
}

// ==================== 内容结果 ====================

// ContentResult 表示纯文本响应。
type ContentResult struct {
	StatusCode  int
	Content     string
	ContentType string
}

// ExecuteResult 实现 IActionResult 接口。
func (r ContentResult) ExecuteResult(c *gin.Context) {
	if r.ContentType != "" {
		c.Data(r.StatusCode, r.ContentType, []byte(r.Content))
	} else {
		c.String(r.StatusCode, r.Content)
	}
}

// Content 创建纯文本结果。
func Content(statusCode int, content string) IActionResult {
	return ContentResult{StatusCode: statusCode, Content: content}
}

// ContentWithType 创建内容结果，带有自定义内容类型。
func ContentWithType(statusCode int, content, contentType string) IActionResult {
	return ContentResult{StatusCode: statusCode, Content: content, ContentType: contentType}
}

// ==================== 文件结果 ====================

// FileResult 表示文件下载响应。
type FileResult struct {
	FilePath string
	FileName string
}

// ExecuteResult 实现 IActionResult 接口。
func (r FileResult) ExecuteResult(c *gin.Context) {
	if r.FileName != "" {
		c.FileAttachment(r.FilePath, r.FileName)
	} else {
		c.File(r.FilePath)
	}
}

// File 创建文件结果。
func File(filePath string) IActionResult {
	return FileResult{FilePath: filePath}
}

// FileDownload 创建文件下载结果，带有自定义文件名。
func FileDownload(filePath, fileName string) IActionResult {
	return FileResult{FilePath: filePath, FileName: fileName}
}

// ==================== 状态结果 ====================

// StatusResult 表示仅包含状态码的响应。
type StatusResult struct {
	StatusCode int
}

// ExecuteResult 实现 IActionResult 接口。
func (r StatusResult) ExecuteResult(c *gin.Context) {
	c.Status(r.StatusCode)
}

// Status 创建仅包含状态码的结果。
func Status(statusCode int) IActionResult {
	return StatusResult{StatusCode: statusCode}
}

// ==================== 智能错误处理 ====================

// ErrorHandler 错误处理器函数类型。
// 接收错误和默认消息，返回 IActionResult。
// 如果返回 nil，表示该处理器不处理此错误，继续尝试其他处理器。
type ErrorHandler func(err error, defaultMessage ...string) IActionResult

var (
	errorHandlersMu     sync.Mutex
	errorHandlersAtomic atomic.Value // 存储 map[reflect.Type]ErrorHandler
)

func init() {
	// 初始化空 map
	errorHandlersAtomic.Store(make(map[reflect.Type]ErrorHandler))
}

// RegisterErrorHandler 注册错误类型处理器（泛型版本）。
// T: 错误类型，handler: 处理函数。
//
// 使用示例：
//
//	// 注册自定义错误处理器
//	web.RegisterErrorHandler[*MyCustomError](func(err *MyCustomError, msg ...string) web.IActionResult {
//	    return web.BadRequest(err.Details)
//	})
//
//	// 注册数据库错误处理器
//	web.RegisterErrorHandler[*sql.ErrNoRows](func(err *sql.ErrNoRows, msg ...string) web.IActionResult {
//	    return web.NotFound("记录不存在")
//	})
func RegisterErrorHandler[T error](handler func(T, ...string) IActionResult) {
	errorHandlersMu.Lock()
	defer errorHandlersMu.Unlock()

	// 获取错误类型
	var zero T
	errType := reflect.TypeOf(zero)

	// Copy-on-Write：复制现有 map 并添加新处理器
	oldMap := errorHandlersAtomic.Load().(map[reflect.Type]ErrorHandler)
	newMap := make(map[reflect.Type]ErrorHandler, len(oldMap)+1)
	for k, v := range oldMap {
		newMap[k] = v
	}

	// 包装为通用 ErrorHandler
	newMap[errType] = func(err error, msg ...string) IActionResult {
		typedErr, ok := err.(T)
		if !ok {
			return nil
		}
		return handler(typedErr, msg...)
	}

	// 原子替换
	errorHandlersAtomic.Store(newMap)
}

// ClearErrorHandlers 清除所有自定义错误处理器。
// 主要用于测试场景。
func ClearErrorHandlers() {
	errorHandlersMu.Lock()
	defer errorHandlersMu.Unlock()
	errorHandlersAtomic.Store(make(map[reflect.Type]ErrorHandler))
}

// FromError 智能处理各种类型的错误并返回对应的 ActionResult。
// 错误处理优先级：
// 1. 自定义错误处理器（如果已注册）
// 2. *errors.Error：自动映射 HTTP 状态码
// 3. validation.ValidationErrors：返回验证错误响应
// 4. 普通 error：返回内部错误，使用自定义消息
//
// 使用示例：
//
//	user, err := service.GetUser(id)
//	if err != nil {
//	    return web.FromError(err, "获取用户失败")
//	}
func FromError(err error, defaultMessage ...string) IActionResult {
	// nil 检查
	if err == nil {
		return nil
	}

	// 1. 尝试类型匹配处理器（O(1) map 查找）
	handlersMap := errorHandlersAtomic.Load().(map[reflect.Type]ErrorHandler)
	if handler, ok := handlersMap[reflect.TypeOf(err)]; ok {
		if result := handler(err, defaultMessage...); result != nil {
			return result
		}
	}

	// 2. 检查是否为框架 Error
	if fwErr, ok := err.(*errors.Error); ok {
		return FrameworkError(fwErr)
	}

	// 3. 检查是否为 ValidationErrors
	if valErrs, ok := err.(validation.ValidationErrors); ok {
		return ValidationBadRequest(valErrs)
	}

	// 4. 普通 error，返回 500
	msg := "服务器内部错误"
	if len(defaultMessage) > 0 && defaultMessage[0] != "" {
		msg = defaultMessage[0]
	}
	return InternalError(msg)
}

// FromErrorWithStatus 类似 FromError，但允许为普通 error 指定自定义 HTTP 状态码。
// 错误处理优先级：
// 1. 自定义错误处理器（如果已注册）
// 2. *errors.Error：忽略 statusCode，使用自动映射
// 3. validation.ValidationErrors：忽略 statusCode，固定返回 400
// 4. 普通 error：使用指定的 statusCode
//
// 使用示例：
//
//	err := db.Connect()
//	if err != nil {
//	    return web.FromErrorWithStatus(err, 503, "数据库服务暂时不可用")
//	}
func FromErrorWithStatus(err error, statusCode int, defaultMessage ...string) IActionResult {
	// nil 检查
	if err == nil {
		return nil
	}

	// 1. 尝试类型匹配处理器（O(1) map 查找，无循环）
	handlersMap := errorHandlersAtomic.Load().(map[reflect.Type]ErrorHandler)
	if handler, ok := handlersMap[reflect.TypeOf(err)]; ok {
		if result := handler(err, defaultMessage...); result != nil {
			return result
		}
	}

	// 2. 检查是否为框架 Error（忽略 statusCode）
	if fwErr, ok := err.(*errors.Error); ok {
		return FrameworkError(fwErr)
	}

	// 3. 检查是否为 ValidationErrors（忽略 statusCode）
	if valErrs, ok := err.(validation.ValidationErrors); ok {
		return ValidationBadRequest(valErrs)
	}

	// 4. 普通 error，使用指定的 statusCode
	msg := "服务器内部错误"
	if len(defaultMessage) > 0 && defaultMessage[0] != "" {
		msg = defaultMessage[0]
	}

	// 根据状态码确定错误码
	code := "ERROR"
	switch statusCode {
	case 400:
		code = "BAD_REQUEST"
	case 401:
		code = "UNAUTHORIZED"
	case 403:
		code = "FORBIDDEN"
	case 404:
		code = "NOT_FOUND"
	case 409:
		code = "CONFLICT"
	case 500:
		code = "INTERNAL_ERROR"
	case 503:
		code = "SERVICE_UNAVAILABLE"
	}

	return Error(statusCode, code, msg)
}

// ==================== 图片结果 ====================

// Base64ImageResult 表示以 base64 格式在 JSON 中返回的图片响应。
type Base64ImageResult struct {
	StatusCode  int
	ImageData   string // Base64 编码的图片
	ContentType string // 原始图片内容类型（如 image/png）
}

// ExecuteResult 实现 IActionResult 接口。
func (r Base64ImageResult) ExecuteResult(c *gin.Context) {
	c.JSON(r.StatusCode, ApiResponse{
		Success: true,
		Data: M{
			"image":       r.ImageData,
			"contentType": r.ContentType,
		},
	})
}

// Base64Image 创建 base64 图片结果。
// 图片数据将被编码为 base64 并以 JSON 格式返回。
func Base64Image(imageData []byte, contentType string) IActionResult {
	encoded := base64.StdEncoding.EncodeToString(imageData)
	return Base64ImageResult{
		StatusCode:  200,
		ImageData:   encoded,
		ContentType: contentType,
	}
}

// BinaryImageResult 表示二进制图片响应。
type BinaryImageResult struct {
	StatusCode  int
	ImageData   []byte
	ContentType string // 如 image/png、image/jpeg
}

// ExecuteResult 实现 IActionResult 接口。
func (r BinaryImageResult) ExecuteResult(c *gin.Context) {
	c.Data(r.StatusCode, r.ContentType, r.ImageData)
}

// BinaryImage 创建二进制图片结果。
// 图片数据将以指定的内容类型作为原始二进制返回。
func BinaryImage(imageData []byte, contentType string) IActionResult {
	return BinaryImageResult{
		StatusCode:  200,
		ImageData:   imageData,
		ContentType: contentType,
	}
}

// PNG 创建 PNG 图片结果。
// BinaryImage 的便捷方法，内容类型为 image/png。
func PNG(imageData []byte) IActionResult {
	return BinaryImage(imageData, "image/png")
}

// JPEG 创建 JPEG 图片结果。
// BinaryImage 的便捷方法，内容类型为 image/jpeg。
func JPEG(imageData []byte) IActionResult {
	return BinaryImage(imageData, "image/jpeg")
}

// WebP 创建 WebP 图片结果。
// BinaryImage 的便捷方法，内容类型为 image/webp。
func WebP(imageData []byte) IActionResult {
	return BinaryImage(imageData, "image/webp")
}
