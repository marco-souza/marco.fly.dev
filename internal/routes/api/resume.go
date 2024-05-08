package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marco-souza/marco.fly.dev/internal/github"
)

func resumeHandler(c *fiber.Ctx) error {
	resumeContent, err := github.Resume(cfg.ResumeURL)

	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Resume not found")
	}

	return c.Send(resumeContent)
}
