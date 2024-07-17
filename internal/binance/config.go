package binance

import "github.com/marco-souza/marco.fly.dev/internal/env"

var (
	binanceBaseUrl   string
	binanceApiKey    string
	binanceApiSecret string

	accountSnapURL  string
	systemStatusURL string
	tickerURL       string
)

func loadEnvs() {
	binanceBaseUrl = env.Env("BINANCE_BASE_URL", "https://api.binance.com")
	binanceApiKey = env.Env("BINANCE_API_KEY", "")
	binanceApiSecret = env.Env("BINANCE_API_SECRET", "")
	// urls
	accountSnapURL = binanceBaseUrl + "/sapi/v1/accountSnapshot"
	systemStatusURL = binanceBaseUrl + "/sapi/v1/system/status"
	tickerURL = binanceBaseUrl + "/api/v3/ticker/price"
}
