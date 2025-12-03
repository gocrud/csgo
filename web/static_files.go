package web

// StaticFileOptions represents static file serving options.
type StaticFileOptions struct {
	RequestPath string
	FileSystem  string
}

// UseStaticFiles serves static files.
// Corresponds to .NET app.UseStaticFiles().
func (app *WebApplication) UseStaticFiles(configure ...func(*StaticFileOptions)) *WebApplication {
	opts := &StaticFileOptions{
		RequestPath: "/static",
		FileSystem:  "./wwwroot",
	}

	if len(configure) > 0 && configure[0] != nil {
		configure[0](opts)
	}

	app.Engine.Static(opts.RequestPath, opts.FileSystem)
	return app
}

// UseDefaultFiles enables default file mapping.
// Corresponds to .NET app.UseDefaultFiles().
func (app *WebApplication) UseDefaultFiles() *WebApplication {
	app.Engine.StaticFile("/", "./wwwroot/index.html")
	return app
}

