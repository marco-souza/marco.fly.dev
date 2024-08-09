package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/marco-souza/marco.fly.dev/internal/cron"
	"github.com/marco-souza/marco.fly.dev/internal/di"
)

func cronsHandler(c *fiber.Ctx) error {
	return renderCronList(c)
}

func createCronHandler(c *fiber.Ctx) error {
	input := CreateCronInput{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	err := input.Validate()
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	err = di.Invoke(func(cron cron.TaskScheduleService) {
		cron.AddScript(input.Name, input.Cron, input.Snippet)
	})

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return renderCronList(c)
}

func deleteCronHandler(c *fiber.Ctx) error {
	id := c.Params("id", "")
	cronId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).SendString("invalid id")
	}

	err = di.Invoke(func(cron cron.TaskScheduleService) {
		cron.Del(cronId)
	})

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return renderCronList(c)
}

func renderCronList(c *fiber.Ctx) error {
	crons := []cron.Cron{}
	err := di.Invoke(func(cron cron.TaskScheduleService) {
		crons = append(crons, cron.List()...)
	})

	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	props := fiber.Map{"Crons": crons, "Total": len(crons)}
	return c.Render("partials/cron-list", props, "layouts/empty")
}
