package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/marco-souza/marco.fly.dev/internal/di"
	"github.com/marco-souza/marco.fly.dev/internal/lua"
)

func luaHandler(c *fiber.Ctx) error {
	snippet := c.FormValue("snippet")
	if snippet == "" {
		logger.Info("no snippet provided")
		lines := []string{"No output yet."}
		return c.Render("partials/code", fiber.Map{"Lines": lines}, "layouts/empty")
	}

	logger.Info("lua code", "snippet", snippet)
	code := ""
	err := di.Invoke(func(l *lua.LuaService) error {
		output, err := l.Run(snippet)
		code = output
		return err
	})

	if err != nil {
		logger.Error("lua error", "err", err)
		lines := []string{err.Error()}
		return c.Render("partials/code", fiber.Map{"Lines": lines}, "layouts/empty")
	}

	logger.Info("lua output", "code", code)
	lines := strings.Split(code, "\n")
	return c.Render("partials/code", fiber.Map{"Lines": lines}, "layouts/empty")
}
