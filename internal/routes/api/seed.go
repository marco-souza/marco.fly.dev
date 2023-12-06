package api

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/models"
)

func seed(router fiber.Router) {
	router.Get("/orders/seed", func(c *fiber.Ctx) error {
		user := models.User{
			Name:     "Marco",
			Username: "marco-souza",
		}

		db := models.Connect()
		db.Create(&user)

		order := models.Order{
			Name: "Corre de natal",
			CoverUrl: fmt.Sprintf(
				"https://source.unsplash.com/random/?Order&%d",
				rand.Intn(100),
			),
			Author: user,
		}
		db.Create(&order)

		product := models.Product{
			Name: "Corre de natal",
			Price: sql.NullFloat64{
				Float64: 1.99,
			},
			Order: order,
		}
		db.Create(&product)

		for i := range [5]int{} {
			item := models.Item{
				Name:      fmt.Sprintf("%s %d", "cha", i),
				Product:   product,
				Quantity:  10.01,
				Recipient: "someone else",
			}
			db.Create(&item)
		}

		return c.JSON(map[string]interface{}{
			"user":    user,
			"order":   order,
			"product": product,
		})
	})
}
