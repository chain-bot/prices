package api

import (
	"github.com/mochahub/coinprice-scraper/scraper/service/api/binance"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/coinbasepro"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/ftx"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/kucoin"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/okex"
	"go.uber.org/fx"
)

type ExchangeClientResult struct {
	fx.Out
	Client ExchangeAPIClient `group:"exchange_client"`
}

type ExchangeClients struct {
	fx.In
	Clients []ExchangeAPIClient `group:"exchange_client"`
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

func NewBinanaceAPIClient() ExchangeClientResult {
	return ExchangeClientResult{
		Client: binance.NewBinanceAPIClient(),
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
