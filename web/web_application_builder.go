package web

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocrud/csgo/config"
	"github.com/gocrud/csgo/di"
	"github.com/gocrud/csgo/hosting"
	"github.com/gocrud/csgo/web/router"
)

// WebApplicationBuilder is a builder for web applications.
// Corresponds to .NET WebApplicationBuilder.
type WebApplicationBuilder struct {
	Services      di.IServiceCollection
	Configuration config.IConfigurationManager
	Environment   hosting.IHostEnvironment
	Host          *ConfigureHostBuilder
	WebHost       *ConfigureWebHostBuilder

	hostBuilder *hosting.HostBuilder
}

// CreateBuilder creates a new web application builder.
// Corresponds to .NET WebApplication.CreateBuilder(args).
func CreateBuilder(args ...string) *WebApplicationBuilder {
	// Create internal HostBuilder (like .NET's HostApplicationBuilder)
	hostBuilder := hosting.CreateDefaultBuilder(args...)

	// Create WebApplicationBuilder with references to HostBuilder's properties
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

	// Get listen address from configuration (set by WebHost.UseUrls)
	addr := b.getListenAddress()

	// Create a shared pointer for runtime URLs
	runtimeUrls := &[]string{}

	// Register HttpServer as hosted service
	b.Services.AddHostedService(func() hosting.IHostedService {
		return NewHttpServer(addr, engine, func() []string {
			return *runtimeUrls
		})
	})

	// Build host using internal HostBuilder (like .NET's approach)
	host := b.hostBuilder.Build()

	// Get the service provider
	services := host.Services()

	// Create web application with shared URL pointer and handler converters
	app := &WebApplication{
		host:        host,
		engine:      engine,
		Services:    services,
		Environment: b.Environment,
		routes:      make([]*router.RouteBuilder, 0),
		groups:      make([]*router.RouteGroupBuilder, 0),
		runtimeUrls: runtimeUrls, // Shared pointer

		// Initialize handler converters with services injection
		toHandler:  MakeToGinHandler(services),
		toHandlers: MakeToGinHandlers(services),
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

// UseShutdownTimeout configures the shutdown timeout.
// Corresponds to .NET builder.WebHost.UseShutdownTimeout().
func (c *ConfigureWebHostBuilder) UseShutdownTimeout(seconds int) *ConfigureWebHostBuilder {
	c.builder.Configuration.Set("server.shutdownTimeout", strconv.Itoa(seconds))
	return c
}

// HttpServer is a hosted service that runs the HTTP server.
type HttpServer struct {
	*hosting.BackgroundService
	defaultAddr string
	getUrls     func() []string // Function to get runtime URLs
	engine      *gin.Engine
}

// NewHttpServer creates a new HTTP server.
func NewHttpServer(addr string, engine *gin.Engine, getUrls func() []string) *HttpServer {
	server := &HttpServer{
		BackgroundService: hosting.NewBackgroundService(),
		defaultAddr:       addr,
		getUrls:           getUrls,
		engine:            engine,
	}
	server.SetExecuteFunc(server.executeAsync)
	return server
}

func (s *HttpServer) executeAsync(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		// Get actual listen address (runtime URLs override default)
		addr := s.getListenAddr()

		// Check if port is available and find an alternative if needed
		originalAddr := addr
		addr, portChanged := s.ensurePortAvailable(addr)

		displayAddr := addr
		if strings.HasPrefix(addr, ":") {
			displayAddr = "http://localhost" + addr
		} else if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
			displayAddr = "http://" + addr
		}

		fmt.Println("========================================")
		fmt.Println("🚀 Web Application Started")
		fmt.Println("========================================")
		if portChanged {
			originalPort := extractPort(originalAddr)
			newPort := extractPort(addr)
			fmt.Printf("⚠️  Port %s is already in use, using port %s instead\n", originalPort, newPort)
		}
		fmt.Printf("📍 Listening on: %s\n", displayAddr)

		// Check if Swagger UI is registered and print the URL
		if s.hasSwaggerRoute() {
			swaggerURL := displayAddr + "/swagger"
			fmt.Printf("📚 Swagger UI: %s\n", swaggerURL)
		}

		fmt.Println("========================================")
		fmt.Println("")

		if err := s.engine.Run(addr); err != nil {
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

// getListenAddr returns the actual listen address (runtime URLs override default).
func (s *HttpServer) getListenAddr() string {
	// Check if runtime URLs are provided
	if s.getUrls != nil {
		urls := s.getUrls()
		if len(urls) > 0 {
			// Use first runtime URL
			return parseListenAddress(urls[0])
		}
	}
	// Fall back to default address
	return s.defaultAddr
}

// ensurePortAvailable checks if the port is available, and finds an alternative if not.
// Returns the final address and a boolean indicating if the port was changed.
func (s *HttpServer) ensurePortAvailable(addr string) (string, bool) {
	// Extract port from address
	port := extractPort(addr)
	if port == "" {
		return addr, false
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return addr, false
	}

	// Check if current port is available
	if isPortAvailable(portNum) {
		return addr, false
	}

	// Find next available port
	newPort, err := findAvailablePort(portNum+1, 100)
	if err != nil {
		// If we can't find an available port, return original and let it fail with proper error
		return addr, false
	}

	// Replace port in address
	newAddr := replacePort(addr, strconv.Itoa(newPort))
	return newAddr, true
}

// hasSwaggerRoute checks if Swagger routes are registered in the engine.
func (s *HttpServer) hasSwaggerRoute() bool {
	routes := s.engine.Routes()
	for _, route := range routes {
		if strings.HasPrefix(route.Path, "/swagger") {
			return true
		}
	}
	return false
}

// isPortAvailable checks if a port is available for listening.
func isPortAvailable(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

// findAvailablePort finds the next available port starting from startPort.
func findAvailablePort(startPort int, maxAttempts int) (int, error) {
	for i := 0; i < maxAttempts; i++ {
		port := startPort + i
		if port > 65535 {
			break
		}
		if isPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found after %d attempts", maxAttempts)
}

// extractPort extracts the port number from an address string.
func extractPort(addr string) string {
	// Handle :port format
	if strings.HasPrefix(addr, ":") {
		return addr[1:]
	}

	// Handle host:port format
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[idx+1:]
	}

	return ""
}

// replacePort replaces the port in an address string.
func replacePort(addr string, newPort string) string {
	// Handle :port format
	if strings.HasPrefix(addr, ":") {
		return ":" + newPort
	}

	// Handle host:port format
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx+1] + newPort
	}

	// If no port found, append it
	return addr + ":" + newPort
}

// getListenAddress returns the listen address from configuration or default.
// This is an internal method used by Build().
func (b *WebApplicationBuilder) getListenAddress() string {
	// Get URLs from configuration (set by WebHost.UseUrls)
	urls := b.Configuration.Get("server.urls")
	if urls == "" {
		return ":8080" // Default address
	}

	// Parse first URL (支持多个URL用分号分隔)
	urlList := strings.Split(urls, ";")
	if len(urlList) == 0 {
		return ":8080"
	}

	// 解析第一个 URL
	firstUrl := strings.TrimSpace(urlList[0])
	return parseListenAddress(firstUrl)
}

// parseListenAddress extracts the listen address from a URL.
// Examples:
//
//	"http://localhost:5000" -> "localhost:5000"
//	"https://0.0.0.0:8443" -> "0.0.0.0:8443"
//	":8080" -> ":8080"
//	"5000" -> ":5000"
func parseListenAddress(urlStr string) string {
	// 如果已经是 :port 格式，直接返回
	if strings.HasPrefix(urlStr, ":") {
		return urlStr
	}

	// 如果是纯数字端口
	if _, err := strconv.Atoi(urlStr); err == nil {
		return ":" + urlStr
	}

	// 解析完整 URL
	if strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://") {
		// 移除协议
		urlStr = strings.TrimPrefix(urlStr, "http://")
		urlStr = strings.TrimPrefix(urlStr, "https://")

		// 如果没有端口，添加默认端口
		if !strings.Contains(urlStr, ":") {
			return urlStr + ":80"
		}
		return urlStr
	}

	// 默认当作 host:port
	if !strings.Contains(urlStr, ":") {
		return urlStr + ":80"
	}
	return urlStr
}
