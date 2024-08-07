package binance_test

import (
	"os"
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/binance"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	os.Setenv("BINANCE_API_SECRET", "api-secret")
	os.Setenv("BINANCE_API_KEY", "api-key")
	os.Setenv("BINANCE_BASE_URL", "https://base-url.api.com")

	b := binance.New()
	b.Start()

	t.Run("failed to fetch ticket", func(t *testing.T) {
		output, err := b.FetchTicker("BTCUSDT")
		assert.ErrorContains(t, err, "base-url")
		assert.Nil(t, output)
	})

	t.Run("failed to fetch account snapshot", func(t *testing.T) {
		output, err := b.FetchAccountSnapshot("spot")
		assert.ErrorContains(t, err, "base-url")
		assert.Nil(t, output)
	})
}
