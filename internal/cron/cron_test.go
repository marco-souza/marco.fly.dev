package cron

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCronJob(t *testing.T) {
	c := New()
	go c.Start()
	defer c.Stop()

	t.Run("can add expression", func(t *testing.T) {
		err := c.Add("* * * * *", func() {})
		if err != nil {
			t.Fatal(err)
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
		for _, cron := range c.List() {
			assert.Greater(t, cron.id, 0)

			// TODO: make this pass
			// if cron.expression == "" {
			// 	t.Fatal("cron expression invalid", cron)
			// }
		}
	})

	t.Run("do nothing if deleting invalid id", func(t *testing.T) {
		assert.Len(t, c.List(), 1)

		c.Del(-1)

		assert.Len(t, c.List(), 1)
	})

	t.Run("can delete expression", func(t *testing.T) {
		assert.Len(t, c.List(), 1)

		c.Del(1)

		assert.Len(t, c.List(), 0)
	})
}
