package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mochahub/coinprice-scraper/config"
	"github.com/mochahub/coinprice-scraper/data/influxdb"
	"github.com/mochahub/coinprice-scraper/data/psql"
	app "github.com/mochahub/coinprice-scraper/scraper/app"
	"github.com/mochahub/coinprice-scraper/scraper/service/api"
	"go.uber.org/fx"
	"log"
)

func main() {
	// TODO: Find a Better Logging Framework and Pass in Via Uber fx
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
