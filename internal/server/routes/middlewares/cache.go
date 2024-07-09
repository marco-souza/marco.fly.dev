package middlewares

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

var logger = log.New(log.Writer(), "cache: ", log.Flags())

// cache middleware with 15 minutes expiration
var DefaultCache = NewCache(15 * time.Minute)

func NewCache(expiration time.Duration) fiber.Handler {
	return cache.New(cache.Config{
		Next:         next,
		Expiration:   expiration,
		CacheControl: true,
	})
}

func next(c *fiber.Ctx) bool {
	isCacheDisabled := c.Query("noCache") == "true"
	logger.Println("is disabled: ", isCacheDisabled)
	return isCacheDisabled
}
