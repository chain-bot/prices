package api

import (
	"github.com/chain-bot/scraper/app/pkg/api/binance"
	"github.com/chain-bot/scraper/app/pkg/api/coinbasepro"
	"github.com/chain-bot/scraper/app/pkg/api/ftx"
	"github.com/chain-bot/scraper/app/pkg/api/kucoin"
	"github.com/chain-bot/scraper/app/pkg/api/okex"
	"github.com/chain-bot/scraper/app/pkg/models"
	"go.uber.org/fx"
)

type ExchangeClientResult struct {
	fx.Out
	Client models.ExchangeAPIClient `group:"exchange_client"`
}

type ExchangeClients struct {
	fx.In
	Clients []models.ExchangeAPIClient `group:"exchange_client"`
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
