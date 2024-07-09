package cache_test

import (
	"testing"
	"time"

	"github.com/marco-souza/marco.fly.dev/internal/cache"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	c := cache.NewCache()
	if c == nil {
		t.Error("Cache should not be nil")
	}

	expectedValue := []byte("value")

	t.Run("Get a missing item", func(t *testing.T) {
		_, err := c.Get("key")

		assert.ErrorContains(t, err, "cache miss for key: key")
	})

	t.Run("Set a new item", func(t *testing.T) {
		err := c.Set("key", expectedValue, nil)
		assert.NoError(t, err)
	})

	t.Run("Get an existing item", func(t *testing.T) {
		value, err := c.Get("key")

		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	})

	t.Run("Set a new item with ttl", func(t *testing.T) {
		err := c.Set("with_ttl", expectedValue, cache.WithTTL(1))
		assert.NoError(t, err)

		// assert cache hits value
		value, err := c.Get("with_ttl")
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, value)

		time.Sleep(1001 * time.Millisecond)
		_, err = c.Get("with_ttl")
		assert.ErrorContains(t, err, "cache miss for key: with_ttl")

	})
}
