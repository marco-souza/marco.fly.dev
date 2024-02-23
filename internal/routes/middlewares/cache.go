package middlewares

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

var Cache = cache.New(cache.Config{
	Next: func(c *fiber.Ctx) bool {
		isCacheDisabled := c.Query("noCache") == "true"
		log.Println("Is cache disabled: ", isCacheDisabled)
		return isCacheDisabled
	},
	Expiration:   15 * time.Minute,
	CacheControl: true,
})
