package pages

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/marco-souza/marco.fly.dev/internal/routes/middlewares"
)

func Apply(router fiber.Router) {
	router.
		// add caching
		Use(cache.New(cache.Config{
			Next: func(c *fiber.Ctx) bool {
				isCacheDisabled := c.Query("noCache") == "true"
				log.Println("Is cache disabled: ", isCacheDisabled)
				return isCacheDisabled
			},
			Expiration:   15 * time.Minute,
			CacheControl: true,
		}))

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
