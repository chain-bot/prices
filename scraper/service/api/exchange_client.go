package api

import (
	"github.com/mochahub/coinprice-scraper/scraper/models"
	"time"
)

// TODO(Zahin): Get a List of Supported Symbols
type ExchangeAPIClient interface {
	GetExchangeIdentifier() string
	GetSupportedPairs() ([]*models.Symbol, error)
	GetRawMarketData() ([]*models.RawMarketData, error)
	GetAllOHLCMarketData(
		symbol models.Symbol,
		interval time.Duration,
		startTime time.Time,
		endTime time.Time,
	) ([]*models.OHLCMarketData, error)
}
