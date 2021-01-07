package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mochahub/coinprice-scraper/config"
	app "github.com/mochahub/coinprice-scraper/main/app"
	"github.com/mochahub/coinprice-scraper/main/src/service/api"
	"github.com/mochahub/coinprice-scraper/main/src/service/database"
	"github.com/mochahub/coinprice-scraper/main/src/service/influxdb"
	"go.uber.org/fx"
	"log"
)

func main() {
	// TODO: Find a Better Logging Framework
	fxApp := fx.New(
		api.GetAPIProviders(),
		fx.Provide(config.GetSecrets),
		fx.Provide(database.NewDatabase),
		fx.Provide(influxdb.NewInfluxDBClient),
		fx.Invoke(
			database.RunMigrations,
			app.InitScrapper,
		),
	)
	if err := fxApp.Start(context.Background()); err != nil {
		log.Printf("ERROR STARTING APP: %s", err)
	}
	<-fxApp.Done()
}
