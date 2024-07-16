package telegram_test

import (
	"os"
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/telegram"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	msg := "testing-msg"

	t.Run("failed to send chat message", func(t *testing.T) {
		telegram.Start()

		err := telegram.SendChatMessage(msg)
		assert.ErrorContains(t, err, "404")
	})

	t.Run("include failing token", func(t *testing.T) {
		token := "test-token"
		os.Setenv("TELEGRAM_BOT_TOKEN", token)

		telegram.Start()

		err := telegram.SendChatMessage(msg)

		assert.ErrorContains(t, err, "404")
		assert.ErrorContains(t, err, token)
	})

	t.Run("include failing params", func(t *testing.T) {
		telegram.Start()

		err := telegram.SendChatMessage(msg)

		assert.ErrorContains(t, err, msg)
		assert.ErrorContains(t, err, "chat_id=") // chat
	})
}
