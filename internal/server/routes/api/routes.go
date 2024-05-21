package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/sqlite3/v2"

	"github.com/marco-souza/marco.fly.dev/internal/config"
	"github.com/marco-souza/marco.fly.dev/internal/github"
)

var conf = config.Load()

// create auth check middleware
func authMiddleware(c *fiber.Ctx) error {
	token := github.AccessToken(c)
	if token == "" {
		log.Println("Unauthorized access, redirecting to login page.")
		return c.Redirect(conf.Github.LoginPage)
	}
	return c.Next()
}

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
	router.Use(authMiddleware)

	router.Group("/").
		Post("/lua", luaHandler).
		Get("/now", nowHandler).
		Get("/sse", sseHandler)

	router.Group("/orders").
		Get("/", ordersHandler).
		Post("/", createOrderHandler).
		Delete("/:id", deleteOrderHandler)

	router.Group("/cron").
		Get("/", cronsHandler).
		Post("/", createCronHandler).
		Delete("/:id", deleteCronHandler)
}
