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
		{"/login", "Login"},
	}

	if token != "" {
		menu = menu[:len(menu)-1] // remove login item

		menu = append(menu, menuItem{"/app/", "App"})
		menu = append(menu, menuItem{"/app/playground", "Playground"})
		menu = append(menu, menuItem{"/app/orders", "Ordero"})
		menu = append(menu, menuItem{cfg.Github.LogoutUrl, "Logout"})
	}

	params := fiber.Map{"MenuItems": menu}
	return c.Render("partials/menu", params, "layouts/empty")
}
