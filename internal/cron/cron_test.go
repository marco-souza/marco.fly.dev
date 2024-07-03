package cron_test

import (
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/cron"
	"github.com/marco-souza/marco.fly.dev/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestCronJob(t *testing.T) {
	db.Init("")
	defer db.Close()

	cron.Start()
	defer cron.Stop()

	luaScript := "print('hello lua')"
	expressions := []string{
		"0 * * * *",
		"* * * * *",
	}

	t.Run("can add expression", func(t *testing.T) {
		for _, expr := range expressions {
			err := cron.AddScript(expr, expr, luaScript)
			assert.Nil(t, err)
		}
	})

	t.Run("fail if expression is invalid", func(t *testing.T) {
		err := cron.AddScript("first", "* * * * * *", luaScript)
		assert.Contains(t, err.Error(), "6")

		err = cron.AddScript("second", "invalid", luaScript)
		assert.Contains(t, err.Error(), "invalid")

		err = cron.AddScript("third", "1ms", luaScript)
		assert.Contains(t, err.Error(), "1ms")
	})

	t.Run("can list expression", func(t *testing.T) {
		crons := cron.List()
		assert.Equal(t, len(crons), 2)

		for _, cronEntry := range crons {
			assert.Greater(t, int(cronEntry.ID), 0)
			assert.Contains(t, expressions, cronEntry.Expression)
		}
	})

	t.Run("do nothing if deleting invalid id", func(t *testing.T) {
		assert.Len(t, cron.List(), len(expressions))

		cron.Del(-1)

		assert.Len(t, cron.List(), len(expressions))
	})

	t.Run("can delete expression", func(t *testing.T) {
		size := len(expressions)
		assert.Len(t, cron.List(), size)

		cron.Del(1)

		assert.Len(t, cron.List(), size-1)
	})
}
