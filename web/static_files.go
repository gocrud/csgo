package web

// StaticFileOptions 表示静态文件服务选项。
type StaticFileOptions struct {
	RequestPath string
	FileSystem  string
}

// UseStaticFiles 提供静态文件服务。
// 对应 .NET 的 app.UseStaticFiles()。
func (app *WebApplication) UseStaticFiles(configure ...func(*StaticFileOptions)) *WebApplication {
	opts := &StaticFileOptions{
		RequestPath: "/static",
		FileSystem:  "./wwwroot",
	}

	if len(configure) > 0 && configure[0] != nil {
		configure[0](opts)
	}

	app.engine.Static(opts.RequestPath, opts.FileSystem)
	return app
}

// UseDefaultFiles 启用默认文件映射。
// 对应 .NET 的 app.UseDefaultFiles()。
func (app *WebApplication) UseDefaultFiles() *WebApplication {
	app.engine.StaticFile("/", "./wwwroot/index.html")
	return app
}
