package openapi

import (
	"reflect"

	"github.com/gocrud/csgo/web/routing"
)

// OptName returns an option that sets the endpoint name.
func OptName(name string) routing.EndpointOption {
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.SetName(name)
		}
		return b
	}
}

// OptSummary returns an option that sets the OpenAPI summary.
func OptSummary(summary string) routing.EndpointOption {
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.SetSummary(summary)
		}
		return b
	}
}

// OptDescription returns an option that sets the OpenAPI description.
func OptDescription(description string) routing.EndpointOption {
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.SetDescription(description)
		}
		return b
	}
}

// OptTags returns an option that adds OpenAPI tags.
func OptTags(tags ...string) routing.EndpointOption {
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.AddTags(tags...)
		}
		return b
	}
}

// OptResponse returns an option that adds a response type (generic).
func OptResponse[T any](statusCode int) routing.EndpointOption {
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.AddResponseMetadata(routing.ResponseMetadata{
				StatusCode: statusCode,
				Type:       reflect.TypeOf((*T)(nil)).Elem(),
			})
		}
		return b
	}
}

// OptResponseProblem returns an option that adds a problem details response.
func OptResponseProblem(statusCode int) routing.EndpointOption {
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.AddResponseMetadata(routing.ResponseMetadata{
				StatusCode: statusCode,
				IsProblem:  true,
			})
		}
		return b
	}
}

// OptResponseValidationProblem returns an option that adds a validation problem response (422).
func OptResponseValidationProblem() routing.EndpointOption {
	return OptResponseProblem(422)
}

// OptRequest returns an option that adds a request body type (generic).
func OptRequest[T any](contentType ...string) routing.EndpointOption {
	ct := "application/json"
	if len(contentType) > 0 && contentType[0] != "" {
		ct = contentType[0]
	}
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.AddRequestMetadata(routing.RequestMetadata{
				ContentType: ct,
				Type:        reflect.TypeOf((*T)(nil)).Elem(),
			})
		}
		return b
	}
}

// OptAuthorization returns an option that adds authorization requirements.
func OptAuthorization(policies ...string) routing.EndpointOption {
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.SetAuthorizationPolicies(policies)
		}
		return b
	}
}

// OptAnonymous returns an option that allows anonymous access.
func OptAnonymous() routing.EndpointOption {
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.SetAllowAnonymous(true)
		}
		return b
	}
}

// OptParam is a generic option for adding parameters.
func OptParam[T any](name, in, description string, required bool) routing.EndpointOption {
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.AddParameterMetadata(routing.ParameterMetadata{
				Name:        name,
				In:          in,
				Description: description,
				Required:    required,
				Type:        reflect.TypeOf((*T)(nil)).Elem(),
			})
		}
		return b
	}
}

// OptQuery defines a query parameter.
func OptQuery[T any](name, description string) routing.EndpointOption {
	return OptParam[T](name, "query", description, false)
}

// OptQueryRequired defines a required query parameter.
func OptQueryRequired[T any](name, description string) routing.EndpointOption {
	return OptParam[T](name, "query", description, true)
}

// OptPath defines a path parameter.
func OptPath[T any](name, description string) routing.EndpointOption {
	return OptParam[T](name, "path", description, true)
}

// OptHeader defines a header parameter.
func OptHeader[T any](name, description string) routing.EndpointOption {
	return OptParam[T](name, "header", description, false)
}

// OptCookie defines a cookie parameter.
func OptCookie[T any](name, description string) routing.EndpointOption {
	return OptParam[T](name, "cookie", description, false)
}
