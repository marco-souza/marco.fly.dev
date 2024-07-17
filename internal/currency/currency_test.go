package currency_test

import (
	"testing"

	"github.com/marco-souza/marco.fly.dev/internal/currency"
	"github.com/stretchr/testify/assert"
)

func TestCurrencyExchange(t *testing.T) {
	output, err := currency.FetchDolarRealExchangeValue()
	assert.Nil(t, err)
	assert.Greater(t, output, 0.0)
}
