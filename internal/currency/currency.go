package currency

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
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
