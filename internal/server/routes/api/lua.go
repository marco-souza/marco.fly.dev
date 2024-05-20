package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/marco-souza/marco.fly.dev/internal/lua"
)

func luaHandler(c *fiber.Ctx) error {
	snippet := c.FormValue("snippet")
	if snippet == "" {
		log.Println("No snippet provided.")
		return c.Render("partials/code", fiber.Map{"Code": "No output yet."}, "layouts/empty")
	}

	log.Println("Lua code:", snippet)
	code, err := lua.Runtime.Run(snippet)
	if err != nil {
		log.Println("Lua error:", err)
		return c.Render("partials/code", fiber.Map{"Code": err.Error()}, "layouts/empty")
	}

	log.Println("Lua output:", code)
	return c.Render("partials/code", fiber.Map{"Code": code}, "layouts/empty")
}
