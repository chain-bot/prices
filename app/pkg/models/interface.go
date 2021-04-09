package models

import (
	"time"
)

// TODO(Zahin): Get a List of Supported Symbols
type ExchangeAPIClient interface {
	GetExchangeIdentifier() string
	GetSupportedPairs() ([]*Symbol, error)
	GetRawMarketData() ([]*RawMarketData, error)
	GetAllOHLCVMarketData(
		symbol Symbol,
		interval time.Duration,
		startTime time.Time,
		endTime time.Time,
	) ([]*OHLCVMarketData, error)
}
