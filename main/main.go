package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mochahub/coinprice-scraper/main/api"
	app "github.com/mochahub/coinprice-scraper/main/app"
	"github.com/mochahub/coinprice-scraper/main/config"
	"github.com/mochahub/coinprice-scraper/main/database"
	"go.uber.org/fx"
	"log"
)

func main() {
	// TODO: Find a Better Logging Framework
	log.Println("Scraper Config:")
	// TODO: fx.provide the code for the influx connection
	fxApp := fx.New(
		api.GetAPIProviders(),
		fx.Provide(config.GetSecrets),
		fx.Provide(database.NewDatabase),
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
