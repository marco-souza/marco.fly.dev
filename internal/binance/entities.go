package binance

type AccountAssets struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

type AccountData struct {
	TotalBtcAsset string          `json:"totalAssetOfBtc"`
	Balances      []AccountAssets `json:"balances"`
}

type AccountSnapshot struct {
	Type       string      `json:"type"`
	UpdateTime int         `json:"updateTime"`
	Data       AccountData `json:"data"`
}

type AccountSnapshotResponse struct {
	Code        int               `json:"code"`
	Msg         string            `json:"msg"`
	SnapshotVos []AccountSnapshot `json:"snapshotVos"`
}

type Ticker struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}
