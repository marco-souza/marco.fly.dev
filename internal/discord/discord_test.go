package discord_test

import (
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/discord"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	dc := discord.New()

	t.Run("should create a new discord service", func(t *testing.T) {
		assert.NotNil(t, dc)
	})
}
