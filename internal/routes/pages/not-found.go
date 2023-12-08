package pages

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/config"
)

func notFoundHandler(c *fiber.Ctx) error {
	log.Println("Page not found")
	return c.Status(fiber.StatusNotFound).Render("404", config.DefaultPageParams)
}
