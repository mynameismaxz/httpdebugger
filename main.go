package main

import (
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

//go:embed views/*
var viewsFs embed.FS

func main() {
	engine := html.NewFileSystem(http.FS(viewsFs), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// making handler
	app.Get("/", func(c *fiber.Ctx) error {
		// response to show method, header, query params, body
		// get header
		headers := c.GetReqHeaders()

		return c.Render("views/index", fiber.Map{
			"title": "Root pages",
			"body":  headers,
		})
	})

	app.Listen(":3000")
}
