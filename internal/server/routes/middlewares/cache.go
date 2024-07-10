package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

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
	logger.Info("cache", "status", isCacheDisabled)
	return isCacheDisabled
}
