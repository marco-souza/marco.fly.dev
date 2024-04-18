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
	loggedUser := github.User("", token)
	pageParams := config.MakePageParams(token != "")
	// TODO: cache user info

	props := dashboardProps{
		pageParams,
		loggedUser,
		loggedUser.Bio,
		cfg.Github.LogoutUrl,
		breadcrumbs,
	}

	return c.Render("dashboard", props)
}

// user playground
func playgroundHandler(c *fiber.Ctx) error {
	token := github.AccessToken(c)
	pageParams := config.MakePageParams(token != "")
	return c.Render("playground", pageParams)
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

	token := github.AccessToken(c)
	pageParams := config.MakePageParams(token != "")

	props := ordersProps{
		PageParams: pageParams,
		Title:      "All Orders",
		Orders:     orders,
		Total:      result.RowsAffected,
	}

	return c.Render("orders", props)
}
