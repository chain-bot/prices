package api

import (
	"github.com/mochahub/coinprice-scraper/scraper/service/api/binance"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/coinbasepro"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/ftx"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/kucoin"
	"github.com/mochahub/coinprice-scraper/scraper/service/api/okex"
	"go.uber.org/fx"
)

type RestExchangeClientResult struct {
	fx.Out
	Client RestExchangeAPIClient `group:"rest_exchange_client"`
}

type RestExchangeClients struct {
	fx.In
	Clients []RestExchangeAPIClient `group:"rest_exchange_client"`
}

type SocketExchangeClientResult struct {
	fx.Out
	Client SocketExchangeAPIClient `group:"socket_exchange_client"`
}

type SocketExchangeClients struct {
	fx.In
	Clients []SocketExchangeAPIClient `group:"socket_exchange_client"`
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

func NewBinanaceAPIClient() RestExchangeClientResult {
	return RestExchangeClientResult{
		Client: binance.NewBinanceAPIClient(),
	}
}

func NewCoinbaseProAPIClient() RestExchangeClientResult {
	return RestExchangeClientResult{
		Client: coinbasepro.NewCoinbaseProAPIClient(),
	}
}

func NewKucoinAPIClient() RestExchangeClientResult {
	return RestExchangeClientResult{
		Client: kucoin.NewKucoinAPIClient(),
	}
}

func NewOkexAPIClient() RestExchangeClientResult {
	return RestExchangeClientResult{
		Client: okex.NewOkexAPIClient(),
	}
}

func NewFtxAPIClient() RestExchangeClientResult {
	return RestExchangeClientResult{
		Client: ftx.NewFtxAPIClient(),
	}
}

func GetSocketAPIProviders() fx.Option {
	return fx.Options(
		fx.Provide(NewBinanaceSocketAPIClient),
	)
}

func NewBinanaceSocketAPIClient() SocketExchangeClientResult {
	return SocketExchangeClientResult{
		Client: binance.NewBinanceAPIClient(),
	}
}
