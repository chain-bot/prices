package common

import "time"

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// ENUM
//////////////////////////////////////////////////////////////////////////////////////////////////////////
type Interval string

const (
	Minute Interval = "1m"
	Hour   Interval = "1h"
	Day    Interval = "1D"
)

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
