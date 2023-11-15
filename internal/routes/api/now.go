package api

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func nowHandler(c *fiber.Ctx) error {
	now := time.Now().Format(time.RFC1123Z)
	log.Println("Now is", now)

	props := fiber.Map{"Time": now}

	return c.Render("partials/now", props, "layouts/empty")
}

func now(router fiber.Router) {
	router.Get("/now", nowHandler)
}
