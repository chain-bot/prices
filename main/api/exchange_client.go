package api

import "time"

// TODO(Zahin): Get a List of Supported Symbols
type ExchangeAPIClient interface {
	GetExchangeIdentifier() string
	GetSupportedPairs() ([]*Symbol, error)
	GetRawMarketData() ([]*RawMarketData, error)
	GetAllOHLCMarketData(
		baseSymbol string,
		quoteSymbol string,
		interval Interval,
		startTime time.Time,
		endTime time.Time,
	) ([]*OHLCMarketData, error)
}

type RawMarketData struct {
	MarketData
	StartTime time.Time
	Value     float64
	Volume    float64
}

type Symbol struct {
	RawBase         string
	NormalizedBase  string
	RawQuote        string
	NormalizedQuote string
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
