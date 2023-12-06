package api

import (
	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/models"
)

func ordersHandler(c *fiber.Ctx) error {
	db := models.Connect()
	orders := []models.Order{}
	result := db.Preload("Author").Find(&orders)

	props := fiber.Map{"Orders": orders, "Total": result.RowsAffected}
	return c.Render("partials/order-list", props, "layouts/empty")
}

func orders(router fiber.Router) {
	router.Get("/orders", ordersHandler)
}
