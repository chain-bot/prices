package models

import "time"

type RawMarketData struct {
	MarketData
	StartTime time.Time
	Value     float64
	Volume    float64
}

// TODO: Add "product", which represents how the exchange recognizes the base-quote pair
// Each exchange handles it slightly differently, we shouldn't try and reconstruct it using the base/quote string
type Symbol struct {
	RawBase         string
	NormalizedBase  string
	RawQuote        string
	NormalizedQuote string
}

type OHLCMarketData struct {
	MarketData
	StartTime  time.Time
	EndTime    time.Time
	OpenPrice  float64
	HighPrice  float64
	LowPrice   float64
	ClosePrice float64
	Volume     float64
}

type MarketData struct {
	Source        string
	BaseCurrency  string
	QuoteCurrency string
}
