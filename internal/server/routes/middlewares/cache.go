package middlewares

import (
	"log"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

var logger = log.New(log.Writer(), "cache: ", log.Flags())

var Cache = cache.New(cache.Config{
	Next: func(c *fiber.Ctx) bool {
		if cacheControl, ok := c.GetReqHeaders()["Cache-Control"]; !ok {
			if slices.Contains(cacheControl, "no-cache") {
				logger.Println("disabled by cache control", cacheControl)
				return false
			}
		}

		if c.Query("noCache") == "true" {
			logger.Println("disabled by query param")
			return false
		}

		// skip middleware if cache enable
		return true
	},
	Expiration:   15 * time.Minute,
	CacheControl: true,
})
