package pages

import (
	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/config"
	"github.com/marco-souza/marco.fly.dev/internal/github"
	"github.com/marco-souza/marco.fly.dev/internal/models"
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
	loggedUser, _ := github.User("", token)
	// TODO: cache user info

	props := dashboardProps{
		config.DefaultPageParams,
		*loggedUser,
		loggedUser.Bio,
		cfg.Github.LogoutUrl,
		breadcrumbs,
	}

	return c.Render("dashboard", props)
}

// user playground
func playgroundHandler(c *fiber.Ctx) error {
	return c.Render("playground", config.DefaultPageParams)
}

type ordersProps struct {
	config.PageParams
	Title  string
	Total  int64
	Orders []models.Order
}

// create user orders
func ordersHandler(c *fiber.Ctx) error {
	db := models.Connect()
	orders := []models.Order{}
	result := db.Preload("Author").Find(&orders)

	props := ordersProps{
		PageParams: config.DefaultPageParams,
		Title:      "All Orders",
		Orders:     orders,
		Total:      result.RowsAffected,
	}

	return c.Render("orders", props)
}

// create cronjobs
func cronHandler(c *fiber.Ctx) error {
	return c.Render("cronjobs", config.DefaultPageParams)
}
