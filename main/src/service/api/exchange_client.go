package api

import (
	"github.com/mochahub/coinprice-scraper/main/src/service/api/common"
	"time"
)

// TODO(Zahin): Get a List of Supported Symbols
type ExchangeAPIClient interface {
	GetExchangeIdentifier() string
	GetSupportedPairs() ([]*common.Symbol, error)
	GetRawMarketData() ([]*common.RawMarketData, error)
	GetAllOHLCMarketData(
		baseSymbol string,
		quoteSymbol string,
		// TODO(Zahin): Makes more sense to make this a duration
		interval common.Interval,
		startTime time.Time,
		endTime time.Time,
	) ([]*common.OHLCMarketData, error)
}
