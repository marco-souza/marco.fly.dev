package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/models"
)

func ordersHandler(c *fiber.Ctx) error {
	db := models.Connect()
	orders := []models.Order{}
	result := db.Preload("Author").Find(&orders)

	fmt.Println(result.RowsAffected)

	props := fiber.Map{"Orders": orders}
	return c.Render("partials/order-list", props, "layouts/empty")
}

func orders(router fiber.Router) {
	router.Get("/orders", ordersHandler)
}
