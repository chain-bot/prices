package api

import (
	"github.com/mochahub/coinprice-scraper/main/api/binance"
	"go.uber.org/fx"
)

func GetAPIProviders() fx.Option {
	return fx.Options(
		fx.Provide(binance.NewBinanceAPIClient),
	)
}
