package api

import "time"

// TODO(Zahin): Get a List of Supported Symbols
type ExchangeAPIClient interface {
	GetRawMarketData() ([]*RawMarketData, error)
	GetOHLCMarketData() ([]*OHLCMarketData, error)
}

type RawMarketData struct {
	MarketData
	StartTime time.Time
	Value     float64
	Volume    float64
}

type OHLCMarketData struct {
	MarketData
	StartTime time.Time
	EndTime   time.Time
	HighPrice float64
	LowPrice  float64
	Volume    float64
}

type MarketData struct {
	Source        string
	BaseCurrency  string
	QuoteCurrency string
}
