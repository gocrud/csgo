package router

import (
	"reflect"

	"github.com/gocrud/csgo/openapi"
)

// ResponseMetadata represents response metadata.
type ResponseMetadata struct {
	StatusCode      int
	Type            reflect.Type
	ContentType     string          // Custom response content type (e.g., image/png)
	Format          string          // OpenAPI format (byte, binary)
	Schema          *openapi.Schema // Manually defined Schema
	IsProblem       bool
	IsApiResponse   bool // Indicates if the response should be wrapped in web.ApiResponse
	IsErrorResponse bool // Indicates if this is an error response (only error field populated)
}

// RequestMetadata represents request metadata.
type RequestMetadata struct {
	ContentType string
	Type        reflect.Type
	Schema      *openapi.Schema // Manually defined Schema
}

// ParameterMetadata represents parameter metadata.
type ParameterMetadata struct {
	Name        string
	In          string // path, query, header, cookie
	Description string
	Required    bool
	Type        reflect.Type
}
