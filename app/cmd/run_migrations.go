package main

import (
	"context"
	"github.com/chain-bot/scraper/app/configs"
	"github.com/chain-bot/scraper/app/internal/data/psql"
	_ "github.com/joho/godotenv/autoload"
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
