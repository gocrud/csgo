package swagger

import (
	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/openapi"
	"github.com/gocrud/csgo/web"
)

// UseSwagger enables the Swagger JSON endpoint.
// Corresponds to .NET app.UseSwagger().
func UseSwagger(app *web.WebApplication) {
	// Get Swagger options from DI
	opts := di.GetOr[*SwaggerGenOptions](app.Services, NewSwaggerGenOptions())

	// Create OpenAPI generator
	generator := openapi.NewGenerator(opts.Title, opts.Version).
		WithDescription(opts.Description)

	// Add security schemes
	for name, scheme := range opts.SecurityDefinitions {
		generator.WithSecurityScheme(name, scheme)
	}

	// Register Swagger JSON endpoint
	app.MapGet("/swagger/v1/swagger.json", func(c *gin.Context) {
		routes := app.GetRoutes()
		routeInfos := make([]openapi.RouteInfo, len(routes))
		for i, r := range routes {
			routeInfos[i] = r
		}

		spec := generator.Generate(routeInfos)
		c.JSON(200, spec)
	})
}

// UseSwaggerUI enables the Swagger UI.
// Corresponds to .NET app.UseSwaggerUI().
func UseSwaggerUI(app *web.WebApplication, configure ...func(*SwaggerUIOptions)) {
	opts := NewSwaggerUIOptions()
	if len(configure) > 0 && configure[0] != nil {
		configure[0](opts)
	}

	// Register Swagger UI endpoints
	handler := func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, getSwaggerUIHTML(opts))
	}

	app.MapGet(opts.RoutePrefix+"/index.html", handler)
	app.MapGet(opts.RoutePrefix+"/", handler)
	app.MapGet(opts.RoutePrefix, func(c *gin.Context) {
		c.Redirect(301, opts.RoutePrefix+"/index.html")
	})
}

// getSwaggerUIHTML returns the Swagger UI HTML.
func getSwaggerUIHTML(opts *SwaggerUIOptions) string {
	return `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>` + opts.Title + `</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5/favicon-32x32.png" sizes="32x32" />
  <style>
    html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin:0; background: #fafafa; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js" charset="UTF-8"> </script>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js" charset="UTF-8"> </script>
  <script>
    window.onload = function() {
      const ui = SwaggerUIBundle({
        url: "` + opts.SpecURL + `",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      });
      window.ui = ui;
    };
  </script>
</body>
</html>
`
}
