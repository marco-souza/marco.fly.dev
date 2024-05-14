package middlewares

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/marco-souza/marco.fly.dev/internal/config"
	"github.com/marco-souza/marco.fly.dev/internal/github"
)

func MustBeLoged(c *fiber.Ctx) error {
	hasAccessToken := github.HasAccessToken(c)
	log.Println("Is user logged: ", hasAccessToken)

	if !hasAccessToken {
		return c.Redirect(
			config.Load().Github.LoginPage,
			fiber.StatusTemporaryRedirect,
		)
	}

	return c.Next()
}
