package routing

// IEndpointConventionBuilder is used to configure endpoints.
type IEndpointConventionBuilder interface {
	// WithName sets the endpoint name.
	WithName(name string) IEndpointConventionBuilder

	// WithDisplayName sets the display name for the endpoint.
	WithDisplayName(displayName string) IEndpointConventionBuilder

	// WithMetadata adds metadata to the endpoint.
	WithMetadata(metadata ...interface{}) IEndpointConventionBuilder

	// WithSummary sets the OpenAPI summary.
	WithSummary(summary string) IEndpointConventionBuilder

	// WithDescription sets the OpenAPI description.
	WithDescription(description string) IEndpointConventionBuilder

	// WithTags adds OpenAPI tags.
	WithTags(tags ...string) IEndpointConventionBuilder

	// RequireAuthorization adds authorization requirements.
	RequireAuthorization(policies ...string) IEndpointConventionBuilder

	// AllowAnonymous allows anonymous access.
	AllowAnonymous() IEndpointConventionBuilder
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

// WithName sets the endpoint name.
func (b *RouteBuilder) WithName(name string) IEndpointConventionBuilder {
	b.name = name
	return b
}

// WithDisplayName sets the display name.
func (b *RouteBuilder) WithDisplayName(displayName string) IEndpointConventionBuilder {
	b.displayName = displayName
	return b
}

// WithMetadata adds metadata to the endpoint.
func (b *RouteBuilder) WithMetadata(metadata ...interface{}) IEndpointConventionBuilder {
	b.metadata = append(b.metadata, metadata...)
	return b
}

// WithSummary sets the OpenAPI summary.
func (b *RouteBuilder) WithSummary(summary string) IEndpointConventionBuilder {
	b.summary = summary
	return b
}

// WithDescription sets the OpenAPI description.
func (b *RouteBuilder) WithDescription(description string) IEndpointConventionBuilder {
	b.description = description
	return b
}

// WithTags adds OpenAPI tags.
func (b *RouteBuilder) WithTags(tags ...string) IEndpointConventionBuilder {
	b.tags = append(b.tags, tags...)
	return b
}

// RequireAuthorization adds authorization requirements.
func (b *RouteBuilder) RequireAuthorization(policies ...string) IEndpointConventionBuilder {
	b.authPolicies = append(b.authPolicies, policies...)
	b.allowAnonymous = false
	return b
}

// AllowAnonymous allows anonymous access.
func (b *RouteBuilder) AllowAnonymous() IEndpointConventionBuilder {
	b.allowAnonymous = true
	b.authPolicies = nil
	return b
}

// Produces adds a response type to the endpoint.
func Produces[T any](b IEndpointConventionBuilder, statusCode int) IEndpointConventionBuilder {
	// Add response metadata
	return b.WithMetadata(&ResponseMetadata{
		StatusCode: statusCode,
		Type:       new(T),
	})
}

// ProducesProblem adds a problem details response.
func ProducesProblem(b IEndpointConventionBuilder, statusCode int) IEndpointConventionBuilder {
	return b.WithMetadata(&ResponseMetadata{
		StatusCode: statusCode,
		IsProblem:  true,
	})
}

// Accepts adds a request body type to the endpoint.
func Accepts[T any](b IEndpointConventionBuilder, contentType string) IEndpointConventionBuilder {
	return b.WithMetadata(&RequestMetadata{
		ContentType: contentType,
		Type:        new(T),
	})
}

// ResponseMetadata represents response metadata.
type ResponseMetadata struct {
	StatusCode int
	Type       interface{}
	IsProblem  bool
}

// RequestMetadata represents request metadata.
type RequestMetadata struct {
	ContentType string
	Type        interface{}
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
