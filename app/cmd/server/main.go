package main

import (
	"context"
	"log"

	"github.com/chain-bot/prices/app/configs"
	"github.com/chain-bot/prices/app/internal/data/influxdb"
	"github.com/chain-bot/prices/app/internal/data/psql"
	"github.com/chain-bot/prices/app/internal/repository"
	"github.com/chain-bot/prices/app/pkg/server"
	"github.com/chain-bot/prices/app/pkg/server/routes"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
)

func main() {
	// TODO: Find a Better Logging Framework and Pass in Via Uber fx
	fxApp := fx.New(
		fx.Provide(configs.GetSecrets),
		fx.Provide(psql.NewDatabase),
		fx.Provide(influxdb.NewInfluxDBClient),
		fx.Provide(repository.NewRepository),
		fx.Provide(routes.NewHandler),
		fx.Invoke(
			psql.RunMigrations,
			server.Run,
		),
	)
	if err := fxApp.Start(context.Background()); err != nil {
		log.Printf("ERROR STARTING Server: %s", err)
	}
	<-fxApp.Done()
}
