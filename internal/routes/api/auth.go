package api

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/github"
)

var auth = github.Auth{}

func redirectGithubAuth(c *fiber.Ctx) error {
	origin := c.BaseURL()
	redirectUrl := auth.RedirectLink(origin)
	return c.Redirect(redirectUrl, 302)
}

func callbackGithubAuth(c *fiber.Ctx) error {
	queries := c.Queries()

	code := queries["code"]
	accessToken, _ := auth.FetchAuthToken(code)

	log.Print(fiber.Map{
		"a": accessToken,
	})

	return nowHandler(c)
}
