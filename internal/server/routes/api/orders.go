package api

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

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
	Name string `json:"name" validate:"required,gte=0,lte=130"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func createOrderHandler(c *fiber.Ctx) error {
	input := CreateOrderInput{}
	if err := c.BodyParser(&input); err != nil {
		fmt.Println("error = ", err)
		return c.SendStatus(400)
	}

	err := validate.Struct(input)
	if err != nil {
		fmt.Println("validation error = ", err)
		return c.SendStatus(400)
	}

	fmt.Println("input test: {?:}", input)

	db := models.Connect()

	user := models.User{
		Name:     "Marco",
		Username: "marco-souza",
	}
	db.Find(&user)

	// generate random img based on order text
	coverUrl := fmt.Sprintf(
		"https://source.unsplash.com/random/256x128?%s&%d",
		strings.ReplaceAll(input.Name, " ", ","),
		rand.Intn(100),
	)
	order := models.Order{
		Name:     input.Name,
		CoverUrl: coverUrl,
		Author:   user,
	}
	db.Create(&order)

	return renderOrderList(db, c)
}

func deleteOrderHandler(c *fiber.Ctx) error {
	orderId := c.Params("id", "")
	if orderId == "" {
		fmt.Println("no order id found")
		return c.SendStatus(400)
	}

	db := models.Connect()
	db.Delete(&models.Order{}, orderId)

	return renderOrderList(db, c)
}

func renderOrderList(db *gorm.DB, c *fiber.Ctx) error {
	orders := []models.Order{}
	result := db.Preload("Author").Find(&orders)
	props := fiber.Map{"Orders": orders, "Total": result.RowsAffected}

	return c.Render("partials/order-list", props, "layouts/empty")

}
