package api

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/marco-souza/marco.fly.dev/internal/lua"
)

func luaHandler(c *fiber.Ctx) error {
	snippet := c.FormValue("snippet")
	if snippet == "" {
		log.Println("No snippet provided.")
		lines := []string{"No output yet."}
		return c.Render("partials/code", fiber.Map{"Lines": lines}, "layouts/empty")
	}

	log.Println("Lua code:", snippet)
	code, err := lua.Run(snippet)
	if err != nil {
		log.Println("Lua error:", err)
		lines := []string{err.Error()}
		return c.Render("partials/code", fiber.Map{"Lines": lines}, "layouts/empty")
	}

	log.Println("Lua output:", code)
	lines := strings.Split(code, "\n")
	return c.Render("partials/code", fiber.Map{"Lines": lines}, "layouts/empty")
}
