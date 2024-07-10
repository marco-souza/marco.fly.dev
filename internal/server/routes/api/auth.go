package api

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/config"
	"github.com/marco-souza/marco.fly.dev/internal/constants"
	"github.com/marco-souza/marco.fly.dev/internal/github"
)

var (
	auth = github.Auth{
		AllowedEmails: map[string]bool{
			"ma.souza.junior@gmail.com": true,
		},
	}
	logger = slog.With("service", "api")
	cfg    = config.Load()
)

func redirectGithubAuth(c *fiber.Ctx) error {
	origin := c.BaseURL()
	redirectUrl := auth.RedirectLink(origin)

	return c.Redirect(redirectUrl, fiber.StatusTemporaryRedirect)
}

func callbackGithubAuth(c *fiber.Ctx) error {
	queries := c.Queries()
	cookies := github.AuthCookies{
		Ctx:             c,
		AccessTokenKey:  constants.ACCESS_TOKEN_KEY,
		RefreshTokenKey: constants.REFRESH_TOKEN_KEY,
	}

	code := queries["code"]
	authToken, _ := auth.FetchAuthToken(code)

	// check if the user is allowed
	if !auth.IsUserAllowed(authToken.AccessToken) {
		logger.Warn("user is not allowed")
		return c.Redirect(cfg.Github.LoginPage+"?error=unauthorized", fiber.StatusTemporaryRedirect)
	}

	logger.Info("settings auth cookies")
	cookies.SetAuthCookies(authToken)

	return c.Redirect(cfg.Github.DashboardPage, fiber.StatusTemporaryRedirect)
}

func logoutGithubAuth(c *fiber.Ctx) error {
	cookies := github.AuthCookies{
		Ctx:             c,
		AccessTokenKey:  constants.ACCESS_TOKEN_KEY,
		RefreshTokenKey: constants.REFRESH_TOKEN_KEY,
	}

	logger.Info("sign out")
	cookies.DeleteAuthCookies()

	return c.Redirect("/", fiber.StatusTemporaryRedirect)
}
