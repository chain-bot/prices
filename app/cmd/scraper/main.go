package main

import (
	"context"
	"log"

	"github.com/chain-bot/prices/app/configs"
	"github.com/chain-bot/prices/app/internal/data/influxdb"
	"github.com/chain-bot/prices/app/internal/data/psql"
	"github.com/chain-bot/prices/app/internal/repository"
	"github.com/chain-bot/prices/app/pkg/api"
	scraper "github.com/chain-bot/prices/app/pkg/scraper"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
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
			psql.RunMigrations,
			scraper.Run,
		),
	)
	if err := fxApp.Start(context.Background()); err != nil {
		log.Printf("ERROR STARTING APP: %s", err)
	}
	<-fxApp.Done()
}
