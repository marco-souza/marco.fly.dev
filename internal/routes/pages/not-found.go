package pages

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/config"
)

func notFound(router fiber.Router) {
	router.Use(func(c *fiber.Ctx) error {
		log.Println("Page not found")
		return c.Status(fiber.StatusNotFound).Render("404", config.DefaultPageParams)
	})
}
