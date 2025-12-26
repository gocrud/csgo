package web

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/hosting"
	"github.com/gocrud/csgo/web/router"
)

// WebApplication 表示已配置的 Web 应用程序。
type WebApplication struct {
	host        hosting.IHost
	engine      *gin.Engine
	Services    di.IServiceProvider // ✅ 直接暴露，强类型
	Environment hosting.IHostEnvironment
	routes      []*router.RouteBuilder
	groups      []*router.RouteGroupBuilder
	runtimeUrls *[]string // 指向运行时 URL 的指针（与 HttpServer 共享）

	// 带有服务注入的处理器转换器
	toHandler  func(Handler) gin.HandlerFunc
	toHandlers func(...Handler) []gin.HandlerFunc
}

// Run 运行 Web 应用程序并阻塞直到关闭。
// 如果提供了 urls 参数，它们将覆盖配置的监听地址。
// 对应 .NET 的 app.Run(url)。
func (app *WebApplication) Run(urls ...string) error {
	if len(urls) > 0 && app.runtimeUrls != nil {
		*app.runtimeUrls = urls
	}
	return app.host.Run()
}

// RunWithContext 使用自定义上下文运行 Web 应用程序并阻塞直到关闭。
func (app *WebApplication) RunWithContext(ctx context.Context) error {
	return app.host.RunWithContext(ctx)
}

// Start 启动 Web 应用程序。
func (app *WebApplication) Start(ctx context.Context) error {
	return app.host.Start(ctx)
}

// Stop 停止 Web 应用程序。
func (app *WebApplication) Stop(ctx context.Context) error {
	return app.host.Stop(ctx)
}

// Use 向管道添加中间件。
func (app *WebApplication) Use(middleware ...gin.HandlerFunc) {
	app.engine.Use(middleware...)
}

// GET 注册 GET 端点。
// 处理器必须是 ActionHandlerFunc: func(*HttpContext) IActionResult
func (app *WebApplication) GET(pattern string, handlers ...Handler) router.IEndpointConventionBuilder {
	return app.mapRoute("GET", pattern, handlers...)
}

// POST 注册 POST 端点。
// 处理器必须是 ActionHandlerFunc: func(*HttpContext) IActionResult
func (app *WebApplication) POST(pattern string, handlers ...Handler) router.IEndpointConventionBuilder {
	return app.mapRoute("POST", pattern, handlers...)
}

// PUT 注册 PUT 端点。
// 处理器必须是 ActionHandlerFunc: func(*HttpContext) IActionResult
func (app *WebApplication) PUT(pattern string, handlers ...Handler) router.IEndpointConventionBuilder {
	return app.mapRoute("PUT", pattern, handlers...)
}

// DELETE 注册 DELETE 端点。
// 处理器必须是 ActionHandlerFunc: func(*HttpContext) IActionResult
func (app *WebApplication) DELETE(pattern string, handlers ...Handler) router.IEndpointConventionBuilder {
	return app.mapRoute("DELETE", pattern, handlers...)
}

// PATCH 注册 PATCH 端点。
// 处理器必须是 ActionHandlerFunc: func(*HttpContext) IActionResult
func (app *WebApplication) PATCH(pattern string, handlers ...Handler) router.IEndpointConventionBuilder {
	return app.mapRoute("PATCH", pattern, handlers...)
}

// Group 创建路由组。
// 支持 ActionHandlerFunc: func(*HttpContext) IActionResult
func (app *WebApplication) Group(prefix string, handlers ...Handler) *router.RouteGroupBuilder {
	// 使用服务感知转换器转换组中间件的处理器
	ginHandlers := app.toHandlers(handlers...)

	ginGroup := app.engine.Group(prefix, ginHandlers...)
	group := router.NewRouteGroupBuilder(ginGroup, prefix)

	// 创建适配器，将 router.Handler (any) 转换为 web.Handler (ActionHandlerFunc)
	routerHandlerConverter := func(h router.Handler) gin.HandlerFunc {
		// 将 router.Handler 转换为 web.Handler
		webHandler, ok := h.(Handler)
		if !ok {
			panic("handler must be ActionHandlerFunc: func(*HttpContext) IActionResult")
		}
		return app.toHandler(webHandler)
	}
	group.SetHandlerConverter(routerHandlerConverter)

	app.groups = append(app.groups, group)
	return group
}

// mapRoute 是注册路由的内部方法。
// 所有处理器必须是 ActionHandlerFunc: func(*HttpContext) IActionResult
func (app *WebApplication) mapRoute(method, pattern string, handlers ...Handler) router.IEndpointConventionBuilder {
	// 使用服务感知转换器将处理器转换为 gin.HandlerFunc
	ginHandlers := app.toHandlers(handlers...)

	// 在 Gin 中注册
	app.engine.Handle(method, pattern, ginHandlers...)

	// 创建路由构建器
	rb := router.NewRouteBuilder(method, pattern)
	app.routes = append(app.routes, rb)

	return rb
}

// GetRoutes 返回所有已注册的路由。
func (app *WebApplication) GetRoutes() []*router.RouteBuilder {
	allRoutes := make([]*router.RouteBuilder, 0)

	// 添加顶层路由
	allRoutes = append(allRoutes, app.routes...)

	// 添加来自组的路由
	for _, group := range app.groups {
		allRoutes = append(allRoutes, group.GetRoutes()...)
	}

	return allRoutes
}
