package main

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	embeded "github.com/mynameismaxz/httpdebugger"
	"github.com/mynameismaxz/httpdebugger/handlers"
)

func main() {
	engine := html.NewFileSystem(http.FS(embeded.ViewFs), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	h := handlers.New()

	// making handler
	app.All("/", h.MainRouteHandlers)
	app.Listen(":3000")
}
