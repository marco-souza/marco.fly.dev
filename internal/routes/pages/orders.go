package pages

import (
	"github.com/gofiber/fiber/v2"
)

func ordersHandler(c *fiber.Ctx) error {
	props := fiber.Map{
		"Title": "Orders List",
	}
	return c.Render("orders", props)
}

func orders(router fiber.Router) {
	router.Get("/orders", ordersHandler)
}
