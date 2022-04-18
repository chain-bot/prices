package coinbasepro

import (
	"encoding/json"
)

// CandleStickData https://docs.pro.coinbase.com/?python#get-historic-rates
type CandleStickData struct {
	CloseTime  float64
	LowPrice   float64
	HighPrice  float64
	OpenPrice  float64
	ClosePrice float64
	Volume     float64
}

func (candleStickResponse *CandleStickData) UnmarshalJSON(
	data []byte,
) (err error) {
	var responseSlice []interface{}
	if err := json.Unmarshal(data, &responseSlice); err != nil {
		return err
	}
	candleStickResponse.CloseTime = responseSlice[0].(float64)
	candleStickResponse.LowPrice = responseSlice[1].(float64)
	candleStickResponse.HighPrice = responseSlice[2].(float64)
	candleStickResponse.OpenPrice = responseSlice[3].(float64)
	candleStickResponse.ClosePrice = responseSlice[4].(float64)
	candleStickResponse.Volume = responseSlice[5].(float64)
	return nil
}

// ProductsResponse https://docs.pro.coinbase.com/?python#get-products
type ProductsResponse []struct {
	ID              string `json:"id"`
	Message         string `json:"message"` // ERROR MESSAGE
	DisplayName     string `json:"display_name"`
	BaseCurrency    string `json:"base_currency"`
	QuoteCurrency   string `json:"quote_currency"`
	BaseIncrement   string `json:"base_increment"`
	QuoteIncrement  string `json:"quote_increment"`
	BaseMinSize     string `json:"base_min_size"`
	BaseMaxSize     string `json:"base_max_size"`
	MinMarketFunds  string `json:"min_market_funds"`
	MaxMarketFunds  string `json:"max_market_funds"`
	Status          string `json:"status"`
	StatusMessage   string `json:"status_message"`
	CancelOnly      bool   `json:"cancel_only"`
	LimitOnly       bool   `json:"limit_only"`
	PostOnly        bool   `json:"post_only"`
	TradingDisabled bool   `json:"trading_disabled"`
}
