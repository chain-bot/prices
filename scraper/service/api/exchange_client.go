package api

import (
	"context"
	"github.com/mochahub/coinprice-scraper/scraper/models"
	"time"
)

// TODO(Zahin): Get a List of Supported Symbols
type RestExchangeAPIClient interface {
	GetExchangeIdentifier() string
	GetSupportedPairs() ([]*models.Symbol, error)
	GetAllOHLCMarketData(
		symbol models.Symbol,
		interval time.Duration,
		startTime time.Time,
		endTime time.Time,
	) ([]*models.OHLCMarketData, error)
}

type SocketExchangeAPIClient interface {
	RestExchangeAPIClient // Is embedded the best approach?
	GetOHLCMarketDataChannel(
		ctx context.Context,
		symbol models.Symbol,
		interval time.Duration,
	) (chan *models.OHLCMarketData, error)
}
