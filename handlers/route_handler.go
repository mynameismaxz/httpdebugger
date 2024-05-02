package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
}

func New() *Handlers {
	return &Handlers{}
}

// function to received method request
func (h *Handlers) MainRouteHandlers(c *fiber.Ctx) error {
	// response to show method, header, query params, body
	headers := c.GetReqHeaders()
	switch c.Method() {
	case fiber.MethodGet:
		return c.Render("views/index", fiber.Map{
			"title": "[GET] / Method",
			"body":  headers,
		})
	case fiber.MethodPost:
		return h.postRouteMethod(c)
	case fiber.MethodPut:
		return nil
	case fiber.MethodDelete:
		return nil
	}

	return nil
}

// handle post method request
func (h *Handlers) postRouteMethod(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"code":    http.StatusOK,
		"message": "Post method request",
	})
}
