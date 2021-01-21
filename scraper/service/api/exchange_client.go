package api

import (
	"github.com/mochahub/coinprice-scraper/scraper/models"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/common"
	"time"
)

// TODO(Zahin): Get a List of Supported Symbols
type ExchangeAPIClient interface {
	GetExchangeIdentifier() string
	GetSupportedPairs() ([]*models.Symbol, error)
	GetRawMarketData() ([]*models.RawMarketData, error)
	GetAllOHLCMarketData(
		baseSymbol string,
		quoteSymbol string,
		// TODO(Zahin): Makes more sense to make this a duration
		interval common.Interval,
		startTime time.Time,
		endTime time.Time,
	) ([]*models.OHLCMarketData, error)
}
