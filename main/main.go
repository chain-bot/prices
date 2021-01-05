package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mochahub/coinprice-scraper/main/api"
	app "github.com/mochahub/coinprice-scraper/main/app"
	"go.uber.org/fx"
	"log"
)

func main() {
	// TODO: Find a Better Logging Framework
	log.Println("Scraper Config:")
	// TODO: fx.provide the code for the influx connection
	fxApp := fx.New(
		api.GetAPIProviders(),
		fx.Invoke(app.StartScrapperCron),
	)
	if err := fxApp.Start(context.Background()); err != nil {
		log.Printf("ERROR STARTING APP: %s", err)
	}
	<-fxApp.Done()
}
