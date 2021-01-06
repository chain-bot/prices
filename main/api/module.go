package api

import (
	"github.com/mochahub/coinprice-scraper/main/api/binance"
	"github.com/mochahub/coinprice-scraper/main/config"
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
	)
}

func NewBinanaceAPIClient(secrets *config.Secrets) ExchangeClientResult {
	return ExchangeClientResult{
		Client: binance.NewBinanceAPIClient(secrets.BinanceApiKey),
	}
}

//func NewServer(p ServerParams) *Server {
//	server := newServer()
//	for _, h := range p.Handlers {
//		server.Register(h)
//	}
//	return server
//}
