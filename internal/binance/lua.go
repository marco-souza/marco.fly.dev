package binance

import (
	"strconv"

	"github.com/Shopify/go-lua"
	"github.com/marco-souza/marco.fly.dev/internal/di"
)

type binanceLuaFuncs struct {
	*BinanceService
}

func newBinanceLuaFuncs() *binanceLuaFuncs {
	b := di.MustInject(BinanceService{})
	return &binanceLuaFuncs{b}
}

func (blf *binanceLuaFuncs) fetchTicker(s *lua.State) int {
	currencyPair, ok := s.ToString(1)
	if !ok {
		logger.Error("failed to get ticker", "pair", currencyPair)
		s.PushNil()
		return 1
	}

	logger.Info("fetching ticker for: " + currencyPair)

	tick, err := blf.FetchTicker(currencyPair)
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

func (blf *binanceLuaFuncs) fetchAccountSnapshot(s *lua.State) int {
	logger.Info("fetching account snapshot")

	accType, ok := s.ToString(1)
	if !ok {
		logger.Error("failed to get account type")
		s.PushNil()
		return 1
	}

	snapshot, err := blf.FetchAccountSnapshot(accType)
	if err != nil {
		logger.Error("error fetching account snapshot", "err", err)
		s.PushNil()
		return 1
	}

	logger.Info("account snapshot", "type", accType, "snapshot", snapshot)
	s.PushString(snapshot.Msg)
	return 1
}

func (blf *binanceLuaFuncs) generateWalletReport(s *lua.State) int {
	logger.Info("generating wallet report")

	report, err := blf.GenerateWalletReport()
	if err != nil {
		logger.Error("error generating wallet report", "err", err)
		s.PushNil()
		return 1
	}

	logger.Info("wallet report", "report", report)
	s.PushString(report)
	return 1
}

func PushClient(l *lua.State) {
	l.NewTable()

	f := newBinanceLuaFuncs()

	l.PushString("pair_price")
	l.PushGoFunction(f.fetchTicker)
	l.SetTable(-3)

	l.PushString("account_snapshot")
	l.PushGoFunction(f.fetchAccountSnapshot)
	l.SetTable(-3)

	l.PushString("wallet_report")
	l.PushGoFunction(f.generateWalletReport)
	l.SetTable(-3)

	l.SetGlobal("binance")
}
