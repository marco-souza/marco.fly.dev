package pages

import (
	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/models"
)

func ordersHandler(c *fiber.Ctx) error {
	db := models.Connect()
	orders := []models.Order{}
	result := db.Preload("Author").Find(&orders)

	props := fiber.Map{
		"Title":  "All Orders",
		"Orders": orders,
		"Total":  result.RowsAffected,
	}

	return c.Render("orders", props)
}

func orders(router fiber.Router) {
	router.Get("/orders", ordersHandler)
}
