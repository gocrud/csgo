package router

// EndpointOption 是配置端点的函数。
type EndpointOption func(IEndpointConventionBuilder) IEndpointConventionBuilder

// IEndpointConventionBuilder 用于配置端点。
// 所有配置应通过 WithOpenApi() 和端点选项完成。
type IEndpointConventionBuilder interface {
	// WithOpenApi 为此端点启用 OpenAPI 文档并应用选项。
	// 对应 .NET 的 endpoint.WithOpenApi()。
	WithOpenApi(configure func(*OpenApiBuilder)) IEndpointConventionBuilder
}

// RouteBuilder 实现 IEndpointConventionBuilder 接口。
type RouteBuilder struct {
	method                  string
	path                    string
	name                    string
	displayName             string
	summary                 string
	description             string
	tags                    []string
	metadata                []interface{}
	authPolicies            []string
	allowAnonymous          bool
	openApiEnabled          bool
	apiSecurityRequirements []map[string][]string // OpenAPI 安全要求
}

// NewRouteBuilder 创建新的 RouteBuilder。
func NewRouteBuilder(method, path string) *RouteBuilder {
	return &RouteBuilder{
		method:                  method,
		path:                    path,
		tags:                    make([]string, 0),
		metadata:                make([]interface{}, 0),
		authPolicies:            make([]string, 0),
		apiSecurityRequirements: make([]map[string][]string, 0),
	}
}

// WithOpenApi 为此端点启用 OpenAPI 文档并应用选项。
// 对应 .NET 的 endpoint.WithOpenApi()。
func (b *RouteBuilder) WithOpenApi(configure func(*OpenApiBuilder)) IEndpointConventionBuilder {
	b.openApiEnabled = true

	// Create builder and apply configuration
	builder := &OpenApiBuilder{builder: b}
	configure(builder)

	return b
}

// Setter methods for OpenApiBuilder to configure endpoints

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

// AddParameterMetadata adds parameter metadata.
func (b *RouteBuilder) AddParameterMetadata(metadata ParameterMetadata) {
	b.metadata = append(b.metadata, &metadata)
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

// SetApiSecurityRequirements sets the API security requirements for OpenAPI documentation.
func (b *RouteBuilder) SetApiSecurityRequirements(requirements []map[string][]string) {
	b.apiSecurityRequirements = requirements
}

// GetApiSecurityRequirements returns the API security requirements.
func (b *RouteBuilder) GetApiSecurityRequirements() []map[string][]string {
	return b.apiSecurityRequirements
}

// AddApiSecurityRequirement adds a single security requirement.
func (b *RouteBuilder) AddApiSecurityRequirement(name string, scopes []string) {
	if scopes == nil {
		scopes = []string{}
	}
	b.apiSecurityRequirements = append(b.apiSecurityRequirements, map[string][]string{
		name: scopes,
	})
}
