package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	app "github.com/mochahub/coinprice-scraper/main/app"
	"go.uber.org/fx"
	"log"
)

func main() {
	// TODO: Find a Better Logging Framework
	log.Println("Scraper Config:")
	// fx.provide the code for the influx connection
	fxApp := fx.New(fx.Invoke(app.StartScraper))
	_ = fxApp.Start(context.Background())
	_ = fxApp.Done()
}
