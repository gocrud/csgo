package web

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/configuration"
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/hosting"
	"github.com/gocrud/csgo/web/routing"
)

// WebApplicationBuilder is a builder for web applications.
// Corresponds to .NET WebApplicationBuilder.
type WebApplicationBuilder struct {
	// Public properties (exposed like .NET)
	Services      di.IServiceCollection
	Configuration configuration.IConfiguration
	Environment   hosting.IHostEnvironment
	Host          *ConfigureHostBuilder
	WebHost       *ConfigureWebHostBuilder

	hostBuilder *hosting.HostBuilder
}

// CreateBuilder creates a new web application builder.
// Corresponds to .NET WebApplication.CreateBuilder(args).
func CreateBuilder(args ...string) *WebApplicationBuilder {
	// Create generic host builder
	hostBuilder := hosting.CreateDefaultBuilder(args...)

	builder := &WebApplicationBuilder{
		Services:      hostBuilder.Services,
		Configuration: hostBuilder.Configuration,
		Environment:   hostBuilder.Environment,
		hostBuilder:   hostBuilder,
	}

	builder.Host = &ConfigureHostBuilder{builder: builder}
	builder.WebHost = &ConfigureWebHostBuilder{builder: builder}

	return builder
}

// Build builds the web application.
func (b *WebApplicationBuilder) Build() *WebApplication {
	// Configure Gin mode
	gin.SetMode(gin.ReleaseMode)
	if b.Environment.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin engine
	engine := gin.New()

	// Register HttpServer as hosted service
	b.Services.AddHostedService(func() hosting.IHostedService {
		return NewHttpServer(":8080", engine)
	})

	// Build host
	host := b.hostBuilder.Build()

	// Create web application
	app := &WebApplication{
		host:     host,
		Engine:   engine,
		Services: host.Services().(di.IServiceProvider), // ✅ 构建时转换
		routes:   make([]*routing.RouteBuilder, 0),
		groups:   make([]*routing.RouteGroupBuilder, 0),
	}

	return app
}

// ConfigureHostBuilder allows configuring the generic host.
type ConfigureHostBuilder struct {
	builder *WebApplicationBuilder
}

// ConfigureServices configures services for the host.
func (c *ConfigureHostBuilder) ConfigureServices(configure func(di.IServiceCollection)) *ConfigureHostBuilder {
	configure(c.builder.Services)
	return c
}

// ConfigureWebHostBuilder allows configuring the web host.
type ConfigureWebHostBuilder struct {
	builder *WebApplicationBuilder
}

// UseUrls configures the URLs the web server listens on.
func (c *ConfigureWebHostBuilder) UseUrls(urls ...string) *ConfigureWebHostBuilder {
	if len(urls) > 0 {
		c.builder.Configuration.Set("server.urls", strings.Join(urls, ";"))
	}
	return c
}

// HttpServer is a hosted service that runs the HTTP server.
type HttpServer struct {
	*hosting.BackgroundService
	addr   string
	engine *gin.Engine
}

// NewHttpServer creates a new HTTP server.
func NewHttpServer(addr string, engine *gin.Engine) *HttpServer {
	server := &HttpServer{
		BackgroundService: hosting.NewBackgroundService(),
		addr:              addr,
		engine:            engine,
	}
	server.SetExecuteFunc(server.executeAsync)
	return server
}

func (s *HttpServer) executeAsync(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		addr := s.addr
		if strings.HasPrefix(addr, ":") {
			addr = "http://localhost" + addr
		} else if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
			addr = "http://" + addr
		}

		fmt.Println("========================================")
		fmt.Println("🚀 Web Application Started")
		fmt.Println("========================================")
		fmt.Printf("📍 Listening on: %s\n", addr)
		fmt.Println("========================================")
		fmt.Println("")

		if err := s.engine.Run(s.addr); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-s.StoppingToken():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
