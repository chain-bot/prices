package main

import (
	"context"
	"go.uber.org/fx"
	"log"
)

func main() {
	// TODO: Find a Better Logging Framework
	log.Println("Scraper Config:")
	//app.StartScraper()
	app := fx.New()
	_ = app.Start(context.Background())
}
