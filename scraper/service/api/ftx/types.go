package ftx

import (
	"time"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// https://docs.ftx.com/#get-historical-prices
//////////////////////////////////////////////////////////////////////////////////////////////////////////
type HistoricalPricesResponse struct {
	Success bool     `json:"success"`
	Result  []Candle `json:"result"`
}
type Candle struct {
	Close     float64   `json:"close"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Open      float64   `json:"open"`
	StartTime time.Time `json:"startTime"`
	Volume    float64   `json:"volume"`
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// https://docs.ftx.com/#get-markets
//////////////////////////////////////////////////////////////////////////////////////////////////////////
type MarketsResponse struct {
	Success bool     `json:"success"`
	Result  []Market `json:"result"`
}

type Market struct {
	Name          string `json:"name"`
	BaseCurrency  string `json:"baseCurrency"`
	QuoteCurrency string `json:"quoteCurrency"`
	Type          string `json:"type"`
	Underlying    string `json:"underlying"`
}
