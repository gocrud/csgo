package swagger

import "github.com/gocrud/csgo/openapi"

// SwaggerGenOptions configures Swagger generation.
type SwaggerGenOptions struct {
	Title               string
	Version             string
	Description         string
	SecurityDefinitions map[string]openapi.SecurityScheme
}

// NewSwaggerGenOptions creates a new SwaggerGenOptions with defaults.
func NewSwaggerGenOptions() *SwaggerGenOptions {
	return &SwaggerGenOptions{
		Title:               "API",
		Version:             "v1",
		SecurityDefinitions: make(map[string]openapi.SecurityScheme),
	}
}

// AddSecurityDefinition adds a security scheme definition.
func (o *SwaggerGenOptions) AddSecurityDefinition(name string, scheme openapi.SecurityScheme) {
	if o.SecurityDefinitions == nil {
		o.SecurityDefinitions = make(map[string]openapi.SecurityScheme)
	}
	o.SecurityDefinitions[name] = scheme
}

// SwaggerUIOptions configures Swagger UI.
type SwaggerUIOptions struct {
	RoutePrefix string
	SpecURL     string
	Title       string
}

// NewSwaggerUIOptions creates a new SwaggerUIOptions with defaults.
func NewSwaggerUIOptions() *SwaggerUIOptions {
	return &SwaggerUIOptions{
		RoutePrefix: "/swagger",
		SpecURL:     "/swagger/v1/swagger.json",
		Title:       "Swagger UI",
	}
}
