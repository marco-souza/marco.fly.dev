package pages

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/config"
)

func notFoundHandler(c *fiber.Ctx) error {
	log.Println("Page not found")
	return c.
		Status(http.StatusNotFound).
		Render("404", config.DefaultPageParams)
}
