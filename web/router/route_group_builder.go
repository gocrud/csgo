package router

import (
	"path"

	"github.com/gin-gonic/gin"
)

// Handler represents a unified handler type.
// Uses 'any' to allow any handler type, converted via handlerConvertFn.
type Handler = any

// IEndpointRouteBuilder defines a contract for a route builder in an application.
type IEndpointRouteBuilder interface {
	// MapGet registers a GET endpoint.
	MapGet(pattern string, handlers ...Handler) IEndpointConventionBuilder

	// MapPost registers a POST endpoint.
	MapPost(pattern string, handlers ...Handler) IEndpointConventionBuilder

	// MapPut registers a PUT endpoint.
	MapPut(pattern string, handlers ...Handler) IEndpointConventionBuilder

	// MapDelete registers a DELETE endpoint.
	MapDelete(pattern string, handlers ...Handler) IEndpointConventionBuilder

	// MapPatch registers a PATCH endpoint.
	MapPatch(pattern string, handlers ...Handler) IEndpointConventionBuilder

	// MapGroup creates a route group.
	MapGroup(prefix string, handlers ...gin.HandlerFunc) *RouteGroupBuilder
}

// RouteGroupBuilder represents a group of endpoints with a common prefix.
type RouteGroupBuilder struct {
	ginGroup         *gin.RouterGroup
	prefix           string
	metadata         []interface{}
	routes           []*RouteBuilder
	handlerConvertFn func(Handler) gin.HandlerFunc
}

// NewRouteGroupBuilder creates a new RouteGroupBuilder.
func NewRouteGroupBuilder(ginGroup *gin.RouterGroup, prefix string) *RouteGroupBuilder {
	return &RouteGroupBuilder{
		ginGroup: ginGroup,
		prefix:   prefix,
		metadata: make([]interface{}, 0),
		routes:   make([]*RouteBuilder, 0),
	}
}

// SetHandlerConverter sets the handler converter function.
// This is used to convert custom handler types to gin.HandlerFunc.
func (g *RouteGroupBuilder) SetHandlerConverter(fn func(Handler) gin.HandlerFunc) {
	g.handlerConvertFn = fn
}

// convertHandler converts a Handler to gin.HandlerFunc.
func (g *RouteGroupBuilder) convertHandler(h Handler) gin.HandlerFunc {
	// If a converter is set, use it
	if g.handlerConvertFn != nil {
		return g.handlerConvertFn(h)
	}

	// Default: assume it's already a gin.HandlerFunc
	if ginHandler, ok := h.(gin.HandlerFunc); ok {
		return ginHandler
	}
	if ginHandler, ok := h.(func(*gin.Context)); ok {
		return ginHandler
	}

	panic("unsupported handler type: set handler converter or use gin.HandlerFunc")
}

// convertHandlers converts multiple Handlers to gin.HandlerFunc slice.
func (g *RouteGroupBuilder) convertHandlers(handlers ...Handler) []gin.HandlerFunc {
	result := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		result[i] = g.convertHandler(h)
	}
	return result
}

// MapGet registers a GET endpoint.
// Supports multiple handler types when handler converter is set.
func (g *RouteGroupBuilder) MapGet(pattern string, handlers ...Handler) IEndpointConventionBuilder {
	return g.mapRoute("GET", pattern, handlers...)
}

// MapPost registers a POST endpoint.
func (g *RouteGroupBuilder) MapPost(pattern string, handlers ...Handler) IEndpointConventionBuilder {
	return g.mapRoute("POST", pattern, handlers...)
}

// MapPut registers a PUT endpoint.
func (g *RouteGroupBuilder) MapPut(pattern string, handlers ...Handler) IEndpointConventionBuilder {
	return g.mapRoute("PUT", pattern, handlers...)
}

// MapDelete registers a DELETE endpoint.
func (g *RouteGroupBuilder) MapDelete(pattern string, handlers ...Handler) IEndpointConventionBuilder {
	return g.mapRoute("DELETE", pattern, handlers...)
}

// MapPatch registers a PATCH endpoint.
func (g *RouteGroupBuilder) MapPatch(pattern string, handlers ...Handler) IEndpointConventionBuilder {
	return g.mapRoute("PATCH", pattern, handlers...)
}

// MapGroup creates a nested route group.
func (g *RouteGroupBuilder) MapGroup(prefix string, handlers ...gin.HandlerFunc) *RouteGroupBuilder {
	newGinGroup := g.ginGroup.Group(prefix, handlers...)
	newPrefix := path.Join(g.prefix, prefix)

	newGroup := NewRouteGroupBuilder(newGinGroup, newPrefix)

	// Inherit parent metadata
	newGroup.metadata = append([]interface{}{}, g.metadata...)

	// Inherit handler converter
	newGroup.handlerConvertFn = g.handlerConvertFn

	return newGroup
}

// mapRoute is the internal method to register a route.
func (g *RouteGroupBuilder) mapRoute(method, pattern string, handlers ...Handler) IEndpointConventionBuilder {
	// Convert handlers
	ginHandlers := g.convertHandlers(handlers...)

	// Register with Gin
	g.ginGroup.Handle(method, pattern, ginHandlers...)

	// Calculate full path
	fullPath := path.Join(g.prefix, pattern)

	// Create route builder
	rb := NewRouteBuilder(method, fullPath)

	// Inherit group metadata
	rb.metadata = append([]interface{}{}, g.metadata...)

	// Inherit OpenAPI setting from group
	// If the group has enabled OpenAPI, all child routes automatically inherit it
	for _, meta := range g.metadata {
		if openApiMeta, ok := meta.(*OpenApiMetadata); ok && openApiMeta.Enabled {
			rb.openApiEnabled = true
		}
		if groupConfig, ok := meta.(*GroupOpenApiConfig); ok && groupConfig.Configure != nil {
			// Apply group's OpenAPI configuration to this route
			rb.openApiEnabled = true
			builder := &OpenApiBuilder{builder: rb}
			groupConfig.Configure(builder)
		}
	}

	// Store route
	g.routes = append(g.routes, rb)

	return rb
}

// WithOpenApi enables OpenAPI documentation for this group.
// Configuration will be applied to all routes created in this group.
func (g *RouteGroupBuilder) WithOpenApi(configure func(*OpenApiBuilder)) *RouteGroupBuilder {
	g.metadata = append(g.metadata, &OpenApiMetadata{Enabled: true})

	// Store configuration in metadata to be applied to child routes
	if configure != nil {
		g.metadata = append(g.metadata, &GroupOpenApiConfig{Configure: configure})
	}

	return g
}

// GetRoutes returns all routes in this group.
func (g *RouteGroupBuilder) GetRoutes() []*RouteBuilder {
	return g.routes
}

// OpenApiMetadata represents OpenAPI metadata.
type OpenApiMetadata struct {
	Enabled bool
}

// GroupOpenApiConfig holds OpenAPI configuration to be applied to group routes.
type GroupOpenApiConfig struct {
	Configure func(*OpenApiBuilder)
}

// AuthorizationMetadata represents authorization metadata.
type AuthorizationMetadata struct {
	Policies []string
}

// TagsMetadata represents OpenAPI tags metadata.
type TagsMetadata struct {
	Tags []string
}
