package pages

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/config"
)

var cfg = config.Load()

type loginProps struct {
	config.PageParams
	SignInUrl string
}

func loginHandler(c *fiber.Ctx) error {
	if c.Cookies("access_token", "") == "" {
		return c.Redirect(cfg.Github.DashboardPage, http.StatusTemporaryRedirect)
	}
	return c.Render("login", loginProps{
		PageParams: config.DefaultPageParams,
		SignInUrl:  cfg.Github.SignInUrl,
	})
}
