package pages

import (
	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/config"
)

func playgroundHandler(c *fiber.Ctx) error {
	return c.Render("playground", config.DefaultPageParams)
}

func playground(router fiber.Router) {
	router.Get("/playground", playgroundHandler)
}
