package openapi

import (
	"reflect"

	"github.com/gocrud/csgo/web/routing"
)

// Name returns an option that sets the endpoint name.
func Name(name string) routing.EndpointOption {
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.SetName(name)
		}
		return b
	}
}

// Summary returns an option that sets the OpenAPI summary.
func Summary(summary string) routing.EndpointOption {
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.SetSummary(summary)
		}
		return b
	}
}

// Description returns an option that sets the OpenAPI description.
func Description(description string) routing.EndpointOption {
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.SetDescription(description)
		}
		return b
	}
}

// Tags returns an option that adds OpenAPI tags.
func Tags(tags ...string) routing.EndpointOption {
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.AddTags(tags...)
		}
		return b
	}
}

// Produces returns an option that adds a response type (generic).
func Produces[T any](statusCode int) routing.EndpointOption {
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

// ProducesProblem returns an option that adds a problem details response.
func ProducesProblem(statusCode int) routing.EndpointOption {
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

// ProducesValidationProblem returns an option that adds a validation problem response (422).
func ProducesValidationProblem() routing.EndpointOption {
	return ProducesProblem(422)
}

// Accepts returns an option that adds a request body type (generic).
func Accepts[T any](contentType string) routing.EndpointOption {
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.AddRequestMetadata(routing.RequestMetadata{
				ContentType: contentType,
				Type:        reflect.TypeOf((*T)(nil)).Elem(),
			})
		}
		return b
	}
}

// Authorization returns an option that adds authorization requirements.
func Authorization(policies ...string) routing.EndpointOption {
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.SetAuthorizationPolicies(policies)
		}
		return b
	}
}

// Anonymous returns an option that allows anonymous access.
func Anonymous() routing.EndpointOption {
	return func(b routing.IEndpointConventionBuilder) routing.IEndpointConventionBuilder {
		if rb, ok := b.(*routing.RouteBuilder); ok {
			rb.SetAllowAnonymous(true)
		}
		return b
	}
}
