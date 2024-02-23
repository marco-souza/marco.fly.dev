package pages

import (
	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/config"
	"github.com/marco-souza/marco.fly.dev/internal/constants"
	"github.com/marco-souza/marco.fly.dev/internal/github"
)

type Breadcrumb struct {
	Url, Text, Icon string
}

type dashboardProps struct {
	config.PageParams
	Profile     github.GitHubUser
	Description string
	Logout      string
	Breadcrumbs []Breadcrumb
}

func dashboardHandler(c *fiber.Ctx) error {
	breadcrumbs := []Breadcrumb{
		{"/", "Home", "🏠"},
		{"/app", "Dashboard", "🏂"},
	}

	token := github.AccessToken(c)
	loggedUser := github.User("", token)
	props := dashboardProps{
		config.DefaultPageParams,
		loggedUser,
		loggedUser.Bio,
		cfg.Github.LogoutUrl,
		breadcrumbs,
	}

	return c.Render("dashboard", props)
}
