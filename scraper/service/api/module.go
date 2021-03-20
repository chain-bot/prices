package api

import (
	"context"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/binance"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/coinbasepro"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/ftx"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/kucoin"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/okex"
	"go.uber.org/fx"
)

type ExchangeClientResult struct {
	fx.Out
	Client RestExchangeAPIClient `group:"exchange_client"`
}

type ExchangeClients struct {
	fx.In
	Clients []RestExchangeAPIClient `group:"exchange_client"`
}

func GetAPIProviders() fx.Option {
	return fx.Options(
		fx.Provide(NewBinanaceAPIClient),
		fx.Provide(NewCoinbaseProAPIClient),
		fx.Provide(NewKucoinAPIClient),
		fx.Provide(NewOkexAPIClient),
		fx.Provide(NewFtxAPIClient),
	)
}

func NewBinanaceAPIClient(
	ctx context.Context,
) ExchangeClientResult {
	return ExchangeClientResult{
		Client: binance.NewBinanceAPIClient(ctx),
	}
}

func NewCoinbaseProAPIClient() ExchangeClientResult {
	return ExchangeClientResult{
		Client: coinbasepro.NewCoinbaseProAPIClient(),
	}
}

func NewKucoinAPIClient() ExchangeClientResult {
	return ExchangeClientResult{
		Client: kucoin.NewKucoinAPIClient(),
	}
}

func NewOkexAPIClient() ExchangeClientResult {
	return ExchangeClientResult{
		Client: okex.NewOkexAPIClient(),
	}
}

func NewFtxAPIClient() ExchangeClientResult {
	return ExchangeClientResult{
		Client: ftx.NewFtxAPIClient(),
	}
}
