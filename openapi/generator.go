package openapi

import (
	"reflect"
	"strings"
	"unicode"
)

// Generator generates OpenAPI specifications from route metadata.
type Generator struct {
	title           string
	version         string
	description     string
	servers         []Server
	security        []map[string][]string
	securitySchemes map[string]SecurityScheme
}

// NewGenerator creates a new OpenAPI generator.
func NewGenerator(title, version string) *Generator {
	return &Generator{
		title:           title,
		version:         version,
		servers:         make([]Server, 0),
		security:        make([]map[string][]string, 0),
		securitySchemes: make(map[string]SecurityScheme),
	}
}

// WithDescription sets the API description.
func (g *Generator) WithDescription(description string) *Generator {
	g.description = description
	return g
}

// WithServer adds a server to the specification.
func (g *Generator) WithServer(url, description string) *Generator {
	g.servers = append(g.servers, Server{
		URL:         url,
		Description: description,
	})
	return g
}

// WithBearerAuth adds bearer authentication.
func (g *Generator) WithBearerAuth() *Generator {
	g.securitySchemes["Bearer"] = SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	}
	return g
}

// WithSecurityScheme adds a custom security scheme.
func (g *Generator) WithSecurityScheme(name string, scheme SecurityScheme) *Generator {
	g.securitySchemes[name] = scheme
	return g
}

// RouteInfo represents route information for OpenAPI generation.
type RouteInfo interface {
	GetMethod() string
	GetPath() string
	GetName() string
	GetSummary() string
	GetDescription() string
	GetTags() []string
	GetMetadata() []interface{}
	IsOpenApiEnabled() bool
}

// Generate generates the OpenAPI specification.
func (g *Generator) Generate(routes []RouteInfo) *Specification {
	spec := &Specification{
		OpenAPI: "3.0.3",
		Info: Info{
			Title:       g.title,
			Version:     g.version,
			Description: g.description,
		},
		Servers: g.servers,
		Paths:   make(map[string]PathItem),
		Components: Components{
			Schemas:         make(map[string]Schema),
			SecuritySchemes: g.securitySchemes,
		},
		Security: g.security,
	}

	// Only include routes that have explicitly enabled OpenAPI documentation
	// Corresponds to .NET behavior where endpoints need .WithOpenApi() to be included
	for _, route := range routes {
		if route.IsOpenApiEnabled() {
			g.addRoute(spec, route)
		}
	}

	return spec
}

// addRoute adds a route to the specification.
func (g *Generator) addRoute(spec *Specification, route RouteInfo) {
	path := normalizePath(route.GetPath())
	method := strings.ToLower(route.GetMethod())

	if spec.Paths[path] == nil {
		spec.Paths[path] = make(PathItem)
	}

	op := Operation{
		Summary:     route.GetSummary(),
		Description: route.GetDescription(),
		OperationID: generateOperationID(method, path),
		Tags:        route.GetTags(),
		Responses:   make(map[string]Response),
	}

	// Process metadata
	for _, meta := range route.GetMetadata() {
		g.processMetadata(spec, &op, meta)
	}

	// Ensure at least one response
	if len(op.Responses) == 0 {
		op.Responses["200"] = Response{Description: "OK"}
	}

	spec.Paths[path][method] = op
}

// processMetadata processes route metadata.
func (g *Generator) processMetadata(spec *Specification, op *Operation, meta interface{}) {
	// This will be implemented based on metadata types
	// For now, skip
}

// normalizePath converts Gin path format to OpenAPI format.
// Example: /users/:id -> /users/{id}
func normalizePath(p string) string {
	parts := strings.Split(p, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			parts[i] = "{" + part[1:] + "}"
		} else if strings.HasPrefix(part, "*") {
			parts[i] = "{" + part[1:] + "}"
		}
	}
	return strings.Join(parts, "/")
}

// generateOperationID creates a camelCase operation ID from method and path.
// Example: GET /users/{id} -> getUsersId
func generateOperationID(method, path string) string {
	var sb strings.Builder
	sb.WriteString(strings.ToLower(method))

	parts := strings.Split(path, "/")
	for _, part := range parts {
		if part == "" {
			continue
		}

		// Handle path params {id} -> Id
		cleanPart := strings.Trim(part, "{}")
		cleanPart = strings.TrimPrefix(cleanPart, ":")

		// Capitalize first letter
		runes := []rune(cleanPart)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		sb.WriteString(string(runes))
	}

	return sb.String()
}

// registerSchema registers a type as a schema in the components.
func registerSchema(spec *Specification, obj interface{}) string {
	t := reflect.TypeOf(obj)
	if t == nil {
		return "Unknown"
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return t.Name()
	}

	name := t.Name()
	if name == "" {
		name = "AnonymousStruct"
	}

	if _, exists := spec.Components.Schemas[name]; exists {
		return name
	}

	schema := Schema{
		Type:       "object",
		Properties: make(map[string]Schema),
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		jsonTag := f.Tag.Get("json")
		fieldName := strings.Split(jsonTag, ",")[0]
		if fieldName == "" {
			fieldName = f.Name
		}
		if fieldName == "-" {
			continue
		}

		fieldType := "string"
		switch f.Type.Kind() {
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Uint, reflect.Uint64:
			fieldType = "integer"
		case reflect.Float32, reflect.Float64:
			fieldType = "number"
		case reflect.Bool:
			fieldType = "boolean"
		}

		schema.Properties[fieldName] = Schema{
			Type:        fieldType,
			Description: f.Tag.Get("description"),
		}
	}

	spec.Components.Schemas[name] = schema
	return name
}
