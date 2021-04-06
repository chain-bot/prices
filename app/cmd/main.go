package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mochahub/coinprice-scraper/app/configs"
	"github.com/mochahub/coinprice-scraper/app/internal/data/influxdb"
	"github.com/mochahub/coinprice-scraper/app/internal/data/psql"
	"github.com/mochahub/coinprice-scraper/app/internal/repository"
	"github.com/mochahub/coinprice-scraper/app/internal/scraper"
	"github.com/mochahub/coinprice-scraper/app/pkg/api"
	"go.uber.org/fx"
	"log"
)

func main() {
	// TODO: Find a Better Logging Framework and Pass in Via Uber fx
	fxApp := fx.New(
		api.GetAPIProviders(),
		fx.Provide(configs.GetSecrets),
		fx.Provide(psql.NewDatabase),
		fx.Provide(influxdb.NewInfluxDBClient),
		fx.Provide(repository.NewRepository),
		fx.Invoke(
			RunMigrations,
			scraper.InitScrapper,
		),
	)
	if err := fxApp.Start(context.Background()); err != nil {
		log.Printf("ERROR STARTING APP: %s", err)
	}
	<-fxApp.Done()
}
