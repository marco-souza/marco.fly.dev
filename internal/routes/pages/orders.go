package pages

import (
	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/config"
	"github.com/marco-souza/marco.fly.dev/internal/models"
)

type ordersProps struct {
	config.PageParams
	Title  string
	Total  int64
	Orders []models.Order
}

func ordersHandler(c *fiber.Ctx) error {
	db := models.Connect()
	orders := []models.Order{}
	result := db.Preload("Author").Find(&orders)

	props := ordersProps{
		PageParams: config.DefaultPageParams,
		Title:      "All Orders",
		Orders:     orders,
		Total:      result.RowsAffected,
	}

	return c.Render("orders", props)
}
