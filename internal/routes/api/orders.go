package api

import (
	"fmt"
	"math/rand"

	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/models"
)

func ordersHandler(c *fiber.Ctx) error {
	db := models.Connect()
	orders := []models.Order{}
	result := db.Preload("Author").Find(&orders)
	props := fiber.Map{"Orders": orders, "Total": result.RowsAffected}

	return c.Render("partials/order-list", props, "layouts/empty")
}

type CreateOrderInput struct {
	Name string `json: "name"`
}

func createOrderHandler(c *fiber.Ctx) error {
	input := CreateOrderInput{}
	if err := c.BodyParser(&input); err != nil {
		fmt.Println("error = ", err)
		return c.SendStatus(400)
	}

	fmt.Println("input: {?:}", input)

	db := models.Connect()

	user := models.User{
		Name:     "Marco",
		Username: "marco-souza",
	}
	db.Find(&user)

	order := models.Order{
		Name: input.Name,
		CoverUrl: fmt.Sprintf(
			"https://source.unsplash.com/random/?%s&%d&w=%d",
			input.Name,
			rand.Intn(100),
			250,
		),
		Author: user,
	}
	db.Create(&order)

	orders := []models.Order{}
	result := db.Preload("Author").Find(&orders)
	props := fiber.Map{"Orders": orders, "Total": result.RowsAffected}

	return c.Render("partials/order-list", props, "layouts/empty")
}

func orders(router fiber.Router) {
	router.Get("/orders", ordersHandler)
	router.Post("/orders", createOrderHandler)
}
