package currency

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Shopify/go-lua"
)

var (
	USDBRL_EXCHANGE_URL = "https://dolarhoje.com/cotacao.txt"
	logger              = slog.With("service", "currency")
)

func FetchDolarRealExchangeValue() (float64, error) {
	logger.Info("Fetching dolar real exchange value")

	req, err := http.Get(USDBRL_EXCHANGE_URL)
	if err != nil {
		logger.Error("error creating request", "err", err)
		return -1, err
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Error("error reading response body", "err", err)
		return -1, err
	}

	// parse output
	exchangeValueString := strings.Replace(string(body), ",", ".", 1)
	exchangeRate, err := strconv.ParseFloat(exchangeValueString, 64)
	if err != nil {
		logger.Error("error parsing exchange rate", "err", err)
		return -1, err
	}

	return exchangeRate, nil
}

func fetchDolBrlWrapper(s *lua.State) int {
	logger.Info("fetching exchange rate")

	exchangeRate, err := FetchDolarRealExchangeValue()
	if err != nil {
		logger.Error("error fetching exchange rate", "err", err)
		s.PushNumber(-1)
		return 1 // number of results
	}

	logger.Info("exchange rate", "rate", exchangeRate)
	s.PushNumber(exchangeRate)
	return 1 // number of results
}

func PushClient(l *lua.State) {
	l.NewTable() // {}

	l.PushString("usd_brl")              // {}, "usd_brl"
	l.PushGoFunction(fetchDolBrlWrapper) // {}, "usd_brl", fetchDolBrlWrapper
	l.SetTable(-3)                       // {usd_brl: fetchDolBrlWrapper}

	// make it available globaly
	l.SetGlobal("currency")
}
