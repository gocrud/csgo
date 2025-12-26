package router

import (
	"path"

	"github.com/gin-gonic/gin"
)

// Handler 表示统一的处理器类型。
// 使用 'any' 允许任何处理器类型，通过 handlerConvertFn 转换。
type Handler = any

// IEndpointRouteBuilder 定义应用程序中路由构建器的契约。
type IEndpointRouteBuilder interface {
	// GET 注册 GET 端点。
	GET(pattern string, handlers ...Handler) IEndpointConventionBuilder

	// POST 注册 POST 端点。
	POST(pattern string, handlers ...Handler) IEndpointConventionBuilder

	// PUT 注册 PUT 端点。
	PUT(pattern string, handlers ...Handler) IEndpointConventionBuilder

	// DELETE 注册 DELETE 端点。
	DELETE(pattern string, handlers ...Handler) IEndpointConventionBuilder

	// PATCH 注册 PATCH 端点。
	PATCH(pattern string, handlers ...Handler) IEndpointConventionBuilder

	// Group 创建路由组。
	Group(prefix string, handlers ...gin.HandlerFunc) *RouteGroupBuilder
}

// RouteGroupBuilder 表示具有公共前缀的端点组。
type RouteGroupBuilder struct {
	ginGroup         *gin.RouterGroup
	prefix           string
	metadata         []interface{}
	routes           []*RouteBuilder
	childGroups      []*RouteGroupBuilder // 跟踪子组
	handlerConvertFn func(Handler) gin.HandlerFunc
}

// NewRouteGroupBuilder 创建新的 RouteGroupBuilder。
func NewRouteGroupBuilder(ginGroup *gin.RouterGroup, prefix string) *RouteGroupBuilder {
	return &RouteGroupBuilder{
		ginGroup:    ginGroup,
		prefix:      prefix,
		metadata:    make([]interface{}, 0),
		routes:      make([]*RouteBuilder, 0),
		childGroups: make([]*RouteGroupBuilder, 0),
	}
}

// SetHandlerConverter 设置处理器转换函数。
// 用于将自定义处理器类型转换为 gin.HandlerFunc。
func (g *RouteGroupBuilder) SetHandlerConverter(fn func(Handler) gin.HandlerFunc) {
	g.handlerConvertFn = fn
}

// convertHandler 将 Handler 转换为 gin.HandlerFunc。
func (g *RouteGroupBuilder) convertHandler(h Handler) gin.HandlerFunc {
	// 如果设置了转换器，使用它
	if g.handlerConvertFn != nil {
		return g.handlerConvertFn(h)
	}

	// 默认：假设它已经是 gin.HandlerFunc
	if ginHandler, ok := h.(gin.HandlerFunc); ok {
		return ginHandler
	}
	if ginHandler, ok := h.(func(*gin.Context)); ok {
		return ginHandler
	}

	panic("不支持的处理器类型：请设置处理器转换器或使用 gin.HandlerFunc")
}

// convertHandlers 将多个 Handler 转换为 gin.HandlerFunc 切片。
func (g *RouteGroupBuilder) convertHandlers(handlers ...Handler) []gin.HandlerFunc {
	result := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		result[i] = g.convertHandler(h)
	}
	return result
}

// GET 注册 GET 端点。
// 设置处理器转换器时支持多种处理器类型。
func (g *RouteGroupBuilder) GET(pattern string, handlers ...Handler) IEndpointConventionBuilder {
	return g.mapRoute("GET", pattern, handlers...)
}

// POST 注册 POST 端点。
func (g *RouteGroupBuilder) POST(pattern string, handlers ...Handler) IEndpointConventionBuilder {
	return g.mapRoute("POST", pattern, handlers...)
}

// PUT 注册 PUT 端点。
func (g *RouteGroupBuilder) PUT(pattern string, handlers ...Handler) IEndpointConventionBuilder {
	return g.mapRoute("PUT", pattern, handlers...)
}

