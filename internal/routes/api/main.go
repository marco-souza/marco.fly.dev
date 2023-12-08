package api

import (
	"github.com/gofiber/fiber/v2"
)

func Apply(router fiber.Router) {
	router.Group("/orders").
		Get("/", ordersHandler).
		Post("/", createOrderHandler).
		Delete("/:id", deleteOrderHandler)

	router.Get("/now", nowHandler)
}
