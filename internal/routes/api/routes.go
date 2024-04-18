package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/sqlite3/v2"

	"github.com/marco-souza/marco.fly.dev/internal/config"
)

var conf = config.Load()

func Apply(router fiber.Router) {
	// https://docs.gofiber.io/api/middleware/limiter
	router.Use(limiter.New(limiter.Config{
		Max:     conf.RateLimit,
		Storage: sqlite3.New(),
	}))

	router.Group("/orders").
		Get("/", ordersHandler).
		Post("/", createOrderHandler).
		Delete("/:id", deleteOrderHandler)

	router.Group("/auth/github").
		Get("/", redirectGithubAuth).
		Get("/callback", callbackGithubAuth).
		Get("/refresh", logoutGithubAuth).
		Get("/logout", logoutGithubAuth)

	router.Get("/now", nowHandler)
	router.Get("/sse", sseHandler)

	if conf.Env == "development" {
		router.Get("/reload", sseReloadHandler)
	}
}
