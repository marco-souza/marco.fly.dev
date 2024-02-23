package pages

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marco-souza/marco.fly.dev/internal/routes/middlewares"
)

func Apply(router fiber.Router) {
	// private routes
	router.Group("/app").
		Use(middlewares.MustBeLoged).
		Get("/", dashboardHandler)

	router.Get("/", rootHandler).
		Get("/playground", playgroundHandler).
		Get("/orders", ordersHandler).
		Get("/login", loginHandler).
		Use(notFoundHandler)
}
