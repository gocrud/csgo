package main

import (
	"github.com/gocrud/csgo/web"
)

func main() {
	builder := web.CreateBuilder()

	app := builder.Build()

	app.MapGet("/", func(c *web.HttpContext) web.IActionResult {
		return c.Ok(web.M{"message": "Hello, CSGO!"})
	})

	app.Run()
}
