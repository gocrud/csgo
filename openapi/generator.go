package openapi

import (
	"net/http"
	"reflect"
	"strconv"
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
	// Check if it's our ResponseMetadata or RequestMetadata
	// We need to import the routing package, so we'll use type assertion with interface{}

	// Try to access metadata fields using reflection
	v := reflect.ValueOf(meta)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return
	}

	typeName := v.Type().Name()

	switch typeName {
	case "ResponseMetadata":
		g.addResponse(spec, op, meta)
	case "RequestMetadata":
		g.addRequestBody(spec, op, meta)
	}
}

// addResponse adds a response to the operation.
func (g *Generator) addResponse(spec *Specification, op *Operation, meta interface{}) {
	v := reflect.ValueOf(meta)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	statusCode := int(v.FieldByName("StatusCode").Int())
	isProblem := v.FieldByName("IsProblem").Bool()
	typeField := v.FieldByName("Type")

	statusCodeStr := strings.TrimSpace(strings.Split(strings.TrimPrefix(http.StatusText(statusCode), "HTTP "), " ")[0])
	if statusCodeStr == "" {
		statusCodeStr = "200"
	}
	// Convert status code to string
	statusCodeStr = ""
	if statusCode >= 100 && statusCode < 600 {
		statusCodeStr = strconv.Itoa(statusCode)
	} else {
		statusCodeStr = "200"
	}

	if isProblem {
		op.Responses[statusCodeStr] = Response{
			Description: getDefaultResponseDescription(statusCode),
			Content: map[string]MediaType{
				"application/problem+json": {
					Schema: g.getProblemDetailsSchema(),
				},
			},
		}
		return
	}

	if typeField.IsValid() && !typeField.IsZero() {
		reflectType := typeField.Interface().(reflect.Type)
		schema := g.generateSchemaFromReflectType(spec, reflectType)
		op.Responses[statusCodeStr] = Response{
			Description: getDefaultResponseDescription(statusCode),
			Content: map[string]MediaType{
				"application/json": {Schema: schema},
			},
		}
	}
}

// addRequestBody adds a request body to the operation.
func (g *Generator) addRequestBody(spec *Specification, op *Operation, meta interface{}) {
	v := reflect.ValueOf(meta)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	contentType := v.FieldByName("ContentType").String()
	typeField := v.FieldByName("Type")

	if typeField.IsValid() && !typeField.IsZero() {
		reflectType := typeField.Interface().(reflect.Type)
		schema := g.generateSchemaFromReflectType(spec, reflectType)
		op.RequestBody = &RequestBody{
			Required: true,
			Content: map[string]MediaType{
				contentType: {Schema: schema},
			},
		}
	}
}

// generateSchemaFromReflectType generates a schema from a reflect.Type.
func (g *Generator) generateSchemaFromReflectType(spec *Specification, t reflect.Type) Schema {
	if t == nil {
		return Schema{Type: "object"}
	}

	switch t.Kind() {
	case reflect.String:
		return Schema{Type: "string"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Schema{Type: "integer"}
	case reflect.Float32, reflect.Float64:
		return Schema{Type: "number"}
	case reflect.Bool:
		return Schema{Type: "boolean"}
	case reflect.Struct:
		return g.generateStructSchema(spec, t)
	case reflect.Slice, reflect.Array:
		elemType := t.Elem()
		itemSchema := g.generateSchemaFromReflectType(spec, elemType)
		return Schema{
			Type:  "array",
			Items: &itemSchema,
		}
	case reflect.Ptr:
		return g.generateSchemaFromReflectType(spec, t.Elem())
	default:
		return Schema{Type: "object"}
	}
}

// generateStructSchema generates a schema for a struct type.
func (g *Generator) generateStructSchema(spec *Specification, t reflect.Type) Schema {
	schemaName := t.Name()
	if schemaName == "" {
		schemaName = "Anonymous"
	}

	// If already exists, return reference
	if _, exists := spec.Components.Schemas[schemaName]; exists {
		return Schema{Ref: "#/components/schemas/" + schemaName}
	}

	// Parse tags
	fieldInfos := ParseStructTags(t)

	properties := make(map[string]Schema)
	required := []string{}

	for fieldName, fieldInfo := range fieldInfos {
		fieldSchema := g.generateFieldSchema(spec, t, fieldInfo)
		properties[fieldName] = fieldSchema

		if fieldInfo.Required {
			required = append(required, fieldName)
		}
	}

	schema := Schema{
		Type:       "object",
		Properties: properties,
	}
	if len(required) > 0 {
		schema.Required = required
	}

	// Add to components
	spec.Components.Schemas[schemaName] = schema

	// Return reference
	return Schema{Ref: "#/components/schemas/" + schemaName}
}

// generateFieldSchema generates a schema for a struct field.
func (g *Generator) generateFieldSchema(spec *Specification, structType reflect.Type, fieldInfo FieldInfo) Schema {
	// Find field type by name
	var fieldType reflect.Type
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
		if jsonName == fieldInfo.Name || field.Name == fieldInfo.Name {
			fieldType = field.Type
			break
		}
	}

	if fieldType == nil {
		return Schema{Type: "string"}
	}

	// Generate base schema from field type
	schema := g.generateSchemaFromReflectType(spec, fieldType)

	// Apply tag information
	if fieldInfo.Description != "" {
		schema.Description = fieldInfo.Description
	}
	if fieldInfo.Format != "" {
		schema.Format = fieldInfo.Format
	}
	if fieldInfo.Example != nil {
		schema.Example = fieldInfo.Example
	}
	if fieldInfo.Minimum != nil {
		schema.Minimum = fieldInfo.Minimum
	}
	if fieldInfo.Maximum != nil {
		schema.Maximum = fieldInfo.Maximum
	}
	if fieldInfo.MinLength != nil {
		schema.MinLength = fieldInfo.MinLength
	}
	if fieldInfo.MaxLength != nil {
		schema.MaxLength = fieldInfo.MaxLength
	}
	if fieldInfo.Pattern != "" {
		schema.Pattern = fieldInfo.Pattern
	}
	if len(fieldInfo.Enum) > 0 {
		schema.Enum = fieldInfo.Enum
	}

	return schema
}

// getProblemDetailsSchema returns the schema for RFC 7807 Problem Details.
func (g *Generator) getProblemDetailsSchema() Schema {
	return Schema{
		Type: "object",
		Properties: map[string]Schema{
			"type":     {Type: "string"},
			"title":    {Type: "string"},
			"status":   {Type: "integer"},
			"detail":   {Type: "string"},
			"instance": {Type: "string"},
		},
	}
}

// getDefaultResponseDescription returns a default description for a status code.
func getDefaultResponseDescription(statusCode int) string {
	switch statusCode {
	case 200:
		return "Success"
	case 201:
		return "Created"
	case 204:
		return "No Content"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 422:
		return "Validation Error"
	case 500:
		return "Internal Server Error"
	default:
		return "Response"
	}
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
