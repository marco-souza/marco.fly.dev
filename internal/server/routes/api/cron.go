package api

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/cron"
)

func cronsHandler(c *fiber.Ctx) error {
	return renderCronList(c)
}

type CreateCronInput struct {
	Name    string `json:"name" validate:"required,gte=0,lte=130"`
	Cron    string `json:"cron" validate:"required,gte=9,lte=130"`
	Snippet string `json:"snippet" validate:"required,gte=0"`
}

func createCronHandler(c *fiber.Ctx) error {
	input := CreateCronInput{}
	if err := c.BodyParser(&input); err != nil {
		fmt.Println("error = ", err)
		return c.SendStatus(400)
	}

	err := validate.Struct(input)
	if err != nil {
		fmt.Println("validation error = ", err)
		return c.SendStatus(400)
	}

	cron.AddScript(input.Name, input.Cron, input.Snippet)

	return renderCronList(c)
}

func deleteCronHandler(c *fiber.Ctx) error {
	id := c.Params("id", "")
	cronId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("no cron id found", id)
		return c.SendStatus(400)
	}

	cron.Del(cronId)
	return renderCronList(c)
}

func renderCronList(c *fiber.Ctx) error {
	crons := cron.List()
	props := fiber.Map{"Crons": crons, "Total": len(crons)}
	return c.Render("partials/cron-list", props, "layouts/empty")
}
