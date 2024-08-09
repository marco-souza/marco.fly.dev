package telegram_test

import (
	"os"
	"testing"

	"github.com/Shopify/go-lua"
	"github.com/marco-souza/marco.fly.dev/internal/di"
	"github.com/marco-souza/marco.fly.dev/internal/telegram"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	msg := "testing-msg"

	t.Run("failed to send chat message", func(t *testing.T) {
		tg := telegram.New()

		err := tg.SendChatMessage(msg)
		assert.ErrorContains(t, err, "404")
	})

	t.Run("include failing token", func(t *testing.T) {
		token := "test-token"
		os.Setenv("TELEGRAM_BOT_TOKEN", token)
		tg := telegram.New()

		err := tg.SendChatMessage(msg)

		assert.ErrorContains(t, err, "404")
		assert.ErrorContains(t, err, token)
	})

	t.Run("include failing params", func(t *testing.T) {
		tg := telegram.New()

		err := tg.SendChatMessage(msg)

		assert.ErrorContains(t, err, msg)
		assert.ErrorContains(t, err, "chat_id=") // chat
	})
}

func TestLuaPushClient(t *testing.T) {
	tg := telegram.New()
	di.Injectable(tg)

	t.Run("push client to lua", func(t *testing.T) {
		l := lua.NewState()

		lua.OpenLibraries(l)
		telegram.PushClient(l)

		err := lua.DoString(l, "print(not_existent_mod.send_message ~= nil)")
		assert.Error(t, err)

		err = lua.DoString(l, "print(telegram.send_message ~= nil)")
		assert.NoError(t, err)

	})
}
