package main

import (
	"github.com/gocrud/csgo/openapi"
	"github.com/gocrud/csgo/swagger"
	"github.com/gocrud/csgo/web"
)

func main() {
	builder := web.CreateBuilder()

	swagger.AddSwaggerGen(builder.Services, func(sgo *swagger.SwaggerGenOptions) {
		sgo.Title = "CSGO Example"
		sgo.Version = "1.0.0"
		sgo.Description = "CSGO Example API"
		sgo.AddSecurityDefinition("Bearer", openapi.SecurityScheme{
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		})
		sgo.AddSecurityDefinition("Basic", openapi.SecurityScheme{
			Type:   "http",
			Scheme: "basic",
		})
	})
	app := builder.Build()

	swagger.UseSwagger(app)
	swagger.UseSwaggerUI(app)

	app.MapGet("/", Hello).WithOpenApi(
		openapi.OptSummary("Hello"),
		openapi.OptApiResponse[HelloResponse](),
		openapi.OptApiErrorResponse(400),
		openapi.OptApiErrorResponse(500),
		openapi.OptApiAuth("Bearer"),
	)

	app.Run("9998")
}

func Hello(c *web.HttpContext) web.IActionResult {
	response := HelloResponse{
		Message: "Hello, CSGO!",
	}
	return c.Ok(response)
}

type HelloResponse struct {
	Message string `json:"message"`
}
