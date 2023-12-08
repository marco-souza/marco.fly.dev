package pages

import (
	"github.com/gofiber/fiber/v2"
)

func Apply(router fiber.Router) {
	router.Get("/", rootHandler).
		Get("/playground", playgroundHandler).
		Get("/orders", ordersHandler).
		Use(notFoundHandler)
}
