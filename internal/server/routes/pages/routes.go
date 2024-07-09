package pages

import (
	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/server/routes/middlewares"
)

func Apply(router fiber.Router) {
	// add caching
	router.Use(middlewares.DefaultCache)

	// private routes
	router.Group("/app").
		Use(middlewares.MustBeLoged).
		Get("/", dashboardHandler).
		Get("/playground", playgroundHandler).
		Get("/cronjobs", cronHandler)

	router.Get("/", rootHandler).
		Get("/resume", resumeHandler).
		Get("/login", loginHandler).
		Get("/blog", blogHandler).
		Use(notFoundHandler)
}
