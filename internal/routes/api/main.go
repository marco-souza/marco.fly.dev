package api

import (
	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/entities"
)

var registers = []entities.Register{
	now,
	test,
	orders,
}

func Apply(router fiber.Router) {
	for _, register := range registers {
		register(router)
	}
}
