package web

import (
	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/di"
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

// MakeToGinHandler creates a handler converter that injects services into HttpContext.
// This factory function captures the services and returns a converter function.
func MakeToGinHandler(services di.IServiceProvider) func(Handler) gin.HandlerFunc {
	return func(handler Handler) gin.HandlerFunc {
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
				ctx := &HttpContext{
					gin:      c,
					Services: services,
				}
				h(ctx)
			}

		case func(*HttpContext):
			// HttpContext handler without return value (type alias)
			return func(c *gin.Context) {
				ctx := &HttpContext{
					gin:      c,
					Services: services,
				}
				h(ctx)
			}

		case ActionHandlerFunc:
			// ActionResult handler
			return func(c *gin.Context) {
				ctx := &HttpContext{
					gin:      c,
					Services: services,
				}
				result := h(ctx)
				if result != nil {
					result.ExecuteResult(c)
				}
			}

		case func(*HttpContext) IActionResult:
			// ActionResult handler (type alias)
			return func(c *gin.Context) {
				ctx := &HttpContext{
					gin:      c,
					Services: services,
				}
				result := h(ctx)
				if result != nil {
					result.ExecuteResult(c)
				}
			}

		default:
			// Fallback: panic with clear error message
			panic("unsupported handler type: must be gin.HandlerFunc, func(*HttpContext), or func(*HttpContext) IActionResult")
		}
	}
}

// ToGinHandler converts various handler types to gin.HandlerFunc.
// Deprecated: Use MakeToGinHandler factory function instead to inject services.
// This function creates HttpContext without services, for backward compatibility only.
func ToGinHandler(handler Handler) gin.HandlerFunc {
	return MakeToGinHandler(nil)(handler)
}

// MakeToGinHandlers creates a function that converts multiple handlers with services injection.
func MakeToGinHandlers(services di.IServiceProvider) func(...Handler) []gin.HandlerFunc {
	converter := MakeToGinHandler(services)
	return func(handlers ...Handler) []gin.HandlerFunc {
		result := make([]gin.HandlerFunc, len(handlers))
		for i, h := range handlers {
			result[i] = converter(h)
		}
		return result
	}
}

// ToGinHandlers converts multiple handlers to gin.HandlerFunc slice.
// Deprecated: Use MakeToGinHandlers factory function instead to inject services.
func ToGinHandlers(handlers ...Handler) []gin.HandlerFunc {
	return MakeToGinHandlers(nil)(handlers...)
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

