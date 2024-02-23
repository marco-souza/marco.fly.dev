package pages

import (
	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/config"
)

type Breadcrumb struct {
	Url, Text, Icon string
}

type dashboardProps struct {
	config.PageParams
	Breadcrumbs []Breadcrumb
}

func dashboardHandler(c *fiber.Ctx) error {
	breadcrumbs := []Breadcrumb{
		{"/", "Home", "🏠"},
	}

	props := dashboardProps{
		config.DefaultPageParams,
		breadcrumbs,
	}

	return c.Render("dashboard", props)
}
