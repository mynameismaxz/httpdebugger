package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()
	// making handler
	app.Get("/", func(c *fiber.Ctx) error {
		// response to show method, header, query params, body
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}