// DELETE 注册 DELETE 端点。
func (g *RouteGroupBuilder) DELETE(pattern string, handlers ...Handler) IEndpointConventionBuilder {
	return g.mapRoute("DELETE", pattern, handlers...)
}

// PATCH 注册 PATCH 端点。
func (g *RouteGroupBuilder) PATCH(pattern string, handlers ...Handler) IEndpointConventionBuilder {
	return g.mapRoute("PATCH", pattern, handlers...)
}

// Group 创建嵌套路由组。
func (g *RouteGroupBuilder) Group(prefix string, handlers ...gin.HandlerFunc) *RouteGroupBuilder {
	newGinGroup := g.ginGroup.Group(prefix, handlers...)
	newPrefix := path.Join(g.prefix, prefix)

	newGroup := NewRouteGroupBuilder(newGinGroup, newPrefix)

	// 继承父级元数据
	newGroup.metadata = append([]interface{}{}, g.metadata...)

	// 继承处理器转换器
	newGroup.handlerConvertFn = g.handlerConvertFn

	// 跟踪子组
	g.childGroups = append(g.childGroups, newGroup)

	return newGroup
}

// mapRoute 是注册路由的内部方法。
func (g *RouteGroupBuilder) mapRoute(method, pattern string, handlers ...Handler) IEndpointConventionBuilder {
	// 转换处理器
	ginHandlers := g.convertHandlers(handlers...)

	// 在 Gin 中注册
	g.ginGroup.Handle(method, pattern, ginHandlers...)

	// 计算完整路径
	fullPath := path.Join(g.prefix, pattern)

	// 创建路由构建器
	rb := NewRouteBuilder(method, fullPath)

	// 继承组元数据
	rb.metadata = append([]interface{}{}, g.metadata...)

	// 从组继承 OpenAPI 设置
	// 如果组启用了 OpenAPI，所有子路由将自动继承它
	for _, meta := range g.metadata {
		if openApiMeta, ok := meta.(*OpenApiMetadata); ok && openApiMeta.Enabled {
			rb.openApiEnabled = true
		}
		if groupConfig, ok := meta.(*GroupOpenApiConfig); ok && groupConfig.Configure != nil {
			// 将组的 OpenAPI 配置应用到此路由
			rb.openApiEnabled = true
			builder := &OpenApiBuilder{builder: rb}
			groupConfig.Configure(builder)
		}
	}

	// 存储路由
	g.routes = append(g.routes, rb)

	return rb
}

// WithOpenApi 为此组启用 OpenAPI 文档。
// 配置将应用于此组中创建的所有路由。
func (g *RouteGroupBuilder) WithOpenApi(configure func(*OpenApiBuilder)) *RouteGroupBuilder {
	g.metadata = append(g.metadata, &OpenApiMetadata{Enabled: true})

	// 将配置存储在元数据中以应用于子路由
	if configure != nil {
		g.metadata = append(g.metadata, &GroupOpenApiConfig{Configure: configure})
	}

	return g
}

// GetRoutes 递归返回此组及其子组中的所有路由。
func (g *RouteGroupBuilder) GetRoutes() []*RouteBuilder {
	allRoutes := make([]*RouteBuilder, 0)

	// 添加此组中的路由
	allRoutes = append(allRoutes, g.routes...)

	// 递归添加子组中的路由
	for _, childGroup := range g.childGroups {
		allRoutes = append(allRoutes, childGroup.GetRoutes()...)
	}

	return allRoutes
}

// OpenApiMetadata 表示 OpenAPI 元数据。
type OpenApiMetadata struct {
	Enabled bool
}

// GroupOpenApiConfig 保存要应用于组路由的 OpenAPI 配置。
type GroupOpenApiConfig struct {
	Configure func(*OpenApiBuilder)
}

// AuthorizationMetadata 表示授权元数据。
type AuthorizationMetadata struct {
	Policies []string
}

// TagsMetadata 表示 OpenAPI 标签元数据。
type TagsMetadata struct {
	Tags []string
}
