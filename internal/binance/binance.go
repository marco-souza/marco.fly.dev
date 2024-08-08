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

	"github.com/marco-souza/marco.fly.dev/internal/currency"
)

var logger = slog.With("service", "binance")

type BinanceService struct {
}

func New() *BinanceService {
	return &BinanceService{}
}

func (b *BinanceService) Start() error {
	loadEnvs()
	return nil
}

func (b *BinanceService) Stop() error {
	return nil
}

func (b *BinanceService) FetchAccountSnapshot(walletType string) (*AccountSnapshotResponse, error) {
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

var pairMap = map[string]string{
	"BRL": "USDTBRL",
	"BTC": "BTCUSDT",
	"ETH": "ETHUSDT",
	"SOL": "SOLUSDT",
	"BNB": "BNBUSDT",
}

func (b *BinanceService) GenerateWalletReport() (string, error) {
	// fetch account snapshot
	snapshot, err := b.FetchAccountSnapshot("SPOT")
	if err != nil {
		logger.Error("error fetching account snapshot", "err", err)
		return "", err
	}

	// generate report
	latestSnapshot := snapshot.SnapshotVos[len(snapshot.SnapshotVos)-1]
	// for _, snapshot := range snapshot.SnapshotVos {
	date := time.Unix(0, int64(latestSnapshot.UpdateTime)*int64(time.Millisecond))
	total := latestSnapshot.Data.TotalBtcAsset

	totalFloat, err := strconv.ParseFloat(total, 64)
	if err != nil {
		logger.Error("error parsing total", "err", err)
		return "", err
	}

	usdRateTick, err := b.FetchTicker("BTCUSDT")
	if err != nil {
		logger.Error("error fetching asset ticker ", "err", err)
		return "", err
	}

	usdRateFloat, err := strconv.ParseFloat(usdRateTick.Price, 64)
	if err != nil {
		logger.Error("error parsing price", "err", err)
		return "", err
	}

	brlUsdRate, err := currency.FetchDolarRealExchangeValue()
	if err != nil {
		logger.Error("error fetching exchange rate", "err", err)
		return "", err
	}

	formatedDate := date.Format("02.01.2006")
	totalUsd := totalFloat * usdRateFloat
	report := fmt.Sprintf(
		"*Wallet Report - %s*\n\n*Total*: `$%.2f x %.2f = R$%.2f `\n---\n",
		formatedDate,
		totalUsd,
		brlUsdRate,
		totalUsd*brlUsdRate,
	)

	for _, balance := range latestSnapshot.Data.Balances {
		if balance.Free == "0" || balance.Asset == "ETHW" {
			continue
		}

		price := 1.0
		if balance.Asset[:3] != "USD" {
			pair, ok := pairMap[balance.Asset]
			if !ok {
				logger.Error("error fetching asset ticker ", "err", "pair not found", "pair", balance.Asset)
				return "", fmt.Errorf("pair '%s' not found", balance.Asset)
			}

			t, err := b.FetchTicker(pair)
			if err != nil {
				logger.Error("error fetching asset ticker ", "err", err, "pair", "BTC"+balance.Asset)
				return "", err
			}

			price, err = strconv.ParseFloat(t.Price, 64)
			if err != nil {
				logger.Error("error parsing price", "err", err)
				return "", err
			}
		}

		free, err := strconv.ParseFloat(balance.Free, 64)
		if err != nil {
			logger.Error("error parsing free", "err", err)
			return "", err
		}

		if balance.Asset == "BRL" {
			// fix pair order (BRL to USD)
			price = 1 / price
		}

		total := free * price
		report += fmt.Sprintf(
			"- %s %s: `$ %07.4f`\n",
			balance.Asset[:3],
			upOrDown(balance.Asset, total),
			total,
		)

	}

	return report, nil
}

var latestTotal = map[string]float64{}

func upOrDown(asset string, current float64) string {
	previous, ok := latestTotal[asset]
	if !ok {
		logger.Warn("missing previous value")
	}

	// persist latest values
	latestTotal[asset] = current

	if previous == 0 {
		return "🔵 (+0.000%)"
	}

	percentage := (current - previous) / previous * 100
	if percentage > 0 {
		return fmt.Sprintf("🔥 (%+4.3f%%)", percentage)
	}

	if percentage < 0 {
		return fmt.Sprintf("⛈️  (%+4.3f%%)", percentage)
	}

	return fmt.Sprintf("🔵 (%+4.3f%%)", percentage)
}

func (b *BinanceService) FetchTicker(currencyPair string) (*Ticker, error) {
	// API ref: https://binance-docs.github.io/apidocs/spot/en/#symbol-price-ticker
	req, err := http.NewRequest("GET", tickerURL, nil)
	if err != nil {
		logger.Error("error creating request", "err", err)
		return nil, err
	}

	logger.Info("fetching ticker", "pair", currencyPair)

	// set params
	params := url.Values{}
	params.Add("symbol", currencyPair)

	req.URL.RawQuery = params.Encode()

	logger.Info("fetching ticker", "url", req.URL.String())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("error fetching ticker", "err", err)
		return nil, err
	}

	logger.Info("ticker response", "status", res.Status)

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("error reading response body", "err", err, "res", res)
		return nil, err
	}

	logger.Info("ticker response", "body", string(responseBody))

	var ticker Ticker
	if err := json.Unmarshal(responseBody, &ticker); err != nil {
		logger.Error("error parsing ticker", "err", err, "body", string(responseBody))
		return nil, err
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
