package middlewares

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/marco-souza/marco.fly.dev/internal/config"
	"github.com/marco-souza/marco.fly.dev/internal/github"
)

var logger = slog.With("service", "middleware")

func MustBeLoged(c *fiber.Ctx) error {
	hasAccessToken := github.HasAccessToken(c)
	logger.Info("is user logged: ", "hasAccessToken", hasAccessToken)

	if !hasAccessToken {
		return c.Redirect(
			config.Load().Github.LoginPage,
			fiber.StatusTemporaryRedirect,
		)
	}

	return c.Next()
}

// create auth check middleware
func MustHaveToken(c *fiber.Ctx) error {
	token := github.AccessToken(c)
	if token == "" {
		logger.Warn("unauthorized access, redirecting to login page")
		return c.Redirect(config.Load().Github.LoginPage)
	}
	return c.Next()
}
