package swagger

import "github.com/gocrud/csgo/di"

// AddSwaggerGen adds Swagger generation services to the service collection.
// Corresponds to .NET services.AddSwaggerGen().
func AddSwaggerGen(services di.IServiceCollection, configure ...func(*SwaggerGenOptions)) di.IServiceCollection {
	services.Add(func() *SwaggerGenOptions {
		opts := NewSwaggerGenOptions()
		if len(configure) > 0 && configure[0] != nil {
			configure[0](opts)
		}
		return opts
	})
	return services
}
