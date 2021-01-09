package api

import (
	"github.com/mochahub/coinprice-scraper/app/coinpricescraper/service/api/binance"
	"github.com/mochahub/coinprice-scraper/app/coinpricescraper/service/api/coinbasepro"
	"github.com/mochahub/coinprice-scraper/config"
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
		// Coinbase PR
		fx.Provide(NewCoinbaseProAPIClient),
	)
}

func NewBinanaceAPIClient(secrets *config.Secrets) ExchangeClientResult {
	return ExchangeClientResult{
		Client: binance.NewBinanceAPIClient(secrets.BinanceApiKey),
	}
}

func NewCoinbaseProAPIClient(secrets *config.Secrets) ExchangeClientResult {
	return ExchangeClientResult{
		Client: coinbasepro.NewCoinbaseProAPIClient(
			secrets.CoinbaseProApiKey, secrets.CoinbaseProApiSecret, secrets.CoinbaseProApiPassphrase),
	}
}
