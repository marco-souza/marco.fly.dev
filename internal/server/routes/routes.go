package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/server/routes/api"
	"github.com/marco-souza/marco.fly.dev/internal/server/routes/pages"
)

func SetupRoutes(app *fiber.App) {
	api.Apply(app.Group("/api"))
	pages.Apply(app.Group("/"))
}
