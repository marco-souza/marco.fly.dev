package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marco-souza/marco.fly.dev/internal/github"
)

type menuItem struct {
	Href string
	Name string
}

func menuHandler(c *fiber.Ctx) error {
	token := github.AccessToken(c)

	menu := []menuItem{
		{"/", "Home"},
		{"/resume", "Resume"},
		{"/blog", "Blog"},
		{"/login", "Login"},
	}

	if token != "" {
		menu = []menuItem{
			{"/app/", "Dashboard"},
			{"/app/playground", "Playground"},
			{"/app/cronjobs", "Task Scheduler"},
			{"/logout", "Logout"},
		}
	}

	params := fiber.Map{"MenuItems": menu}
	return c.Render("partials/menu", params, "layouts/empty")
}
