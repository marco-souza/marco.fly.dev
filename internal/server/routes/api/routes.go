package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/sqlite3/v2"

	"github.com/marco-souza/marco.fly.dev/internal/config"
	"github.com/marco-souza/marco.fly.dev/internal/server/routes/middlewares"
)

var (
	conf = config.Load()
)

func Apply(router fiber.Router) {
	// https://docs.gofiber.io/api/middleware/limiter
	if conf.Env == "production" {
		router.Use(limiter.New(limiter.Config{
			Max:     conf.RateLimit,
			Storage: sqlite3.New(),
		}))
	}

	// public routes
	router.Group("/").
		Get("/resume", resumeHandler).
		Get("/menu", menuHandler)

	router.Group("/auth/github").
		Get("/", redirectGithubAuth).
		Get("/callback", callbackGithubAuth).
		Get("/refresh", logoutGithubAuth).
		Get("/logout", logoutGithubAuth)

	if conf.Env == "development" {
		router.Get("/reload", sseReloadHandler)
	}

	// private routes
	router.Use(middlewares.MustHaveToken)

	router.Group("/").
		Post("/lua", luaHandler).
		Get("/now", nowHandler).
		Get("/sse", sseHandler)

	router.Group("/cron").
		Get("/", cronsHandler).
		Post("/", createCronHandler).
		Delete("/:id", deleteCronHandler)
}
