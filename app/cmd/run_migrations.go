package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mochahub/coinprice-scraper/app/configs"
	"github.com/mochahub/coinprice-scraper/app/internal/data/psql"
	"go.uber.org/fx"
	"log"
)

func main() {
	fxApp := fx.New(
		fx.Provide(configs.GetSecrets),
		fx.Provide(psql.NewDatabase),
		fx.Invoke(psql.RunMigrations),
	)
	if err := fxApp.Start(context.Background()); err != nil {
		log.Printf("ERROR STARTING APP FOR MIGRATIONS: %s", err)
	}
}
