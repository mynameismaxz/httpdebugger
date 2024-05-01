package main

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	embeded "github.com/mynameismaxz/httpdebugger"
)

func main() {
	engine := html.NewFileSystem(http.FS(embeded.ViewFs), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// making handler
	app.All("/", func(c *fiber.Ctx) error {
		// response to show method, header, query params, body
		// get header
		headers := c.GetReqHeaders()
		method := c.Method()
		fmt.Println("Method: ", method)

		return c.Render("views/index", fiber.Map{
			"title": "Root pages",
			"body":  headers,
		})
	})

	app.Listen(":3000")
}
