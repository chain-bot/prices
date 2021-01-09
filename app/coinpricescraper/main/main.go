package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	app "github.com/mochahub/coinprice-scraper/app/coinpricescraper/app"
	"github.com/mochahub/coinprice-scraper/app/coinpricescraper/service/api"
	"github.com/mochahub/coinprice-scraper/app/data/influxdb"
	"github.com/mochahub/coinprice-scraper/app/data/psql"
	"github.com/mochahub/coinprice-scraper/config"
	"go.uber.org/fx"
	"log"
)

func main() {
	// TODO: Find a Better Logging Framework
	fxApp := fx.New(
		api.GetAPIProviders(),
		fx.Provide(config.GetSecrets),
		fx.Provide(psql.NewDatabase),
		fx.Provide(influxdb.NewInfluxDBClient),
		fx.Invoke(
			psql.RunMigrations,
			app.InitScrapper,
		),
	)
	if err := fxApp.Start(context.Background()); err != nil {
		log.Printf("ERROR STARTING APP: %s", err)
	}
	<-fxApp.Done()
}
