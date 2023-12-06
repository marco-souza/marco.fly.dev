package pages

import (
	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/entities"
)

var registers = []entities.Register{
	root, playground, orders,
	notFound,
}

func Apply(router fiber.Router) {
	for _, register := range registers {
		register(router)
	}
}
