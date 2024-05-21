package cron_test

import (
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/cron"
	"github.com/stretchr/testify/assert"
)

func TestCronJob(t *testing.T) {
	c := cron.CronService
	go c.Start()
	defer c.Stop()

	expressions := []string{
		"0 * * * *",
		"* * * * *",
	}

	t.Run("can add expression", func(t *testing.T) {
		for _, expr := range expressions {
			err := c.Add(expr, func() {})
			assert.Nil(t, err)
		}
	})

	t.Run("fail if expression is invalid", func(t *testing.T) {
		err := c.Add("* * * * * *", func() {})
		assert.Contains(t, err.Error(), "6")

		err = c.Add("invalid", func() {})
		assert.Contains(t, err.Error(), "invalid")

		err = c.Add("1ms", func() {})
		assert.Contains(t, err.Error(), "1ms")
	})

	t.Run("can list expression", func(t *testing.T) {
		for _, entry := range c.List() {
			assert.Greater(t, entry.ID, 0)
			assert.Contains(t, expressions, entry.Expression)
		}
	})

	t.Run("do nothing if deleting invalid id", func(t *testing.T) {
		assert.Len(t, c.List(), len(expressions))

		c.Del(-1)

		assert.Len(t, c.List(), len(expressions))
	})

	t.Run("can delete expression", func(t *testing.T) {
		size := len(expressions)
		assert.Len(t, c.List(), size)

		c.Del(1)

		assert.Len(t, c.List(), size-1)
	})
}
