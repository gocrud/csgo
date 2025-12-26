package web

import (
	"time"

	"github.com/gin-contrib/cors"
)

// CorsOptions 表示 CORS 配置选项。
// 对应 .NET 的 CorsOptions。
type CorsOptions struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// NewCorsOptions 创建默认的 CORS 选项。
func NewCorsOptions() *CorsOptions {
	return &CorsOptions{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		MaxAge:       12 * time.Hour,
	}
}

// UseCors 向应用程序添加 CORS 中间件。
// 对应 .NET 的 app.UseCors()。
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

// AddCors 向服务集合添加 CORS 服务。
// 对应 .NET 的 services.AddCors()。
func (b *WebApplicationBuilder) AddCors(configure ...func(*CorsOptions)) *WebApplicationBuilder {
	// 将 CORS 选项存储在 DI 中
	b.Services.Add(func() *CorsOptions {
		opts := NewCorsOptions()
		if len(configure) > 0 && configure[0] != nil {
			configure[0](opts)
		}
		return opts
	})
	return b
}
