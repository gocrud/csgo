package web

import (
	"github.com/gin-gonic/gin"
)

// Handler represents a unified handler type that can be:
// - gin.HandlerFunc
// - func(*HttpContext)
// - func(*HttpContext) IActionResult
type Handler = any

// HandlerFunc is a handler function that uses HttpContext.
type HandlerFunc func(*HttpContext)

// ActionHandlerFunc is a handler function that returns IActionResult.
type ActionHandlerFunc func(*HttpContext) IActionResult

// ToGinHandler converts various handler types to gin.HandlerFunc.
// Supported types:
// - gin.HandlerFunc: used as-is
// - func(*HttpContext): wrapped to gin.HandlerFunc
// - func(*HttpContext) IActionResult: wrapped and executes result
func ToGinHandler(handler Handler) gin.HandlerFunc {
	switch h := handler.(type) {
	case gin.HandlerFunc:
		// Already a gin.HandlerFunc, use as-is
		return h

	case func(*gin.Context):
		// Raw gin handler function
		return h

	case HandlerFunc:
		// HttpContext handler without return value
		return func(c *gin.Context) {
			h(NewHttpContext(c))
		}

	case func(*HttpContext):
		// HttpContext handler without return value (type alias)
		return func(c *gin.Context) {
			h(NewHttpContext(c))
		}

	case ActionHandlerFunc:
		// ActionResult handler
		return func(c *gin.Context) {
			result := h(NewHttpContext(c))
			if result != nil {
				result.ExecuteResult(c)
			}
		}

	case func(*HttpContext) IActionResult:
		// ActionResult handler (type alias)
		return func(c *gin.Context) {
			result := h(NewHttpContext(c))
			if result != nil {
				result.ExecuteResult(c)
			}
		}

	default:
		// Fallback: panic with clear error message
		panic("unsupported handler type: must be gin.HandlerFunc, func(*HttpContext), or func(*HttpContext) IActionResult")
	}
}

// ToGinHandlers converts multiple handlers to gin.HandlerFunc slice.
func ToGinHandlers(handlers ...Handler) []gin.HandlerFunc {
	result := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		result[i] = ToGinHandler(h)
	}
	return result
}

// WrapHttpContext wraps a HttpContext handler to gin.HandlerFunc.
func WrapHttpContext(handler func(*HttpContext)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(NewHttpContext(c))
	}
}

// WrapActionResult wraps an ActionResult handler to gin.HandlerFunc.
func WrapActionResult(handler func(*HttpContext) IActionResult) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := handler(NewHttpContext(c))
		if result != nil {
			result.ExecuteResult(c)
		}
	}
}

