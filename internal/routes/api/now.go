package api

import (
	"log"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/gx/views"
)

func now(router fiber.Router) {
	router.Get("/now",
		func(c *fiber.Ctx) error {
			now := time.Now().Format(time.RFC1123Z)
			log.Println("Now is", now)

			comp := views.Now(now)
			handler := templ.Handler(comp)

			return handler.Component.Render(
				c.Context(),
				c.Response().BodyWriter(),
			)
		},
	)
}
