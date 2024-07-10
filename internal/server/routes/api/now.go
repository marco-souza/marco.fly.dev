package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func nowHandler(c *fiber.Ctx) error {
	now := time.Now().Format(time.RFC1123Z)
	props := fiber.Map{"Time": now}
	return c.Render("partials/now", props, "layouts/empty")
}
