package web

import (
	"time"

	"github.com/gin-contrib/cors"
)

// CorsOptions represents CORS configuration options.
// Corresponds to .NET CorsOptions.
type CorsOptions struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// NewCorsOptions creates default CORS options.
func NewCorsOptions() *CorsOptions {
	return &CorsOptions{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		MaxAge:       12 * time.Hour,
	}
}

// UseCors adds CORS middleware to the application.
// Corresponds to .NET app.UseCors().
func (app *WebApplication) UseCors(configure ...func(*CorsOptions)) *WebApplication {
	opts := NewCorsOptions()
	if len(configure) > 0 && configure[0] != nil {
		configure[0](opts)
	}

	config := cors.Config{
		AllowOrigins:     opts.AllowOrigins,
		AllowMethods:     opts.AllowMethods,
		AllowHeaders:     opts.AllowHeaders,
		ExposeHeaders:    opts.ExposeHeaders,
		AllowCredentials: opts.AllowCredentials,
		MaxAge:           opts.MaxAge,
	}

	app.Use(cors.New(config))
	return app
}

// AddCors adds CORS services to the service collection.
// Corresponds to .NET services.AddCors().
func (b *WebApplicationBuilder) AddCors(configure ...func(*CorsOptions)) *WebApplicationBuilder {
	// Store CORS options in DI
	b.Services.AddSingleton(func() *CorsOptions {
		opts := NewCorsOptions()
		if len(configure) > 0 && configure[0] != nil {
			configure[0](opts)
		}
		return opts
	})
	return b
}

