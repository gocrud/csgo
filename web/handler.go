package web

import (
	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/di"
)

// Handler 表示统一的处理器类型，可以是 ActionHandlerFunc 或 gin.HandlerFunc。
// 支持两种类型：
//   - ActionHandlerFunc: func(*HttpContext) IActionResult
//   - gin.HandlerFunc: func(*gin.Context)
type Handler interface{}

// ActionHandlerFunc 是返回 IActionResult 的处理器函数。
type ActionHandlerFunc func(*HttpContext) IActionResult

// MakeToGinHandler 创建一个处理器转换器，将服务注入到 HttpContext 中。
// 此工厂函数捕获服务并返回转换器函数。
// 支持 ActionHandlerFunc 和 gin.HandlerFunc 两种类型：
//   - ActionHandlerFunc: func(*HttpContext) IActionResult
//   - gin.HandlerFunc: func(*gin.Context)
func MakeToGinHandler(services di.IServiceProvider) func(Handler) gin.HandlerFunc {
	return func(handler Handler) gin.HandlerFunc {
		// 检查是否是 gin.HandlerFunc
		if ginHandler, ok := handler.(func(*gin.Context)); ok {
			return ginHandler
		}

		// 检查是否是 ActionHandlerFunc
		if actionHandler, ok := handler.(func(*HttpContext) IActionResult); ok {
			return func(c *gin.Context) {
				ctx := &HttpContext{
					gin:      c,
					Services: services,
				}
				result := actionHandler(ctx)
				if result != nil {
					result.ExecuteResult(c)
				}
			}
		}

		// 如果都不是，panic
		panic("handler must be web.ActionHandlerFunc or gin.HandlerFunc")
	}
}

// MakeToGinHandlers 创建一个函数，用于转换多个处理器并注入服务。
// 支持混合使用 ActionHandlerFunc 和 gin.HandlerFunc。
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
