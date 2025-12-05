package routing

import "reflect"

// EndpointOption is a function that configures an endpoint.
type EndpointOption func(IEndpointConventionBuilder) IEndpointConventionBuilder

// IEndpointConventionBuilder is used to configure endpoints.
// All configuration should be done through WithOpenApi() with endpoint options.
type IEndpointConventionBuilder interface {
	// WithOpenApi enables OpenAPI documentation for this endpoint and applies options.
	// Corresponds to .NET endpoint.WithOpenApi().
	WithOpenApi(options ...EndpointOption) IEndpointConventionBuilder
}

// RouteBuilder implements IEndpointConventionBuilder.
type RouteBuilder struct {
	method         string
	path           string
	name           string
	displayName    string
	summary        string
	description    string
	tags           []string
	metadata       []interface{}
	authPolicies   []string
	allowAnonymous bool
	openApiEnabled bool
}

// NewRouteBuilder creates a new RouteBuilder.
func NewRouteBuilder(method, path string) *RouteBuilder {
	return &RouteBuilder{
		method:       method,
		path:         path,
		tags:         make([]string, 0),
		metadata:     make([]interface{}, 0),
		authPolicies: make([]string, 0),
	}
}

// WithOpenApi enables OpenAPI documentation for this endpoint and applies options.
// Corresponds to .NET endpoint.WithOpenApi().
func (b *RouteBuilder) WithOpenApi(options ...EndpointOption) IEndpointConventionBuilder {
	b.openApiEnabled = true

	// Apply all options
	for _, option := range options {
		option(b)
	}

	return b
}

// Setter methods for openapi package to configure endpoints

// SetName sets the endpoint name.
func (b *RouteBuilder) SetName(name string) {
	b.name = name
}

// SetSummary sets the OpenAPI summary.
func (b *RouteBuilder) SetSummary(summary string) {
	b.summary = summary
}

// SetDescription sets the OpenAPI description.
func (b *RouteBuilder) SetDescription(description string) {
	b.description = description
}

// SetTags sets the OpenAPI tags.
func (b *RouteBuilder) SetTags(tags []string) {
	b.tags = tags
}

// AddTags adds OpenAPI tags.
func (b *RouteBuilder) AddTags(tags ...string) {
	b.tags = append(b.tags, tags...)
}

// AddResponseMetadata adds response metadata.
func (b *RouteBuilder) AddResponseMetadata(metadata ResponseMetadata) {
	b.metadata = append(b.metadata, &metadata)
}

// AddRequestMetadata adds request metadata.
func (b *RouteBuilder) AddRequestMetadata(metadata RequestMetadata) {
	b.metadata = append(b.metadata, &metadata)
}

// SetAuthorizationPolicies sets authorization policies.
func (b *RouteBuilder) SetAuthorizationPolicies(policies []string) {
	b.authPolicies = policies
	b.allowAnonymous = false
}

// SetAllowAnonymous sets whether anonymous access is allowed.
func (b *RouteBuilder) SetAllowAnonymous(allow bool) {
	b.allowAnonymous = allow
	if allow {
		b.authPolicies = nil
	}
}

// ResponseMetadata represents response metadata.
type ResponseMetadata struct {
	StatusCode int
	Type       reflect.Type
	IsProblem  bool
}

// RequestMetadata represents request metadata.
type RequestMetadata struct {
	ContentType string
	Type        reflect.Type
}

// GetMethod returns the HTTP method.
func (b *RouteBuilder) GetMethod() string {
	return b.method
}

// GetPath returns the route path.
func (b *RouteBuilder) GetPath() string {
	return b.path
}

// GetName returns the endpoint name.
func (b *RouteBuilder) GetName() string {
	return b.name
}

// GetSummary returns the OpenAPI summary.
func (b *RouteBuilder) GetSummary() string {
	return b.summary
}

// GetDescription returns the OpenAPI description.
func (b *RouteBuilder) GetDescription() string {
	return b.description
}

// GetTags returns the OpenAPI tags.
func (b *RouteBuilder) GetTags() []string {
	return b.tags
}

// GetMetadata returns all metadata.
func (b *RouteBuilder) GetMetadata() []interface{} {
	return b.metadata
}

// IsOpenApiEnabled returns whether OpenAPI documentation is enabled for this endpoint.
func (b *RouteBuilder) IsOpenApiEnabled() bool {
	return b.openApiEnabled
}
