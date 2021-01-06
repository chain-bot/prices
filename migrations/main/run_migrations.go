package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mochahub/coinprice-scraper/config"
	"github.com/mochahub/coinprice-scraper/main/src/service/database"
	"go.uber.org/fx"
	"log"
)

func main() {
	fxApp := fx.New(
		fx.Provide(config.GetSecrets),
		fx.Provide(database.NewDatabase),
		fx.Invoke(database.RunMigrations),
	)
	if err := fxApp.Start(context.Background()); err != nil {
		log.Printf("ERROR STARTING APP FOR MIGRATIONS: %s", err)
	}
}
