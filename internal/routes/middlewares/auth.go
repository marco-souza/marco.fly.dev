package middlewares

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/marco-souza/marco.fly.dev/internal/config"
)

func MustBeLoged(c *fiber.Ctx) error {
	hasAccessToken := c.Cookies("access_token", "") != ""
	log.Println("Is user logged: ", hasAccessToken)
	if !hasAccessToken {
		return c.Redirect(config.Load().Github.LoginPage, 302)
	}

	return c.Next()
}
