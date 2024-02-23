package api

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/config"
	"github.com/marco-souza/marco.fly.dev/internal/github"
)

var (
	auth            = github.Auth{}
	cfg             = config.Load()
	refreshTokenKey = "refresh_token"
	accessTokenKey  = "access_token"
)

func redirectGithubAuth(c *fiber.Ctx) error {
	origin := c.BaseURL()
	redirectUrl := auth.RedirectLink(origin)

	return c.Redirect(redirectUrl, http.StatusTemporaryRedirect)
}

func callbackGithubAuth(c *fiber.Ctx) error {
	queries := c.Queries()
	cookies := github.AuthCookies{
		Ctx:             c,
		AccessTokenKey:  accessTokenKey,
		RefreshTokenKey: refreshTokenKey,
	}

	code := queries["code"]
	authToken, _ := auth.FetchAuthToken(code)

	log.Print("settings auth cookies")
	cookies.SetAuthCookies(authToken)

	return c.Redirect(cfg.Github.DashboardPage, http.StatusTemporaryRedirect)
}

func logoutGithubAuth(c *fiber.Ctx) error {
	cookies := github.AuthCookies{
		Ctx:             c,
		AccessTokenKey:  accessTokenKey,
		RefreshTokenKey: refreshTokenKey,
	}

	log.Print("logging out")
	cookies.DeleteAuthCookies()

	return c.Redirect("/", http.StatusTemporaryRedirect)
}
