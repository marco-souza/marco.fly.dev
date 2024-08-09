package cron_test

import (
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/config"
	"github.com/marco-souza/marco.fly.dev/internal/cron"
	"github.com/marco-souza/marco.fly.dev/internal/db"
	"github.com/marco-souza/marco.fly.dev/internal/di"
	"github.com/marco-souza/marco.fly.dev/internal/lua"
	"github.com/stretchr/testify/assert"
)

func TestCronJob(t *testing.T) {
	di.Injectables(
		config.Config{SqliteUrl: ":memory:"},
		db.New,
		lua.LuaService{},
		cron.New,
	)

	taskScheduler := di.MustInject(cron.TaskScheduleService{})

	luaScript := "print('hello lua')"
	expressions := []string{
		"0 * * * *",
		"* * * * *",
	}

	t.Run("can add expression", func(t *testing.T) {
		for _, expr := range expressions {
			err := taskScheduler.AddScript(expr, expr, luaScript)
			assert.Nil(t, err)
		}
	})

	t.Run("fail if expression is invalid", func(t *testing.T) {
		err := taskScheduler.AddScript("first", "* * * * * *", luaScript)
		assert.Contains(t, err.Error(), "6")

		err = taskScheduler.AddScript("second", "invalid", luaScript)
		assert.Contains(t, err.Error(), "invalid")

		err = taskScheduler.AddScript("third", "1ms", luaScript)
		assert.Contains(t, err.Error(), "1ms")
	})

	t.Run("can list expression", func(t *testing.T) {
		crons := taskScheduler.List()
		assert.Equal(t, len(crons), 2)

		for _, cronEntry := range crons {
			assert.Greater(t, int(cronEntry.ID), 0)
			assert.Contains(t, expressions, cronEntry.Expression)
		}
	})

	t.Run("do nothing if deleting invalid id", func(t *testing.T) {
		assert.Len(t, taskScheduler.List(), len(expressions))

		taskScheduler.Del(-1)

		assert.Len(t, taskScheduler.List(), len(expressions))
	})

	t.Run("can delete expression", func(t *testing.T) {
		size := len(expressions)
		assert.Len(t, taskScheduler.List(), size)

		taskScheduler.Del(1)

		assert.Len(t, taskScheduler.List(), size-1)
	})
}
