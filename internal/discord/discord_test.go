package discord_test

import (
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/discord"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("should create a new discord service", func(t *testing.T) {
		service := discord.DiscordService
		assert.NotNil(t, service)
	})
}
