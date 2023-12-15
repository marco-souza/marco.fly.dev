package pages

import (
	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/config"
)

var cfg = config.Load()

type loginProps struct {
	config.PageParams
	SignInUrl string
}

func loginHandler(c *fiber.Ctx) error {
	return c.Render("login", loginProps{
		PageParams: config.DefaultPageParams,
		SignInUrl:  cfg.Github.SignInUrl,
	})
}
