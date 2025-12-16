package router

import (
	"reflect"
	"strings"

	"github.com/gocrud/csgo/openapi"
)

// TypeOf returns the reflect.Type of T without instantiation.
// This is a helper function for use with OpenApiBuilder methods.
// Example: api.Body(router.TypeOf[CreateUserRequest]())
func TypeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// OpenApiBuilder provides a fluent API for configuring OpenAPI metadata.
type OpenApiBuilder struct {
	builder *RouteBuilder
}

// Name sets the endpoint name.
func (o *OpenApiBuilder) Name(name string) *OpenApiBuilder {
	o.builder.SetName(name)
	return o
}

// Summary sets the OpenAPI summary.
func (o *OpenApiBuilder) Summary(summary string) *OpenApiBuilder {
	o.builder.SetSummary(summary)
	return o
}

// Description sets the OpenAPI description.
func (o *OpenApiBuilder) Description(description string) *OpenApiBuilder {
	o.builder.SetDescription(description)
	return o
}

// Tags adds OpenAPI tags.
func (o *OpenApiBuilder) Tags(tags ...string) *OpenApiBuilder {
	o.builder.AddTags(tags...)
	return o
}

// Body sets the request body type.
// Use router.TypeOf[T]() to specify the type.
// Example: api.Body(router.TypeOf[CreateUserRequest]())
func (o *OpenApiBuilder) Body(typ reflect.Type, contentType ...string) *OpenApiBuilder {
	ct := "application/json"
	if len(contentType) > 0 && contentType[0] != "" {
		ct = contentType[0]
	}
	o.builder.AddRequestMetadata(RequestMetadata{
		ContentType: ct,
		Type:        typ,
	})
	return o
}

// BodySchema uses a custom schema for request body.
func (o *OpenApiBuilder) BodySchema(schema openapi.Schema, contentType ...string) *OpenApiBuilder {
	ct := "application/json"
	if len(contentType) > 0 && contentType[0] != "" {
		ct = contentType[0]
	}
	schemaCopy := schema
	o.builder.AddRequestMetadata(RequestMetadata{
		ContentType: ct,
		Schema:      &schemaCopy,
	})
	return o
}

// Response adds a response type.
// Use router.TypeOf[T]() to specify the type.
// Example: api.Response(router.TypeOf[User](), 201)
func (o *OpenApiBuilder) Response(typ reflect.Type, statusCode ...int) *OpenApiBuilder {
	code := 200
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	o.builder.AddResponseMetadata(ResponseMetadata{
		StatusCode: code,
		Type:       typ,
	})
	return o
}

// ApiResponse adds an ApiResponse-wrapped response.
// Use router.TypeOf[T]() to specify the type.
// Example: api.ApiResponse(router.TypeOf[User](), 200)
func (o *OpenApiBuilder) ApiResponse(typ reflect.Type, statusCode ...int) *OpenApiBuilder {
	code := 200
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	o.builder.AddResponseMetadata(ResponseMetadata{
		StatusCode:    code,
		Type:          typ,
		IsApiResponse: true,
	})
	return o
}

// ApiError adds an API error response.
func (o *OpenApiBuilder) ApiError(statusCode int) *OpenApiBuilder {
	o.builder.AddResponseMetadata(ResponseMetadata{
		StatusCode:      statusCode,
		IsApiResponse:   true,
		IsErrorResponse: true,
	})
	return o
}

// ValidationError adds a validation error response (422).
func (o *OpenApiBuilder) ValidationError() *OpenApiBuilder {
	return o.Problem(422)
}

// Problem adds a problem details response.
func (o *OpenApiBuilder) Problem(statusCode int) *OpenApiBuilder {
	o.builder.AddResponseMetadata(ResponseMetadata{
		StatusCode: statusCode,
		IsProblem:  true,
	})
	return o
}

// Query adds a query parameter.
// Use router.TypeOf[T]() to specify the type.
// Example: api.Query(router.TypeOf[int](), "page", "页码", false)
func (o *OpenApiBuilder) Query(typ reflect.Type, name, description string, required bool) *OpenApiBuilder {
	o.builder.AddParameterMetadata(ParameterMetadata{
		Name:        name,
		In:          "query",
		Description: description,
		Required:    required,
		Type:        typ,
	})
	return o
}

// Path adds a path parameter (always required).
// Use router.TypeOf[T]() to specify the type.
// Example: api.Path(router.TypeOf[string](), "id", "用户ID")
func (o *OpenApiBuilder) Path(typ reflect.Type, name, description string) *OpenApiBuilder {
	o.builder.AddParameterMetadata(ParameterMetadata{
		Name:        name,
		In:          "path",
		Description: description,
		Required:    true,
		Type:        typ,
	})
	return o
}

