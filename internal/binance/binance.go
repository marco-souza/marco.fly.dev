package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var logger = slog.With("service", "binance")

func Start() {
	logger.Info("starting binance service")
	loadEnvs()
}

func Stop() {
	logger.Info("stopping binance service")
}

func FetchAccountSnapshot(walletType string) (*AccountSnapshotResponse, error) {
	req, err := http.NewRequest("GET", accountSnapURL, nil)
	if err != nil {
		logger.Error("error creating request", "err", err)
		return nil, err
	}

	// set params
	params := req.URL.Query()
	params.Add("type", walletType)
	params.Add("endTime", fmt.Sprint(time.Now().Unix()*1000))

	signedParams := signParams(params)
	req.URL.RawQuery = signedParams.Encode()

	// set headers
	req.Header.Set("X-MBX-APIKEY", binanceApiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("error fetching account snapshot", "err", err)
		return nil, err
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("error fetching account snapshot", "err", err)
		return nil, err
	}

	var accountSnapshot AccountSnapshotResponse
	if err := json.Unmarshal(responseBody, &accountSnapshot); err != nil {
		logger.Error("error fetching account snapshot", "err", err)
		return nil, err
	}

	return &accountSnapshot, nil
}

func FetchTicker(currencyPair string) (*Ticker, error) {
	// API ref: https://binance-docs.github.io/apidocs/spot/en/#symbol-price-ticker
	req, err := http.NewRequest("GET", tickerURL, nil)
	if err != nil {
		logger.Error("error creating request", "err", err)
		return nil, err
	}

	// set params
	params := url.Values{}
	params.Add("symbol", currencyPair)

	req.URL.RawQuery = params.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("error fetching ticker", "err", err)
		return nil, err
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("error reading response body", "err", err)
		return nil, err
	}

	var ticker Ticker
	if err := json.Unmarshal(responseBody, &ticker); err != nil {
		panic(err)
	}

	return &ticker, nil
}

func signParams(params url.Values) url.Values {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	params.Add("timestamp", timestamp)

	signature := sign(params.Encode())
	params.Add("signature", signature)

	return params
}

func sign(text string) string {
	hash := hmac.New(sha256.New, []byte(binanceApiSecret))
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
