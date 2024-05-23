package env_test

import (
	"os"
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/env"
	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	t.Run("should return value from env", func(t *testing.T) {
		os.Setenv("TEST_ENV", "test")
		assert.Equal(t, "test", env.Env("TEST_ENV", "default"))
	})

	t.Run("should return default value", func(t *testing.T) {
		assert.Equal(t, "default", env.Env("NOT_FOUND", "default"))
	})
}