// Header adds a header parameter.
// Use router.TypeOf[T]() to specify the type.
// Example: api.Header(router.TypeOf[string](), "X-Token", "认证令牌", false)
func (o *OpenApiBuilder) Header(typ reflect.Type, name, description string, required bool) *OpenApiBuilder {
	o.builder.AddParameterMetadata(ParameterMetadata{
		Name:        name,
		In:          "header",
		Description: description,
		Required:    required,
		Type:        typ,
	})
	return o
}

// Cookie adds a cookie parameter.
// Use router.TypeOf[T]() to specify the type.
// Example: api.Cookie(router.TypeOf[string](), "session", "会话ID", false)
func (o *OpenApiBuilder) Cookie(typ reflect.Type, name, description string, required bool) *OpenApiBuilder {
	o.builder.AddParameterMetadata(ParameterMetadata{
		Name:        name,
		In:          "cookie",
		Description: description,
		Required:    required,
		Type:        typ,
	})
	return o
}

// Params adds multiple parameters from a struct type.
// Fields should be tagged with in:"<location>" where location is query/header/path/cookie.
// Supported tags:
//   - in: Parameter location (required) - query/header/path/cookie
//   - desc: Parameter description (optional)
//   - required: Whether the parameter is required (optional, default: false, path params are always required)
//   - example: Example value (optional)
//   - enum: Comma-separated enum values (optional)
//   - default: Default value (optional)
//
// Example:
//
//	type SearchParams struct {
//	    CategoryID string `in:"path" desc:"分类ID"`
//	    Keyword    string `in:"query" desc:"搜索关键词" required:"true"`
//	    Page       int    `in:"query" desc:"页码" default:"1"`
//	    SortBy     string `in:"query" desc:"排序字段" enum:"name,price,date"`
//	}
//	api.Params(router.TypeOf[SearchParams]())
func (o *OpenApiBuilder) Params(typ reflect.Type) *OpenApiBuilder {
	if typ.Kind() != reflect.Struct {
		return o
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Parse in tag (parameter location)
		location := field.Tag.Get("in")
		if location == "" {
			continue // Skip fields without in tag
		}

		// Validate location value
		if location != "query" && location != "header" && location != "path" && location != "cookie" {
			continue
		}

		// Get field name (prefer json tag)
		name := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			jsonName := strings.Split(jsonTag, ",")[0]
			if jsonName != "" && jsonName != "-" {
				name = jsonName
			}
		}

		// Parse other tags
		description := field.Tag.Get("desc")
		required := field.Tag.Get("required") == "true"

		// Path parameters are always required
		if location == "path" {
			required = true
		}

		o.builder.AddParameterMetadata(ParameterMetadata{
			Name:        name,
			In:          location,
			Description: description,
			Required:    required,
			Type:        field.Type,
		})
	}

	return o
}

// Auth adds authentication requirement.
func (o *OpenApiBuilder) Auth(name string, scopes ...string) *OpenApiBuilder {
	o.builder.AddApiSecurityRequirement(name, scopes)
	return o
}

// Anonymous allows anonymous access.
func (o *OpenApiBuilder) Anonymous() *OpenApiBuilder {
	o.builder.SetAllowAnonymous(true)
	return o
}

// ResponseSchema uses a custom schema for response.
func (o *OpenApiBuilder) ResponseSchema(schema openapi.Schema, statusCode int, contentType ...string) *OpenApiBuilder {
	ct := "application/json"
	if len(contentType) > 0 && contentType[0] != "" {
		ct = contentType[0]
	}
	schemaCopy := schema
	o.builder.AddResponseMetadata(ResponseMetadata{
		StatusCode:  statusCode,
		ContentType: ct,
		Schema:      &schemaCopy,
	})
	return o
}

// ApiResponseSchema uses a custom schema wrapped in ApiResponse.
func (o *OpenApiBuilder) ApiResponseSchema(schema openapi.Schema, statusCode ...int) *OpenApiBuilder {
	code := 200
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	schemaCopy := schema
	o.builder.AddResponseMetadata(ResponseMetadata{
		StatusCode:    code,
		Schema:        &schemaCopy,
		IsApiResponse: true,
	})
	return o
}

// BinaryImage adds a binary image response.
func (o *OpenApiBuilder) BinaryImage(contentType string, statusCode ...int) *OpenApiBuilder {
	code := 200
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	o.builder.AddResponseMetadata(ResponseMetadata{
		StatusCode:  code,
		ContentType: contentType,
		Format:      "binary",
	})
	return o
}

// Base64Image adds a base64 image response.
func (o *OpenApiBuilder) Base64Image(statusCode ...int) *OpenApiBuilder {
	code := 200
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	o.builder.AddResponseMetadata(ResponseMetadata{
		StatusCode:    code,
		Type:          reflect.TypeOf(""),
		Format:        "byte",
		IsApiResponse: true,
	})
	return o
}
