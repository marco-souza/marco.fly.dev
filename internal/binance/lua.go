package binance

import (
	"strconv"

	"github.com/Shopify/go-lua"
)

func fetchTicker(s *lua.State) int {
	currencyPair, ok := s.ToString(1)
	if !ok {
		logger.Error("failed to get ticker", "pair", currencyPair)
		s.PushNil()
		return 1
	}

	logger.Info("fetching ticker for: " + currencyPair)

	tick, err := FetchTicker(currencyPair)
	if err != nil {
		logger.Error("error fetching ticker", "err", err)
		s.PushNil()
		return 1
	}

	logger.Info("ticker", "pair", currencyPair, "tick", tick)

	price, err := strconv.ParseFloat(tick.Price, 64)
	if err != nil {
		logger.Error("error parsing price", "err", err)
		s.PushNil()
		return 1
	}

	logger.Info("ticker", "pair", tick.Symbol, "rate", tick.Price)
	s.PushNumber(price)
	return 1
}

func fetchAccountSnapshot(s *lua.State) int {
	logger.Info("fetching account snapshot")

	accType, ok := s.ToString(1)
	if !ok {
		logger.Error("failed to get account type")
		s.PushNil()
		return 1
	}

	snapshot, err := FetchAccountSnapshot(accType)
	if err != nil {
		logger.Error("error fetching account snapshot", "err", err)
		s.PushNil()
		return 1
	}

	logger.Info("account snapshot", "type", accType, "snapshot", snapshot)
	s.PushString(snapshot.Msg)
	return 1
}

func generateWalletReport(s *lua.State) int {
	logger.Info("generating wallet report")

	report, err := GenerateWalletReport()
	if err != nil {
		logger.Error("error generating wallet report", "err", err.Error())
		s.PushNil()
		return 1
	}

	logger.Info("wallet report", "report", report)
	s.PushString(report)
	return 1
}

func PushClient(l *lua.State) {
	l.NewTable()

	l.PushString("pair_price")
	l.PushGoFunction(fetchTicker)
	l.SetTable(-3)

	l.PushString("account_snapshot")
	l.PushGoFunction(fetchAccountSnapshot)
	l.SetTable(-3)

	l.PushString("wallet_report")
	l.PushGoFunction(generateWalletReport)
	l.SetTable(-3)

	l.SetGlobal("binance")
}
